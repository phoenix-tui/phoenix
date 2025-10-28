package ansi

import (
	"bytes"
	"testing"

	"github.com/phoenix-tui/phoenix/render/internal/domain/model"
	value2 "github.com/phoenix-tui/phoenix/render/internal/domain/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWriter(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	assert.NotNil(t, w)
	assert.Equal(t, 0, w.BufferedSize())
}

func TestNewWriterWithBuffer(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriterWithBuffer(&buf, 1024)

	assert.NotNil(t, w)
}

func TestWriter_MoveCursor(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	err := w.MoveCursor(10, 5)
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, CSI)
	assert.Contains(t, output, "H") // Home position command
}

func TestWriter_MoveCursor_SamePosition(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	// Move to position
	err := w.MoveCursor(10, 5)
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	buf.Reset()

	// Move to same position (should be no-op)
	err = w.MoveCursor(10, 5)
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	assert.Empty(t, buf.String())
}

func TestWriter_MoveCursorRelative(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	err := w.MoveCursorRelative(5, 3)
	require.NoError(t, err)

	x, y := w.CurrentPosition()
	assert.Equal(t, 5, x)
	assert.Equal(t, 3, y)
}

func TestWriter_SetStyle(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	style := value2.NewStyleWithFg(value2.ColorRed)
	err := w.SetStyle(style)
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, CSI)
	assert.Contains(t, output, "38;2") // RGB foreground
}

func TestWriter_SetStyle_SameStyle(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	style := value2.NewStyleWithFg(value2.ColorRed)

	// Set style
	err := w.SetStyle(style)
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	buf.Reset()

	// Set same style (should be no-op)
	err = w.SetStyle(style)
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	assert.Empty(t, buf.String())
}

func TestWriter_WriteCell(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	cell := model.NewCell('A', value2.NewStyleWithFg(value2.ColorRed))
	err := w.WriteCell(cell)
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "A")
}

func TestWriter_WriteCellAt(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	pos := value2.NewPosition(10, 5)
	cell := model.NewCell('X', value2.NewStyle())

	err := w.WriteCellAt(pos, cell)
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "X")
	assert.Contains(t, output, CSI) // Should contain cursor move
}

func TestWriter_WriteString(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	style := value2.NewStyleWithFg(value2.ColorBlue)
	err := w.WriteString("Hello", style)
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Hello")
}

func TestWriter_WriteStringAt(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	pos := value2.NewPosition(5, 2)
	style := value2.NewStyle()

	err := w.WriteStringAt(pos, "Test", style)
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Test")
}

func TestWriter_Clear(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	err := w.Clear()
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)

	output := buf.String()
	assert.Equal(t, ClearScreen, output)
}

func TestWriter_ClearLine(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	err := w.ClearLine()
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)

	output := buf.String()
	assert.Equal(t, ClearLine, output)
}

func TestWriter_Reset(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	// Set a style
	style := value2.NewStyleWithFg(value2.ColorRed)
	err := w.SetStyle(style)
	require.NoError(t, err)

	// Reset
	err = w.Reset()
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)

	// Current style should be empty
	currentStyle := w.CurrentStyle()
	assert.True(t, currentStyle.IsEmpty())

	output := buf.String()
	assert.Contains(t, output, Reset)
}

func TestWriter_HideCursor(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	err := w.HideCursor()
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)

	output := buf.String()
	assert.Equal(t, CursorHide, output)
}

func TestWriter_ShowCursor(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	err := w.ShowCursor()
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)

	output := buf.String()
	assert.Equal(t, CursorShow, output)
}

func TestWriter_WriteRaw(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	data := []byte("raw data")
	err := w.WriteRaw(data)
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)

	output := buf.String()
	assert.Equal(t, "raw data", output)
}

func TestWriter_WriteEscape(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	err := w.WriteEscape(CursorHome)
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)

	output := buf.String()
	assert.Equal(t, CursorHome, output)
}

func TestWriter_CurrentPosition(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	// Initial position
	x, y := w.CurrentPosition()
	assert.Equal(t, 0, x)
	assert.Equal(t, 0, y)

	// After move
	err := w.MoveCursor(10, 5)
	require.NoError(t, err)

	x, y = w.CurrentPosition()
	assert.Equal(t, 10, x)
	assert.Equal(t, 5, y)
}

func TestWriter_CurrentStyle(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	// Initial style
	style := w.CurrentStyle()
	assert.True(t, style.IsEmpty())

	// After setting style
	newStyle := value2.NewStyleWithFg(value2.ColorRed)
	err := w.SetStyle(newStyle)
	require.NoError(t, err)

	currentStyle := w.CurrentStyle()
	assert.True(t, currentStyle.Equals(newStyle))
}

func TestWriter_SetPosition(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	w.SetPosition(20, 10)

	x, y := w.CurrentPosition()
	assert.Equal(t, 20, x)
	assert.Equal(t, 10, y)

	// Should not write anything to buffer
	err := w.Flush()
	require.NoError(t, err)
	assert.Empty(t, buf.String())
}

func TestWriter_BufferedSize(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	// Initially zero
	assert.Equal(t, 0, w.BufferedSize())

	// After writing
	err := w.WriteString("test", value2.NewStyle())
	require.NoError(t, err)

	assert.Greater(t, w.BufferedSize(), 0)

	// After flush
	err = w.Flush()
	require.NoError(t, err)

	assert.Equal(t, 0, w.BufferedSize())
}

func TestWriter_Close(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	err := w.WriteString("test", value2.NewStyle())
	require.NoError(t, err)

	err = w.Close()
	require.NoError(t, err)

	// Data should be flushed
	assert.Contains(t, buf.String(), "test")
}

func TestWriter_MultipleCells(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	style := value2.NewStyleWithFg(value2.ColorRed)
	cells := []model.Cell{
		model.NewCell('H', style),
		model.NewCell('e', style),
		model.NewCell('l', style),
		model.NewCell('l', style),
		model.NewCell('o', style),
	}

	for _, cell := range cells {
		err := w.WriteCell(cell)
		require.NoError(t, err)
	}

	err := w.Flush()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Hello")
}

// Benchmarks
func BenchmarkWriter_WriteCell(b *testing.B) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	cell := model.NewCell('A', value2.NewStyleWithFg(value2.ColorRed))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = w.WriteCell(cell)
	}
	_ = w.Flush()
}

func BenchmarkWriter_WriteString(b *testing.B) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	style := value2.NewStyle()
	text := "The quick brown fox"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = w.WriteString(text, style)
	}
	_ = w.Flush()
}

func BenchmarkWriter_MoveCursor(b *testing.B) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = w.MoveCursor(i%80, i%24)
	}
	_ = w.Flush()
}

func BenchmarkWriter_SetStyle(b *testing.B) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	style := value2.NewStyleWithFg(value2.ColorRed).WithBold(true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = w.SetStyle(style)
	}
	_ = w.Flush()
}
