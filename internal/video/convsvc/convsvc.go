package convsvc

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/internal/storage"
	"github.com/grigagod/compresso/internal/video"
	"github.com/grigagod/compresso/internal/video/repo"
	"github.com/grigagod/compresso/pkg/converter"
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

			func(msg amqp.Delivery) {
				defer func() {
					if rval := recover(); rval != nil {
						s.logger.Errorw("Panic when processing rmq message", "rval", rval, "msg", msg.Body)
						err = msg.Nack(false, false)
						if err != nil {
							s.logger.Errorw("Errorw nack message", "err", err, "msg", msg.Body)
						}
					}
				}()

				err := s.handleMsg(ctx, msg.Body)
				if err != nil {
					s.logger.Errorw("Error occured while processing msg", "error", err.Error(), "msg", msg.Body)

					err := msg.Nack(false, true)
					if err != nil {
						s.logger.Errorw("Errorw nack message", "err", err, "msg", msg.Body)
					}
					return
				}

				err = msg.Ack(false)
				if err != nil {
					s.logger.Errorw("Errorw ack message", "err", err, "msg", msg.Body)
				}

			}(msg)
		}
	}

	return err
}

func (s *Service) handleMsg(ctx context.Context, body []byte) error {
	var msg models.QueueVideoMsg
	// decode incoming msg
	err := json.NewDecoder(bytes.NewReader(body)).Decode(&msg)
	if err != nil {
		return err
	}

	s.logger.Infow("Start ticket handling", "ID", msg.TicketID.String())

	// create ticket for future updates in DB
	ticket := &models.VideoTicket{
		Ticket: models.Ticket{
			ID:  msg.TicketID,
			URL: msg.ProcessedURL,
		},
		TargetFormat: msg.TargetFormat,
		CRF:          msg.CRF,
	}

	// update ticket status in DB.
	err = s.updateTicketState(ctx, ticket, models.Processing)
	if err != nil {
		return err
	}

	s.logger.Infow("Downloading original ticket", "ID", msg.TicketID.String())
	src, err := s.storage.GetObject(ctx, msg.OriginURL)
	defer src.Close()
	if err != nil {
		upderr := s.updateTicketState(ctx, ticket, models.Failed)
		return errors.Wrap(upderr, err.Error())
	}

	s.logger.Infow("Start converting ticket", "ID", msg.TicketID.String())

	// buffer is used as Writer for converter output and as Reader for upload to the storage
	buf := new(bytes.Buffer)
	err = converter.ProcessVideo(ctx, src, buf, msg.TargetFormat, msg.CRF)
	if err != nil {
		upderr := s.updateTicketState(ctx, ticket, models.Failed)
		return errors.Wrap(upderr, err.Error())
	}

	s.logger.Infow("Uploading processed ticket", "ID", msg.TicketID.String())
	err = s.storage.PutObject(ctx, buf, msg.ProcessedURL)
	if err != nil {
		ticket.State = models.Failed
		_, err = s.repo.UpdateTicket(ctx, ticket)
		return err
	}

	err = s.updateTicketState(ctx, ticket, models.Done)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) updateTicketState(ctx context.Context, ticket *models.VideoTicket, state models.ProcessingState) error {
	ticket.State = state
	_, err := s.repo.UpdateTicket(ctx, ticket)
	return err
}
