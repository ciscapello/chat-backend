package main

import (
	"fmt"
	"log/slog"

	"github.com/ciscapello/lib/contracts"
	"github.com/ciscapello/message_service/internal/application/config"
	"github.com/ciscapello/message_service/internal/application/db"
	"github.com/ciscapello/message_service/internal/common/logger"
	messageservice "github.com/ciscapello/message_service/internal/domain/service"
	"github.com/ciscapello/message_service/internal/infrastructure/rabbitmq"
	"github.com/ciscapello/message_service/internal/infrastructure/repository"
	"github.com/ciscapello/message_service/internal/infrastructure/wsClient"
)

func main() {
	fmt.Println("message service")

	config := config.New()

	wsConn, err := wsClient.New(*config)
	if err != nil {
		slog.Error("cannot create websocket connection")
		panic(err)
	}

	logger := logger.GetLogger(config)
	defer logger.Sync()

	db := db.New()
	database := db.Start(config)

	messagesRepo := repository.NewMessagesRepository(database, logger)

	messagesService := messageservice.New(messagesRepo, wsConn, logger)

	consumer := rabbitmq.NewConsumer(config, logger, messagesService)

	never := make(chan bool, 1)

	go func() {
		consumer.Consume(contracts.MessageCreatedTopic, never)
	}()

	<-never

}