package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/internal/utils"
	"github.com/grigagod/compresso/internal/video/api/mock"
	"github.com/grigagod/compresso/pkg/converter"
	"github.com/stretchr/testify/require"
)

func TestAPIHandlers_CreateVideo(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAPIUC := mock.NewMockUseCase(ctrl)

	apiHandlers := NewAPIHandlers(mockAPIUC)

	body := []byte("")

	t.Run("Main case", func(t *testing.T) {
		video := &models.Video{
			AuthorID: uuid.New(),
			Format:   converter.MKV,
		}
		mimeType, ok := utils.DetectVideoMIMEType(video.Format)
		require.True(t, ok)

		req := httptest.NewRequest("POST", "/videos", bytes.NewReader(body))
		req = req.WithContext(utils.ContextWithUserID(req.Context(), video.AuthorID.String()))
		req = req.WithContext(utils.ContextWithContentType(req.Context(), mimeType))
		rec := httptest.NewRecorder()

		mockAPIUC.EXPECT().CreateVideo(gomock.Any(), gomock.Any(), gomock.Any()).Return(video, nil)

		mux := chi.NewMux()
		mux.Method("POST", "/videos", apiHandlers.CreateVideo())

		mux.ServeHTTP(rec, req)

		require.Equal(t, rec.Code, http.StatusCreated)
		require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})
	t.Run("No token in context", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/videos", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		mux := chi.NewMux()
		mux.Method("POST", "/videos", apiHandlers.CreateVideo())

		mux.ServeHTTP(rec, req)

		require.Equal(t, rec.Code, http.StatusUnauthorized)
	})
	t.Run("No content type in context", func(t *testing.T) {

	})
	t.Run("Wrong content type in context", func(t *testing.T) {

	})
}
