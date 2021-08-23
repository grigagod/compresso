package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/auth/config"
	"github.com/grigagod/compresso/internal/auth/mock"
	"github.com/grigagod/compresso/internal/auth/models"
	"github.com/stretchr/testify/require"
)

func TestAuthHandlers_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUC := mock.NewMockUseCase(ctrl)

	cfg := &config.Config{
		Auth: config.Auth{
			JwtExpires: time.Second * 10,
		},
	}

	authHandlers := NewAuthHandlers(cfg, mockAuthUC)

	user := &models.User{
		Username: "test",
		Password: "test",
	}

	buf, err := converter.AnyToBytesBuffer(user)
	require.NoError(t, err)
	require.NotNil(t, buf)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(buf.String()))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	userID := uuid.New()

	userWithToken := &models.UserWithToken{
		User: &models.User{
			ID: userID,
		},
	}

	mockAuthUC.EXPECT().Register(gomock.Eq(user)).Return(userWithToken, nil)

	err = authHandlers.Register()

}

func TestAuthHandlers_Login(t *testing.T) {

}
