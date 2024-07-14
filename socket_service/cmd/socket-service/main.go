package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserId string
}

type Message struct {
	Type           string `json:"type"`
	ConversationId int    `json:"conversation_id,omitempty"`
	FromUserID     string `json:"from_user_id,omitempty"`
	ToUserID       string `json:"to_user_id,omitempty"`
	MessageBody    string `json:"message_body,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*Client]bool)
var broadcast = make(chan *Message)

func receiver(client *Client) {
	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		m := &Message{}
		err = json.Unmarshal(p, m)

		if err != nil {
			fmt.Println("error while unmarshaling chat", err)
			continue
		}

		fmt.Println("host", client.Conn.RemoteAddr())
		if m.Type == "bootup" {
			// do mapping on bootup
			client.UserId = m.FromUserID
			fmt.Println("client successfully mapped", &client, client, client.UserId)
			continue
		}

		fmt.Println("received message", m.Type, m.MessageBody)
		broadcast <- m
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
