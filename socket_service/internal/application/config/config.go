package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	HTTP_PORT = "HTTP_PORT"
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
		log.Fatal(ErrNoEnvs)
	}

	httpPort := os.Getenv(HTTP_PORT)

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		HttpPort: httpPort,
		LogPath:  path + "/logs",
	}
}
