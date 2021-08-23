package config

import (
	"log"
	"time"

	"github.com/grigagod/compresso/pkg/utils"
)

// Config stores auth service config.
type Config struct {
	Server
	Auth
}

// Auth stores jwt config.
type Auth struct {
	JwtSecretKey string
	JwtExpires   time.Duration
}

// Server stores http server config.
type Server struct {
	Addr         string
	PprofPort    string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
}

func LoadConfig(filepath string) (*Config, error) {
	v, err := utils.LoadConfig(filepath)
	if err != nil {
		log.Printf("unable to load auth config, %v", err)
	}

	var c Config

	err = v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
