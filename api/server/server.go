package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/rlawnsxo131/ws-placeholder/pkg/lib/logger"
)

type server struct {
	r   *chi.Mux
	srv *http.Server
}

func New() *server {
	r := chi.NewRouter()
	return &server{
		r: r,
		srv: &http.Server{
			Handler:      r,
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
			IdleTimeout:  time.Second * 10,
		},
	}
}

func (s *server) Router() *chi.Mux {
	return s.r
}

func (s *server) Run(port string) {
	s.srv.Addr = "0.0.0.0" + ":" + port
	logger.Default().Info().Msgf("set Addr %s", s.srv.Addr)

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, shutdownCtxCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCtxCancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Default().Fatal().Msg("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		if err := s.srv.Shutdown(shutdownCtx); err != nil {
			logger.Default().Fatal().Err(err)
		}
		serverStopCtx()
	}()

	logger.Default().Info().Msgf("start server at %s", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Default().Fatal().Err(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	logger.Default().Info().Msg("Server gracefully stopped")
}
