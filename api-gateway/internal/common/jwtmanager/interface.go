package jwtmanager

import (
	"context"

	userEntity "github.com/ciscapello/api-gateway/internal/domain/entity/user_entity"
	"github.com/google/uuid"
)

type IJwtManager interface {
	Generate(uid uuid.UUID, role userEntity.Role) (ReturnTokenType, error)
	VerifyRefreshToken(refreshTokenStr string) (string, error)
	GetUserId(ctx context.Context) (string, error)
	GetUserRole(ctx context.Context) (userEntity.Role, error)
}
