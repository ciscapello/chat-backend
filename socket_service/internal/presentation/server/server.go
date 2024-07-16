package server

import (
	"net/http"

	"github.com/ciscapello/socket_service/internal/infrastructure/broadcast"
)

func StartWebsocketServer() {

	go broadcast.Broadcaster()
	http.ListenAndServe(":3002", nil)
}
