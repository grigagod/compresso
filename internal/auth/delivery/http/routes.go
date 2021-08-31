package http

import (
	"github.com/gorilla/mux"
	"github.com/grigagod/compresso/internal/auth"
)

func MapAuthRoutes(router *mux.Router, h auth.Handlers) {
	router.Handle("/register", h.Register()).Methods("POST")
	router.Handle("/login", h.Login()).Methods("POST")
}
