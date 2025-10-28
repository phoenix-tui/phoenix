// Package value defines value objects for rendering (Color, Style).
package value

import "fmt"

// Color represents a terminal color in RGB format.
// Color is immutable and uses value semantics.
type Color struct {
	r, g, b uint8
}

// NewColor creates a new RGB color.
func NewColor(r, g, b uint8) Color {
	return Color{r: r, g: g, b: b}
}

// RGB returns the RGB components.
func (c Color) RGB() (r, g, b uint8) {
	return c.r, c.g, c.b
}

// R returns the red component.
func (c Color) R() uint8 {
	return c.r
}

// G returns the green component.
func (c Color) G() uint8 {
	return c.g
}

// B returns the blue component.
func (c Color) B() uint8 {
	return c.b
}

// Equals checks if two colors are equal.
func (c Color) Equals(other Color) bool {
	return c.r == other.r && c.g == other.g && c.b == other.b
}

// ToANSI256 converts RGB to ANSI 256 color index (approximate).
// Uses 6x6x6 color cube (colors 16-231).
func (c Color) ToANSI256() uint8 {
	// Convert to 6-level scale.
	r := uint8(float64(c.r) / 255.0 * 5.0)
	g := uint8(float64(c.g) / 255.0 * 5.0)
	b := uint8(float64(c.b) / 255.0 * 5.0)

	// ANSI 256 color cube formula: 16 + 36*r + 6*g + b.
	return 16 + 36*r + 6*g + b
}

// String returns a string representation for debugging.
func (c Color) String() string {
	return fmt.Sprintf("Color(%d, %d, %d)", c.r, c.g, c.b)
}

// Common colors (predefined constants).
var (
	ColorBlack   = NewColor(0, 0, 0)
	ColorRed     = NewColor(255, 0, 0)
	ColorGreen   = NewColor(0, 255, 0)
	ColorYellow  = NewColor(255, 255, 0)
	ColorBlue    = NewColor(0, 0, 255)
	ColorMagenta = NewColor(255, 0, 255)
	ColorCyan    = NewColor(0, 255, 255)
	ColorWhite   = NewColor(255, 255, 255)

	// Gray scale.
	ColorGray        = NewColor(128, 128, 128)
	ColorDarkGray    = NewColor(64, 64, 64)
	ColorLightGray   = NewColor(192, 192, 192)
	ColorBrightWhite = NewColor(255, 255, 255)
)
