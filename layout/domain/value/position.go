package value

import "fmt"

// Position represents a 2D position in terminal space.
// Coordinates are 0-based and represent character cells.
//
// Design Philosophy:
//   - Immutable value object
//   - X is horizontal (column), Y is vertical (row)
//   - 0-based indexing (matches terminal conventions)
//   - Non-negative values only (validated on construction)
//
// Coordinate System:
//
//	(0,0) ─────────> X (columns)
//	  │
//	  │
//	  v
//	  Y (rows)
//
// Example:
//
//	pos := NewPosition(10, 5)  // Column 10, Row 5
//	pos = pos.Add(2, 3)        // Column 12, Row 8
type Position struct {
	x int // Horizontal position (column, 0-based)
	y int // Vertical position (row, 0-based)
}

// NewPosition creates a Position with the given coordinates.
// Negative values are clamped to 0.
//
// Example:
//
//	pos := NewPosition(10, 5)   // Valid position
//	pos := NewPosition(-1, -5)  // Clamped to (0, 0)
func NewPosition(x, y int) Position {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	return Position{x: x, y: y}
}

// Origin returns the origin position (0, 0).
func Origin() Position {
	return Position{x: 0, y: 0}
}

// X returns the horizontal position (column).
func (p Position) X() int {
	return p.x
}

// Y returns the vertical position (row).
func (p Position) Y() int {
	return p.y
}

// Add returns a new Position offset by the given delta.
// Negative results are clamped to 0.
//
// Example:
//
//	pos := NewPosition(10, 5)
//	pos = pos.Add(2, 3)  // (12, 8)
//	pos = pos.Add(-20, -10)  // (0, 0) - clamped
func (p Position) Add(dx, dy int) Position {
	return NewPosition(p.x+dx, p.y+dy)
}

// Sub returns a new Position offset by subtracting the delta.
// This is equivalent to Add(-dx, -dy).
// Negative results are clamped to 0.
//
// Example:
//
//	pos := NewPosition(10, 5)
//	pos = pos.Sub(2, 3)  // (8, 2)
func (p Position) Sub(dx, dy int) Position {
	return NewPosition(p.x-dx, p.y-dy)
}

// Offset returns a new Position offset by another Position.
// This is useful for relative positioning.
//
// Example:
//
//	base := NewPosition(10, 5)
//	offset := NewPosition(2, 3)
//	result := base.Offset(offset)  // (12, 8)
func (p Position) Offset(other Position) Position {
	return p.Add(other.x, other.y)
}

// Distance calculates the Manhattan distance to another Position.
// Manhattan distance = |x1 - x2| + |y1 - y2|
//
// Example:
//
//	p1 := NewPosition(10, 5)
//	p2 := NewPosition(15, 8)
//	dist := p1.Distance(p2)  // 5 + 3 = 8
func (p Position) Distance(other Position) int {
	dx := p.x - other.x
	if dx < 0 {
		dx = -dx
	}
	dy := p.y - other.y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

// Equals returns true if both positions are equal.
func (p Position) Equals(other Position) bool {
	return p.x == other.x && p.y == other.y
}

// IsOrigin returns true if position is (0, 0).
func (p Position) IsOrigin() bool {
	return p.x == 0 && p.y == 0
}

// String returns a human-readable representation.
func (p Position) String() string {
	return fmt.Sprintf("Position{x=%d, y=%d}", p.x, p.y)
}
