// Package value provides value objects for textarea domain model.
package value

// Range represents a range of text (start position to end position).
// This is a value object - immutable and comparable by value.
type Range struct {
	start Position
	end   Position
}

// NewRange creates a new range from start to end.
// The range is normalized so start is always before or equal to end.
func NewRange(start, end Position) Range {
	// Normalize: ensure start <= end.
	if end.IsBefore(start) {
		start, end = end, start
	}
	return Range{start: start, end: end}
}

// Start returns the start position.
func (r Range) Start() Position {
	return r.start
}

// End returns the end position.
func (r Range) End() Position {
	return r.end
}

// StartRowCol returns start position as (row, col).
func (r Range) StartRowCol() (row, col int) {
	return r.start.Row(), r.start.Col()
}

// EndRowCol returns end position as (row, col).
func (r Range) EndRowCol() (row, col int) {
	return r.end.Row(), r.end.Col()
}

// Contains returns true if position is within range.
func (r Range) Contains(pos Position) bool {
	return !pos.IsBefore(r.start) && !pos.IsAfter(r.end)
}

// IsEmpty returns true if start equals end.
func (r Range) IsEmpty() bool {
	return r.start.Equals(r.end)
}

// IsSingleLine returns true if range is within single line.
func (r Range) IsSingleLine() bool {
	return r.start.Row() == r.end.Row()
}
