package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgres(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	// Run migration
	if err := runMigration(db); err != nil {
		log.Fatal("Failed to run migration:", err)
	}

	return db
}

func runMigration(db *sql.DB) error {
	// Create ses_events table
	sesEventsQuery := `
		CREATE TABLE IF NOT EXISTS ses_events (
			id BIGSERIAL PRIMARY KEY,
			message_id VARCHAR(100),
			email VARCHAR(255),
			event_type VARCHAR(50),
			status VARCHAR(50),
			reason TEXT,
			created_at TIMESTAMP DEFAULT now()
		);
	`
	if _, err := db.Exec(sesEventsQuery); err != nil {
		return err
	}

	// Create suppressions table
	suppressionsQuery := `
		CREATE TABLE IF NOT EXISTS suppressions (
			email VARCHAR(255) PRIMARY KEY,
			reason TEXT NOT NULL,
			source VARCHAR(50) NOT NULL DEFAULT 'AWS',
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_suppressions_source ON suppressions(source);
		CREATE INDEX IF NOT EXISTS idx_suppressions_updated_at ON suppressions(updated_at);
	`
	if _, err := db.Exec(suppressionsQuery); err != nil {
		return err
	}

	return nil
}
