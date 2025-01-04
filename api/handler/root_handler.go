package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type rootHandler struct{}

func NewRootHandler() *rootHandler {
	return &rootHandler{}
}

func (h *rootHandler) ApplyRoutes(r chi.Router) {
	r.Get("/ping", h.getPing())
}

func (h *rootHandler) getPing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}
}
