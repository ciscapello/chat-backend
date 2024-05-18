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
	"github.com/ciscapello/chat-backend/internal/presentation/handlers"
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

	userService := userservice.New()

	handlers := handlers.New(userService)

	server := httpServer.New(config, handlers, logger)

	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
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

	if err := server.Stop(ctx); err != nil {
		fmt.Printf("failed to stop server: %v", err)
	}

	if err := database.Close(); err != nil {
		fmt.Println(err.Error())
	}
}
