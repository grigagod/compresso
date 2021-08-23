package auth

import "github.com/grigagod/compresso/internal/auth/models"

type Repository interface {
	Create(user *models.User) (*models.User, error)
	FindByName(username string) (*models.User, error)
}
