package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/rlawnsxo131/ws-placeholder/api/handler"
	"github.com/rlawnsxo131/ws-placeholder/api/server"
	"github.com/rlawnsxo131/ws-placeholder/pkg/constants"
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
	r.Use(middleware.HTTPCors(middleware.HTTPCorsConfig{
		AllowOrigins: []string{"https://*", "http://*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{
			constants.HeaderContentType,
			constants.HeaderAccept,
			constants.HeaderAuthorization,
			constants.HeaderXRequestID,
			constants.HeaderXForwardedFor,
		},
		AllowCredentials: true,
		// MaxAge:           300,
	}))
	r.Use(middleware.HTTPRecovery)

	// chi default handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	})

	// handlers
	root := chi.NewRouter()
	handler.NewRootHandler().ApplyRoutes(root)
	r.Mount("/", root)

	ws := chi.NewRouter()
	handler.NewWSHandler().ApplyRoutes(ws)
	r.Mount("/ws", ws)

	api := chi.NewRouter().With(middleware.HTTPContentType(middleware.HeaderJson))
	r.Mount("/api", api)

	chat := chi.NewRouter()
	handler.NewChatHandler().ApplyRoutes(chat)
	api.Mount("/chat", chat)

	srv.Run(port)
}
