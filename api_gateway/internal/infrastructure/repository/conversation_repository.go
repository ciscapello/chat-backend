package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ConversationRepository struct {
	logger *zap.Logger
	db     *sql.DB
}

type ConversationsWithUser struct {
	ID                   int
	UserID               string
	Username             string
	Email                string
	LastMessageBody      string
	LastMessageCreatedAt time.Time
	LastMessageSenderId  string
}

func NewConversationRepository(db *sql.DB, logger *zap.Logger) *ConversationRepository {
	return &ConversationRepository{
		logger: logger,
		db:     db,
	}
}

func (cr *ConversationRepository) CheckIfConversationBelongsToUser(userId uuid.UUID, conversationId int) bool {
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

func (cr *ConversationRepository) CheckIfConversationExists(creatorId uuid.UUID, secondUserId uuid.UUID) bool {
	query := `
	SELECT p1.conversation_id FROM participants p1
	JOIN participants p2 ON p1.conversation_id = p2.conversation_id
	WHERE (p1.user_id = $1 AND p2.user_id = $2) OR (p1.user_id = $2 AND p2.user_id = $1)
	`

	row := cr.db.QueryRow(query, creatorId, secondUserId)

	var conversationId int
	err := row.Scan(&conversationId)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		cr.logger.Error(err.Error())
		return false
	}

	return true
}

func (cr *ConversationRepository) CreateConversation(creatorId uuid.UUID, secondUserId uuid.UUID) error {
	if cr.CheckIfConversationExists(creatorId, secondUserId) {
		return fmt.Errorf("conversation with users %v and %v already exists", creatorId, secondUserId)
	}

	tx, err := cr.db.Begin()
	if err != nil {
		return err
	}

	statement := "INSERT INTO conversations (creator_id) VALUES ($1) RETURNING id"
	var conversationId int
	err = tx.QueryRow(statement, creatorId).Scan(&conversationId)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO participants (user_id, conversation_id) VALUES ($1, $2), ($3, $4)", creatorId, conversationId, secondUserId, conversationId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (cr *ConversationRepository) GetConversationsList(userId uuid.UUID) ([]ConversationsWithUser, error) {
	query := `select c.id, u2.id, u2.username, u2.email, m.message, m.created_at, m.sender_id from conversations c 
	join participants p on p.conversation_id = c.id 
	join users u2 on u2.id = p.user_id
	left join (
		select 
            m1.conversation_id, 
            m1.message, 
            m1.created_at,
            m1.sender_id
        from 
            messages m1
            join (
                select 
                    conversation_id, 
                    max(created_at) AS latest_message_time
                from 
                    messages
                group by
                    conversation_id
            ) m2 ON m1.conversation_id = m2.conversation_id AND m1.created_at = m2.latest_message_time
	) m on c.id = m.conversation_id
	where p.conversation_id in (
		select c.id from conversations c
		join participants p ON p.conversation_id = c.id 
		where p.user_id = $1
	)
	and u2.id <> $1
	`

	convs := []ConversationsWithUser{}

	nullableUsername := sql.NullString{}
	nullableMessageBody := sql.NullString{}
	nullableCreatedAt := sql.NullTime{}
	nullableSenderId := sql.NullString{}

	rows, err := cr.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		conv := ConversationsWithUser{}

		err := rows.Scan(&conv.ID, &conv.UserID, &nullableUsername, &conv.Email, &nullableMessageBody, &nullableCreatedAt, &nullableSenderId)
		conv.Username = nullableUsername.String
		conv.LastMessageBody = nullableMessageBody.String
		conv.LastMessageCreatedAt = nullableCreatedAt.Time
		conv.LastMessageSenderId = nullableSenderId.String
		convs = append(convs, conv)
		if err != nil {
			return nil, err
		}
	}

	return convs, nil
}
