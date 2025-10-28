//go:build windows

package windows

import (
	"os"
	"testing"

	"golang.org/x/term"
)

// ┌─────────────────────────────────────────────────────────────────┐
// │ Raw Mode Tests (Windows)                                        │
// └─────────────────────────────────────────────────────────────────┘

// Note: Raw mode tests require actual terminal (can't use pipes).
// These tests check the state tracking logic, not actual terminal behavior.

func TestConsole_RawMode_Lifecycle(t *testing.T) {
	// Create console (uses real stdin/stdout)
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Initially not in raw mode
	if console.IsInRawMode() {
		t.Error("Console should not be in raw mode initially")
	}

	// Skip actual raw mode test if not in a terminal (CI environment)
	if !isTerminal(os.Stdin) {
		t.Skip("Skipping raw mode test - not running in a terminal")
		return
	}

	// Enter raw mode
	err = console.EnterRawMode()
	if err != nil {
		t.Fatalf("EnterRawMode failed: %v", err)
	}

	// Should be in raw mode now
	if !console.IsInRawMode() {
		t.Error("Console should be in raw mode after EnterRawMode")
	}

	// Exit raw mode
	err = console.ExitRawMode()
	if err != nil {
		t.Fatalf("ExitRawMode failed: %v", err)
	}

	// Should not be in raw mode anymore
	if console.IsInRawMode() {
		t.Error("Console should not be in raw mode after ExitRawMode")
	}
}

func TestConsole_RawMode_DoubleEnter(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Skip if not in a terminal
	if !isTerminal(os.Stdin) {
		t.Skip("Skipping raw mode test - not running in a terminal")
		return
	}

	// Enter raw mode
	err = console.EnterRawMode()
	if err != nil {
		t.Fatalf("First EnterRawMode failed: %v", err)
	}
	defer console.ExitRawMode() // Cleanup

	// Second enter should fail
	err = console.EnterRawMode()
	if err == nil {
		t.Error("Second EnterRawMode should fail")
	}

	expectedMsg := "already in raw mode"
	if err != nil && err.Error() != expectedMsg {
		t.Errorf("Expected error %q, got %q", expectedMsg, err.Error())
	}
}

func TestConsole_RawMode_ExitWithoutEnter(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Exit without enter should fail
	err = console.ExitRawMode()
	if err == nil {
		t.Error("ExitRawMode without EnterRawMode should fail")
	}

	expectedMsg := "not in raw mode"
	if err != nil && err.Error() != expectedMsg {
		t.Errorf("Expected error %q, got %q", expectedMsg, err.Error())
	}
}

func TestConsole_RawMode_StateTracking(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Initially not in raw mode
	if console.IsInRawMode() {
		t.Error("IsInRawMode should return false initially")
	}

	// Manually set state (simulating successful EnterRawMode)
	console.inRawMode = true

	// Should report as in raw mode
	if !console.IsInRawMode() {
		t.Error("IsInRawMode should return true after setting state")
	}

	// Manually clear state (simulating successful ExitRawMode)
	console.inRawMode = false

	// Should report as not in raw mode
	if console.IsInRawMode() {
		t.Error("IsInRawMode should return false after clearing state")
	}
}

// isTerminal checks if file descriptor is a terminal.
// Used to skip raw mode tests in CI environments.
func isTerminal(f *os.File) bool {
	// Use golang.org/x/term.IsTerminal for robust check
	return term.IsTerminal(int(f.Fd()))
}
