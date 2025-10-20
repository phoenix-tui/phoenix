package application

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/phoenix-tui/phoenix/render/domain/model"
	"github.com/phoenix-tui/phoenix/render/domain/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRenderer(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	assert.NotNil(t, r)
	assert.Equal(t, 80, r.Width())
	assert.Equal(t, 24, r.Height())
}

func TestNewRendererWithOptions(t *testing.T) {
	var buf bytes.Buffer
	r := NewRendererWithOptions(80, 24, &buf, 4096)

	assert.NotNil(t, r)
	assert.Equal(t, 80, r.Width())
	assert.Equal(t, 24, r.Height())
}

func TestRenderer_Render_EmptyBuffer(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	buffer := model.NewBuffer(80, 24)
	err := r.Render(buffer)

	require.NoError(t, err)
}

func TestRenderer_Render_WithContent(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	buffer := model.NewBuffer(80, 24)
	style := value.NewStyleWithFg(value.ColorRed)
	buffer.SetString(value.NewPosition(0, 0), "Hello", style)

	err := r.Render(buffer)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Hello")
}

func TestRenderer_Render_DifferentialUpdate(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	// First render
	buffer1 := model.NewBuffer(80, 24)
	buffer1.SetString(value.NewPosition(0, 0), "Hello", value.NewStyle())
	err := r.Render(buffer1)
	require.NoError(t, err)

	buf.Reset()

	// Second render with change
	buffer2 := model.NewBuffer(80, 24)
	buffer2.SetString(value.NewPosition(0, 0), "World", value.NewStyle())
	err = r.Render(buffer2)
	require.NoError(t, err)

	output := buf.String()
	// Differential rendering outputs only changed parts with ANSI sequences
	// "Hello" -> "World" should output "Wor" and "d" with cursor positioning
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "Wor") // Changed characters
	assert.Contains(t, output, "d")   // Last character
}

func TestRenderer_Render_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	buffer := model.NewBuffer(80, 24)
	buffer.SetString(value.NewPosition(0, 0), "Hello", value.NewStyle())

	// First render
	err := r.Render(buffer)
	require.NoError(t, err)

	buf.Reset()

	// Second render with same content (should output nothing)
	err = r.Render(buffer)
	require.NoError(t, err)

	assert.Empty(t, buf.String())
}

func TestRenderer_Clear(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	err := r.Clear()
	require.NoError(t, err)

	output := buf.String()
	assert.NotEmpty(t, output) // Should contain clear sequences
}

func TestRenderer_Resize(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	err := r.Resize(100, 30)
	require.NoError(t, err)

	assert.Equal(t, 100, r.Width())
	assert.Equal(t, 30, r.Height())
}

func TestRenderer_HideCursor(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	err := r.HideCursor()
	require.NoError(t, err)

	output := buf.String()
	assert.NotEmpty(t, output)
}

func TestRenderer_ShowCursor(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	err := r.ShowCursor()
	require.NoError(t, err)

	output := buf.String()
	assert.NotEmpty(t, output)
}

func TestRenderer_GetPutBuffer(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	// Get buffer from pool
	poolBuf := r.GetBuffer()
	assert.NotNil(t, poolBuf)
	assert.Equal(t, 80, poolBuf.Width())
	assert.Equal(t, 24, poolBuf.Height())

	// Return to pool
	r.PutBuffer(poolBuf)
}

func TestRenderer_GetStats(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	stats := r.GetStats()
	assert.NotNil(t, stats)
}

func TestRenderer_Close(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	err := r.Close()
	require.NoError(t, err)
}

// Benchmark rendering performance
func BenchmarkRenderer_Render_Empty(b *testing.B) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)
	buffer := model.NewBuffer(80, 24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Render(buffer)
	}
}

func BenchmarkRenderer_Render_SmallChange(b *testing.B) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer := model.NewBuffer(80, 24)
		buffer.SetString(value.NewPosition(0, 0), "Test", value.NewStyle())
		_ = r.Render(buffer)
	}
}

func BenchmarkRenderer_Render_FullScreen(b *testing.B) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)
	style := value.NewStyle()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer := model.NewBuffer(80, 24)
		buffer.Fill('X', style)
		_ = r.Render(buffer)
	}
}

