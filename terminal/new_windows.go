//go:build windows

package terminal

import (
	"github.com/phoenix-tui/phoenix/terminal/internal/infrastructure/unix"
	"github.com/phoenix-tui/phoenix/terminal/internal/infrastructure/windows"
)

// newWindowsTerminal creates Windows-optimized terminal with auto-fallback.
//
// Detection algorithm:
//  1. Try Windows Console API first (windows.NewConsole())
//  2. If Console API fails → Auto-fallback to ANSI (unix.NewANSI())
//
// Failure cases that trigger fallback:
//   - Git Bash (GetConsoleScreenBufferInfo returns error)
//   - WSL (no Windows console handle)
//   - Redirected I/O (stdout is pipe/file, not console)
//   - SSH sessions without PTY
func newWindowsTerminal() Terminal {
	// Try Windows Console API first (10x faster!)
	term, err := windows.NewConsole()
	if err == nil {
		return &terminalAdapter{internal: term}
	}

	// Fallback to ANSI (Git Bash, WSL, pipes)
	return &terminalAdapter{internal: unix.NewANSI()}
}

// detectWindowsPlatform detects Windows terminal type.
//
// Returns:
//   - PlatformWindowsConsole: Native Console API available
//   - PlatformWindowsANSI: ANSI mode (Git Bash, WSL, etc.)
func detectWindowsPlatform() Platform {
	// Try to detect Console API availability
	_, err := windows.NewConsole()
	if err == nil {
		return PlatformWindowsConsole
	}
	return PlatformWindowsANSI
}
