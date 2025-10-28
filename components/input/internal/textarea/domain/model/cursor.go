// Package model provides rich domain models for textarea.
package model

// Cursor tracks current editing position.
// This is a rich domain model that encapsulates cursor behavior.
// All operations return new instances (immutable).
type Cursor struct {
	row int // Line number (0-based)
	col int // Column number (0-based, rune offset)
}

// NewCursor creates cursor at position.
func NewCursor(row, col int) *Cursor {
	return &Cursor{row: row, col: col}
}

// Row returns current row.
func (c *Cursor) Row() int {
	return c.row
}

// Col returns current column.
func (c *Cursor) Col() int {
	return c.col
}

// Position returns (row, col) as tuple.
func (c *Cursor) Position() (row, col int) {
	return c.row, c.col
}

// MoveTo returns new cursor at position.
func (c *Cursor) MoveTo(row, col int) *Cursor {
	return &Cursor{row: row, col: col}
}

// MoveBy returns new cursor offset by delta.
func (c *Cursor) MoveBy(deltaRow, deltaCol int) *Cursor {
	return &Cursor{
		row: c.row + deltaRow,
		col: c.col + deltaCol,
	}
}

// Copy returns copy of cursor.
func (c *Cursor) Copy() *Cursor {
	return &Cursor{row: c.row, col: c.col}
}
