package main

import (
	"testing"
	"time"
)

func TestNewFilter(t *testing.T) {
	records := []Record{
		{Command: "ls", Arguments: "-la"},
		{Command: "grep", Arguments: "foo bar.txt"},
	}

	filter := NewFilter(records)

	if filter.Filter() != "" {
		t.Errorf("Expected empty filter text, got %s", filter.Filter())
	}

	if len(filter.FilteredRecords()) != len(records) {
		t.Errorf("Expected %d filtered records, got %d", len(records), len(filter.FilteredRecords()))
	}
}

func TestUpdateFilter(t *testing.T) {
	now := time.Now()
	records := []Record{
		{ID: 1, Command: "ls", Arguments: "-la", Timestamp: now, WorkingDirectory: "/home", ExitStatus: 0},
		{ID: 2, Command: "grep", Arguments: "foo bar.txt", Timestamp: now, WorkingDirectory: "/home", ExitStatus: 0},
		{ID: 3, Command: "find", Arguments: ". -name '*.go'", Timestamp: now, WorkingDirectory: "/home", ExitStatus: 0},
	}

	filter := NewFilter(records)

	// Test with empty filter
	filter.UpdateFilter("")
	if len(filter.FilteredRecords()) != len(records) {
		t.Errorf("Expected %d records with empty filter, got %d", len(records), len(filter.FilteredRecords()))
	}

	// Test filtering by command
	filter.UpdateFilter("grep")
	if len(filter.FilteredRecords()) != 1 {
		t.Errorf("Expected 1 record when filtering by 'grep', got %d", len(filter.FilteredRecords()))
	}
	if filter.FilteredRecords()[0].Command != "grep" {
		t.Errorf("Expected command 'grep', got '%s'", filter.FilteredRecords()[0].Command)
	}

	// Test filtering by arguments
	filter.UpdateFilter("name")
	if len(filter.FilteredRecords()) != 1 {
		t.Errorf("Expected 1 record when filtering by 'name', got %d", len(filter.FilteredRecords()))
	}
	if filter.FilteredRecords()[0].Command != "find" {
		t.Errorf("Expected command 'find', got '%s'", filter.FilteredRecords()[0].Command)
	}

	// Test case insensitivity
	filter.UpdateFilter("LS")
	if len(filter.FilteredRecords()) != 1 {
		t.Errorf("Expected 1 record when filtering by 'LS', got %d", len(filter.FilteredRecords()))
	}
	if filter.FilteredRecords()[0].Command != "ls" {
		t.Errorf("Expected command 'ls', got '%s'", filter.FilteredRecords()[0].Command)
	}

	// Test no matches
	filter.UpdateFilter("nonexistent")
	if len(filter.FilteredRecords()) != 0 {
		t.Errorf("Expected 0 records with non-matching filter, got %d", len(filter.FilteredRecords()))
	}
}

func TestTextManipulation(t *testing.T) {
	records := []Record{
		{Command: "ls", Arguments: "-la"},
		{Command: "grep", Arguments: "foo bar.txt"},
	}

	filter := NewFilter(records)

	// Test inserting text
	filter.InsertTextAtCursor("hello", 0)
	if filter.Filter() != "hello" {
		t.Errorf("Expected filter text 'hello', got '%s'", filter.Filter())
	}

	// Test inserting at specific position
	filter.InsertTextAtCursor(" world", 5)
	if filter.Filter() != "hello world" {
		t.Errorf("Expected filter text 'hello world', got '%s'", filter.Filter())
	}

	// Test removing character before cursor
	filter.RemoveCharBeforeCursor(5)
	if filter.Filter() != "hell world" {
		t.Errorf("Expected filter text 'hell world', got '%s'", filter.Filter())
	}

	// Test removing text before cursor
	filter.RemoveTextBeforeCursor(0, 4)
	if filter.Filter() != " world" {
		t.Errorf("Expected filter text ' world', got '%s'", filter.Filter())
	}

	// Test removing text after cursor
	filter.RemoveTextAfterCursor(1)
	if filter.Filter() != " " {
		t.Errorf("Expected filter text ' ', got '%s'", filter.Filter())
	}
}
