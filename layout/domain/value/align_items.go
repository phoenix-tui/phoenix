package value

// AlignItems defines how items are aligned along the cross axis.
//
// Design Philosophy:
//   - Simplified from CSS Flexbox (most common cases only)
//   - Controls alignment perpendicular to main axis
//   - Applied to cross axis (row = vertical, column = horizontal)
//
// Supported Values:
//   - Start: Items aligned at start
//   - End: Items aligned at end
//   - Center: Items centered
//   - Stretch: Items stretched to fill container (default)
//
// Deferred to v0.2.0+:
//   - ❌ Baseline alignment
//
// Example:
//
//	align := AlignItemsCenter
//	if align == AlignItemsStretch {
//	    // Make item fill available cross-axis space
//	}
type AlignItems int

const (
	// AlignItemsStretch stretches items to fill the cross axis.
	// This is the default value.
	//
	// Visual (Row, items with different heights):
	//   ┌───┐ ┌───┐ ┌───┐
	//   │ 1 │ │ 2 │ │ 3 │  ← All stretched to same height
	//   └───┘ └───┘ └───┘
	AlignItemsStretch AlignItems = iota

	// AlignItemsStart aligns items at the start of the cross axis.
	//
	// Visual (Row, items with different heights):
	//   ┌───┐ ┌───┐ ┌───┐
	//   │ 1 │ │ 2 │ │ 3 │  ← Aligned to top
	//   │   │ └───┘ │   │
	//   └───┘       └───┘
	AlignItemsStart

	// AlignItemsEnd aligns items at the end of the cross axis.
	//
	// Visual (Row, items with different heights):
	//   ┌───┐       ┌───┐
	//   │ 1 │ ┌───┐ │ 3 │
	//   │   │ │ 2 │ │   │  ← Aligned to bottom
	//   └───┘ └───┘ └───┘
	AlignItemsEnd

	// AlignItemsCenter centers items along the cross axis.
	//
	// Visual (Row, items with different heights):
	//   ┌───┐
	//   │ 1 │ ┌───┐ ┌───┐
	//   │   │ │ 2 │ │ 3 │  ← Centered vertically
	//   └───┘ └───┘ │   │
	//               └───┘
	AlignItemsCenter
)

// String returns a human-readable representation.
func (a AlignItems) String() string {
	switch a {
	case AlignItemsStretch:
		return "stretch"
	case AlignItemsStart:
		return "start"
	case AlignItemsEnd:
		return "end"
	case AlignItemsCenter:
		return "center"
	default:
		return "unknown"
	}
}

// Validate checks if the align items value is valid.
func (a AlignItems) Validate() bool {
	return a >= AlignItemsStretch && a <= AlignItemsCenter
}

// IsDefault returns true if this is the default value (Stretch).
func (a AlignItems) IsDefault() bool {
	return a == AlignItemsStretch
}

// RequiresStretching returns true if items should be stretched.
func (a AlignItems) RequiresStretching() bool {
	return a == AlignItemsStretch
}

// RequiresAlignment returns true if items need cross-axis positioning.
// Stretch doesn't need alignment (items fill space), others do.
func (a AlignItems) RequiresAlignment() bool {
	return a != AlignItemsStretch
}
