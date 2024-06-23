package userhandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	userEntity "github.com/ciscapello/api-gateway/internal/domain/entity/user_entity"
	"github.com/google/uuid"
)

// @Summary Update user
// @Description Update user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body userEntity.UpdateUserRequest true "User with optional fields"
// @Success 200 {object} response.Response{data=userEntity.PublicUser}
// @Failure 400 {object} response.Response{error=string}
// @Router /users [put]
func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		uh.responder.SendError(w, http.StatusBadRequest, "Unable to read request body")
		uh.logErrorInRequest(r, "Unable to read request body")
		return
	}
	var requestBody userEntity.UpdateUserRequest

	if err := json.Unmarshal(body, &requestBody); err != nil {
		uh.responder.SendError(w, http.StatusBadRequest, "Unable to unmarshal request body")
		uh.logErrorInRequest(r, "Unable to unmarshal request body")
		return
	}

	fmt.Println(requestBody)

	if requestBody.Email != nil && !uh.isValidEmail(*requestBody.Email) {
		uh.responder.SendError(w, http.StatusBadRequest, "Invalid email")
		uh.logErrorInRequest(r, "Invalid email")
		return
	}

	us, err := uh.userService.UpdateUser(uid, requestBody)
	if err != nil {
		uh.responder.SendError(w, http.StatusBadRequest, err.Error())
		uh.logErrorInRequest(r, err.Error())
		return
	}

	uh.responder.SendSuccess(w, http.StatusOK, us)
}
