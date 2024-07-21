package defaulthandler

import (
	"log/slog"
	"net/http"
)

func (dh *DefaultHandler) MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	dh.logger.Warn("method not allowed",
		slog.String("method", r.Method),
		slog.String("url", r.URL.String()))
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("method not allowed"))
}
