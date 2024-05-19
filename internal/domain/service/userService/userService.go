package userservice

import (
	"github.com/ciscapello/chat-backend/internal/domain/entity/user"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo UserRepo
}

type UserRepo interface {
	GetUserById(id uuid.UUID) (user.User, error)
}

func New(userRepo UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (us *UserService) Login() {}

func (us *UserService) Registration() {}

func (us *UserService) GetUser(uuid uuid.UUID) {
	us.userRepo.GetUserById(uuid)
}

func (us *UserService) GetAllUsers() {}

func (us *UserService) UpdateUser() {}
