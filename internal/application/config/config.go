package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
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

	return &Config{
		DbPassword: dbPassword,
		DbHost:     dbHost,
		DbPort:     dbPort,
		DbName:     dbName,
		DbUser:     dbUser,
		HttpPort:   httpPort,
	}
}
