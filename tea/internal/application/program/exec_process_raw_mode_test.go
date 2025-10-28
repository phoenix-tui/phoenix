package program

import (
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	phoenixtesting "github.com/phoenix-tui/phoenix/testing"
)

// ┌─────────────────────────────────────────────────────────────────┐
// │ Raw Mode Management Tests (CRITICAL BUG FIX)                    │
// └─────────────────────────────────────────────────────────────────┘

// TestProgram_ExecProcess_RawModeExited verifies raw mode is exited before command.
func TestProgram_ExecProcess_RawModeExited(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Enter raw mode manually to simulate TUI state
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)
	assert.True(t, mockTerm.IsInRawMode(), "should be in raw mode initially")

	// Reset calls to track ExecProcess operations only
	mockTerm.Reset()

	// Create simple command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo", "test")
	} else {
		cmd = exec.Command("echo", "test")
	}

	// Execute command
	err = p.ExecProcess(cmd)
	assert.NoError(t, err)

	// Verify raw mode was exited before command
	assert.Greater(t, mockTerm.CallCount("ExitRawMode"), 0, "ExitRawMode should be called")
}

// TestProgram_ExecProcess_RawModeRestored verifies raw mode is restored after command.
func TestProgram_ExecProcess_RawModeRestored(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Enter raw mode manually to simulate TUI state
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)

	// Reset calls
	mockTerm.Reset()

	// Create simple command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo", "test")
	} else {
		cmd = exec.Command("echo", "test")
	}

	// Execute command
	err = p.ExecProcess(cmd)
	assert.NoError(t, err)

	// Verify raw mode was exited and re-entered
	assert.Greater(t, mockTerm.CallCount("ExitRawMode"), 0, "ExitRawMode should be called")
	assert.Greater(t, mockTerm.CallCount("EnterRawMode"), 0, "EnterRawMode should be called to restore")

	// Verify we're back in raw mode after command
	assert.True(t, mockTerm.IsInRawMode(), "should be in raw mode after ExecProcess")
}

// TestProgram_ExecProcess_RawModeStatePreserved verifies raw mode state is preserved correctly.
func TestProgram_ExecProcess_RawModeStatePreserved(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Scenario 1: Start NOT in raw mode, should stay NOT in raw mode
	assert.False(t, mockTerm.IsInRawMode())

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo", "test")
	} else {
		cmd = exec.Command("echo", "test")
	}

	err := p.ExecProcess(cmd)
	assert.NoError(t, err)

	// Should remain NOT in raw mode after command
	assert.False(t, mockTerm.IsInRawMode(), "should not be in raw mode if not initially in it")

	// Scenario 2: Enter raw mode, run command, should return to raw mode
	err = mockTerm.EnterRawMode()
	require.NoError(t, err)
	assert.True(t, mockTerm.IsInRawMode())

	mockTerm.Reset() // Clear call history

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo", "test2")
	} else {
		cmd = exec.Command("echo", "test2")
	}

	err = p.ExecProcess(cmd)
	assert.NoError(t, err)

	// Should be back in raw mode after command
	assert.True(t, mockTerm.IsInRawMode(), "should return to raw mode after ExecProcess")
}

// TestProgram_ExecProcess_RawModeWithAltScreen verifies raw mode + alt screen interaction.
func TestProgram_ExecProcess_RawModeWithAltScreen(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm), WithAltScreen[TestModel]())

	// Enter both raw mode and alt screen (typical TUI state)
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)
	err = mockTerm.EnterAltScreen()
	require.NoError(t, err)

	assert.True(t, mockTerm.IsInRawMode())
	assert.True(t, mockTerm.IsInAltScreen())

	// Reset calls
	mockTerm.Reset()

	// Create simple command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo", "test")
	} else {
		cmd = exec.Command("echo", "test")
	}

	// Execute command
	err = p.ExecProcess(cmd)
	assert.NoError(t, err)

	// Verify BOTH raw mode and alt screen were managed
	assert.Greater(t, mockTerm.CallCount("ExitRawMode"), 0, "ExitRawMode called")
	assert.Greater(t, mockTerm.CallCount("ExitAltScreen"), 0, "ExitAltScreen called")
	assert.Greater(t, mockTerm.CallCount("EnterAltScreen"), 0, "EnterAltScreen called to restore")
	assert.Greater(t, mockTerm.CallCount("EnterRawMode"), 0, "EnterRawMode called to restore")

	// Verify both are restored after command
	assert.True(t, mockTerm.IsInRawMode(), "raw mode should be restored")
	assert.True(t, mockTerm.IsInAltScreen(), "alt screen should be restored")
}

