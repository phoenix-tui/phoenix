// Package value provides immutable value objects for the layout domain.
// These are pure data structures with no behavior beyond validation and equality.
package value

import "fmt"

// Size represents dimensional constraints for layout elements.
// It supports minimum, maximum, and exact size constraints.
//
// Design Philosophy:
//   - Immutable value object
//   - Validation on construction
//   - Negative values represent "not set" or "unlimited"
//   - Width/Height are exact constraints
//   - Min/Max are soft constraints (can be overridden)
//
// Example:
//
//	// Exact size: 80x24
//	size := NewSize(80, 24, -1, -1, -1, -1)
//
//	// Min width 40, max width 120, no height constraints
//	size := NewSize(-1, -1, 40, 120, -1, -1)
//
//	// Only min constraints
//	size := NewSize(-1, -1, 20, -1, 10, -1)
type Size struct {
	width     int // Exact width (-1 = not set)
	height    int // Exact height (-1 = not set)
	minWidth  int // Minimum width (-1 = unlimited)
	maxWidth  int // Maximum width (-1 = unlimited)
	minHeight int // Minimum height (-1 = unlimited)
	maxHeight int // Maximum height (-1 = unlimited)
}

// NewSize creates a Size with the given constraints.
// Negative values (-1) mean "not set" or "unlimited".
//
// Validation rules:
//   - If width is set, it must be >= 0
//   - If height is set, it must be >= 0
//   - If minWidth is set, it must be >= 0
//   - If maxWidth is set, it must be >= minWidth (if both set)
//   - Same rules apply for height
//   - If width is set, min/max are ignored
//   - If height is set, min/max are ignored
func NewSize(width, height, minWidth, maxWidth, minHeight, maxHeight int) Size {
	// Normalize -1 for "not set"
	if width < 0 {
		width = -1
	}
	if height < 0 {
		height = -1
	}
	if minWidth < 0 {
		minWidth = -1
	}
	if maxWidth < 0 {
		maxWidth = -1
	}
	if minHeight < 0 {
		minHeight = -1
	}
	if maxHeight < 0 {
		maxHeight = -1
	}

	// Validate min <= max (if both set)
	if minWidth != -1 && maxWidth != -1 && minWidth > maxWidth {
		panic(fmt.Sprintf("minWidth (%d) cannot be greater than maxWidth (%d)", minWidth, maxWidth))
	}
	if minHeight != -1 && maxHeight != -1 && minHeight > maxHeight {
		panic(fmt.Sprintf("minHeight (%d) cannot be greater than maxHeight (%d)", minHeight, maxHeight))
	}

	return Size{
		width:     width,
		height:    height,
		minWidth:  minWidth,
		maxWidth:  maxWidth,
		minHeight: minHeight,
		maxHeight: maxHeight,
	}
}

// NewSizeExact creates a Size with exact width and height.
// This is a convenience constructor for fixed-size elements.
//
// Example:
//
//	size := NewSizeExact(80, 24) // Exactly 80x24
func NewSizeExact(width, height int) Size {
	return NewSize(width, height, -1, -1, -1, -1)
}

// NewSizeUnconstrained creates a Size with no constraints.
// The element can be any size (determined by content or parent).
func NewSizeUnconstrained() Size {
	return NewSize(-1, -1, -1, -1, -1, -1)
}

// Width returns the exact width constraint (-1 if not set).
func (s Size) Width() int {
	return s.width
}

// Height returns the exact height constraint (-1 if not set).
func (s Size) Height() int {
	return s.height
}

// MinWidth returns the minimum width constraint (-1 if not set).
func (s Size) MinWidth() int {
	return s.minWidth
}

// MaxWidth returns the maximum width constraint (-1 if not set).
func (s Size) MaxWidth() int {
	return s.maxWidth
}

// MinHeight returns the minimum height constraint (-1 if not set).
func (s Size) MinHeight() int {
	return s.minHeight
}

// MaxHeight returns the maximum height constraint (-1 if not set).
func (s Size) MaxHeight() int {
	return s.maxHeight
}

// HasWidth returns true if exact width is set.
func (s Size) HasWidth() bool {
	return s.width != -1
}

// HasHeight returns true if exact height is set.
func (s Size) HasHeight() bool {
	return s.height != -1
}

// HasMinWidth returns true if minimum width is set.
func (s Size) HasMinWidth() bool {
	return s.minWidth != -1
}

