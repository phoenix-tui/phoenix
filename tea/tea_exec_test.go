package tea_test

import (
	"os/exec"
	"runtime"
	"testing"

	"github.com/phoenix-tui/phoenix/tea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	phoenixtesting "github.com/phoenix-tui/phoenix/testing"
)

// TestProgram_ExecProcess_APIWrapper verifies API wrapper calls internal program.ExecProcess.
func TestProgram_ExecProcess_APIWrapper(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{value: 0}
	p := tea.New(m, tea.WithTerminal[TestModel](mockTerm))

	// Create simple command.
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo", "API test")
	} else {
		cmd = exec.Command("echo", "API test")
	}

	// Execute via API.
	err := p.ExecProcess(cmd)

	// Should succeed.
	assert.NoError(t, err, "API ExecProcess should succeed")

	// Verify terminal operations were called (proves wrapper works).
	assert.Greater(t, mockTerm.CallCount("ShowCursor"), 0, "ShowCursor should be called via API wrapper")
}

// TestProgram_ExecProcess_API_NilCommand verifies API wrapper error handling.
func TestProgram_ExecProcess_API_NilCommand(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{value: 0}
	p := tea.New(m, tea.WithTerminal[TestModel](mockTerm))

	// Execute nil command via API.
	err := p.ExecProcess(nil)

	// Should return error.
	assert.Error(t, err, "API ExecProcess should return error for nil command")
	assert.Contains(t, err.Error(), "cmd is nil", "error message should be propagated from internal implementation")
}

// TestProgram_ExecProcess_API_CommandFailure verifies API error propagation.
func TestProgram_ExecProcess_API_CommandFailure(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{value: 0}
	p := tea.New(m, tea.WithTerminal[TestModel](mockTerm))

	// Create failing command.
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "exit", "1")
	} else {
		cmd = exec.Command("sh", "-c", "exit 1")
	}

	// Execute via API.
	err := p.ExecProcess(cmd)

	// Should propagate error.
	assert.Error(t, err, "API ExecProcess should propagate command error")
}

// TestProgram_ExecProcess_API_WithAltScreen verifies API alt screen handling.
func TestProgram_ExecProcess_API_WithAltScreen(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{value: 0}
	p := tea.New(m, tea.WithTerminal[TestModel](mockTerm), tea.WithAltScreen[TestModel]())

	// Enter alt screen.
	err := mockTerm.EnterAltScreen()
	require.NoError(t, err)
	mockTerm.Reset()

	// Create simple command.
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo", "test")
	} else {
		cmd = exec.Command("echo", "test")
	}

	// Execute via API.
	err = p.ExecProcess(cmd)
	assert.NoError(t, err)

	// Verify alt screen operations via API.
	assert.Greater(t, mockTerm.CallCount("ExitAltScreen"), 0, "ExitAltScreen via API")
	assert.Greater(t, mockTerm.CallCount("EnterAltScreen"), 0, "EnterAltScreen via API")
}

// TestWithTerminal_Option verifies WithTerminal option sets terminal correctly.
func TestWithTerminal_Option(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{value: 0}
	p := tea.New(m, tea.WithTerminal[TestModel](mockTerm))

	// Program should use the provided mock terminal.
	// We can verify this by calling a terminal operation via program.

	// Create simple command.
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo", "test")
	} else {
		cmd = exec.Command("echo", "test")
	}

	// Execute command (uses terminal from WithTerminal option).
	err := p.ExecProcess(cmd)
	assert.NoError(t, err)

	// Verify mock terminal was used (calls were recorded).
	assert.Greater(t, len(mockTerm.Calls), 0, "mock terminal should have recorded calls")
}

