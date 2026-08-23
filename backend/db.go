package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const defaultDatabaseURL = "postgres://postgres:postgres@localhost:5433/nueces?sslmode=disable"

const schema = `
CREATE TABLE IF NOT EXISTS productos (
	id SERIAL PRIMARY KEY,
	nombre TEXT NOT NULL,
	descripcion TEXT,
	categoria TEXT,
	precio_por_kg NUMERIC(10,2) NOT NULL,
	stock_kg NUMERIC(10,2) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS pedidos (
	id SERIAL PRIMARY KEY,
	cliente_nombre TEXT NOT NULL,
	cliente_contacto TEXT NOT NULL,
	fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
	estado TEXT NOT NULL DEFAULT 'pendiente'
);

CREATE TABLE IF NOT EXISTS pedido_items (
	id SERIAL PRIMARY KEY,
	pedido_id INT NOT NULL REFERENCES pedidos(id),
	producto_id INT NOT NULL REFERENCES productos(id),
	cantidad_kg NUMERIC(10,2) NOT NULL,
	precio_unitario NUMERIC(10,2) NOT NULL
);
`

// connectDB abre la conexion a Postgres y crea el esquema si no existe.
// No crea la base de datos: se asume que ya existe (ver README/instrucciones).
func connectDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = defaultDatabaseURL
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("abriendo conexion a la base de datos: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("conectando a la base de datos: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("creando esquema: %w", err)
	}

	return db, nil
}
