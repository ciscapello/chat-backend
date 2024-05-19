package userhandler

import (
	"net/http"

	"go.uber.org/zap"
)

func (uh *UserHandler) logErrorInRequest(r *http.Request, msg string) {
	uh.logger.Error(msg, zap.String("url", r.URL.String()), zap.String("method", r.Method))
}
