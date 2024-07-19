package wsClient

import (
	"encoding/json"
	"fmt"

	"github.com/ciscapello/chat-lib/contracts"
	"github.com/ciscapello/message_service/internal/application/config"
	"github.com/gorilla/websocket"
)

type WsClient struct {
	conn *websocket.Conn
}

type Message struct {
	Type           string `json:"type"`
	ConversationId int    `json:"conversation_id,omitempty"`
	FromUserID     string `json:"from_user_id,omitempty"`
	ToUserID       string `json:"to_user_id,omitempty"`
	MessageBody    string `json:"message_body,omitempty"`
}

func New(conf config.Config) (*WsClient, error) {
	fmt.Println(conf.SocketServicePort)
	conn, _, err := websocket.DefaultDialer.Dial(conf.SocketServicePort, nil)
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
