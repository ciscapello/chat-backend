package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ciscapello/api-gateway/internal/domain/entity/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

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

func (ur *UserRepository) CheckUserIfExistsByEmail(email string) bool {
	var user user.User
	query := "SELECT * FROM users WHERE email = $1"
	row := ur.db.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Username, &user.Code, &user.Role, &user.Enabled)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return true
	}
	return true
}

func (ur *UserRepository) CheckUserIfExistsByUsername(username string) bool {
	fmt.Println(username)
	var user user.User
	query := "SELECT * FROM users WHERE username = $1"
	row := ur.db.QueryRow(query, username)
	err := row.Scan(&user.ID, &user.Username, &user.Code, &user.Role, &user.Enabled)
	fmt.Println(row)
	fmt.Println(user)
	if err == sql.ErrNoRows {
		fmt.Println(err)
		return false
	} else if err != nil {
		fmt.Println(err)
		return true
	}
	return true
}

func (ur *UserRepository) GetUserById(id uuid.UUID) (user.User, error) {
	var us user.User
	query := "SELECT * FROM users WHERE id = $1"
	row := ur.db.QueryRow(query, id)

	var createdAt string
	var updatedAt string
	var roleString string

	err := row.Scan(&us.ID, &us.Username, &us.Enabled, &roleString, &createdAt, &updatedAt, &us.Code, &us.Email)
	us.Role = user.ParseRole(roleString)

	if err == sql.ErrNoRows {
		ur.logger.Error("User not found", zap.String("id", id.String()))
		return us, ErrUserNotFound
	} else if err != nil {
		ur.logger.Error(err.Error(), zap.String("id", id.String()))
		return us, err
	}
	return us, nil
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

func (ur *UserRepository) GetAllUsers() ([]user.User, error) {
	query := "SELECT * FROM users"

	rows, err := ur.db.Query(query)
	if err != nil {
		ur.logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var users []user.User

	var createdAt string
	var updatedAt string
	var roleString string

	for rows.Next() {
		var us user.User
		err := rows.Scan(&us.ID, &us.Username, &us.Enabled, &roleString, &createdAt, &updatedAt, &us.Code, &us.Email)
		if err != nil {
			ur.logger.Error(err.Error())
			return nil, err
		}
		us.Role = user.ParseRole(roleString)

		users = append(users, us)
	}

	if err = rows.Err(); err != nil {
		ur.logger.Error(err.Error())
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) UpdateUser(u user.User) error {
	query := "UPDATE users SET username = $1, email = $2, enabled = $3, role = $4, code = $5, updated_at = $6 WHERE id = $7"
	res, err := ur.db.Exec(query, u.Username, u.Email, u.Enabled, u.Role.String(), u.Code, time.Now(), u.ID)
	fmt.Println("here?")
	if err != nil {
		ur.logger.Error(err.Error())
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		ur.logger.Error(err.Error())
		return err
	}

	if count == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (ur *UserRepository) Registration() {}

func (ur *UserRepository) Login() {}

func (ur *UserRepository) DeleteUser() {}
