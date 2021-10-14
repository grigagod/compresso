package models

import (
	"github.com/google/uuid"
	"github.com/grigagod/compresso/pkg/converter"
)

type QueueVideoMsg struct {
	TicketID     uuid.UUID
	CRF          int
	TargetFormat converter.VideoFormat
	OriginURL    string
	ProcessedURL string
}
