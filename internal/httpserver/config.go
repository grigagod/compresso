package httpserver

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// ConfigPrefix environment variables prefix.
const ConfigPrefix = "HTTP"

// Config stores http server config.
type Config struct {
	Addr         string        `required:"true"`
	WriteTimeout time.Duration `split_words:"true"`
	ReadTimeout  time.Duration `split_words:"true"`
	IdleTimeout  time.Duration `split_words:"true"`
	Pprof        bool          `default:"false"`
	Swagger      bool          `default:"false"`
	SwaggerURL   string        `split_words:"true"`
}

// GetConfigFromEnv return http server config from environment.
func GetConfigFromEnv() (*Config, error) {
	c := new(Config)
	err := envconfig.Process(ConfigPrefix, c)
	return c, err
}
