package auth

import (
	"net/http"

	"nueces-backend/internal/httpx"
)

// RequireAuth envuelve un handler para que solo se ejecute si el request
// trae una cookie de sesión válida; si no, responde 401. Se aplica solo a
// las rutas de admin (ej: listar pedidos), nunca a las rutas públicas.
func RequireAuth(service *Service) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)
			if err != nil || !service.Validate(cookie.Value) {
				httpx.WriteError(w, http.StatusUnauthorized, "no autenticado")
				return
			}
			next(w, r)
		}
	}
}
