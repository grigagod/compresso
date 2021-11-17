package usecase

import (
	"bytes"
	"context"

	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/internal/storage"
	"github.com/grigagod/compresso/internal/utils"
	"github.com/grigagod/compresso/internal/video"
	"github.com/grigagod/compresso/pkg/converter"
	"github.com/grigagod/compresso/pkg/logger"
	"github.com/pkg/errors"
)

type SVCUseCase struct {
	repo    video.Repository
	storage storage.Storage
	logger  logger.Logger
}

func NewSVCUseCase(repo video.Repository, storage storage.Storage, logger logger.Logger) *SVCUseCase {
	return &SVCUseCase{
		repo:    repo,
		storage: storage,
		logger:  logger,
	}
}

func (u *SVCUseCase) ProcessVideo(ctx context.Context, msg *models.ProcessVideoMsg) error {
	// create ticket for future updates in DB
	ticket := &models.VideoTicket{
		Ticket: models.Ticket{
			ID:  msg.TicketID,
			URL: msg.ProcessedURL,
		},
		TargetFormat: msg.TargetFormat,
		CRF:          msg.CRF,
	}

	// determine MIME type for the target video
	format, ok := utils.DetectVideoMIMEType(ticket.TargetFormat)
	if !ok {
		return u.updateTicketState(ctx, ticket, models.Failed)
	}

	// update ticket status in DB
	if err := u.updateTicketState(ctx, ticket, models.Processing); err != nil {
		return err
	}

	u.logger.Infow("Downloading original ticket", "ID", msg.TicketID.String())

	src, err := u.storage.GetObject(ctx, msg.OriginURL)
	if err != nil {
		upderr := u.updateTicketState(ctx, ticket, models.Failed)
		return errors.Wrap(upderr, err.Error())
	}

	u.logger.Infow("Start converting ticket", "ID", msg.TicketID.String())

	// buffer is used as Writer for converter output and as Reader for upload to the storage
	buf := new(bytes.Buffer)
	err = converter.ProcessVideo(ctx, src, buf, msg.TargetFormat, msg.CRF)
	src.Close()
	if err != nil {
		upderr := u.updateTicketState(ctx, ticket, models.Failed)
		return errors.Wrap(upderr, err.Error())
	}

	u.logger.Infow("Uploading processed ticket", "ID", msg.TicketID.String())

	err = u.storage.PutObject(ctx, buf, ticket.URL, format)
	if err != nil {
		return err
	}

	err = u.updateTicketState(ctx, ticket, models.Done)
	if err != nil {
		upderr := u.updateTicketState(ctx, ticket, models.Failed)
		return errors.Wrap(upderr, err.Error())
	}

	u.logger.Infow("Ticket handling finished", "ID", msg.TicketID.String())

	return nil
}

func (u *SVCUseCase) updateTicketState(ctx context.Context, ticket *models.VideoTicket, state models.ProcessingState) error {
	ticket.State = state
	_, err := u.repo.UpdateTicket(ctx, ticket)
	return err
}
