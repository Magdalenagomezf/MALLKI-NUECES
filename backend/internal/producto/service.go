package producto

import (
	"database/sql"
	"errors"
)

var (
	ErrNombreRequerido = errors.New("nombre es requerido")
	ErrPrecioInvalido  = errors.New("precio_por_kg debe ser mayor a 0")
	ErrStockNegativo   = errors.New("stock_kg no puede ser negativo")
	ErrNotFound        = errors.New("producto no encontrado")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func validate(p Producto) error {
	if p.Nombre == "" {
		return ErrNombreRequerido
	}
	if p.PrecioPorKg <= 0 {
		return ErrPrecioInvalido
	}
	if p.StockKg < 0 {
		return ErrStockNegativo
	}
	return nil
}

func (s *Service) List() ([]Producto, error) {
	return s.repo.List()
}

func (s *Service) Get(id int) (Producto, error) {
	return s.repo.Get(id)
}

func (s *Service) Create(p Producto) (Producto, error) {
	if err := validate(p); err != nil {
		return Producto{}, err
	}
	return s.repo.Create(p)
}

func (s *Service) Update(id int, p Producto) (Producto, error) {
	if err := validate(p); err != nil {
		return Producto{}, err
	}
	ok, err := s.repo.Update(id, p)
	if err != nil {
		return Producto{}, err
	}
	if !ok {
		return Producto{}, ErrNotFound
	}
	p.ID = id
	return p, nil
}

func (s *Service) UpdateStock(id int, stockKg float64) (Producto, error) {
	if stockKg < 0 {
		return Producto{}, ErrStockNegativo
	}
	p, err := s.repo.UpdateStock(id, stockKg)
	if err != nil {
		if err == sql.ErrNoRows {
			return Producto{}, ErrNotFound
		}
		return Producto{}, err
	}
	return p, nil
}

func (s *Service) Delete(id int) error {
	ok, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}
