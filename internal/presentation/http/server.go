package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ciscapello/chat-backend/internal/application/config"
	"github.com/ciscapello/chat-backend/internal/presentation/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func New(cfg *config.Config, handlers *handlers.Handlers, logger *zap.Logger) *Server {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowedHandler)

	userRouter := router.PathPrefix("/users").Subrouter()
	ConfigureUserRoutes(userRouter, handlers)

	router.Handle("", userRouter)

	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		if pathTemplate != "" {
			fmt.Printf("[%s]: %s\n", strings.Join(methods, ","), pathTemplate)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return &Server{
		httpServer: &http.Server{
			Addr:         "0.0.0.0:" + cfg.HttpPort,
			Handler:      router,
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
		},
		logger: logger,
	}
}

func (s *Server) Run() error {
	s.logger.Info("Server is running", zap.String("port", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func ConfigureUserRoutes(subrouter *mux.Router, handlers *handlers.Handlers) {
	subrouter.HandleFunc("", handlers.GetAllUsers).Methods(http.MethodGet)
	subrouter.HandleFunc("/{id}", handlers.GetUser).Methods(http.MethodGet)
	subrouter.HandleFunc("", handlers.UpdateUser).Methods(http.MethodPut)
	subrouter.HandleFunc("/registration", handlers.Registration).Methods(http.MethodPost)
	subrouter.HandleFunc("/login", handlers.Login).Methods(http.MethodPost)
}
