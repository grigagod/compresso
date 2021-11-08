package convsvc

import (
	"context"
	"time"

	"github.com/grigagod/compresso/internal/storage"
	"github.com/grigagod/compresso/internal/video"
	"github.com/grigagod/compresso/internal/video/repo"
	"github.com/grigagod/compresso/pkg/logger"
	"github.com/grigagod/compresso/pkg/rmq"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Service struct {
	channel *rmq.Channel
	repo    video.Repository
	storage storage.Storage
	logger  logger.Logger
}

func NewService(channel *rmq.Channel, db *sqlx.DB, storage storage.Storage, logger logger.Logger) *Service {
	return &Service{
		channel: channel,
		repo:    repo.NewVideoRepository(db),
		storage: storage,
		logger:  logger,
	}
}

func (s *Service) Run(ctx context.Context, qrCfg *rmq.QueueReadConfig) error {
	cn := rmq.NewConsumer(s.channel)
	d, err := cn.Receiving(qrCfg)
	if err != nil {
		return err
	}

LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case msg := <-d:
			if msg.Acknowledger == nil {
				s.logger.Errorw("Acknowledger is nil. Channel Was broken", "msg", msg.Body)
				break LOOP
			}

			func(ctx context.Context, msg amqp.Delivery) {
				defer func() {
					if rval := recover(); rval != nil {
						s.logger.Errorw("Panic when processing rmq message", "rval", rval, "msg", msg.Body)
						err = msg.Nack(false, false)
						if err != nil {
							s.logger.Errorw("Errorw nack message", "err", err, "msg", msg.Body)
						}
					}
				}()
				ctx, cancel := context.WithTimeout(ctx, time.Minute)
				defer cancel()
				err := s.handleMsg(ctx, msg.Body)
				if errors.Is(err, context.Canceled) {
					// requeue msg when err occurred because of program termination
					err := msg.Nack(false, true)
					if err != nil {
						s.logger.Errorw("Errorw nack message", "err", err, "msg", msg.Body)
					}
					return
				}
				if err != nil {
					s.logger.Errorw("Error occurred while processing msg", "error", err.Error(), "msg", msg.Body)
					err := msg.Nack(false, false)
					if err != nil {
						s.logger.Errorw("Errorw nack message", "err", err, "msg", msg.Body)
					}
					return
				}

				err = msg.Ack(false)
				if err != nil {
					s.logger.Errorw("Errorw ack message", "err", err, "msg", msg.Body)
				}
			}(ctx, msg)
		}
	}

	s.logger.Infof("Gracefully shut down")

	return err
}
