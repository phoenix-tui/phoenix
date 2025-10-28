//go:build windows
// +build windows

package windows

import (
	"testing"

	"github.com/phoenix-tui/phoenix/terminal/types"
)

// TestNewConsole verifies Console creation.
func TestNewConsole(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		// This is expected in Git Bash or redirected I/O.
		t.Skipf("Not running in Windows Console: %v", err)
	}

	if console == nil {
		t.Fatal("NewConsole() returned nil console")
	}

	// Verify handles are valid.
	if console.stdout == 0 {
		t.Error("stdout handle is invalid (0)")
	}
	if console.stdin == 0 {
		t.Error("stdin handle is invalid (0)")
	}
}

// TestConsole_Platform verifies platform detection.
func TestConsole_Platform(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	platform := console.Platform()
	if platform != types.PlatformWindowsConsole {
		t.Errorf("Platform() = %v, want %v", platform, types.PlatformWindowsConsole)
	}
}

// TestConsole_Size verifies terminal size detection.
func TestConsole_Size(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	width, height, err := console.Size()
	if err != nil {
		t.Fatalf("Size() error: %v", err)
	}

	// Reasonable terminal size bounds.
	if width < 20 || width > 1000 {
		t.Errorf("Width = %d, expected 20-1000", width)
	}
	if height < 10 || height > 500 {
		t.Errorf("Height = %d, expected 10-500", height)
	}
}

// TestConsole_SetCursorPosition verifies absolute cursor positioning.
func TestConsole_SetCursorPosition(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Set cursor to position (10, 5).
	err = console.SetCursorPosition(10, 5)
	if err != nil {
		t.Fatalf("SetCursorPosition(10, 5) error: %v", err)
	}

	// Verify position was set correctly.
	x, y, err := console.GetCursorPosition()
	if err != nil {
		t.Fatalf("GetCursorPosition() error: %v", err)
	}

	if x != 10 || y != 5 {
		t.Errorf("GetCursorPosition() = (%d, %d), want (10, 5)", x, y)
	}
}

// TestConsole_GetCursorPosition verifies cursor readback.
func TestConsole_GetCursorPosition(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	x, y, err := console.GetCursorPosition()
	if err != nil {
		t.Fatalf("GetCursorPosition() error: %v", err)
	}

	// Cursor position should be within screen bounds.
	width, height, _ := console.Size()
	if x < 0 || x >= width {
		t.Errorf("Cursor X = %d, expected 0-%d", x, width-1)
	}
	if y < 0 || y >= height {
		t.Errorf("Cursor Y = %d, expected 0-%d", y, height-1)
	}
}

// TestConsole_MoveCursor verifies relative cursor movements.
func TestConsole_MoveCursor(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Set known position.
	err = console.SetCursorPosition(20, 10)
	if err != nil {
		t.Fatalf("SetCursorPosition() error: %v", err)
	}

	// Move right 5.
	err = console.MoveCursorRight(5)
	if err != nil {
		t.Fatalf("MoveCursorRight(5) error: %v", err)
	}

	x, y, _ := console.GetCursorPosition()
	if x != 25 {
		t.Errorf("After MoveCursorRight(5), X = %d, want 25", x)
	}
	if y != 10 {
		t.Errorf("After MoveCursorRight(5), Y = %d, want 10", y)
	}

	// Move down 3.
	err = console.MoveCursorDown(3)
	if err != nil {
		t.Fatalf("MoveCursorDown(3) error: %v", err)
	}

	x, y, _ = console.GetCursorPosition()
	if x != 25 {
		t.Errorf("After MoveCursorDown(3), X = %d, want 25", x)
	}
	if y != 13 {
		t.Errorf("After MoveCursorDown(3), Y = %d, want 13", y)
	}

	// Move left 10.
	err = console.MoveCursorLeft(10)
	if err != nil {
		t.Fatalf("MoveCursorLeft(10) error: %v", err)
	}

	x, y, _ = console.GetCursorPosition()
	if x != 15 {
		t.Errorf("After MoveCursorLeft(10), X = %d, want 15", x)
	}

	// Move up 5.
	err = console.MoveCursorUp(5)
	if err != nil {
		t.Fatalf("MoveCursorUp(5) error: %v", err)
	}

	x, y, _ = console.GetCursorPosition()
	if y != 8 {
		t.Errorf("After MoveCursorUp(5), Y = %d, want 8", y)
	}
}

