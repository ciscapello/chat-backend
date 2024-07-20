package messageservice

import (
	"errors"
	"fmt"

	"github.com/ciscapello/api_gateway/internal/common/jwtmanager"
	"github.com/ciscapello/api_gateway/internal/infrastructure/repository"
	"github.com/ciscapello/api_gateway/internal/infrastructure/wsClient"
	"github.com/ciscapello/api_gateway/pkg/dto"
	"github.com/ciscapello/chat-lib/contracts"
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
	ws                 *wsClient.WsClient
}

func New(messagesStorer MessagesStorer, wsClient *wsClient.WsClient, conversationStorer ConversationStorer, logger *zap.Logger, jwtManager *jwtmanager.JwtManager) *MessagesService {
	return &MessagesService{
		messagesStorer:     messagesStorer,
		conversationStorer: conversationStorer,
		logger:             logger,
		jwtManager:         jwtManager,
		ws:                 wsClient,
	}
}

func (ms *MessagesService) CreateMessage(senderId uuid.UUID, receiverId uuid.UUID, conversationId int, message string) error {

	err := ms.messagesStorer.CreateMessage(senderId, conversationId, message)
	if err != nil {
		return err
	}

	body := contracts.MessageSocketBody{
		FromUserID:     senderId.String(),
		ToUserID:       receiverId.String(),
		ConversationId: conversationId,
		MessageBody:    message,
	}

	fmt.Println(body)

	err = ms.ws.SendMessage(body)
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
