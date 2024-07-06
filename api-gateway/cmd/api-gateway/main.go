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

	"github.com/ciscapello/api-gateway/internal/application/config"
	"github.com/ciscapello/api-gateway/internal/application/db"
	"github.com/ciscapello/api-gateway/internal/common/jwtmanager"
	"github.com/ciscapello/api-gateway/internal/common/logger"
	conversationservice "github.com/ciscapello/api-gateway/internal/domain/service/conversation_service"
	messageservice "github.com/ciscapello/api-gateway/internal/domain/service/message_service"
	userservice "github.com/ciscapello/api-gateway/internal/domain/service/user_service"
	"github.com/ciscapello/api-gateway/internal/infrastructure/rabbitmq"
	"github.com/ciscapello/api-gateway/internal/infrastructure/repository"
	conversationhandler "github.com/ciscapello/api-gateway/internal/presentation/handlers/conversation_handler"
	defaulthandler "github.com/ciscapello/api-gateway/internal/presentation/handlers/default_handler"
	messagehandler "github.com/ciscapello/api-gateway/internal/presentation/handlers/message_handler"
	userhandler "github.com/ciscapello/api-gateway/internal/presentation/handlers/user_handler"
	httpServer "github.com/ciscapello/api-gateway/internal/presentation/http"
	"github.com/ciscapello/api-gateway/internal/presentation/response"

	_ "github.com/ciscapello/api-gateway/docs"
	_ "github.com/swaggo/http-swagger" // http-swagger middleware
)

// @title API Gateway
// @version 1.0

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	run()
}

func run() {
	config := config.New()

	logger := logger.GetLogger(config)
	defer logger.Sync()

	db := db.New()
	database := db.Start(config)

	producer := rabbitmq.NewProducer(config, logger)
	defer producer.Close()

	userRepository := repository.NewUserRepository(database, logger)
	conversationRepo := repository.NewConversationRepository(database, logger)
	messagesRepo := repository.NewMessagesRepository(database, logger)

	jwtMan := jwtmanager.NewJwtManager(config, logger)

	userService := userservice.New(userRepository, logger, producer, jwtMan)
	conversationService := conversationservice.New(conversationRepo, logger, jwtMan)
	messagesService := messageservice.New(messagesRepo, conversationRepo, logger, jwtMan)

	responder := response.Responder{}

	userHandler := userhandler.New(userService, logger, jwtMan, responder)
	defaulthandler := defaulthandler.New(logger)
	conversationhandler := conversationhandler.New(conversationService, logger, jwtMan, responder)
	messagesHandler := messagehandler.New(messagesService, logger, jwtMan, responder)

	httpServer := httpServer.New(config, &httpServer.Handlers{
		UserHandler:         userHandler,
		DefaultHandler:      defaulthandler,
		ConversationHandler: conversationhandler,
		MessageHandler:      messagesHandler,
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
