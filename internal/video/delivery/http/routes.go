package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/grigagod/compresso/internal/video"
)

func MapVideoRoutes(router chi.Router, h video.Handlers) {
	router.Method("POST", "/", h.UploadVideo())
	router.Method("POST", "/tickets", h.CreateTicket())
}
