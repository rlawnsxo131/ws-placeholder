package handler

import "github.com/go-chi/chi/v5"

type Handler interface {
	ApplyRoutes(r chi.Router)
}
