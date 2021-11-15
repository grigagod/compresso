//go:build unit

package usecase

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/internal/storage/mock"
	"github.com/grigagod/compresso/internal/utils"
	videoMock "github.com/grigagod/compresso/internal/video/mock"
	"github.com/grigagod/compresso/pkg/converter"
	"github.com/stretchr/testify/require"
)

func TestAPIUseCase_MainCases(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVideoRepo := videoMock.NewMockRepository(ctrl)
	mockVideoPublisher := videoMock.NewMockPublisher(ctrl)
	mockStorage := mock.NewMockStorage(ctrl)
	apiUC := NewAPIUseCase(mockVideoRepo, mockStorage, mockVideoPublisher)

	video := &models.Video{
		ID:        uuid.New(),
		AuthorID:  uuid.New(),
		Format:    converter.MKV,
		CreatedAt: time.Now(),
	}

	file := bytes.NewReader([]byte("it's video"))

	url, err := utils.GenerateURL(video.AuthorID, video.ID)
	require.NoError(t, err)

	video.URL = url

	ticket := &models.VideoTicket{
		Ticket: models.Ticket{
			ID:       uuid.New(),
			AuthorID: video.AuthorID,
			State:    models.Queued,
		},
		VideoID:      video.ID,
		CRF:          10,
		TargetFormat: converter.WebM,
	}

	format, err := utils.DetectVideoMIMEType(video.Format)
	require.NoError(t, err)

	url, err = utils.GenerateURL(ticket.AuthorID, ticket.ID)
	require.NoError(t, err)

	msg := &models.ProcessVideoMsg{
		TicketID:     ticket.ID,
		CRF:          ticket.CRF,
		TargetFormat: ticket.TargetFormat,
		OriginURL:    video.URL,
		ProcessedURL: url,
	}

	t.Run("Create video", func(t *testing.T) {
		mockVideoRepo.EXPECT().InsertVideo(context.Background(), gomock.Eq(video)).Return(video, nil)
		mockStorage.EXPECT().PutObject(context.Background(), file, video.URL, format).Return(nil)
		mockStorage.EXPECT().GetDownloadURL(video.URL).Return(video.URL, nil)

		createdVideo, err := apiUC.CreateVideo(context.Background(), video, file)

		require.NoError(t, err)
		require.NotNil(t, createdVideo)
	})

	t.Run("Create ticket", func(t *testing.T) {
		mockVideoRepo.EXPECT().SelectVideoByID(context.Background(), ticket.AuthorID, ticket.VideoID).Return(video, nil)
		mockVideoRepo.EXPECT().InsertTicket(context.Background(), gomock.Eq(ticket)).Return(ticket, nil)
		mockVideoPublisher.EXPECT().SendMsg(msg).Return(nil)

		createdTicket, err := apiUC.CreateTicket(context.Background(), ticket)

		require.NoError(t, err)
		require.NotNil(t, createdTicket)
	})
	t.Run("Get video by ID", func(t *testing.T) {
		mockVideoRepo.EXPECT().SelectVideoByID(context.Background(), video.AuthorID, video.ID).Return(video, nil)
		mockStorage.EXPECT().GetDownloadURL(video.URL).Return(video.URL, nil)

		gotVideo, err := apiUC.GetVideoByID(context.Background(), video.AuthorID, video.ID)

		require.NoError(t, err)
		require.NotNil(t, gotVideo)
	})
	t.Run("Get ticket by ID", func(t *testing.T) {
		mockVideoRepo.EXPECT().SelectTicketByID(context.Background(), ticket.AuthorID, ticket.ID).Return(ticket, nil)

		gotTicket, err := apiUC.GetTicketByID(context.Background(), ticket.AuthorID, ticket.ID)

		require.NoError(t, err)
		require.NotNil(t, gotTicket)
	})
	t.Run("Get videos by author ID", func(t *testing.T) {
		mockVideoRepo.EXPECT().SelectVideos(context.Background(), video.AuthorID).Return([]*models.Video{video}, nil)
		mockStorage.EXPECT().GetDownloadURL(video.URL).Return(video.URL, nil)

		gotVideos, err := apiUC.GetVideos(context.Background(), video.AuthorID)

		require.NoError(t, err)
		require.NotNil(t, gotVideos)
		require.Equal(t, 1, len(gotVideos))

	})
	t.Run("Get tickets by author ID", func(t *testing.T) {
		mockVideoRepo.EXPECT().SelectTickets(context.Background(), ticket.AuthorID).Return([]*models.VideoTicket{ticket}, nil)

		gotTickets, err := apiUC.GetTickets(context.Background(), ticket.AuthorID)

		require.NoError(t, err)
		require.NotNil(t, gotTickets)
		require.Equal(t, 1, len(gotTickets))

	})
}
