package telegram

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ciscapello/notification-service/application/config"
	"go.uber.org/zap"
)

type TelegramManager struct {
	logger   *zap.Logger
	client   http.Client
	botToken string
	chatId   string
}

func New(logger *zap.Logger, config config.Config) *TelegramManager {
	client := http.Client{Timeout: 10 * time.Second}

	return &TelegramManager{
		logger:   logger,
		client:   client,
		botToken: config.BotToken,
		chatId:   config.ChatId,
	}
}

func (tm *TelegramManager) SendMessage(text string) {

	baseUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s", tm.botToken, tm.chatId)

	params := url.Values{}
	params.Add("text", text)
	params.Add("parse_mode", "HTML")
	params.Add("disable_notification", "false")

	urlWithParams := fmt.Sprintf("%s&%s", baseUrl, params.Encode())

	req, err := http.NewRequest(http.MethodGet, urlWithParams, nil)
	if err != nil {
		tm.logger.Error(err.Error())
	}

	fmt.Println(urlWithParams)
	res, err := tm.client.Do(req)
	if err != nil {
		tm.logger.Error(err.Error())
	}

	if res.StatusCode != http.StatusOK {
		tm.logger.Error("failed to send message")
	}

	fmt.Println("tg message sended")
}
