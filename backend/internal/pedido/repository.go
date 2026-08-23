package pedido

import "database/sql"

// Repository is the only place with SQL for the pedidos and pedido_items
// tables.
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// InsertPedido creates the pedidos row inside an existing transaction and
// returns it with its generated id, fecha_creacion and estado.
func (r *Repository) InsertPedido(tx *sql.Tx, clienteNombre, clienteContacto string) (Pedido, error) {
	var p Pedido
	err := tx.QueryRow(
		`INSERT INTO pedidos (cliente_nombre, cliente_contacto) VALUES ($1, $2)
		 RETURNING id, cliente_nombre, cliente_contacto, fecha_creacion, estado`,
		clienteNombre, clienteContacto,
	).Scan(&p.ID, &p.ClienteNombre, &p.ClienteContacto, &p.FechaCreacion, &p.Estado)
	return p, err
}

// InsertItem creates a pedido_items row inside an existing transaction and
// returns its generated id.
func (r *Repository) InsertItem(tx *sql.Tx, pedidoID, productoID int, cantidadKg, precioUnitario float64) (int, error) {
	var id int
	err := tx.QueryRow(
		`INSERT INTO pedido_items (pedido_id, producto_id, cantidad_kg, precio_unitario)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		pedidoID, productoID, cantidadKg, precioUnitario,
	).Scan(&id)
	return id, err
}

// List reads all pedidos together with their items, including the producto
// name for display. This is a read-model concern (presenting pedido data),
// not a write concern, so it legitimately joins against productos here.
func (r *Repository) List() ([]Pedido, error) {
	rows, err := r.db.Query(`SELECT id, cliente_nombre, cliente_contacto, fecha_creacion, estado FROM pedidos ORDER BY id`)
	if err != nil {
		return nil, err
	}

	pedidos := []Pedido{}
	pedidosByID := make(map[int]*Pedido)
	var order []int
	for rows.Next() {
		var p Pedido
		if err := rows.Scan(&p.ID, &p.ClienteNombre, &p.ClienteContacto, &p.FechaCreacion, &p.Estado); err != nil {
			rows.Close()
			return nil, err
		}
		p.Items = []PedidoItem{}
		pedidosByID[p.ID] = &p
		order = append(order, p.ID)
	}
	rows.Close()

	if len(order) > 0 {
		itemRows, err := r.db.Query(`
			SELECT pi.id, pi.pedido_id, pi.producto_id, pr.nombre, pi.cantidad_kg, pi.precio_unitario
			FROM pedido_items pi
			JOIN productos pr ON pr.id = pi.producto_id
			ORDER BY pi.id`)
		if err != nil {
			return nil, err
		}
		defer itemRows.Close()

		for itemRows.Next() {
			var item PedidoItem
			if err := itemRows.Scan(&item.ID, &item.PedidoID, &item.ProductoID, &item.ProductoNombre, &item.CantidadKg, &item.PrecioUnitario); err != nil {
				return nil, err
			}
			if p, ok := pedidosByID[item.PedidoID]; ok {
				p.Items = append(p.Items, item)
			}
		}
	}

	for _, id := range order {
		pedidos = append(pedidos, *pedidosByID[id])
	}

	return pedidos, nil
}
