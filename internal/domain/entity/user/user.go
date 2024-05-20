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

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Code     *string `json:"code,omitempty"`
	Role     *Role   `json:"role,omitempty"`
	Enabled  *bool   `json:"enabled,omitempty"`
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

func NewPublicUser(u User) PublicUser {
	return PublicUser{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}

func (u *User) Update(params UpdateUserRequest) {
	if params.Username != nil {
		u.Username = *params.Username
	}
	if params.Email != nil {
		u.Email = *params.Email
	}
	if params.Code != nil {
		u.Code = *params.Code
	}
	if params.Role != nil {
		u.Role = *params.Role
	}
	if params.Enabled != nil {
		u.Enabled = *params.Enabled
	}
}
