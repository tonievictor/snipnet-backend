package database

import (
	"context"
	"database/sql"
	"time"
)

func Init(driver, datasource string) (*sql.DB, error) {
	db, err := sql.Open(driver, datasource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	err = createTables(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	query := tables()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func tables() string {
	query := `
		CREATE TABLE IF NOT EXISTS users(
			id TEXT PRIMARY KEY NOT NULL UNIQUE,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS snippets (
			id TEXT PRIMARY KEY NOT NULL UNIQUE,
			user_id TEXT NOT NULL UNIQUE,
			title TEXT NOT NULL DEFAULT 'Untitled',
			description TEXT NOT NULL DEFAULT '',
			code TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id)
		);

		CREATE TABLE IF NOT EXISTS tokens (
			id SERIAL PRIMARY KEY,
			token TEXT NOT NULL UNIQUE,
			user_id TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id)
		);
	`
	return query
}
