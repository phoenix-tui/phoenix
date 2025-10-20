package render

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	r := New(80, 24, &buf)

	assert.NotNil(t, r)
	assert.Equal(t, 80, r.Width())
	assert.Equal(t, 24, r.Height())
}

func TestNewDefault(t *testing.T) {
	var buf bytes.Buffer
	r := NewDefault(&buf)

	assert.NotNil(t, r)
	assert.Equal(t, 80, r.Width())
	assert.Equal(t, 24, r.Height())
}

func TestNewStdout(t *testing.T) {
	r := NewStdout(80, 24)

	assert.NotNil(t, r)
	assert.Equal(t, 80, r.Width())
	assert.Equal(t, 24, r.Height())
}

func TestWithBufferSize(t *testing.T) {
	var buf bytes.Buffer
	r := WithBufferSize(80, 24, &buf, 8192)

	assert.NotNil(t, r)
}

func TestRenderer_Buffer(t *testing.T) {
	var buf bytes.Buffer
	r := New(80, 24, &buf)

	buffer := r.Buffer()
	assert.NotNil(t, buffer)
	assert.Equal(t, 80, buffer.Width())
	assert.Equal(t, 24, buffer.Height())
}

func TestRenderer_Render(t *testing.T) {
	var buf bytes.Buffer
	r := New(80, 24, &buf)

	buffer := r.Buffer()
	buffer.SetString(0, 0, "Hello", StyleDefault())

	err := r.Render(buffer)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Hello")
}

func TestRenderer_Clear(t *testing.T) {
	var buf bytes.Buffer
	r := New(80, 24, &buf)

	err := r.Clear()
	require.NoError(t, err)

	output := buf.String()
	assert.NotEmpty(t, output)
}

func TestRenderer_Resize(t *testing.T) {
	var buf bytes.Buffer
	r := New(80, 24, &buf)

	err := r.Resize(100, 30)
	require.NoError(t, err)

	assert.Equal(t, 100, r.Width())
	assert.Equal(t, 30, r.Height())
}

func TestRenderer_HideShowCursor(t *testing.T) {
	var buf bytes.Buffer
	r := New(80, 24, &buf)

	err := r.HideCursor()
	require.NoError(t, err)

	err = r.ShowCursor()
	require.NoError(t, err)
}

func TestRenderer_Close(t *testing.T) {
	var buf bytes.Buffer
	r := New(80, 24, &buf)

	err := r.Close()
	require.NoError(t, err)
}

func TestBuffer_Set(_ *testing.T) {
	buffer := NewBuffer(80, 24)
	style := StyleDefault()

	buffer.Set(10, 5, 'A', style)

	// Buffer should accept the operation without error
}

func TestBuffer_SetString(t *testing.T) {
	buffer := NewBuffer(80, 24)
	style := StyleFg(255, 0, 0)

	written := buffer.SetString(0, 0, "Hello", style)
	assert.Equal(t, 5, written)
}

func TestBuffer_SetLine(_ *testing.T) {
	buffer := NewBuffer(80, 24)
	style := StyleDefault()

	buffer.SetLine(5, "Test Line", style)

	// Should execute without error
}

func TestBuffer_Clear(_ *testing.T) {
	buffer := NewBuffer(80, 24)

	buffer.SetString(0, 0, "Test", StyleDefault())
	buffer.Clear()

	// Should execute without error
}

func TestBuffer_Fill(_ *testing.T) {
	buffer := NewBuffer(80, 24)
	style := StyleBg(0, 0, 255)

	buffer.Fill('X', style)

	// Should execute without error
}

func TestBuffer_Release(_ *testing.T) {
	var buf bytes.Buffer
	r := New(80, 24, &buf)

	buffer := r.Buffer()
	buffer.Release()

	// Should execute without error
}

func TestStyleDefault(t *testing.T) {
	style := StyleDefault()
	assert.NotNil(t, style)
}

func TestStyleFg(t *testing.T) {
	style := StyleFg(255, 0, 0)
	assert.NotNil(t, style)
}

func TestStyleBg(t *testing.T) {
	style := StyleBg(0, 0, 255)
	assert.NotNil(t, style)
}

func TestStyleColors(t *testing.T) {
	style := StyleColors(255, 0, 0, 0, 0, 255)
	assert.NotNil(t, style)
}

func TestStyle_WithMethods(t *testing.T) {
	style := StyleDefault()

	style = style.WithFg(255, 0, 0)
	assert.NotNil(t, style)

	style = style.WithBg(0, 0, 255)
	assert.NotNil(t, style)

	style = style.WithBold(true)
	assert.NotNil(t, style)

	style = style.WithItalic(true)
	assert.NotNil(t, style)

	style = style.WithUnderline(true)
	assert.NotNil(t, style)

	style = style.WithReverse(true)
	assert.NotNil(t, style)
}

func TestPredefinedStyles(t *testing.T) {
	// Test predefined styles don't panic
	assert.NotNil(t, Bold)
	assert.NotNil(t, Italic)
	assert.NotNil(t, Underline)
	assert.NotNil(t, Reverse)

	assert.NotNil(t, FgRed)
	assert.NotNil(t, FgGreen)
	assert.NotNil(t, FgBlue)

	assert.NotNil(t, BgRed)
	assert.NotNil(t, BgGreen)
	assert.NotNil(t, BgBlue)
}

func TestPredefinedColors(t *testing.T) {
	// Test predefined colors exist
	assert.NotNil(t, ColorBlack)
	assert.NotNil(t, ColorRed)
	assert.NotNil(t, ColorGreen)
	assert.NotNil(t, ColorBlue)
}

// Example usage test
func TestExampleUsage(t *testing.T) {
	var buf bytes.Buffer
	renderer := New(80, 24, &buf)
	defer renderer.Close()

	// Create buffer
	buffer := renderer.Buffer()
	defer buffer.Release()

	// Write some text
	buffer.SetString(0, 0, "Hello, World!", FgRed)
	buffer.SetString(0, 1, "Phoenix TUI", FgGreen.WithBold(true))

	// Render
	err := renderer.Render(buffer)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Hello, World!")
	assert.Contains(t, output, "Phoenix TUI")
}

// Benchmarks
func BenchmarkRenderer_Render(b *testing.B) {
	var buf bytes.Buffer
	r := New(80, 24, &buf)
	buffer := r.Buffer()
	buffer.SetString(0, 0, "Benchmark test", StyleDefault())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Render(buffer)
	}
}

func BenchmarkBuffer_SetString(b *testing.B) {
	buffer := NewBuffer(80, 24)
	style := StyleDefault()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.SetString(0, 0, "The quick brown fox", style)
	}
}

func BenchmarkStyle_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StyleFg(255, 128, 64).WithBold(true).WithItalic(true)
	}
}
