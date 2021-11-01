package video

import (
	"context"

	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/models"
)

type Repository interface {
	InsertVideo(ctx context.Context, video *models.Video) (*models.Video, error)
	InsertTicket(ctx context.Context, ticket *models.VideoTicket) (*models.VideoTicket, error)
	UpdateTicket(ctx context.Context, ticket *models.VideoTicket) (*models.VideoTicket, error)
	SelectVideoByID(ctx context.Context, authorID, id uuid.UUID) (*models.Video, error)
	SelectTicketByID(ctx context.Context, authorID, id uuid.UUID) (*models.VideoTicket, error)
	SelectVideos(ctx context.Context, authorID uuid.UUID) ([]*models.Video, error)
	SelectTickets(ctx context.Context, authorID uuid.UUID) ([]*models.VideoTicket, error)
}
