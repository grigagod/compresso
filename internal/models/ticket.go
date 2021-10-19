package models

import (
	"time"

	"github.com/google/uuid"
)

type ProcessingState string

const (
	Queued     ProcessingState = "queued"
	Processing ProcessingState = "processing"
	Done       ProcessingState = "done"
	Failed     ProcessingState = "failed"
)

type Ticket struct {
	ID        uuid.UUID       `db:"ticket_id" json:"ticket_id"`
	AuthorID  uuid.UUID       `db:"author_id" json:"author_id"`
	State     ProcessingState `db:"state" json:"state"`
	URL       string          `db:"url" json:"url"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
}
