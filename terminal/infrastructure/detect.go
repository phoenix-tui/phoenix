// Package infrastructure provides platform-specific terminal implementations.
package infrastructure

import (
	"runtime"

	"github.com/phoenix-tui/phoenix/terminal/api"
	"github.com/phoenix-tui/phoenix/terminal/infrastructure/unix"
)

// NewTerminal creates platform-optimized terminal with auto-detection.
//
// Week 16 Implementation (Windows Console API + ANSI):
//   - Detects Windows platform (runtime.GOOS == "windows")
//   - Tries Windows Console API first (windows.NewConsole())
//   - Auto-fallbacks to ANSI on error (Git Bash, redirected I/O)
//   - Unix platforms always use ANSI
//
// Detection Algorithm:
//
//	if runtime.GOOS == "windows" {
//	    return newWindowsTerminal() // Tries Win32 API, falls back to ANSI
//	}
//	return newUnixTerminal() // Linux/macOS - always ANSI
//
// Auto-fallback ensures compatibility:
//   - cmd.exe → Windows Console API (10x faster!)
//   - PowerShell → Windows Console API (10x faster!)
//   - Git Bash → ANSI (GetConsoleScreenBufferInfo fails, auto-fallback)
//   - WSL → ANSI
//   - Redirected output → ANSI (no console handle, auto-fallback)
//
// Platform-specific implementations:
//   - detect_windows.go: Windows detection logic (build tag: windows)
//   - detect_unix.go: Unix stub (build tag: !windows)
func NewTerminal() api.Terminal {
	if runtime.GOOS == "windows" {
		return newWindowsTerminal() // Implemented in detect_windows.go
	}
	return newUnixTerminal()
}

// NewANSITerminal forces ANSI implementation regardless of platform.
//
// Use cases:
//   - Testing ANSI code generation
//   - Forcing fallback behavior
//   - Redirected output (pipes, files)
//   - Environments where Win32 API shouldn't be used
//   - Cross-platform consistency testing
func NewANSITerminal() api.Terminal {
	return unix.NewANSI()
}

// newUnixTerminal creates ANSI terminal for Unix-like systems.
// Always returns ANSI implementation (Linux, macOS, BSD).
func newUnixTerminal() api.Terminal {
	return unix.NewANSI()
}

// DetectPlatform identifies current terminal platform type.
//
// Week 16 Implementation:
//   - Detects Windows Console vs Windows ANSI vs Unix
//   - Returns platform-specific identifier for diagnostics
//
// Detection logic:
//
//	if runtime.GOOS != "windows" {
//	    return api.PlatformUnix
//	}
//	return detectWindowsPlatform() // Try Win32 API, check TERM env
//
// Platform-specific implementations:
//   - detect_windows.go: Windows detection (build tag: windows)
//   - detect_unix.go: Unix stub (build tag: !windows)
func DetectPlatform() api.Platform {
	if runtime.GOOS != "windows" {
		return api.PlatformUnix
	}
	return detectWindowsPlatform() // Implemented in detect_windows.go
}
