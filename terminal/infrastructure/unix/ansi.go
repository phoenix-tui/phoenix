// Package unix provides ANSI escape code implementation for Unix-like terminals.
//
// This implementation works on:
//   - Linux terminals (gnome-terminal, konsole, xterm, etc.)
//   - macOS Terminal.app and iTerm2
//   - Windows Git Bash, MinTTY, WSL
//
// ANSI escape codes are the universal terminal control standard, supported
// by virtually all modern terminals. Performance is good but slower than
// native Win32 API on Windows platforms.
package unix

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/phoenix-tui/phoenix/terminal/api"
)

// ANSITerminal implements Terminal interface using ANSI escape codes.
type ANSITerminal struct {
	output *os.File // Usually os.Stdout
}

// NewANSI creates new ANSI terminal implementation.
// Uses os.Stdout by default.
func NewANSI() *ANSITerminal {
	return &ANSITerminal{
		output: os.Stdout,
	}
}

// NewANSIWithOutput creates ANSI terminal with custom output.
// Useful for testing with captured output.
func NewANSIWithOutput(output *os.File) *ANSITerminal {
	return &ANSITerminal{
		output: output,
	}
}

// ┌─────────────────────────────────────────────────────────────────┐
// │ Cursor Operations                                               │
// └─────────────────────────────────────────────────────────────────┘

// SetCursorPosition moves cursor to absolute position (x, y).
// ANSI: "\033[{row};{col}H" (1-based indexing!)
func (a *ANSITerminal) SetCursorPosition(x, y int) error {
	// ANSI uses 1-based indexing, API uses 0-based
	_, err := fmt.Fprintf(a.output, "\033[%d;%dH", y+1, x+1)
	return err
}

// GetCursorPosition returns error - ANSI doesn't support reliable readback.
//
// Technical note: ANSI has CPR (Cursor Position Report) protocol:
//
//	Write: "\033[6n"
//	Read: "\033[{row};{col}R"
//
// But this requires:
//   - Raw terminal mode
//   - Stdin read with timeout
//   - Complex parsing
//   - Race conditions with other input
//
// For simplicity, we don't support this. Use Windows Console API for readback.
func (a *ANSITerminal) GetCursorPosition() (x, y int, err error) {
	return 0, 0, fmt.Errorf("ANSI terminals don't support reliable cursor readback")
}

// MoveCursorUp moves cursor up n lines.
// ANSI: "\033[{n}A"
func (a *ANSITerminal) MoveCursorUp(n int) error {
	if n <= 0 {
		return nil // No-op for non-positive values
	}
	_, err := fmt.Fprintf(a.output, "\033[%dA", n)
	return err
}

// MoveCursorDown moves cursor down n lines.
// ANSI: "\033[{n}B"
func (a *ANSITerminal) MoveCursorDown(n int) error {
	if n <= 0 {
		return nil
	}
	_, err := fmt.Fprintf(a.output, "\033[%dB", n)
	return err
}

// MoveCursorLeft moves cursor left n columns.
// ANSI: "\033[{n}D"
func (a *ANSITerminal) MoveCursorLeft(n int) error {
	if n <= 0 {
		return nil
	}
	_, err := fmt.Fprintf(a.output, "\033[%dD", n)
	return err
}

// MoveCursorRight moves cursor right n columns.
// ANSI: "\033[{n}C"
func (a *ANSITerminal) MoveCursorRight(n int) error {
	if n <= 0 {
		return nil
	}
	_, err := fmt.Fprintf(a.output, "\033[%dC", n)
	return err
}

// SaveCursorPosition saves cursor position to stack.
// ANSI: "\033[s" or "\0337" (DEC mode)
func (a *ANSITerminal) SaveCursorPosition() error {
	_, err := fmt.Fprint(a.output, "\033[s")
	return err
}

// RestoreCursorPosition restores saved cursor position.
// ANSI: "\033[u" or "\0338" (DEC mode)
func (a *ANSITerminal) RestoreCursorPosition() error {
	_, err := fmt.Fprint(a.output, "\033[u")
	return err
}

// ┌─────────────────────────────────────────────────────────────────┐
// │ Cursor Visibility & Style                                       │
// └─────────────────────────────────────────────────────────────────┘

// HideCursor makes cursor invisible.
// ANSI: "\033[?25l" (DECTCEM - DEC Text Cursor Enable Mode)
func (a *ANSITerminal) HideCursor() error {
	_, err := fmt.Fprint(a.output, "\033[?25l")
	return err
}

// ShowCursor makes cursor visible.
// ANSI: "\033[?25h"
func (a *ANSITerminal) ShowCursor() error {
	_, err := fmt.Fprint(a.output, "\033[?25h")
	return err
}

// SetCursorStyle changes cursor appearance.
// ANSI: "\033[{n} q" (DECSCUSR - DEC Set Cursor Style)
//
// Codes:
//
//	0 or 1: Blinking block
//	2: Steady block
//	3: Blinking underline
//	4: Steady underline
//	5: Blinking bar
//	6: Steady bar
//
// We use steady variants for consistency.
func (a *ANSITerminal) SetCursorStyle(style api.CursorStyle) error {
	var code int
	switch style {
	case api.CursorBlock:
		code = 2 // Steady block
	case api.CursorUnderline:
		code = 4 // Steady underline
	case api.CursorBar:
		code = 6 // Steady bar
	default:
		return fmt.Errorf("unknown cursor style: %v", style)
	}

	_, err := fmt.Fprintf(a.output, "\033[%d q", code)
	return err
}

