package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ciscapello/api_gateway/internal/application/config"
	_ "github.com/lib/pq"
)

var ErrCannotConnect = errors.New("cannot connect to database")

type Database struct {
}

func New() *Database {
	return &Database{}
}

func (d *Database) Start(config *config.Config) *sql.DB {

	db, err := sql.Open("postgres", config.DbUrl)

	if err != nil {
		log.Fatal(err, ErrCannotConnect)
	}

	if err != nil {
		log.Fatal(err, ErrCannotConnect)
	}

	fmt.Println("Successfully connected!")
	return db
}
