package userhandler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ciscapello/api-gateway/internal/presentation/response"
	"github.com/google/uuid"
)

type checkCodeReq struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

func (uh *UserHandler) CheckCode(w http.ResponseWriter, r *http.Request) {

	var body checkCodeReq
	b, err := io.ReadAll(r.Body)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Unable to read request body")
		uh.logErrorInRequest(r, "Unable to read request body")
		return
	}

	err = json.Unmarshal(b, &body)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Unable to unmarshal request body")
		uh.logErrorInRequest(r, "Unable to unmarshal request body")
		return
	}

	id, err := uuid.Parse(body.ID)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid id")
		uh.logErrorInRequest(r, "invalid id")
		return
	}

	isEqual, err := uh.userService.CheckCode(id, body.Code)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, err.Error())
		uh.logErrorInRequest(r, err.Error())
		return
	}

	if !isEqual {
		response.SendError(w, http.StatusBadRequest, "Invalid code")
		uh.logErrorInRequest(r, "Invalid code")
		return
	}

	response.SendSuccess(w, http.StatusOK, "Code is valid")
}
