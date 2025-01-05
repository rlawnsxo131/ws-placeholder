package handler

import "github.com/go-chi/chi/v5"

type internalHandler struct{}

func NewInternalHander() Handler {
	return &internalHandler{}
}

func (h *internalHandler) ApplyRoutes(r chi.Router) {
	r.Route("/chat", func(r chi.Router) {})
}
