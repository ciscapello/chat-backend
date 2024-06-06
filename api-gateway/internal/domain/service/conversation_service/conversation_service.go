package conversationservice

import (
	"github.com/ciscapello/api-gateway/internal/common/jwtmanager"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ConversationStorer interface {
	CreateConversation(creatorId uuid.UUID, secondUserId uuid.UUID) error
	// GetConversations(userId uuid.UUID) ([]uuid.UUID, error)
}

type ConversationService struct {
	conversationStorer ConversationStorer
	logger             *zap.Logger
	jwtManager         *jwtmanager.JwtManager
}

func New(conversationStorer ConversationStorer, logger *zap.Logger, jwtManager *jwtmanager.JwtManager) *ConversationService {
	return &ConversationService{
		conversationStorer: conversationStorer,
		logger:             logger,
		jwtManager:         jwtManager,
	}
}

func (us *ConversationService) CreateConversation(creatorId uuid.UUID, secondUserId uuid.UUID) error {
	return us.conversationStorer.CreateConversation(creatorId, secondUserId)
}
