package main

import (
	"net/http"

	"github.com/ciscapello/socket_service/internal/application/config"
	"github.com/ciscapello/socket_service/internal/infrastructure/broadcast"
)

func main() {
	config := config.New()
	http.HandleFunc("/ws", broadcast.ServeWs)
	go broadcast.Broadcaster()

	http.ListenAndServe(":"+config.HttpPort, nil)
}
