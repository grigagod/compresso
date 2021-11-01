package config

import (
	"github.com/grigagod/compresso/internal/httpserver"
	"github.com/grigagod/compresso/pkg/rmq"
	"github.com/kelseyhightower/envconfig"
)

const HTTPConfigPrefix = "VIDEO"

// API stores config for video API service.
type API struct {
	*rmq.QueueWriteConfig `ignored:"true"`
	JwtSecretKey          string `split_words:"true"`
}

func GetHTTPConfigFromEnv() (*httpserver.Config, error) {
	c := new(httpserver.Config)
	err := envconfig.Process(HTTPConfigPrefix, c)
	return c, err
}

func GetAPIConfigFromEnv() (*API, error) {
	qwCfg, err := rmq.GetQueueWriteConfigFromEnv()
	if err != nil {
		return nil, err
	}

	c := new(API)
	err = envconfig.Process("", c)
	c.QueueWriteConfig = qwCfg

	return c, err
}
