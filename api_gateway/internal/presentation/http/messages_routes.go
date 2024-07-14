package http

import (
	"net/http"

	messagehandler "github.com/ciscapello/api_gateway/internal/presentation/handlers/message_handler"
	"github.com/gorilla/mux"
)

func ConfigureMessagesRoutes(subrouter *mux.Router,
	handlers *messagehandler.MessagesHandler,
	jwtMiddleware mux.MiddlewareFunc) {
	subrouter.Handle("", jwtMiddleware.Middleware(http.HandlerFunc(handlers.CreateMessage))).Methods(http.MethodPost)
	subrouter.Handle("/{conversation_id}", jwtMiddleware.Middleware(http.HandlerFunc(handlers.GetMessages))).Methods(http.MethodGet)
}
