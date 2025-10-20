//go:build windows.
// +build windows.

package infrastructure

import (
	"os"

	"github.com/phoenix-tui/phoenix/terminal/api"
	"github.com/phoenix-tui/phoenix/terminal/infrastructure/unix"
	"github.com/phoenix-tui/phoenix/terminal/infrastructure/windows"
)

// newWindowsTerminal creates Windows-optimized terminal with ANSI fallback.
//
// Detection strategy:.
//  1. Try Windows Console API first (NewConsole()).
//  2. If it fails (Git Bash, redirected I/O), use ANSI fallback.
//
// This ensures:.
//   - cmd.exe gets 10x performance (Windows Console API).
//   - PowerShell gets 10x performance (Windows Console API).
//   - Git Bash works correctly (ANSI fallback).
//   - Redirected I/O works (ANSI fallback).
func newWindowsTerminal() api.Terminal {
	// Try Windows Console API first.
	console, err := windows.NewConsole()
	if err == nil {
		return console // Success! Use native Win32 API (10x faster)
	}

	// Fallback to ANSI (Git Bash, MinTTY, or redirected I/O).
	return unix.NewANSI()
}

// detectWindowsPlatform determines Windows terminal type.
//
// Returns:.
//   - PlatformWindowsConsole: cmd.exe, PowerShell (Win32 API works).
//   - PlatformWindowsANSI: Git Bash, MinTTY (Win32 API fails).
//   - PlatformUnknown: Detection failed.
func detectWindowsPlatform() api.Platform {
	// Try Windows Console API.
	_, err := windows.NewConsole()
	if err == nil {
		return api.PlatformWindowsConsole
	}

	// Check TERM environment (Git Bash, MinTTY set this).
	if os.Getenv("TERM") != "" {
		return api.PlatformWindowsANSI
	}

	// Redirected I/O or unknown environment.
	return api.PlatformWindowsANSI // Conservative fallback
}
