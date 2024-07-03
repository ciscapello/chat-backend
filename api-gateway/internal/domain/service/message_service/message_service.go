package messageservice

import (
	"github.com/ciscapello/api-gateway/internal/common/jwtmanager"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MessagesStorer interface {
	CreateMessage(senderId uuid.UUID, conversationId int, messageText string) error
}

type MessagesService struct {
	messagesStorer MessagesStorer
	logger         *zap.Logger
	jwtManager     *jwtmanager.JwtManager
}

func New(messagesStorer MessagesStorer, logger *zap.Logger, jwtManager *jwtmanager.JwtManager) *MessagesService {
	return &MessagesService{
		messagesStorer: messagesStorer,
		logger:         logger,
		jwtManager:     jwtManager,
	}
}

func (ms *MessagesService) CreateMessage(senderId uuid.UUID, conversationId int, message string) error {

	return ms.messagesStorer.CreateMessage(senderId, conversationId, message)
}
