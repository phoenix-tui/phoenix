// Package api provides a platform-optimized terminal operations abstraction layer.
//
// Phoenix Terminal implements hybrid approach for best performance on each platform:
//   - Windows Console: Direct Win32 API calls (10x faster)
//   - Windows Git Bash: ANSI escape codes (automatic fallback)
//   - Unix (Linux/macOS): ANSI escape codes (universal)
//
// Example:
//
//	term := terminal.New() // Auto-detects best implementation
//	term.HideCursor()
//	term.SetCursorPosition(10, 5)
//	term.Write("Hello, Phoenix!")
//	term.ShowCursor()
//
// Platform detection is automatic - no configuration needed:
//
//	if term.SupportsDirectPositioning() {
//		// Use fast absolute positioning
//		term.WriteAt(x, y, content)
//	} else {
//		// Use ANSI relative movements
//		term.Write(content)
//	}
package api

// No imports needed - API layer is independent

// Terminal provides platform-optimized terminal operations.
//
// All methods are safe to call from any goroutine, though performance
// is best when called from the main event loop thread.
//
// Error handling: Most operations return error for robustness, but in
// typical usage errors are rare (write to stdout). Check errors in
// critical sections (e.g., before major rendering).
type Terminal interface {
	// ┌─────────────────────────────────────────────────────────────┐
	// │ Cursor Operations                                           │
	// └─────────────────────────────────────────────────────────────┘

	// SetCursorPosition moves the cursor to absolute position (x, y).
	// Coordinates are 0-based (top-left is 0,0).
	//
	// Windows Console API: Direct Win32 call (~10μs)
	// ANSI: Escape code "\033[{row};{col}H" (~100μs)
	//
	// Returns error if position is out of bounds or write fails.
	SetCursorPosition(x, y int) error

	// GetCursorPosition returns current cursor position (x, y).
	// Coordinates are 0-based (top-left is 0,0).
	//
	// Windows Console API: Instant readback via GetConsoleScreenBufferInfo
	// ANSI: Returns error (requires CPR protocol, unreliable)
	//
	// Use SupportsReadback() to check if this is available.
	GetCursorPosition() (x, y int, err error)

	// MoveCursorUp moves cursor up n lines (relative movement).
	MoveCursorUp(n int) error

	// MoveCursorDown moves cursor down n lines (relative movement).
	MoveCursorDown(n int) error

	// MoveCursorLeft moves cursor left n columns (relative movement).
	MoveCursorLeft(n int) error

	// MoveCursorRight moves cursor right n columns (relative movement).
	MoveCursorRight(n int) error

	// SaveCursorPosition saves current cursor position to stack.
	// Must be paired with RestoreCursorPosition().
	SaveCursorPosition() error

	// RestoreCursorPosition restores previously saved cursor position.
	// Must be called after SaveCursorPosition().
	RestoreCursorPosition() error

	// ┌─────────────────────────────────────────────────────────────┐
	// │ Cursor Visibility & Style                                   │
	// └─────────────────────────────────────────────────────────────┘

	// HideCursor makes the cursor invisible.
	// IMPORTANT: Always pair with ShowCursor() to restore visibility!
	HideCursor() error

	// ShowCursor makes the cursor visible.
	ShowCursor() error

	// SetCursorStyle changes cursor appearance.
	// Not all terminals support all styles - check terminal documentation.
	SetCursorStyle(style CursorStyle) error

	// ┌─────────────────────────────────────────────────────────────┐
	// │ Screen Operations                                           │
	// └─────────────────────────────────────────────────────────────┘

	// Clear clears the entire screen.
	// Cursor position is typically moved to top-left (0,0).
	Clear() error

	// ClearLine clears the current line (where cursor is located).
	// Cursor position remains unchanged.
	ClearLine() error

	// ClearFromCursor clears from cursor to end of screen.
	// Useful for clearing stale content after rendering.
	ClearFromCursor() error

	// ClearLines clears N lines starting from current cursor position.
	//
	// CRITICAL for multiline input (like GoSh shell):
	//   - Efficiently clears multiple lines of previous content
	//   - Positions cursor at start of cleared region
	//
	// Windows Console API: FillConsoleOutputCharacter (~50μs for 10 lines)
	// ANSI: Move up + clear to end (~500μs for 10 lines)
	ClearLines(count int) error

	// ┌─────────────────────────────────────────────────────────────┐
	// │ Output                                                      │
	// └─────────────────────────────────────────────────────────────┘

	// Write writes string to terminal at current cursor position.
	// Cursor advances automatically.
	Write(s string) error

	// WriteAt writes string to terminal at specific position (x, y).
	//
	// Equivalent to:
	//   SetCursorPosition(x, y)
	//   Write(s)
	//
	// But optimized on platforms that support direct positioning.
	WriteAt(x, y int, s string) error

	// ┌─────────────────────────────────────────────────────────────┐
	// │ Screen Buffer (Windows Console API only)                    │
	// └─────────────────────────────────────────────────────────────┘

	// ReadScreenBuffer reads entire screen buffer content.
	//
	// Enables differential rendering (like PSReadLine):
	//   oldBuffer := term.ReadScreenBuffer()
	//   // ... compute changes ...
	//   term.WriteOnlyDiff(oldBuffer, newBuffer)
	//
	// Windows Console API: Supported via ReadConsoleOutput
	// ANSI: Returns error (not supported)
	//
	// Use SupportsReadback() to check if this is available.
	ReadScreenBuffer() ([][]rune, error)

	// ┌─────────────────────────────────────────────────────────────┐
	// │ Terminal Info                                               │
	// └─────────────────────────────────────────────────────────────┘

	// Size returns current terminal dimensions (width, height).
	// Returns (80, 24) as fallback if detection fails.
	Size() (width, height int, err error)

	// ColorDepth returns number of colors supported.
	//   - 16: Basic ANSI colors
	//   - 256: Extended ANSI colors
	//   - 16777216: True color (24-bit RGB)
	ColorDepth() int

	// ┌─────────────────────────────────────────────────────────────┐
	// │ Capabilities Discovery                                      │
	// └─────────────────────────────────────────────────────────────┘

	// SupportsDirectPositioning returns true if terminal supports
	// fast absolute cursor positioning (Windows Console API).
	//
	// If false, prefer relative movements (MoveCursorUp/Down/Left/Right).
	SupportsDirectPositioning() bool

	// SupportsReadback returns true if terminal supports reading
	// cursor position and screen buffer (Windows Console API).
	//
	// If false, GetCursorPosition() and ReadScreenBuffer() will fail.
	SupportsReadback() bool

	// SupportsTrueColor returns true if terminal supports 24-bit RGB colors.
	SupportsTrueColor() bool

	// Platform returns the detected terminal platform type.
	Platform() Platform
}

