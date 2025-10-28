// Package value contains value objects for the modal component.
package value

// Position represents modal position on screen.
// Positions can be either centered (automatic) or custom (specific x, y coordinates).
type Position struct {
	x      int  // X coordinate (or -1 for center)
	y      int  // Y coordinate (or -1 for center)
	center bool // Auto-center?
}

// NewPositionCenter creates a position that auto-centers the modal.
func NewPositionCenter() *Position {
	return &Position{
		x:      -1,
		y:      -1,
		center: true,
	}
}

// NewPositionCustom creates a position at specific x, y coordinates.
func NewPositionCustom(x, y int) *Position {
	return &Position{
		x:      x,
		y:      y,
		center: false,
	}
}

// X returns the X coordinate (-1 if centered).
func (p *Position) X() int {
	return p.x
}

// Y returns the Y coordinate (-1 if centered).
func (p *Position) Y() int {
	return p.y
}

// IsCenter returns true if the position is auto-centered.
func (p *Position) IsCenter() bool {
	return p.center
}