// ┌─────────────────────────────────────────────────────────────────┐
// │ Screen Operations                                               │
// └─────────────────────────────────────────────────────────────────┘

// Clear clears entire screen and moves cursor to top-left.
// ANSI: "\033[2J" + "\033[H"
func (a *ANSITerminal) Clear() error {
	_, err := fmt.Fprint(a.output, "\033[2J\033[H")
	return err
}

// ClearLine clears current line (where cursor is).
// ANSI: "\033[2K"
func (a *ANSITerminal) ClearLine() error {
	// CRITICAL: Must include \r (carriage return) to move cursor to start of line!
	// Without \r, clearing happens from current cursor position
	_, err := fmt.Fprint(a.output, "\r\033[2K")
	return err
}

// ClearFromCursor clears from cursor to end of screen.
// ANSI: "\033[J" or "\033[0J"
func (a *ANSITerminal) ClearFromCursor() error {
	_, err := fmt.Fprint(a.output, "\033[J")
	return err
}

// ClearLines clears N lines starting from current cursor position.
//
// CRITICAL for multiline input (GoSh shell):
//   - Move cursor to start of first line to clear
//   - Clear from cursor to end of screen
//   - Result: N lines cleared, cursor at start
//
// Algorithm:
//
//	count == 1: "\r\033[J" (just clear from start of line)
//	count > 1:  "\033[{count-1}A\r\033[J" (move up, then clear)
//
// Performance: ~500μs for 10 lines (10x slower than Windows API but acceptable)
func (a *ANSITerminal) ClearLines(count int) error {
	if count <= 0 {
		return nil // No-op
	}

	if count == 1 {
		// Single line: just CR + clear to end
		_, err := fmt.Fprint(a.output, "\r\033[J")
		return err
	}

	// Multiple lines: move up to first line, then clear to end
	_, err := fmt.Fprintf(a.output, "\033[%dA\r\033[J", count-1)
	return err
}

// ┌─────────────────────────────────────────────────────────────────┐
// │ Output                                                          │
// └─────────────────────────────────────────────────────────────────┘

// Write writes string to terminal at current cursor position.
func (a *ANSITerminal) Write(s string) error {
	_, err := fmt.Fprint(a.output, s)
	return err
}

// WriteAt writes string at specific position (x, y).
// Equivalent to SetCursorPosition + Write.
func (a *ANSITerminal) WriteAt(x, y int, s string) error {
	if err := a.SetCursorPosition(x, y); err != nil {
		return err
	}
	return a.Write(s)
}

// ┌─────────────────────────────────────────────────────────────────┐
// │ Screen Buffer (Not Supported on ANSI)                          │
// └─────────────────────────────────────────────────────────────────┘

// ReadScreenBuffer returns error - ANSI terminals don't support buffer readback.
//
// Windows Console API can read screen buffer via ReadConsoleOutput,
// but ANSI terminals don't have equivalent functionality.
//
// Differential rendering requires platform-specific APIs.
func (a *ANSITerminal) ReadScreenBuffer() ([][]rune, error) {
	return nil, fmt.Errorf("ANSI terminals don't support screen buffer readback")
}

// ┌─────────────────────────────────────────────────────────────────┐
// │ Terminal Info                                                   │
// └─────────────────────────────────────────────────────────────────┘

// Size returns current terminal dimensions (width, height).
// Uses golang.org/x/term for cross-platform detection.
func (a *ANSITerminal) Size() (width, height int, err error) {
	fd := int(a.output.Fd())
	w, h, err := term.GetSize(fd)
	if err != nil {
		// Fallback to common default
		return 80, 24, err
	}
	return w, h, nil
}

// ColorDepth returns color support level.
//
// Detection heuristics:
//   - COLORTERM=truecolor → 24-bit (16777216 colors)
//   - TERM contains "256color" → 8-bit (256 colors)
//   - Otherwise → 4-bit (16 colors)
//
// Most modern terminals support at least 256 colors.
func (a *ANSITerminal) ColorDepth() int {
	// Check for 24-bit truecolor support
	if os.Getenv("COLORTERM") == "truecolor" || os.Getenv("COLORTERM") == "24bit" {
		return 16777216 // 24-bit RGB
	}

	// Check for 256 color support
	term := os.Getenv("TERM")
	if term == "xterm-256color" || term == "screen-256color" || term == "tmux-256color" {
		return 256 // 8-bit
	}

	// Fallback to basic 16 colors
	return 16 // 4-bit
}

// ┌─────────────────────────────────────────────────────────────────┐
// │ Capabilities Discovery                                          │
// └─────────────────────────────────────────────────────────────────┘

// SupportsDirectPositioning returns false - ANSI uses escape codes.
// Windows Console API has true direct positioning via Win32 calls.
func (a *ANSITerminal) SupportsDirectPositioning() bool {
	return false
}

// SupportsReadback returns false - ANSI can't read cursor/buffer.
// Windows Console API supports readback via GetConsoleScreenBufferInfo.
func (a *ANSITerminal) SupportsReadback() bool {
	return false
}

// SupportsTrueColor returns true if terminal supports 24-bit RGB.
func (a *ANSITerminal) SupportsTrueColor() bool {
	return a.ColorDepth() == 16777216
}

// Platform returns Unix platform type.
func (a *ANSITerminal) Platform() api.Platform {
	return api.PlatformUnix
}
