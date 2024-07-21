package conversationhandler

import (
	"log/slog"
	"net/http"
)

func (ch *ConversationHandler) logErrorInRequest(r *http.Request, msg string) {
	ch.logger.Error(msg, slog.String("url", r.URL.String()), slog.String("method", r.Method))
}
