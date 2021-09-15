package rmq

import "github.com/streadway/amqp"

// Consumer describe consumer.
type Consumer struct {
	Channel *Channel
}

// NewConsumer return new consumer.
func NewConsumer(ch *Channel) *Consumer {
	return &Consumer{
		Channel: ch,
	}
}

// Receiving receive message from rmq queue.
func (c *Consumer) Receiving(queueConfig *QueueReadConfig) (<-chan amqp.Delivery, error) {
	q, err := c.Channel.QueueDeclare(
		queueConfig.Name,
		queueConfig.Durable,
		queueConfig.AutoDelete,
		queueConfig.Exclusive,
		queueConfig.NoWait,
		queueConfig.Args,
	)
	if err != nil {
		return nil, err
	}

	for _, bindConfig := range queueConfig.Bind {
		err = c.Channel.QueueBind(queueConfig.Name, bindConfig.Key, bindConfig.Exchange, bindConfig.NoWait, bindConfig.Args)
		if err != nil {
			return nil, err
		}
	}

	return c.Channel.Consume(
		q.Name,                // queue
		queueConfig.Consumer,  // consumer
		queueConfig.AutoAck,   // auto-ack
		queueConfig.Exclusive, // exclusive
		queueConfig.NoLocal,   // no-local
		queueConfig.NoWait,    // no-wait
		queueConfig.Args,      // args
	)
}
