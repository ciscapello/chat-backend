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
	GetAllUsers() ([]user.User, error)
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

func (us *UserService) Registration(username, email string) (uuid.UUID, error) {
	if isExists := us.userRepo.CheckUserIfExistsByUsername(username); isExists {
		return uuid.UUID{}, ErrUserWithThisUsernameExists
	}
	if isExists := us.userRepo.CheckUserIfExistsByEmail(email); isExists {
		return uuid.UUID{}, ErrUserWithThisEmailExists
	}

	code, err := utils.GenerateOneTimeCode(6)
	if err != nil {
		us.logger.Error("failed to generate code", zap.Error(err))
		return uuid.UUID{}, ErrCannotCreateUser
	}
	user := user.NewUser(username, email, code)

	err = us.userRepo.CreateUser(*user)
	if err != nil {
		us.logger.Error("failed to create user", zap.Error(err))
		return uuid.UUID{}, ErrCannotCreateUser
	}
	return user.ID, nil
}

func (us *UserService) GetUser(uuid uuid.UUID) (user.PublicUser, error) {
	u, err := us.userRepo.GetUserById(uuid)
	if err != nil {
		return user.PublicUser{}, err
	}

	return user.PublicUser{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}, nil
}

func (us *UserService) GetAllUsers() ([]user.PublicUser, error) {
	users, err := us.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	var publicUsers []user.PublicUser

	for _, us := range users {
		publicUsers = append(publicUsers, user.PublicUser{
			ID:       us.ID,
			Username: us.Username,
			Email:    us.Email,
		})
	}
	return publicUsers, nil

}

func (us *UserService) UpdateUser() {}
