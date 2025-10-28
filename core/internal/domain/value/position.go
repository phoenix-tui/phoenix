// Package value provides immutable value objects for terminal operations.
package value

// Position represents a position in the terminal (row, column).
// This is an immutable value object using 0-based indexing.
//
// Invariants:
//   - Row and Col are always >= 0
//   - Position is immutable after creation
type Position struct {
	Row int // 0-based row index
	Col int // 0-based column index
}

// NewPosition creates a position with validation.
// Negative values are clamped to 0.
func NewPosition(row, col int) Position {
	if row < 0 {
		row = 0
	}
	if col < 0 {
		col = 0
	}
	return Position{Row: row, Col: col}
}

// Add returns a new Position offset by the given deltas.
func (p Position) Add(deltaRow, deltaCol int) Position {
	return NewPosition(p.Row+deltaRow, p.Col+deltaCol)
}

// IsZero returns true if position is at origin (0, 0).
func (p Position) IsZero() bool {
	return p.Row == 0 && p.Col == 0
}

// Equal returns true if positions are equal.
func (p Position) Equal(other Position) bool {
	return p.Row == other.Row && p.Col == other.Col
}
