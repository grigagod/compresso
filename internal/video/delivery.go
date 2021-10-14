package video

import "net/http"

type Handlers interface {
	UploadVideo() http.Handler
	CreateTicket() http.Handler
}
