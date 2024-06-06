package conversationhandler

import (
	"github.com/ciscapello/api-gateway/internal/common/jwtmanager"
	conversationservice "github.com/ciscapello/api-gateway/internal/domain/service/conversation_service"
	"go.uber.org/zap"
)

type ConversationHandler struct {
	conversationService *conversationservice.ConversationService
	logger              *zap.Logger
	jwtManager          *jwtmanager.JwtManager
}

func New(conversationservice *conversationservice.ConversationService, logger *zap.Logger, jwtManager *jwtmanager.JwtManager) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationservice,
		logger:              logger,
		jwtManager:          jwtManager,
	}
}
