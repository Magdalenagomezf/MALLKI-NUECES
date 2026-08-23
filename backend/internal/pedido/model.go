package pedido

import "time"

type Pedido struct {
	ID              int          `json:"id"`
	ClienteNombre   string       `json:"cliente_nombre"`
	ClienteContacto string       `json:"cliente_contacto"`
	FechaCreacion   time.Time    `json:"fecha_creacion"`
	Estado          string       `json:"estado"`
	Items           []PedidoItem `json:"items"`
}

type PedidoItem struct {
	ID             int     `json:"id"`
	PedidoID       int     `json:"pedido_id"`
	ProductoID     int     `json:"producto_id"`
	ProductoNombre string  `json:"producto_nombre,omitempty"`
	CantidadKg     float64 `json:"cantidad_kg"`
	PrecioUnitario float64 `json:"precio_unitario"`
}

type PedidoItemInput struct {
	ProductoID int     `json:"producto_id"`
	CantidadKg float64 `json:"cantidad_kg"`
}

type PedidoInput struct {
	ClienteNombre   string            `json:"cliente_nombre"`
	ClienteContacto string            `json:"cliente_contacto"`
	Items           []PedidoItemInput `json:"items"`
}
