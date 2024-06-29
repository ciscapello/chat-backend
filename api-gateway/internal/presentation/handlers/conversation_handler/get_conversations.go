package conversationhandler

import (
	"net/http"

	"github.com/google/uuid"
)

// @Summary Get conversations
// @Description Get conversations that belongs to user
// @Tags conversations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param skip query integer false "Skip items count"
// @Success 200 {object} response.Response{data=dto.ConversationsDTO}
// @Failure 400 {object} response.Response{error=string}
// @Router /api/v1/conversations [get]
func (ch *ConversationHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	userId, err := ch.jwtManager.GetUserId(r.Context())

	if err != nil {
		ch.responder.SendError(w, http.StatusBadRequest, "invalid token")
		ch.logErrorInRequest(r, "invalid token")
		return
	}

	uid, err := uuid.Parse(userId)
	if err != nil {
		ch.responder.SendError(w, http.StatusBadRequest, "invalid creator_id")
		ch.logErrorInRequest(r, "invalid creator_id")
		return
	}

	conversations, err := ch.conversationService.GetUserConversations(uid)
	if err != nil {
		ch.responder.SendError(w, http.StatusInternalServerError, err.Error())
		ch.logErrorInRequest(r, err.Error())
		return
	}

	ch.responder.SendSuccess(w, http.StatusOK, conversations)
}
