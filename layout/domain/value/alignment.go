package value

import "fmt"

// HorizontalAlignment represents horizontal text/content alignment.
type HorizontalAlignment int

const (
	// AlignLeft aligns content to the left edge.
	AlignLeft HorizontalAlignment = iota
	// AlignCenter centers content horizontally.
	AlignCenter
	// AlignRight aligns content to the right edge.
	AlignRight
)

// String returns a human-readable representation of HorizontalAlignment.
func (h HorizontalAlignment) String() string {
	switch h {
	case AlignLeft:
		return "Left"
	case AlignCenter:
		return "Center"
	case AlignRight:
		return "Right"
	default:
		return "Unknown"
	}
}

// VerticalAlignment represents vertical text/content alignment.
type VerticalAlignment int

const (
	// AlignTop aligns content to the top edge.
	AlignTop VerticalAlignment = iota
	// AlignMiddle centers content vertically.
	AlignMiddle
	// AlignBottom aligns content to the bottom edge.
	AlignBottom
)

// String returns a human-readable representation of VerticalAlignment.
func (v VerticalAlignment) String() string {
	switch v {
	case AlignTop:
		return "Top"
	case AlignMiddle:
		return "Middle"
	case AlignBottom:
		return "Bottom"
	default:
		return "Unknown"
	}
}

// Alignment combines horizontal and vertical alignment.
//
// Design Philosophy:
//   - Immutable value object
//   - Separate horizontal and vertical components
//   - CSS-like semantics (justify-content, align-items)
//   - Default is top-left (most common for TUI)
//
// Example:
//
//	// Center both axes
//	align := NewAlignment(AlignCenter, AlignMiddle)
//
//	// Top-left (default)
//	align := NewAlignment(AlignLeft, AlignTop)
//
//	// Bottom-right
//	align := NewAlignment(AlignRight, AlignBottom)
type Alignment struct {
	horizontal HorizontalAlignment
	vertical   VerticalAlignment
}

// NewAlignment creates an Alignment with the given horizontal and vertical components.
func NewAlignment(horizontal HorizontalAlignment, vertical VerticalAlignment) Alignment {
	return Alignment{
		horizontal: horizontal,
		vertical:   vertical,
	}
}

// NewAlignmentDefault creates an Alignment with default values (top-left).
// This is the most common alignment for TUI applications.
func NewAlignmentDefault() Alignment {
	return Alignment{
		horizontal: AlignLeft,
		vertical:   AlignTop,
	}
}

// NewAlignmentCenter creates an Alignment centered on both axes.
func NewAlignmentCenter() Alignment {
	return Alignment{
		horizontal: AlignCenter,
		vertical:   AlignMiddle,
	}
}

// Horizontal returns the horizontal alignment component.
func (a Alignment) Horizontal() HorizontalAlignment {
	return a.horizontal
}

// Vertical returns the vertical alignment component.
func (a Alignment) Vertical() VerticalAlignment {
	return a.vertical
}

// IsHorizontal checks if horizontal alignment matches the given value.
func (a Alignment) IsHorizontal(h HorizontalAlignment) bool {
	return a.horizontal == h
}

// IsVertical checks if vertical alignment matches the given value.
func (a Alignment) IsVertical(v VerticalAlignment) bool {
	return a.vertical == v
}

// IsCenter returns true if both axes are centered.
func (a Alignment) IsCenter() bool {
	return a.horizontal == AlignCenter && a.vertical == AlignMiddle
}

// IsDefault returns true if alignment is top-left (default).
func (a Alignment) IsDefault() bool {
	return a.horizontal == AlignLeft && a.vertical == AlignTop
}

// WithHorizontal returns a new Alignment with the given horizontal component.
func (a Alignment) WithHorizontal(h HorizontalAlignment) Alignment {
	return Alignment{
		horizontal: h,
		vertical:   a.vertical,
	}
}

// WithVertical returns a new Alignment with the given vertical component.
func (a Alignment) WithVertical(v VerticalAlignment) Alignment {
	return Alignment{
		horizontal: a.horizontal,
		vertical:   v,
	}
}

// Equals returns true if both alignments are equal.
func (a Alignment) Equals(other Alignment) bool {
	return a.horizontal == other.horizontal && a.vertical == other.vertical
}

// String returns a human-readable representation.
func (a Alignment) String() string {
	return fmt.Sprintf("Alignment{%s, %s}", a.horizontal, a.vertical)
}

// CalculateHorizontalOffset calculates the offset needed to align content within a container.
// This is a utility function used by layout engine.
//
// Parameters:
//   - contentSize: Size of the content to align (width or height)
//   - containerSize: Size of the container (width or height)
//
// Returns:
//   - offset: Number of cells to offset the content
//
// Example (horizontal alignment):
//
//	offset := CalculateHorizontalOffset(AlignCenter, 20, 80)
//	// Content width 20, container width 80
//	// Center alignment: offset = (80 - 20) / 2 = 30
func CalculateHorizontalOffset(h HorizontalAlignment, contentWidth, containerWidth int) int {
	if contentWidth >= containerWidth {
		return 0 // Content fills or exceeds container
	}

	switch h {
	case AlignLeft:
		return 0
	case AlignCenter:
		return (containerWidth - contentWidth) / 2
	case AlignRight:
		return containerWidth - contentWidth
	default:
		return 0
	}
}

// CalculateVerticalOffset calculates the vertical offset for alignment.
// See CalculateHorizontalOffset for details.
func CalculateVerticalOffset(v VerticalAlignment, contentHeight, containerHeight int) int {
	if contentHeight >= containerHeight {
		return 0 // Content fills or exceeds container
	}

	switch v {
	case AlignTop:
		return 0
	case AlignMiddle:
		return (containerHeight - contentHeight) / 2
	case AlignBottom:
		return containerHeight - contentHeight
	default:
		return 0
	}
}

// CalculateOffsets calculates both horizontal and vertical offsets.
// This is a convenience function combining both calculations.
//
// Returns:
//   - xOffset: Horizontal offset
//   - yOffset: Vertical offset
//
// Example:
//
//	align := NewAlignmentCenter()
//	xOffset, yOffset := align.CalculateOffsets(20, 10, 80, 24)
//	// Content 20x10, container 80x24
//	// Center alignment: xOffset=30, yOffset=7
func (a Alignment) CalculateOffsets(contentWidth, contentHeight, containerWidth, containerHeight int) (xOffset, yOffset int) {
	xOffset = CalculateHorizontalOffset(a.horizontal, contentWidth, containerWidth)
	yOffset = CalculateVerticalOffset(a.vertical, contentHeight, containerHeight)
	return
}
