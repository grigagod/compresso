package http

import (
	"github.com/gorilla/mux"
	"github.com/grigagod/compresso/internal/auth"
)

func MapAuthRoutes(router *mux.Router, h auth.Handlers) {
	router.HandleFunc("/register", h.Register()).Methods("POST")
	router.HandleFunc("/login", h.Login()).Methods("POST")
}
