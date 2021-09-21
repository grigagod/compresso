package rmq

import (
	"strings"

	"github.com/streadway/amqp"
)

// Config contains url for connection to rmq and additional options.
type Config struct {
	Proto    string `default:"amqp://"`
	Host     string `default:"127.0.0.1"`
	Port     string `default:"5672"`
	User     string
	Password string
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
	Name      string
	Mandatory bool
	Immediate bool
}

// QueueReadConfig describe read queue configuration.
type QueueReadConfig struct {
	Args       amqp.Table
	Bind       []BindQueueConfig
	Name       string
	Consumer   string
	Durable    bool
	AutoDelete bool
	AutoAck    bool
	Exclusive  bool
	NoLocal    bool
	NoWait     bool
}

// BindQueueConfig describe bind queue config.
type BindQueueConfig struct {
	Key      string
	Exchange string
	NoWait   bool
	Args     amqp.Table
}
