package value

// Border represents a border style with 8 characters for drawing box borders.
// This is a value object in DDD terms - immutable and defined by its values.
// Each field represents a Unicode character used to draw that part of the border.
type Border struct {
	// Top is the character used for the top horizontal edge.
	Top string

	// Bottom is the character used for the bottom horizontal edge.
	Bottom string

	// Left is the character used for the left vertical edge.
	Left string

	// Right is the character used for the right vertical edge.
	Right string

	// TopLeft is the character used for the top-left corner.
	TopLeft string

	// TopRight is the character used for the top-right corner.
	TopRight string

	// BottomLeft is the character used for the bottom-left corner.
	BottomLeft string

	// BottomRight is the character used for the bottom-right corner.
	BottomRight string
}

// Pre-defined border styles using Unicode box-drawing characters.
// These are commonly used border styles that work in most modern terminals.

var (
	// RoundedBorder uses rounded corners (─ ╭ ╮ ╯ ╰).
	// This is a modern, friendly border style.
	RoundedBorder = Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	// ThickBorder uses thick/bold box-drawing characters (━ ┃).
	// This creates a more prominent, emphasized border.
	ThickBorder = Border{
		Top:         "━",
		Bottom:      "━",
		Left:        "┃",
		Right:       "┃",
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
	}

	// DoubleBorder uses double-line box-drawing characters (═ ║).
	// This creates a classic, formal border style.
	DoubleBorder = Border{
		Top:         "═",
		Bottom:      "═",
		Left:        "║",
		Right:       "║",
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╚",
		BottomRight: "╝",
	}

	// NormalBorder uses standard single-line box-drawing characters (─ │ ┌ ┐ └ ┘).
	// This is the most compatible border style, works in all terminals.
	NormalBorder = Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
	}

	// HiddenBorder uses empty strings for all characters.
	// This effectively creates no border (useful for conditional borders).
	HiddenBorder = Border{
		Top:         "",
		Bottom:      "",
		Left:        "",
		Right:       "",
		TopLeft:     "",
		TopRight:    "",
		BottomLeft:  "",
		BottomRight: "",
	}

	// ASCIIBorder uses ASCII characters (+-|) for maximum compatibility.
	// This works even in terminals without Unicode support.
	ASCIIBorder = Border{
		Top:         "-",
		Bottom:      "-",
		Left:        "|",
		Right:       "|",
		TopLeft:     "+",
		TopRight:    "+",
		BottomLeft:  "+",
		BottomRight: "+",
	}
)

// Equal returns true if this border equals another border.
// Value objects are compared by value, not identity.
func (b Border) Equal(other Border) bool {
	return b.Top == other.Top &&
		b.Bottom == other.Bottom &&
		b.Left == other.Left &&
		b.Right == other.Right &&
		b.TopLeft == other.TopLeft &&
		b.TopRight == other.TopRight &&
		b.BottomLeft == other.BottomLeft &&
		b.BottomRight == other.BottomRight
}

// IsHidden returns true if this is a hidden border (all empty strings).
func (b Border) IsHidden() bool {
	return b.Top == "" &&
		b.Bottom == "" &&
		b.Left == "" &&
		b.Right == "" &&
		b.TopLeft == "" &&
		b.TopRight == "" &&
		b.BottomLeft == "" &&
		b.BottomRight == ""
}

// String returns a human-readable representation of the border.
func (b Border) String() string {
	if b.IsHidden() {
		return "Border(hidden)"
	}
	return "Border(" + b.TopLeft + b.Top + b.TopRight + " / " +
		b.Left + " " + b.Right + " / " +
		b.BottomLeft + b.Bottom + b.BottomRight + ")"
}
