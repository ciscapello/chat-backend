package messageservice

import (
	"errors"
	"fmt"

	"github.com/ciscapello/chat-lib/contracts"
	"github.com/ciscapello/message_service/internal/infrastructure/repository"
	"github.com/ciscapello/message_service/internal/infrastructure/wsClient"
	"github.com/ciscapello/message_service/pkg/dto"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MessagesStorer interface {
	CreateMessage(senderId uuid.UUID, conversationId int, messageText string) error
	GetMessagesByConversationID(id int) ([]repository.MessagesRow, error)
	CheckIfConversationBelongsToUser(userId uuid.UUID, conversationId int) bool
}

type MessagesService struct {
	messagesStorer MessagesStorer
	ws             *wsClient.WsClient

	logger *zap.Logger
}

func New(messagesStorer MessagesStorer, ws *wsClient.WsClient, logger *zap.Logger) *MessagesService {
	return &MessagesService{
		messagesStorer: messagesStorer,
		ws:             ws,
		logger:         logger,
	}
}

func (ms *MessagesService) CreateMessage(senderId string, receiverId string, conversationId int, message string) error {

	senderUid, err := uuid.Parse(senderId)
	if err != nil {
		return err
	}

	err = ms.messagesStorer.CreateMessage(senderUid, conversationId, message)
	if err != nil {
		return err
	}

	body := contracts.MessageSocketBody{
		FromUserID:     senderId,
		ToUserID:       receiverId,
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
	isConvBelongsToUser := ms.messagesStorer.CheckIfConversationBelongsToUser(userId, conversationId)
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
