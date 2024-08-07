package messagehandler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type createMessageRequestBody struct {
	MessageTo      string `json:"message_to"`
	ConversationId int    `json:"conversation_id"`
	Text           string `json:"text"`
}

// @Summary Create message
// @Description Create message
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createMessageRequestBody true "Request body containing conversation_id and text of the message"
// @Success 200
// @Failure 400 {object} response.Response{error=string}
// @Router /api/v1/messages [post]
func (mh *MessagesHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {

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

	var body createMessageRequestBody

	b, err := io.ReadAll(r.Body)
	if err != nil {
		mh.responder.SendError(w, http.StatusBadRequest, "Unable to read request body")
		mh.logErrorInRequest(r, "Unable to read request body")
		return
	}

	err = json.Unmarshal(b, &body)
	if err != nil {
		mh.responder.SendError(w, http.StatusBadRequest, "Unable to unmarshal request body")
		mh.logErrorInRequest(r, "Unable to unmarshal request body")
		return
	}

	receiverId, err := uuid.Parse(body.MessageTo)
	if err != nil {
		mh.responder.SendError(w, http.StatusBadRequest, "Unable to parse receiver id")
		mh.logErrorInRequest(r, "Unable to parse receiver id")
		return
	}

	err = mh.messagesService.CreateMessage(uid, receiverId, body.ConversationId, body.Text)
	if err != nil {
		mh.responder.SendError(w, http.StatusInternalServerError, err.Error())
		mh.logErrorInRequest(r, err.Error())
		return
	}

	mh.responder.SendSuccess(w, http.StatusCreated, "ok")
}
