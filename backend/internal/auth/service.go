package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ErrCredencialesInvalidas se devuelve cuando el usuario o la contraseña no
// coinciden con la credencial de admin configurada por variables de entorno.
var ErrCredencialesInvalidas = errors.New("usuario o contraseña incorrectos")

// ErrDemasiadosIntentos se devuelve cuando una IP superó el límite de
// intentos fallidos de login y todavía está bloqueada.
var ErrDemasiadosIntentos = errors.New("demasiados intentos, probá de nuevo en unos minutos")

// sessionTTL es cuanto dura una sesión activa desde el login.
const sessionTTL = 24 * time.Hour

// maxIntentosFallidos es la cantidad de logins fallidos que tolera una IP
// antes de quedar bloqueada. loginLockoutTTL es cuanto dura ese bloqueo.
const (
	maxIntentosFallidos = 5
	loginLockoutTTL     = 15 * time.Minute
)

// loginAttempts lleva la cuenta de intentos fallidos de una IP y, si superó
// el límite, hasta cuándo queda bloqueada.
type loginAttempts struct {
	count       int
	lockedUntil time.Time
}

// Service valida las credenciales del único usuario admin (configurado por
// env vars) y mantiene las sesiones activas en memoria. No hay tabla de
// usuarios: no hace falta persistir nada en la base de datos.
type Service struct {
	adminUser         string
	adminPasswordHash []byte

	mu       sync.Mutex
	sessions map[string]time.Time      // token -> expiración
	attempts map[string]*loginAttempts // IP -> intentos fallidos
}

// NewService crea el servicio de auth a partir de la credencial de admin.
// adminPasswordHash es un hash bcrypt (nunca la contraseña en texto plano).
func NewService(adminUser, adminPasswordHash string) *Service {
	return &Service{
		adminUser:         adminUser,
		adminPasswordHash: []byte(adminPasswordHash),
		sessions:          make(map[string]time.Time),
		attempts:          make(map[string]*loginAttempts),
	}
}

// Login valida usuario y contraseña contra la credencial configurada y, si
// son correctos, crea una sesión nueva y devuelve su token. ip es la IP del
// cliente, usada para el bloqueo por fuerza bruta.
func (s *Service) Login(user, password, ip string) (string, error) {
	if s.estaBloqueada(ip) {
		return "", ErrDemasiadosIntentos
	}

	if user == "" || password == "" || user != s.adminUser {
		s.registrarIntentoFallido(ip)
		return "", ErrCredencialesInvalidas
	}
	if err := bcrypt.CompareHashAndPassword(s.adminPasswordHash, []byte(password)); err != nil {
		s.registrarIntentoFallido(ip)
		return "", ErrCredencialesInvalidas
	}

	token, err := generateToken()
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	s.sessions[token] = time.Now().Add(sessionTTL)
	delete(s.attempts, ip) // login exitoso: se resetea el contador de esa IP
	s.mu.Unlock()

	return token, nil
}

// estaBloqueada reporta si la IP superó el límite de intentos fallidos y
// todavía sigue dentro de la ventana de bloqueo. Si el bloqueo ya venció, lo
// limpia (mismo patrón check-on-access que usa Validate con las sesiones).
func (s *Service) estaBloqueada(ip string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	a, ok := s.attempts[ip]
	if !ok || a.lockedUntil.IsZero() {
		return false
	}
	if time.Now().After(a.lockedUntil) {
		delete(s.attempts, ip)
		return false
	}
	return true
}

// registrarIntentoFallido suma un intento fallido para la IP y, si llega al
// límite, la bloquea por loginLockoutTTL.
func (s *Service) registrarIntentoFallido(ip string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	a, ok := s.attempts[ip]
	if !ok {
		a = &loginAttempts{}
		s.attempts[ip] = a
	}
	a.count++
	if a.count >= maxIntentosFallidos {
		a.lockedUntil = time.Now().Add(loginLockoutTTL)
	}
}

// Logout invalida la sesión asociada al token, si existe.
func (s *Service) Logout(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

// Validate reporta si el token corresponde a una sesión activa y no vencida.
func (s *Service) Validate(token string) bool {
	if token == "" {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	expiracion, ok := s.sessions[token]
	if !ok {
		return false
	}
	if time.Now().After(expiracion) {
		delete(s.sessions, token)
		return false
	}
	return true
}

// generateToken genera un token de sesión opaco y aleatorio.
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
