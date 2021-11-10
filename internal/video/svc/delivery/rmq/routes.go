package rmq

import (
	"github.com/grigagod/compresso/internal/video"
	"github.com/grigagod/compresso/internal/video/svc"
	"github.com/grigagod/compresso/pkg/rmq"
)

func MapHandlers(router *rmq.Router, h svc.Handlers) {
	router.Add(video.ProcessVideoHeader, h.ProcessVideo())
}
