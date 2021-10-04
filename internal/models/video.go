package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/grigagod/compresso/pkg/converter"
)

type Video struct {
	ID        uuid.UUID             `db:"video_id" json:"video_id"`
	AuthorID  uuid.UUID             `db:"author_id" json:"author_id"`
	Format    converter.VideoFormat `db:"video_format" json:"video_format"`
	URL       string                `db:"url" json:"url"`
	CreatedAt time.Time             `db:"created_at" json:"created_at"`
}

type VideoTicket struct {
	Ticket
	VideoID uuid.UUID             `db:"video_id" json:"video_id"`
	CRF     int                   `db:"crf" json:"crf"`
	Format  converter.VideoFormat `db:"video_format" json:"video_format"`
}
