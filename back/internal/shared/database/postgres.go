package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// NewPostgresConnection creates a new PostgreSQL database connection pool
func NewPostgresConnection(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Connection pool configuration
	db.SetMaxOpenConns(25)                 // Maximum connections
	db.SetMaxIdleConns(5)                  // Idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Connection lifetime

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
