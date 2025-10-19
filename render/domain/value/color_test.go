package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewColor(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b uint8
	}{
		{"black", 0, 0, 0},
		{"white", 255, 255, 255},
		{"red", 255, 0, 0},
		{"green", 0, 255, 0},
		{"blue", 0, 0, 255},
		{"gray", 128, 128, 128},
		{"custom", 100, 150, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := NewColor(tt.r, tt.g, tt.b)
			assert.Equal(t, tt.r, color.R())
			assert.Equal(t, tt.g, color.G())
			assert.Equal(t, tt.b, color.B())

			r, g, b := color.RGB()
			assert.Equal(t, tt.r, r)
			assert.Equal(t, tt.g, g)
			assert.Equal(t, tt.b, b)
		})
	}
}

func TestColor_Equals(t *testing.T) {
	tests := []struct {
		name     string
		c1, c2   Color
		expected bool
	}{
		{"same color", NewColor(100, 150, 200), NewColor(100, 150, 200), true},
		{"different red", NewColor(100, 150, 200), NewColor(101, 150, 200), false},
		{"different green", NewColor(100, 150, 200), NewColor(100, 151, 200), false},
		{"different blue", NewColor(100, 150, 200), NewColor(100, 150, 201), false},
		{"black vs white", ColorBlack, ColorWhite, false},
		{"same predefined", ColorRed, NewColor(255, 0, 0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.c1.Equals(tt.c2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestColor_ToANSI256(t *testing.T) {
	tests := []struct {
		name     string
		color    Color
		expected uint8
	}{
		{"black", ColorBlack, 16},    // 16 + 36*0 + 6*0 + 0
		{"white", ColorWhite, 231},   // 16 + 36*5 + 6*5 + 5
		{"red", ColorRed, 196},       // 16 + 36*5 + 6*0 + 0
		{"green", ColorGreen, 46},    // 16 + 36*0 + 6*5 + 0
		{"blue", ColorBlue, 21},      // 16 + 36*0 + 6*0 + 5
		{"mid gray", ColorGray, 102}, // ~middle of cube
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.color.ToANSI256()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestColor_String(t *testing.T) {
	tests := []struct {
		name     string
		color    Color
		expected string
	}{
		{"black", ColorBlack, "Color(0, 0, 0)"},
		{"white", ColorWhite, "Color(255, 255, 255)"},
		{"custom", NewColor(100, 150, 200), "Color(100, 150, 200)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.color.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPredefinedColors(t *testing.T) {
	// Test that predefined colors have expected RGB values
	tests := []struct {
		name    string
		color   Color
		r, g, b uint8
	}{
		{"ColorBlack", ColorBlack, 0, 0, 0},
		{"ColorRed", ColorRed, 255, 0, 0},
		{"ColorGreen", ColorGreen, 0, 255, 0},
		{"ColorYellow", ColorYellow, 255, 255, 0},
		{"ColorBlue", ColorBlue, 0, 0, 255},
		{"ColorMagenta", ColorMagenta, 255, 0, 255},
		{"ColorCyan", ColorCyan, 0, 255, 255},
		{"ColorWhite", ColorWhite, 255, 255, 255},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.r, tt.color.R())
			assert.Equal(t, tt.g, tt.color.G())
			assert.Equal(t, tt.b, tt.color.B())
		})
	}
}
