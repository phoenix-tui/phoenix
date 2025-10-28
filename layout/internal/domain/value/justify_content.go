package value

// JustifyContent defines how items are distributed along the main axis.
//
// Design Philosophy:
//   - Simplified from CSS Flexbox (most common cases only)
//   - Controls spacing between items
//   - Applied along main axis (row = horizontal, column = vertical)
//
// Supported Values:
//   - Start: Items packed at start (default)
//   - End: Items packed at end
//   - Center: Items centered
//   - SpaceBetween: Equal spacing between items (no space at edges)
//
// Deferred to v0.2.0+:
//   - ❌ SpaceAround
//   - ❌ SpaceEvenly
//
// Example:
//
//	justify := JustifyContentCenter
//	if justify == JustifyContentSpaceBetween {
//	    // Calculate equal gaps between items
//	}
type JustifyContent int

const (
	// JustifyContentStart packs items at the start of the container.
	// This is the default value.
	//
	// Visual (Row):
	//   [1][2][3]         (remaining space)
	//
	// Visual (Column):
	//   [1]
	//   [2]
	//   [3]
	//   (remaining space)
	JustifyContentStart JustifyContent = iota

	// JustifyContentEnd packs items at the end of the container.
	//
	// Visual (Row):
	//   (remaining space)         [1][2][3]
	//
	// Visual (Column):
	//   (remaining space)
	//   [1]
	//   [2]
	//   [3]
	JustifyContentEnd

	// JustifyContentCenter centers items in the container.
	//
	// Visual (Row):
	//   (space)    [1][2][3]    (space)
	//
	// Visual (Column):
	//   (space)
	//   [1]
	//   [2]
	//   [3]
	//   (space)
	JustifyContentCenter

	// JustifyContentSpaceBetween distributes items with equal spacing.
	// First item at start, last item at end, equal gaps between.
	//
	// Visual (Row):
	//   [1]    (gap)    [2]    (gap)    [3]
	//
	// Visual (Column):
	//   [1]
	//   (gap)
	//   [2]
	//   (gap)
	//   [3]
	JustifyContentSpaceBetween
)

// String returns a human-readable representation.
func (j JustifyContent) String() string {
	switch j {
	case JustifyContentStart:
		return "start"
	case JustifyContentEnd:
		return "end"
	case JustifyContentCenter:
		return "center"
	case JustifyContentSpaceBetween:
		return "space-between"
	default:
		return "unknown"
	}
}

// Validate checks if the justify content value is valid.
func (j JustifyContent) Validate() bool {
	return j >= JustifyContentStart && j <= JustifyContentSpaceBetween
}

// IsDefault returns true if this is the default value (Start).
func (j JustifyContent) IsDefault() bool {
	return j == JustifyContentStart
}

// NeedsDistribution returns true if spacing calculation is needed.
// Start/End don't need distribution, Center/SpaceBetween do.
func (j JustifyContent) NeedsDistribution() bool {
	return j == JustifyContentCenter || j == JustifyContentSpaceBetween
}
