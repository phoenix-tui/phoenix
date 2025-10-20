// Package ansi provides ANSI escape code generation for terminal styling.
package ansi

import "fmt"

// ANSICodeGenerator generates low-level ANSI escape codes for terminal control.
// This is infrastructure code that knows the technical details of ANSI sequences.
// The domain layer (ColorAdapter) decides WHAT to generate, this just knows HOW.
//
//nolint:revive // ANSICodeGenerator is intentional - distinguishes from other generators
type ANSICodeGenerator struct{}

// NewANSICodeGenerator creates a new ANSICodeGenerator.
func NewANSICodeGenerator() *ANSICodeGenerator {
	return &ANSICodeGenerator{}
}

// Foreground generates an ANSI escape code for 24-bit RGB foreground color.
// Format: ESC[38;2;R;G;Bm.
func (gen *ANSICodeGenerator) Foreground(r, g, b uint8) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

// Background generates an ANSI escape code for 24-bit RGB background color.
// Format: ESC[48;2;R;G;Bm.
func (gen *ANSICodeGenerator) Background(r, g, b uint8) string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
}

// Foreground256 generates an ANSI escape code for 256-color foreground.
// Format: ESC[38;5;Nm where N is 0-255.
func (gen *ANSICodeGenerator) Foreground256(code uint8) string {
	return fmt.Sprintf("\x1b[38;5;%dm", code)
}

// Background256 generates an ANSI escape code for 256-color background.
// Format: ESC[48;5;Nm where N is 0-255.
func (gen *ANSICodeGenerator) Background256(code uint8) string {
	return fmt.Sprintf("\x1b[48;5;%dm", code)
}

// Foreground16 generates an ANSI escape code for 16-color foreground.
// Format: ESC[30-37m for normal colors, ESC[90-97m for bright colors.
func (gen *ANSICodeGenerator) Foreground16(code uint8) string {
	if code >= 8 {
		// Bright colors (90-97).
		return fmt.Sprintf("\x1b[%dm", 90+(code-8))
	}
	// Normal colors (30-37).
	return fmt.Sprintf("\x1b[%dm", 30+code)
}

// Background16 generates an ANSI escape code for 16-color background.
// Format: ESC[40-47m for normal colors, ESC[100-107m for bright colors.
func (gen *ANSICodeGenerator) Background16(code uint8) string {
	if code >= 8 {
		// Bright colors (100-107).
		return fmt.Sprintf("\x1b[%dm", 100+(code-8))
	}
	// Normal colors (40-47).
	return fmt.Sprintf("\x1b[%dm", 40+code)
}

// Reset generates an ANSI escape code to reset all attributes.
// Format: ESC[0m.
func (gen *ANSICodeGenerator) Reset() string {
	return "\x1b[0m"
}

// Bold generates an ANSI escape code for bold text.
// Format: ESC[1m.
func (gen *ANSICodeGenerator) Bold() string {
	return "\x1b[1m"
}

// Italic generates an ANSI escape code for italic text.
// Format: ESC[3m.
func (gen *ANSICodeGenerator) Italic() string {
	return "\x1b[3m"
}

// Underline generates an ANSI escape code for underlined text.
// Format: ESC[4m.
func (gen *ANSICodeGenerator) Underline() string {
	return "\x1b[4m"
}

// Strikethrough generates an ANSI escape code for strikethrough text.
// Format: ESC[9m.
func (gen *ANSICodeGenerator) Strikethrough() string {
	return "\x1b[9m"
}

// BoldOff generates an ANSI escape code to turn off bold.
// Format: ESC[22m.
func (gen *ANSICodeGenerator) BoldOff() string {
	return "\x1b[22m"
}

// ItalicOff generates an ANSI escape code to turn off italic.
// Format: ESC[23m.
func (gen *ANSICodeGenerator) ItalicOff() string {
	return "\x1b[23m"
}

// UnderlineOff generates an ANSI escape code to turn off underline.
// Format: ESC[24m.
func (gen *ANSICodeGenerator) UnderlineOff() string {
	return "\x1b[24m"
}

// StrikethroughOff generates an ANSI escape code to turn off strikethrough.
// Format: ESC[29m.
func (gen *ANSICodeGenerator) StrikethroughOff() string {
	return "\x1b[29m"
}
