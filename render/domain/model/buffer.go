// Package model defines domain models for terminal rendering (Buffer, Cell, Position).
package model

import (
	"strings"

	"github.com/phoenix-tui/phoenix/render/domain/value"
	"github.com/rivo/uniseg"
)

// Buffer represents a virtual terminal buffer (2D array of cells).
// Buffer is mutable but provides methods for safe operations.
type Buffer struct {
	width  int
	height int
	cells  [][]Cell // 2D array [y][x] of cells
}

// NewBuffer creates a new buffer with specified dimensions.
// All cells are initialized as empty.
func NewBuffer(width, height int) *Buffer {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}

	cells := make([][]Cell, height)
	for y := 0; y < height; y++ {
		cells[y] = make([]Cell, width)
		for x := 0; x < width; x++ {
			cells[y][x] = NewEmptyCell()
		}
	}

	return &Buffer{
		width:  width,
		height: height,
		cells:  cells,
	}
}

// Width returns the buffer width.
func (b *Buffer) Width() int {
	return b.width
}

// Height returns the buffer height.
func (b *Buffer) Height() int {
	return b.height
}

// Cells returns the 2D cell array (read-only access recommended).
func (b *Buffer) Cells() [][]Cell {
	return b.cells
}

// Get returns the cell at position. Returns empty cell if out of bounds.
func (b *Buffer) Get(pos value.Position) Cell {
	x, y := pos.X(), pos.Y()
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return NewEmptyCell()
	}
	return b.cells[y][x]
}

// Set sets the cell at position. Does nothing if out of bounds.
func (b *Buffer) Set(pos value.Position, cell Cell) {
	x, y := pos.X(), pos.Y()
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return
	}
	b.cells[y][x] = cell
}

// Clear resets all cells to empty.
func (b *Buffer) Clear() {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			b.cells[y][x] = NewEmptyCell()
		}
	}
}

// Fill fills all cells with character and style.
func (b *Buffer) Fill(char rune, style value.Style) {
	cell := NewCell(char, style)
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			b.cells[y][x] = cell
		}
	}
}

// SetString writes a string at position with style.
// Handles Unicode grapheme clusters correctly.
// Returns the number of cells written.
func (b *Buffer) SetString(pos value.Position, text string, style value.Style) int {
	x, y := pos.X(), pos.Y()
	if y < 0 || y >= b.height || x < 0 {
		return 0
	}

	cellsWritten := 0
	state := -1

	for text != "" {
		if x >= b.width {
			break
		}

		// Extract grapheme cluster.
		var cluster string
		cluster, text, _, state = uniseg.FirstGraphemeClusterInString(text, state)

		// Get rune and calculate width.
		runes := []rune(cluster)
		if len(runes) == 0 {
			continue
		}

		char := runes[0]
		width := uniseg.StringWidth(cluster)

		// Create and set cell.
		cell := NewCellWithWidth(char, style, width)
		b.Set(value.NewPosition(x, y), cell)

		x += width
		cellsWritten += width
	}

	return cellsWritten
}

// SetLine writes a string at the beginning of line y with style.
// Clears the rest of the line after the text.
func (b *Buffer) SetLine(y int, text string, style value.Style) {
	if y < 0 || y >= b.height {
		return
	}

	// Write text.
	pos := value.NewPosition(0, y)
	written := b.SetString(pos, text, style)

	// Clear rest of line.
	for x := written; x < b.width; x++ {
		b.Set(value.NewPosition(x, y), NewEmptyCell())
	}
}

// Clone creates a deep copy of the buffer.
func (b *Buffer) Clone() *Buffer {
	clone := NewBuffer(b.width, b.height)
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			clone.cells[y][x] = b.cells[y][x]
		}
	}
	return clone
}

// Resize creates a new buffer with different dimensions.
// Preserves content that fits in new dimensions.
func (b *Buffer) Resize(width, height int) *Buffer {
	newBuffer := NewBuffer(width, height)

	// Copy cells that fit.
	minHeight := minInt(b.height, height)
	minWidth := minInt(b.width, width)

	for y := 0; y < minHeight; y++ {
		for x := 0; x < minWidth; x++ {
			newBuffer.cells[y][x] = b.cells[y][x]
		}
	}

	return newBuffer
}

// GetLine returns all cells in line y as a slice.
func (b *Buffer) GetLine(y int) []Cell {
	if y < 0 || y >= b.height {
		return nil
	}
	// Return copy to prevent external modification.
	line := make([]Cell, b.width)
	copy(line, b.cells[y])
	return line
}

// SetCells sets multiple cells at once (batch operation).
func (b *Buffer) SetCells(positions []value.Position, cells []Cell) {
	if len(positions) != len(cells) {
		return
	}

	for i, pos := range positions {
		b.Set(pos, cells[i])
	}
}

// IsEmpty returns true if buffer is completely empty.
func (b *Buffer) IsEmpty() bool {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			if !b.cells[y][x].IsEmpty() {
				return false
			}
		}
	}
	return true
}

// String returns a string representation (for debugging).
func (b *Buffer) String() string {
	var sb strings.Builder
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			sb.WriteString(b.cells[y][x].String())
		}
		if y < b.height-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

// min returns the minimum of two integers.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
