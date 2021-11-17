//go:build unit

package repo

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/grigagod/compresso/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func NewVideoRows(t *testing.T, video *models.Video) *sqlmock.Rows {
	t.Helper()
	return sqlmock.NewRows([]string{"video_id", "author_id", "format", "url", "created_at"}).
		AddRow(&video.ID, &video.AuthorID, &video.Format, &video.URL, &video.CreatedAt)
}

func NewTicketRows(t *testing.T, ticket *models.VideoTicket) *sqlmock.Rows {
	t.Helper()
	return sqlmock.NewRows([]string{"ticket_id", "video_id", "author_id", "target_format",
		"state", "url", "created_at"}).
		AddRow(&ticket.ID, &ticket.VideoID, &ticket.AuthorID, &ticket.TargetFormat,
			&ticket.State, &ticket.URL, &ticket.CreatedAt)
}

func TestVideoRepo_MainCases(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	videoRepo := NewVideoRepository(sqlxDB)

	video := &models.Video{}

	ticket := &models.VideoTicket{}

	t.Run("Insert video", func(t *testing.T) {
		rows := NewVideoRows(t, video)
		query := `INSERT INTO
        svc.user_videos(video_id, author_id, format, url)
        VALUES($1, $2, $3, $4) RETURNING *`
		mock.ExpectQuery(query).
			WithArgs(&video.ID, &video.AuthorID, &video.Format, &video.URL).
			WillReturnRows(rows)

		createdVideo, err := videoRepo.InsertVideo(context.Background(), video)
		require.NoError(t, err)
		require.Equal(t, video, createdVideo)

	})
	t.Run("Insert ticket", func(t *testing.T) {
		rows := NewTicketRows(t, ticket)
		query := `INSERT INTO
        svc.video_tickets(ticket_id, video_id, author_id, target_format, state, crf, url)
        VALUES($1, $2, $3, $4, $5, $6, COALESCE(NULLIF($7, ''), $7)) RETURNING *`
		mock.ExpectQuery(query).
			WithArgs(&ticket.ID, &ticket.VideoID, &ticket.AuthorID, &ticket.TargetFormat,
				&ticket.State, &ticket.CRF, &ticket.URL).
			WillReturnRows(rows)

		createdTicket, err := videoRepo.InsertTicket(context.Background(), ticket)
		require.NoError(t, err)
		require.Equal(t, ticket, createdTicket)
	})
	t.Run("Select video by ID", func(t *testing.T) {
		rows := NewVideoRows(t, video)
		query := `SELECT * FROM svc.user_videos WHERE author_id = $1 AND video_id = $2`
		mock.ExpectQuery(query).WithArgs(&video.AuthorID, &video.ID).WillReturnRows(rows)

		gotVideo, err := videoRepo.SelectVideoByID(context.Background(), video.AuthorID, video.ID)
		require.NoError(t, err)
		require.Equal(t, gotVideo, video)
	})
	t.Run("Select videos by author ID", func(t *testing.T) {
		rows := NewVideoRows(t, video)
		query := `SELECT * FROM svc.user_videos WHERE author_id = $1`
		mock.ExpectQuery(query).WithArgs(&video.AuthorID).WillReturnRows(rows)

		gotVideos, err := videoRepo.SelectVideos(context.Background(), video.AuthorID)
		require.NoError(t, err)
		require.NotNil(t, gotVideos)
		require.Equal(t, 1, len(gotVideos))
	})
	t.Run("Select ticket by ID", func(t *testing.T) {
		rows := NewTicketRows(t, ticket)
		query := `SELECT * FROM svc.video_tickets WHERE author_id = $1 AND ticket_id = $2`
		mock.ExpectQuery(query).WithArgs(&ticket.AuthorID, &ticket.ID).WillReturnRows(rows)

		gotTicket, err := videoRepo.SelectTicketByID(context.Background(), ticket.AuthorID, ticket.ID)
		require.NoError(t, err)
		require.Equal(t, gotTicket, ticket)
	})
	t.Run("Select tickets by author ID", func(t *testing.T) {
		rows := NewTicketRows(t, ticket)
		query := `SELECT * FROM svc.video_tickets WHERE author_id = $1`
		mock.ExpectQuery(query).WithArgs(&ticket.AuthorID).WillReturnRows(rows)

		gotTickets, err := videoRepo.SelectTickets(context.Background(), ticket.AuthorID)
		require.NoError(t, err)
		require.NotNil(t, gotTickets)
		require.Equal(t, 1, len(gotTickets))
	})
}
