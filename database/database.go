package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "tasks.db")
	if err != nil {
		return nil, err
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		done INTEGER NOT NULL DEFAULT 0
	);
	`

	_, err = db.Exec(createTable)
	if err != nil {
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
		('Learn Go', 0),
		('Build Task API', 0),
		('Connect API to SQLite', 0);
	`

	_, err = db.Exec(seedSQL)
	return err
}
