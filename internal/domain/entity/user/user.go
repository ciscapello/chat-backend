package user

import "github.com/google/uuid"

const (
	Admin Role = iota
	Regular
)

type Role int

type User struct {
	ID       uuid.UUID
	username string
	password string
	role     Role
	enabled  bool
}

func NewUser(username, password string, role Role) *User {
	return &User{
		ID:       uuid.New(),
		username: username,
		password: password,
		role:     role,
		enabled:  true,
	}
}
