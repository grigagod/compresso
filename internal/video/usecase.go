//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package video

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/models"
)

type UseCase interface {
	CreateVideo(ctx context.Context, video *models.Video, file io.Reader) (*models.Video, error)
	CreateTicket(ctx context.Context, ticket *models.VideoTicket) (*models.VideoTicket, error)
	GetVideoByID(ctx context.Context, authorID, id uuid.UUID) (*models.Video, error)
	GetTicketByID(ctx context.Context, authorID, id uuid.UUID) (*models.VideoTicket, error)
	GetVideos(ctx context.Context, authorID uuid.UUID) ([]*models.Video, error)
	GetTickets(ctx context.Context, authorID uuid.UUID) ([]*models.VideoTicket, error)
}
