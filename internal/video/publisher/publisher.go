package publisher

import (
	"encoding/json"

	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/pkg/rmq"
	"github.com/streadway/amqp"
)

type VideoPublisher struct {
	qwCfg *rmq.QueueWriteConfig
	ch    *rmq.Channel
}

func NewVideoPublisher(cfg *rmq.QueueWriteConfig, ch *rmq.Channel) *VideoPublisher {
	return &VideoPublisher{
		qwCfg: cfg,
		ch:    ch,
	}
}

func (p *VideoPublisher) SendMsg(msg *models.ProcessVideoMsg) error {
	pub := rmq.NewPublisher(p.ch)

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = pub.Send(rmq.NewMessage(amqp.Publishing{
		Headers:      map[string]interface{}{},
		ContentType:  rmq.JSONContentType,
		Body:         body,
		DeliveryMode: amqp.Persistent,
	}, p.qwCfg))

	return err

}
