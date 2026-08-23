package producto

import "database/sql"

// Repository is the only place with SQL for the productos table.
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List() ([]Producto, error) {
	rows, err := r.db.Query(`SELECT id, nombre, descripcion, categoria, precio_por_kg, stock_kg FROM productos ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productos := []Producto{}
	for rows.Next() {
		var p Producto
		var descripcion, categoria sql.NullString
		if err := rows.Scan(&p.ID, &p.Nombre, &descripcion, &categoria, &p.PrecioPorKg, &p.StockKg); err != nil {
			return nil, err
		}
		p.Descripcion = descripcion.String
		p.Categoria = categoria.String
		productos = append(productos, p)
	}

	return productos, nil
}

func (r *Repository) Get(id int) (Producto, error) {
	var p Producto
	var descripcion, categoria sql.NullString
	row := r.db.QueryRow(`SELECT id, nombre, descripcion, categoria, precio_por_kg, stock_kg FROM productos WHERE id = $1`, id)
	err := row.Scan(&p.ID, &p.Nombre, &descripcion, &categoria, &p.PrecioPorKg, &p.StockKg)
	p.Descripcion = descripcion.String
	p.Categoria = categoria.String
	return p, err
}

func (r *Repository) Create(p Producto) (Producto, error) {
	err := r.db.QueryRow(
		`INSERT INTO productos (nombre, descripcion, categoria, precio_por_kg, stock_kg)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		p.Nombre, p.Descripcion, p.Categoria, p.PrecioPorKg, p.StockKg,
	).Scan(&p.ID)
	return p, err
}

func (r *Repository) Update(id int, p Producto) (bool, error) {
	result, err := r.db.Exec(
		`UPDATE productos SET nombre = $1, descripcion = $2, categoria = $3, precio_por_kg = $4, stock_kg = $5 WHERE id = $6`,
		p.Nombre, p.Descripcion, p.Categoria, p.PrecioPorKg, p.StockKg, id,
	)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}

func (r *Repository) Delete(id int) (bool, error) {
	result, err := r.db.Exec(`DELETE FROM productos WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}

// producto's stock. Kept here (rather than in pedido's repository) so this
// file remains the only place with SQL for the productos table, even though
// it is called from within pedido's order-creation transaction.

// LockForUpdate reads a producto's name, stock and price while holding a
// row-level lock (SELECT ... FOR UPDATE), for use inside an existing
// transaction such as pedido's order-creation flow.
func (r *Repository) LockForUpdate(tx *sql.Tx, id int) (Producto, error) {
	var p Producto
	p.ID = id
	row := tx.QueryRow(`SELECT nombre, stock_kg, precio_por_kg FROM productos WHERE id = $1 FOR UPDATE`, id)
	err := row.Scan(&p.Nombre, &p.StockKg, &p.PrecioPorKg)
	return p, err
}

// DecrementStock reduces stock_kg by cantidad for the given producto, for use
// inside an existing transaction such as pedido's order-creation flow.
func (r *Repository) DecrementStock(tx *sql.Tx, id int, cantidad float64) error {
	_, err := tx.Exec(`UPDATE productos SET stock_kg = stock_kg - $1 WHERE id = $2`, cantidad, id)
	return err
}
