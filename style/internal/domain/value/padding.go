//nolint:dupl // Padding and Margin are separate domain concepts - duplication is intentional
package value

import "fmt"

// Padding represents an immutable CSS-like padding value object.
// Follows the CSS box model: top, right, bottom, left.
// This is a value object in DDD terms - immutable and defined by its values.
type Padding struct {
	top, right, bottom, left int
}

// NewPadding creates a new Padding with individual values for each side.
// CSS order: top, right, bottom, left (clockwise from top).
// Negative values are clamped to 0.
func NewPadding(top, right, bottom, left int) Padding {
	return Padding{
		top:    max(0, top),
		right:  max(0, right),
		bottom: max(0, bottom),
		left:   max(0, left),
	}
}

// UniformPadding creates a new Padding with the same value for all sides.
// Negative values are clamped to 0.
func UniformPadding(all int) Padding {
	all = max(0, all)
	return Padding{
		top:    all,
		right:  all,
		bottom: all,
		left:   all,
	}
}

// VerticalHorizontal creates a new Padding with vertical (top/bottom) and horizontal (left/right) values.
// Negative values are clamped to 0.
func VerticalHorizontal(vertical, horizontal int) Padding {
	vertical = max(0, vertical)
	horizontal = max(0, horizontal)
	return Padding{
		top:    vertical,
		right:  horizontal,
		bottom: vertical,
		left:   horizontal,
	}
}

// Top returns the top padding value.
func (p Padding) Top() int {
	return p.top
}

// Right returns the right padding value.
func (p Padding) Right() int {
	return p.right
}

// Bottom returns the bottom padding value.
func (p Padding) Bottom() int {
	return p.bottom
}

// Left returns the left padding value.
func (p Padding) Left() int {
	return p.left
}

// Horizontal returns the total horizontal padding (left + right).
func (p Padding) Horizontal() int {
	return p.left + p.right
}

// Vertical returns the total vertical padding (top + bottom).
func (p Padding) Vertical() int {
	return p.top + p.bottom
}

// Total returns the total vertical and horizontal padding.
func (p Padding) Total() (vertical, horizontal int) {
	return p.Vertical(), p.Horizontal()
}

// Equal returns true if this padding equals another padding.
// Value objects are compared by value, not identity.
func (p Padding) Equal(other Padding) bool {
	return p.top == other.top &&
		p.right == other.right &&
		p.bottom == other.bottom &&
		p.left == other.left
}

// String returns a human-readable representation of the padding.
func (p Padding) String() string {
	return fmt.Sprintf("Padding(top=%d, right=%d, bottom=%d, left=%d)", p.top, p.right, p.bottom, p.left)
}

// --- Private helpers ---.
// Note: Using Go 1.21+ builtin max() function.
