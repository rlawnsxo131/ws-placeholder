package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type chatHandler struct{}

func NewChatHandler() Handler {
	return &chatHandler{}
}

func (h *chatHandler) ApplyRoutes(r chi.Router) {
	r.Route("/room", func(r chi.Router) {
		r.Post("/", h.postRoom())
		r.Delete("/{roomId}", h.deleteRoom())
	})
	r.Route("/rooms", func(r chi.Router) {
		r.Get("/", h.getRooms())
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

func (h *chatHandler) getRooms() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 3800 ~ 5000 / 1sec
		// res := make([]string, 100)

		// for i := 0; i < 100; i++ {
		// 	res[i] = "roomList"
		// }
		// w.WriteHeader(http.StatusOK)
		// json.NewEncoder(w).Encode(res)

		// 3800 ~ 7700 / 1sec
		t := make(chan []string, 100)
		res := make([]string, 100)
		go func() {
			for i := 0; i < 100; i++ {
				res[i] = "roomList"
			}
		}()
		t <- res

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(<-t)
	}
}
