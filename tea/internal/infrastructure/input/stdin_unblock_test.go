package input

import (
	"runtime"
	"testing"
)

// TestUnblockStdinRead_NoError verifies UnblockStdinRead doesn't return errors.
// This is platform-specific but should work on all platforms.
func TestUnblockStdinRead_NoError(t *testing.T) {
	err := UnblockStdinRead()
	if err != nil {
		t.Errorf("UnblockStdinRead() returned error: %v", err)
	}
}

// TestUnblockStdinRead_MultipleCalls verifies it's safe to call multiple times.
// Important for Cancel() which might be called multiple times.
func TestUnblockStdinRead_MultipleCalls(t *testing.T) {
	// First call
	err := UnblockStdinRead()
	if err != nil {
		t.Errorf("First UnblockStdinRead() returned error: %v", err)
	}

	// Second call - should also succeed
	err = UnblockStdinRead()
	if err != nil {
		t.Errorf("Second UnblockStdinRead() returned error: %v", err)
	}

	// Third call - verify idempotency
	err = UnblockStdinRead()
	if err != nil {
		t.Errorf("Third UnblockStdinRead() returned error: %v", err)
	}
}

// TestUnblockStdinRead_WindowsOnly verifies Windows-specific behavior.
// This test only runs on Windows and ensures WriteConsoleInputW works.
func TestUnblockStdinRead_WindowsOnly(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	// On Windows, this should inject a key event successfully
	err := UnblockStdinRead()
	if err != nil {
		t.Errorf("Windows UnblockStdinRead() failed: %v", err)
	}

	// Verify it doesn't panic or cause issues when called multiple times
	for i := 0; i < 10; i++ {
		err := UnblockStdinRead()
		if err != nil {
			t.Errorf("Windows UnblockStdinRead() call %d failed: %v", i+1, err)
		}
	}
}

// TestUnblockStdinRead_NonWindowsNoOp verifies non-Windows platforms return nil.
// This ensures the build tag separation works correctly.
func TestUnblockStdinRead_NonWindowsNoOp(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Non-Windows test")
	}

	// On non-Windows, this should be a no-op returning nil
	err := UnblockStdinRead()
	if err != nil {
		t.Errorf("Non-Windows UnblockStdinRead() should be no-op but returned: %v", err)
	}
}
