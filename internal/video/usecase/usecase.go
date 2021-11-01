// Package usecase implements video service usecases.
package usecase

import (
	"context"
	"encoding/json"
	"io"

	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/internal/storage"
	"github.com/grigagod/compresso/internal/utils"
	"github.com/grigagod/compresso/internal/video"
	"github.com/grigagod/compresso/pkg/rmq"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// VideoUseCase implement video.UseCase interface.
type VideoUseCase struct {
	qwCfg     *rmq.QueueWriteConfig
	repo      video.Repository
	storage   storage.Storage
	publisher *rmq.Publisher
}

func NewVideoUseCase(qwCfg *rmq.QueueWriteConfig, repo video.Repository, storage storage.Storage, publisher *rmq.Publisher) *VideoUseCase {
	return &VideoUseCase{
		qwCfg:     qwCfg,
		repo:      repo,
		storage:   storage,
		publisher: publisher,
	}
}

// CreateVideo upload video to the storage and make DB insert.
func (u *VideoUseCase) CreateVideo(ctx context.Context, video *models.Video, file io.Reader) (*models.Video, error) {
	url, err := utils.GenerateURL(video.AuthorID, video.ID)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.CreateVideo.GenerateURL")
	}
	video.URL = url

	err = u.storage.PutObject(ctx, file, video.URL)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.CreateVideo.PutObject")
	}

	v, err := u.repo.InsertVideo(ctx, video)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.CreateVideo.InsertVideo")
	}

	if err := u.signVideoURL(v); err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.CreateVideo.GetDownloadURL")
	}

	return v, nil
}

// CreateTicket find video in DB, send message for processing to the broker and update DB.
func (u *VideoUseCase) CreateTicket(ctx context.Context, ticket *models.VideoTicket) (*models.VideoTicket, error) {
	video, err := u.repo.SelectVideoByID(ctx, ticket.AuthorID, ticket.VideoID)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.CreateTicket.SelectVideoByID")
	}

	url, err := utils.GenerateURL(video.AuthorID, ticket.ID)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.CreateTicket.GenerateURL")
	}

	msg := &models.QueueVideoMsg{
		TicketID:     ticket.ID,
		CRF:          ticket.CRF,
		TargetFormat: ticket.TargetFormat,
		OriginURL:    video.URL,
		ProcessedURL: url,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.CreateTicket.Marshal")
	}

	ticket.State = models.Queued
	t, err := u.repo.InsertTicket(ctx, ticket)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.CreateTicket.InsertTicket")
	}

	err = u.publisher.Send(rmq.NewMessage(amqp.Publishing{
		Headers:     map[string]interface{}{},
		ContentType: rmq.JSONContentType,
		Body:        body,
	}, u.qwCfg))
	if err != nil {
		ticket.State = models.Failed
		_, err := u.repo.UpdateTicket(ctx, ticket)

		return nil, errors.Wrap(err, "VideoUseCase.CreateTicket.Send")
	}

	return t, nil
}

// GetVideoByID return author`s video with the given ID.
func (u *VideoUseCase) GetVideoByID(ctx context.Context, authorID, id uuid.UUID) (*models.Video, error) {
	video, err := u.repo.SelectVideoByID(ctx, authorID, id)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.GetVideoByID.SelectVideoByID")
	}

	if err := u.signVideoURL(video); err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.GetVideoByID.GetDownloadURL")
	}

	return video, nil
}

// GetTicketByID return author`s video ticket with the given ID.
func (u *VideoUseCase) GetTicketByID(ctx context.Context, authorID, id uuid.UUID) (*models.VideoTicket, error) {
	ticket, err := u.repo.SelectTicketByID(ctx, authorID, id)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.GetTicketByID.SelectTicketByID")
	}

	if err := u.signTicketURL(ticket); err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.GetTicketByID.GetDownloadURL")
	}

	return ticket, nil
}

// GetVideos return author`s videos.
func (u *VideoUseCase) GetVideos(ctx context.Context, authorID uuid.UUID) ([]*models.Video, error) {
	videos, err := u.repo.SelectVideos(ctx, authorID)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.GetVideos.SelectVideos")
	}

	// sign urls for videos
	for _, v := range videos {
		if err := u.signVideoURL(v); err != nil {
			return nil, errors.Wrap(err, "VideoUseCase.GetVideos.GetDownloadURL")
		}
	}

	return videos, nil
}

// GetTickets return author`s tickets.
func (u *VideoUseCase) GetTickets(ctx context.Context, authorID uuid.UUID) ([]*models.VideoTicket, error) {
	tickets, err := u.repo.SelectTickets(ctx, authorID)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.GetTickets.SelectTickets")
	}

	// sign urls for processed tickets
	for _, t := range tickets {
		if err := u.signTicketURL(t); err != nil {
			return nil, errors.Wrap(err, "VideoUseCase.GetVideos.GetDownloadURL")
		}
	}

	return tickets, nil
}

// signVideoURL sign video URL.
func (u *VideoUseCase) signVideoURL(video *models.Video) error {
	url, err := u.storage.GetDownloadURL(video.URL)
	if err == nil {
		video.URL = url
	}

	return err
}

// signTicketURL sign video ticket URL.
func (u *VideoUseCase) signTicketURL(ticket *models.VideoTicket) error {
	if ticket.State == models.Done {
		url, err := u.storage.GetDownloadURL(ticket.URL)
		if err == nil {
			ticket.URL = url
		}

		return err
	} else {
		ticket.URL = ""

		return nil
	}
}
