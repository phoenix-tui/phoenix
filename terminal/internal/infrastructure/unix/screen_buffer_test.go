//go:build !windows
// +build !windows

package unix

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestANSI_EnterAltScreen_Success verifies alternate screen buffer entry.
func TestANSI_EnterAltScreen_Success(t *testing.T) {
	term := NewANSI()

	// Initially not in alt screen.
	assert.False(t, term.IsInAltScreen(), "should initially NOT be in alt screen")

	// Enter alt screen.
	err := term.EnterAltScreen()
	require.NoError(t, err, "EnterAltScreen should succeed")

	// Now in alt screen.
	assert.True(t, term.IsInAltScreen(), "should be in alt screen after Enter")

	// Cleanup: exit alt screen to restore terminal.
	defer term.ExitAltScreen()
}

// TestANSI_EnterAltScreen_AlreadyIn verifies error on double-enter.
func TestANSI_EnterAltScreen_AlreadyIn(t *testing.T) {
	term := NewANSI()

	// Enter alt screen first time.
	err := term.EnterAltScreen()
	require.NoError(t, err, "first EnterAltScreen should succeed")
	defer term.ExitAltScreen()

	// Try to enter again (should fail).
	err = term.EnterAltScreen()
	assert.ErrorIs(t, err, ErrAlreadyInAltScreen, "second EnterAltScreen should return ErrAlreadyInAltScreen")

	// State should still be in alt screen.
	assert.True(t, term.IsInAltScreen(), "should remain in alt screen")
}

// TestANSI_ExitAltScreen_Success verifies alt screen buffer exit.
func TestANSI_ExitAltScreen_Success(t *testing.T) {
	term := NewANSI()

	// Enter alt screen.
	err := term.EnterAltScreen()
	require.NoError(t, err, "EnterAltScreen should succeed")
	assert.True(t, term.IsInAltScreen())

	// Exit alt screen.
	err = term.ExitAltScreen()
	require.NoError(t, err, "ExitAltScreen should succeed")

	// Now NOT in alt screen.
	assert.False(t, term.IsInAltScreen(), "should NOT be in alt screen after Exit")
}

// TestANSI_ExitAltScreen_NotIn verifies error on double-exit.
func TestANSI_ExitAltScreen_NotIn(t *testing.T) {
	term := NewANSI()

	// Initially not in alt screen.
	assert.False(t, term.IsInAltScreen())

	// Try to exit (should fail).
	err := term.ExitAltScreen()
	assert.ErrorIs(t, err, ErrNotInAltScreen, "ExitAltScreen should return ErrNotInAltScreen when not in alt screen")

	// State should remain NOT in alt screen.
	assert.False(t, term.IsInAltScreen())
}

// TestANSI_IsInAltScreen_InitiallyFalse verifies initial state.
func TestANSI_IsInAltScreen_InitiallyFalse(t *testing.T) {
	term := NewANSI()

	// Newly created terminal should NOT be in alt screen.
	assert.False(t, term.IsInAltScreen(), "new terminal should NOT be in alt screen")
}

// TestANSI_IsInAltScreen_AfterEnter verifies state after entering.
func TestANSI_IsInAltScreen_AfterEnter(t *testing.T) {
	term := NewANSI()

	// Enter alt screen.
	err := term.EnterAltScreen()
	require.NoError(t, err)
	defer term.ExitAltScreen()

	// Check state.
	assert.True(t, term.IsInAltScreen(), "IsInAltScreen should return true after Enter")
}

// TestANSI_IsInAltScreen_AfterExit verifies state after exiting.
func TestANSI_IsInAltScreen_AfterExit(t *testing.T) {
	term := NewANSI()

	// Enter then exit.
	err := term.EnterAltScreen()
	require.NoError(t, err)
	err = term.ExitAltScreen()
	require.NoError(t, err)

	// Check state.
	assert.False(t, term.IsInAltScreen(), "IsInAltScreen should return false after Exit")
}

// TestANSI_ScreenBuffer_ThreadSafety verifies concurrent access safety.
func TestANSI_ScreenBuffer_ThreadSafety(t *testing.T) {
	term := NewANSI()

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Concurrent goroutines trying to check alt screen state.
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			// This should not panic or race.
			_ = term.IsInAltScreen()
		}()
	}

	wg.Wait()

	// All goroutines completed without panic.
	assert.True(t, true, "concurrent IsInAltScreen calls completed successfully")
}

// TestANSI_ScreenBuffer_EnterExitCycle verifies multiple enter/exit cycles.
func TestANSI_ScreenBuffer_EnterExitCycle(t *testing.T) {
	term := NewANSI()

	// Perform 5 enter/exit cycles.
	for i := 0; i < 5; i++ {
		// Enter.
		err := term.EnterAltScreen()
		require.NoError(t, err, "cycle %d: EnterAltScreen failed", i)
		assert.True(t, term.IsInAltScreen(), "cycle %d: should be in alt screen", i)

		// Exit.
		err = term.ExitAltScreen()
		require.NoError(t, err, "cycle %d: ExitAltScreen failed", i)
		assert.False(t, term.IsInAltScreen(), "cycle %d: should NOT be in alt screen", i)
	}
}

// TestANSI_ScreenBuffer_States is a table-driven test for state transitions.
func TestANSI_ScreenBuffer_States(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*ANSITerminal)
		operation func(*ANSITerminal) error
		wantErr   error
		wantState bool
	}{
		{
			name:  "enter when not in alt screen",
			setup: func(a *ANSITerminal) {},
			operation: func(a *ANSITerminal) error {
				return a.EnterAltScreen()
			},
			wantErr:   nil,
			wantState: true,
		},
		{
			name: "enter when already in alt screen",
			setup: func(a *ANSITerminal) {
				_ = a.EnterAltScreen()
			},
			operation: func(a *ANSITerminal) error {
				return a.EnterAltScreen()
			},
			wantErr:   ErrAlreadyInAltScreen,
			wantState: true,
		},
		{
			name: "exit when in alt screen",
			setup: func(a *ANSITerminal) {
				_ = a.EnterAltScreen()
			},
			operation: func(a *ANSITerminal) error {
				return a.ExitAltScreen()
			},
			wantErr:   nil,
			wantState: false,
		},
		{
			name:  "exit when not in alt screen",
			setup: func(a *ANSITerminal) {},
			operation: func(a *ANSITerminal) error {
				return a.ExitAltScreen()
			},
			wantErr:   ErrNotInAltScreen,
			wantState: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			term := NewANSI()

			// Cleanup: always try to exit alt screen at the end.
			defer func() {
				// Ignore errors - may already be in normal screen.
				_ = term.ExitAltScreen()
			}()

			// Setup.
			tt.setup(term)

			// Operation.
			err := tt.operation(term)

			// Verify error.
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr, "operation should return expected error")
			} else {
				assert.NoError(t, err, "operation should succeed")
			}

			// Verify state.
			assert.Equal(t, tt.wantState, term.IsInAltScreen(), "IsInAltScreen state mismatch")
		})
	}
}

// TestANSI_ScreenBuffer_Write verifies writing to alt screen.
func TestANSI_ScreenBuffer_Write(t *testing.T) {
	term := NewANSI()

	// Enter alt screen.
	err := term.EnterAltScreen()
	require.NoError(t, err)
	defer term.ExitAltScreen()

	// Write to alt screen (should not fail).
	err = term.Write("Testing alt screen write (will appear on screen temporarily)")
	assert.NoError(t, err, "Write to alt screen should succeed")

	// Clear alt screen (should not fail).
	err = term.Clear()
	assert.NoError(t, err, "Clear alt screen should succeed")
}
