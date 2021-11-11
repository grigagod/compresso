package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/auth"
	"github.com/grigagod/compresso/internal/auth/mock"
	"github.com/grigagod/compresso/internal/httper"
	"github.com/stretchr/testify/require"
)

func TestAuthHandlers_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUC := mock.NewMockUseCase(ctrl)

	authHandlers := NewAuthHandlers(mockAuthUC)

	t.Run("Main case", func(t *testing.T) {
		user := &auth.User{
			Username: "test",
			Password: "test",
		}

		userWithToken := &auth.UserWithToken{User: &auth.User{
			ID:        uuid.New(),
			Username:  user.Username,
			CreatedAt: time.Now(),
		}, Token: "token"}

		body, err := json.Marshal(user)
		require.NoError(t, err)
		require.NotNil(t, body)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mockAuthUC.EXPECT().Register(req.Context(), gomock.Eq(user)).Return(userWithToken, nil)

		handlerFunc := authHandlers.Register()
		handlerFunc.ServeHTTP(rec, req)

		require.Equal(t, rec.Code, http.StatusCreated)
		require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})
	t.Run("Invalid password", func(t *testing.T) {
		user := &auth.User{
			Username: "test",
			Password: "t",
		}

		body, err := json.Marshal(user)
		require.NoError(t, err)
		require.NotNil(t, body)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handlerFunc := authHandlers.Register()
		handlerFunc.ServeHTTP(rec, req)

		require.Equal(t, rec.Code, http.StatusBadRequest)
		require.Equal(t, "text/plain; charset=utf-8", rec.Header().Get("Content-Type"))
		require.EqualValues(t, rec.Body.String(), string(httper.InvalidPasswordMsg))
	})
}

func TestAuthHandlers_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUC := mock.NewMockUseCase(ctrl)

	authHandlers := NewAuthHandlers(mockAuthUC)
	t.Run("Main case", func(t *testing.T) {
		user := &auth.User{
			Username: "test",
			Password: "test",
		}
		userWithToken := &auth.UserWithToken{User: &auth.User{
			ID:        uuid.New(),
			Username:  user.Username,
			CreatedAt: time.Now(),
		}, Token: "token"}

		body, err := json.Marshal(user)
		require.NoError(t, err)
		require.NotNil(t, body)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mockAuthUC.EXPECT().Login(req.Context(), gomock.Eq(user)).Return(userWithToken, nil)

		handlerFunc := authHandlers.Login()
		handlerFunc.ServeHTTP(rec, req)

		require.Equal(t, rec.Code, http.StatusOK)
		require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})
	t.Run("Wrong credentials", func(t *testing.T) {
		user := &auth.User{
			Username: "test",
			Password: "test2",
		}

		body, err := json.Marshal(user)
		require.NoError(t, err)
		require.NotNil(t, body)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mockAuthUC.EXPECT().Login(req.Context(), gomock.Eq(user)).Return(nil, httper.NewWrongCredentialsMsg())

		handlerFunc := authHandlers.Login()
		handlerFunc.ServeHTTP(rec, req)

		require.Equal(t, rec.Code, http.StatusUnauthorized)
		require.Equal(t, "text/plain; charset=utf-8", rec.Header().Get("Content-Type"))
		require.EqualValues(t, rec.Body.String(), string(httper.WrongCredentialsMsg))
	})
}
