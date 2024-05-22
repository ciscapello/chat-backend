package main

import (
	"fmt"
	"time"

	"github.com/ciscapello/notification-service/internal/infrastructure/rabbitmq"
)

func main() {
	fmt.Println("hello from not service")

	rabbitmq.Init()

	time.Sleep(time.Hour)
}
