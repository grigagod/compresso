//go:build unit

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
)

func TestAPIHandlers(t *testing.T) {
	type testCase struct {
		name                string
		reqPath             string
		reqBody             []byte
		reqContentType      string
		reqUserID           uuid.UUID
		mockExpect          func(uc *mock.MockUseCase)
		expectedStatusCode  int
		expectedContentType string
		requirements        func(t *testing.T, rec *httptest.ResponseRecorder)
	}
	testTable := []struct {
		handler  string
		method   string
		basePath string
		cases    []testCase
	}{
		{
			handler:  "Create video",
			method:   "POST",
			basePath: "/videos",
			cases: []testCase{
				{
					name:           "main case",
					reqContentType: videoMIMEType,
					reqUserID:      video.AuthorID,
					mockExpect: func(uc *mock.MockUseCase) {
						uc.EXPECT().CreateVideo(gomock.Any(), gomock.Any(), gomock.Any()).Return(video, nil)
					},
					expectedStatusCode:  http.StatusCreated,
					expectedContentType: utils.JSONContentType,
					requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
						resp := new(models.Video)
						require.NoError(t, json.Unmarshal(rec.Body.Bytes(), resp))
						require.Equal(t, video, resp)

					},
				},
				{
					name:                "no user ID in context",
					reqUserID:           uuid.UUID{},
					mockExpect:          func(uc *mock.MockUseCase) {},
					expectedStatusCode:  http.StatusUnauthorized,
					expectedContentType: utils.TextContentType,
					requirements:        func(t *testing.T, rec *httptest.ResponseRecorder) {},
				},
				{
					name:                "no content type in context",
					reqUserID:           video.AuthorID,
					mockExpect:          func(uc *mock.MockUseCase) {},
					expectedStatusCode:  http.StatusBadRequest,
					expectedContentType: utils.TextContentType,
					requirements:        func(t *testing.T, rec *httptest.ResponseRecorder) {},
				},
			},
		},
		{
			handler:  "Create Ticket",
			method:   "POST",
			basePath: "/tickets",
			cases: []testCase{
				{
					name:      "main case",
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
					name:                "no user ID in context",
					reqUserID:           uuid.UUID{},
					mockExpect:          func(uc *mock.MockUseCase) {},
					expectedStatusCode:  http.StatusUnauthorized,
					expectedContentType: utils.TextContentType,
					requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
						require.Equal(t, string(httper.WrongCredentialsMsg), rec.Body.String())
					},
				},
				{
					name:                "validation error",
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
			},
		},
		{
			handler:  "Get video by ID",
			method:   "GET",
			basePath: "/videos/",
			cases: []testCase{
				{
					name:      "Get video by ID main case",
					reqPath:   video.ID.String(),
					reqUserID: video.AuthorID,
					mockExpect: func(uc *mock.MockUseCase) {
						uc.EXPECT().GetVideoByID(gomock.Any(), video.AuthorID, video.ID).Return(video, nil)
					},
					expectedStatusCode:  http.StatusOK,
					expectedContentType: utils.JSONContentType,
					requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
						resp := new(models.Video)
						require.NoError(t, json.Unmarshal(rec.Body.Bytes(), resp))
					},
				},
				{
					name:                "Get video by ID no user ID in context",
					reqPath:             video.ID.String(),
					reqUserID:           uuid.UUID{},
					mockExpect:          func(uc *mock.MockUseCase) {},
					expectedStatusCode:  http.StatusUnauthorized,
					expectedContentType: utils.TextContentType,
					requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
						require.Equal(t, string(httper.WrongCredentialsMsg), rec.Body.String())
					},
				},
				{
					name:                "Get video by ID bad ID ",
					reqPath:             "video_id",
					reqUserID:           video.AuthorID,
					mockExpect:          func(uc *mock.MockUseCase) {},
					expectedStatusCode:  http.StatusBadRequest,
					expectedContentType: utils.TextContentType,
					requirements:        func(t *testing.T, rec *httptest.ResponseRecorder) {},
				},
				{
					name:      "Get video by ID not found",
					reqPath:   video.ID.String(),
					reqUserID: video.AuthorID,
					mockExpect: func(uc *mock.MockUseCase) {
						uc.EXPECT().GetVideoByID(gomock.Any(), video.AuthorID, video.ID).Return(nil, httper.NewNotFoundMsg())
					},
					expectedStatusCode:  http.StatusNotFound,
					expectedContentType: utils.TextContentType,
					requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
						require.Equal(t, string(httper.NotFoundMsg), rec.Body.String())
					},
				},
			},
		},
		{
			handler:  "Get ticket by ID",
			method:   "GET",
			basePath: "/tickets/",
			cases: []testCase{
				{
					name:      "Get ticket by ID main case",
					reqPath:   ticket.ID.String(),
					reqUserID: ticket.AuthorID,
					mockExpect: func(uc *mock.MockUseCase) {
						uc.EXPECT().GetTicketByID(gomock.Any(), ticket.AuthorID, ticket.ID).Return(ticket, nil)
					},
					expectedStatusCode:  http.StatusOK,
					expectedContentType: utils.JSONContentType,
					requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
						resp := new(models.VideoTicket)
						require.NoError(t, json.Unmarshal(rec.Body.Bytes(), resp))
					},
				},
				{
					name:                "Get ticket by ID no user ID in context",
					reqPath:             ticket.ID.String(),
					reqUserID:           uuid.UUID{},
					mockExpect:          func(uc *mock.MockUseCase) {},
					expectedStatusCode:  http.StatusUnauthorized,
					expectedContentType: utils.TextContentType,
					requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
						require.Equal(t, string(httper.WrongCredentialsMsg), rec.Body.String())
					},
				},
				{
					name:                "Get ticket by ID bad ID",
					reqPath:             "ticket_id",
					reqUserID:           ticket.AuthorID,
					mockExpect:          func(uc *mock.MockUseCase) {},
					expectedStatusCode:  http.StatusBadRequest,
					expectedContentType: utils.TextContentType,
					requirements:        func(t *testing.T, rec *httptest.ResponseRecorder) {},
				},
				{
					name:      "Get ticket by ID not found",
					reqPath:   ticket.ID.String(),
					reqUserID: ticket.AuthorID,
					mockExpect: func(uc *mock.MockUseCase) {
						uc.EXPECT().GetTicketByID(gomock.Any(), ticket.AuthorID, ticket.ID).Return(nil, httper.NewNotFoundMsg())
					},
					expectedStatusCode:  http.StatusNotFound,
					expectedContentType: utils.TextContentType,
					requirements: func(t *testing.T, rec *httptest.ResponseRecorder) {
						require.Equal(t, string(httper.NotFoundMsg), rec.Body.String())
					},
				},
			},
		},
	}

	for _, h := range testTable {
		t.Run(h.handler, func(t *testing.T) {
			for _, tc := range h.cases {
				t.Run(tc.name, func(t *testing.T) {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()

					mockAPIUC := mock.NewMockUseCase(ctrl)

					apiHandlers := NewAPIHandlers(mockAPIUC)
					tc.mockExpect(mockAPIUC)
					req := httptest.NewRequest(h.method, h.basePath+tc.reqPath, bytes.NewReader(tc.reqBody))
					if !reflect.DeepEqual(tc.reqUserID, uuid.UUID{}) {
						req = req.WithContext(utils.ContextWithUserID(req.Context(), tc.reqUserID.String()))
					}
					if tc.reqContentType != "" {
						req = req.WithContext(utils.ContextWithContentType(req.Context(), tc.reqContentType))
					}
					rec := httptest.NewRecorder()

					mux := chi.NewMux()
					MapVideoRoutes(mux, apiHandlers)
					mux.ServeHTTP(rec, req)

					require.Equal(t, tc.expectedStatusCode, rec.Code)
					tc.requirements(t, rec)
				})
			}
		})
	}
}
