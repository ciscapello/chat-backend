package userservice

import (
	"errors"

	"github.com/ciscapello/api-gateway/internal/common/utils"
	"github.com/ciscapello/api-gateway/internal/domain/entity/userEntity"
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
	GetUserById(id uuid.UUID) (userEntity.User, error)
	CreateUser(user userEntity.User) error
	CheckUserIfExistsByUsername(username string) bool
	CheckUserIfExistsByEmail(email string) bool
	GetAllUsers() ([]userEntity.User, error)
	UpdateUser(u userEntity.User) error
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
	user := userEntity.NewUser(username, email, code)

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

func (us *UserService) GetUser(uuid uuid.UUID) (userEntity.PublicUser, error) {
	u, err := us.userRepo.GetUserById(uuid)
	if err != nil {
		return userEntity.PublicUser{}, err
	}

	return userEntity.PublicUser{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}, nil
}

func (us *UserService) GetAllUsers() ([]userEntity.PublicUser, error) {
	users, err := us.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	var publicUsers []userEntity.PublicUser

	for _, us := range users {
		publicUsers = append(publicUsers, userEntity.PublicUser{
			ID:       us.ID,
			Username: us.Username,
			Email:    us.Email,
		})
	}
	return publicUsers, nil

}

func (us *UserService) UpdateUser(uuid uuid.UUID, fields userEntity.UpdateUserRequest) (userEntity.PublicUser, error) {
	u, err := us.userRepo.GetUserById(uuid)
	if err != nil {
		return userEntity.PublicUser{}, err
	}

	u.Update(fields)

	err = us.userRepo.UpdateUser(u)
	if err != nil {
		return userEntity.PublicUser{}, err
	}

	publicUser := userEntity.NewPublicUser(u)

	return publicUser, nil
}
