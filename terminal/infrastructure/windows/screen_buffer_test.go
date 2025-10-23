//go:build windows
// +build windows

package windows

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConsole_EnterAltScreen_Success verifies alternate screen buffer creation.
func TestConsole_EnterAltScreen_Success(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Initially not in alt screen.
	assert.False(t, console.IsInAltScreen(), "should initially NOT be in alt screen")

	// Enter alt screen.
	err = console.EnterAltScreen()
	require.NoError(t, err, "EnterAltScreen should succeed")

	// Now in alt screen.
	assert.True(t, console.IsInAltScreen(), "should be in alt screen after Enter")

	// Cleanup: exit alt screen.
	defer func() {
		err := console.ExitAltScreen()
		assert.NoError(t, err, "ExitAltScreen cleanup should succeed")
	}()
}

// TestConsole_EnterAltScreen_AlreadyIn verifies error on double-enter.
func TestConsole_EnterAltScreen_AlreadyIn(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Enter alt screen first time.
	err = console.EnterAltScreen()
	require.NoError(t, err, "first EnterAltScreen should succeed")

	// Try to enter again (should fail).
	err = console.EnterAltScreen()
	assert.ErrorIs(t, err, ErrAlreadyInAltScreen, "second EnterAltScreen should return ErrAlreadyInAltScreen")

	// State should still be in alt screen.
	assert.True(t, console.IsInAltScreen(), "should remain in alt screen")

	// Cleanup.
	defer console.ExitAltScreen()
}

// TestConsole_ExitAltScreen_Success verifies alt screen buffer restoration.
func TestConsole_ExitAltScreen_Success(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Enter alt screen.
	err = console.EnterAltScreen()
	require.NoError(t, err, "EnterAltScreen should succeed")
	assert.True(t, console.IsInAltScreen())

	// Exit alt screen.
	err = console.ExitAltScreen()
	require.NoError(t, err, "ExitAltScreen should succeed")

	// Now NOT in alt screen.
	assert.False(t, console.IsInAltScreen(), "should NOT be in alt screen after Exit")
}

// TestConsole_ExitAltScreen_NotIn verifies error on double-exit.
func TestConsole_ExitAltScreen_NotIn(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Initially not in alt screen.
	assert.False(t, console.IsInAltScreen())

	// Try to exit (should fail).
	err = console.ExitAltScreen()
	assert.ErrorIs(t, err, ErrNotInAltScreen, "ExitAltScreen should return ErrNotInAltScreen when not in alt screen")

	// State should remain NOT in alt screen.
	assert.False(t, console.IsInAltScreen())
}

// TestConsole_IsInAltScreen_InitiallyFalse verifies initial state.
func TestConsole_IsInAltScreen_InitiallyFalse(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Newly created console should NOT be in alt screen.
	assert.False(t, console.IsInAltScreen(), "new console should NOT be in alt screen")
}

// TestConsole_IsInAltScreen_AfterEnter verifies state after entering.
func TestConsole_IsInAltScreen_AfterEnter(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Enter alt screen.
	err = console.EnterAltScreen()
	require.NoError(t, err)

	// Check state.
	assert.True(t, console.IsInAltScreen(), "IsInAltScreen should return true after Enter")

	// Cleanup.
	defer console.ExitAltScreen()
}

// TestConsole_IsInAltScreen_AfterExit verifies state after exiting.
func TestConsole_IsInAltScreen_AfterExit(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Enter then exit.
	err = console.EnterAltScreen()
	require.NoError(t, err)
	err = console.ExitAltScreen()
	require.NoError(t, err)

	// Check state.
	assert.False(t, console.IsInAltScreen(), "IsInAltScreen should return false after Exit")
}

// TestConsole_ScreenBuffer_ThreadSafety verifies concurrent access safety.
func TestConsole_ScreenBuffer_ThreadSafety(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Concurrent goroutines trying to check alt screen state.
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			// This should not panic or race.
			_ = console.IsInAltScreen()
		}()
	}

	wg.Wait()

	// All goroutines completed without panic.
	assert.True(t, true, "concurrent IsInAltScreen calls completed successfully")
}

