package models

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDb(url string) (*sql.DB, error) {
	var err error

	db, err = sql.Open("postgres", url)
	if err != nil {
		return db, err
	}

	return db, db.Ping()
}

func StartTx() *sql.Tx {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	return tx
}
