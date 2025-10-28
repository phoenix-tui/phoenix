// Package benchmarks provides comprehensive performance benchmarks for phoenix/render.
//
// Target: 60 FPS (< 16.67ms per frame)
// Target: 10x faster than Charm ecosystem
//
// Run benchmarks:
//
//	go test -bench=. -benchmem -benchtime=10s ./benchmarks
package benchmarks

import (
	"bytes"
	"testing"

	"github.com/phoenix-tui/phoenix/render"
)

// Benchmark 60 FPS target: Full screen render < 16.67ms
func BenchmarkFullScreen_60FPS(b *testing.B) {
	var buf bytes.Buffer
	r := render.New(80, 24, &buf)
	defer r.Close()

	buffer := r.Buffer()
	buffer.Fill('X', render.StyleDefault())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Render(buffer)
	}

	// Report ns/op for FPS calculation
	b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N), "ns/op")
	fps := 1000000000.0 / (float64(b.Elapsed().Nanoseconds()) / float64(b.N))
	b.ReportMetric(fps, "fps")
}

// Benchmark differential rendering (typical case)
func BenchmarkDifferential_SmallChange(b *testing.B) {
	var buf bytes.Buffer
	r := render.New(80, 24, &buf)
	defer r.Close()

	// Initial render
	buffer1 := r.Buffer()
	buffer1.SetString(0, 0, "Hello, World!", render.StyleDefault())
	_ = r.Render(buffer1)

	// Second buffer with small change
	buffer2 := r.Buffer()
	buffer2.SetString(0, 0, "Hello, Phoenix!", render.StyleDefault())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Render(buffer2)
	}
}

// Benchmark differential rendering (10% change - typical TUI scenario)
func BenchmarkDifferential_10Percent(b *testing.B) {
	var buf bytes.Buffer
	r := render.New(80, 24, &buf)
	defer r.Close()

	// Initial render
	buffer1 := r.Buffer()
	for y := 0; y < 24; y++ {
		buffer1.SetLine(y, "Line content here", render.StyleDefault())
	}
	_ = r.Render(buffer1)

	// Buffer with 10% change (~2-3 lines)
	buffer2 := r.Buffer()
	for y := 0; y < 24; y++ {
		if y < 3 {
			buffer2.SetLine(y, "CHANGED LINE HERE", render.FgRed)
		} else {
			buffer2.SetLine(y, "Line content here", render.StyleDefault())
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Render(buffer2)
	}
}

// Benchmark large screen (modern terminals)
func BenchmarkLargeScreen_200x50(b *testing.B) {
	var buf bytes.Buffer
	r := render.New(200, 50, &buf)
	defer r.Close()

	buffer := r.Buffer()
	buffer.Fill('X', render.StyleDefault())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Render(buffer)
	}
}

// Benchmark Unicode rendering (emoji, CJK)
func BenchmarkUnicode_Emoji(b *testing.B) {
	var buf bytes.Buffer
	r := render.New(80, 24, &buf)
	defer r.Close()

	buffer := r.Buffer()
	text := "😀😃😄😁😆😅🤣😂" // Emoji (2 cells each)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.SetString(0, 0, text, render.StyleDefault())
		_ = r.Render(buffer)
	}
}

// Benchmark Unicode rendering (CJK characters)
func BenchmarkUnicode_CJK(b *testing.B) {
	var buf bytes.Buffer
	r := render.New(80, 24, &buf)
	defer r.Close()

	buffer := r.Buffer()
	text := "日本語の文字列テスト" // Japanese text

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.SetString(0, 0, text, render.StyleDefault())
		_ = r.Render(buffer)
	}
}

// Benchmark styled text rendering
func BenchmarkStyled_MultiColor(b *testing.B) {
	var buf bytes.Buffer
	r := render.New(80, 24, &buf)
	defer r.Close()

	buffer := r.Buffer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.SetString(0, 0, "Red", render.FgRed)
		buffer.SetString(4, 0, "Green", render.FgGreen)
		buffer.SetString(10, 0, "Blue", render.FgBlue)
		buffer.SetString(15, 0, "Yellow", render.FgYellow)
		_ = r.Render(buffer)
	}
}

// Benchmark buffer operations
func BenchmarkBuffer_SetString_ASCII(b *testing.B) {
	buffer := render.NewBuffer(80, 24)
	text := "The quick brown fox jumps over the lazy dog"
	style := render.StyleDefault()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.SetString(0, 0, text, style)
	}
}

