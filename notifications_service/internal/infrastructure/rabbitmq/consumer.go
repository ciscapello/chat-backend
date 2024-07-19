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

			var userCreatedMsg contracts.UserCreatedMessage
			log.Printf("Received in %s a message: %s", queueName, d.Body)
			err := json.Unmarshal(d.Body, &userCreatedMsg)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(userCreatedMsg)
			c.emailservice.SendCodeToUser(userCreatedMsg.Code, userCreatedMsg.Email)
			c.telegramManager.SendMessage(fmt.Sprintf("code %s", userCreatedMsg.Code))
		}
		doneCh <- true
	}()

	return nil
}
