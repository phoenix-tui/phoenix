package program

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	phoenixtesting "github.com/phoenix-tui/phoenix/testing"
)

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Suspend Tests                                                   │.
// └─────────────────────────────────────────────────────────────────┘.

// TestProgram_Suspend_Basic verifies basic Suspend functionality.
func TestProgram_Suspend_Basic(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Enter raw mode first (simulating running TUI).
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)

	// Suspend.
	err = p.Suspend()
	assert.NoError(t, err, "Suspend should succeed")

	// Verify suspended state.
	assert.True(t, p.IsSuspended(), "Should be suspended after Suspend")

	// Verify terminal restored to cooked mode.
	assert.False(t, mockTerm.IsInRawMode(), "Should exit raw mode")
	assert.Greater(t, mockTerm.CallCount("ShowCursor"), 0, "ShowCursor should be called during Suspend")
}

// TestProgram_Suspend_Idempotent verifies Suspend is idempotent.
func TestProgram_Suspend_Idempotent(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Enter raw mode.
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)

	// First Suspend.
	err = p.Suspend()
	assert.NoError(t, err)
	assert.True(t, p.IsSuspended())

	// Second Suspend should be no-op.
	err = p.Suspend()
	assert.NoError(t, err, "Second Suspend should succeed (no-op)")
	assert.True(t, p.IsSuspended(), "Should still be suspended")
}

// TestProgram_Suspend_WithAltScreen verifies Suspend exits alt screen.
func TestProgram_Suspend_WithAltScreen(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm), WithAltScreen[TestModel]())

	// Enter raw mode and alt screen.
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)
	err = mockTerm.EnterAltScreen()
	require.NoError(t, err)

	// Suspend.
	err = p.Suspend()
	assert.NoError(t, err)

	// Verify both raw mode and alt screen exited.
	assert.False(t, mockTerm.IsInRawMode(), "Should exit raw mode")
	assert.False(t, mockTerm.IsInAltScreen(), "Should exit alt screen")
	assert.True(t, p.IsSuspended())
}

// TestProgram_Suspend_NoRawMode verifies Suspend works without raw mode.
func TestProgram_Suspend_NoRawMode(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// No raw mode - should still work.
	err := p.Suspend()
	assert.NoError(t, err, "Suspend should succeed without raw mode")
	assert.True(t, p.IsSuspended())
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Resume Tests                                                    │.
// └─────────────────────────────────────────────────────────────────┘.

// TestProgram_Resume_Basic verifies basic Resume functionality.
func TestProgram_Resume_Basic(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Setup: Enter raw mode and then suspend.
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)
	err = p.Suspend()
	require.NoError(t, err)

	// Resume.
	err = p.Resume()
	assert.NoError(t, err, "Resume should succeed")

	// Verify no longer suspended.
	assert.False(t, p.IsSuspended(), "Should not be suspended after Resume")

	// Verify raw mode restored.
	assert.True(t, mockTerm.IsInRawMode(), "Raw mode should be restored")
}

// TestProgram_Resume_Idempotent verifies Resume is idempotent.
func TestProgram_Resume_Idempotent(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Not suspended - Resume should be no-op.
	err := p.Resume()
	assert.NoError(t, err, "Resume should succeed when not suspended (no-op)")
	assert.False(t, p.IsSuspended())

	// Suspend and Resume.
	err = mockTerm.EnterRawMode()
	require.NoError(t, err)
	err = p.Suspend()
	require.NoError(t, err)
	err = p.Resume()
	require.NoError(t, err)

	// Second Resume should be no-op.
	err = p.Resume()
	assert.NoError(t, err, "Second Resume should succeed (no-op)")
}

// TestProgram_Resume_WithAltScreen verifies Resume restores alt screen.
func TestProgram_Resume_WithAltScreen(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm), WithAltScreen[TestModel]())

	// Enter raw mode and alt screen, then suspend.
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)
	err = mockTerm.EnterAltScreen()
	require.NoError(t, err)
	err = p.Suspend()
	require.NoError(t, err)

	// Verify exited.
	assert.False(t, mockTerm.IsInRawMode())
	assert.False(t, mockTerm.IsInAltScreen())

	// Resume.
	err = p.Resume()
	assert.NoError(t, err)

	// Verify both restored.
	assert.True(t, mockTerm.IsInRawMode(), "Raw mode should be restored")
	assert.True(t, mockTerm.IsInAltScreen(), "Alt screen should be restored")
}

// TestProgram_Resume_CursorHidden verifies Resume hides cursor.
func TestProgram_Resume_CursorHidden(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Enter raw mode and suspend.
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)
	err = p.Suspend()
	require.NoError(t, err)

	// ShowCursor should be called after suspend.
	assert.Greater(t, mockTerm.CallCount("ShowCursor"), 0, "ShowCursor should be called during Suspend")

	// Reset to track Resume calls only.
	mockTerm.Reset()

	// Resume.
	err = p.Resume()
	require.NoError(t, err)

	// HideCursor should be called after resume (TUI state).
	assert.Greater(t, mockTerm.CallCount("HideCursor"), 0, "HideCursor should be called during Resume")
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Suspend/Resume Cycle Tests                                      │.
// └─────────────────────────────────────────────────────────────────┘.

