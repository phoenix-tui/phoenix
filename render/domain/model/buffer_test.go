package model

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/render/domain/value"
	"github.com/stretchr/testify/assert"
)

func TestNewBuffer(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"small", 10, 5},
		{"large", 80, 24},
		{"zero width", 0, 10},
		{"zero height", 10, 0},
		{"negative width", -5, 10},
		{"negative height", 10, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBuffer(tt.width, tt.height)
			assert.NotNil(t, buf)

			expectedWidth := tt.width
			expectedHeight := tt.height
			if expectedWidth < 0 {
				expectedWidth = 0
			}
			if expectedHeight < 0 {
				expectedHeight = 0
			}

			assert.Equal(t, expectedWidth, buf.Width())
			assert.Equal(t, expectedHeight, buf.Height())
			assert.True(t, buf.IsEmpty())
		})
	}
}

func TestBuffer_GetSet(t *testing.T) {
	buf := NewBuffer(10, 5)
	pos := value.NewPosition(3, 2)
	cell := NewCell('A', value.NewStyleWithFg(value.ColorRed))

	// Set cell
	buf.Set(pos, cell)

	// Get cell
	retrieved := buf.Get(pos)
	assert.True(t, cell.Equals(retrieved))
}

func TestBuffer_GetSet_OutOfBounds(t *testing.T) {
	buf := NewBuffer(10, 5)
	cell := NewCell('A', value.NewStyle())

	tests := []struct {
		name string
		pos  value.Position
	}{
		{"negative x", value.NewPosition(-1, 2)},
		{"negative y", value.NewPosition(3, -1)},
		{"x too large", value.NewPosition(10, 2)},
		{"y too large", value.NewPosition(3, 5)},
		{"both negative", value.NewPosition(-1, -1)},
		{"both too large", value.NewPosition(10, 5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set should do nothing
			buf.Set(tt.pos, cell)

			// Get should return empty cell
			retrieved := buf.Get(tt.pos)
			assert.True(t, retrieved.IsEmpty())
		})
	}
}

func TestBuffer_Clear(t *testing.T) {
	buf := NewBuffer(10, 5)

	// Fill with non-empty cells
	cell := NewCell('A', value.NewStyleWithFg(value.ColorRed))
	for y := 0; y < buf.Height(); y++ {
		for x := 0; x < buf.Width(); x++ {
			buf.Set(value.NewPosition(x, y), cell)
		}
	}

	assert.False(t, buf.IsEmpty())

	// Clear
	buf.Clear()

	// Verify all empty
	assert.True(t, buf.IsEmpty())
	for y := 0; y < buf.Height(); y++ {
		for x := 0; x < buf.Width(); x++ {
			cell := buf.Get(value.NewPosition(x, y))
			assert.True(t, cell.IsEmpty())
		}
	}
}

func TestBuffer_Fill(t *testing.T) {
	buf := NewBuffer(10, 5)
	style := value.NewStyleWithFg(value.ColorRed)

	buf.Fill('X', style)

	// Verify all cells filled
	for y := 0; y < buf.Height(); y++ {
		for x := 0; x < buf.Width(); x++ {
			cell := buf.Get(value.NewPosition(x, y))
			assert.Equal(t, 'X', cell.Char())
			assert.True(t, style.Equals(cell.Style()))
		}
	}
}

func TestBuffer_SetString(t *testing.T) {
	buf := NewBuffer(20, 5)
	style := value.NewStyleWithFg(value.ColorRed)

	tests := []struct {
		name          string
		pos           value.Position
		text          string
		expectedCells int
	}{
		{"ascii text", value.NewPosition(0, 0), "Hello", 5},
		{"with spaces", value.NewPosition(0, 1), "Hello World", 11},
		{"emoji", value.NewPosition(0, 2), "😀😃", 4},  // 2 emojis = 4 cells
		{"cjk", value.NewPosition(0, 3), "中文", 4},    // 2 CJK chars = 4 cells
		{"mixed", value.NewPosition(0, 4), "Hi😀", 4}, // Hi(2) + emoji(2) = 4
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			written := buf.SetString(tt.pos, tt.text, style)
			assert.Equal(t, tt.expectedCells, written)

			// Verify cells have style
			y := tt.pos.Y()
			for x := tt.pos.X(); x < tt.pos.X()+written; x++ {
				cell := buf.Get(value.NewPosition(x, y))
				if !cell.IsEmpty() {
					assert.True(t, style.Equals(cell.Style()))
				}
			}
		})
	}
}