// HasMaxWidth returns true if maximum width is set.
func (s Size) HasMaxWidth() bool {
	return s.maxWidth != -1
}

// HasMinHeight returns true if minimum height is set.
func (s Size) HasMinHeight() bool {
	return s.minHeight != -1
}

// HasMaxHeight returns true if maximum height is set.
func (s Size) HasMaxHeight() bool {
	return s.maxHeight != -1
}

// IsUnconstrained returns true if no constraints are set.
func (s Size) IsUnconstrained() bool {
	return !s.HasWidth() && !s.HasHeight() &&
		!s.HasMinWidth() && !s.HasMaxWidth() &&
		!s.HasMinHeight() && !s.HasMaxHeight()
}

// WithWidth returns a new Size with the given exact width.
func (s Size) WithWidth(width int) Size {
	return NewSize(width, s.height, s.minWidth, s.maxWidth, s.minHeight, s.maxHeight)
}

// WithHeight returns a new Size with the given exact height.
func (s Size) WithHeight(height int) Size {
	return NewSize(s.width, height, s.minWidth, s.maxWidth, s.minHeight, s.maxHeight)
}

// WithMinWidth returns a new Size with the given minimum width.
func (s Size) WithMinWidth(minWidth int) Size {
	return NewSize(s.width, s.height, minWidth, s.maxWidth, s.minHeight, s.maxHeight)
}

// WithMaxWidth returns a new Size with the given maximum width.
func (s Size) WithMaxWidth(maxWidth int) Size {
	return NewSize(s.width, s.height, s.minWidth, maxWidth, s.minHeight, s.maxHeight)
}

// WithMinHeight returns a new Size with the given minimum height.
func (s Size) WithMinHeight(minHeight int) Size {
	return NewSize(s.width, s.height, s.minWidth, s.maxWidth, minHeight, s.maxHeight)
}

// WithMaxHeight returns a new Size with the given maximum height.
func (s Size) WithMaxHeight(maxHeight int) Size {
	return NewSize(s.width, s.height, s.minWidth, s.maxWidth, s.minHeight, maxHeight)
}

// Constrain applies the constraints to the given dimensions.
// This is used by layout engine to clamp sizes within bounds.
//
// Rules:
//  1. If exact width/height is set, use it
//  2. Otherwise, clamp between min/max (if set)
//  3. If unconstrained, return input as-is
//
// Example:
//
//	size := NewSize(-1, -1, 40, 120, 10, 30)
//	w, h := size.Constrain(150, 5)  // Returns (120, 10) - clamped to max/min
func (s Size) Constrain(width, height int) (int, int) {
	// Apply width constraints
	if s.HasWidth() {
		width = s.width
	} else {
		if s.HasMinWidth() && width < s.minWidth {
			width = s.minWidth
		}
		if s.HasMaxWidth() && width > s.maxWidth {
			width = s.maxWidth
		}
	}

	// Apply height constraints
	if s.HasHeight() {
		height = s.height
	} else {
		if s.HasMinHeight() && height < s.minHeight {
			height = s.minHeight
		}
		if s.HasMaxHeight() && height > s.maxHeight {
			height = s.maxHeight
		}
	}

	return width, height
}

// String returns a human-readable representation of the Size.
func (s Size) String() string {
	if s.IsUnconstrained() {
		return "Size{unconstrained}"
	}

	result := "Size{"
	if s.HasWidth() {
		result += fmt.Sprintf("width=%d", s.width)
	} else {
		if s.HasMinWidth() {
			result += fmt.Sprintf("minW=%d", s.minWidth)
		}
		if s.HasMaxWidth() {
			if s.HasMinWidth() {
				result += ","
			}
			result += fmt.Sprintf("maxW=%d", s.maxWidth)
		}
	}

	if s.HasHeight() {
		if result != "Size{" {
			result += " "
		}
		result += fmt.Sprintf("height=%d", s.height)
	} else {
		if s.HasMinHeight() {
			if result != "Size{" {
				result += " "
			}
			result += fmt.Sprintf("minH=%d", s.minHeight)
		}
		if s.HasMaxHeight() {
			if result != "Size{" && !s.HasMinHeight() {
				result += " "
			} else if s.HasMinHeight() {
				result += ","
			}
			result += fmt.Sprintf("maxH=%d", s.maxHeight)
		}
	}

	result += "}"
	return result
}
