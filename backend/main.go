package main

import (
	"log"
	"net/http"
	"os"

	"nueces-backend/internal/auth"
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

	adminUser := os.Getenv("ADMIN_USER")
	adminPasswordHash := os.Getenv("ADMIN_PASSWORD_HASH")
	if adminUser == "" || adminPasswordHash == "" {
		log.Fatal("ADMIN_USER y ADMIN_PASSWORD_HASH son requeridos")
	}

	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:3000"
	}

	// COOKIE_SECURE debe ser "true" en producción (https). En desarrollo
	// local sobre http se deja en false (default).
	secureCookie := os.Getenv("COOKIE_SECURE") == "true"

	authService := auth.NewService(adminUser, adminPasswordHash)
	authHandler := auth.NewHandler(authService, secureCookie)
	requireAuth := auth.RequireAuth(authService)

	productoRepo := producto.NewRepository(db)
	productoService := producto.NewService(productoRepo)
	productoHandler := producto.NewHandler(productoService)

	pedidoRepo := pedido.NewRepository(db)
	pedidoService := pedido.NewService(db, pedidoRepo, productoRepo)
	pedidoHandler := pedido.NewHandler(pedidoService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handleHealth)

	mux.HandleFunc("POST /api/login", authHandler.Login)
	mux.HandleFunc("POST /api/logout", authHandler.Logout)

	mux.HandleFunc("GET /api/productos", productoHandler.List)
	mux.HandleFunc("GET /api/productos/{id}", productoHandler.Get)
	mux.HandleFunc("POST /api/productos", productoHandler.Create)
	mux.HandleFunc("PUT /api/productos/{id}", productoHandler.Update)
	mux.HandleFunc("DELETE /api/productos/{id}", productoHandler.Delete)

	// POST /api/pedidos (armar pedido) es público a propósito: lo usan los
	// clientes desde /armar-pedido. GET /api/pedidos (listar) es solo para
	// el admin panel, protegido con la sesión de login.
	mux.HandleFunc("POST /api/pedidos", pedidoHandler.Create)
	mux.HandleFunc("GET /api/pedidos", requireAuth(pedidoHandler.List))

	handler := httpx.CORS(frontendOrigin)(mux)

	log.Println("servidor escuchando en :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("error en el servidor: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