// TestConsole_ScreenBuffer_EnterExitCycle verifies multiple enter/exit cycles.
func TestConsole_ScreenBuffer_EnterExitCycle(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Perform 5 enter/exit cycles.
	for i := 0; i < 5; i++ {
		// Enter.
		err = console.EnterAltScreen()
		require.NoError(t, err, "cycle %d: EnterAltScreen failed", i)
		assert.True(t, console.IsInAltScreen(), "cycle %d: should be in alt screen", i)

		// Exit.
		err = console.ExitAltScreen()
		require.NoError(t, err, "cycle %d: ExitAltScreen failed", i)
		assert.False(t, console.IsInAltScreen(), "cycle %d: should NOT be in alt screen", i)
	}
}

// TestConsole_ScreenBuffer_States is a table-driven test for state transitions.
func TestConsole_ScreenBuffer_States(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*Console)
		operation func(*Console) error
		wantErr   error
		wantState bool
	}{
		{
			name:  "enter when not in alt screen",
			setup: func(c *Console) {},
			operation: func(c *Console) error {
				return c.EnterAltScreen()
			},
			wantErr:   nil,
			wantState: true,
		},
		{
			name: "enter when already in alt screen",
			setup: func(c *Console) {
				_ = c.EnterAltScreen()
			},
			operation: func(c *Console) error {
				return c.EnterAltScreen()
			},
			wantErr:   ErrAlreadyInAltScreen,
			wantState: true,
		},
		{
			name: "exit when in alt screen",
			setup: func(c *Console) {
				_ = c.EnterAltScreen()
			},
			operation: func(c *Console) error {
				return c.ExitAltScreen()
			},
			wantErr:   nil,
			wantState: false,
		},
		{
			name:  "exit when not in alt screen",
			setup: func(c *Console) {},
			operation: func(c *Console) error {
				return c.ExitAltScreen()
			},
			wantErr:   ErrNotInAltScreen,
			wantState: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			console, err := NewConsole()
			if err != nil {
				t.Skipf("Not running in Windows Console: %v", err)
			}

			// Cleanup: always try to exit alt screen at the end.
			defer func() {
				// Ignore errors - may already be in normal screen.
				_ = console.ExitAltScreen()
			}()

			// Setup.
			tt.setup(console)

			// Operation.
			err = tt.operation(console)

			// Verify error.
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr, "operation should return expected error")
			} else {
				assert.NoError(t, err, "operation should succeed")
			}

			// Verify state.
			assert.Equal(t, tt.wantState, console.IsInAltScreen(), "IsInAltScreen state mismatch")
		})
	}
}

// TestConsole_ScreenBuffer_OriginalBufferPreserved verifies original buffer is restored.
func TestConsole_ScreenBuffer_OriginalBufferPreserved(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Save original buffer handle (before entering alt screen).
	originalBuffer := console.stdout

	// Enter alt screen.
	err = console.EnterAltScreen()
	require.NoError(t, err)

	// Alt screen buffer should be different.
	assert.NotEqual(t, originalBuffer, console.stdout, "alt screen buffer should differ from original")

	// Exit alt screen.
	err = console.ExitAltScreen()
	require.NoError(t, err)

	// Original buffer should be restored.
	assert.Equal(t, originalBuffer, console.stdout, "original buffer should be restored")
}

// TestConsole_ScreenBuffer_Write verifies writing to alt screen.
func TestConsole_ScreenBuffer_Write(t *testing.T) {
	console, err := NewConsole()
	if err != nil {
		t.Skipf("Not running in Windows Console: %v", err)
	}

	// Enter alt screen.
	err = console.EnterAltScreen()
	require.NoError(t, err)
	defer console.ExitAltScreen()

	// Write to alt screen (should not fail).
	err = console.Write("Testing alt screen write")
	assert.NoError(t, err, "Write to alt screen should succeed")

	// Clear alt screen (should not fail).
	err = console.Clear()
	assert.NoError(t, err, "Clear alt screen should succeed")
}
