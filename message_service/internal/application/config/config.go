package config

import (
	"errors"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

const (
	DATABASE_URL       = "DATABASE_URL"
	RABBITMQ_URL       = "RABBITMQ_URL"
	HTTP_PORT          = "HTTP_PORT"
	SOCKET_SERVICE_URL = "SOCKET_SERVICE_URL"
)

var ErrNoEnvs = errors.New("there's no environment variables")

type Config struct {
	DbUrl             string
	HttpPort          string
	LogPath           string
	RmqConnStr        string
	EmailAddress      string
	EmailPassword     string
	BotToken          string
	ChatId            string
	SocketServicePort string
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		slog.Warn(ErrNoEnvs.Error())
	}

	dbUrl := os.Getenv(DATABASE_URL)
	httpPort := os.Getenv(HTTP_PORT)
	socketServicePort := os.Getenv(SOCKET_SERVICE_URL)

	botToken := os.Getenv("BOT_TOKEN")
	chatId := os.Getenv("CHAT_ID")

	rmqConnStr := os.Getenv(RABBITMQ_URL)

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		DbUrl:             dbUrl,
		HttpPort:          httpPort,
		LogPath:           path + "/logs",
		RmqConnStr:        rmqConnStr,
		BotToken:          botToken,
		ChatId:            chatId,
		SocketServicePort: socketServicePort,
	}
}
