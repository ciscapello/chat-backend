package main

import (
	"log/slog"
	"net/http"

	"github.com/ciscapello/socket-service/internal/application/config"
	"github.com/ciscapello/socket-service/internal/infrastructure/broadcast"
)

func main() {
	config := config.New()
	http.HandleFunc("/ws", broadcast.ServeWs)
	go broadcast.Broadcaster()

	err := http.ListenAndServe(":"+config.HttpPort, nil)
	if err != nil {
		slog.Error("cannot listen and serve on", config.HttpPort)
	}
}
