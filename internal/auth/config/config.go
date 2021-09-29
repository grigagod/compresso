package config

import (
	"time"

	"github.com/grigagod/compresso/internal/httpserver"
	"github.com/grigagod/compresso/pkg/db/postgres"
)

// Config stores auth service config.
type Config struct {
	HTTP httpserver.Config
	DB   postgres.Config
	Auth
}

// Auth stores jwt config.
type Auth struct {
	JwtSecretKey string
	JwtExpires   time.Duration
}
