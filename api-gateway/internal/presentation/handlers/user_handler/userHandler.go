package userhandler

import (
	"github.com/ciscapello/api-gateway/internal/common/jwtmanager"
	userservice "github.com/ciscapello/api-gateway/internal/domain/service/userService"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService *userservice.UserService
	logger      *zap.Logger
	jwtManager  *jwtmanager.JwtManager
}

func New(userService *userservice.UserService, logger *zap.Logger, jwtManager *jwtmanager.JwtManager) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
		jwtManager:  jwtManager,
	}
}
