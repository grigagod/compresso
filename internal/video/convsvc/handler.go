package convsvc

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/pkg/converter"
	"github.com/pkg/errors"
)

// handleMsg handle incoming messages(process video tickets).
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
	if err != nil {
		upderr := s.updateTicketState(ctx, ticket, models.Failed)
		return errors.Wrap(upderr, err.Error())
	}

	s.logger.Infow("Start converting ticket", "ID", msg.TicketID.String())

	// buffer is used as Writer for converter output and as Reader for upload to the storage
	buf := new(bytes.Buffer)
	err = converter.ProcessVideo(ctx, src, buf, msg.TargetFormat, msg.CRF)
	src.Close()
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

	s.logger.Infow("Ticket handling finished", "ID", msg.TicketID.String())
	return nil
}

func (s *Service) updateTicketState(ctx context.Context, ticket *models.VideoTicket, state models.ProcessingState) error {
	ticket.State = state
	_, err := s.repo.UpdateTicket(ctx, ticket)
	return err
}
