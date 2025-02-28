package main_test

import (
	"os"
	"testing"
	"time"

	rt "github.com/nuchs/retour"
)

func TestDB(t *testing.T) {
	// Create a temporary database file
	tmpFile, err := os.CreateTemp("", "retour-test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Open the database
	database, err := rt.NewDB(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()

	// Test inserting a record
	record := &rt.Record{
		Command:          "ls -la",
		Timestamp:        time.Now(),
		WorkingDirectory: "/home/user",
		ExitStatus:       0,
		Arguments:        "-la",
	}

	if err := database.Insert(record); err != nil {
		t.Errorf("Failed to insert record: %v", err)
	}

	// Test querying records
	records, err := database.Query("SELECT * FROM history")
	if err != nil {
		t.Errorf("Failed to query records: %v", err)
	}

	if len(records) != 1 {
		t.Errorf("Expected 1 record, got %d", len(records))
	}

	if records[0].Command != record.Command {
		t.Errorf("Expected command %q, got %q", record.Command, records[0].Command)
	}

	// Test filtered query
	records, err = database.QueryFiltered(24*time.Hour, "success", "/home/user", 10)
	if err != nil {
		t.Errorf("Failed to query filtered records: %v", err)
	}

	if len(records) != 1 {
		t.Errorf("Expected 1 record, got %d", len(records))
	}

	// Test no results
	records, err = database.QueryFiltered(24*time.Hour, "failed", "/home/user", 10)
	if err != nil {
		t.Errorf("Failed to query filtered records: %v", err)
	}

	if len(records) != 0 {
		t.Errorf("Expected 0 records, got %d", len(records))
	}
}
