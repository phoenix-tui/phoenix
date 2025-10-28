// Package value provides value objects for textarea domain model.
package value

// Position represents a position in the text buffer (row, col).
// This is a value object - immutable and comparable by value.
type Position struct {
	row int // Line number (0-based)
	col int // Column number (0-based, rune offset)
}

// NewPosition creates a new position.
func NewPosition(row, col int) Position {
	return Position{row: row, col: col}
}

// Row returns the row (line number).
func (p Position) Row() int {
	return p.row
}

// Col returns the column (rune offset).
func (p Position) Col() int {
	return p.col
}

// IsBefore returns true if this position is before other position.
func (p Position) IsBefore(other Position) bool {
	if p.row < other.row {
		return true
	}
	if p.row == other.row && p.col < other.col {
		return true
	}
	return false
}

// IsAfter returns true if this position is after other position.
func (p Position) IsAfter(other Position) bool {
	return other.IsBefore(p)
}

// Equals returns true if positions are equal.
func (p Position) Equals(other Position) bool {
	return p.row == other.row && p.col == other.col
}

// String returns string representation for debugging.
func (p Position) String() string {
	return "(" + string(rune(p.row+'0')) + "," + string(rune(p.col+'0')) + ")"
}
