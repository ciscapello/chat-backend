package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ciscapello/api-gateway/internal/application/config"
	_ "github.com/lib/pq"
)

var ErrCannotConnect = errors.New("cannot connect to database")

type Database struct {
}

func New() *Database {
	return &Database{}
}

func (d *Database) Start(config *config.Config) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.DbHost, config.DbPort, config.DbUser, config.DbPassword, config.DbName)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err, ErrCannotConnect)
	}

	if err != nil {
		log.Fatal(err, ErrCannotConnect)
	}

	fmt.Println("Successfully connected!")
	return db
}
