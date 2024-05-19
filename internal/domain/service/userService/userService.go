package userservice

import (
	"errors"

	"github.com/ciscapello/chat-backend/internal/common/utils"
	"github.com/ciscapello/chat-backend/internal/domain/entity/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo UserRepo
	logger   *zap.Logger
}

type UserRepo interface {
	GetUserById(id uuid.UUID) (user.User, error)
	CreateUser(user user.User) error
	CheckUserIfExistsByUsername(username string) bool
	CheckUserIfExistsByEmail(email string) bool
}

var (
	ErrCannotCreateUser           = errors.New("cannot create user")
	ErrUserWithThisUsernameExists = errors.New("user with this username already exists")
	ErrUserWithThisEmailExists    = errors.New("user with this email already exists")
)

func New(userRepo UserRepo, logger *zap.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (us *UserService) Login() {}

func (us *UserService) Registration(username, email string) error {
	if isExists := us.userRepo.CheckUserIfExistsByUsername(username); isExists {
		return ErrUserWithThisUsernameExists
	}
	if isExists := us.userRepo.CheckUserIfExistsByEmail(email); isExists {
		return ErrUserWithThisEmailExists
	}

	code, err := utils.GenerateOneTimeCode(6)
	if err != nil {
		us.logger.Error("failed to generate code", zap.Error(err))
		return ErrCannotCreateUser
	}
	user := user.NewUser(username, email, code)

	err = us.userRepo.CreateUser(*user)
	if err != nil {
		us.logger.Error("failed to create user", zap.Error(err))
		return ErrCannotCreateUser
	}
	return nil
}

func (us *UserService) GetUser(uuid uuid.UUID) {
	us.userRepo.GetUserById(uuid)
}

func (us *UserService) GetAllUsers() {}

func (us *UserService) UpdateUser() {}
