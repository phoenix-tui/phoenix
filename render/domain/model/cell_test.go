package model

import (
	"testing"

	"github.com/phoenix-tui/phoenix/render/domain/value"
	"github.com/stretchr/testify/assert"
)

func TestNewCell(t *testing.T) {
	tests := []struct {
		name          string
		char          rune
		style         value.Style
		expectedWidth int
	}{
		{"ascii", 'A', value.NewStyle(), 1},
		{"space", ' ', value.NewStyle(), 1},
		{"digit", '5', value.NewStyle(), 1},
		{"punctuation", '!', value.NewStyle(), 1},
		{"emoji", '😀', value.NewStyle(), 2},
		{"cjk", '中', value.NewStyle(), 2},
		{"zero width", '\u200B', value.NewStyle(), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := NewCell(tt.char, tt.style)
			assert.Equal(t, tt.char, cell.Char())
			assert.True(t, tt.style.Equals(cell.Style()))
			assert.Equal(t, tt.expectedWidth, cell.Width())
		})
	}
}

func TestNewEmptyCell(t *testing.T) {
	cell := NewEmptyCell()
	assert.Equal(t, ' ', cell.Char())
	assert.True(t, cell.Style().IsEmpty())
	assert.Equal(t, 1, cell.Width())
	assert.True(t, cell.IsEmpty())
}

func TestNewCellWithWidth(t *testing.T) {
	style := value.NewStyleWithFg(value.ColorRed)
	cell := NewCellWithWidth('A', style, 1)

	assert.Equal(t, 'A', cell.Char())
	assert.True(t, style.Equals(cell.Style()))
	assert.Equal(t, 1, cell.Width())
}

func TestCell_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		cell     Cell
		expected bool
	}{
		{"empty cell", NewEmptyCell(), true},
		{"space with no style", NewCell(' ', value.NewStyle()), true},
		{"space with style", NewCell(' ', value.NewStyleWithFg(value.ColorRed)), false},
		{"char with no style", NewCell('A', value.NewStyle()), false},
		{"char with style", NewCell('A', value.NewStyleWithFg(value.ColorRed)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cell.IsEmpty()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCell_Equals(t *testing.T) {
	style1 := value.NewStyleWithFg(value.ColorRed)
	style2 := value.NewStyleWithFg(value.ColorBlue)

	tests := []struct {
		name     string
		c1, c2   Cell
		expected bool
	}{
		{
			"same cell",
			NewCell('A', value.NewStyle()),
			NewCell('A', value.NewStyle()),
			true,
		},
		{
			"different char",
			NewCell('A', value.NewStyle()),
			NewCell('B', value.NewStyle()),
			false,
		},
		{
			"different style",
			NewCell('A', style1),
			NewCell('A', style2),
			false,
		},
		{
			"same styled cell",
			NewCell('A', style1),
			NewCell('A', style1),
			true,
		},
		{
			"empty cells",
			NewEmptyCell(),
			NewEmptyCell(),
			true,
		},
		{
			"different width chars",
			NewCell('A', value.NewStyle()),
			NewCell('中', value.NewStyle()),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.c1.Equals(tt.c2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCell_WithChar(t *testing.T) {
	original := NewCell('A', value.NewStyleWithFg(value.ColorRed))
	modified := original.WithChar('B')

	// Modified should have new char
	assert.Equal(t, 'B', modified.Char())
	// Should preserve style
	assert.True(t, original.Style().Equals(modified.Style()))
	// Original should be unchanged
	assert.Equal(t, 'A', original.Char())
}

func TestCell_WithStyle(t *testing.T) {
	style1 := value.NewStyleWithFg(value.ColorRed)
	style2 := value.NewStyleWithFg(value.ColorBlue)

	original := NewCell('A', style1)
	modified := original.WithStyle(style2)

	// Modified should have new style
	assert.True(t, style2.Equals(modified.Style()))
	// Should preserve char
	assert.Equal(t, original.Char(), modified.Char())
	// Original should be unchanged
	assert.True(t, style1.Equals(original.Style()))
}

func TestCell_String(t *testing.T) {
	tests := []struct {
		name     string
		cell     Cell
		expected string
	}{
		{"empty", NewEmptyCell(), " "},
		{"ascii", NewCell('A', value.NewStyle()), "A"},
		{"emoji", NewCell('😀', value.NewStyle()), "😀"},
		{"cjk", NewCell('中', value.NewStyle()), "中"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cell.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateWidth(t *testing.T) {
	tests := []struct {
		name     string
		char     rune
		expected int
	}{
		{"ascii letter", 'A', 1},
		{"ascii digit", '5', 1},
		{"space", ' ', 1},
		{"emoji", '😀', 2},
		{"cjk", '中', 2},
		{"korean", '한', 2},
		{"japanese", 'あ', 2},
		{"zero width", '\u200B', 0},
		{"null", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateWidth(tt.char)
			assert.Equal(t, tt.expected, result, "Width mismatch for %c (U+%04X)", tt.char, tt.char)
		})
	}
}

// Benchmark cell creation
func BenchmarkNewCell(b *testing.B) {
	style := value.NewStyleWithFg(value.ColorRed)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = NewCell('A', style)
	}
}

// Benchmark cell equality check
func BenchmarkCell_Equals(b *testing.B) {
	c1 := NewCell('A', value.NewStyleWithFg(value.ColorRed))
	c2 := NewCell('A', value.NewStyleWithFg(value.ColorRed))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = c1.Equals(c2)
	}
}
