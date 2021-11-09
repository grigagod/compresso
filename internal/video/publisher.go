package video

import "github.com/grigagod/compresso/internal/models"

type Publisher interface {
	SendMsg(msg *models.ProcessVideoMsg) error
}
