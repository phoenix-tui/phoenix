package platform

import (
	"os"
	"runtime"
	"testing"
)

func TestDetector_GetPlatform(t *testing.T) {
	detector := NewDetector()
	platform := detector.GetPlatform()

	// Verify it matches runtime.GOOS
	expected := Type(runtime.GOOS)
	if runtime.GOOS == "windows" {
		expected = TypeWindows
	} else if runtime.GOOS == "darwin" {
		expected = TypeDarwin
	} else if runtime.GOOS == "linux" {
		expected = TypeLinux
	} else {
		expected = TypeUnknown
	}

	if platform != expected {
		t.Errorf("expected platform %s, got %s", expected, platform)
	}
}

func TestDetector_IsSSH(t *testing.T) {
	detector := NewDetector()

	// Save original values
	origSSHTTY := os.Getenv("SSH_TTY")
	origSSHClient := os.Getenv("SSH_CLIENT")
	origSSHConnection := os.Getenv("SSH_CONNECTION")

	// Restore after test
	defer func() {
		os.Setenv("SSH_TTY", origSSHTTY)
		os.Setenv("SSH_CLIENT", origSSHClient)
		os.Setenv("SSH_CONNECTION", origSSHConnection)
	}()

	tests := []struct {
		name      string
		sshTTY    string
		sshClient string
		sshConn   string
		wantSSH   bool
	}{
		{"SSH_TTY set", "/dev/pts/0", "", "", true},
		{"SSH_CLIENT set", "", "127.0.0.1", "", true},
		{"SSH_CONNECTION set", "", "", "127.0.0.1 22", true},
		{"no SSH vars", "", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("SSH_TTY", tt.sshTTY)
			os.Setenv("SSH_CLIENT", tt.sshClient)
			os.Setenv("SSH_CONNECTION", tt.sshConn)

			if got := detector.IsSSH(); got != tt.wantSSH {
				t.Errorf("IsSSH() = %v, want %v", got, tt.wantSSH)
			}
		})
	}
}

func TestDetector_ShouldUseOSC52(t *testing.T) {
	detector := NewDetector()

	// Save original values
	origSSHTTY := os.Getenv("SSH_TTY")
	defer os.Setenv("SSH_TTY", origSSHTTY)

	tests := []struct {
		name   string
		sshTTY string
		want   bool
	}{
		{"SSH session", "/dev/pts/0", true},
		{"no SSH", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("SSH_TTY", tt.sshTTY)

			// ShouldUseOSC52 also checks HasTTY, which depends on stdout
			// so we can't reliably test it in all environments
			// Just verify it doesn't panic
			_ = detector.ShouldUseOSC52()
		})
	}
}

func TestDetector_HasTTY(t *testing.T) {
	detector := NewDetector()

	// Just verify it doesn't panic
	// Actual result depends on test environment
	_ = detector.HasTTY()
}

func TestDetector_IsHeadless(t *testing.T) {
	detector := NewDetector()

	// Save original values
	origDisplay := os.Getenv("DISPLAY")
	origWaylandDisplay := os.Getenv("WAYLAND_DISPLAY")

	defer func() {
		os.Setenv("DISPLAY", origDisplay)
		os.Setenv("WAYLAND_DISPLAY", origWaylandDisplay)
	}()

	if detector.GetPlatform() != TypeLinux {
		// Test non-Linux platforms (should return false)
		t.Run("non-Linux platforms", func(t *testing.T) {
			if detector.IsHeadless() {
				t.Error("IsHeadless() should return false on non-Linux platforms")
			}
		})
		return
	}

	tests := []struct {
		name           string
		display        string
		waylandDisplay string
		want           bool
	}{
		{"X11 display", ":0", "", false},
		{"Wayland display", "", "wayland-0", false},
		{"both displays", ":0", "wayland-0", false},
		{"no display", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("DISPLAY", tt.display)
			os.Setenv("WAYLAND_DISPLAY", tt.waylandDisplay)

			if got := detector.IsHeadless(); got != tt.want {
				t.Errorf("IsHeadless() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetector_GetPlatform_AllCases(t *testing.T) {
	detector := NewDetector()
	platform := detector.GetPlatform()

	// Verify it returns one of the expected types
	validTypes := []Type{TypeWindows, TypeDarwin, TypeLinux, TypeUnknown}
	valid := false
	for _, vt := range validTypes {
		if platform == vt {
			valid = true
			break
		}
	}

	if !valid {
		t.Errorf("GetPlatform() returned unexpected type: %s", platform)
	}

	// Verify it matches current runtime
	expected := Type(runtime.GOOS)
	switch runtime.GOOS {
	case "windows":
		expected = TypeWindows
	case "darwin":
		expected = TypeDarwin
	case "linux":
		expected = TypeLinux
	default:
		expected = TypeUnknown
	}

	if platform != expected {
		t.Errorf("GetPlatform() = %s, want %s", platform, expected)
	}
}

func TestDetector_HasTTY_ErrorHandling(t *testing.T) {
	detector := NewDetector()

	// Just verify it doesn't panic even if Stat fails
	// The actual result depends on whether stdout is a TTY
	hasTTY := detector.HasTTY()

	// Log result for debugging
	t.Logf("HasTTY() = %v", hasTTY)
}
