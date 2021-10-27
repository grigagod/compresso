package rmq

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/streadway/amqp"
)

// Environment variables prefixes.
const (
	ConfigPrefix           = "RMQ"
	QueueWriteConfigPrefix = "QWC"
	QueueReadConfigPrefix  = "QRC"
)

// Config contains url for connection to rmq and additional options.
type Config struct {
	Proto    string `default:"amqp://"`
	Host     string `default:"127.0.0.1"`
	Port     string `default:"5672"`
	User     string `required:"true"`
	Password string `required:"true"`
}

// GetConfigFromEnv return rmq config from environment.
func GetConfigFromEnv() (*Config, error) {
	c := new(Config)
	err := envconfig.Process(ConfigPrefix, c)
	return c, err
}

// NewAMQPConfig return new amqp connection config.
func NewAMQPConfig(proto, host, port, user, password string) *Config {
	return &Config{
		Proto:    proto,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}
}

// URL return connection url.
func (r *Config) URL() string {
	var parts []string

	parts = append(parts, r.Proto)

	if r.Password != "" {
		parts = append(parts, r.User, ":", r.Password, "@")
	}

	parts = append(parts, r.Host, ":", r.Port)

	return strings.Join(parts, "")
}

// QueueWriteConfig describe write queue configuration.
type QueueWriteConfig struct {
	Exchange  string
	Name      string `required:"true"`
	Mandatory bool   `default:"true"`
	Immediate bool   `default:"false"`
}

// GetConfigFromEnv return queue write config from environment.
func GetQueueWriteConfigFromEnv() (*QueueWriteConfig, error) {
	c := new(QueueWriteConfig)
	err := envconfig.Process(QueueWriteConfigPrefix, c)
	return c, err
}

// QueueReadConfig describe read queue configuration.
type QueueReadConfig struct {
	Args       amqp.Table
	Bind       []BindQueueConfig
	Name       string `required:"true"`
	Consumer   string `required:"true"`
	Durable    bool   `required:"true" default:"true"`
	AutoDelete bool   `split_words:"true" default:"false"`
	AutoAck    bool   `split_words:"true" default:"false"`
	Exclusive  bool   `default:"false"`
	NoLocal    bool   `split_words:"true" default:"false"`
	NoWait     bool   `split_words:"true" default:"false"`
}

// GetConfigFromEnv return queue read config from environment.
func GetQueueReadConfigFromEnv() (*QueueReadConfig, error) {
	c := new(QueueReadConfig)
	err := envconfig.Process(QueueReadConfigPrefix, c)
	return c, err
}

// BindQueueConfig describe bind queue configuration.
type BindQueueConfig struct {
	Key      string
	Exchange string
	NoWait   bool `split_words:"true"`
	Args     amqp.Table
}
