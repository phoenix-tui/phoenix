package terminal

import (
	"runtime"

	"github.com/phoenix-tui/phoenix/terminal/internal/infrastructure/unix"
	"github.com/phoenix-tui/phoenix/terminal/types"
)

// internalTerminal is the internal interface used by infrastructure implementations.
// It uses types.Platform and types.CursorStyle to avoid cyclic dependencies.
type internalTerminal interface {
	SetCursorPosition(x, y int) error
	GetCursorPosition() (x, y int, err error)
	MoveCursorUp(n int) error
	MoveCursorDown(n int) error
	MoveCursorLeft(n int) error
	MoveCursorRight(n int) error
	SaveCursorPosition() error
	RestoreCursorPosition() error
	HideCursor() error
	ShowCursor() error
	SetCursorStyle(style types.CursorStyle) error // Uses types
	Clear() error
	ClearLine() error
	ClearFromCursor() error
	ClearLines(count int) error
	Write(s string) error
	WriteAt(x, y int, s string) error
	ReadScreenBuffer() ([][]rune, error)
	Size() (width, height int, err error)
	ColorDepth() int
	SupportsDirectPositioning() bool
	SupportsReadback() bool
	SupportsTrueColor() bool
	Platform() types.Platform // Uses types
	EnterAltScreen() error
	ExitAltScreen() error
	IsInAltScreen() bool
	IsInRawMode() bool
	EnterRawMode() error
	ExitRawMode() error
}

// terminalAdapter wraps internal terminal and converts types.
type terminalAdapter struct {
	internal internalTerminal
}

func (t *terminalAdapter) SetCursorPosition(x, y int) error {
	return t.internal.SetCursorPosition(x, y)
}
func (t *terminalAdapter) GetCursorPosition() (int, int, error) {
	return t.internal.GetCursorPosition()
}
func (t *terminalAdapter) MoveCursorUp(n int) error     { return t.internal.MoveCursorUp(n) }
func (t *terminalAdapter) MoveCursorDown(n int) error   { return t.internal.MoveCursorDown(n) }
func (t *terminalAdapter) MoveCursorLeft(n int) error   { return t.internal.MoveCursorLeft(n) }
func (t *terminalAdapter) MoveCursorRight(n int) error  { return t.internal.MoveCursorRight(n) }
func (t *terminalAdapter) SaveCursorPosition() error    { return t.internal.SaveCursorPosition() }
func (t *terminalAdapter) RestoreCursorPosition() error { return t.internal.RestoreCursorPosition() }
func (t *terminalAdapter) HideCursor() error            { return t.internal.HideCursor() }
func (t *terminalAdapter) ShowCursor() error            { return t.internal.ShowCursor() }
func (t *terminalAdapter) SetCursorStyle(style CursorStyle) error {
	return t.internal.SetCursorStyle(types.CursorStyle(style)) // Convert
}
func (t *terminalAdapter) Clear() error                        { return t.internal.Clear() }
func (t *terminalAdapter) ClearLine() error                    { return t.internal.ClearLine() }
func (t *terminalAdapter) ClearFromCursor() error              { return t.internal.ClearFromCursor() }
func (t *terminalAdapter) ClearLines(count int) error          { return t.internal.ClearLines(count) }
func (t *terminalAdapter) Write(s string) error                { return t.internal.Write(s) }
func (t *terminalAdapter) WriteAt(x, y int, s string) error    { return t.internal.WriteAt(x, y, s) }
func (t *terminalAdapter) ReadScreenBuffer() ([][]rune, error) { return t.internal.ReadScreenBuffer() }
func (t *terminalAdapter) Size() (int, int, error)             { return t.internal.Size() }
func (t *terminalAdapter) ColorDepth() int                     { return t.internal.ColorDepth() }
func (t *terminalAdapter) SupportsDirectPositioning() bool {
	return t.internal.SupportsDirectPositioning()
}
func (t *terminalAdapter) SupportsReadback() bool  { return t.internal.SupportsReadback() }
func (t *terminalAdapter) SupportsTrueColor() bool { return t.internal.SupportsTrueColor() }
func (t *terminalAdapter) Platform() Platform {
	return Platform(t.internal.Platform()) // Convert
}
func (t *terminalAdapter) EnterAltScreen() error { return t.internal.EnterAltScreen() }
func (t *terminalAdapter) ExitAltScreen() error  { return t.internal.ExitAltScreen() }
func (t *terminalAdapter) IsInAltScreen() bool   { return t.internal.IsInAltScreen() }
func (t *terminalAdapter) IsInRawMode() bool     { return t.internal.IsInRawMode() }
func (t *terminalAdapter) EnterRawMode() error   { return t.internal.EnterRawMode() }
func (t *terminalAdapter) ExitRawMode() error    { return t.internal.ExitRawMode() }

// New creates platform-optimized terminal with auto-detection.
//
// Platform detection and optimization:
//   - Windows: Tries Console API first, falls back to ANSI (Git Bash)
//   - Unix (Linux/macOS): Always uses ANSI
//
// Auto-fallback ensures compatibility:
//   - cmd.exe → Windows Console API (10x faster!)
//   - PowerShell → Windows Console API (10x faster!)
//   - Git Bash → ANSI (GetConsoleScreenBufferInfo fails, auto-fallback)
//   - WSL → ANSI
//   - Redirected output → ANSI (no console handle, auto-fallback)
//
// Example:
//
//	term := terminal.New()
//	term.HideCursor()
//	term.SetCursorPosition(10, 5)
//	term.Write("Hello, Phoenix!")
//	term.ShowCursor()
func New() Terminal {
	if runtime.GOOS == "windows" {
		return newWindowsTerminal() // Implemented in new_windows.go
	}
	return &terminalAdapter{internal: unix.NewANSI()}
}

// NewANSI forces ANSI implementation regardless of platform.
//
// Use cases:
//   - Testing ANSI code generation
//   - Forcing fallback behavior
//   - Redirected output (pipes, files)
//   - Environments where Win32 API shouldn't be used
//   - Cross-platform consistency testing
//
// Example:
//
//	term := terminal.NewANSI()
//	// Guaranteed to use ANSI escape codes, even on Windows
func NewANSI() Terminal {
	return &terminalAdapter{internal: unix.NewANSI()}
}

// DetectPlatform identifies current terminal platform type.
//
// Detection logic:
//   - Non-Windows: Returns PlatformUnix
//   - Windows: Attempts to detect Console API vs ANSI mode
//
// Example:
//
//	platform := terminal.DetectPlatform()
//	if platform == terminal.PlatformWindowsConsole {
//	    fmt.Println("Using native Windows Console API")
//	}
func DetectPlatform() Platform {
	if runtime.GOOS != "windows" {
		return PlatformUnix
	}
	return detectWindowsPlatform() // Implemented in new_windows.go
}
