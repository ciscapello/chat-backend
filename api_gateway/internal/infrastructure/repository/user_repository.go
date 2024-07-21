package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	userEntity "github.com/ciscapello/api_gateway/internal/domain/entity/user_entity"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	logger *slog.Logger
	db     *sql.DB
}

func NewUserRepository(
	db *sql.DB,
	logger *slog.Logger,
) *UserRepository {
	return &UserRepository{
		logger: logger,
		db:     db,
	}
}

func (ur *UserRepository) CheckUserIfExistsByEmail(email string) bool {
	var user userEntity.User
	query := "SELECT * FROM users WHERE email = $1"
	row := ur.db.QueryRow(query, email)
	var username sql.NullString
	err := row.Scan(&user.ID, &username, &user.Code, &user.Role, &user.Enabled)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return true
	}
	return true
}

func (ur *UserRepository) CheckUserIfExistsByUsername(username string) bool {
	fmt.Println(username)
	var user userEntity.User
	var name sql.NullString
	query := "SELECT * FROM users WHERE username = $1"
	row := ur.db.QueryRow(query, username)
	err := row.Scan(&user.ID, &name, &user.Code, &user.Role, &user.Enabled)
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

func (ur *UserRepository) GetUserById(id uuid.UUID) (userEntity.User, error) {
	var us userEntity.User
	query := "SELECT * FROM users WHERE id = $1"
	row := ur.db.QueryRow(query, id)

	var createdAt string
	var updatedAt string
	var roleString string

	var username sql.NullString

	err := row.Scan(&us.ID, &username, &us.Enabled, &roleString, &createdAt, &updatedAt, &us.Code, &us.Email, &us.LastCodeUpdate)
	us.Role = userEntity.ParseRole(roleString)

	if username.Valid {
		us.Username = username.String
	} else {
		us.Username = ""
	}

	if err == sql.ErrNoRows {
		ur.logger.Error("User not found", slog.String("id", id.String()))
		return us, ErrUserNotFound
	} else if err != nil {
		ur.logger.Error(err.Error(), slog.String("id", id.String()))
		return us, err
	}
	return us, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (userEntity.User, error) {
	var us userEntity.User
	query := "SELECT * FROM users WHERE email = $1"
	row := ur.db.QueryRow(query, email)

	var createdAt string
	var updatedAt string
	var roleString string
	var username sql.NullString

	err := row.Scan(&us.ID, &username, &us.Enabled, &roleString, &createdAt, &updatedAt, &us.Code, &us.Email, &us.LastCodeUpdate)
	if err == sql.ErrNoRows {
		ur.logger.Error("User not found", slog.String("email", email))
		return us, ErrUserNotFound
	} else if err != nil {
		ur.logger.Error(err.Error(), slog.String("email", email))
		return us, err
	}
	if username.Valid {
		us.Username = username.String
	} else {
		us.Username = ""
	}
	us.Role = userEntity.ParseRole(roleString)

	return us, nil
}

func (ur *UserRepository) CreateUser(u userEntity.User) error {
	query := "INSERT INTO users (id, username, enabled, role, code, email, last_code_update) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	res, err := ur.db.Exec(query, u.ID, nil, u.Enabled, u.Role.String(), u.Code, u.Email, time.Now().UTC())
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

func (ur *UserRepository) GetAllUsers() ([]userEntity.User, error) {
	query := "SELECT * FROM users"

	rows, err := ur.db.Query(query)
	if err != nil {
		ur.logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var users []userEntity.User

	var createdAt string
	var updatedAt string
	var roleString string

	for rows.Next() {
		var us userEntity.User
		err := rows.Scan(&us.ID, &us.Username, &us.Enabled, &roleString, &createdAt, &updatedAt, &us.Code, &us.Email, &us.LastCodeUpdate)
		if err != nil {
			ur.logger.Error(err.Error())
			return nil, err
		}
		us.Role = userEntity.ParseRole(roleString)

		users = append(users, us)
	}

	if err = rows.Err(); err != nil {
		ur.logger.Error(err.Error())
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) UpdateUser(u userEntity.User) error {
	query := "UPDATE users SET username = $1, email = $2, enabled = $3, role = $4, code = $5, updated_at = $6 WHERE id = $7"
	res, err := ur.db.Exec(query, u.Username, u.Email, u.Enabled, u.Role.String(), u.Code, time.Now(), u.ID)
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

func (ur *UserRepository) UpdateCode(id uuid.UUID, code string) error {
	query := "UPDATE users SET code = $1, last_code_update = $2 WHERE id = $3"
	res, err := ur.db.Exec(query, code, time.Now().UTC(), id)
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

func (ur *UserRepository) GetUserRole(id uuid.UUID) userEntity.Role {
	fmt.Println(id)
	query := "SELECT role FROM users WHERE id = $1"
	var role string

	row := ur.db.QueryRow(query, id)
	fmt.Println(row)
	err := row.Scan(&role)
	if err != nil {
		fmt.Println(err)
		ur.logger.Error(err.Error())
		return userEntity.Regular
	}

	fmt.Println(role)

	if role != userEntity.Admin.String() {
		return userEntity.Regular
	}

	return userEntity.Admin
}

func (ur *UserRepository) FindUsersByUsername(username string, uid uuid.UUID) ([]userEntity.PublicUser, error) {
	query := `SELECT id, username, email FROM users WHERE id != $1 AND username IS NOT NULL AND LOWER(username) LIKE '%' || LOWER($2) || '%'`

	rows, err := ur.db.Query(query, uid, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []userEntity.PublicUser

	for rows.Next() {
		var us userEntity.PublicUser
		err := rows.Scan(&us.ID, &us.Username, &us.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, us)
	}

	return users, nil
}

func (ur *UserRepository) Login() {}

func (ur *UserRepository) DeleteUser() {}