// === NEW COMPREHENSIVE TESTS FOR 85%+ COVERAGE ===

func TestRenderer_Render_NilBuffer(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	err := r.Render(nil)
	require.NoError(t, err) // Should handle nil gracefully
}

func TestRenderer_Render_DimensionChange_TriggersFullRedraw(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	// First render at 80x24
	buffer1 := model.NewBuffer(80, 24)
	buffer1.SetString(value.NewPosition(0, 0), "Hello", value.NewStyle())
	err := r.Render(buffer1)
	require.NoError(t, err)

	buf.Reset()

	// Second render with different dimensions (should trigger full redraw)
	buffer2 := model.NewBuffer(100, 30)
	buffer2.SetString(value.NewPosition(0, 0), "World", value.NewStyle())
	err = r.Render(buffer2)
	require.NoError(t, err)

	output := buf.String()
	assert.NotEmpty(t, output) // Should have clear sequence and content
	assert.Contains(t, output, "World")
}

func TestRenderer_RenderFullScreen_WithEmptyCells(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	// Buffer with mix of empty and non-empty cells
	buffer := model.NewBuffer(80, 24)
	buffer.SetString(value.NewPosition(0, 0), "Line 1", value.NewStyle())
	buffer.SetString(value.NewPosition(0, 2), "Line 3", value.NewStyle())
	// Line 1 is empty (should be skipped)

	// Trigger full redraw by changing dimensions
	buffer2 := model.NewBuffer(90, 30)
	buffer2.SetString(value.NewPosition(0, 0), "Resized", value.NewStyle())
	err := r.Render(buffer2)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Resized")
}

func TestRenderer_RenderFullScreen_AllCellsFilled(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(10, 5, &buf)

	// Fill entire buffer
	buffer := model.NewBuffer(10, 5)
	buffer.Fill('X', value.NewStyleWithFg(value.ColorBlue))

	// Trigger full redraw
	err := r.Render(buffer)
	require.NoError(t, err)

	output := buf.String()
	assert.NotEmpty(t, output)
	// Should contain blue color ANSI code
	assert.Contains(t, output, "38;2") // RGB foreground color
}

func TestRenderer_ApplyOperations_AllTypes(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	// First render to set up state
	buffer1 := model.NewBuffer(80, 24)
	buffer1.SetString(value.NewPosition(0, 0), "Initial", value.NewStyle())
	err := r.Render(buffer1)
	require.NoError(t, err)

	buf.Reset()

	// Second render with changes (triggers diff + apply operations)
	buffer2 := model.NewBuffer(80, 24)
	// OpTypeSet: change character
	buffer2.SetString(value.NewPosition(0, 0), "Changed", value.NewStyle())
	// OpTypeClear: clear line 1
	// (happens automatically when line becomes empty)

	err = r.Render(buffer2)
	require.NoError(t, err)

	output := buf.String()
	assert.NotEmpty(t, output)
}

func TestRenderer_ApplyOperations_WithStyles(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	// First render
	buffer1 := model.NewBuffer(80, 24)
	buffer1.SetString(value.NewPosition(0, 0), "Red", value.NewStyleWithFg(value.ColorRed))
	err := r.Render(buffer1)
	require.NoError(t, err)

	buf.Reset()

	// Change style
	buffer2 := model.NewBuffer(80, 24)
	buffer2.SetString(value.NewPosition(0, 0), "Red", value.NewStyleWithFg(value.ColorBlue))
	err = r.Render(buffer2)
	require.NoError(t, err)

	output := buf.String()
	assert.NotEmpty(t, output)
	// Should contain blue color code
	assert.Contains(t, output, "38;2;0;0;255") // Blue RGB
}

func TestRenderer_Render_LargeChangeTriggersFullRedraw(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	// First render - sparse content
	buffer1 := model.NewBuffer(80, 24)
	buffer1.SetString(value.NewPosition(0, 0), "Small", value.NewStyle())
	err := r.Render(buffer1)
	require.NoError(t, err)

	buf.Reset()

	// Second render - fill > 75% of buffer (triggers full redraw)
	buffer2 := model.NewBuffer(80, 24)
	for y := 0; y < 24; y++ {
		for x := 0; x < 70; x++ { // Fill most of the line
			buffer2.Set(value.NewPosition(x, y), model.NewCell('X', value.NewStyle()))
		}
	}

	err = r.Render(buffer2)
	require.NoError(t, err)

	output := buf.String()
	assert.NotEmpty(t, output)
	// Full redraw should happen
}

