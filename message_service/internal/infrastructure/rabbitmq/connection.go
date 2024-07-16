package rabbitmq

import (
	"log"

	"github.com/ciscapello/message_service/internal/application/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type IMessagesService interface {
	CreateMessage(senderId string, receiverId string, conversationId int, message string) error
}

type Consumer struct {
	Channel         *amqp.Channel
	Connection      *amqp.Connection
	logger          *zap.Logger
	MessagesService IMessagesService
}

func NewConsumer(config *config.Config,
	logger *zap.Logger,
	messagesService IMessagesService,
) *Consumer {
	conn, err := amqp.Dial(config.RmqConnStr)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	return &Consumer{
		Channel:         ch,
		Connection:      conn,
		logger:          logger,
		MessagesService: messagesService,
	}
}

func Init() {
}

func (c *Consumer) Close() {
	if c.Channel != nil {
		c.Channel.Close()
	}
	if c.Connection != nil {
		c.Connection.Close()
	}
}
