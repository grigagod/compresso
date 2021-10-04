package video

import (
	"context"

	"github.com/grigagod/compresso/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type VideoRepo struct {
	db *sqlx.DB
}

func NewVideoRepository(db *sqlx.DB) *VideoRepo {
	return &VideoRepo{
		db: db,
	}
}

func (r *VideoRepo) Create(ctx context.Context, video *models.Video) (*models.Video, error) {
	query := `INSERT INTO
        svc.user_videos(author_id, format, url, created_at)
        VALUES($1, $2, $3, now()) RETURNING *`

	var v models.Video

	err := r.db.QueryRowxContext(ctx, query, &video.AuthorID, &video.Format, &video.URL).StructScan(&v)
	if err != nil {
		return nil, errors.Wrap(err, "VideoRepo.Create.StructScan")
	}

	return &v, nil
}

func (r *VideoRepo) CreateTicket(ctx context.Context, ticket *models.VideoTicket) (*models.VideoTicket, error) {
	query := `INSERT INTO
        svc.video_tickets(video_id, author_id, format, state, crf, url, created_at)
        VALUES($1, $2, $3, $4, $5, $6, COALESCE(NULLIF($8, ''), $8), now()) RETURNING *`

	var t models.VideoTicket

	err := r.db.QueryRowxContext(ctx, query, &ticket.VideoID, &t.AuthorID, &ticket.Format,
		&ticket.State, &ticket.CRF, &ticket.URL).StructScan(&t)
	if err != nil {
		return nil, errors.Wrap(err, "VideoRepo.Create.StructScan")

	}
	return &t, nil
}

func (r *VideoRepo) UpdateTicket(ctx context.Context, ticket *models.VideoTicket) (*models.VideoTicket, error) {
	query := `UPDATE svc.video_tickets
        SET format = $1,
            state  = $2,
            url = COALESCE(NULLIF($3, ''), $3
        WHERE ticket_id = $4 AND author_id = $5`

	var t models.VideoTicket

	err := r.db.GetContext(ctx, &t, query, &ticket.Format, &ticket.State, &ticket.URL,
		&ticket.ID, &ticket.AuthorID)
	if err != nil {
		return nil, errors.Wrap(err, "VideoRepo.Update.GetContext")
	}

	return &t, nil
}
