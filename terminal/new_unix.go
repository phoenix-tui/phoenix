//go:build !windows

package terminal

import (
	"github.com/phoenix-tui/phoenix/terminal/internal/infrastructure/unix"
)

// newWindowsTerminal stub for non-Windows platforms.
// This function is never called on Unix (runtime.GOOS check in new.go),
// but must exist for compilation.
func newWindowsTerminal() Terminal {
	// Fallback to ANSI (should never be reached due to runtime.GOOS check).
	return &terminalAdapter{internal: unix.NewANSI()}
}

// detectWindowsPlatform stub for non-Windows platforms.
// This function is never called on Unix (runtime.GOOS check in new.go),
// but must exist for compilation.
func detectWindowsPlatform() Platform {
	// Return PlatformUnix (should never be reached due to runtime.GOOS check).
	return PlatformUnix
}
