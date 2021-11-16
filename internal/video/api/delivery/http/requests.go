package http

import (
	"github.com/google/uuid"
)

type CreateTicketRequest struct {
	VideoID      uuid.UUID `json:"video_id" validate:"required,gt=0"`
	TargetFormat string    `json:"target_format" validate:"required,gt=0"`
	CRF          int       `json:"crf" validate:"required,gt=0,lt=51"`
}
