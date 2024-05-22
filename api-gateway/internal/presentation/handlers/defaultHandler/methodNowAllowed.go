package defaulthandler

import (
	"net/http"

	"go.uber.org/zap"
)

func (dh *DefaultHandler) MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	dh.logger.Warn("method not allowed",
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()))
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("method not allowed"))
}
