package video

import "net/http"

type Delivery interface {
	UploadVideo() http.Handler
	CreateTicket() http.Handler
}
