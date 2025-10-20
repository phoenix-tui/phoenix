// Package value contains value objects for text alignment.
package value

import "fmt"

// HorizontalAlignment represents horizontal text alignment.
type HorizontalAlignment int

// Horizontal alignment constants.
const (
	AlignLeft   HorizontalAlignment = iota // AlignLeft aligns text to the left.
	AlignCenter                            // AlignCenter centers text horizontally.
	AlignRight                             // AlignRight aligns text to the right.
)

// String returns a human-readable representation of the horizontal alignment.
func (h HorizontalAlignment) String() string {
	switch h {
	case AlignLeft:
		return "Left"
	case AlignCenter:
		return "Center"
	case AlignRight:
		return "Right"
	default:
		//nolint:goconst // "Unknown" string literal used in switch defaults - constant not beneficial
		return "Unknown"
	}
}

// VerticalAlignment represents vertical text alignment.
type VerticalAlignment int

// Vertical alignment constants.
const (
	AlignTop    VerticalAlignment = iota // AlignTop aligns text to the top.
	AlignMiddle                          // AlignMiddle centers text vertically.
	AlignBottom                          // AlignBottom aligns text to the bottom.
)

// String returns a human-readable representation of the vertical alignment.
func (v VerticalAlignment) String() string {
	switch v {
	case AlignTop:
		return "Top"
	case AlignMiddle:
		return "Middle"
	case AlignBottom:
		return "Bottom"
	default:
		return "Unknown"
	}
}

// Alignment represents an immutable combination of horizontal and vertical alignment.
// This is a value object in DDD terms - immutable and defined by its values.
type Alignment struct {
	horizontal HorizontalAlignment
	vertical   VerticalAlignment
}

// NewAlignment creates a new Alignment with specified horizontal and vertical alignment.
func NewAlignment(h HorizontalAlignment, v VerticalAlignment) Alignment {
	return Alignment{
		horizontal: h,
		vertical:   v,
	}
}

// --- Convenience constructors for common alignments ---.

// LeftTop creates an alignment with left-top positioning.
func LeftTop() Alignment {
	return Alignment{horizontal: AlignLeft, vertical: AlignTop}
}

// LeftMiddle creates an alignment with left-middle positioning.
func LeftMiddle() Alignment {
	return Alignment{horizontal: AlignLeft, vertical: AlignMiddle}
}

// LeftBottom creates an alignment with left-bottom positioning.
func LeftBottom() Alignment {
	return Alignment{horizontal: AlignLeft, vertical: AlignBottom}
}

// CenterTop creates an alignment with center-top positioning.
func CenterTop() Alignment {
	return Alignment{horizontal: AlignCenter, vertical: AlignTop}
}

// CenterMiddle creates an alignment with center-middle positioning (centered both ways).
func CenterMiddle() Alignment {
	return Alignment{horizontal: AlignCenter, vertical: AlignMiddle}
}

// CenterBottom creates an alignment with center-bottom positioning.
func CenterBottom() Alignment {
	return Alignment{horizontal: AlignCenter, vertical: AlignBottom}
}

// RightTop creates an alignment with right-top positioning.
func RightTop() Alignment {
	return Alignment{horizontal: AlignRight, vertical: AlignTop}
}

// RightMiddle creates an alignment with right-middle positioning.
func RightMiddle() Alignment {
	return Alignment{horizontal: AlignRight, vertical: AlignMiddle}
}

// RightBottom creates an alignment with right-bottom positioning.
func RightBottom() Alignment {
	return Alignment{horizontal: AlignRight, vertical: AlignBottom}
}

// Horizontal returns the horizontal alignment component.
func (a Alignment) Horizontal() HorizontalAlignment {
	return a.horizontal
}

// Vertical returns the vertical alignment component.
func (a Alignment) Vertical() VerticalAlignment {
	return a.vertical
}

// Equal returns true if this alignment equals another alignment.
// Value objects are compared by value, not identity.
func (a Alignment) Equal(other Alignment) bool {
	return a.horizontal == other.horizontal && a.vertical == other.vertical
}

// String returns a human-readable representation of the alignment.
func (a Alignment) String() string {
	return fmt.Sprintf("Alignment(%s, %s)", a.horizontal.String(), a.vertical.String())
}
