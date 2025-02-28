package main

import (
	"strings"
)

// Filter represents a fuzzy matcher for Record objects
type Filter struct {
	records         []Record // All available records
	filteredRecords []Record // Records after filtering
	filter          string   // Current filter text
}

// NewFilter creates a new Filter with the given records
func NewFilter(records []Record) *Filter {
	return &Filter{
		records:         records,
		filteredRecords: records, // Initially show all records
		filter:          "",      // Initially empty filter
	}
}

// FilteredRecords returns the current set of filtered records
func (f *Filter) FilteredRecords() []Record {
	return f.filteredRecords
}

// Filter returns the current filter text
func (f *Filter) Filter() string {
	return f.filter
}

// UpdateFilter updates the filter text and refreshes the filtered records
func (f *Filter) UpdateFilter(filterText string) {
	f.filter = filterText

	// If filter is empty, show all records
	if filterText == "" {
		f.filteredRecords = f.records
		return
	}

	// Naive implementation: check if record contains the filter string
	// in either the command or arguments (case insensitive)
	var filtered []Record
	lowerFilter := strings.ToLower(filterText)

	for _, record := range f.records {
		lowerCommand := strings.ToLower(record.Command)
		lowerArgs := strings.ToLower(record.Arguments)

		// Check if command or arguments contain the filter string
		if strings.Contains(lowerCommand, lowerFilter) ||
			strings.Contains(lowerArgs, lowerFilter) {
			filtered = append(filtered, record)
		}
	}

	f.filteredRecords = filtered
}

// InsertTextAtCursor inserts text at the specified cursor position
func (f *Filter) InsertTextAtCursor(text string, cursorPos int) {
	if len(text) == 0 {
		return
	}

	// Ensure cursor position is valid
	if cursorPos < 0 {
		cursorPos = 0
	}
	if cursorPos > len(f.filter) {
		cursorPos = len(f.filter)
	}

	// Insert text at cursor position
	newFilter := f.filter[:cursorPos] + text + f.filter[cursorPos:]
	f.UpdateFilter(newFilter)
}

// InsertCharAtCursor inserts a single character at the specified cursor position
func (f *Filter) InsertCharAtCursor(char rune, cursorPos int) {
	f.InsertTextAtCursor(string(char), cursorPos)
}

// RemoveCharBeforeCursor removes the character before the specified cursor position
func (f *Filter) RemoveCharBeforeCursor(cursorPos int) {
	if cursorPos > 0 && cursorPos <= len(f.filter) {
		newFilter := f.filter[:cursorPos-1] + f.filter[cursorPos:]
		f.UpdateFilter(newFilter)
	}
}

// RemoveTextBeforeCursor removes text from newPos to the specified cursor position
func (f *Filter) RemoveTextBeforeCursor(newPos int, cursorPos int) {
	if newPos < 0 {
		newPos = 0
	}
	if cursorPos > len(f.filter) {
		cursorPos = len(f.filter)
	}

	if newPos < cursorPos {
		newFilter := f.filter[:newPos] + f.filter[cursorPos:]
		f.UpdateFilter(newFilter)
	}
}

// RemoveTextAfterCursor removes all text after the specified cursor position
func (f *Filter) RemoveTextAfterCursor(cursorPos int) {
	if cursorPos >= 0 && cursorPos <= len(f.filter) {
		newFilter := f.filter[:cursorPos]
		f.UpdateFilter(newFilter)
	}
}