func TestBuffer_SetString_OutOfBounds(t *testing.T) {
	buf := NewBuffer(10, 5)
	style := value.NewStyle()

	tests := []struct {
		name string
		pos  value.Position
		text string
	}{
		{"y negative", value.NewPosition(0, -1), "test"},
		{"y too large", value.NewPosition(0, 5), "test"},
		{"x negative", value.NewPosition(-1, 0), "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			written := buf.SetString(tt.pos, tt.text, style)
			assert.Equal(t, 0, written)
		})
	}
}

func TestBuffer_SetString_Truncate(t *testing.T) {
	buf := NewBuffer(10, 5)
	style := value.NewStyle()

	// Text longer than width
	written := buf.SetString(value.NewPosition(0, 0), "This is a very long text", style)

	// Should truncate at buffer width
	assert.LessOrEqual(t, written, buf.Width())
}

func TestBuffer_SetLine(t *testing.T) {
	buf := NewBuffer(20, 5)
	style := value.NewStyleWithFg(value.ColorRed)
	text := "Hello"

	buf.SetLine(2, text, style)

	// Verify text written
	for x := 0; x < 5; x++ {
		cell := buf.Get(value.NewPosition(x, 2))
		assert.False(t, cell.IsEmpty())
		assert.True(t, style.Equals(cell.Style()))
	}

	// Verify rest of line is empty
	for x := 5; x < buf.Width(); x++ {
		cell := buf.Get(value.NewPosition(x, 2))
		assert.True(t, cell.IsEmpty())
	}
}

func TestBuffer_Clone(t *testing.T) {
	buf := NewBuffer(10, 5)
	style := value.NewStyleWithFg(value.ColorRed)
	buf.SetString(value.NewPosition(0, 0), "Test", style)

	clone := buf.Clone()

	// Verify clone has same dimensions
	assert.Equal(t, buf.Width(), clone.Width())
	assert.Equal(t, buf.Height(), clone.Height())

	// Verify clone has same content
	for y := 0; y < buf.Height(); y++ {
		for x := 0; x < buf.Width(); x++ {
			pos := value.NewPosition(x, y)
			assert.True(t, buf.Get(pos).Equals(clone.Get(pos)))
		}
	}

	// Verify clone is independent (modify original)
	buf.Set(value.NewPosition(0, 0), NewCell('X', style))
	assert.False(t, buf.Get(value.NewPosition(0, 0)).Equals(clone.Get(value.NewPosition(0, 0))))
}

func TestBuffer_Resize(t *testing.T) {
	buf := NewBuffer(10, 5)
	style := value.NewStyleWithFg(value.ColorRed)
	buf.SetString(value.NewPosition(0, 0), "Test", style)

	tests := []struct {
		name      string
		newWidth  int
		newHeight int
	}{
		{"larger", 20, 10},
		{"smaller", 5, 3},
		{"same", 10, 5},
		{"wider", 20, 5},
		{"taller", 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resized := buf.Resize(tt.newWidth, tt.newHeight)

			assert.Equal(t, tt.newWidth, resized.Width())
			assert.Equal(t, tt.newHeight, resized.Height())

			// Verify content preserved where it fits
			minWidth := min(buf.Width(), tt.newWidth)
			minHeight := min(buf.Height(), tt.newHeight)

			for y := 0; y < minHeight; y++ {
				for x := 0; x < minWidth; x++ {
					pos := value.NewPosition(x, y)
					assert.True(t, buf.Get(pos).Equals(resized.Get(pos)))
				}
			}
		})
	}
}

