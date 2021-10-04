package video

import (
	"context"

	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/internal/storage"
)

type UseCase interface {
	UploadVideo(ctx context.Context, author_id uuid.UUID, input storage.UploadInput) (*models.Video, error)
	CreateTicket(ctx context.Context, ticket *models.VideoTicket) (*models.VideoTicket, error)
}
