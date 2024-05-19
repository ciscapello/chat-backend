package repository

import (
	"database/sql"
	"errors"

	"github.com/ciscapello/chat-backend/internal/domain/entity/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewUserRepository(
	db *sql.DB,
	logger *zap.Logger,
) *UserRepository {
	return &UserRepository{
		logger: logger,
		db:     db,
	}
}

func (ur *UserRepository) GetUserById(id uuid.UUID) (user.User, error) {
	var user user.User
	query := "SELECT * FROM users WHERE id = $1"
	row := ur.db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.Enabled)
	if err == sql.ErrNoRows {
		ur.logger.Error("User not found", zap.String("id", id.String()))
		return user, ErrUserNotFound
	} else if err != nil {
		ur.logger.Error(err.Error(), zap.String("id", id.String()))
		return user, err
	}
	return user, nil
}

func (ur *UserRepository) GetAllUsers() {}

func (ur *UserRepository) UpdateUser() {}

func (ur *UserRepository) Registration() {}

func (ur *UserRepository) Login() {}

func (ur *UserRepository) DeleteUser() {}
