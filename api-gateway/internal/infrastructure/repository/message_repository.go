package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MessageRepository struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewMessagesRepository(
	db *sql.DB,
	logger *zap.Logger,
) *MessageRepository {
	return &MessageRepository{
		logger: logger,
		db:     db,
	}
}

func (mr *MessageRepository) CreateMessage(senderId uuid.UUID, conversationId int, messageText string) error {
	senderIdStr := senderId.String()
	fmt.Println(senderIdStr)
	fmt.Println(conversationId)
	fmt.Println(messageText)
	stmt := `
	insert into messages (sender_id, conversation_id, message)
	select $1, $2, $3
	where exists (
		select id from participants p
		where p.conversation_id = $4
		and p.user_id = $5
	)`

	res, err := mr.db.Exec(stmt, senderId, conversationId, messageText, conversationId, senderId)
	if err != nil {
		fmt.Println(res)
		return err
	}

	return nil
}
