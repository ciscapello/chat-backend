package messagehandler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// @Summary Get messages
// @Description Get messages by conversation id
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Conversation ID"
// @Success 200 {object} response.Response{data=[]dto.MessageDTO}
// @Failure 400 {object} response.Response{error=string}
// @Router /api/v1/messages/{conversation_id} [get]
func (mh *MessagesHandler) GetMessages(w http.ResponseWriter, r *http.Request) {

	userId, err := mh.jwtManager.GetUserId(r.Context())

	if err != nil {
		mh.responder.SendError(w, http.StatusBadRequest, "cannot get user id")
		mh.logErrorInRequest(r, "cannot get user id")
		return
	}

	uid, err := uuid.Parse(userId)
	if err != nil {
		mh.responder.SendError(w, http.StatusBadRequest, "Unable to parse user id")
		mh.logErrorInRequest(r, "Unable to parse user id")
		return
	}

	params := mux.Vars(r)
	conversationIdStr, ok := params["conversation_id"]
	if !ok {
		mh.responder.SendError(w, http.StatusBadRequest, "cannot found conversation_id in request params")
		mh.logErrorInRequest(r, "cannot found conversation_id in request params")
		return
	}

	conversationId, err := strconv.Atoi(conversationIdStr)
	if err != nil {
		mh.responder.SendError(w, http.StatusBadRequest, "cannot parse conversation_id in request params")
		mh.logErrorInRequest(r, "cannot parse conversation_id in request params")
		return
	}

	messages, err := mh.messagesService.GetMessages(conversationId, uid)
	if err != nil {
		mh.responder.SendError(w, http.StatusInternalServerError, err.Error())
		mh.logErrorInRequest(r, err.Error())
		return
	}

	mh.responder.SendSuccess(w, http.StatusOK, messages)
}
