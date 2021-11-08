//go:generate mockgen -source storage.go -destination mock/storage_mock.go -package mock
package video

import (
	"context"
	"io"

	"github.com/grigagod/compresso/internal/models"
)

type Storage interface {
	GetVideo(ctx context.Context, video *models.Video) (io.ReadCloser, error)
	PutVideo(ctx context.Context, video *models.Video, file io.Reader) error
	SignVideoURL(video *models.Video) error
	SignTicketURL(ticket *models.VideoTicket) error
}
