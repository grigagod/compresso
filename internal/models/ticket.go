package models

import (
	"time"

	"github.com/google/uuid"
)

// ProcessingState represent tiket's processing state.
type ProcessingState string

// Package defined processing states.
const (
	Queued     ProcessingState = "queued"
	Processing ProcessingState = "processing"
	Done       ProcessingState = "done"
	Failed     ProcessingState = "failed"
)

// Ticket is a skeleton for all kinds of media tickets.
type Ticket struct {
	ID        uuid.UUID       `db:"ticket_id" json:"ticket_id"`
	AuthorID  uuid.UUID       `db:"author_id" json:"author_id"`
	State     ProcessingState `db:"state" json:"state"`
	URL       string          `db:"url" json:"url"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
}