// TestConsole_SaveRestoreCursorPosition verifies cursor save/restore.
func TestConsole_SaveRestoreCursorPosition(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Set known position.
	err = console.SetCursorPosition(30, 15)
	if err != nil {
		t.Fatalf("SetCursorPosition() error: %v", err)
	}

	// Save position.
	err = console.SaveCursorPosition()
	if err != nil {
		t.Fatalf("SaveCursorPosition() error: %v", err)
	}

	// Move to different position.
	err = console.SetCursorPosition(5, 5)
	if err != nil {
		t.Fatalf("SetCursorPosition(5, 5) error: %v", err)
	}

	// Restore position.
	err = console.RestoreCursorPosition()
	if err != nil {
		t.Fatalf("RestoreCursorPosition() error: %v", err)
	}

	// Verify we're back at saved position.
	x, y, _ := console.GetCursorPosition()
	if x != 30 || y != 15 {
		t.Errorf("After restore, position = (%d, %d), want (30, 15)", x, y)
	}
}

// TestConsole_HideShowCursor verifies cursor visibility control.
func TestConsole_HideShowCursor(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Hide cursor.
	err = console.HideCursor()
	if err != nil {
		t.Fatalf("HideCursor() error: %v", err)
	}

	// Show cursor (restore visibility).
	err = console.ShowCursor()
	if err != nil {
		t.Fatalf("ShowCursor() error: %v", err)
	}

	// Note: We can't easily verify visibility state without additional Win32 API calls.
	// The fact that no errors occurred is a good sign.
}

// TestConsole_SetCursorStyle verifies cursor style changes.
func TestConsole_SetCursorStyle(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	styles := []types.CursorStyle{
		types.CursorBlock,
		types.CursorUnderline,
		types.CursorBar,
	}

	for _, style := range styles {
		err = console.SetCursorStyle(style)
		if err != nil {
			t.Errorf("SetCursorStyle(%v) error: %v", style, err)
		}
	}

	// Restore to block (default).
	err = console.SetCursorStyle(types.CursorBlock)
	if err != nil {
		t.Fatalf("SetCursorStyle(Block) error: %v", err)
	}
}

// TestConsole_Write verifies text output.
func TestConsole_Write(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Save cursor position.
	err = console.SaveCursorPosition()
	if err != nil {
		t.Fatalf("SaveCursorPosition() error: %v", err)
	}

	// Write some text.
	err = console.Write("Phoenix Terminal Test")
	if err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	// Restore position (cleanup).
	err = console.RestoreCursorPosition()
	if err != nil {
		t.Fatalf("RestoreCursorPosition() error: %v", err)
	}
}

// TestConsole_WriteAt verifies positioned text output.
func TestConsole_WriteAt(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Write at specific position.
	err = console.WriteAt(10, 5, "Test")
	if err != nil {
		t.Fatalf("WriteAt(10, 5, \"Test\") error: %v", err)
	}

	// Verify cursor moved to that position (after text).
	x, y, _ := console.GetCursorPosition()
	if y != 5 {
		t.Errorf("After WriteAt, Y = %d, want 5", y)
	}
	// X will be 10 + len("Test") = 14.
	if x != 14 {
		t.Errorf("After WriteAt, X = %d, want 14", x)
	}
}

// TestConsole_ClearLine verifies line clearing.
func TestConsole_ClearLine(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Write some text.
	err = console.Write("This line will be cleared")
	if err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	// Clear the line.
	err = console.ClearLine()
	if err != nil {
		t.Fatalf("ClearLine() error: %v", err)
	}
}

// TestConsole_ClearFromCursor verifies clearing from cursor to end.
func TestConsole_ClearFromCursor(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Set cursor position.
	err = console.SetCursorPosition(0, 10)
	if err != nil {
		t.Fatalf("SetCursorPosition() error: %v", err)
	}

	// Clear from cursor to end of screen.
	err = console.ClearFromCursor()
	if err != nil {
		t.Fatalf("ClearFromCursor() error: %v", err)
	}
}

