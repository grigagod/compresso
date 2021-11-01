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

// CreateVideo insert video model in db.
func (r *VideoRepo) CreateVideo(ctx context.Context, video *models.Video) (*models.Video, error) {
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

// CreateTicket insert video ticket model in DB.
func (r *VideoRepo) CreateTicket(ctx context.Context, ticket *models.VideoTicket) (*models.VideoTicket, error) {
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

// GetVideoByID select by ID video model form DB.
func (r *VideoRepo) GetVideoByID(ctx context.Context, authorID uuid.UUID, id uuid.UUID) (*models.Video, error) {
	query := `SELECT * FROM svc.user_videos WHERE author_id = $1 AND video_id = $2`

	var v models.Video

	err := r.db.QueryRowxContext(ctx, query, authorID, id).StructScan(&v)
	if err != nil {
		return nil, errors.Wrap(err, "VideoRepo.GetVideoByID.StructScan")
	}

	return &v, nil
}

// GetTicketByID select by ID video ticket model from DB.
func (r *VideoRepo) GetTicketByID(ctx context.Context, authorID, id uuid.UUID) (*models.VideoTicket, error) {
	query := `SELECT * FROM svc.video_tickets WHERE authorID =$1 AND ticket_id = $2`

	var t models.VideoTicket

	err := r.db.QueryRowxContext(ctx, query, authorID, id).StructScan(&t)
	if err != nil {
		return nil, errors.Wrap(err, "VideoRepo.GetTicketByID.StructScan")
	}

	return &t, nil
}

// GetVideos select videos by author ID.
func (r *VideoRepo) GetVideos(ctx context.Context, authorID uuid.UUID) ([]*models.Video, error) {
	query := `SELECT * FROM svc.user_videos WHERE author_id = $1`
	var videos []*models.Video

	err := r.db.SelectContext(ctx, videos, query, authorID)
	if err != nil {
		return nil, errors.Wrap(err, "VideoRepo.GetVideos.SelectContext")
	}

	return videos, nil
}

// GetTickets select tickets by author ID.
func (r *VideoRepo) GetTickets(ctx context.Context, authorID uuid.UUID) ([]*models.VideoTicket, error) {
	query := `SELECT * FROM svc.video_tickets WHERE author_id = $1`
	var tickets []*models.VideoTicket

	err := r.db.SelectContext(ctx, tickets, query, authorID)
	if err != nil {
		return nil, errors.Wrap(err, "VideoRepo.GetTickets.SelectContext")
	}

	return tickets, nil
}
