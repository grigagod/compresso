package usecase

import (
	"github.com/grigagod/compresso/internal/auth"
	"github.com/grigagod/compresso/internal/auth/config"
	"github.com/grigagod/compresso/internal/httper"
	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/pkg/utils"
)

type authUseCase struct {
	cfg      *config.Auth
	authRepo auth.Repository
}

func NewAuthUseCase(cfg *config.Auth, authRepo auth.Repository) auth.UseCase {
	return &authUseCase{cfg: cfg, authRepo: authRepo}
}

func (u *authUseCase) Register(user *models.User) (*models.UserWithToken, error) {
	existsUser, err := u.authRepo.FindByName(user.Username)
	if existsUser != nil || err == nil {
		return nil, httper.NewBadRequestMsg(httper.UserExistsMsg)
	}

	if err = user.HashPassword(); err != nil {
		return nil, err
	}

	createdUser, err := u.authRepo.Create(user)
	if err != nil {
		return nil, httper.ParseSqlError(err)
	}

	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser.ID, u.cfg.JwtExpires, u.cfg.JwtSecretKey)
	if err != nil {
		return nil, err
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

func (u *authUseCase) Login(user *models.User) (*models.UserWithToken, error) {
	foundUser, err := u.authRepo.FindByName(user.Username)
	if err != nil {
		return nil, httper.ParseSqlError(err)
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, httper.NewBadRequestMsg(httper.WrongCredentialsMsg)
	}

	foundUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(foundUser.ID, u.cfg.JwtExpires, u.cfg.JwtSecretKey)
	if err != nil {
		return nil, err
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}
