package jwtmanager

import (
	"context"
	"errors"
	"net/http"

	"github.com/ciscapello/api-gateway/internal/domain/entity/userEntity"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	requiredRole userEntity.Role
	logger       *zap.Logger
	jwtManager   *JwtManager
}

type contextKey string

const (
	authorizationHeader = "Authorization"

	userIdCtx contextKey = "userId"
)

func NewAuthMiddleware(requiredRole userEntity.Role, logger *zap.Logger, j *JwtManager) *AuthMiddleware {
	return &AuthMiddleware{
		requiredRole: requiredRole,
		logger:       logger,
		jwtManager:   j,
	}
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authorizationHeader)
		if authHeader == "" {
			am.logger.Error("Authorization header missing")
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		claims, err := am.jwtManager.verifyToken(authHeader)
		if err != nil {
			am.logger.Error("Invalid token", zap.Error(err))
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims.role != am.requiredRole {
			am.logger.Error("Invalid role", zap.String("role", claims.role.String()))
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), userIdCtx, claims.id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserId(ctx context.Context) (string, error) {
	userId, ok := ctx.Value(userIdCtx).(uuid.UUID)
	if !ok {
		return "", errors.New("user id is not in the context")
	}
	return userId.String(), nil
}
