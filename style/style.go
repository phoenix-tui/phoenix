// Package style provides CSS-like styling for Phoenix TUI framework with correct Unicode handling.
//
// # Overview
//
// Package style offers a comprehensive styling system for terminal UIs:
//   - Colors (foreground, background) with TrueColor/ANSI256/ANSI16 support
//   - Text attributes (bold, italic, underline, strikethrough, reverse, blink, dim)
//   - Borders (solid, rounded, double, thick, custom) with Unicode box-drawing
//   - Spacing (padding, margin) with fine-grained control (top, right, bottom, left)
//   - Alignment (horizontal, vertical) with width/height constraints
//   - 8-stage rendering pipeline for optimal performance
//
// # Features
//
//   - CSS-like fluent API with method chaining
//   - Unicode-correct rendering (fixes Lipgloss #562 - correct emoji/CJK width)
//   - Immutable styles (thread-safe, safe for concurrent use)
//   - Terminal capability adaptation (degrades gracefully: TrueColor → ANSI256 → ANSI16)
//   - Zero allocations in hot paths (optimized for performance)
//   - DDD architecture (clean separation of concerns)
//
// # Quick Start
//
// Basic styling:
//
//	import "github.com/phoenix-tui/phoenix/style"
//
//	// Create a style
//	s := style.New().
//	    Foreground(style.RGB(255, 255, 255)).
//	    Background(style.RGB(0, 0, 255)).
//	    Bold(true)
//
//	// Render styled content
//	output := s.Render("Hello, World!")
//	fmt.Println(output)
//
// Box with borders and padding:
//
//	s := style.New().
//	    Border(style.RoundedBorder).
//	    BorderForeground(style.RGB(100, 200, 255)).
//	    Padding(1, 2).  // vertical, horizontal
//	    Width(40)
//
//	fmt.Println(s.Render("Content inside box"))
//
// # Architecture
//
// This package follows Domain-Driven Design (DDD):
//   - internal/domain/model    - Style aggregate, rendering rules
//   - internal/domain/value    - Color, Border, Spacing value objects
//   - internal/domain/service  - Rendering engine, Unicode handling
//   - internal/application     - Use case commands (RenderCommand)
//   - internal/infrastructure  - ANSI escape sequences, caching
//   - style.go (this file)     - Public API (wrapper types)
//
// The 8-stage rendering pipeline:
//  1. Content preparation (Unicode normalization)
//  2. Alignment (horizontal/vertical)
//  3. Padding application
//  4. Border rendering (Unicode box-drawing)
//  5. Margin application
//  6. Color application (with capability adaptation)
//  7. Text attributes (bold, italic, etc.)
//  8. Final ANSI sequence generation
//
// # Performance
//
// Rendering is optimized for speed:
//   - Zero allocations on cached paths
//   - Pre-computed Unicode widths
//   - Efficient string building
package style

import (
	"github.com/phoenix-tui/phoenix/style/internal/application/command"
	"github.com/phoenix-tui/phoenix/style/internal/domain/model"
	service2 "github.com/phoenix-tui/phoenix/style/internal/domain/service"
	value2 "github.com/phoenix-tui/phoenix/style/internal/domain/value"
	"github.com/phoenix-tui/phoenix/style/internal/infrastructure/ansi"
)

// Style is an alias for model.Style, the main styling configuration.
//
// Zero value: Style with all nil pointers and false booleans is valid but applies no styling.
// Use New() for explicit initialization (recommended).
//
//	var s style.Style      // Zero value - valid, no styling applied
//	s2 := style.New()      // Explicit - same as zero value but more readable
//
// Thread safety: Style is immutable and safe for concurrent use.
// All setter methods return new instances.
type Style = model.Style

// Color is an alias for value.Color, representing terminal colors.
//
// Zero value: Color with RGB(0, 0, 0) represents black.
// Use RGB(), Hex(), or Color256() constructors for explicit colors.
//
//	var c style.Color          // Zero value - black (0, 0, 0)
//	c2 := style.RGB(255, 0, 0) // Explicit - red
//
// Thread safety: Color is immutable and safe for concurrent use.
type Color = value2.Color

// Aliases for struct value types (these are fine as aliases - methods are visible).
type (
	// Border is an alias for value.Border, representing box borders.
	//
	// Zero value: Border with empty strings is not useful.
	// Use predefined borders (NormalBorder, RoundedBorder, etc.) or NewBorder().
	Border = value2.Border

	// Padding is an alias for value.Padding, representing box padding.
	//
	// Zero value: Padding{Top: 0, Right: 0, Bottom: 0, Left: 0} is valid (no padding).
	// Use NewPadding() for explicit values.
	Padding = value2.Padding

	// Margin is an alias for value.Margin, representing box margins.
	//
	// Zero value: Margin{Top: 0, Right: 0, Bottom: 0, Left: 0} is valid (no margin).
	// Use NewMargin() for explicit values.
	Margin = value2.Margin

	// Size is an alias for value.Size, representing box dimensions.
	//
	// Zero value: Size with no constraints is valid (content-sized).
	// Use NewSize() for explicit constraints.
	Size = value2.Size

	// Alignment is an alias for value.Alignment, representing text alignment.
	//
	// Zero value: Alignment with left-top is valid (default alignment).
	// Use NewAlignment() for explicit alignment.
	Alignment = value2.Alignment
)

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
	internal := value2.HorizontalAlignment(h)
	return internal.String()
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
	internal := value2.VerticalAlignment(v)
	return internal.String()
}

