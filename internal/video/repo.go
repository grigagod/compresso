package video

import (
	"context"

	"github.com/grigagod/compresso/internal/models"
)

type Repository interface {
	Create(ctx context.Context, video *models.Video) (*models.Video, error)
	CreateTicket(ctx context.Context, ticket *models.VideoTicket) (*models.VideoTicket, error)
	UpdateTicket(ctx context.Context, ticket *models.VideoTicket) error
}
