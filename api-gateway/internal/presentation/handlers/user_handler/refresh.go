package userhandler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ciscapello/api-gateway/internal/presentation/response"
	"github.com/google/uuid"
)

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" example:"your_refresh_token_here"`
}

// @Summary Refresh
// @Description Refresh token
// @Tags users
// @Accept json
// @Produce json
// @Param request body refreshRequest true "Request body containing refresh token"
// @Success 200 {object} response.Response{data=jwtmanager.ReturnTokenType}
// @Failure 400 {object} response.Response{error=string}
// @Router /users/refresh [post]
func (uh *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {

	var body refreshRequest
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

	if body.RefreshToken == "" {
		response.SendError(w, http.StatusBadRequest, "refresh token is required")
		uh.logErrorInRequest(r, "refresh token is required")
		return
	}

	idStr, err := uh.jwtManager.VerifyRefreshToken(body.RefreshToken)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid refresh token")
		uh.logErrorInRequest(r, "invalid refresh token")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid id")
		uh.logErrorInRequest(r, "invalid id")
		return
	}

	role := uh.userService.GetUserRole(id)

	tokens, err := uh.jwtManager.Generate(id, role)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "unable to generate tokens")
		uh.logErrorInRequest(r, "unable to generate tokens")
		return
	}

	response.SendSuccess(w, http.StatusOK, tokens)
}
