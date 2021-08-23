package http

import "github.com/gorilla/mux"

func MapAuthRoutes(router *mux.Router, h *authHandlers) {
	router.HandleFunc("/register", h.Register())
	router.HandleFunc("/login", h.Login())
}
