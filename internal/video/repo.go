package video

import (
	"context"

	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/models"
)

type Repository interface {
	CreateVideo(ctx context.Context, video *models.Video) (*models.Video, error)
	CreateTicket(ctx context.Context, ticket *models.VideoTicket) (*models.VideoTicket, error)
	UpdateTicket(ctx context.Context, ticket *models.VideoTicket) (*models.VideoTicket, error)
	GetVideoByID(ctx context.Context, id uuid.UUID) (*models.Video, error)
	GetTicketByID(ctx context.Context, id uuid.UUID) (*models.VideoTicket, error)
	GetVideos(ctx context.Context, authorID uuid.UUID) ([]*models.Video, error)
	GetTickets(ctx context.Context, authorID uuid.UUID) ([]*models.VideoTicket, error)
}
