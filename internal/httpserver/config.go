package httpserver

import (
	"time"
)

// Config stores http server config.
type Config struct {
	Addr         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
	Pprof        bool
	Swagger      bool
	SwaggerUrl   string
}
