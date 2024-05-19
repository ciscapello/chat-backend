package userhandler

import (
	userservice "github.com/ciscapello/chat-backend/internal/domain/service/userService"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService *userservice.UserService
	logger      *zap.Logger
}

func New(userService *userservice.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}
