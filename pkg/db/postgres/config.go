package postgres

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// ConfigPrefix environment variables prefix.
const ConfigPrefix = "DB"

// Config could be used for configuring new postgres sqlx.DB instance.
type Config struct {
	URL             string        `required:"true"`
	Driver          string        `default:"pgx"`
	MaxOpenConns    int           `split_words:"true"`
	MaxIdleConns    int           `split_words:"true"`
	ConnMaxLifetime time.Duration `split_words:"true"`
	ConnMaxIdleTime time.Duration `split_words:"true"`
}

// GetConfigFromEnv return postgres config from environment.
func GetConfigFromEnv() (*Config, error) {
	c := new(Config)
	err := envconfig.Process(ConfigPrefix, c)
	return c, err
}
