package value

// Position represents a 2D coordinate in the terminal buffer.
// Position is immutable and uses value semantics.
type Position struct {
	x int
	y int
}

// NewPosition creates a new Position at (x, y).
func NewPosition(x, y int) Position {
	return Position{x: x, y: y}
}

// X returns the x coordinate.
func (p Position) X() int {
	return p.x
}

// Y returns the y coordinate.
func (p Position) Y() int {
	return p.y
}

// Add returns a new Position offset by (dx, dy).
func (p Position) Add(dx, dy int) Position {
	return Position{x: p.x + dx, y: p.y + dy}
}

// Equals checks if two positions are equal.
func (p Position) Equals(other Position) bool {
	return p.x == other.x && p.y == other.y
}

// IsZero returns true if position is at origin (0, 0).
func (p Position) IsZero() bool {
	return p.x == 0 && p.y == 0
}

// WithX returns a new Position with x changed.
func (p Position) WithX(x int) Position {
	return Position{x: x, y: p.y}
}

// WithY returns a new Position with y changed.
func (p Position) WithY(y int) Position {
	return Position{x: p.x, y: y}
}
