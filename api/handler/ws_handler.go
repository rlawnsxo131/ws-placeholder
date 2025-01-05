package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type wsHandler struct{}

func NewWSHandler() Handler {
	return &wsHandler{}
}

func (h *wsHandler) ApplyRoutes(r chi.Router) {
	r.Get("/echo", h.getEcho())
	r.Get("/chat", h.getChat())
}

func (h *wsHandler) getEcho() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("/ws/echo"))
	}
}

func (h *wsHandler) getChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("/ws/chat"))
	}
}
