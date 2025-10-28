// Package model provides rich domain models for textarea.
package model

// CursorPos represents a cursor position in the text buffer.
// This is a value object used for cursor movement validation and observation.
type CursorPos struct {
	Row int // Line number (0-based)
	Col int // Column number (0-based, rune offset)
}

// NewCursorPos creates a new cursor position.
func NewCursorPos(row, col int) CursorPos {
	return CursorPos{Row: row, Col: col}
}

// Equals returns true if two positions are equal.
func (p CursorPos) Equals(other CursorPos) bool {
	return p.Row == other.Row && p.Col == other.Col
}
