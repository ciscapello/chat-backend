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
	Email    string
	Code     string
	Role     Role
	Enabled  bool
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
