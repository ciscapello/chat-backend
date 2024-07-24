package userservice

import (
	"github.com/ciscapello/api-gateway/internal/common/jwtmanager"
	userEntity "github.com/ciscapello/api-gateway/internal/domain/entity/user_entity"
	"github.com/google/uuid"
)

type IUserService interface {
	Authentication(email string) (uuid.UUID, error)
	GetUser(uuid uuid.UUID) (userEntity.PublicUser, error)
	GetAllUsers() ([]userEntity.PublicUser, error)
	UpdateUser(uuid uuid.UUID, fields userEntity.UpdateUserRequest) (userEntity.PublicUser, error)
	CheckCode(uuid uuid.UUID, code string) (bool, error)
	GetUserRole(id uuid.UUID) userEntity.Role
	GetTokens(id uuid.UUID, role userEntity.Role) (jwtmanager.ReturnTokenType, error)
	FindUsersByUsername(username string, id uuid.UUID) ([]userEntity.PublicUser, error)
}