// Platform identifies the terminal platform type.
type Platform int

const (
	// PlatformUnknown means platform detection failed.
	PlatformUnknown Platform = iota

	// PlatformWindowsConsole means Windows Console API (cmd.exe, PowerShell).
	// - Fast direct positioning (Win32 API calls)
	// - Supports cursor readback
	// - Supports screen buffer readback
	PlatformWindowsConsole

	// PlatformWindowsANSI means Windows terminal using ANSI (Git Bash, MinTTY).
	// - ANSI escape codes only
	// - No cursor/buffer readback
	PlatformWindowsANSI

	// PlatformUnix means Linux/macOS terminal.
	// - ANSI escape codes only
	// - No cursor/buffer readback
	PlatformUnix
)

// String returns human-readable platform name.
func (p Platform) String() string {
	switch p {
	case PlatformWindowsConsole:
		return "Windows Console (Win32 API)"
	case PlatformWindowsANSI:
		return "Windows ANSI (Git Bash)"
	case PlatformUnix:
		return "Unix (ANSI)"
	default:
		return "Unknown"
	}
}

// CursorStyle defines cursor appearance.
type CursorStyle int

const (
	// CursorBlock shows cursor as filled block.
	CursorBlock CursorStyle = iota

	// CursorUnderline shows cursor as underline.
	CursorUnderline

	// CursorBar shows cursor as vertical bar (|).
	CursorBar
)

// String returns human-readable cursor style name.
func (c CursorStyle) String() string {
	switch c {
	case CursorBlock:
		return "Block"
	case CursorUnderline:
		return "Underline"
	case CursorBar:
		return "Bar"
	default:
		return "Unknown"
	}
}

// ┌─────────────────────────────────────────────────────────────────┐
// │ Constructors (implemented in infrastructure package)           │
// └─────────────────────────────────────────────────────────────────┘
//
// Note: New() and NewANSI() are implemented in the infrastructure
// package to avoid import cycles. Import from there:
//
//	import "github.com/phoenix-tui/phoenix/terminal/infrastructure"
//	term := infrastructure.NewTerminal()
//
// Or use the convenience wrapper in a higher-level package.
