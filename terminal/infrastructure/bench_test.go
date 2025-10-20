package infrastructure

import (
	"testing"

	"github.com/phoenix-tui/phoenix/terminal/infrastructure/unix"
)

// BenchmarkANSI_SetCursorPosition benchmarks ANSI cursor positioning.
func BenchmarkANSI_SetCursorPosition(b *testing.B) {
	term := unix.NewANSI()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = term.SetCursorPosition(10, 5)
	}
}

// BenchmarkANSI_MoveCursorUp benchmarks ANSI cursor movement.
func BenchmarkANSI_MoveCursorUp(b *testing.B) {
	term := unix.NewANSI()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = term.MoveCursorUp(5)
	}
}

// BenchmarkANSI_ClearLine benchmarks ANSI line clearing.
func BenchmarkANSI_ClearLine(b *testing.B) {
	term := unix.NewANSI()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = term.ClearLine()
	}
}

// BenchmarkANSI_ClearLines benchmarks ANSI multiline clearing.
func BenchmarkANSI_ClearLines(b *testing.B) {
	term := unix.NewANSI()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = term.ClearLines(10)
	}
}

// BenchmarkANSI_Write benchmarks ANSI text output.
func BenchmarkANSI_Write(b *testing.B) {
	term := unix.NewANSI()
	text := "Phoenix Terminal Framework - Next-gen TUI for Go"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = term.Write(text)
	}
}

// BenchmarkANSI_WriteAt benchmarks ANSI positioned text output.
func BenchmarkANSI_WriteAt(b *testing.B) {
	term := unix.NewANSI()
	text := "Phoenix Terminal Framework - Next-gen TUI for Go"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = term.WriteAt(10, 5, text)
	}
}

// BenchmarkANSI_HideShowCursor benchmarks cursor visibility toggling.
func BenchmarkANSI_HideShowCursor(b *testing.B) {
	term := unix.NewANSI()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = term.HideCursor()
		_ = term.ShowCursor()
	}
}

// BenchmarkANSI_SaveRestoreCursorPosition benchmarks cursor save/restore.
func BenchmarkANSI_SaveRestoreCursorPosition(b *testing.B) {
	term := unix.NewANSI()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = term.SaveCursorPosition()
		_ = term.RestoreCursorPosition()
	}
}

// Windows-specific benchmarks (only run on Windows).

// The following benchmarks are implemented but will only run on Windows builds.
// On Unix systems, they'll be skipped with a clear message.
//
// To compare performance on Windows:.
//   go test -bench=. -benchmem ./infrastructure.
//
// Expected results (Windows Console API vs ANSI):.
//   - SetCursorPosition: 10x faster.
//   - ClearLines: 10x faster.
//   - GetCursorPosition: Only Windows supports this (ANSI fails).
//   - ReadScreenBuffer: Only Windows supports this (ANSI fails).

// Note: These benchmarks require the windows package to be built.
// They are implemented in bench_windows_test.go with build tags.
