package value

// ViewportSize represents the visible area dimensions of a viewport.
// It is immutable and ensures dimensions are non-negative.
type ViewportSize struct {
	width  int
	height int
}

// NewViewportSize creates a new ViewportSize with the given dimensions.
// Width and height are clamped to be non-negative.
func NewViewportSize(width, height int) *ViewportSize {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	return &ViewportSize{width: width, height: height}
}

// Width returns the viewport width.
func (s *ViewportSize) Width() int {
	return s.width
}

// Height returns the viewport height.
func (s *ViewportSize) Height() int {
	return s.height
}

// WithWidth returns a new ViewportSize with the given width.
// The width is clamped to be non-negative.
func (s *ViewportSize) WithWidth(width int) *ViewportSize {
	if width < 0 {
		width = 0
	}
	return &ViewportSize{width: width, height: s.height}
}

// WithHeight returns a new ViewportSize with the given height.
// The height is clamped to be non-negative.
func (s *ViewportSize) WithHeight(height int) *ViewportSize {
	if height < 0 {
		height = 0
	}
	return &ViewportSize{width: s.width, height: height}
}
