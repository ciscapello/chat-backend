package userhandler

import (
	"net/http"

	"github.com/ciscapello/api-gateway/internal/presentation/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// @Summary Search users
// @Description Search users
// @Tags users
// @Accept json
// @Produce json
// @Param username query string false "Search query"
// @Success 200 {object} response.Response{data=[]userEntity.PublicUser}
// @Failure 400 {object} response.Response{error=string}
// @Router /users/search [get]
func (uh *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("username")

	id, err := uh.jwtManager.GetUserId(r.Context())
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid token")
		uh.logErrorInRequest(r, "invalid token")
		return
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid id")
		uh.logErrorInRequest(r, "invalid id")
		return
	}

	if username == "" {
		response.SendError(w, http.StatusBadRequest, "username is required")
		uh.logErrorInRequest(r, "username is required")
		return
	}

	users, err := uh.userService.FindUsersByUsername(username, uid)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "cannot get users")
		uh.logger.Error("cannot get users", zap.Error(err))
		return
	}

	response.SendSuccess(w, http.StatusOK, users)
}
