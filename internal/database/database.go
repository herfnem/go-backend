package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"learn/internal/config"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite3", config.AppConfig.DBPath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	createTables()
	log.Println("Database initialized successfully")
}

func createTables() {
	usersTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := DB.Exec(usersTable); err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	postsTable := `CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		slug TEXT NOT NULL UNIQUE,
		content TEXT,
		user_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	if _, err := DB.Exec(postsTable); err != nil {
		log.Fatal("Failed to create posts table:", err)
	}
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