// TestProgram_ExecProcess_RawModeRestoreOnError verifies raw mode restored even on command error.
func TestProgram_ExecProcess_RawModeRestoreOnError(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Enter raw mode
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)

	// Reset calls
	mockTerm.Reset()

	// Create command that fails
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "exit", "1")
	} else {
		cmd = exec.Command("sh", "-c", "exit 1")
	}

	// Execute command
	err = p.ExecProcess(cmd)

	// Command should fail
	assert.Error(t, err, "command should fail")

	// Raw mode should still be restored
	assert.Greater(t, mockTerm.CallCount("ExitRawMode"), 0, "ExitRawMode called")
	assert.Greater(t, mockTerm.CallCount("EnterRawMode"), 0, "EnterRawMode called to restore")
	assert.True(t, mockTerm.IsInRawMode(), "raw mode should be restored even on error")
}

// TestProgram_ExecProcess_RawModeOrder verifies correct order of operations.
func TestProgram_ExecProcess_RawModeOrder(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm), WithAltScreen[TestModel]())

	// Enter both raw mode and alt screen
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)
	err = mockTerm.EnterAltScreen()
	require.NoError(t, err)

	// Reset calls to track order
	mockTerm.Reset()

	// Create simple command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo", "test")
	} else {
		cmd = exec.Command("echo", "test")
	}

	// Execute command
	err = p.ExecProcess(cmd)
	assert.NoError(t, err)

	// Find positions of key operations
	var exitRawModePos, exitAltScreenPos, enterAltScreenPos, enterRawModePos int
	for i, call := range mockTerm.Calls {
		if call == "ExitRawMode" && exitRawModePos == 0 {
			exitRawModePos = i + 1
		}
		if call == "ExitAltScreen" && exitAltScreenPos == 0 {
			exitAltScreenPos = i + 1
		}
		if call == "EnterAltScreen" && enterAltScreenPos == 0 {
			enterAltScreenPos = i + 1
		}
		if call == "EnterRawMode" && enterRawModePos == 0 {
			enterRawModePos = i + 1
		}
	}

	// Verify order (before command):
	// 1. ExitRawMode (restore cooked mode for external command)
	// 2. ExitAltScreen (show normal screen)
	assert.Greater(t, exitRawModePos, 0, "ExitRawMode should be called")
	assert.Greater(t, exitAltScreenPos, 0, "ExitAltScreen should be called")
	assert.Less(t, exitRawModePos, exitAltScreenPos, "ExitRawMode should come BEFORE ExitAltScreen")

	// Verify order (after command):
	// 1. EnterAltScreen (restore TUI screen)
	// 2. EnterRawMode (restore TUI input mode)
	assert.Greater(t, enterAltScreenPos, 0, "EnterAltScreen should be called")
	assert.Greater(t, enterRawModePos, 0, "EnterRawMode should be called")
	assert.Less(t, enterAltScreenPos, enterRawModePos, "EnterAltScreen should come BEFORE EnterRawMode")
}

// TestProgram_ExecProcess_NoRawModeWithoutEnter verifies no raw mode exit if not entered.
func TestProgram_ExecProcess_NoRawModeWithoutEnter(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Ensure NOT in raw mode
	assert.False(t, mockTerm.IsInRawMode())

	// Reset calls
	mockTerm.Reset()

	// Create simple command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo", "test")
	} else {
		cmd = exec.Command("echo", "test")
	}

	// Execute command
	err := p.ExecProcess(cmd)
	assert.NoError(t, err)

	// Verify no raw mode operations (terminal not in raw mode initially)
	// Note: IsInRawMode is called to check state, but Exit/Enter should not be called
	assert.Equal(t, 0, mockTerm.CallCount("ExitRawMode"), "ExitRawMode should NOT be called")
	assert.Equal(t, 0, mockTerm.CallCount("EnterRawMode"), "EnterRawMode should NOT be called")
}

// TestProgram_ExecProcess_RawModeSequentialCalls verifies multiple sequential calls work.
func TestProgram_ExecProcess_RawModeSequentialCalls(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Enter raw mode once
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)

	// Run 3 commands sequentially
	for i := 0; i < 3; i++ {
		mockTerm.Reset()

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/c", "echo", "iteration")
		} else {
			cmd = exec.Command("echo", "iteration")
		}

		err := p.ExecProcess(cmd)
		assert.NoError(t, err, "iteration %d should succeed", i)

		// Each call should exit and re-enter raw mode
		assert.Greater(t, mockTerm.CallCount("ExitRawMode"), 0, "iteration %d: ExitRawMode", i)
		assert.Greater(t, mockTerm.CallCount("EnterRawMode"), 0, "iteration %d: EnterRawMode", i)

		// Should be in raw mode after each iteration
		assert.True(t, mockTerm.IsInRawMode(), "iteration %d: should be in raw mode", i)
	}
}
