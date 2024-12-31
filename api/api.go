package api

import (
	"net/http"

	chi_middleware "github.com/go-chi/chi/v5/middleware"

	"github.com/rlawnsxo131/ws-placeholder/api/server"
	"github.com/rlawnsxo131/ws-placeholder/pkg/middleware"
)

func Run() {
	srv := server.New("8080")

	r := srv.Router()
	r.Use(chi_middleware.RequestID)
	r.Use(middleware.HTTPXRequestID)
	r.Use(chi_middleware.RealIP)
	r.Use(middleware.HTTPLogger(middleware.DefaultHTTPServeLogger))
	r.Use(chi_middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	srv.Run()
}
