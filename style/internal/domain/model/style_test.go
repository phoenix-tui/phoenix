package model

import (
	"testing"

	value2 "github.com/phoenix-tui/phoenix/style/internal/domain/value"
	"github.com/stretchr/testify/assert"
)

// TestNewStyle tests that NewStyle creates a style with default values.
func TestNewStyle(t *testing.T) {
	s := NewStyle()

	// Check default values
	_, hasFg := s.GetForeground()
	assert.False(t, hasFg, "default style should have no foreground color")

	_, hasBg := s.GetBackground()
	assert.False(t, hasBg, "default style should have no background color")

	_, hasBorder := s.GetBorder()
	assert.False(t, hasBorder, "default style should have no border")

	_, hasPadding := s.GetPadding()
	assert.False(t, hasPadding, "default style should have no padding")

	_, hasMargin := s.GetMargin()
	assert.False(t, hasMargin, "default style should have no margin")

	_, hasSize := s.GetSize()
	assert.False(t, hasSize, "default style should have no size constraints")

	_, hasAlignment := s.GetAlignment()
	assert.False(t, hasAlignment, "default style should have no alignment")

	assert.False(t, s.GetBold(), "default style should not be bold")
	assert.False(t, s.GetItalic(), "default style should not be italic")
	assert.False(t, s.GetUnderline(), "default style should not be underlined")
	assert.False(t, s.GetStrikethrough(), "default style should not have strikethrough")

	assert.Equal(t, value2.TrueColor, s.GetTerminalCapability(), "default terminal capability should be TrueColor")
}

// TestStyle_Foreground tests setting foreground color.
func TestStyle_Foreground(t *testing.T) {
	color := value2.RGB(255, 0, 0)
	s := NewStyle().Foreground(color)

	fg, hasFg := s.GetForeground()
	assert.True(t, hasFg, "foreground should be set")
	assert.Equal(t, color, fg, "foreground color should match")
}

// TestStyle_Background tests setting background color.
func TestStyle_Background(t *testing.T) {
	color := value2.RGB(0, 0, 255)
	s := NewStyle().Background(color)

	bg, hasBg := s.GetBackground()
	assert.True(t, hasBg, "background should be set")
	assert.Equal(t, color, bg, "background color should match")
}

// TestStyle_Border tests setting border.
func TestStyle_Border(t *testing.T) {
	border := value2.RoundedBorder
	s := NewStyle().Border(border)

	b, hasBorder := s.GetBorder()
	assert.True(t, hasBorder, "border should be set")
	assert.Equal(t, border, b, "border should match")

	// Border should enable all sides by default
	assert.True(t, s.GetBorderTop(), "top border should be enabled")
	assert.True(t, s.GetBorderBottom(), "bottom border should be enabled")
	assert.True(t, s.GetBorderLeft(), "left border should be enabled")
	assert.True(t, s.GetBorderRight(), "right border should be enabled")
}

// TestStyle_BorderColor tests setting border color.
func TestStyle_BorderColor(t *testing.T) {
	color := value2.RGB(128, 128, 128)
	s := NewStyle().BorderColor(color)

	bc, hasBorderColor := s.GetBorderColor()
	assert.True(t, hasBorderColor, "border color should be set")
	assert.Equal(t, color, bc, "border color should match")
}

// TestStyle_BorderSides tests enabling/disabling individual border sides.
func TestStyle_BorderSides(t *testing.T) {
	border := value2.NormalBorder
	s := NewStyle().
		Border(border).
		BorderTop(true).
		BorderBottom(false).
		BorderLeft(true).
		BorderRight(false)

	assert.True(t, s.GetBorderTop(), "top border should be enabled")
	assert.False(t, s.GetBorderBottom(), "bottom border should be disabled")
	assert.True(t, s.GetBorderLeft(), "left border should be enabled")
	assert.False(t, s.GetBorderRight(), "right border should be disabled")
}

// TestStyle_Padding tests setting padding.
func TestStyle_Padding(t *testing.T) {
	padding := value2.NewPadding(1, 2, 3, 4)
	s := NewStyle().Padding(padding)

	p, hasPadding := s.GetPadding()
	assert.True(t, hasPadding, "padding should be set")
	assert.Equal(t, padding, p, "padding should match")
}

