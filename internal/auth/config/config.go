package config

import (
	"time"

	"github.com/grigagod/compresso/internal/httpserver"
	"github.com/kelseyhightower/envconfig"
)

const HTTPConfigPrefix = "AUTH"

// Auth stores jwt config.
type Auth struct {
	JwtSecretKey string        `split_words:"true"`
	JwtExpires   time.Duration `split_words:"true" default:"15m"`
}

func GetAuthConfigFromEnv() (*Auth, error) {
	c := new(Auth)
	err := envconfig.Process("", c)
	return c, err
}

func GetHTTPConfigFromEnv() (*httpserver.Config, error) {
	c := new(httpserver.Config)
	err := envconfig.Process(HTTPConfigPrefix, c)
	return c, err
}
