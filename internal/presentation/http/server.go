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
)

type Server struct {
	httpServer *http.Server
}

func New(cfg *config.Config, handlers *handlers.Handlers) *Server {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowedHandler)

	userRouter := router.PathPrefix("/users/").Subrouter()
	ConfigureUserRoutes(userRouter, handlers)

	defer func() {
		fmt.Printf("Server gonna start on port %s\n", cfg.HttpPort)
	}()

	router.Handle("", userRouter)

	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
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
	}
}

func (s *Server) Run() error {
	fmt.Println("listen and serve")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func ConfigureUserRoutes(subrouter *mux.Router, handlers *handlers.Handlers) {
	subrouter.HandleFunc("", handlers.GetUsers).Methods(http.MethodGet)
	subrouter.HandleFunc("/{id}", handlers.GetUser).Methods(http.MethodGet)
	subrouter.HandleFunc("", handlers.UpdateUser).Methods(http.MethodPut)
	subrouter.HandleFunc("/registration", handlers.Registration).Methods(http.MethodPost)
	subrouter.HandleFunc("/login", handlers.Login).Methods(http.MethodPost)
}
