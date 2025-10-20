// Package platform detects platform type and SSH session environment.
package platform

import (
	"os"
	"runtime"
)

// Type represents the platform type.
type Type string

const (
	// TypeWindows represents Windows platform.
	TypeWindows Type = "windows"

	// TypeDarwin represents macOS platform.
	TypeDarwin Type = "darwin"

	// TypeLinux represents Linux platform.
	TypeLinux Type = "linux"

	// TypeUnknown represents an unknown platform.
	TypeUnknown Type = "unknown"
)

// Detector detects the current platform and environment.
type Detector struct{}

// NewDetector creates a new platform detector.
func NewDetector() *Detector {
	return &Detector{}
}

// GetPlatform returns the current platform type.
func (d *Detector) GetPlatform() Type {
	switch runtime.GOOS {
	case "windows":
		return TypeWindows
	case "darwin":
		return TypeDarwin
	case "linux":
		return TypeLinux
	default:
		return TypeUnknown
	}
}

// IsSSH returns true if running in an SSH session.
func (d *Detector) IsSSH() bool {
	// Check for SSH environment variables
	if os.Getenv("SSH_TTY") != "" {
		return true
	}

	if os.Getenv("SSH_CLIENT") != "" {
		return true
	}

	if os.Getenv("SSH_CONNECTION") != "" {
		return true
	}

	return false
}

// IsHeadless returns true if running in a headless environment.
func (d *Detector) IsHeadless() bool {
	// Check for DISPLAY on Linux (X11)
	if d.GetPlatform() == TypeLinux {
		if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
			return true
		}
	}

	return false
}

// HasTTY returns true if stdout is a terminal.
func (d *Detector) HasTTY() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	// Check if stdout is a character device (terminal)
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// ShouldUseOSC52 returns true if OSC 52 should be preferred.
func (d *Detector) ShouldUseOSC52() bool {
	// Use OSC 52 if in SSH session and has a TTY
	return d.IsSSH() && d.HasTTY()
}
