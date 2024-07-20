package config

import (
	"errors"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

const (
	SOCKET_SERVICE_HTTP_PORT = "SOCKET_SERVICE_HTTP_PORT"
)

var ErrNoEnvs = errors.New("there's no environment variables")

type Config struct {
	HttpPort   string
	LogPath    string
	RmqConnStr string
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		slog.Warn(ErrNoEnvs.Error())
	}

	httpPort := os.Getenv(SOCKET_SERVICE_HTTP_PORT)

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		HttpPort: httpPort,
		LogPath:  path + "/logs",
	}
}
