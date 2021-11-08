package storage

import (
	"context"
	"io"

	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/internal/storage"
	"github.com/grigagod/compresso/internal/utils"
	"github.com/pkg/errors"
)

type VideoStorage struct {
	storage storage.Storage
}

func NewVideoStorage(storage storage.Storage) *VideoStorage {
	return &VideoStorage{storage: storage}
}

func (s *VideoStorage) GetVideo(ctx context.Context, video *models.Video) (io.ReadCloser, error) {
	src, err := s.storage.GetObject(ctx, video.URL)
	if err != nil {
		return nil, errors.Wrap(err, "VideoStorage.GetVideo.GetObject")
	}

	return src, nil
}

func (s *VideoStorage) PutVideo(ctx context.Context, video *models.Video, file io.Reader) error {
	fileType, err := utils.DetectVideoMIMEType(video.Format)
	if err != nil {
		return errors.Wrap(err, "VideoStorage.PutVideo.DetectVideoMIMEType")
	}
	if err := s.storage.PutObject(ctx, file, video.URL, fileType); err != nil {
		return errors.Wrap(err, "VideoStorage.PutVideo.PutObject")
	}

	return nil
}

func (s *VideoStorage) SignVideoURL(video *models.Video) error {
	url, err := s.storage.GetDownloadURL(video.URL)
	if err != nil {
		return errors.Wrap(err, "VideoStorage.SignVideoURL.GetDownloadURL")
	}
	video.URL = url

	return nil
}

func (s *VideoStorage) SignTicketURL(ticket *models.VideoTicket) (err error) {
	var url string
	if ticket.State == models.Done {
		url, err = s.storage.GetDownloadURL(ticket.URL)
		if err != nil {
			return errors.Wrap(err, "VideoStorage.SignTicketURL.GetDownloadURL")
		}
	}

	ticket.URL = url

	return nil
}
