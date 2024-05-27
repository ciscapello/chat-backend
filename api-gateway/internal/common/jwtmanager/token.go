package jwtmanager

import "github.com/google/uuid"

type accessToken struct {
	id uuid.UUID
}

type ReturnTokenType struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
