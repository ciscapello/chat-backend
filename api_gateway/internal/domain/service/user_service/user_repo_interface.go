package userservice

import (
	userEntity "github.com/ciscapello/api_gateway/internal/domain/entity/user_entity"
	"github.com/google/uuid"
)

type UserRepo interface {
	GetUserById(id uuid.UUID) (userEntity.User, error)
	GetUserByEmail(email string) (userEntity.User, error)
	CreateUser(user userEntity.User) error
	CheckUserIfExistsByUsername(username string) bool
	CheckUserIfExistsByEmail(email string) bool
	GetAllUsers() ([]userEntity.User, error)
	UpdateUser(u userEntity.User) error
	UpdateCode(id uuid.UUID, code string) error
	GetUserRole(id uuid.UUID) userEntity.Role
	FindUsersByUsername(username string, uid uuid.UUID) ([]userEntity.PublicUser, error)
}
