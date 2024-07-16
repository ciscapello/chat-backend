package broadcast

import (
	"fmt"
	"log"

	"github.com/ciscapello/lib/contracts"
)

var broadcast = make(chan *contracts.MessageSocketBody)

func BroadcastMessage(message *contracts.MessageSocketBody) {
	broadcast <- message
}

func Broadcaster() {
	for {
		message := <-broadcast
		fmt.Println("new message in broadcast")

		for cl := range clients {
			fmt.Println("client user id:", cl.UserId,
				"from:", message.FromUserID,
				"to:", message.ToUserID)

			if cl.UserId == message.FromUserID || cl.UserId == message.ToUserID {
				err := cl.Conn.WriteJSON(message)
				if err != nil {
					log.Printf("Websocket error: %s", err)
					cl.Conn.Close()
					delete(clients, cl)
				}
			}
		}
	}
}
