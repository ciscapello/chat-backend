package userhandler

import (
	"net/http"

	"github.com/ciscapello/api-gateway/internal/domain/entity/userEntity"
	"github.com/ciscapello/api-gateway/internal/presentation/response"
	"go.uber.org/zap"
)

// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]userEntity.PublicUser}
// @Failure 400 {object} response.Response{error=string}
// @Router /users [get]
func (uh *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	role, err := uh.jwtManager.GetUserRole(r.Context())
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid token")
		uh.logErrorInRequest(r, "invalid token")
		return
	}

	if role != userEntity.Admin {
		response.SendError(w, http.StatusForbidden, "forbidden")
		uh.logErrorInRequest(r, "forbidden")
		return
	}

	users, err := uh.userService.GetAllUsers()
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "cannot get users")
		uh.logger.Error("cannot get users", zap.Error(err))
	}

	response.SendSuccess(w, http.StatusOK, users)
}
