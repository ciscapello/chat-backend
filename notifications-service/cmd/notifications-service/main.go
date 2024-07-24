package main

import (
	"github.com/ciscapello/chat-lib/contracts"
	"github.com/ciscapello/notification-service/application/config"
	"github.com/ciscapello/notification-service/common/logger"
	emailservice "github.com/ciscapello/notification-service/internal/domain/service/emailService"
	"github.com/ciscapello/notification-service/internal/infrastructure/rabbitmq"
	"github.com/ciscapello/notification-service/internal/infrastructure/telegram"
)

func main() {

	conf := config.New()

	logger := logger.GetLogger(conf)

	es := emailservice.New(*conf, logger)
	telegramManager := telegram.New(logger, *conf)
	cons := rabbitmq.NewConsumer(conf, logger, es, telegramManager)

	never := make(chan bool, 1)

	go func() {
		cons.Consume(contracts.UserCreatedTopic, never)
	}()

	<-never
}
