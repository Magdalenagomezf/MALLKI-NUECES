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
	mux.HandleFunc("GET /api/me", requireAuth(authHandler.Me))

	// Los productos (nombre, descripción, categoría, foto, precio) son fijos y
	// se cargan directamente en la base: el catálogo los lee sin login (GET
	// público), pero crearlos/editarlos/borrarlos es una acción de admin.
	mux.HandleFunc("GET /api/productos", productoHandler.List)
	mux.HandleFunc("GET /api/productos/{id}", productoHandler.Get)
	mux.HandleFunc("POST /api/productos", requireAuth(productoHandler.Create))
	mux.HandleFunc("PUT /api/productos/{id}", requireAuth(productoHandler.Update))
	mux.HandleFunc("DELETE /api/productos/{id}", requireAuth(productoHandler.Delete))
	// PATCH /api/productos/{id}/stock es lo único que el admin puede editar de
	// un producto desde el panel: el resto de los campos son fijos.
	mux.HandleFunc("PATCH /api/productos/{id}/stock", requireAuth(productoHandler.UpdateStock))

	// POST /api/pedidos (armar pedido) es público a propósito: lo usan los
	// clientes desde /armar-pedido. GET /api/pedidos (listar) es solo para
	// el admin panel, protegido con la sesión de login.
	mux.HandleFunc("POST /api/pedidos", pedidoHandler.Create)
	mux.HandleFunc("GET /api/pedidos", requireAuth(pedidoHandler.List))
	// PATCH /api/pedidos/{id}/estado también es solo para el admin panel:
	// cambiar el estado de un pedido es una acción de gestión, no algo que
	// haga el cliente.
	mux.HandleFunc("PATCH /api/pedidos/{id}/estado", requireAuth(pedidoHandler.UpdateEstado))

	handler := httpx.CORS(frontendOrigin)(mux)

	log.Println("servidor escuchando en :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("error en el servidor: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
