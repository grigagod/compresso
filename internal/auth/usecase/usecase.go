package usecase

import (
	"context"

	"github.com/grigagod/compresso/internal/auth"
	"github.com/grigagod/compresso/internal/auth/config"
	"github.com/grigagod/compresso/internal/httper"
	"github.com/grigagod/compresso/pkg/utils"
)

type AuthUseCase struct {
	cfg      *config.Auth
	authRepo auth.Repository
}

func NewAuthUseCase(cfg *config.Auth, authRepo auth.Repository) *AuthUseCase {
	return &AuthUseCase{cfg: cfg, authRepo: authRepo}
}

func (u *AuthUseCase) Register(ctx context.Context, user *auth.User) (*auth.UserWithToken, error) {
	existsUser, err := u.authRepo.FindByName(ctx, user.Username)
	if existsUser != nil || err == nil {
		return nil, httper.NewBadRequestMsg(httper.UserExistsMsg)
	}

	if err = user.HashPassword(); err != nil {
		return nil, err
	}

	createdUser, err := u.authRepo.Create(ctx, user)
	if err != nil {
		return nil, httper.ParseSqlError(err)
	}

	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser.ID, u.cfg.JwtExpires, u.cfg.JwtSecretKey)
	if err != nil {
		return nil, err
	}

	return &auth.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

func (u *AuthUseCase) Login(ctx context.Context, user *auth.User) (*auth.UserWithToken, error) {
	foundUser, err := u.authRepo.FindByName(ctx, user.Username)
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

	return &auth.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}
