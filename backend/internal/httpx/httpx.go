package httpx

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// WriteJSON writes body as a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

// WriteError writes a JSON error response with the given status code.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

// IsNoRows reports whether err is sql.ErrNoRows.
func IsNoRows(err error) bool {
	return err == sql.ErrNoRows
}
