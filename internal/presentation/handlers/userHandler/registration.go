package userhandler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ciscapello/chat-backend/internal/presentation/response"
)

type requestBody struct {
	Username string
	Email    string
}

type resp struct {
	ID string `json:"id"`
}

func (uh *UserHandler) Registration(w http.ResponseWriter, r *http.Request) {

	var rb requestBody
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Unable to read request body")
		uh.logErrorInRequest(r, "Unable to read request body")
		return
	}

	if err := json.Unmarshal(body, &rb); err != nil {
		response.SendError(w, http.StatusBadRequest, "Unable to unmarshal request body")
		uh.logErrorInRequest(r, "Unable to unmarshal request body")
		return
	}

	if !uh.isValidEmail(rb.Email) {
		response.SendError(w, http.StatusBadRequest, "Invalid email")
		uh.logErrorInRequest(r, "Invalid email")
		return
	}

	uid, err := uh.userService.Registration(rb.Username, rb.Email)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, err.Error())
		uh.logErrorInRequest(r, err.Error())
		return
	}

	response.SendSuccess(w, http.StatusOK, resp{
		ID: uid.String(),
	})
}
