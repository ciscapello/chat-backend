package userhandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ciscapello/api-gateway/internal/domain/entity/userEntity"
	"github.com/ciscapello/api-gateway/internal/presentation/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// @Summary Update user
// @Description Update user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body userEntity.UpdateUserRequest true "User with optional fields"
// @Success 200 {object} response.Response{data=userEntity.PublicUser}
// @Failure 400 {object} response.Response{error=string}
// @Router /users/{id} [put]
func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Context().Value("userId"))

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		response.SendError(w, http.StatusBadRequest, "Missing user ID parameter")
		uh.logErrorInRequest(r, "Missing user ID parameter")
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid id")
		uh.logErrorInRequest(r, "invalid id")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Unable to read request body")
		uh.logErrorInRequest(r, "Unable to read request body")
		return
	}
	var requestBody userEntity.UpdateUserRequest

	if err := json.Unmarshal(body, &requestBody); err != nil {
		response.SendError(w, http.StatusBadRequest, "Unable to unmarshal request body")
		uh.logErrorInRequest(r, "Unable to unmarshal request body")
		return
	}

	if requestBody.Email != nil && !uh.isValidEmail(*requestBody.Email) {
		response.SendError(w, http.StatusBadRequest, "Invalid email")
		uh.logErrorInRequest(r, "Invalid email")
		return
	}

	us, err := uh.userService.UpdateUser(uid, requestBody)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, err.Error())
		uh.logErrorInRequest(r, err.Error())
		return
	}

	response.SendSuccess(w, http.StatusOK, us)
}
