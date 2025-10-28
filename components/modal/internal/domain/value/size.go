package value

// Size represents modal dimensions (width and height).
type Size struct {
	width  int
	height int
}

// NewSize creates a new size with the given width and height.
// Negative or zero values are allowed (will be clamped by rendering logic).
func NewSize(width, height int) *Size {
	return &Size{
		width:  width,
		height: height,
	}
}

// Width returns the width in characters.
func (s *Size) Width() int {
	return s.width
}

// Height returns the height in rows.
func (s *Size) Height() int {
	return s.height
}

// WithWidth returns a new size with the specified width.
func (s *Size) WithWidth(width int) *Size {
	return &Size{
		width:  width,
		height: s.height,
	}
}

// WithHeight returns a new size with the specified height.
func (s *Size) WithHeight(height int) *Size {
	return &Size{
		width:  s.width,
		height: height,
	}
}