func TestBuffer_GetLine(t *testing.T) {
	buf := NewBuffer(10, 5)
	style := value.NewStyleWithFg(value.ColorRed)
	buf.SetLine(2, "Test", style)

	line := buf.GetLine(2)
	assert.Equal(t, buf.Width(), len(line))

	// Verify content
	assert.Equal(t, 'T', line[0].Char())
	assert.Equal(t, 'e', line[1].Char())
	assert.Equal(t, 's', line[2].Char())
	assert.Equal(t, 't', line[3].Char())

	// Rest should be empty
	for i := 4; i < len(line); i++ {
		assert.True(t, line[i].IsEmpty())
	}
}

func TestBuffer_GetLine_OutOfBounds(t *testing.T) {
	buf := NewBuffer(10, 5)

	tests := []struct {
		name string
		y    int
	}{
		{"negative", -1},
		{"too large", 5},
		{"way too large", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := buf.GetLine(tt.y)
			assert.Nil(t, line)
		})
	}
}

func TestBuffer_SetCells(t *testing.T) {
	buf := NewBuffer(10, 5)
	style := value.NewStyleWithFg(value.ColorRed)

	positions := []value.Position{
		value.NewPosition(0, 0),
		value.NewPosition(1, 1),
		value.NewPosition(2, 2),
	}
	cells := []Cell{
		NewCell('A', style),
		NewCell('B', style),
		NewCell('C', style),
	}

	buf.SetCells(positions, cells)

	// Verify cells set
	assert.Equal(t, 'A', buf.Get(positions[0]).Char())
	assert.Equal(t, 'B', buf.Get(positions[1]).Char())
	assert.Equal(t, 'C', buf.Get(positions[2]).Char())
}

func TestBuffer_SetCells_MismatchedLength(t *testing.T) {
	buf := NewBuffer(10, 5)
	style := value.NewStyleWithFg(value.ColorRed)

	positions := []value.Position{
		value.NewPosition(0, 0),
	}
	cells := []Cell{
		NewCell('A', style),
		NewCell('B', style),
	}

	// Should do nothing with mismatched lengths
	buf.SetCells(positions, cells)
	assert.True(t, buf.IsEmpty())
}

func TestBuffer_IsEmpty(t *testing.T) {
	buf := NewBuffer(10, 5)
	assert.True(t, buf.IsEmpty())

	// Add one cell
	buf.Set(value.NewPosition(5, 2), NewCell('A', value.NewStyle()))
	assert.False(t, buf.IsEmpty())

	// Clear
	buf.Clear()
	assert.True(t, buf.IsEmpty())
}

func TestBuffer_String(t *testing.T) {
	buf := NewBuffer(5, 3)
	buf.SetLine(0, "Hello", value.NewStyle())
	buf.SetLine(1, "World", value.NewStyle())

	result := buf.String()
	lines := strings.Split(result, "\n")

	assert.Equal(t, 3, len(lines))
	assert.Equal(t, "Hello", lines[0])
	assert.Equal(t, "World", lines[1])
}

// Benchmarks
func BenchmarkBuffer_NewBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewBuffer(80, 24)
	}
}

func BenchmarkBuffer_Set(b *testing.B) {
	buf := NewBuffer(80, 24)
	cell := NewCell('A', value.NewStyle())
	pos := value.NewPosition(40, 12)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.Set(pos, cell)
	}
}

func BenchmarkBuffer_SetString(b *testing.B) {
	buf := NewBuffer(80, 24)
	style := value.NewStyle()
	text := "The quick brown fox jumps over the lazy dog"
	pos := value.NewPosition(0, 0)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.SetString(pos, text, style)
	}
}

func BenchmarkBuffer_Clone(b *testing.B) {
	buf := NewBuffer(80, 24)
	buf.Fill('X', value.NewStyle())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = buf.Clone()
	}
}
