package rmq

import (
	"bytes"
	"context"

	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/internal/utils"
	"github.com/grigagod/compresso/internal/video/svc"
	"github.com/grigagod/compresso/pkg/rmq"
)

type SVCHandlers struct {
	svcUC svc.UseCase
}

func NewSVCHandlers(useCase svc.UseCase) *SVCHandlers {
	return &SVCHandlers{
		svcUC: useCase,
	}
}

func (h *SVCHandlers) ProcessVideo() rmq.Handler {
	fn := func(ctx context.Context, body []byte) error {
		var msg models.ProcessVideoMsg

		err := utils.StructScan(bytes.NewReader(body), &msg)
		if err != nil {
			return err
		}

		err = h.svcUC.ProcessVideo(ctx, &msg)

		return err
	}

	return rmq.HandlerFunc(fn)
}
