//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package auth

import "github.com/grigagod/compresso/internal/models"

type UseCase interface {
	Register(user *models.User) (*models.UserWithToken, error)
	Login(user *models.User) (*models.UserWithToken, error)
}
