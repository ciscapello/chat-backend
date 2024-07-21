package userhandler

import (
	"log/slog"
	"net/http"

	"github.com/ciscapello/api_gateway/internal/common/jwtmanager"
	userservice "github.com/ciscapello/api_gateway/internal/domain/service/user_service"
)

type Responder interface {
	SendSuccess(w http.ResponseWriter, statusCode int, payload interface{})
	SendError(w http.ResponseWriter, statusCode int, message string)
}

type UserHandler struct {
	userService userservice.IUserService
	logger      *slog.Logger
	jwtManager  jwtmanager.IJwtManager
	responder   Responder
}

func New(userService userservice.IUserService, logger *slog.Logger, jwtManager jwtmanager.IJwtManager, responder Responder) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
		jwtManager:  jwtManager,
		responder:   responder,
	}
}
