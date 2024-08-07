package userhandler

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// @Summary Search users
// @Description Search users
// @Tags users
// @Accept json
// @Produce json
// @Param username query string false "Search query"
// @Success 200 {object} response.Response{data=[]userEntity.PublicUser}
// @Failure 400 {object} response.Response{error=string}
// @Router /api/v1/users/search [get]
func (uh *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("username")

	id, err := uh.jwtManager.GetUserId(r.Context())
	if err != nil {
		uh.responder.SendError(w, http.StatusBadRequest, "invalid token")
		uh.logErrorInRequest(r, "invalid token")
		return
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		uh.responder.SendError(w, http.StatusBadRequest, "invalid id")
		uh.logErrorInRequest(r, "invalid id")
		return
	}

	if username == "" {
		uh.responder.SendError(w, http.StatusBadRequest, "username is required")
		uh.logErrorInRequest(r, "username is required")
		return
	}

	users, err := uh.userService.FindUsersByUsername(username, uid)
	if err != nil {
		uh.responder.SendError(w, http.StatusBadRequest, "cannot get users")
		uh.logger.Error("cannot get users", slog.String("message", err.Error()))
		return
	}

	uh.responder.SendSuccess(w, http.StatusOK, users)
}
