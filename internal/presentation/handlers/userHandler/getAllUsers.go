package userhandler

import (
	"net/http"

	"github.com/ciscapello/chat-backend/internal/presentation/response"
	"go.uber.org/zap"
)

func (uh *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uh.userService.GetAllUsers()
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "cannot get users")
		uh.logger.Error("cannot get users", zap.Error(err))
	}

	response.SendSuccess(w, http.StatusOK, users)
}
