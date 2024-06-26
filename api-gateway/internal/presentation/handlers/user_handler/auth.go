package userhandler

import (
	"encoding/json"
	"io"
	"net/http"
)

type requestBody struct {
	Email string `json:"email"`
}

type resp struct {
	ID string `json:"id"`
}

// @Summary Authentication
// @Description Authentication by username and email
// @Tags users
// @Accept json
// @Produce json
// @Param request body requestBody true "Request body containing username and email"
// @Success 200 {object} response.Response{data=resp}
// @Failure 400 {object} response.Response{error=string}
// @Router /api/v1/users/auth [post]
func (uh *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {

	var rb requestBody
	body, err := io.ReadAll(r.Body)
	if err != nil {
		uh.responder.SendError(w, http.StatusBadRequest, "Unable to read request body")
		uh.logErrorInRequest(r, "Unable to read request body")
		return
	}

	if err := json.Unmarshal(body, &rb); err != nil {
		uh.responder.SendError(w, http.StatusBadRequest, "Unable to unmarshal request body")
		uh.logErrorInRequest(r, "Unable to unmarshal request body")
		return
	}

	if !uh.isValidEmail(rb.Email) {
		uh.responder.SendError(w, http.StatusBadRequest, "Invalid email")
		uh.logErrorInRequest(r, "Invalid email")
		return
	}

	uid, err := uh.userService.Authentication(rb.Email)
	if err != nil {
		uh.responder.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	uh.responder.SendSuccess(w, http.StatusOK, resp{
		ID: uid.String(),
	})
}
