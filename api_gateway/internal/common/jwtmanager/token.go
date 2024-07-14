package jwtmanager

import (
	userEntity "github.com/ciscapello/api_gateway/internal/domain/entity/user_entity"
	"github.com/google/uuid"
)

type tokenClaims struct {
	id   uuid.UUID
	role userEntity.Role
}

type ReturnTokenType struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