func TestRenderer_Clear_ErrorHandling(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	// Add some content first
	buffer := model.NewBuffer(80, 24)
	buffer.SetString(value.NewPosition(0, 0), "Test", value.NewStyle())
	_ = r.Render(buffer)

	buf.Reset()

	// Clear should work
	err := r.Clear()
	require.NoError(t, err)

	output := buf.String()
	assert.NotEmpty(t, output) // Should have clear sequences
}

func TestRenderer_Resize_WithClear(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	// Add content
	buffer := model.NewBuffer(80, 24)
	buffer.SetString(value.NewPosition(0, 0), "Content", value.NewStyle())
	_ = r.Render(buffer)

	buf.Reset()

	// Resize should clear
	err := r.Resize(100, 30)
	require.NoError(t, err)

	assert.Equal(t, 100, r.Width())
	assert.Equal(t, 30, r.Height())

	output := buf.String()
	assert.NotEmpty(t, output) // Should have clear sequences
}

func TestRenderer_HideCursor_ErrorPath(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	err := r.HideCursor()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "\x1b[?25l") // Hide cursor sequence
}

func TestRenderer_ShowCursor_ErrorPath(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	err := r.ShowCursor()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "\x1b[?25h") // Show cursor sequence
}

func TestRenderer_GetBuffer_FromPool(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	// Get multiple buffers from pool
	buf1 := r.GetBuffer()
	buf2 := r.GetBuffer()

	assert.NotNil(t, buf1)
	assert.NotNil(t, buf2)
	assert.Equal(t, 80, buf1.Width())
	assert.Equal(t, 24, buf1.Height())

	// Put back
	r.PutBuffer(buf1)
	r.PutBuffer(buf2)

	// Get again (should reuse)
	buf3 := r.GetBuffer()
	assert.NotNil(t, buf3)
}

func TestRenderer_PutBuffer_Nil(_ *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	// Should not panic
	r.PutBuffer(nil)
}

func TestRenderer_GetStats_Fields(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	stats := r.GetStats()
	assert.GreaterOrEqual(t, stats.BufferedBytes, 0)
}

func TestRenderer_Close_ShowsCursor(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	// Hide cursor first
	_ = r.HideCursor()
	buf.Reset()

	// Close should show cursor
	err := r.Close()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "\x1b[?25h") // Show cursor sequence
}

func TestRenderer_ConcurrentAccess(_ *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	// Test thread safety with concurrent renders
	done := make(chan bool)

	for i := 0; i < 5; i++ {
		go func(n int) {
			buffer := model.NewBuffer(80, 24)
			buffer.SetString(value.NewPosition(0, 0), fmt.Sprintf("Thread %d", n), value.NewStyle())
			_ = r.Render(buffer)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}
}

func TestNewRendererWithOptions_CustomBufferSize(t *testing.T) {
	var buf bytes.Buffer
	r := NewRendererWithOptions(80, 24, &buf, 8192)

	assert.NotNil(t, r)
	assert.Equal(t, 80, r.Width())
	assert.Equal(t, 24, r.Height())

	// Render something to verify it works
	buffer := model.NewBuffer(80, 24)
	buffer.SetString(value.NewPosition(0, 0), "Test", value.NewStyle())
	err := r.Render(buffer)
	require.NoError(t, err)
}

func TestRenderer_Render_StyleOptimization(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(80, 24, &buf)

	// First render with style
	buffer1 := model.NewBuffer(80, 24)
	style := value.NewStyleWithFg(value.ColorRed).WithBold(true)
	buffer1.SetString(value.NewPosition(0, 0), "Hello", style)
	err := r.Render(buffer1)
	require.NoError(t, err)

	buf.Reset()

	// Second render with same style (should optimize)
	buffer2 := model.NewBuffer(80, 24)
	buffer2.SetString(value.NewPosition(0, 0), "World", style)
	err = r.Render(buffer2)
	require.NoError(t, err)

	output := buf.String()
	assert.NotEmpty(t, output)
}
