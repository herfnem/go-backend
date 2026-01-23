package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // Blank import for the driver
	"log"
)

var DB *sql.DB

func initDB() {
	var err error
	// Open connection to a file named 'app.db'
	DB, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create a simple table if it doesn't exist
	statement := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		email TEXT
	);`
	_, err = DB.Exec(statement)
	if err != nil {
		log.Fatal("Table creation failed:", err)
	}
}
