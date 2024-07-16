package messageservice

import (
	"errors"

	"github.com/ciscapello/api_gateway/internal/common/jwtmanager"
	"github.com/ciscapello/api_gateway/internal/infrastructure/repository"
	"github.com/ciscapello/api_gateway/pkg/dto"
	"github.com/ciscapello/lib/contracts"
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

type MessageBroker interface {
	Publish(topic string, msg interface{}) error
}

type MessagesService struct {
	messagesStorer     MessagesStorer
	conversationStorer ConversationStorer
	logger             *zap.Logger
	jwtManager         *jwtmanager.JwtManager
	messageBroker      MessageBroker
}

func New(messagesStorer MessagesStorer, messageBroker MessageBroker, conversationStorer ConversationStorer, logger *zap.Logger, jwtManager *jwtmanager.JwtManager) *MessagesService {
	return &MessagesService{
		messagesStorer:     messagesStorer,
		conversationStorer: conversationStorer,
		logger:             logger,
		jwtManager:         jwtManager,
		messageBroker:      messageBroker,
	}
}

func (ms *MessagesService) CreateMessage(senderId uuid.UUID, receiverId uuid.UUID, conversationId int, message string) error {

	var body contracts.MessageSocketBody

	body.ConversationId = conversationId
	body.FromUserID = senderId.String()
	body.ToUserID = receiverId.String()
	body.MessageBody = message

	err := ms.messageBroker.Publish(contracts.MessageCreatedTopic, body)
	if err != nil {
		return err
	}

	return nil
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
