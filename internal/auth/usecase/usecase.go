package usecase

import (
	"github.com/grigagod/compresso/internal/auth"
	"github.com/grigagod/compresso/internal/auth/config"
	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/pkg/utils"
	"github.com/pkg/errors"
)

type authUseCase struct {
	cfg      *config.Config
	authRepo auth.Repository
}

func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository) auth.UseCase {
	return &authUseCase{cfg: cfg, authRepo: authRepo}
}

func (u *authUseCase) Register(user *models.User) (*models.UserWithToken, error) {
	existsUser, err := u.authRepo.FindByName(user.Username)
	if existsUser != nil || err == nil {
		return nil, errors.New("authUseCase.Register.UserExists")
	}

	if err = user.HashPassword(); err != nil {
		return nil, errors.Wrap(err, "authUseCase.Register.HashPassword")
	}

	createdUser, err := u.authRepo.Create(user)
	if err != nil {
		return nil, errors.Wrap(err, "authUseCase.Register.Create")
	}

	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser.ID, u.cfg.JwtExpires, u.cfg.JwtSecretKey)

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

func (u *authUseCase) Login(user *models.User) (*models.UserWithToken, error) {
	foundUser, err := u.authRepo.FindByName(user.Username)
	if err != nil {
		return nil, errors.Wrap(err, "authUseCase.Login.FindByName")
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, errors.Wrap(err, "authUseCase.Login.ComparePasswords")
	}

	foundUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(foundUser.ID, u.cfg.JwtExpires, u.cfg.JwtSecretKey)

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}
