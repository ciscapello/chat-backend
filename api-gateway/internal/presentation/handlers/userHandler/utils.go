package userhandler

import (
	"net/http"
	"regexp"
	"unicode"
	"unicode/utf8"

	"go.uber.org/zap"
)

func (uh *UserHandler) logErrorInRequest(r *http.Request, msg string) {
	uh.logger.Error(msg, zap.String("url", r.URL.String()), zap.String("method", r.Method))
}

func (uh *UserHandler) isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email) && utf8.RuneCountInString(email) <= 40
}

func isAlphabetic(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func (uh *UserHandler) isValidUsername(username string) bool {
	if !isAlphabetic(username) {
		return false
	}

	return utf8.RuneCountInString(username) >= 6 && utf8.RuneCountInString(username) <= 24
}
