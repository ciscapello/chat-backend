package broadcast

import (
	"encoding/json"
	"fmt"

	"github.com/ciscapello/chat-lib/contracts"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserId string
}

var clients = make(map[*Client]bool)

func (c *Client) Register() {
	clients[c] = true
	fmt.Println("clients", len(clients), clients, c.Conn.RemoteAddr())
}

func (c *Client) Unregister() {
	c.Conn.Close()
	delete(clients, c)
	fmt.Println("exiting", c.Conn.RemoteAddr().String())
}

func ReceiveMessages(client *Client) {
	for {
		msgType, p, err := client.Conn.ReadMessage()
		fmt.Println("Message type:", msgType, "Payload:", string(p), "Error:", err)
		if err != nil {
			fmt.Println(err)
			return
		}

		var body contracts.MessageSocketBody
		err = json.Unmarshal(p, &body)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			continue
		}
		fmt.Println("Unmarshaled body:", body)

		fmt.Println("host", client.Conn.RemoteAddr())
		if body.Type == "bootup" {
			// do mapping on bootup
			client.UserId = body.FromUserID
			fmt.Println("client successfully mapped", &client, client, client.UserId)
			continue
		}

		fmt.Println("received message", body.Type, body.MessageBody)
		BroadcastMessage(&body)
	}
}
