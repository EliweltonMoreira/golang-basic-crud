package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Connection driver with postgres
)

const (
	host     = "localhost"
	port     = 5435
	user     = "devbook"
	password = "devbook00"
	dbname   = "devbook"
)

// Connect open the connection with database
func Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
