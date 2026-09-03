package pedido

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"nueces-backend/internal/producto"
)

// Error carries the HTTP status and exact message a handler should respond
// with, so the service can own validation and transaction-failure wording
// while the handler stays free of business logic.
type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func badRequest(message string) *Error {
	return &Error{Status: http.StatusBadRequest, Message: message}
}

func notFound(message string) *Error {
	return &Error{Status: http.StatusNotFound, Message: message}
}

func internalError(message string) *Error {
	return &Error{Status: http.StatusInternalServerError, Message: message}
}

// estadosValidos son los únicos valores aceptados para pedidos.estado. Las
// transiciones son libres a propósito: el admin puede corregir errores
// poniendo cualquier estado en cualquier momento, sin máquina de estados.
var estadosValidos = map[string]bool{
	"pendiente":      true,
	"confirmado":     true,
	"en_preparacion": true,
	"entregado":      true,
}

type Service struct {
	db           *sql.DB
	repo         *Repository
	productoRepo *producto.Repository
}

func NewService(db *sql.DB, repo *Repository, productoRepo *producto.Repository) *Service {
	return &Service{db: db, repo: repo, productoRepo: productoRepo}
}

func (s *Service) List() ([]Pedido, error) {
	return s.repo.List()
}

func (s *Service) UpdateEstado(id int, estado string) (Pedido, error) {
	if !estadosValidos[estado] {
		return Pedido{}, badRequest("estado debe ser uno de: pendiente, confirmado, en_preparacion, entregado")
	}

	pedido, err := s.repo.UpdateEstado(id, estado)
	if err != nil {
		if err == sql.ErrNoRows {
			return Pedido{}, notFound("pedido no encontrado")
		}
		return Pedido{}, internalError("error actualizando estado del pedido")
	}
	return pedido, nil
}

func (s *Service) Create(input PedidoInput) (Pedido, error) {
	if input.ClienteNombre == "" {
		return Pedido{}, badRequest("cliente_nombre es requerido")
	}
	if input.ClienteContacto == "" {
		return Pedido{}, badRequest("cliente_contacto es requerido")
	}
	if input.MetodoPago != "transferencia" && input.MetodoPago != "efectivo" {
		return Pedido{}, badRequest("metodo_pago debe ser 'transferencia' o 'efectivo'")
	}
	if len(input.Items) == 0 {
		return Pedido{}, badRequest("el pedido debe tener al menos un item")
	}
	for _, item := range input.Items {
		if item.CantidadKg <= 0 {
			return Pedido{}, badRequest("cantidad_kg debe ser mayor a 0 en todos los items")
		}
	}

	tx, err := s.db.Begin()
	if err != nil {
		return Pedido{}, internalError("error iniciando transaccion")
	}
	defer tx.Rollback()

	locks := make(map[int]producto.Producto)
	var faltantes []string

	for _, item := range input.Items {
		if _, ok := locks[item.ProductoID]; ok {
			continue
		}
		pl, err := s.productoRepo.LockForUpdate(tx, item.ProductoID)
		if err != nil {
			if err == sql.ErrNoRows {
				return Pedido{}, badRequest(fmt.Sprintf("producto id %d no existe", item.ProductoID))
			}
			return Pedido{}, internalError("error consultando stock")
		}
		locks[item.ProductoID] = pl
	}

	pedidoPorProducto := make(map[int]float64)
	for _, item := range input.Items {
		pedidoPorProducto[item.ProductoID] += item.CantidadKg
	}
	for productoID, cantidadPedida := range pedidoPorProducto {
		pl := locks[productoID]
		if cantidadPedida > pl.StockKg {
			faltantes = append(faltantes, fmt.Sprintf("%s (pedido: %.2fkg, disponible: %.2fkg)", pl.Nombre, cantidadPedida, pl.StockKg))
		}
	}
	if len(faltantes) > 0 {
		return Pedido{}, badRequest("stock insuficiente para: " + strings.Join(faltantes, ", "))
	}

	for productoID, cantidadPedida := range pedidoPorProducto {
		if err := s.productoRepo.DecrementStock(tx, productoID, cantidadPedida); err != nil {
			return Pedido{}, internalError("error actualizando stock")
		}
	}

	pedido, err := s.repo.InsertPedido(tx, input.ClienteNombre, input.ClienteContacto, input.MetodoPago)
	if err != nil {
		return Pedido{}, internalError("error creando pedido")
	}

	for _, item := range input.Items {
		pl := locks[item.ProductoID]
		itemID, err := s.repo.InsertItem(tx, pedido.ID, item.ProductoID, item.CantidadKg, pl.PrecioPorKg)
		if err != nil {
			return Pedido{}, internalError("error creando items del pedido")
		}
		pedido.Items = append(pedido.Items, PedidoItem{
			ID:             itemID,
			PedidoID:       pedido.ID,
			ProductoID:     item.ProductoID,
			ProductoNombre: pl.Nombre,
			CantidadKg:     item.CantidadKg,
			PrecioUnitario: pl.PrecioPorKg,
		})
	}

	if err := tx.Commit(); err != nil {
		return Pedido{}, internalError("error confirmando pedido")
	}

	return pedido, nil
}
