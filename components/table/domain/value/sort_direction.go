package value

// SortDirection defines the direction of sorting.
type SortDirection int

const (
	// SortDirectionNone indicates no sorting is applied.
	SortDirectionNone SortDirection = iota
	// SortDirectionAsc sorts in ascending order (A-Z, 0-9).
	SortDirectionAsc
	// SortDirectionDesc sorts in descending order (Z-A, 9-0).
	SortDirectionDesc
)

// String returns the string representation of the sort direction.
func (d SortDirection) String() string {
	switch d {
	case SortDirectionNone:
		return "none"
	case SortDirectionAsc:
		return "asc"
	case SortDirectionDesc:
		return "desc"
	default:
		return "unknown"
	}
}

// IsAscending returns true if direction is ascending.
func (d SortDirection) IsAscending() bool {
	return d == SortDirectionAsc
}

// IsDescending returns true if direction is descending.
func (d SortDirection) IsDescending() bool {
	return d == SortDirectionDesc
}

// IsNone returns true if no sorting is applied.
func (d SortDirection) IsNone() bool {
	return d == SortDirectionNone
}

// Toggle switches between ascending and descending.
// If current direction is None, returns Asc.
func (d SortDirection) Toggle() SortDirection {
	switch d {
	case SortDirectionAsc:
		return SortDirectionDesc
	case SortDirectionDesc:
		return SortDirectionAsc
	default:
		return SortDirectionAsc
	}
}