// TestStyle_PaddingIndividual tests setting individual padding sides.
func TestStyle_PaddingIndividual(t *testing.T) {
	s := NewStyle().
		PaddingTop(1).
		PaddingRight(2).
		PaddingBottom(3).
		PaddingLeft(4)

	p, hasPadding := s.GetPadding()
	assert.True(t, hasPadding, "padding should be set")
	assert.Equal(t, 1, p.Top(), "top padding should be 1")
	assert.Equal(t, 2, p.Right(), "right padding should be 2")
	assert.Equal(t, 3, p.Bottom(), "bottom padding should be 3")
	assert.Equal(t, 4, p.Left(), "left padding should be 4")
}

// TestStyle_PaddingAll tests setting padding for all sides.
func TestStyle_PaddingAll(t *testing.T) {
	s := NewStyle().PaddingAll(2)

	p, hasPadding := s.GetPadding()
	assert.True(t, hasPadding, "padding should be set")
	assert.Equal(t, 2, p.Top(), "top padding should be 2")
	assert.Equal(t, 2, p.Right(), "right padding should be 2")
	assert.Equal(t, 2, p.Bottom(), "bottom padding should be 2")
	assert.Equal(t, 2, p.Left(), "left padding should be 2")
}

// TestStyle_PaddingHorizontalVertical tests setting horizontal/vertical padding.
func TestStyle_PaddingHorizontalVertical(t *testing.T) {
	s := NewStyle().
		PaddingHorizontal(3).
		PaddingVertical(1)

	p, hasPadding := s.GetPadding()
	assert.True(t, hasPadding, "padding should be set")
	assert.Equal(t, 1, p.Top(), "top padding should be 1")
	assert.Equal(t, 3, p.Right(), "right padding should be 3")
	assert.Equal(t, 1, p.Bottom(), "bottom padding should be 1")
	assert.Equal(t, 3, p.Left(), "left padding should be 3")
}

// TestStyle_Margin tests setting margin.
func TestStyle_Margin(t *testing.T) {
	margin := value2.NewMargin(1, 2, 3, 4)
	s := NewStyle().Margin(margin)

	m, hasMargin := s.GetMargin()
	assert.True(t, hasMargin, "margin should be set")
	assert.Equal(t, margin, m, "margin should match")
}

// TestStyle_MarginIndividual tests setting individual margin sides.
func TestStyle_MarginIndividual(t *testing.T) {
	s := NewStyle().
		MarginTop(1).
		MarginRight(2).
		MarginBottom(3).
		MarginLeft(4)

	m, hasMargin := s.GetMargin()
	assert.True(t, hasMargin, "margin should be set")
	assert.Equal(t, 1, m.Top(), "top margin should be 1")
	assert.Equal(t, 2, m.Right(), "right margin should be 2")
	assert.Equal(t, 3, m.Bottom(), "bottom margin should be 3")
	assert.Equal(t, 4, m.Left(), "left margin should be 4")
}

// TestStyle_MarginAll tests setting margin for all sides.
func TestStyle_MarginAll(t *testing.T) {
	s := NewStyle().MarginAll(2)

	m, hasMargin := s.GetMargin()
	assert.True(t, hasMargin, "margin should be set")
	assert.Equal(t, 2, m.Top(), "top margin should be 2")
	assert.Equal(t, 2, m.Right(), "right margin should be 2")
	assert.Equal(t, 2, m.Bottom(), "bottom margin should be 2")
	assert.Equal(t, 2, m.Left(), "left margin should be 2")
}

// TestStyle_MarginHorizontalVertical tests setting horizontal/vertical margin.
func TestStyle_MarginHorizontalVertical(t *testing.T) {
	s := NewStyle().
		MarginHorizontal(3).
		MarginVertical(1)

	m, hasMargin := s.GetMargin()
	assert.True(t, hasMargin, "margin should be set")
	assert.Equal(t, 1, m.Top(), "top margin should be 1")
	assert.Equal(t, 3, m.Right(), "right margin should be 3")
	assert.Equal(t, 1, m.Bottom(), "bottom margin should be 1")
	assert.Equal(t, 3, m.Left(), "left margin should be 3")
}

// TestStyle_Width tests setting width.
func TestStyle_Width(t *testing.T) {
	s := NewStyle().Width(80)

	size, hasSize := s.GetSize()
	assert.True(t, hasSize, "size should be set")
	w, hasWidth := size.Width()
	assert.True(t, hasWidth, "width should be set")
	assert.Equal(t, 80, w, "width should be 80")
}

