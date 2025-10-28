package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStyle(t *testing.T) {
	style := NewStyle()
	assert.Nil(t, style.Foreground())
	assert.Nil(t, style.Background())
	assert.False(t, style.Bold())
	assert.False(t, style.Italic())
	assert.False(t, style.Underline())
	assert.True(t, style.IsEmpty())
}

func TestNewStyleWithFg(t *testing.T) {
	color := ColorRed
	style := NewStyleWithFg(color)

	assert.NotNil(t, style.Foreground())
	assert.True(t, style.Foreground().Equals(color))
	assert.Nil(t, style.Background())
	assert.False(t, style.IsEmpty())
}

func TestNewStyleWithBg(t *testing.T) {
	color := ColorBlue
	style := NewStyleWithBg(color)

	assert.Nil(t, style.Foreground())
	assert.NotNil(t, style.Background())
	assert.True(t, style.Background().Equals(color))
	assert.False(t, style.IsEmpty())
}

func TestNewStyleWithColors(t *testing.T) {
	fg := ColorRed
	bg := ColorBlue
	style := NewStyleWithColors(fg, bg)

	assert.NotNil(t, style.Foreground())
	assert.True(t, style.Foreground().Equals(fg))
	assert.NotNil(t, style.Background())
	assert.True(t, style.Background().Equals(bg))
	assert.False(t, style.IsEmpty())
}

func TestStyle_WithMethods(t *testing.T) {
	style := NewStyle()

	// Test WithFg
	style = style.WithFg(ColorRed)
	assert.NotNil(t, style.Foreground())
	assert.True(t, style.Foreground().Equals(ColorRed))

	// Test WithBg
	style = style.WithBg(ColorBlue)
	assert.NotNil(t, style.Background())
	assert.True(t, style.Background().Equals(ColorBlue))

	// Test WithBold
	style = style.WithBold(true)
	assert.True(t, style.Bold())

	// Test WithItalic
	style = style.WithItalic(true)
	assert.True(t, style.Italic())

	// Test WithUnderline
	style = style.WithUnderline(true)
	assert.True(t, style.Underline())

	// Test WithReverse
	style = style.WithReverse(true)
	assert.True(t, style.Reverse())

	// Test WithDim
	style = style.WithDim(true)
	assert.True(t, style.Dim())

	// Test WithBlink
	style = style.WithBlink(true)
	assert.True(t, style.Blink())

	// Test WithHidden
	style = style.WithHidden(true)
	assert.True(t, style.Hidden())

	// Test WithStrike
	style = style.WithStrike(true)
	assert.True(t, style.Strike())
}

func TestStyle_Immutability(t *testing.T) {
	original := NewStyle()
	modified := original.WithBold(true).WithFg(ColorRed)

	// Original should be unchanged
	assert.False(t, original.Bold())
	assert.Nil(t, original.Foreground())

	// Modified should have changes
	assert.True(t, modified.Bold())
	assert.NotNil(t, modified.Foreground())
}

func TestStyle_Equals(t *testing.T) {
	tests := []struct {
		name     string
		s1, s2   Style
		expected bool
	}{
		{
			"empty styles",
			NewStyle(),
			NewStyle(),
			true,
		},
		{
			"same foreground",
			NewStyleWithFg(ColorRed),
			NewStyleWithFg(ColorRed),
			true,
		},
		{
			"different foreground",
			NewStyleWithFg(ColorRed),
			NewStyleWithFg(ColorBlue),
			false,
		},
		{
			"same background",
			NewStyleWithBg(ColorRed),
			NewStyleWithBg(ColorRed),
			true,
		},
		{
			"one has fg, other doesn't",
			NewStyleWithFg(ColorRed),
			NewStyle(),
			false,
		},
		{
			"same attributes",
			NewStyle().WithBold(true).WithItalic(true),
			NewStyle().WithBold(true).WithItalic(true),
			true,
		},
		{
			"different attributes",
			NewStyle().WithBold(true),
			NewStyle().WithItalic(true),
			false,
		},
		{
			"complex equal",
			NewStyleWithColors(ColorRed, ColorBlue).WithBold(true).WithUnderline(true),
			NewStyleWithColors(ColorRed, ColorBlue).WithBold(true).WithUnderline(true),
			true,
		},
		{
			"complex different",
			NewStyleWithColors(ColorRed, ColorBlue).WithBold(true),
			NewStyleWithColors(ColorRed, ColorBlue).WithItalic(true),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.s1.Equals(tt.s2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStyle_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		style    Style
		expected bool
	}{
		{"empty", NewStyle(), true},
		{"with foreground", NewStyleWithFg(ColorRed), false},
		{"with background", NewStyleWithBg(ColorBlue), false},
		{"with bold", NewStyle().WithBold(true), false},
		{"with italic", NewStyle().WithItalic(true), false},
		{"with underline", NewStyle().WithUnderline(true), false},
		{"complex", NewStyleWithColors(ColorRed, ColorBlue).WithBold(true), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.style.IsEmpty()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStyle_ToANSI(t *testing.T) {
	tests := []struct {
		name     string
		style    Style
		expected string
	}{
		{
			"empty",
			NewStyle(),
			"",
		},
		{
			"foreground only",
			NewStyleWithFg(ColorRed),
			"\x1b[38;2;255;0;0m",
		},
		{
			"background only",
			NewStyleWithBg(ColorBlue),
			"\x1b[48;2;0;0;255m",
		},
		{
			"bold only",
			NewStyle().WithBold(true),
			"\x1b[1m",
		},
		{
			"italic only",
			NewStyle().WithItalic(true),
			"\x1b[3m",
		},
		{
			"underline only",
			NewStyle().WithUnderline(true),
			"\x1b[4m",
		},
		{
			"foreground and bold",
			NewStyleWithFg(ColorRed).WithBold(true),
			"\x1b[38;2;255;0;0;1m",
		},
		{
			"complex style",
			NewStyleWithColors(ColorRed, ColorBlue).WithBold(true).WithItalic(true),
			"\x1b[38;2;255;0;0;48;2;0;0;255;1;3m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.style.ToANSI()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStyle_String(t *testing.T) {
	tests := []struct {
		name     string
		style    Style
		contains []string
	}{
		{"empty", NewStyle(), []string{"empty"}},
		{"foreground", NewStyleWithFg(ColorRed), []string{"fg:", "255, 0, 0"}},
		{"background", NewStyleWithBg(ColorBlue), []string{"bg:", "0, 0, 255"}},
		{"bold", NewStyle().WithBold(true), []string{"bold"}},
		{"complex", NewStyleWithFg(ColorRed).WithBold(true).WithItalic(true), []string{"fg:", "bold", "italic"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.style.String()
			for _, substr := range tt.contains {
				assert.Contains(t, result, substr)
			}
		})
	}
}
