package auth

import "github.com/grigagod/compresso/internal/auth/models"

type UseCase interface {
	Register(user *models.User) (*models.UserWithToken, error)
	Login(user *models.User) (*models.UserWithToken, error)
}
