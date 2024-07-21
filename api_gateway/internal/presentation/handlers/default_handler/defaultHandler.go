package defaulthandler

import (
	"log/slog"
)

type DefaultHandler struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *DefaultHandler {
	return &DefaultHandler{
		logger: logger,
	}
}
