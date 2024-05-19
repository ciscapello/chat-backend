package user

import (
	"github.com/google/uuid"
)

const (
	Admin Role = iota
	Regular
)

type Role int

func (r Role) String() string {
	return [...]string{"admin", "regular"}[r]
}

func ParseRole(s string) Role {
	switch s {
	case "admin":
		return Admin
	case "regular":
		return Regular
	default:
		return Regular
	}
}

type User struct {
	ID       uuid.UUID
	Username string
	Email    string
	Code     string
	Role     Role
	Enabled  bool
}

type PublicUser struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

func NewUser(username, email, code string) *User {
	return &User{
		ID:       uuid.New(),
		Username: username,
		Email:    email,
		Code:     code,
		Role:     Regular,
		Enabled:  true,
	}
}
