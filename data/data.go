package data

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "dbname=webapp_template sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}
