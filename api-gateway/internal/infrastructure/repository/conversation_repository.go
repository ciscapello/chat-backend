package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ConversationRepository struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewConversationRepository(db *sql.DB, logger *zap.Logger) *ConversationRepository {
	return &ConversationRepository{
		logger: logger,
		db:     db,
	}
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
