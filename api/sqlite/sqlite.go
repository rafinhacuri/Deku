package sqlite

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var SQL *sql.DB

func InitSqlite() (*sql.DB, error) {
	var err error
	SQL, err = sql.Open("sqlite3", os.Getenv("SQLITE_PATH"))
	if err != nil {
		return nil, err
	}

	_, err = SQL.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS salary (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user INTEGER NOT NULL,
			value REAL NOT NULL,
			type TEXT NOT NULL,
			day INTEGER,
			month TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			
			FOREIGN KEY (user) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS expenses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user INTEGER NOT NULL,
			value REAL NOT NULL,
			type TEXT NOT NULL,
			day INTEGER,
			month TEXT NOT NULL,
			paymentMethod TEXT NOT NULL,
			description TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			
			FOREIGN KEY (user) REFERENCES users(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return nil, err
	}

	return SQL, nil
}

func TestSqlite(ctx context.Context) error {
	if err := SQL.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
