package models

import (
	"github.com/google/uuid"
	"github.com/grigagod/compresso/pkg/converter"
)

// QueueVideoMsg represent message that's used in broker for async video processing.
type QueueVideoMsg struct {
	TicketID     uuid.UUID
	CRF          int
	TargetFormat converter.VideoFormat
	OriginURL    string
	ProcessedURL string
}
