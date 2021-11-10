package svc

import "github.com/grigagod/compresso/pkg/rmq"

type Handlers interface {
	ProcessVideo() rmq.Handler
}
