package dto

import (
	"time"

	"github.com/google/uuid"
)

type MessageUserDTO struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type MessageDTO struct {
	SenderId    uuid.UUID      `json:"sender_id"`
	CreatedAt   time.Time      `json:"created_at"`
	MessageBody string         `json:"message_body"`
	User        MessageUserDTO `json:"user"`
}
