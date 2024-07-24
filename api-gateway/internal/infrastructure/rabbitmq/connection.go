package rabbitmq

import (
	"log"
	"log/slog"

	"github.com/ciscapello/api-gateway/internal/application/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	Channel    *amqp.Channel
	Connection *amqp.Connection
	logger     *slog.Logger
}

func NewProducer(config *config.Config, logger *slog.Logger) *Producer {
	var err error
	conn, err := amqp.Dial(config.RmqConnStr)
	if err != nil {
		log.Fatal(err, "cannot connect to rmq")
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err, "cannot create channel")
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
