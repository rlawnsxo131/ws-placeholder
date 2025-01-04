package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type chatHandler struct{}

func NewChatHandler() *chatHandler {
	return &chatHandler{}
}

func (h *chatHandler) ApplyRoutes(r chi.Router) {
	r.Route("/chat", func(r chi.Router) {
		r.Route("/room", func(r chi.Router) {
			r.Post("/", h.postRoom())
			r.Delete("/{roomId}", h.deleteRoom())
			r.Get("/list", h.getRoomList())
		})
	})
}

func (h *chatHandler) postRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("roomId"))
	}
}

func (h *chatHandler) deleteRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomId := chi.URLParam(r, "roomId")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(roomId))
	}
}

func (h *chatHandler) getRoomList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("roomList"))
	}
}
