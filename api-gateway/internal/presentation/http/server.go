package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ciscapello/api-gateway/internal/application/config"
	"github.com/ciscapello/api-gateway/internal/common/jwtmanager"
	conversationhandler "github.com/ciscapello/api-gateway/internal/presentation/handlers/conversation_handler"
	defaulthandler "github.com/ciscapello/api-gateway/internal/presentation/handlers/default_handler"
	messagehandler "github.com/ciscapello/api-gateway/internal/presentation/handlers/message_handler"
	userhandler "github.com/ciscapello/api-gateway/internal/presentation/handlers/user_handler"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handlers struct {
	UserHandler         *userhandler.UserHandler
	DefaultHandler      *defaulthandler.DefaultHandler
	ConversationHandler *conversationhandler.ConversationHandler
	MessageHandler      *messagehandler.MessagesHandler
}

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	handlers   *Handlers
}

func New(cfg *config.Config, handlers *Handlers, logger *slog.Logger) *Server {
	router := mux.NewRouter()
	routerWithPrefix := router.PathPrefix("/api").Subrouter()

	jwtm := jwtmanager.NewJwtManager(cfg, logger)

	jwtMiddleware := jwtmanager.NewAuthMiddleware(logger, jwtm)

	routerWithPrefix.Use(LoggingMiddleware(logger))

	routerWithPrefix.NotFoundHandler = http.HandlerFunc(handlers.DefaultHandler.NotFoundHandler)
	routerWithPrefix.MethodNotAllowedHandler = http.HandlerFunc(handlers.DefaultHandler.MethodNotAllowedHandler)

	routerWithPrefix.HandleFunc("/v1/health", handlers.DefaultHandler.HealthHandler).Methods(http.MethodGet)

	userRouter := routerWithPrefix.PathPrefix("/v1/users").Subrouter()
	ConfigureUserRoutes(userRouter, handlers.UserHandler, jwtMiddleware.Middleware)

	conversationRoutes := routerWithPrefix.PathPrefix("/v1/conversations").Subrouter()
	ConfigureConversationRoutes(conversationRoutes, handlers.ConversationHandler, jwtMiddleware.Middleware)

	messageRouter := routerWithPrefix.PathPrefix("/v1/messages").Subrouter()
	ConfigureMessagesRoutes(messageRouter, handlers.MessageHandler, jwtMiddleware.Middleware)

	err := routerWithPrefix.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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

	router.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	return &Server{
		httpServer: &http.Server{
			Addr:         "0.0.0.0:" + cfg.HttpPort,
			Handler:      router,
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
		},
		logger:   logger,
		handlers: handlers,
	}
}

func (s *Server) Run() error {
	s.logger.Info("Server is running", slog.String("port", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
