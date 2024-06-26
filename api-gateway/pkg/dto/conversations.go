package dto

import userEntity "github.com/ciscapello/api-gateway/internal/domain/entity/user_entity"

type ConversationsDTO struct {
	ID          int                   `json:"id"`
	User        userEntity.PublicUser `json:"user"`
	LastMessage string                `json:"last_message"`
}
