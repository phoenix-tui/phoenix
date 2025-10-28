// Package value provides value objects for the layout domain.
package value

// FlexDirection defines the main axis direction for flexbox layout.
//
// Design Philosophy:
//   - Enum-like type (type-safe)
//   - Only two directions: Row (horizontal) and Column (vertical)
//   - No wrap support in v0.1.0 (deferred to v0.2.0+)
//
// Simplified from CSS Flexbox:
//   - ✅ Row (left-to-right)
//   - ✅ Column (top-to-bottom)
//   - ❌ Row-reverse (v0.2.0+)
//   - ❌ Column-reverse (v0.2.0+)
//   - ❌ Wrap (v0.2.0+)
//
// Example:
//
//	direction := FlexDirectionRow
//	if direction.IsHorizontal() {
//	    // Layout children left-to-right
//	}
type FlexDirection int

const (
	// FlexDirectionRow arranges items horizontally (left-to-right).
	// This is the default direction for flexbox containers.
	//
	// Visual:
	//   [Item1] [Item2] [Item3]
	FlexDirectionRow FlexDirection = iota

	// FlexDirectionColumn arranges items vertically (top-to-bottom).
	//
	// Visual:
	//   [Item1]
	//   [Item2]
	//   [Item3]
	FlexDirectionColumn
)

// IsHorizontal returns true if direction is row-based.
func (d FlexDirection) IsHorizontal() bool {
	return d == FlexDirectionRow
}

// IsVertical returns true if direction is column-based.
func (d FlexDirection) IsVertical() bool {
	return d == FlexDirectionColumn
}

// String returns a human-readable representation.
func (d FlexDirection) String() string {
	switch d {
	case FlexDirectionRow:
		return "row"
	case FlexDirectionColumn:
		return "column"
	default:
		return "unknown"
	}
}

// Validate checks if the direction value is valid.
// Returns true for Row or Column, false otherwise.
func (d FlexDirection) Validate() bool {
	return d == FlexDirectionRow || d == FlexDirectionColumn
}
