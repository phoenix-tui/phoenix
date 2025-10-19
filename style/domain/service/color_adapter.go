package service

import (
	"fmt"
	"math"

	"github.com/phoenix-tui/phoenix/style/domain/value"
)

// ColorAdapter is a domain service that adapts Color value objects
// to ANSI escape codes based on terminal capabilities.
// This is pure business logic with no infrastructure dependencies.
type ColorAdapter interface {
	// ToANSIForeground converts a color to an ANSI foreground code.
	ToANSIForeground(color value.Color, capability value.TerminalCapability) string

	// ToANSIBackground converts a color to an ANSI background code.
	ToANSIBackground(color value.Color, capability value.TerminalCapability) string
}

// DefaultColorAdapter is the default implementation of ColorAdapter.
type DefaultColorAdapter struct{}

// NewColorAdapter creates a new DefaultColorAdapter.
func NewColorAdapter() ColorAdapter {
	return &DefaultColorAdapter{}
}

// ToANSIForeground converts a color to an ANSI foreground code based on terminal capability.
func (a *DefaultColorAdapter) ToANSIForeground(color value.Color, capability value.TerminalCapability) string {
	switch capability {
	case value.NoColor:
		return ""
	case value.ANSI16:
		return a.toANSI16Foreground(color)
	case value.ANSI256:
		return a.toANSI256Foreground(color)
	case value.TrueColor:
		return a.toTrueColorForeground(color)
	default:
		return ""
	}
}

// ToANSIBackground converts a color to an ANSI background code based on terminal capability.
func (a *DefaultColorAdapter) ToANSIBackground(color value.Color, capability value.TerminalCapability) string {
	switch capability {
	case value.NoColor:
		return ""
	case value.ANSI16:
		return a.toANSI16Background(color)
	case value.ANSI256:
		return a.toANSI256Background(color)
	case value.TrueColor:
		return a.toTrueColorBackground(color)
	default:
		return ""
	}
}

// --- TrueColor (24-bit RGB) ---

func (a *DefaultColorAdapter) toTrueColorForeground(color value.Color) string {
	r, g, b := color.RGB()
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

func (a *DefaultColorAdapter) toTrueColorBackground(color value.Color) string {
	r, g, b := color.RGB()
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
}

// --- ANSI 256 colors ---

func (a *DefaultColorAdapter) toANSI256Foreground(color value.Color) string {
	code := a.rgbToANSI256(color)
	return fmt.Sprintf("\x1b[38;5;%dm", code)
}

func (a *DefaultColorAdapter) toANSI256Background(color value.Color) string {
	code := a.rgbToANSI256(color)
	return fmt.Sprintf("\x1b[48;5;%dm", code)
}

// rgbToANSI256 converts RGB color to closest ANSI 256-color code.
// Uses Euclidean distance in RGB space to find the closest match.
func (a *DefaultColorAdapter) rgbToANSI256(color value.Color) uint8 {
	r, g, b := color.RGB()

	// Check if it's close to a grayscale value (232-255)
	if isGrayish(r, g, b) {
		// Map to grayscale range (232-255): 24 shades from near-black to near-white
		// gray = 8 + (code-232)*10, so code = (gray-8)/10 + 232
		gray := (uint16(r) + uint16(g) + uint16(b)) / 3
		if gray < 8 {
			return 16 // Use black from color cube
		}
		if gray > 238 {
			return 231 // Use white from color cube
		}
		code := (gray-8)/10 + 232
		return uint8(code)
	}

	// Use 6x6x6 RGB color cube (16-231)
	// Map RGB (0-255) to cube index (0-5)
	rCube := rgbToCubeIndex(r)
	gCube := rgbToCubeIndex(g)
	bCube := rgbToCubeIndex(b)

	// Calculate ANSI 256 code: 16 + 36*r + 6*g + b
	code := 16 + 36*rCube + 6*gCube + bCube
	return uint8(code)
}

// isGrayish returns true if the color is close to grayscale (low color variance).
func isGrayish(r, g, b uint8) bool {
	avg := (uint16(r) + uint16(g) + uint16(b)) / 3
	tolerance := uint16(10) // Allow small variance

	return abs(uint16(r), avg) <= tolerance &&
		abs(uint16(g), avg) <= tolerance &&
		abs(uint16(b), avg) <= tolerance
}

// rgbToCubeIndex maps RGB value (0-255) to 6x6x6 cube index (0-5).
// Uses closest match to standard intensity levels: [0, 95, 135, 175, 215, 255]
func rgbToCubeIndex(value uint8) uint8 {
	intensities := []uint8{0, 95, 135, 175, 215, 255}
	closestIndex := uint8(0)
	closestDistance := uint8(255)

	for i, intensity := range intensities {
		distance := absDiff(value, intensity)
		if distance < closestDistance {
			closestDistance = distance
			closestIndex = uint8(i)
		}
	}

	return closestIndex
}

// --- ANSI 16 colors ---

func (a *DefaultColorAdapter) toANSI16Foreground(color value.Color) string {
	code := a.rgbToANSI16(color)
	if code >= 8 {
		// Bright colors (90-97)
		return fmt.Sprintf("\x1b[%dm", 90+(code-8))
	}
	// Normal colors (30-37)
	return fmt.Sprintf("\x1b[%dm", 30+code)
}

func (a *DefaultColorAdapter) toANSI16Background(color value.Color) string {
	code := a.rgbToANSI16(color)
	if code >= 8 {
		// Bright colors (100-107)
		return fmt.Sprintf("\x1b[%dm", 100+(code-8))
	}
	// Normal colors (40-47)
	return fmt.Sprintf("\x1b[%dm", 40+code)
}

// rgbToANSI16 converts RGB color to closest ANSI 16-color code (0-15).
// Uses Euclidean distance in RGB space.
func (a *DefaultColorAdapter) rgbToANSI16(color value.Color) uint8 {
	// Standard ANSI 16 colors (same as in color.go)
	ansi16Colors := [16][3]uint8{
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

	r, g, b := color.RGB()
	closestIndex := uint8(0)
	closestDistance := math.MaxFloat64

	for i, ansiColor := range ansi16Colors {
		distance := colorDistance(r, g, b, ansiColor[0], ansiColor[1], ansiColor[2])
		if distance < closestDistance {
			closestDistance = distance
			closestIndex = uint8(i)
		}
	}

	return closestIndex
}

// --- Helpers ---

// colorDistance calculates Euclidean distance between two RGB colors.
func colorDistance(r1, g1, b1, r2, g2, b2 uint8) float64 {
	dr := float64(r1) - float64(r2)
	dg := float64(g1) - float64(g2)
	db := float64(b1) - float64(b2)
	return math.Sqrt(dr*dr + dg*dg + db*db)
}

// abs returns the absolute difference between two uint16 values.
func abs(a, b uint16) uint16 {
	if a > b {
		return a - b
	}
	return b - a
}

// absDiff returns the absolute difference between two uint8 values.
func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}
