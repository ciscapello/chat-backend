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

func (ur *UserRepository) CheckUserIfExists(username string) bool {
	var user user.User
	query := "SELECT * FROM users WHERE username = $1"
	row := ur.db.QueryRow(query, username)
	err := row.Scan(&user.ID, &user.Username, &user.Code, &user.Role, &user.Enabled)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return true
	}
	return true
}

func (ur *UserRepository) GetUserById(id uuid.UUID) (user.User, error) {
	var user user.User
	query := "SELECT * FROM users WHERE id = $1"
	row := ur.db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Username, &user.Code, &user.Role, &user.Enabled)
	if err == sql.ErrNoRows {
		ur.logger.Error("User not found", zap.String("id", id.String()))
		return user, ErrUserNotFound
	} else if err != nil {
		ur.logger.Error(err.Error(), zap.String("id", id.String()))
		return user, err
	}
	return user, nil
}

func (ur *UserRepository) CreateUser(u user.User) error {
	query := "INSERT INTO users (id, username, enabled, role, code, email) VALUES ($1, $2, $3, $4, $5, $6)"

	res, err := ur.db.Exec(query, u.ID, u.Username, u.Enabled, u.Role.String(), u.Code, u.Email)
	if err != nil {
		ur.logger.Error(err.Error())
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		ur.logger.Error(err.Error())
		return err
	}

	return nil
}

func (ur *UserRepository) GetAllUsers() {}

func (ur *UserRepository) UpdateUser() {}

func (ur *UserRepository) Registration() {}

func (ur *UserRepository) Login() {}

func (ur *UserRepository) DeleteUser() {}
