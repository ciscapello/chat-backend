package defaulthandler

import (
	"go.uber.org/zap"
)

type DefaultHandler struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *DefaultHandler {
	return &DefaultHandler{
		logger: logger,
	}
}
