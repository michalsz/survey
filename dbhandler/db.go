package dbhandler

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectToDB() error {
	var err error
	db, err = sql.Open("postgres",
		"postgres://postgres:@localhost:5432/test?sslmode=disable")
	return err
}

func GetDatabase() *sql.DB {
	return db
}
