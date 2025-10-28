package value

// ColorDepth represents terminal color support level.
type ColorDepth int

const (
	// ColorDepthNone indicates no color support (monochrome terminal).
	ColorDepthNone ColorDepth = 0

	// ColorDepth8 indicates 8-color support (3-bit color).
	// Standard ANSI colors: black, red, green, yellow, blue, magenta, cyan, white.
	ColorDepth8 ColorDepth = 8

	// ColorDepth256 indicates 256-color support (8-bit color).
	// 216 colors (6×6×6 cube) + 16 system colors + 24 grayscale.
	ColorDepth256 ColorDepth = 256

	// ColorDepthTrueColor indicates 24-bit RGB color support.
	// 16,777,216 possible colors (256×256×256).
	ColorDepthTrueColor ColorDepth = 16777216
)

// String returns human-readable color depth description.
func (cd ColorDepth) String() string {
	switch cd {
	case ColorDepthNone:
		return "no-color"
	case ColorDepth8:
		return "8-color"
	case ColorDepth256:
		return "256-color"
	case ColorDepthTrueColor:
		return "truecolor"
	default:
		return "unknown"
	}
}

// Capabilities represents terminal capabilities (immutable value object).
//
// This encapsulates what the terminal can do, detected from environment
// variables (TERM, COLORTERM) and terminal queries.
//
// Invariants:
//   - All fields are immutable after creation
//   - ColorDepth is one of the defined constants
//   - Capabilities is determined at terminal initialization
type Capabilities struct {
	ansiSupport   bool       // Terminal supports ANSI escape sequences
	colorDepth    ColorDepth // Color support level
	mouseSupport  bool       // Terminal supports mouse events (SGR mouse mode)
	altScreen     bool       // Terminal supports alternate screen buffer
	cursorControl bool       // Terminal supports cursor positioning/visibility
}

// NewCapabilities creates capabilities with validation.
func NewCapabilities(ansi bool, colors ColorDepth, mouse, alt, cursor bool) *Capabilities {
	// Business rule: mouse/alt/cursor require ANSI support
	if !ansi {
		mouse = false
		alt = false
		cursor = false
	}

	// Business rule: validate color depth
	switch colors {
	case ColorDepthNone, ColorDepth8, ColorDepth256, ColorDepthTrueColor:
		// Valid
	default:
		colors = ColorDepthNone // Invalid depth → no color
	}

	return &Capabilities{
		ansiSupport:   ansi,
		colorDepth:    colors,
		mouseSupport:  mouse,
		altScreen:     alt,
		cursorControl: cursor,
	}
}

// SupportsANSI returns true if terminal supports ANSI escape sequences.
func (c *Capabilities) SupportsANSI() bool {
	return c.ansiSupport
}

// ColorDepth returns terminal color support level.
func (c *Capabilities) ColorDepth() ColorDepth {
	return c.colorDepth
}

// SupportsColor returns true if terminal supports any colors.
func (c *Capabilities) SupportsColor() bool {
	return c.colorDepth > ColorDepthNone
}

// SupportsTrueColor returns true if terminal supports 24-bit RGB colors.
func (c *Capabilities) SupportsTrueColor() bool {
	return c.colorDepth == ColorDepthTrueColor
}

// Supports256Color returns true if terminal supports at least 256 colors.
func (c *Capabilities) Supports256Color() bool {
	return c.colorDepth >= ColorDepth256
}

// SupportsMouse returns true if terminal supports mouse events.
func (c *Capabilities) SupportsMouse() bool {
	return c.mouseSupport
}

// SupportsAltScreen returns true if terminal supports alternate screen buffer.
func (c *Capabilities) SupportsAltScreen() bool {
	return c.altScreen
}

// SupportsCursorControl returns true if terminal supports cursor positioning.
func (c *Capabilities) SupportsCursorControl() bool {
	return c.cursorControl
}

// IsDumbTerminal returns true if terminal has no special capabilities.
func (c *Capabilities) IsDumbTerminal() bool {
	return !c.ansiSupport && c.colorDepth == ColorDepthNone
}

// Equal returns true if capabilities are equal.
func (c *Capabilities) Equal(other *Capabilities) bool {
	if other == nil {
		return false
	}
	return c.ansiSupport == other.ansiSupport &&
		c.colorDepth == other.colorDepth &&
		c.mouseSupport == other.mouseSupport &&
		c.altScreen == other.altScreen &&
		c.cursorControl == other.cursorControl
}
