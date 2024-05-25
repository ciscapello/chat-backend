package main

import (
	"github.com/ciscapello/notification-service/application/config"
	"github.com/ciscapello/notification-service/common/logger"
	emailservice "github.com/ciscapello/notification-service/internal/domain/service/emailService"
	"github.com/ciscapello/notification-service/internal/infrastructure/rabbitmq"
)

func main() {

	conf := config.New()

	logger := logger.GetLogger(conf)

	es := emailservice.New(*conf, logger)
	cons := rabbitmq.NewConsumer(conf, logger, es)

	never := make(chan bool, 1)

	go func() {
		cons.Consume(rabbitmq.UserCreatedTopic, never)
	}()

	<-never
}
