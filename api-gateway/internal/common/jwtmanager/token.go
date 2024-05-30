package jwtmanager

import (
	"github.com/ciscapello/api-gateway/internal/domain/entity/userEntity"
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
