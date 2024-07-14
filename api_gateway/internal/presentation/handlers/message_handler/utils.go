package messagehandler

import (
	"net/http"

	"go.uber.org/zap"
)

func (ch *MessagesHandler) logErrorInRequest(r *http.Request, msg string) {
	ch.logger.Error(msg, zap.String("url", r.URL.String()), zap.String("method", r.Method))
}
