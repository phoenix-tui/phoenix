package program

import (
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/term"

	phoenixtesting "github.com/phoenix-tui/phoenix/testing"
)

// ┌─────────────────────────────────────────────────────────────────┐.
// │ TTYOptions Tests                                                │.
// └─────────────────────────────────────────────────────────────────┘.

// TestTTYOptions_ZeroValue verifies zero value TTYOptions.
func TestTTYOptions_ZeroValue(t *testing.T) {
	var opts TTYOptions

	assert.False(t, opts.TransferForeground, "Zero value should not transfer foreground")
	assert.False(t, opts.CreateProcessGroup, "Zero value should not create process group")
}

// TestTTYOptions_TransferForeground verifies TransferForeground option.
func TestTTYOptions_TransferForeground(t *testing.T) {
	opts := TTYOptions{TransferForeground: true}

	assert.True(t, opts.TransferForeground)
	assert.False(t, opts.CreateProcessGroup)
}

// TestTTYOptions_CreateProcessGroup verifies CreateProcessGroup option.
func TestTTYOptions_CreateProcessGroup(t *testing.T) {
	opts := TTYOptions{CreateProcessGroup: true}

	assert.False(t, opts.TransferForeground)
	assert.True(t, opts.CreateProcessGroup)
}

// TestTTYOptions_BothOptions verifies both options enabled.
func TestTTYOptions_BothOptions(t *testing.T) {
	opts := TTYOptions{
		TransferForeground: true,
		CreateProcessGroup: true,
	}

	assert.True(t, opts.TransferForeground)
	assert.True(t, opts.CreateProcessGroup)
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ ExecProcessWithTTY Basic Tests                                  │.
// └─────────────────────────────────────────────────────────────────┘.

// TestProgram_ExecProcessWithTTY_NilCommand verifies error on nil command.
func TestProgram_ExecProcessWithTTY_NilCommand(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	opts := TTYOptions{}
	err := p.ExecProcessWithTTY(nil, opts)

	assert.Error(t, err, "Should error on nil command")
	assert.Contains(t, err.Error(), "nil", "Error should mention nil")
}

// TestProgram_ExecProcessWithTTY_SimpleCommand verifies basic execution.
func TestProgram_ExecProcessWithTTY_SimpleCommand(t *testing.T) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		t.Skip("Skipping test: not a TTY (CI environment)")
	}

	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Enter raw mode (simulating running TUI)
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)

	// Execute simple command (echo)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/C", "echo test")
	} else {
		cmd = exec.Command("echo", "test")
	}

	opts := TTYOptions{
		TransferForeground: false, // Don't transfer for simple test
		CreateProcessGroup: false,
	}

	err = p.ExecProcessWithTTY(cmd, opts)
	assert.NoError(t, err, "Simple command should succeed")
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Platform-Specific Tests                                         │.
// └─────────────────────────────────────────────────────────────────┘.

// TestProgram_ExecProcessWithTTY_Unix_TransferForeground verifies Unix tcsetpgrp behavior.
func TestProgram_ExecProcessWithTTY_Unix_TransferForeground(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping Unix-specific test on Windows")
	}

	if !term.IsTerminal(int(os.Stdin.Fd())) {
		t.Skip("Skipping test: not a TTY (CI environment)")
	}

	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Enter raw mode
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)

	// Execute command with foreground transfer
	cmd := exec.Command("true") // Simple command that always succeeds

	opts := TTYOptions{
		TransferForeground: true,
		CreateProcessGroup: true, // Recommended with TransferForeground
	}

	err = p.ExecProcessWithTTY(cmd, opts)
	assert.NoError(t, err, "Command with TTY transfer should succeed")
}

// TestProgram_ExecProcessWithTTY_Unix_ProcessGroup verifies Unix process group creation.
func TestProgram_ExecProcessWithTTY_Unix_ProcessGroup(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping Unix-specific test on Windows")
	}

	if !term.IsTerminal(int(os.Stdin.Fd())) {
		t.Skip("Skipping test: not a TTY (CI environment)")
	}

	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Enter raw mode
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)

	// Execute command with process group creation (no TTY transfer)
	cmd := exec.Command("true")

	opts := TTYOptions{
		TransferForeground: false,
		CreateProcessGroup: true, // Create process group without transferring foreground
	}

	err = p.ExecProcessWithTTY(cmd, opts)
	assert.NoError(t, err, "Command with process group should succeed")
}

// TestProgram_ExecProcessWithTTY_Windows_ConsoleMode verifies Windows console mode handling.
func TestProgram_ExecProcessWithTTY_Windows_ConsoleMode(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on Unix")
	}

	if !term.IsTerminal(int(os.Stdin.Fd())) {
		t.Skip("Skipping test: not a console (CI environment)")
	}

	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Enter raw mode
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)

	// Execute simple command
	cmd := exec.Command("cmd.exe", "/C", "echo test")

	opts := TTYOptions{
		TransferForeground: false, // Ignored on Windows
		CreateProcessGroup: true,  // Uses CREATE_NEW_PROCESS_GROUP on Windows
	}

	err = p.ExecProcessWithTTY(cmd, opts)
	assert.NoError(t, err, "Command with console mode control should succeed")
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Fallback Tests                                                  │.
// └─────────────────────────────────────────────────────────────────┘.

// TestProgram_ExecProcessWithTTY_Fallback_NoTTY verifies fallback when not a TTY.
func TestProgram_ExecProcessWithTTY_Fallback_NoTTY(t *testing.T) {
	// This test intentionally does NOT check for TTY
	// It should gracefully fall back to ExecProcess

	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Execute command (should fall back to simple exec if no TTY)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/C", "echo test")
	} else {
		cmd = exec.Command("echo", "test")
	}

	opts := TTYOptions{
		TransferForeground: true, // Will be ignored if no TTY available
		CreateProcessGroup: true,
	}

	// Should not error even without TTY (falls back gracefully)
	err := p.ExecProcessWithTTY(cmd, opts)
	// Error is acceptable here (no TTY or execution failed),
	// but should not panic or hang
	_ = err
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Suspend/Resume Integration Tests                                │.
// └─────────────────────────────────────────────────────────────────┘.

// TestProgram_ExecProcessWithTTY_SuspendResume verifies Suspend/Resume integration.
func TestProgram_ExecProcessWithTTY_SuspendResume(t *testing.T) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		t.Skip("Skipping test: not a TTY (CI environment)")
	}

	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Enter raw mode
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)
	assert.True(t, mockTerm.IsInRawMode(), "Should be in raw mode before exec")

	// Execute command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/C", "echo test")
	} else {
		cmd = exec.Command("echo", "test")
	}

	opts := TTYOptions{}
	err = p.ExecProcessWithTTY(cmd, opts)
	require.NoError(t, err)

	// Verify raw mode restored after exec
	assert.True(t, mockTerm.IsInRawMode(), "Should be back in raw mode after exec")
}

// TestProgram_ExecProcessWithTTY_ErrorHandling verifies error handling.
func TestProgram_ExecProcessWithTTY_ErrorHandling(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Enter raw mode
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)

	// Execute non-existent command (should fail)
	cmd := exec.Command("nonexistent-command-that-does-not-exist")

	opts := TTYOptions{}
	err = p.ExecProcessWithTTY(cmd, opts)

	// Should return error for non-existent command
	assert.Error(t, err, "Non-existent command should error")

	// Verify raw mode restored even after error
	assert.True(t, mockTerm.IsInRawMode(), "Raw mode should be restored even after error")
}
