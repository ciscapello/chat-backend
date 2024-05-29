package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ciscapello/api-gateway/internal/application/config"
	defaulthandler "github.com/ciscapello/api-gateway/internal/presentation/handlers/defaultHandler"
	userhandler "github.com/ciscapello/api-gateway/internal/presentation/handlers/userHandler"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type Handlers struct {
	UserHandler    *userhandler.UserHandler
	DefaultHandler *defaulthandler.DefaultHandler
}

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
	handlers   *Handlers
}

func New(cfg *config.Config, handlers *Handlers, logger *zap.Logger) *Server {
	router := mux.NewRouter()

	// wd, err := os.Getwd()
	// if err != nil {
	// logger.Fatal("Error getting current working directory", zap.Error(err))
	// }

	router.Use(LoggingMiddleware(logger))

	router.NotFoundHandler = http.HandlerFunc(handlers.DefaultHandler.NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.DefaultHandler.MethodNotAllowedHandler)

	userRouter := router.PathPrefix("/users").Subrouter()
	ConfigureUserRoutes(userRouter, handlers.UserHandler)

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
	s.logger.Info("Server is running", zap.String("port", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func ConfigureUserRoutes(subrouter *mux.Router, handlers *userhandler.UserHandler) {
	subrouter.HandleFunc("", handlers.GetAllUsers).Methods(http.MethodGet)
	subrouter.HandleFunc("/{id}", handlers.GetUser).Methods(http.MethodGet)
	subrouter.HandleFunc("/{id}", handlers.UpdateUser).Methods(http.MethodPut)
	subrouter.HandleFunc("/auth", handlers.Auth).Methods(http.MethodPost)
	subrouter.HandleFunc("/check-code", handlers.CheckCode).Methods(http.MethodPost)
}
