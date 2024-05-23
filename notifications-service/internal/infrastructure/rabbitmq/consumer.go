package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
)

func (c *Consumer) Consume(queueName string, doneCh chan<- bool) error {
	q, err := c.Channel.QueueDeclare(
		UserCreatedTopic, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
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

			var userCreatedMsg UserCreatedMessage
			log.Printf("Received in %s a message: %s", UserCreatedTopic, d.Body)
			err := json.Unmarshal(d.Body, &userCreatedMsg)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(userCreatedMsg)
		}
		doneCh <- true
	}()

	return nil
}
