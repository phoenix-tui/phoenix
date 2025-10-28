// Package value provides value objects for the table component.
package value

// Alignment defines how text should be aligned within a cell.
type Alignment int

const (
	// AlignmentLeft aligns text to the left edge (default).
	AlignmentLeft Alignment = iota
	// AlignmentCenter centers text within the cell.
	AlignmentCenter
	// AlignmentRight aligns text to the right edge.
	AlignmentRight
)

// String returns the string representation of the alignment.
func (a Alignment) String() string {
	switch a {
	case AlignmentLeft:
		return "left"
	case AlignmentCenter:
		return "center"
	case AlignmentRight:
		return "right"
	default:
		return "unknown"
	}
}

// IsLeft returns true if alignment is left.
func (a Alignment) IsLeft() bool {
	return a == AlignmentLeft
}

// IsCenter returns true if alignment is center.
func (a Alignment) IsCenter() bool {
	return a == AlignmentCenter
}

// IsRight returns true if alignment is right.
func (a Alignment) IsRight() bool {
	return a == AlignmentRight
}
