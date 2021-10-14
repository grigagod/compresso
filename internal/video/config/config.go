package config

import (
	"github.com/grigagod/compresso/internal/httpserver"
	"github.com/grigagod/compresso/internal/storage"
	"github.com/grigagod/compresso/pkg/db/postgres"
	"github.com/grigagod/compresso/pkg/rmq"
)

// API stores full config for video API service(include initialization and hot used parts).
type API struct {
	HTTP    httpserver.Config
	DB      postgres.Config
	Storage storage.Config
	RMQ     rmq.Config
	APIsvc
}

// APIsvc stores hot used part of config for video API service.
type APIsvc struct {
	rmq.QueueWriteConfig
	JwtSecretKey string
}

// Converer stores config for video converter Service.
type Converter struct {
	DB      postgres.Config
	Storage storage.Config
	RMQ     rmq.Config
	rmq.QueueReadConfig
}
