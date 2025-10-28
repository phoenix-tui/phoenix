package value

// Selection represents a text selection range in grapheme cluster offsets.
// The start and end positions are measured in grapheme clusters, not bytes or runes.
type Selection struct {
	start int // Start grapheme offset (inclusive)
	end   int // End grapheme offset (exclusive)
}

// NewSelection creates a new selection with the specified range.
// The start and end values are normalized so that start <= end.
// Negative values are clamped to 0.
func NewSelection(start, end int) *Selection {
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}

	// Normalize so start <= end.
	if start > end {
		start, end = end, start
	}

	return &Selection{
		start: start,
		end:   end,
	}
}

// Start returns the start position of the selection (inclusive).
func (s *Selection) Start() int {
	return s.start
}

// End returns the end position of the selection (exclusive).
func (s *Selection) End() int {
	return s.end
}

// Length returns the number of grapheme clusters in the selection.
func (s *Selection) Length() int {
	return s.end - s.start
}

// IsEmpty returns true if the selection has zero length.
func (s *Selection) IsEmpty() bool {
	return s.start == s.end
}

// Contains returns true if the given offset is within the selection range.
func (s *Selection) Contains(offset int) bool {
	return offset >= s.start && offset < s.end
}

// Clamp restricts the selection to the given maximum offset.
// Returns a new Selection instance (immutable).
func (s *Selection) Clamp(maxOffset int) *Selection {
	start := s.start
	end := s.end

	if start > maxOffset {
		start = maxOffset
	}
	if end > maxOffset {
		end = maxOffset
	}

	return &Selection{
		start: start,
		end:   end,
	}
}

// Clone creates a copy of this selection.
func (s *Selection) Clone() *Selection {
	return &Selection{
		start: s.start,
		end:   s.end,
	}
}
