//go:generate mockgen -source publisher.go -destination mock/publisher_mock.go -package mock
package video

import "github.com/grigagod/compresso/internal/models"

type Publisher interface {
	SendMsg(msg *models.ProcessVideoMsg) error
}
