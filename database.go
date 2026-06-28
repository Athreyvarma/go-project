package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectDatabase() {

	var err error

	connStr := "host=localhost port=5432 user=postgres password=Araj@1234 dbname=project sslmode=disable"

	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database Connected Successfully")
}
