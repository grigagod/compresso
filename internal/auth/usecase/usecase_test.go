package usecase_test

import (
	"database/sql"
	"testing"

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

	cfg := config.Auth{}

	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := usecase.NewAuthUseCase(&cfg, mockAuthRepo)

	user := &models.User{
		Username: "aa",
		Password: "asd",
	}

	t.Run("Main case", func(t *testing.T) {
		mockAuthRepo.EXPECT().FindByName(gomock.Eq(user.Username)).Return(nil, sql.ErrNoRows)
		mockAuthRepo.EXPECT().Create(gomock.Eq(user)).Return(user, nil)

		createdUser, err := authUC.Register(user)

		require.NoError(t, err)
		require.NotNil(t, createdUser)

	})
	t.Run("Already existst", func(t *testing.T) {
		mockAuthRepo.EXPECT().FindByName(gomock.Eq(user.Username)).Return(user, nil)

		createdUser, err := authUC.Register(user)

		require.Error(t, err)
		require.Nil(t, createdUser)
	})
}

func TestAuthUseCase_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Auth{}

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

	t.Run("Main case", func(t *testing.T) {
		mockAuthRepo.EXPECT().FindByName(gomock.Eq(user.Username)).Return(mockUser, nil)

		userWithToken, err := authUC.Login(user)

		require.NoError(t, err)
		require.NotNil(t, userWithToken)

	})
	t.Run("User not exist", func(t *testing.T) {
		mockAuthRepo.EXPECT().FindByName(gomock.Eq(user.Username)).Return(nil, sql.ErrNoRows)

		userWithToken, err := authUC.Login(user)

		require.Nil(t, userWithToken)
		require.Error(t, err)
	})
	t.Run("Wrong password", func(t *testing.T) {
		mockAuthRepo.EXPECT().FindByName(gomock.Eq(user.Username)).Return(user, nil)

		userWithToken, err := authUC.Login(user)

		require.Nil(t, userWithToken)
		require.Error(t, err)

	})
}
