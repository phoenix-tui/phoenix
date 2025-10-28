package value

import (
	"fmt"
	"strconv"
	"strings"
)

// Color represents an immutable RGB color value.
// Internal representation is always RGB (0-255 for each channel).
// This is a value object in DDD terms - immutable and defined by its values.
type Color struct {
	r, g, b uint8
}

// RGB creates a new Color from RGB values (0-255 each).
// This is the canonical constructor.
func RGB(r, g, b uint8) Color {
	return Color{r: r, g: g, b: b}
}

// Hex creates a new Color from a hex string (e.g., "#FF00FF", "FF00FF", "#F0F").
// Returns error if the hex string is invalid.
// Supports both 3-digit (#RGB) and 6-digit (#RRGGBB) formats.
func Hex(hex string) (Color, error) {
	// Remove leading '#' if present.
	hex = strings.TrimPrefix(hex, "#")

	// Validate length.
	if len(hex) != 3 && len(hex) != 6 {
		return Color{}, fmt.Errorf("invalid hex color: %q (expected 3 or 6 digits)", hex)
	}

	// Expand 3-digit format to 6-digit (#RGB -> #RRGGBB).
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}

	// Parse RGB components.
	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return Color{}, fmt.Errorf("invalid hex color: %q (bad red component)", hex)
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return Color{}, fmt.Errorf("invalid hex color: %q (bad green component)", hex)
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return Color{}, fmt.Errorf("invalid hex color: %q (bad blue component)", hex)
	}

	return Color{r: uint8(r), g: uint8(g), b: uint8(b)}, nil
}

// FromANSI256 creates a new Color from an ANSI 256-color code (0-255).
// Uses standard ANSI 256-color palette conversion to RGB.
// See: https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit.
func FromANSI256(code uint8) Color {
	// Basic 16 colors (0-15).
	if code < 16 {
		return ansi16ToRGB(code)
	}

	// 216-color cube (16-231): 6x6x6 RGB cube.
	if code >= 16 && code <= 231 {
		index := code - 16
		r := (index / 36) % 6
		g := (index / 6) % 6
		b := index % 6

		// Map 0-5 to 0-255 (with standard intensity levels).
		return Color{
			r: intensityMap[r],
			g: intensityMap[g],
			b: intensityMap[b],
		}
	}

	// Grayscale (232-255): 24 shades of gray.
	if code >= 232 {
		gray := 8 + (code-232)*10
		return Color{r: gray, g: gray, b: gray}
	}

	// Fallback (should never reach here).
	return Color{r: 0, g: 0, b: 0}
}

// RGB returns the RGB components of the color.
func (c Color) RGB() (r, g, b uint8) {
	return c.r, c.g, c.b
}

// Hex returns the color as a hex string (e.g., "#FF00FF").
func (c Color) Hex() string {
	return fmt.Sprintf("#%02X%02X%02X", c.r, c.g, c.b)
}

// Equal returns true if this color equals another color.
// Value objects are compared by value, not identity.
func (c Color) Equal(other Color) bool {
	return c.r == other.r && c.g == other.g && c.b == other.b
}

// String returns a human-readable representation of the color.
func (c Color) String() string {
	return fmt.Sprintf("Color(r=%d, g=%d, b=%d, hex=%s)", c.r, c.g, c.b, c.Hex())
}

// --- Private helpers ---.

// intensityMap maps 0-5 cube index to RGB intensity (standard ANSI 256-color mapping).
var intensityMap = [6]uint8{0, 95, 135, 175, 215, 255}

// ansi16ToRGB converts ANSI 16-color codes to RGB.
// Standard ANSI color palette (varies by terminal, but these are common defaults).
func ansi16ToRGB(code uint8) Color {
	// Standard ANSI colors (approximate RGB values).
	ansi16Colors := [16]Color{
		{0, 0, 0},       // 0: Black
		{128, 0, 0},     // 1: Red
		{0, 128, 0},     // 2: Green
		{128, 128, 0},   // 3: Yellow
		{0, 0, 128},     // 4: Blue
		{128, 0, 128},   // 5: Magenta
		{0, 128, 128},   // 6: Cyan
		{192, 192, 192}, // 7: White
		{128, 128, 128}, // 8: Bright Black (Gray)
		{255, 0, 0},     // 9: Bright Red
		{0, 255, 0},     // 10: Bright Green
		{255, 255, 0},   // 11: Bright Yellow
		{0, 0, 255},     // 12: Bright Blue
		{255, 0, 255},   // 13: Bright Magenta
		{0, 255, 255},   // 14: Bright Cyan
		{255, 255, 255}, // 15: Bright White
	}

	if code < 16 {
		return ansi16Colors[code]
	}

	// Fallback for invalid codes.
	return Color{0, 0, 0}
}
