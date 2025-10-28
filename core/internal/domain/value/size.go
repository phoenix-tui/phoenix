package value

// Size represents terminal dimensions (width, height in cells).
// This is an immutable value object.
//
// Invariants:
//   - Width and Height are always >= 1
//   - Size is immutable after creation
type Size struct {
	Width  int // Terminal width in cells
	Height int // Terminal height in cells
}

// NewSize creates size with validation.
// Values < 1 are clamped to 1 (minimum terminal size).
func NewSize(width, height int) Size {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	return Size{Width: width, Height: height}
}

// Area returns total number of cells (width * height).
func (s Size) Area() int {
	return s.Width * s.Height
}

// Contains checks if a position is within terminal bounds.
func (s Size) Contains(pos Position) bool {
	return pos.Row >= 0 && pos.Row < s.Height &&
		pos.Col >= 0 && pos.Col < s.Width
}

// Equal returns true if sizes are equal.
func (s Size) Equal(other Size) bool {
	return s.Width == other.Width && s.Height == other.Height
}

// IsEmpty returns true if size has zero area (shouldn't happen due to validation).
func (s Size) IsEmpty() bool {
	return s.Width == 0 || s.Height == 0
}
