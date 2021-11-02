package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// VideoRepo implement video.Repository interface.
type VideoRepo struct {
	db *sqlx.DB
}

func NewVideoRepository(db *sqlx.DB) *VideoRepo {
	return &VideoRepo{
		db: db,
	}
}

// InsertVideo insert video model in DB.
func (r *VideoRepo) InsertVideo(ctx context.Context, video *models.Video) (*models.Video, error) {
	query := `INSERT INTO
        svc.user_videos(video_id, author_id, format, url, created_at)
        VALUES($1, $2, $3, $4, now()) RETURNING *`

	var v models.Video

	err := r.db.QueryRowxContext(ctx, query, &video.ID, &video.AuthorID, &video.Format, &video.URL).StructScan(&v)
	if err != nil {
		return nil, errors.Wrap(err, "VideoRepo.Create.StructScan")
	}

	return &v, nil
}

// InsertTicket insert video ticket model in DB.
func (r *VideoRepo) InsertTicket(ctx context.Context, ticket *models.VideoTicket) (*models.VideoTicket, error) {
	query := `INSERT INTO
        svc.video_tickets(ticket_id, video_id, author_id, target_format, state, crf, url, created_at)
        VALUES($1, $2, $3, $4, $5, $6, COALESCE(NULLIF($7, ''), $7), now()) RETURNING *`

	var t models.VideoTicket

	err := r.db.QueryRowxContext(ctx, query, &ticket.ID, &ticket.VideoID, &ticket.AuthorID, &ticket.TargetFormat,
		&ticket.State, &ticket.CRF, &ticket.URL).StructScan(&t)
	if err != nil {
		return nil, errors.Wrap(err, "VideoRepo.Create.StructScan")
	}

	return &t, nil
}

// UpdateTicket update video ticket model in DB.
func (r *VideoRepo) UpdateTicket(ctx context.Context, ticket *models.VideoTicket) (*models.VideoTicket, error) {
	query := `UPDATE svc.video_tickets
        SET target_format = $1,
            state  = $2,
            url = COALESCE(NULLIF($3, ''), $3)
        WHERE ticket_id = $4
        RETURNING *
    `

	var t models.VideoTicket

	err := r.db.GetContext(ctx, &t, query, &ticket.TargetFormat, &ticket.State, &ticket.URL,
		&ticket.ID)
	if err != nil {
		return nil, errors.Wrap(err, "VideoRepo.Update.GetContext")
	}

	return &t, nil
}

// SelectVideoByID select by ID video model form DB.
func (r *VideoRepo) SelectVideoByID(ctx context.Context, authorID uuid.UUID, id uuid.UUID) (*models.Video, error) {
	query := `SELECT * FROM svc.user_videos WHERE author_id = $1 AND video_id = $2`

	var v models.Video

	err := r.db.QueryRowxContext(ctx, query, authorID, id).StructScan(&v)
	if err != nil {
		return nil, errors.Wrap(err, "VideoRepo.SelectVideoByID.StructScan")
	}

	return &v, nil
}

// SelectTicketByID select by ID video ticket model from DB.
func (r *VideoRepo) SelectTicketByID(ctx context.Context, authorID, id uuid.UUID) (*models.VideoTicket, error) {
	query := `SELECT * FROM svc.video_tickets WHERE author_id = $1 AND ticket_id = $2`

	var t models.VideoTicket

	err := r.db.QueryRowxContext(ctx, query, authorID, id).StructScan(&t)
	if err != nil {
		return nil, errors.Wrap(err, "VideoRepo.SelectTicketByID.StructScan")
	}

	return &t, nil
}

// SelectVideos select videos by author ID.
func (r *VideoRepo) SelectVideos(ctx context.Context, authorID uuid.UUID) ([]*models.Video, error) {
	query := `SELECT * FROM svc.user_videos WHERE author_id = $1`
	var videos []*models.Video

	err := r.db.SelectContext(ctx, &videos, query, authorID)
	if err != nil {
		return nil, errors.Wrap(err, "VideoRepo.SelectVideos.SelectContext")
	}

	return videos, nil
}

// SelectTickets select tickets by author ID.
func (r *VideoRepo) SelectTickets(ctx context.Context, authorID uuid.UUID) ([]*models.VideoTicket, error) {
	query := `SELECT * FROM svc.video_tickets WHERE author_id = $1`
	var tickets []*models.VideoTicket

	err := r.db.SelectContext(ctx, &tickets, query, authorID)
	if err != nil {
		return nil, errors.Wrap(err, "VideoRepo.SelectTickets.SelectContext")
	}

	return tickets, nil
}
