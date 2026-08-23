package producto

type Producto struct {
	ID          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Categoria   string  `json:"categoria"`
	PrecioPorKg float64 `json:"precio_por_kg"`
	StockKg     float64 `json:"stock_kg"`
}
