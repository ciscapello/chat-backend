package http

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)

			duration := time.Since(start)
			logger.Info("[HTTP]",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.Int("status", rw.status),
				slog.Duration("duration", duration),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

// func (rw *responseWriter) WriteHeader(code int) {
// 	if !rw.wroteHeader {
// 		rw.status = code
// 		rw.ResponseWriter.WriteHeader(code)
// 		rw.wroteHeader = true
// 	}
// }
