package conversationhandler

import (
	"log/slog"

	"github.com/ciscapello/api_gateway/internal/common/jwtmanager"
	conversationservice "github.com/ciscapello/api_gateway/internal/domain/service/conversation_service"
	"github.com/ciscapello/api_gateway/internal/presentation/response"
)

type ConversationHandler struct {
	conversationService *conversationservice.ConversationService
	logger              *slog.Logger
	jwtManager          *jwtmanager.JwtManager
	responder           response.Responder
}

func New(conversationservice *conversationservice.ConversationService, logger *slog.Logger, jwtManager *jwtmanager.JwtManager, responder response.Responder) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationservice,
		logger:              logger,
		jwtManager:          jwtManager,
		responder:           responder,
	}
}
