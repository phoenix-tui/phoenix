// Package types defines shared types for the terminal package.
//
// This package contains type definitions and constants that are used
// by both the terminal package and its internal implementations.
// Separated to avoid import cycles.
package types

// CursorStyle represents the visual appearance of the terminal cursor.
type CursorStyle int

// Cursor style constants.
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

// Platform represents the terminal platform type.
type Platform int

// Platform type constants.
const (
	// PlatformUnknown represents unknown or undetected platform.
	PlatformUnknown Platform = iota

	// PlatformUnix represents Unix-like systems (Linux, macOS, BSD).
	// Uses ANSI escape codes for all operations.
	PlatformUnix

	// PlatformWindowsConsole represents Windows Console API.
	// Uses native Win32 API calls (fastest on Windows cmd.exe/PowerShell).
	PlatformWindowsConsole

	// PlatformWindowsANSI represents Windows ANSI mode.
	// Uses ANSI escape codes (Git Bash, MinTTY, WSL).
	PlatformWindowsANSI
)

// String returns human-readable platform name.
func (p Platform) String() string {
	switch p {
	case PlatformUnix:
		return "Unix (ANSI)"
	case PlatformWindowsConsole:
		return "Windows Console (Win32 API)"
	case PlatformWindowsANSI:
		return "Windows ANSI (Git Bash)"
	default:
		return "Unknown"
	}
}
