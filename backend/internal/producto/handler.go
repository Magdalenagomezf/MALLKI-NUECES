package producto

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

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	productos, err := h.service.List()
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "error consultando productos")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, productos)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	p, err := h.service.Get(id)
	if err != nil {
		if httpx.IsNoRows(err) {
			httpx.WriteError(w, http.StatusNotFound, "producto no encontrado")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "error consultando producto")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, p)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var p Producto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "body invalido")
		return
	}

	created, err := h.service.Create(p)
	if err != nil {
		if isValidationError(err) {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "error creando producto")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, created)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	var p Producto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "body invalido")
		return
	}

	updated, err := h.service.Update(id, p)
	if err != nil {
		if isValidationError(err) {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, "producto no encontrado")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "error actualizando producto")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, updated)
}

func (h *Handler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	var input struct {
		StockKg float64 `json:"stock_kg"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "body invalido")
		return
	}

	updated, err := h.service.UpdateStock(id, input.StockKg)
	if err != nil {
		if isValidationError(err) {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, "producto no encontrado")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "error actualizando stock del producto")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, updated)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, "producto no encontrado")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "error eliminando producto")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func isValidationError(err error) bool {
	return errors.Is(err, ErrNombreRequerido) ||
		errors.Is(err, ErrPrecioInvalido) ||
		errors.Is(err, ErrStockNegativo)
}
