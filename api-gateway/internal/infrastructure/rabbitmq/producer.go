package rabbitmq

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func (p *Producer) Publish(topic string, msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		p.logger.Error("failed to marshal message", zap.Error(err))
	}

	q, err := p.Channel.QueueDeclare(
		topic, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = p.Channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}

	log.Printf(" [x] Sent %s", string(body))
	return nil
}
