package service

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/style/domain/value"
)

// TestNewColorAdapter tests the constructor.
func TestNewColorAdapter(t *testing.T) {
	adapter := NewColorAdapter()
	if adapter == nil {
		t.Error("NewColorAdapter() returned nil")
	}
}

// TestToANSIForeground_TrueColor tests TrueColor foreground conversion.
func TestToANSIForeground_TrueColor(t *testing.T) {
	adapter := NewColorAdapter()

	tests := []struct {
		name  string
		color value.Color
		want  string
	}{
		{"Red", value.RGB(255, 0, 0), "\x1b[38;2;255;0;0m"},
		{"Green", value.RGB(0, 255, 0), "\x1b[38;2;0;255;0m"},
		{"Blue", value.RGB(0, 0, 255), "\x1b[38;2;0;0;255m"},
		{"White", value.RGB(255, 255, 255), "\x1b[38;2;255;255;255m"},
		{"Black", value.RGB(0, 0, 0), "\x1b[38;2;0;0;0m"},
		{"Custom", value.RGB(123, 45, 67), "\x1b[38;2;123;45;67m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adapter.ToANSIForeground(tt.color, value.TrueColor)
			if got != tt.want {
				t.Errorf("ToANSIForeground() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestToANSIBackground_TrueColor tests TrueColor background conversion.
func TestToANSIBackground_TrueColor(t *testing.T) {
	adapter := NewColorAdapter()

	tests := []struct {
		name  string
		color value.Color
		want  string
	}{
		{"Red", value.RGB(255, 0, 0), "\x1b[48;2;255;0;0m"},
		{"Green", value.RGB(0, 255, 0), "\x1b[48;2;0;255;0m"},
		{"Blue", value.RGB(0, 0, 255), "\x1b[48;2;0;0;255m"},
		{"White", value.RGB(255, 255, 255), "\x1b[48;2;255;255;255m"},
		{"Black", value.RGB(0, 0, 0), "\x1b[48;2;0;0;0m"},
		{"Custom", value.RGB(123, 45, 67), "\x1b[48;2;123;45;67m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adapter.ToANSIBackground(tt.color, value.TrueColor)
			if got != tt.want {
				t.Errorf("ToANSIBackground() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestToANSIForeground_ANSI256 tests ANSI256 foreground conversion.
func TestToANSIForeground_ANSI256(t *testing.T) {
	adapter := NewColorAdapter()

	tests := []struct {
		name  string
		color value.Color
	}{
		{"Red", value.RGB(255, 0, 0)},
		{"Green", value.RGB(0, 255, 0)},
		{"Blue", value.RGB(0, 0, 255)},
		{"White", value.RGB(255, 255, 255)},
		{"Black", value.RGB(0, 0, 0)},
		{"Gray", value.RGB(128, 128, 128)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adapter.ToANSIForeground(tt.color, value.ANSI256)

			// Should start with "\x1b[38;5;" and end with "m"
			if !strings.HasPrefix(got, "\x1b[38;5;") {
				t.Errorf("ToANSIForeground() = %q, want prefix %q", got, "\x1b[38;5;")
			}
			if !strings.HasSuffix(got, "m") {
				t.Errorf("ToANSIForeground() = %q, want suffix %q", got, "m")
			}

			// Extract color code
			code := got[7 : len(got)-1]
			if code == "" {
				t.Errorf("ToANSIForeground() has empty color code: %q", got)
			}
		})
	}
}

// TestToANSIBackground_ANSI256 tests ANSI256 background conversion.
func TestToANSIBackground_ANSI256(t *testing.T) {
	adapter := NewColorAdapter()

	tests := []struct {
		name  string
		color value.Color
	}{
		{"Red", value.RGB(255, 0, 0)},
		{"Green", value.RGB(0, 255, 0)},
		{"Blue", value.RGB(0, 0, 255)},
		{"White", value.RGB(255, 255, 255)},
		{"Black", value.RGB(0, 0, 0)},
		{"Gray", value.RGB(128, 128, 128)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adapter.ToANSIBackground(tt.color, value.ANSI256)

			// Should start with "\x1b[48;5;" and end with "m"
			if !strings.HasPrefix(got, "\x1b[48;5;") {
				t.Errorf("ToANSIBackground() = %q, want prefix %q", got, "\x1b[48;5;")
			}
			if !strings.HasSuffix(got, "m") {
				t.Errorf("ToANSIBackground() = %q, want suffix %q", got, "m")
			}

			// Extract color code
			code := got[7 : len(got)-1]
			if code == "" {
				t.Errorf("ToANSIBackground() has empty color code: %q", got)
			}
		})
	}
}

// TestToANSIForeground_ANSI16 tests ANSI16 foreground conversion.
func TestToANSIForeground_ANSI16(t *testing.T) {
	adapter := NewColorAdapter()

	tests := []struct {
		name  string
		color value.Color
	}{
		{"Red", value.RGB(255, 0, 0)},
		{"Green", value.RGB(0, 255, 0)},
		{"Blue", value.RGB(0, 0, 255)},
		{"White", value.RGB(255, 255, 255)},
		{"Black", value.RGB(0, 0, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adapter.ToANSIForeground(tt.color, value.ANSI16)

			// Should start with "\x1b[" and end with "m"
			if !strings.HasPrefix(got, "\x1b[") {
				t.Errorf("ToANSIForeground() = %q, want prefix %q", got, "\x1b[")
			}
			if !strings.HasSuffix(got, "m") {
				t.Errorf("ToANSIForeground() = %q, want suffix %q", got, "m")
			}

			// Should be either 30-37 (normal) or 90-97 (bright)
			code := got[2 : len(got)-1]
			if code == "" {
				t.Errorf("ToANSIForeground() has empty color code: %q", got)
			}
		})
	}
}

// TestToANSIBackground_ANSI16 tests ANSI16 background conversion.
func TestToANSIBackground_ANSI16(t *testing.T) {
	adapter := NewColorAdapter()

	tests := []struct {
		name  string
		color value.Color
	}{
		{"Red", value.RGB(255, 0, 0)},
		{"Green", value.RGB(0, 255, 0)},
		{"Blue", value.RGB(0, 0, 255)},
		{"White", value.RGB(255, 255, 255)},
		{"Black", value.RGB(0, 0, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adapter.ToANSIBackground(tt.color, value.ANSI16)

			// Should start with "\x1b[" and end with "m"
			if !strings.HasPrefix(got, "\x1b[") {
				t.Errorf("ToANSIBackground() = %q, want prefix %q", got, "\x1b[")
			}
			if !strings.HasSuffix(got, "m") {
				t.Errorf("ToANSIBackground() = %q, want suffix %q", got, "m")
			}

			// Should be either 40-47 (normal) or 100-107 (bright)
			code := got[2 : len(got)-1]
			if code == "" {
				t.Errorf("ToANSIBackground() has empty color code: %q", got)
			}
		})
	}
}

// TestToANSI_NoColor tests NoColor mode (should return empty string).
func TestToANSI_NoColor(t *testing.T) {
	adapter := NewColorAdapter()
	color := value.RGB(255, 0, 0)

	fg := adapter.ToANSIForeground(color, value.NoColor)
	if fg != "" {
		t.Errorf("ToANSIForeground(NoColor) = %q, want %q", fg, "")
	}

	bg := adapter.ToANSIBackground(color, value.NoColor)
	if bg != "" {
		t.Errorf("ToANSIBackground(NoColor) = %q, want %q", bg, "")
	}
}

// TestRGBToANSI256_Grayscale tests grayscale conversion.
func TestRGBToANSI256_Grayscale(t *testing.T) {
	adapter := &DefaultColorAdapter{}

	tests := []struct {
		name      string
		r, g, b   uint8
		wantRange string // "cube" or "gray"
	}{
		{"Pure black", 0, 0, 0, "cube"},
		{"Near black", 5, 5, 5, "cube"}, // Too close to black, uses cube
		{"Dark gray", 50, 50, 50, "gray"},
		{"Mid gray", 128, 128, 128, "gray"},
		{"Light gray", 200, 200, 200, "gray"},
		{"Near white", 240, 240, 240, "cube"}, // Too close to white, uses cube
		{"Pure white", 255, 255, 255, "cube"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := value.RGB(tt.r, tt.g, tt.b)
			code := adapter.rgbToANSI256(color)

			// Verify code is in expected range
			if tt.wantRange == "cube" {
				if code < 16 || code > 231 {
					t.Errorf("rgbToANSI256(%d,%d,%d) = %d, want cube range (16-231)",
						tt.r, tt.g, tt.b, code)
				}
			} else if tt.wantRange == "gray" {
				if code < 232 || code > 255 {
					t.Errorf("rgbToANSI256(%d,%d,%d) = %d, want gray range (232-255)",
						tt.r, tt.g, tt.b, code)
				}
			}
		})
	}
}

// TestRGBToANSI256_ColorCube tests color cube conversion.
func TestRGBToANSI256_ColorCube(t *testing.T) {
	adapter := &DefaultColorAdapter{}

	tests := []struct {
		name    string
		r, g, b uint8
	}{
		{"Red", 255, 0, 0},
		{"Green", 0, 255, 0},
		{"Blue", 0, 0, 255},
		{"Yellow", 255, 255, 0},
		{"Cyan", 0, 255, 255},
		{"Magenta", 255, 0, 255},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := value.RGB(tt.r, tt.g, tt.b)
			code := adapter.rgbToANSI256(color)

			// Should be in color cube range (16-231)
			if code < 16 || code > 231 {
				t.Errorf("rgbToANSI256(%d,%d,%d) = %d, want range 16-231",
					tt.r, tt.g, tt.b, code)
			}
		})
	}
}

// TestRGBToANSI16_BasicColors tests basic 16-color conversion.
func TestRGBToANSI16_BasicColors(t *testing.T) {
	adapter := &DefaultColorAdapter{}

	tests := []struct {
		name     string
		r, g, b  uint8
		wantCode uint8
	}{
		{"Black", 0, 0, 0, 0},
		{"Bright white", 255, 255, 255, 15},
		{"Red", 255, 0, 0, 9},        // Bright red
		{"Green", 0, 255, 0, 10},     // Bright green
		{"Blue", 0, 0, 255, 12},      // Bright blue
		{"Yellow", 255, 255, 0, 11},  // Bright yellow
		{"Cyan", 0, 255, 255, 14},    // Bright cyan
		{"Magenta", 255, 0, 255, 13}, // Bright magenta
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := value.RGB(tt.r, tt.g, tt.b)
			code := adapter.rgbToANSI16(color)

			if code != tt.wantCode {
				t.Errorf("rgbToANSI16(%d,%d,%d) = %d, want %d",
					tt.r, tt.g, tt.b, code, tt.wantCode)
			}
		})
	}
}

// TestColorAdaptation_SameColorDifferentCapabilities tests that the same color
// produces different ANSI codes for different terminal capabilities.
func TestColorAdaptation_SameColorDifferentCapabilities(t *testing.T) {
	adapter := NewColorAdapter()
	color := value.RGB(123, 45, 67)

	trueColor := adapter.ToANSIForeground(color, value.TrueColor)
	ansi256 := adapter.ToANSIForeground(color, value.ANSI256)
	ansi16 := adapter.ToANSIForeground(color, value.ANSI16)
	noColor := adapter.ToANSIForeground(color, value.NoColor)

	// All should be different (except noColor which is empty)
	if trueColor == ansi256 {
		t.Error("TrueColor and ANSI256 should produce different codes")
	}
	if trueColor == ansi16 {
		t.Error("TrueColor and ANSI16 should produce different codes")
	}
	if ansi256 == ansi16 {
		t.Error("ANSI256 and ANSI16 should produce different codes")
	}
	if noColor != "" {
		t.Errorf("NoColor should be empty, got %q", noColor)
	}

	// All color modes should produce non-empty codes
	if trueColor == "" {
		t.Error("TrueColor should not be empty")
	}
	if ansi256 == "" {
		t.Error("ANSI256 should not be empty")
	}
	if ansi16 == "" {
		t.Error("ANSI16 should not be empty")
	}
}
