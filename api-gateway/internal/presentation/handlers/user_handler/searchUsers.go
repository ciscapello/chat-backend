package userhandler

import (
	"net/http"

	"github.com/ciscapello/api-gateway/internal/presentation/response"
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

	users, err := uh.userService.FindUsersByUsername(username)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "cannot get users")
		uh.logger.Error("cannot get users", zap.Error(err))
	}

	response.SendSuccess(w, http.StatusOK, users)
}
