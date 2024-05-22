package rabbitmq

import (
	"log"

	"github.com/ciscapello/api-gateway/internal/application/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var ch *amqp.Channel

const (
	UserCreatedTopic = "user.created"
)

func Init(config *config.Config) {
	var err error
	conn, err = amqp.Dial(config.RmqConnStr)
	if err != nil {
		log.Fatal(err)
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
}

func Close() {
	if ch != nil {
		ch.Close()
	}
	if conn != nil {
		conn.Close()
	}
}

func GetChannel() *amqp.Channel {
	return ch
}
