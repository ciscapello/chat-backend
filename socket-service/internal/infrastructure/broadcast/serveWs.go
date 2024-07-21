package broadcast

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host, r.URL.Query())

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &Client{Conn: ws}

	clients[client] = true
	fmt.Println("clients", len(clients), clients, ws.RemoteAddr())

	ReceiveMessages(client)

	fmt.Println("exiting", ws.RemoteAddr().String())
	delete(clients, client)
}
