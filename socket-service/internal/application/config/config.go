package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const (
	DB_PASSWORD = "DB_PASSWORD"
	DB_HOST     = "DB_HOST"
	DB_PORT     = "DB_PORT"
	DB_NAME     = "DB_NAME"
	DB_USER     = "DB_USER"
	HTTP_PORT   = "HTTP_PORT"
)

var ErrNoEnvs = errors.New("there's no environment variables")

type Config struct {
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
	DbUser     string
	HttpPort   string
	LogPath    string
	RmqConnStr string

	AccessTokenExpTime  time.Duration
	RefreshTokenExpTime time.Duration
	AccessTokenSecret   string
	RefreshTokenSecret  string
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(ErrNoEnvs)
	}

	dbPassword := os.Getenv(DB_PASSWORD)
	dbHost := os.Getenv(DB_HOST)
	dbPort := os.Getenv(DB_PORT)
	dbName := os.Getenv(DB_NAME)
	dbUser := os.Getenv(DB_USER)
	httpPort := os.Getenv(HTTP_PORT)

	rmqHost := os.Getenv("RMQ_HOST")
	rmqPost := os.Getenv("RMQ_PORT")

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

	rmqConnStr := "amqp://" + "guest" + ":" + "guest" + "@" + rmqHost + ":" + rmqPost + "/"

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		DbPassword: dbPassword,
		DbHost:     dbHost,
		DbPort:     dbPort,
		DbName:     dbName,
		DbUser:     dbUser,
		HttpPort:   httpPort,
		LogPath:    path + "/logs",
		RmqConnStr: rmqConnStr,

		AccessTokenExpTime:  accessTokenExpTime,
		RefreshTokenExpTime: refreshTokenExpTime,
		AccessTokenSecret:   accessTokenSecret,
		RefreshTokenSecret:  refreshTokenSecret,
	}
}