func BenchmarkBuffer_SetString_Unicode(b *testing.B) {
	buffer := render.NewBuffer(80, 24)
	text := "日本語の文字列テスト😀😃😄" // Mixed CJK + emoji
	style := render.StyleDefault()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.SetString(0, 0, text, style)
	}
}

func BenchmarkBuffer_Fill(b *testing.B) {
	buffer := render.NewBuffer(80, 24)
	style := render.StyleDefault()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Fill('X', style)
	}
}

func BenchmarkBuffer_Clear(b *testing.B) {
	buffer := render.NewBuffer(80, 24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Clear()
	}
}

// Benchmark style operations
func BenchmarkStyle_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = render.StyleFg(255, 128, 64)
	}
}

func BenchmarkStyle_Chaining(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = render.StyleDefault().
			WithFg(255, 0, 0).
			WithBg(0, 0, 255).
			WithBold(true).
			WithItalic(true)
	}
}

// Benchmark memory allocations
func BenchmarkMemory_BufferPooling(b *testing.B) {
	var buf bytes.Buffer
	r := render.New(80, 24, &buf)
	defer r.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer := r.Buffer()
		buffer.SetString(0, 0, "Test", render.StyleDefault())
		_ = r.Render(buffer)
		buffer.Release()
	}
}

// Benchmark real-world scenario: scrolling terminal
func BenchmarkRealWorld_ScrollingTerminal(b *testing.B) {
	var buf bytes.Buffer
	r := render.New(80, 24, &buf)
	defer r.Close()

	// Simulate scrolling by shifting lines
	lines := make([]string, 100)
	for i := range lines {
		lines[i] = "This is line content with some text here"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer := r.Buffer()

		// Render visible window (simulating scroll)
		offset := i % (len(lines) - 24)
		for y := 0; y < 24; y++ {
			buffer.SetLine(y, lines[offset+y], render.StyleDefault())
		}

		_ = r.Render(buffer)
		buffer.Release()
	}
}

// Benchmark real-world scenario: code editor
func BenchmarkRealWorld_CodeEditor(b *testing.B) {
	var buf bytes.Buffer
	r := render.New(120, 40, &buf)
	defer r.Close()

	// Simulate code with syntax highlighting
	codeLines := []struct {
		text  string
		style render.Style
	}{
		{"package main", render.FgMagenta.WithBold(true)},
		{"import \"fmt\"", render.FgGreen},
		{"func main() {", render.FgBlue},
		{"    fmt.Println(\"Hello\")", render.StyleDefault()},
		{"}", render.FgBlue},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer := r.Buffer()

		for y, line := range codeLines {
			if y < 40 {
				buffer.SetString(0, y, line.text, line.style)
			}
		}

		_ = r.Render(buffer)
		buffer.Release()
	}
}

// Benchmark worst case: full screen change every frame
func BenchmarkWorstCase_FullScreenChange(b *testing.B) {
	var buf bytes.Buffer
	r := render.New(80, 24, &buf)
	defer r.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer := r.Buffer()
		// Fill with different content each time
		buffer.Fill(rune('A'+(i%26)), render.StyleDefault())
		_ = r.Render(buffer)
		buffer.Release()
	}
}

// Benchmark best case: no changes
func BenchmarkBestCase_NoChanges(b *testing.B) {
	var buf bytes.Buffer
	r := render.New(80, 24, &buf)
	defer r.Close()

	buffer := r.Buffer()
	buffer.SetString(0, 0, "Static content", render.StyleDefault())
	_ = r.Render(buffer) // Initial render

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Render(buffer) // Should be near-zero cost
	}
}

// Sub-benchmarks for different screen sizes
func BenchmarkScreenSizes(b *testing.B) {
	sizes := []struct {
		name          string
		width, height int
	}{
		{"Small_40x12", 40, 12},
		{"Standard_80x24", 80, 24},
		{"Large_120x40", 120, 40},
		{"XLarge_200x60", 200, 60},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			var buf bytes.Buffer
			r := render.New(size.width, size.height, &buf)
			defer r.Close()

			buffer := r.Buffer()
			buffer.Fill('X', render.StyleDefault())

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = r.Render(buffer)
			}
		})
	}
}
