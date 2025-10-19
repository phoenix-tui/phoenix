//go:build !windows
// +build !windows

package infrastructure

import (
	"github.com/phoenix-tui/phoenix/terminal/api"
	"github.com/phoenix-tui/phoenix/terminal/infrastructure/unix"
)

// newWindowsTerminal is a stub for non-Windows platforms.
// This function is never called on Unix systems, but needs to exist
// for compilation to succeed.
func newWindowsTerminal() api.Terminal {
	// This should never be called on non-Windows systems,
	// but we return ANSI as a safe fallback.
	return unix.NewANSI()
}

// detectWindowsPlatform is a stub for non-Windows platforms.
// This function is never called on Unix systems.
func detectWindowsPlatform() api.Platform {
	// This should never be called on non-Windows systems
	return api.PlatformUnknown
}
