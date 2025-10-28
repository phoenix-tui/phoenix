package value

import "fmt"

// Size represents immutable size constraints for width and height.
// Uses pointers to int to distinguish between "not set" (nil) and "set to 0" (non-nil).
// This is a value object in DDD terms - immutable and defined by its values.
type Size struct {
	width     *int // nil = not set
	height    *int // nil = not set
	minWidth  *int
	maxWidth  *int
	minHeight *int
	maxHeight *int
}

// NewSize creates a new Size with no constraints set.
func NewSize() Size {
	return Size{}
}

// WithWidth creates a new Size with a fixed width.
func WithWidth(w int) Size {
	w = max(0, w)
	return Size{width: &w}
}

// WithHeight creates a new Size with a fixed height.
func WithHeight(h int) Size {
	h = max(0, h)
	return Size{height: &h}
}

// Width returns the width constraint if set.
// The second return value indicates whether the width was set.
func (s Size) Width() (int, bool) {
	if s.width == nil {
		return 0, false
	}
	return *s.width, true
}

// Height returns the height constraint if set.
// The second return value indicates whether the height was set.
func (s Size) Height() (int, bool) {
	if s.height == nil {
		return 0, false
	}
	return *s.height, true
}

// MinWidth returns the minimum width constraint if set.
func (s Size) MinWidth() (int, bool) {
	if s.minWidth == nil {
		return 0, false
	}
	return *s.minWidth, true
}

// MaxWidth returns the maximum width constraint if set.
func (s Size) MaxWidth() (int, bool) {
	if s.maxWidth == nil {
		return 0, false
	}
	return *s.maxWidth, true
}

// MinHeight returns the minimum height constraint if set.
func (s Size) MinHeight() (int, bool) {
	if s.minHeight == nil {
		return 0, false
	}
	return *s.minHeight, true
}

// MaxHeight returns the maximum height constraint if set.
func (s Size) MaxHeight() (int, bool) {
	if s.maxHeight == nil {
		return 0, false
	}
	return *s.maxHeight, true
}

// SetWidth returns a new Size with the width set.
// Negative values are clamped to 0.
// This enforces immutability - returns new instance, doesn't modify existing.
func (s Size) SetWidth(w int) Size {
	w = max(0, w)
	return Size{
		width:     &w,
		height:    s.height,
		minWidth:  s.minWidth,
		maxWidth:  s.maxWidth,
		minHeight: s.minHeight,
		maxHeight: s.maxHeight,
	}
}

// SetHeight returns a new Size with the height set.
// Negative values are clamped to 0.
func (s Size) SetHeight(h int) Size {
	h = max(0, h)
	return Size{
		width:     s.width,
		height:    &h,
		minWidth:  s.minWidth,
		maxWidth:  s.maxWidth,
		minHeight: s.minHeight,
		maxHeight: s.maxHeight,
	}
}

// SetMinWidth returns a new Size with the minimum width set.
// Negative values are clamped to 0.
func (s Size) SetMinWidth(w int) Size {
	w = max(0, w)
	return Size{
		width:     s.width,
		height:    s.height,
		minWidth:  &w,
		maxWidth:  s.maxWidth,
		minHeight: s.minHeight,
		maxHeight: s.maxHeight,
	}
}

// SetMaxWidth returns a new Size with the maximum width set.
// Negative values are clamped to 0.
func (s Size) SetMaxWidth(w int) Size {
	w = max(0, w)
	return Size{
		width:     s.width,
		height:    s.height,
		minWidth:  s.minWidth,
		maxWidth:  &w,
		minHeight: s.minHeight,
		maxHeight: s.maxHeight,
	}
}

// SetMinHeight returns a new Size with the minimum height set.
// Negative values are clamped to 0.
func (s Size) SetMinHeight(h int) Size {
	h = max(0, h)
	return Size{
		width:     s.width,
		height:    s.height,
		minWidth:  s.minWidth,
		maxWidth:  s.maxWidth,
		minHeight: &h,
		maxHeight: s.maxHeight,
	}
}

// SetMaxHeight returns a new Size with the maximum height set.
// Negative values are clamped to 0.
func (s Size) SetMaxHeight(h int) Size {
	h = max(0, h)
	return Size{
		width:     s.width,
		height:    s.height,
		minWidth:  s.minWidth,
		maxWidth:  s.maxWidth,
		minHeight: s.minHeight,
		maxHeight: &h,
	}
}

// ValidateWidth clamps a width value to the min/max constraints (if set).
// Returns the clamped width.
func (s Size) ValidateWidth(w int) int {
	// Apply minimum constraint.
	if s.minWidth != nil && w < *s.minWidth {
		w = *s.minWidth
	}

	// Apply maximum constraint.
	if s.maxWidth != nil && w > *s.maxWidth {
		w = *s.maxWidth
	}

	return w
}

// ValidateHeight clamps a height value to the min/max constraints (if set).
// Returns the clamped height.
func (s Size) ValidateHeight(h int) int {
	// Apply minimum constraint.
	if s.minHeight != nil && h < *s.minHeight {
		h = *s.minHeight
	}

	// Apply maximum constraint.
	if s.maxHeight != nil && h > *s.maxHeight {
		h = *s.maxHeight
	}

	return h
}

// String returns a human-readable representation of the size constraints.
func (s Size) String() string {
	return fmt.Sprintf("Size(width=%s, height=%s, minW=%s, maxW=%s, minH=%s, maxH=%s)",
		formatOptionalInt(s.width),
		formatOptionalInt(s.height),
		formatOptionalInt(s.minWidth),
		formatOptionalInt(s.maxWidth),
		formatOptionalInt(s.minHeight),
		formatOptionalInt(s.maxHeight),
	)
}

// --- Private helpers ---.

// formatOptionalInt formats an optional int pointer for display.
func formatOptionalInt(p *int) string {
	if p == nil {
		return "unset"
	}
	return fmt.Sprintf("%d", *p)
}
