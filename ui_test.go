package main_test

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	rt "github.com/nuchs/retour"
)

func TestNew(t *testing.T) {
	// Create test data
	records := []rt.Record{
		{
			Command:          "ls",
			Arguments:        "-la",
			Timestamp:        time.Now().Add(-1 * time.Hour),
			WorkingDirectory: "/home/user",
			ExitStatus:       0,
		},
		{
			Command:          "git",
			Arguments:        "status",
			Timestamp:        time.Now().Add(-2 * time.Hour),
			WorkingDirectory: "/home/user/project",
			ExitStatus:       0,
		},
		{
			Command:          "make",
			Arguments:        "build",
			Timestamp:        time.Now().Add(-3 * time.Hour),
			WorkingDirectory: "/home/user/project",
			ExitStatus:       1,
		},
	}

	// Create model
	model := rt.NewUI(records)

	// Verify initial state
	if len(model.Records()) != len(records) {
		t.Errorf("Expected %d records, got %d", len(records), len(model.Records()))
	}

	// Verify no selection initially
	if _, ok := model.Selected(); ok {
		t.Error("Expected no selection initially")
	}
}

func TestNavigation(t *testing.T) {
	// Create test data with just two records for simplicity
	records := []rt.Record{
		{
			Command:    "first",
			ExitStatus: 0,
		},
		{
			Command:    "second",
			ExitStatus: 0,
		},
	}

	model := rt.NewUI(records)

	// Test down navigation
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	m := newModel.(rt.Model)
	if m.Cursor() != 1 {
		t.Errorf("Expected cursor at 1 after down, got %d", m.Cursor())
	}

	// Test up navigation
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = newModel.(rt.Model)
	if m.Cursor() != 0 {
		t.Errorf("Expected cursor at 0 after up, got %d", m.Cursor())
	}
}

func TestSelection(t *testing.T) {
	records := []rt.Record{
		{
			Command:    "test",
			ExitStatus: 0,
		},
	}

	model := rt.NewUI(records)

	// Select the item
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := newModel.(rt.Model)

	// Verify selection
	record, ok := m.Selected()
	if !ok {
		t.Error("Expected selection after Enter")
	}
	if record.Command != "test" {
		t.Errorf("Expected selected command 'test', got '%s'", record.Command)
	}
}

func TestFilterStub(t *testing.T) {
	records := []rt.Record{
		{
			Command:    "test1",
			ExitStatus: 0,
		},
		{
			Command:    "test2",
			ExitStatus: 0,
		},
	}

	model := rt.NewUI(records)

	// Add some filter text
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("test")})
	m := newModel.(rt.Model)

	// Verify that filtering is currently a no-op
	if len(m.Records()) != len(records) {
		t.Error("Expected no-op filter to return all records")
	}
}
