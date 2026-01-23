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
	tables := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		// Posts table
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			content TEXT,
			user_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`,

		// Monitors table (Uptime Ninja)
		`CREATE TABLE IF NOT EXISTS monitors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			url TEXT NOT NULL,
			interval_seconds INTEGER DEFAULT 300,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`,

		// Monitor logs table
		`CREATE TABLE IF NOT EXISTS monitor_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			monitor_id INTEGER NOT NULL,
			status TEXT NOT NULL,
			status_code INTEGER,
			response_time_ms INTEGER,
			error_message TEXT,
			checked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (monitor_id) REFERENCES monitors(id) ON DELETE CASCADE
		);`,

		// Snippets table (Pastebin)
		`CREATE TABLE IF NOT EXISTS snippets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			hash TEXT NOT NULL UNIQUE,
			content TEXT NOT NULL,
			password TEXT,
			burn_after_read BOOLEAN DEFAULT 0,
			expires_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			log.Fatal("Failed to create table:", err)
		}
	}
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
