package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/grigagod/compresso/internal/auth"
)

func MapAuthRoutes(router *chi.Mux, h auth.Handlers) {
	router.Method("POST", "/register", h.Register())
	router.Method("POST", "/login", h.Login())
}
