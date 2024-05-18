package handlers

import userservice "github.com/ciscapello/chat-backend/internal/domain/service/userService"

type Handlers struct {
	userService *userservice.UserService
}

func New(userService *userservice.UserService) *Handlers {
	return &Handlers{
		userService: userService,
	}
}
