package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type MessageRepository struct {
	logger *slog.Logger
	db     *sql.DB
}

func NewMessagesRepository(
	db *sql.DB,
	logger *slog.Logger,
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
