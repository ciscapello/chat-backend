package rabbitmq

import (
	"log"

	"github.com/ciscapello/api-gateway/internal/application/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Producer struct {
	Channel    *amqp.Channel
	Connection *amqp.Connection
	logger     *zap.Logger
}

func NewProducer(config *config.Config, logger *zap.Logger) *Producer {
	var err error
	conn, err := amqp.Dial(config.RmqConnStr)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	return &Producer{
		Channel:    ch,
		Connection: conn,
		logger:     logger,
	}
}

func (p *Producer) Close() {
	if p.Channel != nil {
		p.Channel.Close()
	}
	if p.Connection != nil {
		p.Connection.Close()
	}
}
