package value

import "fmt"

// Spacing represents padding or margin around an element.
// Follows CSS box model conventions (top, right, bottom, left).
//
// Design Philosophy:
//   - Immutable value object
//   - CSS-like API (top, right, bottom, left)
//   - Non-negative values only
//   - Convenience constructors for common patterns
//
// Box Model:
//
//	        top
//	    ┌─────────┐
//	    │  ┌───┐  │
//	left│  │   │  │right
//	    │  └───┘  │
//	    └─────────┘
//	       bottom
//
// Example:
//
//	// All sides equal
//	spacing := NewSpacingAll(2)  // 2 on all sides
//
//	// Vertical and horizontal
//	spacing := NewSpacingVH(1, 2)  // 1 top/bottom, 2 left/right
//
//	// Individual sides
//	spacing := NewSpacing(1, 2, 1, 2)  // top, right, bottom, left
type Spacing struct {
	top    int
	right  int
	bottom int
	left   int
}

// NewSpacing creates a Spacing with individual values for each side.
// Negative values are clamped to 0.
//
// Example:
//
//	spacing := NewSpacing(1, 2, 1, 2)  // CSS: padding: 1 2 1 2
func NewSpacing(top, right, bottom, left int) Spacing {
	return Spacing{
		top:    max(0, top),
		right:  max(0, right),
		bottom: max(0, bottom),
		left:   max(0, left),
	}
}

// NewSpacingAll creates a Spacing with the same value on all sides.
// Negative values are clamped to 0.
//
// Example:
//
//	spacing := NewSpacingAll(2)  // CSS: padding: 2
func NewSpacingAll(value int) Spacing {
	v := max(0, value)
	return Spacing{
		top:    v,
		right:  v,
		bottom: v,
		left:   v,
	}
}

// NewSpacingVH creates a Spacing with vertical and horizontal values.
// Negative values are clamped to 0.
//
// Example:
//
//	spacing := NewSpacingVH(1, 2)  // CSS: padding: 1 2 (vertical, horizontal)
func NewSpacingVH(vertical, horizontal int) Spacing {
	v := max(0, vertical)
	h := max(0, horizontal)
	return Spacing{
		top:    v,
		right:  h,
		bottom: v,
		left:   h,
	}
}

// NewSpacingZero creates a Spacing with zero on all sides.
// This is useful as a default/empty value.
func NewSpacingZero() Spacing {
	return Spacing{
		top:    0,
		right:  0,
		bottom: 0,
		left:   0,
	}
}

// Top returns the top spacing.
func (s Spacing) Top() int {
	return s.top
}

// Right returns the right spacing.
func (s Spacing) Right() int {
	return s.right
}

// Bottom returns the bottom spacing.
func (s Spacing) Bottom() int {
	return s.bottom
}

// Left returns the left spacing.
func (s Spacing) Left() int {
	return s.left
}

// Horizontal returns the total horizontal spacing (left + right).
func (s Spacing) Horizontal() int {
	return s.left + s.right
}

// Vertical returns the total vertical spacing (top + bottom).
func (s Spacing) Vertical() int {
	return s.top + s.bottom
}

// IsZero returns true if all sides are zero.
func (s Spacing) IsZero() bool {
	return s.top == 0 && s.right == 0 && s.bottom == 0 && s.left == 0
}

// IsUniform returns true if all sides have the same value.
func (s Spacing) IsUniform() bool {
	return s.top == s.right && s.right == s.bottom && s.bottom == s.left
}

// WithTop returns a new Spacing with the given top value.
func (s Spacing) WithTop(top int) Spacing {
	return NewSpacing(top, s.right, s.bottom, s.left)
}

// WithRight returns a new Spacing with the given right value.
func (s Spacing) WithRight(right int) Spacing {
	return NewSpacing(s.top, right, s.bottom, s.left)
}

// WithBottom returns a new Spacing with the given bottom value.
func (s Spacing) WithBottom(bottom int) Spacing {
	return NewSpacing(s.top, s.right, bottom, s.left)
}

// WithLeft returns a new Spacing with the given left value.
func (s Spacing) WithLeft(left int) Spacing {
	return NewSpacing(s.top, s.right, s.bottom, left)
}

// Add returns a new Spacing with values added.
// This is useful for combining padding and margin.
//
// Example:
//
//	padding := NewSpacingAll(1)
//	margin := NewSpacingAll(2)
//	total := padding.Add(margin)  // All sides = 3
func (s Spacing) Add(other Spacing) Spacing {
	return NewSpacing(
		s.top+other.top,
		s.right+other.right,
		s.bottom+other.bottom,
		s.left+other.left,
	)
}

// Scale returns a new Spacing with values multiplied by a factor.
// Negative factors are treated as 0.
//
// Example:
//
//	spacing := NewSpacingAll(2)
//	doubled := spacing.Scale(2)  // All sides = 4
func (s Spacing) Scale(factor int) Spacing {
	if factor < 0 {
		factor = 0
	}
	return NewSpacing(
		s.top*factor,
		s.right*factor,
		s.bottom*factor,
		s.left*factor,
	)
}

// Equals returns true if both spacings are equal.
func (s Spacing) Equals(other Spacing) bool {
	return s.top == other.top &&
		s.right == other.right &&
		s.bottom == other.bottom &&
		s.left == other.left
}

// String returns a human-readable representation.
// Uses CSS shorthand notation when possible.
func (s Spacing) String() string {
	if s.IsZero() {
		return "Spacing{0}"
	}
	if s.IsUniform() {
		return fmt.Sprintf("Spacing{%d}", s.top)
	}
	if s.top == s.bottom && s.left == s.right {
		return fmt.Sprintf("Spacing{%d %d}", s.top, s.left)
	}
	return fmt.Sprintf("Spacing{%d %d %d %d}", s.top, s.right, s.bottom, s.left)
}

// Helper function (Go 1.25+ has max, but for clarity)
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
