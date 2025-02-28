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
	records     []Record // All available records
	filtered    []Record // Records after filtering
	cursor      int      // Current selection in the list
	filterInput string   // Current filter text
	textCursor  int      // Current cursor position in filter input
	selected    bool     // Whether a selection has been made
	height      int      // Terminal height
}

// Records returns all records (for testing)
func (m Model) Records() []Record {
	return m.filtered
}

// Cursor returns the current cursor position (for testing)
func (m Model) Cursor() int {
	return m.cursor
}

// New creates a new UI model with the given records
func NewUI(records []Record) Model {
	return Model{
		records:    records,
		filtered:   records, // Initially show all records
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
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}

		case tea.KeyEnter:
			m.selected = true
			return m, tea.Quit

		case tea.KeyBackspace:
			if len(m.filterInput) > 0 && m.textCursor > 0 {
				// Remove the character before the cursor
				m.filterInput = m.filterInput[:m.textCursor-1] + m.filterInput[m.textCursor:]
				m.textCursor--
				m.updateFilter()
			}

		case tea.KeyLeft:
			if m.textCursor > 0 {
				m.textCursor--
			}

		case tea.KeyRight:
			if m.textCursor < len(m.filterInput) {
				m.textCursor++
			}

		case tea.KeyCtrlA:
			// Beginning of line
			m.textCursor = 0

		case tea.KeyCtrlE:
			// End of line
			m.textCursor = len(m.filterInput)

		case tea.KeyCtrlW:
			// Kill word backward
			if m.textCursor > 0 {
				newPos := findWordStart(m.filterInput, m.textCursor)
				m.filterInput = m.filterInput[:newPos] + m.filterInput[m.textCursor:]
				m.textCursor = newPos
				m.updateFilter()
			}

		case tea.KeyCtrlK:
			// Kill to end of line
			if m.textCursor < len(m.filterInput) {
				m.filterInput = m.filterInput[:m.textCursor]
				m.updateFilter()
			}

		case tea.KeySpace:
			// Insert space at cursor position
			m.filterInput = m.filterInput[:m.textCursor] + " " + m.filterInput[m.textCursor:]
			m.textCursor++
			m.updateFilter()

		case tea.KeyRunes:
			// Insert the characters at the cursor position
			text := string(msg.Runes)
			m.filterInput = m.filterInput[:m.textCursor] + text + m.filterInput[m.textCursor:]
			m.textCursor += len(text)
			m.updateFilter()
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
	if len(m.filtered) > maxItems && m.cursor >= maxItems {
		start = m.cursor - maxItems + 1
	}
	end := min(start+maxItems, len(m.filtered))

	// Render visible items
	for i, record := range m.filtered[start:end] {
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
	beforeCursor := m.filterInput[:m.textCursor]
	afterCursor := m.filterInput[m.textCursor:]
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
	if !m.selected || len(m.filtered) == 0 {
		return Record{}, false
	}
	return m.filtered[m.cursor], true
}

// formatRecord formats a record for display
func formatRecord(r Record) string {
	status := "✓"
	if r.ExitStatus != 0 {
		status = "✗"
	}
	return status + " " + r.Command + " " + r.Arguments
}

// updateFilter applies the current filter to the records
// Currently a no-op as requested
func (m *Model) updateFilter() {
	// No-op for now
	m.filtered = m.records
	m.cursor = 0
}
