package dto

import (
	"time"

	userEntity "github.com/ciscapello/api_gateway/internal/domain/entity/user_entity"
	"github.com/google/uuid"
)

type ConversationMessageDTO struct {
	SenderId    uuid.UUID `json:"sender_id"`
	CreatedAt   time.Time `json:"created_at"`
	MessageBody string    `json:"message_body"`
}

type ConversationsDTO struct {
	ID          int                     `json:"id"`
	User        userEntity.PublicUser   `json:"user"`
	LastMessage *ConversationMessageDTO `json:"last_message,omitempty"`
}
