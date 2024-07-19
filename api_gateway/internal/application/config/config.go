package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
)

const (
	DATABASE_URL = "DATABASE_URL"
	RABBITMQ_URL = "RABBITMQ_URL"
	HTTP_PORT    = "HTTP_PORT"
)

var ErrNoEnvs = errors.New("there's no environment variables")

type Config struct {
	DbUrl      string
	HttpPort   string
	LogPath    string
	RmqConnStr string

	AccessTokenExpTime  time.Duration
	RefreshTokenExpTime time.Duration
	AccessTokenSecret   string
	RefreshTokenSecret  string
}

func New() *Config {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal(ErrNoEnvs)
	// }
	httpPort := os.Getenv(HTTP_PORT)
	dbUrl := os.Getenv(DATABASE_URL)

	accessTokenExpTimeStr := os.Getenv("ACCESS_TOKEN_EXPIRES_IN")
	accessTokenExpTime, err := time.ParseDuration(accessTokenExpTimeStr)
	if err != nil {
		zap.Error(errors.New("invalid `access token expires in` time"))
	}
	refreshTokenExpTimeStr := os.Getenv("REFRESH_TOKEN_EXPIRES_IN")
	refreshTokenExpTime, err := time.ParseDuration(refreshTokenExpTimeStr)
	if err != nil {
		zap.Error(errors.New("invalid `refresh token expires in` time"))
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

		AccessTokenExpTime:  accessTokenExpTime,
		RefreshTokenExpTime: refreshTokenExpTime,
		AccessTokenSecret:   accessTokenSecret,
		RefreshTokenSecret:  refreshTokenSecret,
	}
}
