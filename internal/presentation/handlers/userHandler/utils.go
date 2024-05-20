package userhandler

import (
	"net/http"
	"regexp"

	"go.uber.org/zap"
)

func (uh *UserHandler) logErrorInRequest(r *http.Request, msg string) {
	uh.logger.Error(msg, zap.String("url", r.URL.String()), zap.String("method", r.Method))
}

func (uh *UserHandler) isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}
