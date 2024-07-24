package wsClient

import (
	"encoding/json"
	"fmt"

	"github.com/ciscapello/api-gateway/internal/application/config"
	"github.com/ciscapello/chat-lib/contracts"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	conn *websocket.Conn
}

func New(conf config.Config) (*WsClient, error) {
	fmt.Println(conf.SocketUrl)
	conn, _, err := websocket.DefaultDialer.Dial(conf.SocketUrl, nil)
	if err != nil {
		return nil, err
	}

	return &WsClient{
		conn: conn,
	}, nil
}

func (ws *WsClient) SendMessage(body contracts.MessageSocketBody) error {

	bytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	err = ws.conn.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		return err
	}

	return nil
}
