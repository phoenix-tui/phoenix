//go:build unix

package unix

import (
	"os"
	"testing"

	"golang.org/x/term"
)

// ┌─────────────────────────────────────────────────────────────────┐
// │ Raw Mode Tests                                                  │
// └─────────────────────────────────────────────────────────────────┘

// Note: Raw mode tests require actual terminal (can't use pipes).
// These tests check the state tracking logic, not actual terminal behavior.

func TestANSI_RawMode_Lifecycle(t *testing.T) {
	// Create terminal with real stdin (required for raw mode)
	term := &ANSITerminal{
		output: os.Stdout,
		input:  os.Stdin,
	}

	// Initially not in raw mode
	if term.IsInRawMode() {
		t.Error("Terminal should not be in raw mode initially")
	}

	// Skip actual raw mode test if not in a terminal (CI environment)
	if !isTerminal(os.Stdin) {
		t.Skip("Skipping raw mode test - not running in a terminal")
		return
	}

	// Enter raw mode
	err := term.EnterRawMode()
	if err != nil {
		t.Fatalf("EnterRawMode failed: %v", err)
	}

	// Should be in raw mode now
	if !term.IsInRawMode() {
		t.Error("Terminal should be in raw mode after EnterRawMode")
	}

	// Exit raw mode
	err = term.ExitRawMode()
	if err != nil {
		t.Fatalf("ExitRawMode failed: %v", err)
	}

	// Should not be in raw mode anymore
	if term.IsInRawMode() {
		t.Error("Terminal should not be in raw mode after ExitRawMode")
	}
}

func TestANSI_RawMode_DoubleEnter(t *testing.T) {
	term := &ANSITerminal{
		output: os.Stdout,
		input:  os.Stdin,
	}

	// Skip if not in a terminal
	if !isTerminal(os.Stdin) {
		t.Skip("Skipping raw mode test - not running in a terminal")
		return
	}

	// Enter raw mode
	err := term.EnterRawMode()
	if err != nil {
		t.Fatalf("First EnterRawMode failed: %v", err)
	}
	defer term.ExitRawMode() // Cleanup

	// Second enter should fail
	err = term.EnterRawMode()
	if err == nil {
		t.Error("Second EnterRawMode should fail")
	}

	expectedMsg := "already in raw mode"
	if err != nil && err.Error() != expectedMsg {
		t.Errorf("Expected error %q, got %q", expectedMsg, err.Error())
	}
}

func TestANSI_RawMode_ExitWithoutEnter(t *testing.T) {
	term := &ANSITerminal{
		output: os.Stdout,
		input:  os.Stdin,
	}

	// Exit without enter should fail
	err := term.ExitRawMode()
	if err == nil {
		t.Error("ExitRawMode without EnterRawMode should fail")
	}

	expectedMsg := "terminal: not in raw mode"
	if err != nil && err.Error() != expectedMsg {
		t.Errorf("Expected error %q, got %q", expectedMsg, err.Error())
	}
}

func TestANSI_RawMode_StateTracking(t *testing.T) {
	// Test state tracking without actual terminal
	term := &ANSITerminal{
		output:    os.Stdout,
		input:     os.Stdin,
		inRawMode: false, // Simulate not in raw mode
	}

	// Initially not in raw mode
	if term.IsInRawMode() {
		t.Error("IsInRawMode should return false initially")
	}

	// Manually set state (simulating successful EnterRawMode)
	term.inRawMode = true

	// Should report as in raw mode
	if !term.IsInRawMode() {
		t.Error("IsInRawMode should return true after setting state")
	}

	// Manually clear state (simulating successful ExitRawMode)
	term.inRawMode = false

	// Should report as not in raw mode
	if term.IsInRawMode() {
		t.Error("IsInRawMode should return false after clearing state")
	}
}

// isTerminal checks if file descriptor is a terminal.
// Used to skip raw mode tests in CI environments.
func isTerminal(f *os.File) bool {
	// Use golang.org/x/term.IsTerminal for robust check
	return term.IsTerminal(int(f.Fd()))
}
