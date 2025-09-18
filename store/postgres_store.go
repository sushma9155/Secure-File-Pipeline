package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Postgres driver
)

type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new Postgres DB connection
func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	// Verify the connection works
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping Postgres: %w", err)
	}

	// Create a table if not exists
	schema := `
	CREATE TABLE IF NOT EXISTS records (
		id SERIAL PRIMARY KEY,
		source     TEXT,
		filename   TEXT,
		data       TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

// SaveRecord inserts a record into the DB
func (s *PostgresStore) SaveRecord(source, filename, data string) error {
	_, err := s.db.Exec(
		`INSERT INTO records (source, filename, data) VALUES ($1, $2, $3)`,
		source, filename, data,
	)
	return err
}

// Close closes the DB connection
func (s *PostgresStore) Close() error {
	return s.db.Close()
}
