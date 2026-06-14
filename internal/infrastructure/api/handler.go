package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/michaelcardoza/dummy-smtp/internal/core/mail"
)

type MailService interface {
	List(ctx context.Context) ([]*mail.Message, error)
	Get(ctx context.Context, id string) (*mail.Message, error)
	DeleteByID(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error
}

type Handler struct {
	mailService MailService
}

func NewHandler(mailService MailService) *Handler {
	return &Handler{mailService: mailService}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/messages", h.list)
	mux.HandleFunc("DELETE /api/v1/messages", h.deleteAll)
	mux.HandleFunc("GET /api/v1/messages/{id}", h.get)
	mux.HandleFunc("DELETE /api/v1/messages/{id}", h.delete)
	return mux
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	messages, err := h.mailService.List(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to list messages")
		return
	}
	h.writeJSON(w, http.StatusOK, messages)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	message, err := h.mailService.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, mail.ErrNotFound) {
			h.writeError(w, http.StatusNotFound, "message not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "failed to get message")
		return
	}
	h.writeJSON(w, http.StatusOK, message)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.mailService.DeleteByID(r.Context(), id); err != nil {
		if errors.Is(err, mail.ErrNotFound) {
			h.writeError(w, http.StatusNotFound, "message not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "failed to delete message")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteAll(w http.ResponseWriter, r *http.Request) {
	if err := h.mailService.DeleteAll(r.Context()); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to delete messages")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *Handler) writeError(w http.ResponseWriter, status int, msg string) {
	h.writeJSON(w, status, map[string]string{"error": msg})
}
