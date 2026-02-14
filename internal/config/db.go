package config

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

func OpenDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func Migrate(ctx context.Context, db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS posts (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			content TEXT,
			user_id TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS monitors (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			name TEXT NOT NULL,
			url TEXT NOT NULL,
			interval_seconds INTEGER DEFAULT 300,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS monitor_logs (
			id TEXT PRIMARY KEY,
			monitor_id TEXT NOT NULL,
			status TEXT NOT NULL,
			status_code INTEGER,
			response_time_ms INTEGER,
			error_message TEXT,
			checked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (monitor_id) REFERENCES monitors(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS snippets (
			id TEXT PRIMARY KEY,
			hash TEXT NOT NULL UNIQUE,
			content TEXT NOT NULL,
			password TEXT,
			burn_after_read BOOLEAN DEFAULT 0,
			expires_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, table := range tables {
		if _, err := db.ExecContext(ctx, table); err != nil {
			return err
		}
	}

	return nil
}
