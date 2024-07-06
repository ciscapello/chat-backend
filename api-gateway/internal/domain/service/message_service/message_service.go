package messageservice

import (
	"errors"

	"github.com/ciscapello/api-gateway/internal/common/jwtmanager"
	"github.com/ciscapello/api-gateway/internal/infrastructure/repository"
	"github.com/ciscapello/api-gateway/pkg/dto"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MessagesStorer interface {
	CreateMessage(senderId uuid.UUID, conversationId int, messageText string) error
	GetMessagesByConversationID(id int) ([]repository.MessagesRow, error)
}

type ConversationStorer interface {
	CheckIfConversationBelongsToUser(userId uuid.UUID, conversationId int) bool
}

type MessagesService struct {
	messagesStorer     MessagesStorer
	conversationStorer ConversationStorer
	logger             *zap.Logger
	jwtManager         *jwtmanager.JwtManager
}

func New(messagesStorer MessagesStorer, conversationStorer ConversationStorer, logger *zap.Logger, jwtManager *jwtmanager.JwtManager) *MessagesService {
	return &MessagesService{
		messagesStorer:     messagesStorer,
		conversationStorer: conversationStorer,
		logger:             logger,
		jwtManager:         jwtManager,
	}
}

func (ms *MessagesService) CreateMessage(senderId uuid.UUID, conversationId int, message string) error {

	return ms.messagesStorer.CreateMessage(senderId, conversationId, message)
}

func (ms *MessagesService) GetMessages(conversationId int, userId uuid.UUID) ([]dto.MessageDTO, error) {
	isConvBelongsToUser := ms.conversationStorer.CheckIfConversationBelongsToUser(userId, conversationId)
	if !isConvBelongsToUser {
		return nil, errors.New("conversation is not belongs to user")
	}

	messages, err := ms.messagesStorer.GetMessagesByConversationID(conversationId)
	if err != nil {
		return nil, err
	}

	res := make([]dto.MessageDTO, len(messages))

	for i, message := range messages {
		msgDto := dto.MessageDTO{}

		uid, err := uuid.Parse(message.SenderId)
		if err != nil {
			return nil, err
		}

		msgDto.CreatedAt = message.CreatedAt
		msgDto.MessageBody = message.Message
		msgDto.SenderId = uid
		msgDto.User.ID = uid
		msgDto.User.Username = message.Username
		res[i] = msgDto
	}

	return res, nil
}
