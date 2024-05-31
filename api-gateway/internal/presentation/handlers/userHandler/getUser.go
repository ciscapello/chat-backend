package userhandler

import (
	"net/http"

	"github.com/ciscapello/api-gateway/internal/common/jwtmanager"
	"github.com/ciscapello/api-gateway/internal/presentation/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// @Summary Get user
// @Description Get user by id
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response{data=userEntity.PublicUser}
// @Failure 400 {object} response.Response{error=string}
// @Router /users/{id} [get]
func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	id, err := jwtmanager.GetUserId(r.Context())
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid token")
		uh.logErrorInRequest(r, "invalid token")
		return
	}

	vars := mux.Vars(r)
	if id != vars["id"] {
		response.SendError(w, http.StatusNotFound, "cannot found user with this id")
		uh.logErrorInRequest(r, "invalid token")
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid id")
		uh.logErrorInRequest(r, "invalid id")
		return
	}

	user, err := uh.userService.GetUser(uid)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "cannot get user")
		uh.logErrorInRequest(r, "cannot get user")
		return
	}

	response.SendSuccess(w, http.StatusOK, user)
}
