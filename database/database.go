package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	// Load environment variables from .env for local development.
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using existing environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		done BOOLEAN NOT NULL DEFAULT FALSE
	);
	`

	if _, err := db.Exec(createTable); err != nil {
		db.Close()
		return nil, err
	}

	if err := seedTasks(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func seedTasks(db *sql.DB) error {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	seedSQL := `
	INSERT INTO tasks (title, done)
	VALUES
		('Learn Go', FALSE),
		('Build Task API', FALSE),
		('Connect API to PostgreSQL', FALSE);
	`

	_, err = db.Exec(seedSQL)
	return err
}