// TestConsole_ClearLines verifies multiline clearing (CRITICAL for GoSh).
func TestConsole_ClearLines(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Set cursor to known position.
	err = console.SetCursorPosition(0, 15)
	if err != nil {
		t.Fatalf("SetCursorPosition() error: %v", err)
	}

	// Clear 5 lines (should move cursor up and clear).
	err = console.ClearLines(5)
	if err != nil {
		t.Fatalf("ClearLines(5) error: %v", err)
	}

	// Verify cursor is at start of cleared region.
	x, y, _ := console.GetCursorPosition()
	if x != 0 {
		t.Errorf("After ClearLines, X = %d, want 0", x)
	}
	if y != 11 { // 15 - 5 + 1 = 11
		t.Errorf("After ClearLines(5), Y = %d, want 11", y)
	}
}

// TestConsole_Clear verifies full screen clearing.
func TestConsole_Clear(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Clear entire screen.
	err = console.Clear()
	if err != nil {
		t.Fatalf("Clear() error: %v", err)
	}

	// Verify cursor is at top-left.
	x, y, _ := console.GetCursorPosition()
	if x != 0 || y != 0 {
		t.Errorf("After Clear(), position = (%d, %d), want (0, 0)", x, y)
	}
}

// TestConsole_ReadScreenBuffer verifies screen buffer readback.
func TestConsole_ReadScreenBuffer(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Read screen buffer.
	buffer, err := console.ReadScreenBuffer()
	if err != nil {
		t.Fatalf("ReadScreenBuffer() error: %v", err)
	}

	// Verify buffer dimensions match terminal size.
	width, height, _ := console.Size()
	if len(buffer) != height {
		t.Errorf("Buffer height = %d, want %d", len(buffer), height)
	}

	if len(buffer) > 0 && len(buffer[0]) != width {
		t.Errorf("Buffer width = %d, want %d", len(buffer[0]), width)
	}
}

// TestConsole_ColorDepth verifies color support detection.
func TestConsole_ColorDepth(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	depth := console.ColorDepth()
	// Windows 10+ should support TrueColor.
	if depth != 16777216 {
		t.Errorf("ColorDepth() = %d, want 16777216 (TrueColor)", depth)
	}
}

// TestConsole_Capabilities verifies capability flags.
func TestConsole_Capabilities(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Windows Console should support all capabilities.
	if !console.SupportsDirectPositioning() {
		t.Error("SupportsDirectPositioning() = false, want true")
	}

	if !console.SupportsReadback() {
		t.Error("SupportsReadback() = false, want true")
	}

	if !console.SupportsTrueColor() {
		t.Error("SupportsTrueColor() = false, want true")
	}
}

// TestConsole_EdgeCases verifies edge case handling.
func TestConsole_EdgeCases(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Test zero movements (should be no-ops).
	tests := []struct {
		name string
		fn   func() error
	}{
		{"MoveCursorUp(0)", func() error { return console.MoveCursorUp(0) }},
		{"MoveCursorDown(0)", func() error { return console.MoveCursorDown(0) }},
		{"MoveCursorLeft(0)", func() error { return console.MoveCursorLeft(0) }},
		{"MoveCursorRight(0)", func() error { return console.MoveCursorRight(0) }},
		{"ClearLines(0)", func() error { return console.ClearLines(0) }},
		{"ClearLines(-1)", func() error { return console.ClearLines(-1) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); err != nil {
				t.Errorf("%s error: %v", tt.name, err)
			}
		})
	}
}

// TestConsole_BoundaryConditions verifies boundary handling.
func TestConsole_BoundaryConditions(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Try to move cursor beyond boundaries.
	width, height, _ := console.Size()

	// Set to top-left.
	console.SetCursorPosition(0, 0)

	// Try to move up from top (should stay at 0).
	err = console.MoveCursorUp(10)
	if err != nil {
		t.Fatalf("MoveCursorUp(10) error: %v", err)
	}
	_, y, _ := console.GetCursorPosition()
	if y != 0 {
		t.Errorf("After MoveCursorUp from top, Y = %d, want 0", y)
	}

	// Set to bottom-right.
	console.SetCursorPosition(width-1, height-1)

	// Try to move down from bottom (should stay at max).
	err = console.MoveCursorDown(10)
	if err != nil {
		t.Fatalf("MoveCursorDown(10) error: %v", err)
	}
	_, y, _ = console.GetCursorPosition()
	if y != height-1 {
		t.Errorf("After MoveCursorDown from bottom, Y = %d, want %d", y, height-1)
	}
}
