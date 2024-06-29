package userhandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// @Summary Get user
// @Description Get user by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.Response{data=userEntity.PublicUser}
// @Failure 400 {object} response.Response{error=string}
// @Router /api/v1/users/{id} [get]
func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		uh.responder.SendError(w, http.StatusNotFound, "id did not provide")
		uh.logErrorInRequest(r, "id did not provide")
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		uh.responder.SendError(w, http.StatusBadRequest, "invalid id")
		uh.logErrorInRequest(r, "invalid id")
		return
	}

	user, err := uh.userService.GetUser(uid)
	if err != nil {
		uh.responder.SendError(w, http.StatusBadRequest, "cannot get user")
		uh.logErrorInRequest(r, "cannot get user")
		return
	}

	uh.responder.SendSuccess(w, http.StatusOK, user)
}
