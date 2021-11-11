package svc

import (
	"context"

	"github.com/grigagod/compresso/internal/models"
)

type UseCase interface {
	ProcessVideo(ctx context.Context, msg *models.ProcessVideoMsg) error
}