// TestStyle_Height tests setting height.
func TestStyle_Height(t *testing.T) {
	s := NewStyle().Height(24)

	size, hasSize := s.GetSize()
	assert.True(t, hasSize, "size should be set")
	h, hasHeight := size.Height()
	assert.True(t, hasHeight, "height should be set")
	assert.Equal(t, 24, h, "height should be 24")
}

// TestStyle_MaxWidth tests setting maximum width.
func TestStyle_MaxWidth(t *testing.T) {
	s := NewStyle().MaxWidth(100)

	size, hasSize := s.GetSize()
	assert.True(t, hasSize, "size should be set")
	maxW, hasMaxWidth := size.MaxWidth()
	assert.True(t, hasMaxWidth, "max width should be set")
	assert.Equal(t, 100, maxW, "max width should be 100")
}

// TestStyle_MinWidth tests setting minimum width.
func TestStyle_MinWidth(t *testing.T) {
	s := NewStyle().MinWidth(10)

	size, hasSize := s.GetSize()
	assert.True(t, hasSize, "size should be set")
	minW, hasMinWidth := size.MinWidth()
	assert.True(t, hasMinWidth, "min width should be set")
	assert.Equal(t, 10, minW, "min width should be 10")
}

// TestStyle_MaxHeight tests setting maximum height.
func TestStyle_MaxHeight(t *testing.T) {
	s := NewStyle().MaxHeight(50)

	size, hasSize := s.GetSize()
	assert.True(t, hasSize, "size should be set")
	maxH, hasMaxHeight := size.MaxHeight()
	assert.True(t, hasMaxHeight, "max height should be set")
	assert.Equal(t, 50, maxH, "max height should be 50")
}

// TestStyle_MinHeight tests setting minimum height.
func TestStyle_MinHeight(t *testing.T) {
	s := NewStyle().MinHeight(5)

	size, hasSize := s.GetSize()
	assert.True(t, hasSize, "size should be set")
	minH, hasMinHeight := size.MinHeight()
	assert.True(t, hasMinHeight, "min height should be set")
	assert.Equal(t, 5, minH, "min height should be 5")
}

// TestStyle_Align tests setting alignment.
func TestStyle_Align(t *testing.T) {
	alignment := value2.NewAlignment(value2.AlignCenter, value2.AlignMiddle)
	s := NewStyle().Align(alignment)

	a, hasAlign := s.GetAlignment()
	assert.True(t, hasAlign, "alignment should be set")
	assert.Equal(t, alignment, a, "alignment should match")
}

// TestStyle_AlignHorizontal tests setting horizontal alignment.
func TestStyle_AlignHorizontal(t *testing.T) {
	s := NewStyle().AlignHorizontal(value2.AlignCenter)

	a, hasAlign := s.GetAlignment()
	assert.True(t, hasAlign, "alignment should be set")
	assert.Equal(t, value2.AlignCenter, a.Horizontal(), "horizontal alignment should be center")
	assert.Equal(t, value2.AlignTop, a.Vertical(), "vertical alignment should default to top")
}

// TestStyle_AlignVertical tests setting vertical alignment.
func TestStyle_AlignVertical(t *testing.T) {
	s := NewStyle().AlignVertical(value2.AlignMiddle)

	a, hasAlign := s.GetAlignment()
	assert.True(t, hasAlign, "alignment should be set")
	assert.Equal(t, value2.AlignLeft, a.Horizontal(), "horizontal alignment should default to left")
	assert.Equal(t, value2.AlignMiddle, a.Vertical(), "vertical alignment should be middle")
}

// TestStyle_Bold tests setting bold.
func TestStyle_Bold(t *testing.T) {
	s := NewStyle().Bold(true)
	assert.True(t, s.GetBold(), "bold should be enabled")

	s = NewStyle().Bold(false)
	assert.False(t, s.GetBold(), "bold should be disabled")
}

// TestStyle_Italic tests setting italic.
func TestStyle_Italic(t *testing.T) {
	s := NewStyle().Italic(true)
	assert.True(t, s.GetItalic(), "italic should be enabled")

	s = NewStyle().Italic(false)
	assert.False(t, s.GetItalic(), "italic should be disabled")
}

