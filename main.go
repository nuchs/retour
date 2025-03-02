package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Create some test data
	records := []Record{
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
		{
			Command:          "vim",
			Arguments:        "main.go",
			Timestamp:        time.Now().Add(-4 * time.Hour),
			WorkingDirectory: "/home/user/project",
			ExitStatus:       0,
		},
		{
			Command:          "go",
			Arguments:        "test ./...",
			Timestamp:        time.Now().Add(-5 * time.Hour),
			WorkingDirectory: "/home/user/project",
			ExitStatus:       0,
		},
	}

	// Create and run the UI
	filter := NewFilter(records)
	p := tea.NewProgram(NewUI(filter))
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

	// Get the selected record if any
	if model, ok := m.(Model); ok {
		if record, ok := model.Selected(); ok {
			fmt.Printf("Selected: %s %s\n", record.Command, record.Arguments)
		}
	}
}
