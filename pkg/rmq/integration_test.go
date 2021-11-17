//go:build integration

package rmq

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_PublisherConsumerExchange(t *testing.T) {
	conf := Config{
		Proto: "amqp://",
		Host:  "127.0.0.1",
		Port:  "5672",
	}

	ch, err := NewChannelFromConfig(&conf)
	assert.Nil(t, err)
	_, err = ch.QueueDeclare("test_queue", true, false, false, false, amqp.Table{})
	assert.Nil(t, err)

	pub := NewPublisher(ch)
	err = pub.Send(NewMessage(amqp.Publishing{
		Headers:     map[string]interface{}{},
		ContentType: PlainContentType,
		Body:        []byte("test1234"),
	}, &QueueWriteConfig{
		Name:      "test_queue",
		Mandatory: true,
	}))
	assert.Nil(t, err)

	cn := NewConsumer(ch)
	delivery, err := cn.Receiving(&QueueReadConfig{
		Name:    "test_queue",
		Durable: true,
	})
	assert.Nil(t, err)

	for msg := range delivery {
		err = msg.Ack(false)
		assert.Nil(t, err)
		break
	}
	err = ch.Close()
	assert.Nil(t, err)
}
