package video

import (
	"context"
	"io"

	"github.com/grigagod/compresso/internal/models"
)

type UseCase interface {
	UploadVideo(ctx context.Context, video *models.Video, file io.Reader) (*models.Video, error)
	CreateTicket(ctx context.Context, ticket *models.VideoTicket) (*models.VideoTicket, error)
}
