package server

import (
	"time"
)

// Config stores http server config.
type Config struct {
	Addr         string
	SwaggerUrl   string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
}
