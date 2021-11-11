package http

import (
	"github.com/google/uuid"
)

type CreateTicketRequest struct {
	VideoID      uuid.UUID `json:"video_id" validate:"required"`
	TargetFormat string    `json:"target_format" validate:"required"`
	CRF          int       `json:"crf" validate:"required"`
}
