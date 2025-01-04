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
	// @TODO Allow 헤더등 표준에 약간 아쉬운게 있어서 직접 구현하기
	r.Use(middleware.CorsHandler(middleware.CorsOptions{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			constants.HeaderContentType,
			constants.HeaderAccessControlAllowCredentials,
			constants.HeaderXForwardedFor,
		},
		MaxAge: 300,
	}))
	r.Use(middleware.HTTPRecoverer)

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

	api := chi.NewRouter().With(middleware.HTTPContentType)
	handler.NewChatHandler().ApplyRoutes(api)
	r.Mount("/api", api)

	srv.Run(port)
}
