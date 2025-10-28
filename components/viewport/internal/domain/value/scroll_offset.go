// Package value contains value objects for viewport domain.
package value

// ScrollOffset represents the current scroll position in a viewport.
// It is immutable and ensures the offset stays within valid bounds.
type ScrollOffset struct {
	offset int
}

// NewScrollOffset creates a new ScrollOffset with the given offset.
// The offset is clamped to be non-negative.
func NewScrollOffset(offset int) *ScrollOffset {
	if offset < 0 {
		offset = 0
	}
	return &ScrollOffset{offset: offset}
}

// Offset returns the current scroll offset value.
func (s *ScrollOffset) Offset() int {
	return s.offset
}

// Add returns a new ScrollOffset with the given delta added.
// The result is clamped between 0 and maxOffset.
func (s *ScrollOffset) Add(delta, maxOffset int) *ScrollOffset {
	newOffset := s.offset + delta
	return s.clamp(newOffset, maxOffset)
}

// Set returns a new ScrollOffset with the offset set to the given value.
// The result is clamped between 0 and maxOffset.
func (s *ScrollOffset) Set(offset, maxOffset int) *ScrollOffset {
	return s.clamp(offset, maxOffset)
}

// Clamp returns a new ScrollOffset with the offset clamped to valid bounds.
func (s *ScrollOffset) Clamp(maxOffset int) *ScrollOffset {
	return s.clamp(s.offset, maxOffset)
}

// clamp ensures the offset is between 0 and maxOffset (inclusive).
func (s *ScrollOffset) clamp(offset, maxOffset int) *ScrollOffset {
	if offset < 0 {
		offset = 0
	}
	if maxOffset < 0 {
		maxOffset = 0
	}
	if offset > maxOffset {
		offset = maxOffset
	}
	return &ScrollOffset{offset: offset}
}
