package value

// TerminalCapability represents the color support level of a terminal.
// This is used to adapt colors to the terminal's capabilities.
type TerminalCapability int

const (
	// NoColor indicates no color support (monochrome terminal).
	// All colors will be rendered as plain text without ANSI codes.
	NoColor TerminalCapability = iota

	// ANSI16 indicates support for 16 basic ANSI colors (8 normal + 8 bright).
	// Uses ANSI codes 30-37 (foreground) and 40-47 (background).
	ANSI16

	// ANSI256 indicates support for 256 colors.
	// Uses ANSI codes 38;5;N (foreground) and 48;5;N (background).
	ANSI256

	// TrueColor indicates support for 24-bit RGB colors (16.7 million colors).
	// Uses ANSI codes 38;2;R;G;B (foreground) and 48;2;R;G;B (background).
	TrueColor
)

// String returns a human-readable name for the terminal capability.
func (tc TerminalCapability) String() string {
	switch tc {
	case NoColor:
		return "NoColor"
	case ANSI16:
		return "ANSI16"
	case ANSI256:
		return "ANSI256"
	case TrueColor:
		return "TrueColor"
	default:
		return "Unknown"
	}
}

// SupportsColor returns true if the terminal supports any color.
func (tc TerminalCapability) SupportsColor() bool {
	return tc != NoColor
}

// SupportsTrueColor returns true if the terminal supports 24-bit RGB colors.
func (tc TerminalCapability) SupportsTrueColor() bool {
	return tc == TrueColor
}

// Supports256Color returns true if the terminal supports 256 colors or better.
func (tc TerminalCapability) Supports256Color() bool {
	return tc == ANSI256 || tc == TrueColor
}

// Supports16Color returns true if the terminal supports at least 16 colors.
func (tc TerminalCapability) Supports16Color() bool {
	return tc == ANSI16 || tc == ANSI256 || tc == TrueColor
}
