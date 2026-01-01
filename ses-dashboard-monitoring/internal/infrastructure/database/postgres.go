package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgres(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db
}