package unix

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/terminal/types"
)

// Note: testWriter and createTestTerminal removed (unused helper functions).
// Tests now use direct ANSITerminal creation with os.Stdout.

// captureANSI executes function and captures ANSI output.
func captureANSI(fn func(*ANSITerminal)) string {
	r, w, _ := os.Pipe()
	term := &ANSITerminal{output: w}

	fn(term)
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	return buf.String()
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Cursor Position Tests                                           │.
// └─────────────────────────────────────────────────────────────────┘.

func TestANSI_SetCursorPosition(t *testing.T) {
	tests := []struct {
		name string
		x, y int
		want string
	}{
		{"origin", 0, 0, "\033[1;1H"},
		{"middle", 10, 5, "\033[6;11H"},
		{"large", 100, 50, "\033[51;101H"},
		{"single_digit", 5, 3, "\033[4;6H"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureANSI(func(term *ANSITerminal) {
				term.SetCursorPosition(tt.x, tt.y)
			})

			if got != tt.want {
				t.Errorf("SetCursorPosition(%d, %d) = %q, want %q", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestANSI_GetCursorPosition_NotSupported(t *testing.T) {
	term := NewANSI()
	x, y, err := term.GetCursorPosition()

	if err == nil {
		t.Error("GetCursorPosition should return error on ANSI")
	}

	if x != 0 || y != 0 {
		t.Errorf("GetCursorPosition on error should return 0,0, got %d,%d", x, y)
	}

	if !strings.Contains(err.Error(), "readback") {
		t.Errorf("Error should mention readback, got: %v", err)
	}
}

func TestANSI_MoveCursorUp(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{"one", 1, "\033[1A"},
		{"ten", 10, "\033[10A"},
		{"zero", 0, ""},      // No-op
		{"negative", -5, ""}, // No-op
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureANSI(func(term *ANSITerminal) {
				term.MoveCursorUp(tt.n)
			})

			if got != tt.want {
				t.Errorf("MoveCursorUp(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}

func TestANSI_MoveCursorDown(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{"one", 1, "\033[1B"},
		{"five", 5, "\033[5B"},
		{"zero", 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureANSI(func(term *ANSITerminal) {
				term.MoveCursorDown(tt.n)
			})

			if got != tt.want {
				t.Errorf("MoveCursorDown(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}

func TestANSI_MoveCursorLeft(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{"one", 1, "\033[1D"},
		{"twenty", 20, "\033[20D"},
		{"zero", 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureANSI(func(term *ANSITerminal) {
				term.MoveCursorLeft(tt.n)
			})

			if got != tt.want {
				t.Errorf("MoveCursorLeft(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}

func TestANSI_MoveCursorRight(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{"one", 1, "\033[1C"},
		{"fifteen", 15, "\033[15C"},
		{"zero", 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureANSI(func(term *ANSITerminal) {
				term.MoveCursorRight(tt.n)
			})

			if got != tt.want {
				t.Errorf("MoveCursorRight(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}

func TestANSI_SaveRestoreCursorPosition(t *testing.T) {
	// Test Save.
	got := captureANSI(func(term *ANSITerminal) {
		term.SaveCursorPosition()
	})
	if got != "\033[s" {
		t.Errorf("SaveCursorPosition = %q, want %q", got, "\033[s")
	}

	// Test Restore.
	got = captureANSI(func(term *ANSITerminal) {
		term.RestoreCursorPosition()
	})
	if got != "\033[u" {
		t.Errorf("RestoreCursorPosition = %q, want %q", got, "\033[u")
	}
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Cursor Visibility Tests                                         │.
// └─────────────────────────────────────────────────────────────────┘.

func TestANSI_HideCursor(t *testing.T) {
	got := captureANSI(func(term *ANSITerminal) {
		term.HideCursor()
	})

	want := "\033[?25l"
	if got != want {
		t.Errorf("HideCursor = %q, want %q", got, want)
	}
}

func TestANSI_ShowCursor(t *testing.T) {
	got := captureANSI(func(term *ANSITerminal) {
		term.ShowCursor()
	})

	want := "\033[?25h"
	if got != want {
		t.Errorf("ShowCursor = %q, want %q", got, want)
	}
}

func TestANSI_SetCursorStyle(t *testing.T) {
	tests := []struct {
		name  string
		style types.CursorStyle
		want  string
	}{
		{"block", types.CursorBlock, "\033[2 q"},
		{"underline", types.CursorUnderline, "\033[4 q"},
		{"bar", types.CursorBar, "\033[6 q"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureANSI(func(term *ANSITerminal) {
				term.SetCursorStyle(tt.style)
			})

			if got != tt.want {
				t.Errorf("SetCursorStyle(%v) = %q, want %q", tt.style, got, tt.want)
			}
		})
	}
}

func TestANSI_SetCursorStyle_Invalid(t *testing.T) {
	term := NewANSI()

	err := term.SetCursorStyle(types.CursorStyle(99))
	if err == nil {
		t.Error("SetCursorStyle with invalid style should return error")
	}
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Screen Operations Tests                                         │.
// └─────────────────────────────────────────────────────────────────┘.

func TestANSI_Clear(t *testing.T) {
	got := captureANSI(func(term *ANSITerminal) {
		term.Clear()
	})

	want := "\033[2J\033[H"
	if got != want {
		t.Errorf("Clear = %q, want %q", got, want)
	}
}

func TestANSI_ClearLine(t *testing.T) {
	got := captureANSI(func(term *ANSITerminal) {
		term.ClearLine()
	})

	// CRITICAL: Must include \r to move cursor to start of line!
	want := "\r\033[2K"
	if got != want {
		t.Errorf("ClearLine = %q, want %q", got, want)
	}
}

func TestANSI_ClearFromCursor(t *testing.T) {
	got := captureANSI(func(term *ANSITerminal) {
		term.ClearFromCursor()
	})

	want := "\033[J"
	if got != want {
		t.Errorf("ClearFromCursor = %q, want %q", got, want)
	}
}

func TestANSI_ClearLines(t *testing.T) {
	tests := []struct {
		name  string
		count int
		want  string
	}{
		{"single", 1, "\r\033[J"},
		{"double", 2, "\033[1A\r\033[J"},
		{"ten", 10, "\033[9A\r\033[J"},
		{"zero", 0, ""},      // No-op
		{"negative", -5, ""}, // No-op
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureANSI(func(term *ANSITerminal) {
				term.ClearLines(tt.count)
			})

			if got != tt.want {
				t.Errorf("ClearLines(%d) = %q, want %q", tt.count, got, tt.want)
			}
		})
	}
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Output Tests                                                    │.
// └─────────────────────────────────────────────────────────────────┘.

func TestANSI_Write(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"simple", "Hello, World!"},
		{"multiline", "Line 1\nLine 2\nLine 3"},
		{"empty", ""},
		{"unicode", "Hello 世界 🚀"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureANSI(func(term *ANSITerminal) {
				term.Write(tt.text)
			})

			if got != tt.text {
				t.Errorf("Write(%q) = %q, want %q", tt.text, got, tt.text)
			}
		})
	}
}

func TestANSI_WriteAt(t *testing.T) {
	got := captureANSI(func(term *ANSITerminal) {
		term.WriteAt(10, 5, "Test")
	})

	want := "\033[6;11HTest"
	if got != want {
		t.Errorf("WriteAt = %q, want %q", got, want)
	}
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Screen Buffer Tests                                             │.
// └─────────────────────────────────────────────────────────────────┘.

func TestANSI_ReadScreenBuffer_NotSupported(t *testing.T) {
	term := NewANSI()
	buffer, err := term.ReadScreenBuffer()

	if err == nil {
		t.Error("ReadScreenBuffer should return error on ANSI")
	}

	if buffer != nil {
		t.Errorf("ReadScreenBuffer on error should return nil, got %v", buffer)
	}

	if !strings.Contains(err.Error(), "readback") {
		t.Errorf("Error should mention readback, got: %v", err)
	}
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Terminal Info Tests                                             │.
// └─────────────────────────────────────────────────────────────────┘.

func TestANSI_Size(t *testing.T) {
	term := NewANSI()
	w, h, err := term.Size()

	// Size detection might fail in test environment, but should return fallback.
	if err != nil { //nolint:nestif // Test validation requires nested checks for fallback behavior
		// Fallback should be 80x24.
		if w != 80 || h != 24 {
			t.Errorf("Size fallback should be 80x24, got %dx%d", w, h)
		}
	} else {
		// Real terminal size should be positive.
		if w <= 0 || h <= 0 {
			t.Errorf("Size should return positive values, got %dx%d", w, h)
		}
	}
}

func TestANSI_ColorDepth(t *testing.T) {
	term := NewANSI()
	depth := term.ColorDepth()

	// Color depth should be one of: 16, 256, or 16777216.
	validDepths := map[int]bool{
		16:       true,
		256:      true,
		16777216: true,
	}

	if !validDepths[depth] {
		t.Errorf("ColorDepth = %d, want one of [16, 256, 16777216]", depth)
	}
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Capabilities Tests                                              │.
// └─────────────────────────────────────────────────────────────────┘.

func TestANSI_SupportsDirectPositioning(t *testing.T) {
	term := NewANSI()
	if term.SupportsDirectPositioning() {
		t.Error("ANSI should not support direct positioning")
	}
}

func TestANSI_SupportsReadback(t *testing.T) {
	term := NewANSI()
	if term.SupportsReadback() {
		t.Error("ANSI should not support readback")
	}
}

func TestANSI_SupportsTrueColor(_ *testing.T) {
	term := NewANSI()
	// TrueColor support depends on environment.
	// Just verify method works.
	_ = term.SupportsTrueColor()
}

func TestANSI_Platform(t *testing.T) {
	term := NewANSI()
	platform := term.Platform()

	if platform != types.PlatformUnix {
		t.Errorf("Platform = %v, want PlatformUnix", platform)
	}
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Benchmark Tests                                                 │.
// └─────────────────────────────────────────────────────────────────┘.

func BenchmarkANSI_SetCursorPosition(b *testing.B) {
	r, w, _ := os.Pipe()
	term := &ANSITerminal{output: w}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		term.SetCursorPosition(10, 5)
	}

	w.Close()
	r.Close()
}

func BenchmarkANSI_ClearLines(b *testing.B) {
	r, w, _ := os.Pipe()
	term := &ANSITerminal{output: w}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		term.ClearLines(10)
	}

	w.Close()
	r.Close()
}

func BenchmarkANSI_Write(b *testing.B) {
	r, w, _ := os.Pipe()
	term := &ANSITerminal{output: w}
	text := "Hello, Phoenix TUI!"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		term.Write(text)
	}

	w.Close()
	r.Close()
}

func BenchmarkANSI_WriteAt(b *testing.B) {
	r, w, _ := os.Pipe()
	term := &ANSITerminal{output: w}
	text := "Benchmark"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		term.WriteAt(10, 5, text)
	}

	w.Close()
	r.Close()
}

func BenchmarkANSI_HideShowCursor(b *testing.B) {
	r, w, _ := os.Pipe()
	term := &ANSITerminal{output: w}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		term.HideCursor()
		term.ShowCursor()
	}

	w.Close()
	r.Close()
}

func BenchmarkANSI_MoveCursor(b *testing.B) {
	r, w, _ := os.Pipe()
	term := &ANSITerminal{output: w}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		term.MoveCursorUp(1)
		term.MoveCursorDown(1)
		term.MoveCursorLeft(1)
		term.MoveCursorRight(1)
	}

	w.Close()
	r.Close()
}
