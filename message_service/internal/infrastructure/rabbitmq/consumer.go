package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ciscapello/chat-lib/contracts"
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
			var body contracts.MessageSocketBody

			fmt.Println(d.Body)

			err := json.Unmarshal(d.Body, &body)
			if err != nil {
				log.Fatal(err)
			}

			c.MessagesService.CreateMessage(body.FromUserID, body.ToUserID, body.ConversationId, body.MessageBody)
		}
		doneCh <- true
	}()

	return nil
}
