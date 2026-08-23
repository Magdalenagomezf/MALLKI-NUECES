package main

import (
	"log"
	"net/http"

	"nueces-backend/internal/httpx"
	"nueces-backend/internal/pedido"
	"nueces-backend/internal/producto"
)

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatalf("no se pudo inicializar la base de datos: %v", err)
	}
	defer db.Close()

	productoRepo := producto.NewRepository(db)
	productoService := producto.NewService(productoRepo)
	productoHandler := producto.NewHandler(productoService)

	pedidoRepo := pedido.NewRepository(db)
	pedidoService := pedido.NewService(db, pedidoRepo, productoRepo)
	pedidoHandler := pedido.NewHandler(pedidoService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handleHealth)

	mux.HandleFunc("GET /api/productos", productoHandler.List)
	mux.HandleFunc("GET /api/productos/{id}", productoHandler.Get)
	mux.HandleFunc("POST /api/productos", productoHandler.Create)
	mux.HandleFunc("PUT /api/productos/{id}", productoHandler.Update)
	mux.HandleFunc("DELETE /api/productos/{id}", productoHandler.Delete)

	mux.HandleFunc("POST /api/pedidos", pedidoHandler.Create)
	mux.HandleFunc("GET /api/pedidos", pedidoHandler.List)

	log.Println("servidor escuchando en :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("error en el servidor: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