// TestAPI_ExecProcess_Sequential verifies multiple sequential API calls.
func TestAPI_ExecProcess_Sequential(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{value: 0}
	p := tea.New(m, tea.WithTerminal[TestModel](mockTerm))

	// Run 3 commands via API.
	for i := 0; i < 3; i++ {
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/c", "echo", "API iteration")
		} else {
			cmd = exec.Command("echo", "API iteration")
		}

		err := p.ExecProcess(cmd)
		assert.NoError(t, err, "API call %d should succeed", i)
	}

	// Verify multiple commands executed successfully.
	// Each command should call ShowCursor/HideCursor.
	cursorCalls := mockTerm.CallCount("ShowCursor") + mockTerm.CallCount("HideCursor")
	assert.Greater(t, cursorCalls, 0, "cursor operations should occur for all API calls")
}

// TestAPI_ExecProcess_Integration verifies end-to-end API usage.
func TestAPI_ExecProcess_Integration(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{value: 0}

	// Create program with multiple options (realistic usage).
	p := tea.New(
		m,
		tea.WithTerminal[TestModel](mockTerm),
		tea.WithAltScreen[TestModel](),
	)

	// Enter alt screen (simulate TUI running).
	err := mockTerm.EnterAltScreen()
	require.NoError(t, err)

	// Reset to track only ExecProcess operations.
	mockTerm.Reset()

	// Execute interactive command via API.
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo", "Integration test")
	} else {
		cmd = exec.Command("echo", "Integration test")
	}

	err = p.ExecProcess(cmd)
	assert.NoError(t, err, "Integration test should succeed")

	// Verify complete flow:
	// 1. Exit alt screen
	// 2. Show cursor
	// 3. (command runs)
	// 4. Re-enter alt screen
	// 5. Hide cursor

	assert.Greater(t, mockTerm.CallCount("ExitAltScreen"), 0, "step 1: exit alt screen")
	assert.Greater(t, mockTerm.CallCount("ShowCursor"), 0, "step 2: show cursor")
	assert.Greater(t, mockTerm.CallCount("EnterAltScreen"), 0, "step 4: re-enter alt screen")
	assert.Greater(t, mockTerm.CallCount("HideCursor"), 0, "step 5: hide cursor")

	// Verify final state: back in alt screen.
	assert.True(t, mockTerm.IsInAltScreen(), "should be back in alt screen after ExecProcess")
}

// TestAPI_ExecProcess_ErrorRecovery verifies API error recovery.
func TestAPI_ExecProcess_ErrorRecovery(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{value: 0}
	p := tea.New(m, tea.WithTerminal[TestModel](mockTerm), tea.WithAltScreen[TestModel]())

	// Enter alt screen.
	err := mockTerm.EnterAltScreen()
	require.NoError(t, err)
	mockTerm.Reset()

	// Run failing command.
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "exit", "1")
	} else {
		cmd = exec.Command("sh", "-c", "exit 1")
	}

	err = p.ExecProcess(cmd)

	// Error should be returned.
	assert.Error(t, err, "failing command should return error")

	// TUI state should be restored despite error.
	assert.Greater(t, mockTerm.CallCount("ExitAltScreen"), 0, "alt screen exited")
	assert.Greater(t, mockTerm.CallCount("EnterAltScreen"), 0, "alt screen restored")
	assert.Greater(t, mockTerm.CallCount("ShowCursor"), 0, "cursor shown")
	assert.Greater(t, mockTerm.CallCount("HideCursor"), 0, "cursor hidden")

	// Should be back in alt screen.
	assert.True(t, mockTerm.IsInAltScreen(), "should recover to alt screen state")
}

// TestAPI_ExecProcess_AutoTerminal verifies API auto-creates terminal.
func TestAPI_ExecProcess_AutoTerminal(t *testing.T) {
	m := TestModel{value: 0}

	// Create program WITHOUT explicit terminal.
	p := tea.New(m)

	// Execute command (should auto-create terminal).
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo", "auto terminal test")
	} else {
		cmd = exec.Command("echo", "auto terminal test")
	}

	err := p.ExecProcess(cmd)

	// Should succeed (terminal auto-created internally).
	assert.NoError(t, err, "API should auto-create terminal if not provided")
}
