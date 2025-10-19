package value

import "fmt"

// Position represents a mouse position in terminal coordinates.
// Coordinates are 0-based (top-left is 0,0).
type Position struct {
	x int
	y int
}

// NewPosition creates a new Position.
func NewPosition(x, y int) Position {
	return Position{x: x, y: y}
}

// X returns the X coordinate.
func (p Position) X() int {
	return p.x
}

// Y returns the Y coordinate.
func (p Position) Y() int {
	return p.y
}

// String returns the string representation of the position.
func (p Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

// Equals checks if two positions are equal.
func (p Position) Equals(other Position) bool {
	return p.x == other.x && p.y == other.y
}

// DistanceTo calculates the Manhattan distance to another position.
func (p Position) DistanceTo(other Position) int {
	dx := p.x - other.x
	dy := p.y - other.y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

// IsWithinTolerance checks if two positions are within a tolerance distance.
func (p Position) IsWithinTolerance(other Position, tolerance int) bool {
	return p.DistanceTo(other) <= tolerance
}
