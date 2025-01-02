package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/rlawnsxo131/ws-placeholder/api/server"
	"github.com/rlawnsxo131/ws-placeholder/pkg/middleware"
)

func Run(port string) {
	srv := server.New()
	r := srv.Router()

	r.Use(middleware.HTTPRequestID)
	r.Use(middleware.HTTPXRequestID)
	r.Use(middleware.HTTPRealIP)
	r.Use(middleware.HTTPCompress(5))
	r.Use(middleware.HTTPLogger(middleware.DefaultHTTPServeLogger))
	r.Use(middleware.HTTPTimeout(time.Second * 3))
	r.Use(middleware.HTTPRecoverer)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	r.With(
		middleware.CORS,
	).Route("/ws", func(r chi.Router) {

		r.Get("/echo", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("/ws/echo"))
		})
		r.Get("/chat", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("/ws/chat"))
		})

	})

	r.With(
		middleware.CORS,
		middleware.HTTPContentType,
	).Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			r.Route("/chat", func(r chi.Router) {
				r.Route("/room", func(r chi.Router) {

					r.Post("/", func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte("roomId"))
					})
					r.Delete("/{roomId}", func(w http.ResponseWriter, r *http.Request) {
						roomId := chi.URLParam(r, "roomId")
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(roomId))
					})
					r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte("roomList"))
					})

				})
			})

		})
	})

	srv.Run(port)
}
