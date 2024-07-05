package dto

import (
	"time"

	"github.com/google/uuid"
)

type MessageDTO struct {
	SenderId    uuid.UUID `json:"sender_id"`
	CreatedAt   time.Time `json:"created_at"`
	MessageBody string    `json:"message_body"`
}
