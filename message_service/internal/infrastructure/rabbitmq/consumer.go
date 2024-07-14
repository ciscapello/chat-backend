package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/ciscapello/lib/contracts"
)

func (c *Consumer) Consume(queueName string, doneCh chan<- bool) error {
	q, err := c.Channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := c.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for d := range msgs {

			var body contracts.MessageCreatedBody

			err := json.Unmarshal(d.Body, &body)
			if err != nil {
				log.Fatal(err)
			}

		}
		doneCh <- true
	}()

	return nil
}
