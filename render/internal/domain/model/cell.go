package model

import (
	"github.com/phoenix-tui/phoenix/render/internal/domain/value"
	"github.com/rivo/uniseg"
)

// Cell represents a single terminal cell (character + style + display width).
// Cell is immutable and uses value semantics.
type Cell struct {
	char  rune
	style value.Style
	width int
}

// NewCell creates a new cell with character and style.
// Width is calculated automatically using grapheme cluster analysis.
func NewCell(char rune, style value.Style) Cell {
	width := calculateWidth(char)
	return Cell{
		char:  char,
		style: style,
		width: width,
	}
}

// NewEmptyCell creates an empty cell (space with no style).
func NewEmptyCell() Cell {
	return Cell{
		char:  ' ',
		style: value.NewStyle(),
		width: 1,
	}
}

// NewCellWithWidth creates a cell with explicit width (for optimization).
func NewCellWithWidth(char rune, style value.Style, width int) Cell {
	return Cell{
		char:  char,
		style: style,
		width: width,
	}
}

// Char returns the character.
func (c Cell) Char() rune {
	return c.char
}

// Style returns the style.
func (c Cell) Style() value.Style {
	return c.style
}

// Width returns the display width.
func (c Cell) Width() int {
	return c.width
}

// IsEmpty returns true if cell is empty (space with no style).
func (c Cell) IsEmpty() bool {
	return c.char == ' ' && c.style.IsEmpty()
}

// Equals checks if two cells are equal (for diff optimization).
func (c Cell) Equals(other Cell) bool {
	return c.char == other.char &&
		c.width == other.width &&
		c.style.Equals(other.style)
}

// WithChar returns a new cell with different character.
func (c Cell) WithChar(char rune) Cell {
	return NewCell(char, c.style)
}

// WithStyle returns a new cell with different style.
func (c Cell) WithStyle(style value.Style) Cell {
	return Cell{
		char:  c.char,
		style: style,
		width: c.width,
	}
}

// String returns a string representation for debugging.
func (c Cell) String() string {
	if c.IsEmpty() {
		return " "
	}
	return string(c.char)
}

// calculateWidth calculates display width for a rune using uniseg.
func calculateWidth(r rune) int {
	if r == 0 {
		return 0
	}
	if r == ' ' {
		return 1
	}

	// Use uniseg for accurate width calculation.
	s := string(r)
	state := -1
	var width int

	for s != "" {
		var cluster string
		cluster, s, _, state = uniseg.FirstGraphemeClusterInString(s, state)
		width += uniseg.StringWidth(cluster)
	}

	return width
}
