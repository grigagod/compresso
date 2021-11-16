package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/httper"
	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/internal/utils"
	"github.com/grigagod/compresso/internal/video/api/mock"
	"github.com/grigagod/compresso/pkg/converter"
	"github.com/stretchr/testify/require"
)

var (
	video = &models.Video{
		ID:       uuid.New(),
		AuthorID: uuid.New(),
		Format:   converter.MKV,
	}
	videoMIMEType = "video/x-matroska"
	ticket        = &models.VideoTicket{
		Ticket: models.Ticket{
			ID:       uuid.New(),
			AuthorID: video.AuthorID,
		},
		TargetFormat: converter.WebM,
		VideoID:      video.ID,
		CRF:          20,
	}
	emptyBody = []byte{}
)

func TestAPIHandlers_CreateVideo(t *testing.T) {
	testCases := []struct {
		name                string
		reqBody             []byte
		reqUserID           uuid.UUID
		reqContentType      string
		mockExpect          func(uc *mock.MockUseCase)
		expectedStatusCode  int
		expectedContentType string
		requirements        func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:           "Main case",
			reqBody:        emptyBody,
			reqUserID:      video.AuthorID,
			reqContentType: videoMIMEType,
			mockExpect: func(uc *mock.MockUseCase) {
				uc.EXPECT().CreateVideo(gomock.Any(), gomock.Any(), gomock.Any()).Return(video, nil)
			},
			expectedStatusCode:  http.StatusCreated,
			expectedContentType: utils.JSONContentType,
			requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {

			},
		},
		{
			name:                "No user ID in context",
			reqBody:             emptyBody,
			reqUserID:           uuid.UUID{},
			mockExpect:          func(uc *mock.MockUseCase) {},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedContentType: utils.TextContentType,
			requirements:        func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
		{
			name:                "No content type in context",
			reqBody:             emptyBody,
			reqUserID:           video.AuthorID,
			reqContentType:      "",
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

			mockAPIUC := mock.NewMockUseCase(ctrl)

			apiHandlers := NewAPIHandlers(mockAPIUC)
			tc.mockExpect(mockAPIUC)

			req := httptest.NewRequest("POST", "/videos", bytes.NewReader(tc.reqBody))
			if !reflect.DeepEqual(tc.reqUserID, uuid.UUID{}) {
				req = req.WithContext(utils.ContextWithUserID(req.Context(), video.AuthorID.String()))
			}
			if tc.reqContentType != "" {
				req = req.WithContext(utils.ContextWithContentType(req.Context(), tc.reqContentType))
			}
			rec := httptest.NewRecorder()

			mux := chi.NewMux()
			mux.Method("POST", "/videos", apiHandlers.CreateVideo())
			mux.ServeHTTP(rec, req)

			require.Equal(t, tc.expectedStatusCode, rec.Code)
			require.Equal(t, tc.expectedContentType, rec.Header().Get("Content-Type"))
			tc.requirements(t, rec)
		})
	}
}

func TestAPIHandlers_CreateTicket(t *testing.T) {
	testCases := []struct {
		name                string
		reqBody             []byte
		reqUserID           uuid.UUID
		mockExpect          func(uc *mock.MockUseCase)
		expectedStatusCode  int
		expectedContentType string
		requirements        func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:      "Main case",
			reqBody:   []byte(`{"target_format": "webm", "crf": 20, "video_id": "` + video.ID.String() + `"}`),
			reqUserID: video.AuthorID,
			mockExpect: func(uc *mock.MockUseCase) {
				uc.EXPECT().CreateTicket(gomock.Any(), gomock.Any()).Return(ticket, nil)
			},
			expectedStatusCode:  http.StatusCreated,
			expectedContentType: utils.JSONContentType,
			requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
				resp := new(models.VideoTicket)
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), resp))
				require.Equal(t, ticket, resp)
			},
		},
		{
			name:                "No user ID in context",
			reqBody:             emptyBody,
			reqUserID:           uuid.UUID{},
			mockExpect:          func(uc *mock.MockUseCase) {},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedContentType: utils.TextContentType,
			requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, string(httper.WrongCredentialsMsg), rec.Body.String())
			},
		},
		{
			name:                "Validation error",
			reqBody:             []byte(`"target_format: "webm", "crf": 12, "video_id": "test"`),
			reqUserID:           video.AuthorID,
			mockExpect:          func(uc *mock.MockUseCase) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: utils.TextContentType,
			requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
			},
		},
		{
			name:                "CRF out of range",
			reqBody:             []byte(`{"target_format": "webm", "crf": 100, "video_id": "` + video.ID.String() + `"}`),
			reqUserID:           ticket.AuthorID,
			mockExpect:          func(uc *mock.MockUseCase) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: utils.TextContentType,
			requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, converter.ErrVideoCRF.Error(), rec.Body.String())
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAPIUC := mock.NewMockUseCase(ctrl)

			apiHandlers := NewAPIHandlers(mockAPIUC)
			tc.mockExpect(mockAPIUC)

			req := httptest.NewRequest("POST", "/tickets", bytes.NewReader(tc.reqBody))
			if !reflect.DeepEqual(tc.reqUserID, uuid.UUID{}) {
				req = req.WithContext(utils.ContextWithUserID(req.Context(), video.AuthorID.String()))
			}
			rec := httptest.NewRecorder()

			mux := chi.NewMux()
			mux.Method("POST", "/tickets", apiHandlers.CreateTicket())
			mux.ServeHTTP(rec, req)

			require.Equal(t, tc.expectedStatusCode, rec.Code)
			require.Equal(t, tc.expectedContentType, rec.Header().Get("Content-Type"))
			tc.requirements(t, rec)
		})
	}

}
