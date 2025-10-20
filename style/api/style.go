// Package style provides the public API for the Phoenix TUI Framework styling library.
// This package exports a simplified interface for end users, hiding DDD layer complexity.
//
// Quick Start:.
//
//	import "github.com/phoenix-tui/phoenix/style/api".
//
//	// Create a style.
//	s := style.New().
//	    Foreground(style.RGB(255, 255, 255)).
//	    Background(style.RGB(0, 0, 255)).
//	    Bold(true).
//
//	// Render styled content.
//	output := style.Render(s, "Hello, World!").
//	fmt.Println(output).
//
// Features:.
//   - Rich styling API (colors, borders, padding, margin, alignment).
//   - Unicode-correct rendering (fixes Lipgloss #562).
//   - Fluent API with method chaining.
//   - Immutable styles (thread-safe).
//   - Terminal capability adaptation (TrueColor → ANSI256 → ANSI16).
package style

import (
	coreService "github.com/phoenix-tui/phoenix/core/domain/service"
	"github.com/phoenix-tui/phoenix/style/application/command"
	"github.com/phoenix-tui/phoenix/style/domain/model"
	"github.com/phoenix-tui/phoenix/style/domain/service"
	"github.com/phoenix-tui/phoenix/style/domain/value"
	"github.com/phoenix-tui/phoenix/style/infrastructure/ansi"
)

// Style is an alias for model.Style, the main styling configuration.
type Style = model.Style

// Color is an alias for value.Color, representing terminal colors.
type Color = value.Color

// Aliases for value types.
type (
	// Border is an alias for value.Border, representing box borders.
	Border = value.Border
	// Padding is an alias for value.Padding, representing box padding.
	Padding = value.Padding
	// Margin is an alias for value.Margin, representing box margins.
	Margin = value.Margin
	// Size is an alias for value.Size, representing box dimensions.
	Size = value.Size
	// Alignment is an alias for value.Alignment, representing text alignment.
	Alignment = value.Alignment
	// HorizontalAlignment is an alias for value.HorizontalAlignment.
	HorizontalAlignment = value.HorizontalAlignment
	// VerticalAlignment is an alias for value.VerticalAlignment.
	VerticalAlignment = value.VerticalAlignment
	// TerminalCapability is an alias for value.TerminalCapability.
	TerminalCapability = value.TerminalCapability
)

// Re-export terminal capabilities.
const (
	TrueColor TerminalCapability = value.TrueColor
	ANSI256   TerminalCapability = value.ANSI256
	ANSI16    TerminalCapability = value.ANSI16
)

// Re-export horizontal alignments.
const (
	AlignLeft   HorizontalAlignment = value.AlignLeft
	AlignCenter HorizontalAlignment = value.AlignCenter
	AlignRight  HorizontalAlignment = value.AlignRight
)

// Re-export vertical alignments.
const (
	AlignTop    VerticalAlignment = value.AlignTop
	AlignMiddle VerticalAlignment = value.AlignMiddle
	AlignBottom VerticalAlignment = value.AlignBottom
)

