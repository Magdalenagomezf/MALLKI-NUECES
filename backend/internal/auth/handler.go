package auth

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	"nueces-backend/internal/httpx"
)

// sessionCookieName es el nombre de la cookie httpOnly que guarda el token
// de sesión.
const sessionCookieName = "session_token"

type Handler struct {
	service      *Service
	secureCookie bool
}

// NewHandler crea el handler de auth. secureCookie controla el flag Secure
// de la cookie de sesión: false para desarrollo local sobre http, true en
// producción sobre https.
func NewHandler(service *Service, secureCookie bool) *Handler {
	return &Handler{service: service, secureCookie: secureCookie}
}

type loginInput struct {
	Usuario    string `json:"usuario"`
	Contrasena string `json:"contrasena"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input loginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "body invalido")
		return
	}

	token, err := h.service.Login(input.Usuario, input.Contrasena, clientIP(r))
	if err != nil {
		if errors.Is(err, ErrDemasiadosIntentos) {
			httpx.WriteError(w, http.StatusTooManyRequests, err.Error())
			return
		}
		httpx.WriteError(w, http.StatusUnauthorized, "usuario o contraseña incorrectos")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(sessionTTL),
		HttpOnly: true,
		Secure:   h.secureCookie,
		SameSite: http.SameSiteLaxMode,
	})

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// clientIP extrae la IP del cliente desde r.RemoteAddr (viene como
// "host:port"). No hay proxy reverso delante del backend (ver
// docker-compose.yml: el puerto 8080 se expone directo), así que no se lee
// X-Forwarded-For. Si en el futuro se agrega un proxy delante, revisar esto.
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookieName); err == nil {
		h.service.Logout(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.secureCookie,
		SameSite: http.SameSiteLaxMode,
	})

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Me devuelve 200 si el request trae una sesión válida. La validación de la
// cookie la hace el middleware requireAuth al envolver este handler (mismo
// patrón que GET /api/pedidos), así que acá no hace falta leerla.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
