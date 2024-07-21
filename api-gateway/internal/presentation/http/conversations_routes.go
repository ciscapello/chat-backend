package http

import (
	"net/http"

	conversationhandler "github.com/ciscapello/api_gateway/internal/presentation/handlers/conversation_handler"
	"github.com/gorilla/mux"
)

func ConfigureConversationRoutes(subrouter *mux.Router,
	handlers *conversationhandler.ConversationHandler,
	jwtMiddleware mux.MiddlewareFunc) {
	subrouter.Handle("", jwtMiddleware.Middleware(http.HandlerFunc(handlers.CreateConversation))).Methods(http.MethodPost)
	subrouter.Handle("", jwtMiddleware.Middleware(http.HandlerFunc(handlers.GetConversations))).Methods(http.MethodGet)
}
