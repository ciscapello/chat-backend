package userhandler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ciscapello/chat-backend/internal/domain/entity/user"
	"github.com/ciscapello/chat-backend/internal/presentation/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
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
	var requestBody user.UpdateUserRequest

	if err := json.Unmarshal(body, &requestBody); err != nil {
		response.SendError(w, http.StatusBadRequest, "Unable to unmarshal request body")
		uh.logErrorInRequest(r, "Unable to unmarshal request body")
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
