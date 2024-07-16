package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ciscapello/lib/contracts"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserId string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*Client]bool)
var broadcast = make(chan *contracts.MessageSocketBody)

func receiver(client *Client) {
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
		broadcast <- &body
	}
}

func broadcaster() {
	for {
		message := <-broadcast
		fmt.Println("new message in broadcast")

		for client := range clients {
			fmt.Println("client user id:", client.UserId,
				"from:", message.FromUserID,
				"to:", message.ToUserID)

			if client.UserId == message.FromUserID || client.UserId == message.ToUserID {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host, r.URL.Query())

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &Client{Conn: ws}

	clients[client] = true
	fmt.Println("clients", len(clients), clients, ws.RemoteAddr())

	receiver(client)

	fmt.Println("exiting", ws.RemoteAddr().String())
	delete(clients, client)
}

func StartWebsocketServer() {

	go broadcaster()
	http.ListenAndServe(":3002", nil)
}

func main() {
	http.HandleFunc("/ws", serveWs)

	StartWebsocketServer()
}
