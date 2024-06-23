package userhandler

import (
	"net/http"

	"github.com/ciscapello/api-gateway/internal/common/jwtmanager"
	userservice "github.com/ciscapello/api-gateway/internal/domain/service/user_service"
	"go.uber.org/zap"
)

type Responder interface {
	SendSuccess(w http.ResponseWriter, statusCode int, payload interface{})
	SendError(w http.ResponseWriter, statusCode int, message string)
}

type UserHandler struct {
	userService userservice.IUserService
	logger      *zap.Logger
	jwtManager  jwtmanager.IJwtManager
	responder   Responder
}

func New(userService userservice.IUserService, logger *zap.Logger, jwtManager jwtmanager.IJwtManager, responder Responder) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
		jwtManager:  jwtManager,
		responder:   responder,
	}
}
