package rabbitmq

import (
	"log"

	"github.com/ciscapello/notification-service/application/config"
	emailservice "github.com/ciscapello/notification-service/internal/domain/service/emailService"
	"github.com/ciscapello/notification-service/internal/infrastructure/telegram"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Consumer struct {
	Channel         *amqp.Channel
	Connection      *amqp.Connection
	logger          *zap.Logger
	emailservice    emailservice.EmailService
	telegramManager *telegram.TelegramManager
}

func NewConsumer(config *config.Config,
	logger *zap.Logger,
	emailService *emailservice.EmailService,
	telegramManager *telegram.TelegramManager) *Consumer {
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
		emailservice:    *emailService,
		telegramManager: telegramManager,
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
