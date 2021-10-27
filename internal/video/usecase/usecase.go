// Package usecase implements video service usecases.
package usecase

import (
	"context"
	"encoding/json"
	"io"

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

// UploadVideo upload video to the storage and make DB insert.
func (u *VideoUseCase) UploadVideo(ctx context.Context, video *models.Video, file io.Reader) (*models.Video, error) {
	url, err := utils.GenerateURL(video.AuthorID, video.ID)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.UploadVideo.GenerateURL")
	}
	video.URL = url

	err = u.storage.PutObject(ctx, file, video.URL)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.UploadVideo.PutObject")
	}

	v, err := u.repo.CreateVideo(ctx, video)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.UploadVideo.CreateVideo")
	}

	url, err = u.storage.GetDownloadURL(v.URL)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.UploadVideo.GetDownloadURL")
	}

	v.URL = url

	return v, nil
}

// CreateTicket find video in db, send message for processing to the broker and update DB.
func (u *VideoUseCase) CreateTicket(ctx context.Context, ticket *models.VideoTicket) (*models.VideoTicket, error) {
	video, err := u.repo.FindVideoByID(ctx, ticket.VideoID)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.CreateTicket.FindVideoByID")
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
	t, err := u.repo.CreateTicket(ctx, ticket)
	if err != nil {
		return nil, errors.Wrap(err, "VideoUseCase.CreateTicket.CreateTicket")
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
