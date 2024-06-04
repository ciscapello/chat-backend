package http

import (
	"net/http"

	userhandler "github.com/ciscapello/api-gateway/internal/presentation/handlers/user_handler"
	"github.com/gorilla/mux"
)

func ConfigureUserRoutes(subrouter *mux.Router,
	handlers *userhandler.UserHandler,
	jwtMiddleware mux.MiddlewareFunc) {

	subrouter.Handle("", jwtMiddleware.Middleware(http.HandlerFunc(handlers.GetAllUsers))).Methods(http.MethodGet)
	subrouter.Handle("/{id}", jwtMiddleware.Middleware(http.HandlerFunc(handlers.GetUser))).Methods(http.MethodGet)
	subrouter.Handle("", jwtMiddleware.Middleware(http.HandlerFunc(handlers.UpdateUser))).Methods(http.MethodPut)
	subrouter.HandleFunc("/auth", handlers.Auth).Methods(http.MethodPost)
	subrouter.HandleFunc("/refresh", handlers.Refresh).Methods(http.MethodPost)
	subrouter.HandleFunc("/check-code", handlers.CheckCode).Methods(http.MethodPost)
}