// New creates a new Style with default values.
//
// Default values:.
//   - No colors set (uses terminal defaults).
//   - No border.
//   - No padding/margin.
//   - No size constraints.
//   - No alignment.
//   - No text decorations.
//   - TrueColor terminal capability.
//
// Example:.
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
// The rendering pipeline:.
//  1. Style validation.
//  2. Size validation (if size constraints set).
//  3. Text alignment (if alignment set).
//  4. Apply padding (if padding set).
//  5. Apply border (if border set).
//  6. Apply margin (if margin set).
//  7. Color adaptation & ANSI generation.
//  8. Text decorations (bold, italic, etc.).
//
// Example:.
//
//	s := style.New().
//	    Foreground(style.RGB(255, 255, 255)).
//	    Background(style.RGB(0, 0, 255)).
//	    Padding(style.NewPadding(1, 2, 1, 2)).
//	    Border(style.RoundedBorder()).
//
//	output := style.Render(s, "Hello, World!").
//	fmt.Println(output).
func Render(s Style, content string) string {
	// Create services.
	unicodeService := coreService.NewUnicodeService()
	colorAdapter := service.NewColorAdapter()
	spacingCalculator := service.NewSpacingCalculator(unicodeService)
	textAligner := service.NewTextAligner(unicodeService)
	ansiGenerator := ansi.NewANSICodeGenerator()

	// Create render command.
	renderCmd := command.NewRenderCommand(
		colorAdapter,
		spacingCalculator,
		textAligner,
		unicodeService,
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
// Example:.
//
//	red := style.RGB(255, 0, 0).
//	white := style.RGB(255, 255, 255).
func RGB(r, g, b uint8) Color {
	return value.RGB(r, g, b)
}

// Hex creates a color from a hex string.
// Supports formats: "#RGB", "#RRGGBB", "RGB", "RRGGBB".
//
// Example:.
//
//	red := style.Hex("#FF0000").
//	blue := style.Hex("0000FF").
//	shortRed := style.Hex("#F00").
func Hex(hex string) (Color, error) {
	return value.Hex(hex)
}

// Color256 creates a color from an ANSI 256-color palette index (0-255).
//
// Palette structure:.
//   - 0-15: Standard colors (black, red, green, yellow, blue, magenta, cyan, white + bright variants).
//   - 16-231: 6x6x6 RGB color cube.
//   - 232-255: Grayscale ramp.
//
// Example:.
//
//	red := style.Color256(196).
//	gray := style.Color256(240).
func Color256(code uint8) Color {
	return value.FromANSI256(code)
}

// Color16 creates a color from an ANSI 16-color palette index (0-15).
//
// Colors:.
//   - 0-7: Normal colors (black, red, green, yellow, blue, magenta, cyan, white).
//   - 8-15: Bright variants.
//
// Example:.
//
//	red := style.Color16(1).
//	brightRed := style.Color16(9).
func Color16(code uint8) Color {
	// ANSI16 (0-15) is a subset of ANSI256 (0-255).
	// We can use FromANSI256 which handles 0-15 correctly.
	return value.FromANSI256(code)
}

// Border constructors.

// Re-export border presets.
var (
	// NormalBorder is a standard single-line box-drawing border (┌─┐ │ └─┘).
	NormalBorder = value.NormalBorder

	// RoundedBorder is a rounded corner border (╭─╮ │ ╰─╯).
	RoundedBorder = value.RoundedBorder

	// ThickBorder is a bold/thick border (┏━┓ ┃ ┗━┛).
	ThickBorder = value.ThickBorder

	// DoubleBorder is a double-line border (╔═╗ ║ ╚═╝).
	DoubleBorder = value.DoubleBorder

	// ASCIIBorder is an ASCII-only border (+-+ | +-+).
	ASCIIBorder = value.ASCIIBorder
)

// NewBorder creates a custom border with specified characters.
//
// Example:.
//
//	border := style.NewBorder("*", "*", "*", "*", "*", "*", "*", "*").
//	s := style.New().Border(border).
func NewBorder(top, bottom, left, right, topLeft, topRight, bottomLeft, bottomRight string) Border {
	return value.Border{
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
// Example:.
//
//	padding := style.NewPadding(1, 2, 1, 2) // top, right, bottom, left.
//	s := style.New().Padding(padding).
func NewPadding(top, right, bottom, left int) Padding {
	return value.NewPadding(top, right, bottom, left)
}

// NewMargin creates margin with individual values for each side.
//
// Example:.
//
//	margin := style.NewMargin(1, 2, 1, 2) // top, right, bottom, left.
//	s := style.New().Margin(margin).
func NewMargin(top, right, bottom, left int) Margin {
	return value.NewMargin(top, right, bottom, left)
}

// Size constructors.

// NewSize creates a new Size with no constraints.
//
// Example:.
//
//	size := style.NewSize().WithWidth(80).WithMaxHeight(24).
//	s := style.New().Width(80).MaxHeight(24).
func NewSize() Size {
	return value.NewSize()
}

// Alignment constructors.

// NewAlignment creates an alignment with horizontal and vertical components.
//
// Example:.
//
//	align := style.NewAlignment(style.AlignCenter, style.AlignMiddle).
//	s := style.New().Align(align).
func NewAlignment(horizontal HorizontalAlignment, vertical VerticalAlignment) Alignment {
	return value.NewAlignment(horizontal, vertical)
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
