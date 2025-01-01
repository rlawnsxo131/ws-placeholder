package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"

	"github.com/rlawnsxo131/ws-placeholder/api/server"
	"github.com/rlawnsxo131/ws-placeholder/pkg/middleware"
)

func Run(port string) {
	srv := server.New()

	r := srv.Router()
	r.Use(chi_middleware.RequestID)
	r.Use(middleware.HTTPXRequestID)
	r.Use(chi_middleware.RealIP)
	r.Use(chi_middleware.Compress(5))
	r.Use(middleware.HTTPLogger(middleware.DefaultHTTPServeLogger))
	r.Use(middleware.HTTPTimeout(time.Second * 3))
	r.Use(chi_middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	r.Route("/ws", func(r chi.Router) {
		r.Get("/echo", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("/ws/echo"))
		})

		r.Get("/chat", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("/ws/chat"))
		})
	})

	r.Route("/api", func(r chi.Router) {

		r.Route("/v1", func(r chi.Router) {
			r.Use(middleware.HTTPContentType)

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
				})

			})
		})

	})

	srv.Run(port)
}
