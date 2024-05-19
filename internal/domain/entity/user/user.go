package user

import "github.com/google/uuid"

const (
	Admin Role = iota
	Regular
)

type Role int

type User struct {
	ID       uuid.UUID
	Username string
	Password string
	Role     Role
	Enabled  bool
}

func NewUser(username, password string, role Role) *User {
	return &User{
		ID:       uuid.New(),
		Username: username,
		Password: password,
		Role:     role,
		Enabled:  true,
	}
}
