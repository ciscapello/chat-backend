package userservice

import (
	"errors"

	"github.com/ciscapello/api-gateway/internal/common/utils"
	"github.com/ciscapello/api-gateway/internal/domain/entity/user"
	"github.com/ciscapello/api-gateway/internal/infrastructure/rabbitmq"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo      UserRepo
	logger        *zap.Logger
	messageBroker MessageBroker
}

type MessageBroker interface {
	Publish(topic string, msg interface{}) error
}

type UserRepo interface {
	GetUserById(id uuid.UUID) (user.User, error)
	CreateUser(user user.User) error
	CheckUserIfExistsByUsername(username string) bool
	CheckUserIfExistsByEmail(email string) bool
	GetAllUsers() ([]user.User, error)
	UpdateUser(u user.User) error
}

var (
	ErrCannotCreateUser           = errors.New("cannot create user")
	ErrUserWithThisUsernameExists = errors.New("user with this username already exists")
	ErrUserWithThisEmailExists    = errors.New("user with this email already exists")
)

func New(userRepo UserRepo, logger *zap.Logger, messageBroker MessageBroker) *UserService {
	return &UserService{
		userRepo:      userRepo,
		logger:        logger,
		messageBroker: messageBroker,
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

	us.messageBroker.Publish(rabbitmq.UserCreatedTopic, rabbitmq.UserCreatedMessage{
		Username: user.Username,
		Email:    user.Email,
		Code:     user.Code,
	})

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

func (us *UserService) UpdateUser(uuid uuid.UUID, fields user.UpdateUserRequest) (user.PublicUser, error) {
	u, err := us.userRepo.GetUserById(uuid)
	if err != nil {
		return user.PublicUser{}, err
	}

	u.Update(fields)

	err = us.userRepo.UpdateUser(u)
	if err != nil {
		return user.PublicUser{}, err
	}

	publicUser := user.NewPublicUser(u)

	return publicUser, nil
}