// TerminalCapability represents the color support level of a terminal.
// This is used to adapt colors to the terminal's capabilities.
type TerminalCapability int

// Terminal capability constants define color support levels.
const (
	NoColor   TerminalCapability = iota // NoColor indicates no color support (monochrome).
	ANSI16                              // ANSI16 supports 16 basic colors (8 normal + 8 bright).
	ANSI256                             // ANSI256 supports 256-color palette.
	TrueColor                           // TrueColor supports 24-bit RGB (16.7 million colors).
)

// String returns a human-readable representation of the terminal capability.
func (tc TerminalCapability) String() string {
	internal := value2.TerminalCapability(tc)
	return internal.String()
}

// SupportsColor returns true if the terminal supports any color.
func (tc TerminalCapability) SupportsColor() bool {
	internal := value2.TerminalCapability(tc)
	return internal.SupportsColor()
}

// SupportsTrueColor returns true if the terminal supports 24-bit RGB colors.
func (tc TerminalCapability) SupportsTrueColor() bool {
	internal := value2.TerminalCapability(tc)
	return internal.SupportsTrueColor()
}

// Supports256Color returns true if the terminal supports 256 colors or better.
func (tc TerminalCapability) Supports256Color() bool {
	internal := value2.TerminalCapability(tc)
	return internal.Supports256Color()
}

// Supports16Color returns true if the terminal supports at least 16 colors.
func (tc TerminalCapability) Supports16Color() bool {
	internal := value2.TerminalCapability(tc)
	return internal.Supports16Color()
}

// New creates a new Style with default values.
//
// Default values:
//   - No colors set (uses terminal defaults)
//   - No border
//   - No padding/margin
//   - No size constraints
//   - No alignment
//   - No text decorations
//   - TrueColor terminal capability
//
// Example:
//
//	s := style.New().
//	    Foreground(style.RGB(255, 0, 0)).
//	    Bold(true).
func New() Style {
	return model.NewStyle()
}

// Render applies a Style to content and returns ANSI-styled output.
// This is the main function for styling content.
//
// The rendering pipeline:
//  1. Style validation
//  2. Size validation (if size constraints set)
//  3. Text alignment (if alignment set)
//  4. Apply padding (if padding set)
//  5. Apply border (if border set)
//  6. Apply margin (if margin set)
//  7. Color adaptation & ANSI generation
//  8. Text decorations (bold, italic, etc.)
//
// Example:
//
//	s := style.New().
//	    Foreground(style.RGB(255, 255, 255)).
//	    Background(style.RGB(0, 0, 255)).
//	    Padding(style.NewPadding(1, 2, 1, 2)).
//	    Border(style.RoundedBorder()).
//
//	output := style.Render(s, "Hello, World!")
//	fmt.Println(output)
func Render(s Style, content string) string {
	// Create services.
	colorAdapter := service2.NewColorAdapter()
	spacingCalculator := service2.NewSpacingCalculator()
	textAligner := service2.NewTextAligner()
	ansiGenerator := ansi.NewANSICodeGenerator()

	// Create render command.
	renderCmd := command.NewRenderCommand(
		colorAdapter,
		spacingCalculator,
		textAligner,
		ansiGenerator,
	)

	// Execute rendering.
	output, err := renderCmd.Execute(s, content)
	if err != nil {
		// For user-facing API, we return content as-is on error.
		// In production, you might want to log the error.
		return content
	}

	return output
}

// Color constructors.

// RGB creates a color from RGB values (0-255).
//
// Example:
//
//	red := style.RGB(255, 0, 0)
//	white := style.RGB(255, 255, 255)
func RGB(r, g, b uint8) Color {
	return value2.RGB(r, g, b)
}

// Hex creates a color from a hex string.
// Supports formats: "#RGB", "#RRGGBB", "RGB", "RRGGBB".
//
// Example:
//
//	red := style.Hex("#FF0000")
//	blue := style.Hex("0000FF")
//	shortRed := style.Hex("#F00")
func Hex(hex string) (Color, error) {
	return value2.Hex(hex)
}

