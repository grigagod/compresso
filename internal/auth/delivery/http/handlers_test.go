package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/auth"
	"github.com/grigagod/compresso/internal/auth/mock"
	"github.com/grigagod/compresso/internal/httper"
	"github.com/grigagod/compresso/internal/utils"
	"github.com/stretchr/testify/require"
)

// Common objects for testing.
var (
	user = &auth.User{
		Username: "test",
		Password: "test",
	}
	authBody      = []byte(`{"username": "test", "password": "test"}`)
	userWithToken = &auth.UserWithToken{
		User: &auth.User{
			ID:        uuid.New(),
			Username:  user.Username,
			Password:  "",
			CreatedAt: time.Now(),
		},
	}
)

func TestAuthHandlers_Register(t *testing.T) {
	testCases := []struct {
		name                string
		reqBody             []byte
		mockExpect          func(uc *mock.MockUseCase)
		expectedStatusCode  int
		expectedContentType string
		requirements        func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:    "Main case",
			reqBody: authBody,
			mockExpect: func(uc *mock.MockUseCase) {
				uc.EXPECT().Register(gomock.Any(), gomock.Eq(user)).Return(userWithToken, nil)
			},
			expectedStatusCode:  http.StatusCreated,
			expectedContentType: utils.JSONContentType,
			requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
				resp := new(auth.UserWithToken)
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), resp))
				require.EqualValues(t, "", resp.Password)
			},
		},
		{
			name:    "User exists",
			reqBody: authBody,
			mockExpect: func(uc *mock.MockUseCase) {
				uc.EXPECT().Register(gomock.Any(), gomock.Eq(user)).Return(nil, httper.NewBadRequestMsg(httper.UserExistsMsg))
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: utils.TextContentType,
			requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, rec.Body.String(), string(httper.UserExistsMsg))
			},
		},
		{
			name:                "Validate password",
			reqBody:             []byte(`{"username": "test", "password": "1"}`),
			mockExpect:          func(uc *mock.MockUseCase) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: utils.TextContentType,
			requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, rec.Body.String(), string(httper.InvalidPasswordMsg))
			},
		},
		{
			name:                "Validate username",
			reqBody:             []byte(`{"username": "t", "password": "1234"}`),
			mockExpect:          func(uc *mock.MockUseCase) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: utils.TextContentType,
			requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, rec.Body.String(), string(httper.InvalidUsernameMsg))
			},
		},
		{
			name:                "Invalid request body",
			reqBody:             []byte(""),
			mockExpect:          func(uc *mock.MockUseCase) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: utils.TextContentType,
			requirements:        func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthUC := mock.NewMockUseCase(ctrl)
			tc.mockExpect(mockAuthUC)

			authHandlers := NewAuthHandlers(mockAuthUC)

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(tc.reqBody))
			req.Header.Set("Content-Type", utils.JSONContentType)
			rec := httptest.NewRecorder()

			mux := chi.NewMux()
			mux.Method("POST", "/register", authHandlers.Register())
			mux.ServeHTTP(rec, req)

			require.Equal(t, tc.expectedStatusCode, rec.Code)
			require.Equal(t, tc.expectedContentType, rec.Header().Get("Content-Type"))
			tc.requirements(t, rec)
		})
	}
}

func TestAuthHandlers_Login(t *testing.T) {
	testCases := []struct {
		name                string
		reqBody             []byte
		mockExpect          func(uc *mock.MockUseCase)
		expectedStatusCode  int
		expectedContentType string
		requirements        func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:    "Main case",
			reqBody: authBody,
			mockExpect: func(uc *mock.MockUseCase) {
				uc.EXPECT().Login(gomock.Any(), gomock.Eq(user)).Return(userWithToken, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedContentType: utils.JSONContentType,
			requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
				resp := new(auth.UserWithToken)
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), resp))
				require.EqualValues(t, "", resp.Password)
			},
		},

		{
			name:    "Wrong credentials",
			reqBody: authBody,
			mockExpect: func(uc *mock.MockUseCase) {
				uc.EXPECT().Login(gomock.Any(), gomock.Eq(user)).Return(nil, httper.NewWrongCredentialsMsg())
			},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedContentType: utils.TextContentType,
			requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, rec.Body.String(), string(httper.WrongCredentialsMsg))
			},
		},
		{
			name:                "Validate password",
			reqBody:             []byte(`{"username": "test", "password": "1"}`),
			mockExpect:          func(uc *mock.MockUseCase) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: utils.TextContentType,
			requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, rec.Body.String(), string(httper.InvalidPasswordMsg))
			},
		},
		{
			name:                "Validate username",
			reqBody:             []byte(`{"username": "t", "password": "1234"}`),
			mockExpect:          func(uc *mock.MockUseCase) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: utils.TextContentType,
			requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, rec.Body.String(), string(httper.InvalidUsernameMsg))
			},
		},
		{
			name:                "Invalid request body",
			reqBody:             []byte(""),
			mockExpect:          func(uc *mock.MockUseCase) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: utils.TextContentType,
			requirements:        func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthUC := mock.NewMockUseCase(ctrl)
			tc.mockExpect(mockAuthUC)

			authHandlers := NewAuthHandlers(mockAuthUC)

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(tc.reqBody))
			req.Header.Set("Content-Type", utils.JSONContentType)
			rec := httptest.NewRecorder()

			mux := chi.NewMux()
			mux.Method("POST", "/login", authHandlers.Login())
			mux.ServeHTTP(rec, req)

			require.Equal(t, tc.expectedStatusCode, rec.Code)
			require.Equal(t, tc.expectedContentType, rec.Header().Get("Content-Type"))
			tc.requirements(t, rec)
		})
	}
}
