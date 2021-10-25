// Package rmq defines set of functions and types for working with RabbitMQ.
package rmq

import "github.com/streadway/amqp"

// Channel channel for rmq.
type Channel struct {
	connection *amqp.Connection
	*amqp.Channel
}

// NewChannelFromConfig open new channel.
func NewChannelFromConfig(config *Config) (*Channel, error) {
	conn, err := amqp.Dial(config.URL())
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Channel{
		connection: conn,
		Channel:    ch,
	}, nil
}

// Close close channel and connection.
func (c *Channel) Close() error {
	err := c.Channel.Close()
	if err != nil {
		return err
	}

	return c.connection.Close()
}
