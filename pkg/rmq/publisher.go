package rmq

import "github.com/streadway/amqp"

// RMQ constants.
const (
	JSONContentType  = "application/json"
	PlainContentType = "text/plain"
)

// Message contains message and config.
type Message struct {
	Publishing amqp.Publishing
	Config     *QueueWriteConfig
}

// NewMessage return new message.
func NewMessage(message amqp.Publishing, queueConfig *QueueWriteConfig) Message {
	return Message{
		Publishing: message,
		Config:     queueConfig,
	}
}

// Publisher describe publisher.
type Publisher struct {
	Channel *Channel
}

// NewPublisher return new publisher.
func NewPublisher(ch *Channel) *Publisher {
	return &Publisher{
		Channel: ch,
	}
}

// Send send message to rmq.
func (c *Publisher) Send(message Message) error {
	return c.Channel.Publish(
		message.Config.Exchange,
		message.Config.Name,
		message.Config.Mandatory,
		message.Config.Immediate,
		message.Publishing)
}
