package config

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

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
	DbUrl      string
	HttpPort   string
	LogPath    string
	RmqConnStr string
	SocketUrl  string

	AccessTokenExpTime  time.Duration
	RefreshTokenExpTime time.Duration
	AccessTokenSecret   string
	RefreshTokenSecret  string
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		slog.Warn(ErrNoEnvs.Error())
	}
	httpPort := os.Getenv(HTTP_PORT)
	socketUrl := os.Getenv(SOCKET_SERVICE_URL)
	dbUrl := os.Getenv(DATABASE_URL)

	accessTokenExpTimeStr := os.Getenv("ACCESS_TOKEN_EXPIRES_IN")
	accessTokenExpTime, err := time.ParseDuration(accessTokenExpTimeStr)
	if err != nil {
		slog.Error(errors.New("invalid `access token expires in` time").Error())
	}
	refreshTokenExpTimeStr := os.Getenv("REFRESH_TOKEN_EXPIRES_IN")
	refreshTokenExpTime, err := time.ParseDuration(refreshTokenExpTimeStr)
	if err != nil {
		slog.Error(errors.New("invalid `refresh token expires in` time").Error())
	}
	fmt.Println(accessTokenExpTime)

	accessTokenSecret := os.Getenv("ACCESS_JWT_SECRET")
	refreshTokenSecret := os.Getenv("REFRESH_JWT_SECRET")

	rmqConnStr := os.Getenv("RABBITMQ_URL")

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		DbUrl:      dbUrl,
		HttpPort:   httpPort,
		LogPath:    path + "/logs",
		RmqConnStr: rmqConnStr,
		SocketUrl:  socketUrl,

		AccessTokenExpTime:  accessTokenExpTime,
		RefreshTokenExpTime: refreshTokenExpTime,
		AccessTokenSecret:   accessTokenSecret,
		RefreshTokenSecret:  refreshTokenSecret,
	}
}
