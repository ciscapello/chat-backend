package server

import (
	"net/http"

	"github.com/ciscapello/socket-service/internal/infrastructure/broadcast"
)

func StartWebsocketServer() {

	go broadcast.Broadcaster()
	http.ListenAndServe(":3002", nil)
}
