package messagehandler

import (
	"github.com/ciscapello/api_gateway/internal/common/jwtmanager"
	messageservice "github.com/ciscapello/api_gateway/internal/domain/service/message_service"
	"github.com/ciscapello/api_gateway/internal/presentation/response"
	"go.uber.org/zap"
)

type MessagesHandler struct {
	messagesService *messageservice.MessagesService
	logger          *zap.Logger
	jwtManager      *jwtmanager.JwtManager
	responder       response.Responder
}

func New(messagesService *messageservice.MessagesService, logger *zap.Logger, jwtManager *jwtmanager.JwtManager, responder response.Responder) *MessagesHandler {
	return &MessagesHandler{
		messagesService: messagesService,
		logger:          logger,
		jwtManager:      jwtManager,
		responder:       responder,
	}
}
