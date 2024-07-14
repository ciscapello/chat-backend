package contracts

import "github.com/google/uuid"

type UserCreatedMessage struct {
	Email    string
	Username string
	Code     string
}

type MessageCreatedBody struct {
	SenderId       uuid.UUID
	ConversationId int
	MessageBody    string
}

const (
	UserCreatedTopic    = "user.created"
	MessageCreatedTopic = "message.created"
)
