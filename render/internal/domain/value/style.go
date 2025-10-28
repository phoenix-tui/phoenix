package value

import (
	"fmt"
	"strings"
)

// Style represents ANSI text styling.
// Style is immutable and uses value semantics.
type Style struct {
	fg        *Color // Foreground color (nil = default)
	bg        *Color // Background color (nil = default)
	bold      bool
	italic    bool
	underline bool
	reverse   bool
	dim       bool
	blink     bool
	hidden    bool
	strike    bool
}

// NewStyle creates a new empty style.
func NewStyle() Style {
	return Style{}
}

// NewStyleWithFg creates a style with foreground color.
func NewStyleWithFg(fg Color) Style {
	return Style{fg: &fg}
}

// NewStyleWithBg creates a style with background color.
func NewStyleWithBg(bg Color) Style {
	return Style{bg: &bg}
}

// NewStyleWithColors creates a style with both colors.
func NewStyleWithColors(fg, bg Color) Style {
	return Style{fg: &fg, bg: &bg}
}

// Foreground returns the foreground color (nil if not set).
func (s Style) Foreground() *Color {
	return s.fg
}

// Background returns the background color (nil if not set).
func (s Style) Background() *Color {
	return s.bg
}

// Bold returns true if bold is enabled.
func (s Style) Bold() bool {
	return s.bold
}

// Italic returns true if italic is enabled.
func (s Style) Italic() bool {
	return s.italic
}

// Underline returns true if underline is enabled.
func (s Style) Underline() bool {
	return s.underline
}

// Reverse returns true if reverse video is enabled.
func (s Style) Reverse() bool {
	return s.reverse
}

// Dim returns true if dim is enabled.
func (s Style) Dim() bool {
	return s.dim
}

// Blink returns true if blink is enabled.
func (s Style) Blink() bool {
	return s.blink
}

// Hidden returns true if hidden is enabled.
func (s Style) Hidden() bool {
	return s.hidden
}

// Strike returns true if strikethrough is enabled.
func (s Style) Strike() bool {
	return s.strike
}

// WithFg returns a new style with foreground color.
func (s Style) WithFg(fg Color) Style {
	s.fg = &fg
	return s
}

// WithBg returns a new style with background color.
func (s Style) WithBg(bg Color) Style {
	s.bg = &bg
	return s
}

// WithBold returns a new style with bold setting.
func (s Style) WithBold(bold bool) Style {
	s.bold = bold
	return s
}

// WithItalic returns a new style with italic setting.
func (s Style) WithItalic(italic bool) Style {
	s.italic = italic
	return s
}

// WithUnderline returns a new style with underline setting.
func (s Style) WithUnderline(underline bool) Style {
	s.underline = underline
	return s
}

// WithReverse returns a new style with reverse video setting.
func (s Style) WithReverse(reverse bool) Style {
	s.reverse = reverse
	return s
}

// WithDim returns a new style with dim setting.
func (s Style) WithDim(dim bool) Style {
	s.dim = dim
	return s
}

// WithBlink returns a new style with blink setting.
func (s Style) WithBlink(blink bool) Style {
	s.blink = blink
	return s
}

// WithHidden returns a new style with hidden setting.
func (s Style) WithHidden(hidden bool) Style {
	s.hidden = hidden
	return s
}

// WithStrike returns a new style with strikethrough setting.
func (s Style) WithStrike(strike bool) Style {
	s.strike = strike
	return s
}

// Equals checks if two styles are equal.
//
//nolint:gocyclo,cyclop // Style equality requires checking all fields
func (s Style) Equals(other Style) bool {
	// Compare colors.
	if (s.fg == nil) != (other.fg == nil) {
		return false
	}
	if s.fg != nil && !s.fg.Equals(*other.fg) {
		return false
	}

	if (s.bg == nil) != (other.bg == nil) {
		return false
	}
	if s.bg != nil && !s.bg.Equals(*other.bg) {
		return false
	}

	// Compare attributes.
	return s.bold == other.bold &&
		s.italic == other.italic &&
		s.underline == other.underline &&
		s.reverse == other.reverse &&
		s.dim == other.dim &&
		s.blink == other.blink &&
		s.hidden == other.hidden &&
		s.strike == other.strike
}

// IsEmpty returns true if style has no attributes set.
func (s Style) IsEmpty() bool {
	return s.fg == nil && s.bg == nil &&
		!s.bold && !s.italic && !s.underline && !s.reverse &&
		!s.dim && !s.blink && !s.hidden && !s.strike
}

// ToANSI generates ANSI escape sequence for this style.
// Returns empty string if style is empty.
func (s Style) ToANSI() string {
	if s.IsEmpty() {
		return ""
	}

	var codes []string

	// Foreground color.
	if s.fg != nil {
		r, g, b := s.fg.RGB()
		codes = append(codes, fmt.Sprintf("38;2;%d;%d;%d", r, g, b))
	}

	// Background color.
	if s.bg != nil {
		r, g, b := s.bg.RGB()
		codes = append(codes, fmt.Sprintf("48;2;%d;%d;%d", r, g, b))
	}

	// Attributes.
	if s.bold {
		codes = append(codes, "1")
	}
	if s.dim {
		codes = append(codes, "2")
	}
	if s.italic {
		codes = append(codes, "3")
	}
	if s.underline {
		codes = append(codes, "4")
	}
	if s.blink {
		codes = append(codes, "5")
	}
	if s.reverse {
		codes = append(codes, "7")
	}
	if s.hidden {
		codes = append(codes, "8")
	}
	if s.strike {
		codes = append(codes, "9")
	}

	return "\x1b[" + strings.Join(codes, ";") + "m"
}

// String returns a string representation for debugging.
func (s Style) String() string {
	var parts []string

	if s.fg != nil {
		parts = append(parts, fmt.Sprintf("fg:%s", s.fg))
	}
	if s.bg != nil {
		parts = append(parts, fmt.Sprintf("bg:%s", s.bg))
	}
	if s.bold {
		parts = append(parts, "bold")
	}
	if s.italic {
		parts = append(parts, "italic")
	}
	if s.underline {
		parts = append(parts, "underline")
	}
	if s.reverse {
		parts = append(parts, "reverse")
	}
	if s.dim {
		parts = append(parts, "dim")
	}
	if s.blink {
		parts = append(parts, "blink")
	}
	if s.hidden {
		parts = append(parts, "hidden")
	}
	if s.strike {
		parts = append(parts, "strike")
	}

	if len(parts) == 0 {
		return "Style(empty)"
	}
	return "Style(" + strings.Join(parts, ", ") + ")"
}
