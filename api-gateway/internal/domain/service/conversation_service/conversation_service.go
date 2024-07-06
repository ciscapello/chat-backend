package conversationservice

import (
	"fmt"

	"github.com/ciscapello/api-gateway/internal/common/jwtmanager"
	userEntity "github.com/ciscapello/api-gateway/internal/domain/entity/user_entity"
	"github.com/ciscapello/api-gateway/internal/infrastructure/repository"
	"github.com/ciscapello/api-gateway/pkg/dto"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ConversationStorer interface {
	CreateConversation(creatorId uuid.UUID, secondUserId uuid.UUID) error
	GetConversationsList(userId uuid.UUID) ([]repository.ConversationsWithUser, error)
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

func (us *ConversationService) GetUserConversations(userId uuid.UUID) ([]dto.ConversationsDTO, error) {
	conversations, err := us.conversationStorer.GetConversationsList(userId)
	dtos := make([]dto.ConversationsDTO, len(conversations))
	if err != nil {
		return nil, err
	}

	for index, conv := range conversations {
		convDTO := dto.ConversationsDTO{}
		usr := userEntity.PublicUser{}
		msgDTO := &dto.ConversationMessageDTO{}

		uid, err := uuid.Parse(conv.UserID)
		if err != nil {
			return nil, err
		}
		convDTO.ID = conv.ID
		usr.ID = uid
		usr.Email = conv.Email
		usr.Username = conv.Username
		convDTO.User = usr

		senderUid, err := uuid.Parse(conv.LastMessageSenderId)
		if err != nil {
			fmt.Println(conv.LastMessageSenderId, "err")
			msgDTO = nil
		} else {
			msgDTO.MessageBody = conv.LastMessageBody
			msgDTO.CreatedAt = conv.LastMessageCreatedAt
			msgDTO.SenderId = senderUid
		}

		convDTO.LastMessage = msgDTO
		dtos[index] = convDTO
	}

	return dtos, nil
}
