package config

import (
	"errors"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

const (
	DB_PASSWORD  = "DB_PASSWORD"
	DB_HOST      = "DB_HOST"
	DB_PORT      = "DB_PORT"
	DB_NAME      = "DB_NAME"
	DB_USER      = "DB_USER"
	HTTP_PORT    = "HTTP_PORT"
	RABBITMQ_URL = "RABBITMQ_URL"
)

var ErrNoEnvs = errors.New("there's no environment variables")

type Config struct {
	LogPath       string
	RmqConnStr    string
	EmailAddress  string
	EmailPassword string
	BotToken      string
	ChatId        string
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		slog.Warn(ErrNoEnvs.Error())
	}

	emailAddr := os.Getenv("EMAIL_ADDRESS")
	emailPass := os.Getenv("EMAIL_PASSWORD")

	botToken := os.Getenv("BOT_TOKEN")
	chatId := os.Getenv("CHAT_ID")

	rmqConnStr := os.Getenv(RABBITMQ_URL)

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		LogPath:       path + "/logs",
		RmqConnStr:    rmqConnStr,
		EmailAddress:  emailAddr,
		EmailPassword: emailPass,
		BotToken:      botToken,
		ChatId:        chatId,
	}
}
