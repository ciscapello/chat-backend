package jwtmanager

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	userEntity "github.com/ciscapello/api-gateway/internal/domain/entity/user_entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	logger     *zap.Logger
	jwtManager *JwtManager
}

type contextKey string

const (
	authorizationHeader = "Authorization"

	userIdCtx   contextKey = "userId"
	userRoleCtx contextKey = "userRole"
)

func NewAuthMiddleware(logger *zap.Logger, j *JwtManager) *AuthMiddleware {
	return &AuthMiddleware{
		logger:     logger,
		jwtManager: j,
	}
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authorizationHeader)
		fmt.Println(authHeader)

		if authHeader == "" {
			am.logger.Error("Authorization header missing")
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		claims, err := am.jwtManager.verifyToken(authHeader)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				am.logger.Error("Token expired", zap.Error(err))
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			} else {
				am.logger.Error("Invalid token", zap.Error(err))
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

		}

		ctx := context.WithValue(r.Context(), userIdCtx, claims.id)
		ctx = context.WithValue(ctx, userRoleCtx, claims.role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (JwtManager) GetUserId(ctx context.Context) (string, error) {
	userId, ok := ctx.Value(userIdCtx).(uuid.UUID)
	if !ok {
		return "", errors.New("user id is not in the context")
	}
	return userId.String(), nil
}

func (JwtManager) GetUserRole(ctx context.Context) (userEntity.Role, error) {
	role, ok := ctx.Value(userRoleCtx).(userEntity.Role)
	if !ok {
		return userEntity.Role(0), errors.New("user id is not in the context")
	}
	return role, nil
}
