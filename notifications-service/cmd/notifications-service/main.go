package main

import (
	"fmt"

	"github.com/ciscapello/notification-service/application/config"
	"github.com/ciscapello/notification-service/internal/infrastructure/rabbitmq"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("hello from not service")

	conf := config.New()
	cons := rabbitmq.NewConsumer(conf, zap.NewNop())

	done := make(chan bool, 1)

	go func() {
		cons.Consume(rabbitmq.UserCreatedTopic, done)
	}()

	<-done
}
