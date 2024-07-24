package userhandler

import (
	"log/slog"
	"net/http"

	userEntity "github.com/ciscapello/api-gateway/internal/domain/entity/user_entity"
)

// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]userEntity.PublicUser}
// @Failure 400 {object} response.Response{error=string}
// @Router /api/v1/users [get]
func (uh *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	role, err := uh.jwtManager.GetUserRole(r.Context())
	if err != nil {
		uh.responder.SendError(w, http.StatusBadRequest, "invalid token")
		uh.logErrorInRequest(r, "invalid token")
		return
	}

	if role != userEntity.Admin {
		uh.responder.SendError(w, http.StatusForbidden, "forbidden")
		uh.logErrorInRequest(r, "forbidden")
		return
	}

	users, err := uh.userService.GetAllUsers()
	if err != nil {
		uh.responder.SendError(w, http.StatusBadRequest, "cannot get users")
		uh.logger.Error("cannot get users", slog.String("message", err.Error()))
		return
	}

	uh.responder.SendSuccess(w, http.StatusOK, users)
}
