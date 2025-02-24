package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Record represents a single command history entry
type Record struct {
	ID              int64
	Command         string
	Timestamp       time.Time
	WorkingDirectory string
	ExitStatus      int
	Arguments       string
}

// DB wraps the SQLite database connection
type DB struct {
	conn *sql.DB
}

// New creates a new database connection and ensures the schema is set up
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

// Close closes the database connection
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

// Insert adds a new command record to the database
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

// Query executes a custom SQL query and returns the results
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

// QueryFiltered returns records based on the provided filters
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
