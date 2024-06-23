package conversationhandler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type createConversationRequest struct {
	SecondUserId string `json:"second_user_id"`
}

// @Summary Create conversation
// @Description Create conversation
// @Tags conversations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createConversationRequest true "Request body containing user_id and second_user_id"
// @Success 200
// @Failure 400 {object} response.Response{error=string}
// @Router /conversations [post]
func (ch *ConversationHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	creatorIdStr, err := ch.jwtManager.GetUserId(r.Context())
	if err != nil {
		ch.responder.SendError(w, http.StatusBadRequest, "invalid token")
		ch.logErrorInRequest(r, "invalid token")
		return
	}

	var body createConversationRequest
	b, err := io.ReadAll(r.Body)
	if err != nil {
		ch.responder.SendError(w, http.StatusBadRequest, "Unable to read request body")
		ch.logErrorInRequest(r, "Unable to read request body")
		return
	}

	err = json.Unmarshal(b, &body)
	if err != nil {
		ch.responder.SendError(w, http.StatusBadRequest, "Unable to unmarshal request body")
		ch.logErrorInRequest(r, "Unable to unmarshal request body")
		return
	}

	if body.SecondUserId == "" {
		ch.responder.SendError(w, http.StatusBadRequest, "second_user_id is required")
		ch.logErrorInRequest(r, "second_user_id is required")
		return
	}

	creatorId, err := uuid.Parse(creatorIdStr)
	if err != nil {
		ch.responder.SendError(w, http.StatusBadRequest, "invalid creator_id")
		ch.logErrorInRequest(r, "invalid creator_id")
		return
	}

	secondUserId, err := uuid.Parse(body.SecondUserId)
	if err != nil {
		ch.responder.SendError(w, http.StatusBadRequest, "invalid second_user_id")
		ch.logErrorInRequest(r, "invalid second_user_id")
		return
	}

	err = ch.conversationService.CreateConversation(creatorId, secondUserId)
	if err != nil {
		ch.responder.SendError(w, http.StatusBadRequest, err.Error())
		ch.logErrorInRequest(r, err.Error())
		return
	}

	ch.responder.SendSuccess(w, http.StatusOK, "ok")
}
