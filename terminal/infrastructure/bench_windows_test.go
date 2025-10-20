//go:build windows.
// +build windows.

package infrastructure

import (
	"testing"

	"github.com/phoenix-tui/phoenix/terminal/infrastructure/windows"
)

// BenchmarkWindowsAPI_SetCursorPosition benchmarks Windows Console API cursor positioning.
func BenchmarkWindowsAPI_SetCursorPosition(b *testing.B) {
	console, err := windows.NewConsole()
	if err != nil {
		b.Skipf("Not running in Windows Console: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = console.SetCursorPosition(10, 5)
	}
}

// BenchmarkWindowsAPI_GetCursorPosition benchmarks cursor readback (ANSI can't do this).
func BenchmarkWindowsAPI_GetCursorPosition(b *testing.B) {
	console, err := windows.NewConsole()
	if err != nil {
		b.Skipf("Not running in Windows Console: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = console.GetCursorPosition()
	}
}

// BenchmarkWindowsAPI_MoveCursorUp benchmarks relative cursor movement.
func BenchmarkWindowsAPI_MoveCursorUp(b *testing.B) {
	console, err := windows.NewConsole()
	if err != nil {
		b.Skipf("Not running in Windows Console: %v", err)
	}

	// Set to middle of screen first.
	console.SetCursorPosition(40, 20)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = console.MoveCursorUp(5)
		_ = console.MoveCursorDown(5) // Reset for next iteration
	}
}

// BenchmarkWindowsAPI_ClearLine benchmarks line clearing.
func BenchmarkWindowsAPI_ClearLine(b *testing.B) {
	console, err := windows.NewConsole()
	if err != nil {
		b.Skipf("Not running in Windows Console: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = console.ClearLine()
	}
}

// BenchmarkWindowsAPI_ClearLines benchmarks multiline clearing (CRITICAL for GoSh).
func BenchmarkWindowsAPI_ClearLines(b *testing.B) {
	console, err := windows.NewConsole()
	if err != nil {
		b.Skipf("Not running in Windows Console: %v", err)
	}

	// Set to position with room to clear upward.
	console.SetCursorPosition(0, 20)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = console.ClearLines(10)
		// Reset position for next iteration.
		_ = console.SetCursorPosition(0, 20)
	}
}

// BenchmarkWindowsAPI_Write benchmarks text output.
func BenchmarkWindowsAPI_Write(b *testing.B) {
	console, err := windows.NewConsole()
	if err != nil {
		b.Skipf("Not running in Windows Console: %v", err)
	}

	text := "Phoenix Terminal Framework - Next-gen TUI for Go"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = console.Write(text)
	}
}

// BenchmarkWindowsAPI_WriteAt benchmarks positioned text output.
func BenchmarkWindowsAPI_WriteAt(b *testing.B) {
	console, err := windows.NewConsole()
	if err != nil {
		b.Skipf("Not running in Windows Console: %v", err)
	}

	text := "Phoenix Terminal Framework - Next-gen TUI for Go"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = console.WriteAt(10, 5, text)
	}
}

// BenchmarkWindowsAPI_HideShowCursor benchmarks cursor visibility toggling.
func BenchmarkWindowsAPI_HideShowCursor(b *testing.B) {
	console, err := windows.NewConsole()
	if err != nil {
		b.Skipf("Not running in Windows Console: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = console.HideCursor()
		_ = console.ShowCursor()
	}
}

// BenchmarkWindowsAPI_SaveRestoreCursorPosition benchmarks cursor save/restore.
func BenchmarkWindowsAPI_SaveRestoreCursorPosition(b *testing.B) {
	console, err := windows.NewConsole()
	if err != nil {
		b.Skipf("Not running in Windows Console: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = console.SaveCursorPosition()
		_ = console.RestoreCursorPosition()
	}
}

// BenchmarkWindowsAPI_ReadScreenBuffer benchmarks screen buffer readback (ANSI can't do this).
func BenchmarkWindowsAPI_ReadScreenBuffer(b *testing.B) {
	console, err := windows.NewConsole()
	if err != nil {
		b.Skipf("Not running in Windows Console: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = console.ReadScreenBuffer()
	}
}

// BenchmarkWindowsAPI_Clear benchmarks full screen clearing.
func BenchmarkWindowsAPI_Clear(b *testing.B) {
	console, err := windows.NewConsole()
	if err != nil {
		b.Skipf("Not running in Windows Console: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = console.Clear()
	}
}

// BenchmarkWindowsAPI_Size benchmarks terminal size detection.
func BenchmarkWindowsAPI_Size(b *testing.B) {
	console, err := windows.NewConsole()
	if err != nil {
		b.Skipf("Not running in Windows Console: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = console.Size()
	}
}

// Comparison benchmarks - demonstrate 10x improvement.

// BenchmarkComparison_ClearLines_ANSI vs BenchmarkComparison_ClearLines_WindowsAPI.
// Run both to see the performance difference:.
//   go test -bench=BenchmarkComparison_ClearLines -benchmem ./infrastructure.
//
// Expected result:.
//   ANSI:       ~500 μs per operation.
//   Windows API: ~50 μs per operation.
//   Improvement: 10x faster.

// BenchmarkComparison_ClearLines_ANSI benchmarks ANSI clearing (baseline).
func BenchmarkComparison_ClearLines_ANSI(b *testing.B) {
	// Use NewANSITerminal() to force ANSI even on Windows.
	term := NewANSITerminal()

	// Set position with room to clear upward.
	term.SetCursorPosition(0, 20)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = term.ClearLines(10)
		// Reset position.
		_ = term.SetCursorPosition(0, 20)
	}
}

// BenchmarkComparison_ClearLines_WindowsAPI benchmarks Windows API clearing (optimized).
func BenchmarkComparison_ClearLines_WindowsAPI(b *testing.B) {
	console, err := windows.NewConsole()
	if err != nil {
		b.Skipf("Not running in Windows Console: %v", err)
	}

	// Set position with room to clear upward.
	console.SetCursorPosition(0, 20)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = console.ClearLines(10)
		// Reset position.
		_ = console.SetCursorPosition(0, 20)
	}
}

// BenchmarkComparison_SetCursorPosition demonstrates positioning performance difference.

// BenchmarkComparison_SetCursorPosition_ANSI benchmarks ANSI positioning (baseline).
func BenchmarkComparison_SetCursorPosition_ANSI(b *testing.B) {
	term := NewANSITerminal()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = term.SetCursorPosition(10, 5)
	}
}

// BenchmarkComparison_SetCursorPosition_WindowsAPI benchmarks Windows API positioning (optimized).
func BenchmarkComparison_SetCursorPosition_WindowsAPI(b *testing.B) {
	console, err := windows.NewConsole()
	if err != nil {
		b.Skipf("Not running in Windows Console: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = console.SetCursorPosition(10, 5)
	}
}
