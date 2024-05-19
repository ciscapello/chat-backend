package user

import "github.com/google/uuid"

const (
	Admin Role = iota
	Regular
)

type Role int

func (r Role) String() string {
	return [...]string{"admin", "regular"}[r]
}

type User struct {
	ID       uuid.UUID
	Username string
	Code     string
	Role     Role
	Enabled  bool
}

func NewUser(username, code string, role Role) *User {
	return &User{
		ID:       uuid.New(),
		Username: username,
		Code:     code,
		Role:     role,
		Enabled:  true,
	}
}
