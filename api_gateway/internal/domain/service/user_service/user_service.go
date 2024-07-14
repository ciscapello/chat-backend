package userservice

import (
	"errors"
	"fmt"
	"time"

	"github.com/ciscapello/api_gateway/internal/common/jwtmanager"
	"github.com/ciscapello/api_gateway/internal/common/utils"
	userEntity "github.com/ciscapello/api_gateway/internal/domain/entity/user_entity"
	"github.com/ciscapello/lib/contracts"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo      UserRepo
	logger        *zap.Logger
	messageBroker MessageBroker
	jwtManager    *jwtmanager.JwtManager
}

type MessageBroker interface {
	Publish(topic string, msg interface{}) error
}

var (
	ErrCannotCreateUser           = errors.New("cannot create user")
	ErrUserWithThisUsernameExists = errors.New("user with this username already exists")
	ErrUserWithThisEmailExists    = errors.New("user with this email already exists")
)

func New(userRepo UserRepo, logger *zap.Logger, messageBroker MessageBroker, jwtmanager *jwtmanager.JwtManager) *UserService {
	return &UserService{
		userRepo:      userRepo,
		logger:        logger,
		messageBroker: messageBroker,
		jwtManager:    jwtmanager,
	}
}

func (us *UserService) Authentication(email string) (uuid.UUID, error) {
	isExists := us.userRepo.CheckUserIfExistsByEmail(email)

	var user userEntity.User
	var err error

	code, _ := utils.GenerateOneTimeCode(6)

	if !isExists {
		user = *userEntity.NewUser("", email, code)
		err = us.userRepo.CreateUser(user)
		if err != nil {
			us.logger.Error("failed to create user", zap.Error(err))
			return uuid.UUID{}, ErrCannotCreateUser
		}
	} else {
		user, err = us.userRepo.GetUserByEmail(email)
		if err != nil {
			us.logger.Error("failed to get user", zap.Error(err))
			return uuid.UUID{}, err
		}

		lastCodeUpdateUTC := user.LastCodeUpdate.UTC()
		currentTimeUTC := time.Now().UTC()

		if lastCodeUpdateUTC.Add(time.Minute).After(currentTimeUTC) {
			return uuid.UUID{}, fmt.Errorf("you can receive code again in %v", time.Until(lastCodeUpdateUTC.Add(time.Minute)))
		}

		err := us.userRepo.UpdateCode(user.ID, code)
		if err != nil {
			us.logger.Error("failed to update code", zap.Error(err))
			return uuid.UUID{}, err
		}

		user.Code = code
	}

	us.messageBroker.Publish(contracts.UserCreatedTopic, contracts.UserCreatedMessage{
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

func (us *UserService) CheckCode(uuid uuid.UUID, code string) (bool, error) {
	usr, err := us.userRepo.GetUserById(uuid)
	if err != nil {
		us.logger.Error("failed to get user", zap.Error(err))
		return false, err
	}

	return code == usr.Code, nil
}

func (us *UserService) GetUserRole(id uuid.UUID) userEntity.Role {
	return us.userRepo.GetUserRole(id)
}

func (us *UserService) GetTokens(id uuid.UUID, role userEntity.Role) (jwtmanager.ReturnTokenType, error) {
	tokens, err := us.jwtManager.Generate(id, role)
	if err != nil {
		us.logger.Error("failed to generate tokens", zap.Error(err))
		return jwtmanager.ReturnTokenType{}, err
	}
	return tokens, nil
}

func (us *UserService) FindUsersByUsername(username string, id uuid.UUID) ([]userEntity.PublicUser, error) {
	return us.userRepo.FindUsersByUsername(username, id)
}
