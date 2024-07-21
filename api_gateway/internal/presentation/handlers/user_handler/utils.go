package userhandler

import (
	"log/slog"
	"net/http"
	"regexp"
	"unicode/utf8"
)

func (uh *UserHandler) logErrorInRequest(r *http.Request, msg string) {
	uh.logger.Error(msg, slog.String("url", r.URL.String()), slog.String("method", r.Method))
}

func (uh *UserHandler) isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email) && utf8.RuneCountInString(email) <= 40
}
