package storage

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// ConfigPrefix environment variables prefix.
const ConfigPrefix = "STORAGE"

// Config should be used for configuring new storage instance.
type Config struct {
	Bucket          string
	PresignDuration time.Duration `split_words:"true"`
}

// GetConfigFromEnv return storage config from environment.
func GetConfigFromEnv() (*Config, error) {
	c := new(Config)
	err := envconfig.Process(ConfigPrefix, c)
	return c, err
}
