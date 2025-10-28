//nolint:dupl // Margin and Padding are separate domain concepts - duplication is intentional
package value

import "fmt"

// Margin represents an immutable CSS-like margin value object.
// Follows the CSS box model: top, right, bottom, left.
// This is a value object in DDD terms - immutable and defined by its values.
//
// Note: Margin is kept as a separate type from Padding for type safety and.
// future extensibility (e.g., margin might support "auto" in the future).
type Margin struct {
	top, right, bottom, left int
}

// NewMargin creates a new Margin with individual values for each side.
// CSS order: top, right, bottom, left (clockwise from top).
// Negative values are clamped to 0.
func NewMargin(top, right, bottom, left int) Margin {
	return Margin{
		top:    max(0, top),
		right:  max(0, right),
		bottom: max(0, bottom),
		left:   max(0, left),
	}
}

// UniformMargin creates a new Margin with the same value for all sides.
// Negative values are clamped to 0.
func UniformMargin(all int) Margin {
	all = max(0, all)
	return Margin{
		top:    all,
		right:  all,
		bottom: all,
		left:   all,
	}
}

// VerticalHorizontalMargin creates a new Margin with vertical (top/bottom) and horizontal (left/right) values.
// Negative values are clamped to 0.
func VerticalHorizontalMargin(vertical, horizontal int) Margin {
	vertical = max(0, vertical)
	horizontal = max(0, horizontal)
	return Margin{
		top:    vertical,
		right:  horizontal,
		bottom: vertical,
		left:   horizontal,
	}
}

// Top returns the top margin value.
func (m Margin) Top() int {
	return m.top
}

// Right returns the right margin value.
func (m Margin) Right() int {
	return m.right
}

// Bottom returns the bottom margin value.
func (m Margin) Bottom() int {
	return m.bottom
}

// Left returns the left margin value.
func (m Margin) Left() int {
	return m.left
}

// Horizontal returns the total horizontal margin (left + right).
func (m Margin) Horizontal() int {
	return m.left + m.right
}

// Vertical returns the total vertical margin (top + bottom).
func (m Margin) Vertical() int {
	return m.top + m.bottom
}

// Total returns the total vertical and horizontal margin.
func (m Margin) Total() (vertical, horizontal int) {
	return m.Vertical(), m.Horizontal()
}

// Equal returns true if this margin equals another margin.
// Value objects are compared by value, not identity.
func (m Margin) Equal(other Margin) bool {
	return m.top == other.top &&
		m.right == other.right &&
		m.bottom == other.bottom &&
		m.left == other.left
}

// String returns a human-readable representation of the margin.
func (m Margin) String() string {
	return fmt.Sprintf("Margin(top=%d, right=%d, bottom=%d, left=%d)", m.top, m.right, m.bottom, m.left)
}
