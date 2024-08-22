package main

import (
	"database/sql"
	"log"
	"time"
)

// Globally shared database connection in this application.
var DB *sql.DB

func getDB(constr string, readonly bool) *sql.DB {
	db, err := sql.Open("mysql", constr)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err.Error())
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if readonly {
		_, err = db.Exec("SET SESSION TRANSACTION READ ONLY")
		if err != nil {
			log.Fatalf("failed to set read-only session: %v", err.Error())
		}
	}

	return db
}
