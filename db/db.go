// Package db provides SQLite-based storage for shell command history.
//
// It handles the persistence of command execution records, including their
// timestamps, working directories, and exit statuses. The package provides
// a simple interface for storing and querying command history, with support
// for filtering by time range, working directory, and command success/failure.
package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Record represents a single command history entry in the database.
// Each Record contains the full context of a command's execution,
// including when and where it was run, and whether it succeeded.
type Record struct {
	// ID is the unique identifier for this record in the database
	ID int64

	// Command is the main command that was executed, without arguments
	Command string

	// Timestamp records when the command was executed
	Timestamp time.Time

	// WorkingDirectory is the directory from which the command was run
	WorkingDirectory string

	// ExitStatus is the command's exit code (0 for success, non-zero for failure)
	ExitStatus int

	// Arguments contains any additional arguments passed to the command
	Arguments string
}

// DB provides an interface to the SQLite database storing command history.
// It handles connection management, schema creation, and provides methods
// for storing and querying command records.
type DB struct {
	conn *sql.DB
}

// New creates a new database connection and ensures the schema is set up.
// It takes a connectionString parameter which should be a valid SQLite
// database path. The function will create the database file if it doesn't
// exist and set up the necessary tables and indexes.
//
// Returns a new DB instance or an error if the connection or schema
// creation fails.
func New(connectionString string) (*DB, error) {
	conn, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.ensureSchema(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ensure schema: %w", err)
	}

	return db, nil
}

// Close closes the database connection and releases any associated resources.
// It should be called when the database is no longer needed to prevent
// resource leaks.
func (db *DB) Close() error {
	return db.conn.Close()
}

// ensureSchema creates the necessary tables and indexes if they don't exist
func (db *DB) ensureSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		command TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		working_directory TEXT,
		exit_status INTEGER NOT NULL,
		arguments TEXT
	);
	
	CREATE INDEX IF NOT EXISTS idx_command ON history(command);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON history(timestamp);
	CREATE INDEX IF NOT EXISTS idx_working_directory ON history(working_directory);
	`

	_, err := db.conn.Exec(schema)
	return err
}

// Insert adds a new command record to the database.
// The Record should contain all required fields: Command, Timestamp,
// WorkingDirectory, ExitStatus, and optionally Arguments.
// The ID field will be automatically set by the database.
//
// Returns an error if the insert operation fails.
func (db *DB) Insert(record *Record) error {
	query := `
	INSERT INTO history (command, timestamp, working_directory, exit_status, arguments)
	VALUES (?, ?, ?, ?, ?)
	`
	
	_, err := db.conn.Exec(query,
		record.Command,
		record.Timestamp,
		record.WorkingDirectory,
		record.ExitStatus,
		record.Arguments,
	)
	
	return err
}

// Query executes a custom SQL query and returns the results as a slice of Records.
// This method allows for custom queries beyond the standard filters provided by
// QueryFiltered. The query must return all fields of the history table in the
// correct order (id, command, timestamp, working_directory, exit_status, arguments).
//
// The args parameter allows for safe parameterization of the query.
// Returns the matching records or an error if the query fails.
func (db *DB) Query(query string, args ...interface{}) ([]Record, error) {
	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var r Record
		err := rows.Scan(
			&r.ID,
			&r.Command,
			&r.Timestamp,
			&r.WorkingDirectory,
			&r.ExitStatus,
			&r.Arguments,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}

	return records, rows.Err()
}

// QueryFiltered returns records based on the provided filters.
// It provides a high-level interface for common query patterns:
//
// - timeRange: how far back to look (e.g., 24h for last day)
// - resultFilter: filter by command success/failure ("success", "failed", "all")
// - workingDir: filter by specific working directory (empty string for all)
// - limit: maximum number of records to return
//
// Returns matching records ordered by timestamp (newest first) or an error if the query fails.
func (db *DB) QueryFiltered(timeRange time.Duration, resultFilter string, workingDir string, limit int) ([]Record, error) {
	query := `
	SELECT id, command, timestamp, working_directory, exit_status, arguments
	FROM history
	WHERE 1=1
	`
	var args []interface{}

	if timeRange > 0 {
		query += " AND timestamp >= ?"
		args = append(args, time.Now().Add(-timeRange))
	}

	if workingDir != "" {
		query += " AND working_directory = ?"
		args = append(args, workingDir)
	}

	switch resultFilter {
	case "success":
		query += " AND exit_status = 0"
	case "failed":
		query += " AND exit_status != 0"
	}

	query += " ORDER BY timestamp DESC"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	return db.Query(query, args...)
}
