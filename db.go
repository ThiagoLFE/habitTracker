package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func openDB() (*sql.DB, error) {

	// Initializing DB
	db, err := sql.Open("sqlite3", "./app.db")

	if err != nil {
		return nil, err
	}

	// We ping the database to establish the connection, connection up not means that's is already to work
	if err := db.Ping(); err != nil {
		// If not already, close the data and return not is possible to use the database
		db.Close()
		return nil, err
	}

	// We using just one connection per server, no multiple people can use in same time
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Use pragma to config our database and allowed
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func initSchema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS habits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		done BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(query)
	return err
}