// TestStyle_Underline tests setting underline.
func TestStyle_Underline(t *testing.T) {
	s := NewStyle().Underline(true)
	assert.True(t, s.GetUnderline(), "underline should be enabled")

	s = NewStyle().Underline(false)
	assert.False(t, s.GetUnderline(), "underline should be disabled")
}

// TestStyle_Strikethrough tests setting strikethrough.
func TestStyle_Strikethrough(t *testing.T) {
	s := NewStyle().Strikethrough(true)
	assert.True(t, s.GetStrikethrough(), "strikethrough should be enabled")

	s = NewStyle().Strikethrough(false)
	assert.False(t, s.GetStrikethrough(), "strikethrough should be disabled")
}

// TestStyle_TerminalCapability tests setting terminal capability.
func TestStyle_TerminalCapability(t *testing.T) {
	s := NewStyle().TerminalCapability(value2.ANSI256)
	assert.Equal(t, value2.ANSI256, s.GetTerminalCapability(), "terminal capability should be ANSI256")

	s = NewStyle().TerminalCapability(value2.ANSI16)
	assert.Equal(t, value2.ANSI16, s.GetTerminalCapability(), "terminal capability should be ANSI16")
}

// TestStyle_Immutability tests that setter methods return new instances.
func TestStyle_Immutability(t *testing.T) {
	original := NewStyle()
	modified := original.Foreground(value2.RGB(255, 0, 0))

	// Original should be unchanged
	_, hasFg := original.GetForeground()
	assert.False(t, hasFg, "original style should not have foreground color")

	// Modified should have the color
	fg, hasFg := modified.GetForeground()
	assert.True(t, hasFg, "modified style should have foreground color")
	r, g, b := fg.RGB()
	assert.Equal(t, uint8(255), r, "red should be 255")
	assert.Equal(t, uint8(0), g, "green should be 0")
	assert.Equal(t, uint8(0), b, "blue should be 0")
}

// TestStyle_FluentAPI tests that methods are chainable.
func TestStyle_FluentAPI(t *testing.T) {
	s := NewStyle().
		Foreground(value2.RGB(255, 255, 255)).
		Background(value2.RGB(0, 0, 255)).
		Bold(true).
		Italic(true).
		Padding(value2.NewPadding(1, 2, 1, 2)).
		Border(value2.RoundedBorder).
		BorderColor(value2.RGB(128, 128, 128))

	// Verify all properties are set
	fg, hasFg := s.GetForeground()
	assert.True(t, hasFg, "foreground should be set")
	r, g, b := fg.RGB()
	assert.Equal(t, uint8(255), r)
	assert.Equal(t, uint8(255), g)
	assert.Equal(t, uint8(255), b)

	bg, hasBg := s.GetBackground()
	assert.True(t, hasBg, "background should be set")
	r, g, b = bg.RGB()
	assert.Equal(t, uint8(0), r)
	assert.Equal(t, uint8(0), g)
	assert.Equal(t, uint8(255), b)

	assert.True(t, s.GetBold(), "bold should be set")
	assert.True(t, s.GetItalic(), "italic should be set")

	_, hasPadding := s.GetPadding()
	assert.True(t, hasPadding, "padding should be set")

	_, hasBorder := s.GetBorder()
	assert.True(t, hasBorder, "border should be set")

	_, hasBorderColor := s.GetBorderColor()
	assert.True(t, hasBorderColor, "border color should be set")
}

// TestStyle_Validate tests style validation.
func TestStyle_Validate(t *testing.T) {
	tests := []struct {
		name    string
		style   Style
		wantErr bool
	}{
		{
			name:    "valid style with no properties",
			style:   NewStyle(),
			wantErr: false,
		},
		{
			name: "valid style with all properties",
			style: NewStyle().
				Foreground(value2.RGB(255, 0, 0)).
				Background(value2.RGB(0, 0, 255)).
				Border(value2.RoundedBorder).
				Padding(value2.NewPadding(1, 2, 1, 2)).
				Margin(value2.NewMargin(1, 1, 1, 1)).
				Width(80).
				Height(24),
			wantErr: false,
		},
		{
			name: "invalid: border sides enabled but no border",
			style: NewStyle().
				BorderTop(true),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.style.Validate()
			if tt.wantErr {
				assert.Error(t, err, "validation should fail")
			} else {
				assert.NoError(t, err, "validation should pass")
			}
		})
	}
}
