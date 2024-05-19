package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ciscapello/chat-backend/internal/application/config"
	"github.com/ciscapello/chat-backend/internal/application/db"
	"github.com/ciscapello/chat-backend/internal/common/logger"
	userservice "github.com/ciscapello/chat-backend/internal/domain/service/userService"
	"github.com/ciscapello/chat-backend/internal/infrastructure/repository"
	defaulthandler "github.com/ciscapello/chat-backend/internal/presentation/handlers/defaultHandler"
	userhandler "github.com/ciscapello/chat-backend/internal/presentation/handlers/userHandler"
	httpServer "github.com/ciscapello/chat-backend/internal/presentation/http"
)

func main() {
	run()
}

func run() {
	config := config.New()

	logger := logger.GetLogger(config)
	defer logger.Sync()

	db := db.New()
	database := db.Start(config)

	userRepository := repository.NewUserRepository(database, logger)

	userService := userservice.New(userRepository)

	userHandler := userhandler.New(userService, logger)
	defaulthandler := defaulthandler.New(logger)

	httpServer := httpServer.New(config, &httpServer.Handlers{
		UserHandler:    userHandler,
		DefaultHandler: defaulthandler,
	}, logger)

	go func() {
		if err := httpServer.Run(); !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	fmt.Println("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := httpServer.Stop(ctx); err != nil {
		fmt.Printf("failed to stop server: %v", err)
	}

	if err := database.Close(); err != nil {
		fmt.Println(err.Error())
	}
}
