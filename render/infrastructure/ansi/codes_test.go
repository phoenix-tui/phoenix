package ansi

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	// Verify ANSI escape sequence structure
	assert.True(t, strings.HasPrefix(Bold, CSI))
	assert.True(t, strings.HasPrefix(ClearScreen, CSI))
	assert.True(t, strings.HasPrefix(CursorHide, CSI))
}

func TestMoveCursor(t *testing.T) {
	tests := []struct {
		name string
		x, y int
		want string
	}{
		{"origin", 0, 0, "\x1b[1;1H"},
		{"middle", 40, 12, "\x1b[13;41H"},
		{"corner", 79, 23, "\x1b[24;80H"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoveCursor(tt.x, tt.y)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestMoveCursorUp(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{"zero", 0, ""},
		{"negative", -1, ""},
		{"one", 1, "\x1b[1A"},
		{"ten", 10, "\x1b[10A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoveCursorUp(tt.n)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestMoveCursorDown(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{"zero", 0, ""},
		{"negative", -1, ""},
		{"one", 1, "\x1b[1B"},
		{"ten", 10, "\x1b[10B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoveCursorDown(tt.n)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestMoveCursorRight(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{"zero", 0, ""},
		{"negative", -1, ""},
		{"one", 1, "\x1b[1C"},
		{"ten", 10, "\x1b[10C"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoveCursorRight(tt.n)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestMoveCursorLeft(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{"zero", 0, ""},
		{"negative", -1, ""},
		{"one", 1, "\x1b[1D"},
		{"ten", 10, "\x1b[10D"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoveCursorLeft(tt.n)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestSetFg256(t *testing.T) {
	tests := []struct {
		name  string
		color uint8
		want  string
	}{
		{"black", 0, "\x1b[38;5;0m"},
		{"red", 196, "\x1b[38;5;196m"},
		{"white", 255, "\x1b[38;5;255m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SetFg256(tt.color)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestSetBg256(t *testing.T) {
	tests := []struct {
		name  string
		color uint8
		want  string
	}{
		{"black", 0, "\x1b[48;5;0m"},
		{"red", 196, "\x1b[48;5;196m"},
		{"white", 255, "\x1b[48;5;255m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SetBg256(tt.color)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestSetFgRGB(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b uint8
		want    string
	}{
		{"black", 0, 0, 0, "\x1b[38;2;0;0;0m"},
		{"red", 255, 0, 0, "\x1b[38;2;255;0;0m"},
		{"green", 0, 255, 0, "\x1b[38;2;0;255;0m"},
		{"blue", 0, 0, 255, "\x1b[38;2;0;0;255m"},
		{"white", 255, 255, 255, "\x1b[38;2;255;255;255m"},
		{"custom", 100, 150, 200, "\x1b[38;2;100;150;200m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SetFgRGB(tt.r, tt.g, tt.b)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestSetBgRGB(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b uint8
		want    string
	}{
		{"black", 0, 0, 0, "\x1b[48;2;0;0;0m"},
		{"red", 255, 0, 0, "\x1b[48;2;255;0;0m"},
		{"white", 255, 255, 255, "\x1b[48;2;255;255;255m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SetBgRGB(tt.r, tt.g, tt.b)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestResetFg(t *testing.T) {
	result := ResetFg()
	assert.Equal(t, "\x1b[39m", result)
}

func TestResetBg(t *testing.T) {
	result := ResetBg()
	assert.Equal(t, "\x1b[49m", result)
}

func TestSetCursorShape(t *testing.T) {
	tests := []struct {
		name  string
		shape int
		want  string
	}{
		{"default", 0, "\x1b[0 q"},
		{"blinking block", 1, "\x1b[1 q"},
		{"steady block", 2, "\x1b[2 q"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SetCursorShape(tt.shape)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestSetScrollRegion(t *testing.T) {
	result := SetScrollRegion(5, 20)
	assert.Equal(t, "\x1b[5;20r", result)
}

func TestSaveCursorPosition(t *testing.T) {
	result := SaveCursorPosition()
	assert.Equal(t, "\x1b7", result)
}

func TestRestoreCursorPosition(t *testing.T) {
	result := RestoreCursorPosition()
	assert.Equal(t, "\x1b8", result)
}

func TestSetTitle(t *testing.T) {
	result := SetTitle("My Terminal")
	assert.Contains(t, result, "My Terminal")
	assert.True(t, strings.HasPrefix(result, ESC))
}

func TestBell(t *testing.T) {
	result := Bell()
	assert.Equal(t, "\x07", result)
}

// Benchmark ANSI code generation
func BenchmarkMoveCursor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MoveCursor(40, 12)
	}
}

func BenchmarkSetFgRGB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SetFgRGB(255, 128, 64)
	}
}

func BenchmarkSetFg256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SetFg256(196)
	}
}