// Color256 creates a color from an ANSI 256-color palette index (0-255).
//
// Palette structure:
//   - 0-15: Standard colors (black, red, green, yellow, blue, magenta, cyan, white + bright variants)
//   - 16-231: 6x6x6 RGB color cube
//   - 232-255: Grayscale ramp
//
// Example:
//
//	red := style.Color256(196)
//	gray := style.Color256(240)
func Color256(code uint8) Color {
	return value2.FromANSI256(code)
}

// Color16 creates a color from an ANSI 16-color palette index (0-15).
//
// Colors:
//   - 0-7: Normal colors (black, red, green, yellow, blue, magenta, cyan, white)
//   - 8-15: Bright variants
//
// Example:
//
//	red := style.Color16(1)
//	brightRed := style.Color16(9)
func Color16(code uint8) Color {
	// ANSI16 (0-15) is a subset of ANSI256 (0-255).
	// We can use FromANSI256 which handles 0-15 correctly.
	return value2.FromANSI256(code)
}

// Border constructors.

// Re-export border presets.
var (
	// NormalBorder is a standard single-line box-drawing border (┌─┐ │ └─┘).
	NormalBorder = value2.NormalBorder

	// RoundedBorder is a rounded corner border (╭─╮ │ ╰─╯).
	RoundedBorder = value2.RoundedBorder

	// ThickBorder is a bold/thick border (┏━┓ ┃ ┗━┛).
	ThickBorder = value2.ThickBorder

	// DoubleBorder is a double-line border (╔═╗ ║ ╚═╝).
	DoubleBorder = value2.DoubleBorder

	// ASCIIBorder is an ASCII-only border (+-+ | +-+).
	ASCIIBorder = value2.ASCIIBorder
)

// NewBorder creates a custom border with specified characters.
//
// Example:
//
//	border := style.NewBorder("*", "*", "*", "*", "*", "*", "*", "*")
//	s := style.New().Border(border)
func NewBorder(top, bottom, left, right, topLeft, topRight, bottomLeft, bottomRight string) Border {
	return value2.Border{
		Top:         top,
		Bottom:      bottom,
		Left:        left,
		Right:       right,
		TopLeft:     topLeft,
		TopRight:    topRight,
		BottomLeft:  bottomLeft,
		BottomRight: bottomRight,
	}
}

// Spacing constructors.

// NewPadding creates padding with individual values for each side.
//
// Example:
//
//	padding := style.NewPadding(1, 2, 1, 2) // top, right, bottom, left
//	s := style.New().Padding(padding)
func NewPadding(top, right, bottom, left int) Padding {
	return value2.NewPadding(top, right, bottom, left)
}

// NewMargin creates margin with individual values for each side.
//
// Example:
//
//	margin := style.NewMargin(1, 2, 1, 2) // top, right, bottom, left
//	s := style.New().Margin(margin)
func NewMargin(top, right, bottom, left int) Margin {
	return value2.NewMargin(top, right, bottom, left)
}

// Size constructors.

// NewSize creates a new Size with no constraints.
//
// Example:
//
//	size := style.NewSize().WithWidth(80).WithMaxHeight(24)
//	s := style.New().Width(80).MaxHeight(24)
func NewSize() Size {
	return value2.NewSize()
}

// Alignment constructors.

// NewAlignment creates an alignment with horizontal and vertical components.
//
// Example:
//
//	align := style.NewAlignment(style.AlignCenter, style.AlignMiddle)
//	s := style.New().Align(align)
func NewAlignment(horizontal HorizontalAlignment, vertical VerticalAlignment) Alignment {
	internalH := value2.HorizontalAlignment(horizontal)
	internalV := value2.VerticalAlignment(vertical)
	return value2.NewAlignment(internalH, internalV)
}

// Pre-defined styles.

var (
	// DefaultStyle is a style with no formatting applied.
	DefaultStyle = New()

	// BoldStyle is a style with bold text.
	BoldStyle = New().Bold(true)

	// ItalicStyle is a style with italic text.
	ItalicStyle = New().Italic(true)

	// UnderlineStyle is a style with underlined text.
	UnderlineStyle = New().Underline(true)

	// StrikethroughStyle is a style with strikethrough text.
	StrikethroughStyle = New().Strikethrough(true)
)

// Common color presets.

var (
	// Black is pure black (#000000).
	Black = RGB(0, 0, 0)

	// White is pure white (#FFFFFF).
	White = RGB(255, 255, 255)

	// Red is pure red (#FF0000).
	Red = RGB(255, 0, 0)

	// Green is pure green (#00FF00).
	Green = RGB(0, 255, 0)

	// Blue is pure blue (#0000FF).
	Blue = RGB(0, 0, 255)

	// Yellow is pure yellow (#FFFF00).
	Yellow = RGB(255, 255, 0)

	// Cyan is pure cyan (#00FFFF).
	Cyan = RGB(0, 255, 255)

	// Magenta is pure magenta (#FF00FF).
	Magenta = RGB(255, 0, 255)

	// Gray is mid gray (#808080).
	Gray = RGB(128, 128, 128)
)
