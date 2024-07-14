package repository

import (
	"database/sql"
	"fmt"
	"time"

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

type MessagesRow struct {
	ID        int
	SenderId  string
	CreatedAt time.Time
	Message   string
	Username  string
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

func (mr *MessageRepository) GetMessagesByConversationID(id int) ([]MessagesRow, error) {
	stmt := `select m.id, m.sender_id, m.created_at, m.message, u.username from messages m 
	join users u ON m.sender_id = u.id 
	where m.conversation_id = $1
	`
	rows, err := mr.db.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	row := MessagesRow{}
	nullableUsername := sql.NullString{}
	res := []MessagesRow{}

	for rows.Next() {
		err := rows.Scan(&row.ID, &row.SenderId, &row.CreatedAt, &row.Message, &nullableUsername)
		if err != nil {
			return nil, err
		}
		row.Username = nullableUsername.String
		res = append(res, row)
	}

	return res, nil
}

func (cr *MessageRepository) CheckIfConversationBelongsToUser(userId uuid.UUID, conversationId int) bool {
	query := `select c.id from conversations c 
	join participants p on p.conversation_id = c.id 
	where p.user_id = $1 and c.id = $2
	`
	row := cr.db.QueryRow(query, userId, conversationId)
	if row.Err() != nil {
		return false
	}
	var id int
	err := row.Scan(&id)
	return err == nil
}
