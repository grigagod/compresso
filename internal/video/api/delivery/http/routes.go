package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/grigagod/compresso/internal/video/api"
)

func MapVideoRoutes(router chi.Router, h api.Handlers) {
	router.Route("/videos", func(r chi.Router) {
		r.Method("GET", "/", h.GetVideos())
		r.Method("GET", "/{id}", h.GetVideoByID())
		r.Method("POST", "/", h.CreateVideo())
	})

	router.Route("/tickets", func(r chi.Router) {
		r.Method("GET", "/", h.GetTickets())
		r.Method("GET", "/{id}", h.GetTicketByID())
		r.Method("POST", "/", h.CreateTicket())
	})
}