// TestProgram_SuspendResume_Cycle verifies full suspend/resume cycle.
func TestProgram_SuspendResume_Cycle(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm), WithAltScreen[TestModel]())

	// Setup initial state.
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)
	err = mockTerm.EnterAltScreen()
	require.NoError(t, err)

	// Initial state check.
	assert.True(t, mockTerm.IsInRawMode())
	assert.True(t, mockTerm.IsInAltScreen())
	assert.False(t, p.IsSuspended())

	// Suspend.
	err = p.Suspend()
	require.NoError(t, err)
	assert.False(t, mockTerm.IsInRawMode(), "Raw mode should be exited")
	assert.False(t, mockTerm.IsInAltScreen(), "Alt screen should be exited")
	assert.True(t, p.IsSuspended())

	// Resume.
	err = p.Resume()
	require.NoError(t, err)
	assert.True(t, mockTerm.IsInRawMode(), "Raw mode should be restored")
	assert.True(t, mockTerm.IsInAltScreen(), "Alt screen should be restored")
	assert.False(t, p.IsSuspended())
}

// TestProgram_SuspendResume_MultipleCycles verifies multiple suspend/resume cycles.
func TestProgram_SuspendResume_MultipleCycles(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Enter raw mode.
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)

	// Multiple cycles.
	for i := 0; i < 3; i++ {
		// Suspend.
		err = p.Suspend()
		require.NoError(t, err, "Suspend cycle %d should succeed", i)
		assert.True(t, p.IsSuspended())
		assert.False(t, mockTerm.IsInRawMode())

		// Resume.
		err = p.Resume()
		require.NoError(t, err, "Resume cycle %d should succeed", i)
		assert.False(t, p.IsSuspended())
		assert.True(t, mockTerm.IsInRawMode())
	}
}

// TestProgram_SuspendResume_StatePreserved verifies state is correctly saved and restored.
func TestProgram_SuspendResume_StatePreserved(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Test 1: Raw mode only.
	err := mockTerm.EnterRawMode()
	require.NoError(t, err)

	err = p.Suspend()
	require.NoError(t, err)
	err = p.Resume()
	require.NoError(t, err)

	assert.True(t, mockTerm.IsInRawMode(), "Raw mode should be restored")
	assert.False(t, mockTerm.IsInAltScreen(), "Alt screen should NOT be restored (wasn't active)")

	// Test 2: Raw mode + Alt screen.
	err = mockTerm.EnterAltScreen()
	require.NoError(t, err)

	err = p.Suspend()
	require.NoError(t, err)
	err = p.Resume()
	require.NoError(t, err)

	assert.True(t, mockTerm.IsInRawMode(), "Raw mode should be restored")
	assert.True(t, mockTerm.IsInAltScreen(), "Alt screen should be restored")
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ IsSuspended Tests                                               │.
// └─────────────────────────────────────────────────────────────────┘.

// TestProgram_IsSuspended_Initial verifies initial suspended state.
func TestProgram_IsSuspended_Initial(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	assert.False(t, p.IsSuspended(), "Should not be suspended initially")
}

// TestProgram_IsSuspended_AfterSuspend verifies suspended state after Suspend.
func TestProgram_IsSuspended_AfterSuspend(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	err := p.Suspend()
	require.NoError(t, err)

	assert.True(t, p.IsSuspended(), "Should be suspended after Suspend")
}

// TestProgram_IsSuspended_AfterResume verifies suspended state after Resume.
func TestProgram_IsSuspended_AfterResume(t *testing.T) {
	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	err := p.Suspend()
	require.NoError(t, err)
	err = p.Resume()
	require.NoError(t, err)

	assert.False(t, p.IsSuspended(), "Should not be suspended after Resume")
}

// ┌─────────────────────────────────────────────────────────────────┐.
// │ Integration with InputReader Tests                              │.
// └─────────────────────────────────────────────────────────────────┘.

// TestProgram_Suspend_StopsInputReader verifies Suspend stops input reader.
func TestProgram_Suspend_StopsInputReader(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows: inputReader state timing is non-deterministic due to stdin blocking")
	}

	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Start input reader.
	p.startInputReader()

	// Verify running.
	p.mu.Lock()
	running := p.inputReaderRunning
	p.mu.Unlock()
	assert.True(t, running, "Input reader should be running")

	// Suspend stops input reader.
	err := p.Suspend()
	require.NoError(t, err)

	// Verify stopped.
	p.mu.Lock()
	running = p.inputReaderRunning
	p.mu.Unlock()
	assert.False(t, running, "Input reader should be stopped after Suspend")
}

// TestProgram_Resume_RestartsInputReader verifies Resume restarts input reader.
func TestProgram_Resume_RestartsInputReader(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows: inputReader restart timing is non-deterministic due to stdin blocking")
	}

	mockTerm := phoenixtesting.NewMockTerminal()
	m := TestModel{}
	p := New(m, WithTerminal[TestModel](mockTerm))

	// Start input reader and suspend.
	p.startInputReader()
	err := p.Suspend()
	require.NoError(t, err)

	// Verify stopped.
	p.mu.Lock()
	running := p.inputReaderRunning
	p.mu.Unlock()
	assert.False(t, running)

	// Resume restarts input reader.
	err = p.Resume()
	require.NoError(t, err)

	// Verify restarted.
	p.mu.Lock()
	running = p.inputReaderRunning
	p.mu.Unlock()
	assert.True(t, running, "Input reader should be restarted after Resume")
}
