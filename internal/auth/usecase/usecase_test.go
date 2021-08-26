package usecase_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/grigagod/compresso/internal/auth/config"
	"github.com/grigagod/compresso/internal/auth/mock"
	"github.com/grigagod/compresso/internal/auth/usecase"
	"github.com/grigagod/compresso/internal/models"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthUseCase_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		Auth: config.Auth{
			JwtSecretKey: "secret",
			JwtExpires:   time.Minute,
		},
	}

	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := usecase.NewAuthUseCase(&cfg, mockAuthRepo)

	user := &models.User{
		Username: "aa",
		Password: "asd",
	}

	mockAuthRepo.EXPECT().FindByName(gomock.Eq(user.Username)).Return(nil, sql.ErrNoRows)
	mockAuthRepo.EXPECT().Create(gomock.Eq(user)).Return(user, nil)

	createdUser, err := authUC.Register(user)

	require.NoError(t, err)
	require.NotNil(t, createdUser)
}

func TestAuthUseCase_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		Auth: config.Auth{
			JwtSecretKey: "secret",
			JwtExpires:   time.Hour,
		},
	}

	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := usecase.NewAuthUseCase(&cfg, mockAuthRepo)

	user := &models.User{
		Username: "aa",
		Password: "asd",
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	require.NoError(t, err)

	mockUser := &models.User{
		Username: "aa",
		Password: string(hashPassword),
	}

	mockAuthRepo.EXPECT().FindByName(user.Username).Return(mockUser, nil)

	userWithToken, err := authUC.Login(user)
	require.NoError(t, err)
	require.NotNil(t, userWithToken)
}
