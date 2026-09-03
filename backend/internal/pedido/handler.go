package pedido

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"nueces-backend/internal/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input PedidoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "body invalido")
		return
	}

	created, err := h.service.Create(input)
	if err != nil {
		var serviceErr *Error
		if errors.As(err, &serviceErr) {
			httpx.WriteError(w, serviceErr.Status, serviceErr.Message)
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "error creando pedido")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, created)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pedidos, err := h.service.List()
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "error consultando pedidos")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, pedidos)
}

func (h *Handler) UpdateEstado(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	var input struct {
		Estado string `json:"estado"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "body invalido")
		return
	}

	updated, err := h.service.UpdateEstado(id, input.Estado)
	if err != nil {
		var serviceErr *Error
		if errors.As(err, &serviceErr) {
			httpx.WriteError(w, serviceErr.Status, serviceErr.Message)
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "error actualizando estado del pedido")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, updated)
}
