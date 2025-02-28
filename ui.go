// Package ui provides the terminal user interface for retour.
// It uses the Bubble Tea framework to create an interactive TUI that
// displays and filters command history.
package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Style definitions
var (
	// Style for the filter input at the bottom
	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	// Style for selected items
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	// Style for normal items
	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))
)

// Model represents the UI state and data
type Model struct {
	filter     *Filter // Filter for records
	cursor     int     // Current selection in the list
	textCursor int     // Current cursor position in filter input
	selected   bool    // Whether a selection has been made
	height     int     // Terminal height
}

// Records returns all records (for testing)
func (m Model) Records() []Record {
	return m.filter.FilteredRecords()
}

// Cursor returns the current cursor position (for testing)
func (m Model) Cursor() int {
	return m.cursor
}

// New creates a new UI model with the given filter
func NewUI(filter *Filter) Model {
	return Model{
		filter:     filter,
		cursor:     0,
		textCursor: 0,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles input and updates the model
// findWordStart finds the start of the word before the given position
func findWordStart(text string, pos int) int {
	// Skip spaces immediately before pos
	for pos > 0 && pos-1 < len(text) && text[pos-1] == ' ' {
		pos--
	}
	// Find start of word
	for pos > 0 && pos-1 < len(text) && text[pos-1] != ' ' {
		pos--
	}
	return pos
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyUp, tea.KeyCtrlP:
			if m.cursor > 0 {
				m.cursor--
			}

		case tea.KeyDown, tea.KeyCtrlN:
			if m.cursor < len(m.filter.FilteredRecords())-1 {
				m.cursor++
			}

		case tea.KeyEnter:
			m.selected = true
			return m, tea.Quit

		case tea.KeyBackspace:
			if len(m.filter.Filter()) > 0 && m.textCursor > 0 {
				// Remove the character before the cursor
				m.filter.RemoveCharBeforeCursor(m.textCursor)
				m.textCursor--
			}

		case tea.KeyLeft:
			if m.textCursor > 0 {
				m.textCursor--
			}

		case tea.KeyRight:
			if m.textCursor < len(m.filter.Filter()) {
				m.textCursor++
			}

		case tea.KeyCtrlA:
			// Beginning of line
			m.textCursor = 0

		case tea.KeyCtrlE:
			// End of line
			m.textCursor = len(m.filter.Filter())

		case tea.KeyCtrlW:
			// Kill word backward
			if m.textCursor > 0 {
				newPos := findWordStart(m.filter.Filter(), m.textCursor)
				m.filter.RemoveTextBeforeCursor(newPos, m.textCursor)
				m.textCursor = newPos
			}

		case tea.KeyCtrlK:
			// Kill to end of line
			if m.textCursor < len(m.filter.Filter()) {
				m.filter.RemoveTextAfterCursor(m.textCursor)
			}

		case tea.KeySpace:
			// Insert space at cursor position
			m.filter.InsertCharAtCursor(' ', m.textCursor)
			m.textCursor++

		case tea.KeyRunes:
			// Insert the characters at the cursor position
			m.filter.InsertTextAtCursor(string(msg.Runes), m.textCursor)
			m.textCursor += len(msg.Runes)
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.height == 0 {
		return "Loading..."
	}

	// Reserve space for input line and padding
	maxItems := m.height - 2
	if maxItems <= 0 {
		return "Window too small"
	}

	// Build the list view
	var s strings.Builder

	// Calculate which items to show
	start := 0
	if len(m.filter.FilteredRecords()) > maxItems && m.cursor >= maxItems {
		start = m.cursor - maxItems + 1
	}
	end := min(start+maxItems, len(m.filter.FilteredRecords()))

	// Render visible items
	for i, record := range m.filter.FilteredRecords()[start:end] {
		// Format the record
		line := formatRecord(record)

		// Style based on selection
		if i+start == m.cursor {
			s.WriteString(selectedStyle.Render("> " + line))
		} else {
			s.WriteString(normalStyle.Render("  " + line))
		}
		s.WriteRune('\n')
	}

	// Add the filter input at the bottom with cursor
	prefix := "Filter: "
	beforeCursor := m.filter.Filter()[:m.textCursor]
	afterCursor := m.filter.Filter()[m.textCursor:]
	cursorChar := "█"
	if len(afterCursor) > 0 {
		cursorChar = string(afterCursor[0])
		afterCursor = afterCursor[1:]
	}
	s.WriteString(inputStyle.Render(prefix + beforeCursor))
	s.WriteString(inputStyle.Reverse(true).Render(cursorChar))
	s.WriteString(inputStyle.Render(afterCursor))

	return s.String()
}

// Selected returns the currently selected record, if any
func (m Model) Selected() (Record, bool) {
	if !m.selected || len(m.filter.FilteredRecords()) == 0 {
		return Record{}, false
	}
	return m.filter.FilteredRecords()[m.cursor], true
}

// formatRecord formats a record for display
func formatRecord(r Record) string {
	status := "✓"
	if r.ExitStatus != 0 {
		status = "✗"
	}
	return status + " " + r.Command + " " + r.Arguments
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
