package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

type (
	Msg struct {
		ID        uuid.UUID `json:"id"`
		Text      string    `json:"text"`
		State     State     `json:"state"`
		CreatedAt time.Time `json:"created_at"`
	}

	State string
)

const (
	StateUnprocessed State = "UNPROCESSED"
	StateProcessing  State = "PROCESSING"
	StateProcessed   State = "PROCESSED"
)
