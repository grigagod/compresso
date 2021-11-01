package video

import "net/http"

type Handlers interface {
	CreateVideo() http.Handler
	CreateTicket() http.Handler
	GetVideoByID() http.Handler
	GetTicketByID() http.Handler
	GetVideos() http.Handler
	GetTickets() http.Handler
}
