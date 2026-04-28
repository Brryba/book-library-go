package author

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"book-library-go/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes(r chi.Router) {
	r.Get("/", h.getAll)
	r.Post("/", h.create)
	r.Get("/{id}", h.getByID)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	authors, err := h.service.GetAll(r.Context())
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	middleware.WriteJSON(w, http.StatusOK, authors)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, errors.New("invalid id"))
		return
	}
	a, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			middleware.WriteError(w, http.StatusNotFound, err)
			return
		}
		middleware.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	middleware.WriteJSON(w, http.StatusOK, a)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	a, err := h.service.Create(r.Context(), req)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, err)
		return
	}
	middleware.WriteJSON(w, http.StatusCreated, a)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, errors.New("invalid id"))
		return
	}
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	a, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			middleware.WriteError(w, http.StatusNotFound, err)
			return
		}
		middleware.WriteError(w, http.StatusBadRequest, err)
		return
	}
	middleware.WriteJSON(w, http.StatusOK, a)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, errors.New("invalid id"))
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			middleware.WriteError(w, http.StatusNotFound, err)
			return
		}
		middleware.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
