package input

import (
	"testing"

	"github.com/phoenix-tui/phoenix/tea"
)

// TestShellBoundaryProtection_PreventCursorBeforePrompt tests the primary use case for GoSh.
func TestShellBoundaryProtection_PreventCursorBeforePrompt(t *testing.T) {
	// Simulate shell with prompt "> ".
	ta := NewTextArea().
		SetValue("> ls -la").
		SetCursorPosition(0, 8). // After "ls -la"
		OnMovement(func(_, to CursorPos) bool {
			// Prevent cursor from moving before prompt (column 2)
			if to.Row == 0 && to.Col < 2 {
				return false
			}
			return true
		})

	// Test 1: Arrow left should stop at column 2.
	for i := 0; i < 10; i++ {
		ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyLeft})
	}

	row, col := ta.CursorPosition()
	if row != 0 || col != 2 {
		t.Errorf("Expected cursor to stop at prompt boundary (0,2), got (%d,%d)", row, col)
	}

	// Test 2: Ctrl+A (Home) should move to column 2, not column 0.
	ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'a', Ctrl: true})
	row, col = ta.CursorPosition()
	if row != 0 || col != 2 {
		t.Errorf("Expected Ctrl+A to move to prompt boundary (0,2), got (%d,%d)", row, col)
	}
}

// TestShellBoundaryProtection_MultipleLines tests boundary protection with command history.
func TestShellBoundaryProtection_MultipleLines(t *testing.T) {
	// Simulate shell with history (multiple prompts)
	ta := NewTextArea().
		SetValue("> pwd\n/home/user\n> ls\nfile1.txt\n> ").
		SetCursorPosition(4, 2). // Current prompt
		OnMovement(func(_, to CursorPos) bool {
			// Only allow editing current line (row 4, col >= 2)
			if to.Row < 4 {
				return false // Block access to history
			}
			if to.Row == 4 && to.Col < 2 {
				return false // Block moving before prompt
			}
			return true
		})

	// Test 1: Arrow up should be blocked (can't edit history)
	ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyUp})
	row, col := ta.CursorPosition()
	if row != 4 || col != 2 {
		t.Errorf("Expected cursor to stay at current prompt (4,2), got (%d,%d)", row, col)
	}

	// Test 2: Arrow left should stop at column 2.
	ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyLeft})
	row, col = ta.CursorPosition()
	if row != 4 || col != 2 {
		t.Errorf("Expected cursor to stop at prompt (4,2), got (%d,%d)", row, col)
	}
}

// TestShellBoundaryProtection_WithFeedback tests user feedback when hitting boundary.
func TestShellBoundaryProtection_WithFeedback(t *testing.T) {
	boundaryHitCount := 0
	var lastReason string

	ta := NewTextArea().
		SetValue("> command").
		SetCursorPosition(0, 2).
		OnMovement(func(_, to CursorPos) bool {
			return to.Col >= 2
		}).
		OnBoundaryHit(func(_ CursorPos, reason string) {
			boundaryHitCount++
			lastReason = reason
		})

	// Try to move left (blocked)
	_, _ = ta.Update(tea.KeyMsg{Type: tea.KeyLeft})

	if boundaryHitCount != 1 {
		t.Errorf("Expected 1 boundary hit, got %d", boundaryHitCount)
	}

	if lastReason == "" {
		t.Errorf("Expected boundary hit reason, got empty string")
	}
}

// TestShellBoundaryProtection_AllowTyping tests that typing is still allowed.
func TestShellBoundaryProtection_AllowTyping(t *testing.T) {
	ta := NewTextArea().
		SetValue("> ").
		SetCursorPosition(0, 2).
		OnMovement(func(_, to CursorPos) bool {
			// Only prevent navigation before prompt.
			// Typing doesn't trigger movement validator.
			if to.Col < 2 {
				return false
			}
			return true
		})

	// Type a character (should work)
	ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'l'})
	ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 's'})

	expectedValue := "> ls"
	if ta.Value() != expectedValue {
		t.Errorf("Expected value '%s', got '%s'", expectedValue, ta.Value())
	}

	row, col := ta.CursorPosition()
	if row != 0 || col != 4 {
		t.Errorf("Expected cursor at (0,4), got (%d,%d)", row, col)
	}
}

// TestShellBoundaryProtection_DynamicPrompt tests shell with changing prompt (e.g., multi-line input).
func TestShellBoundaryProtection_DynamicPrompt(t *testing.T) {
	// Track minimum column based on prompt state.
	minCol := 2 // Initial prompt "> " is 2 chars

	ta := NewTextArea().
		SetValue("> ").
		SetCursorPosition(0, 2).
		OnMovement(func(_, to CursorPos) bool {
			if to.Row == 0 && to.Col < minCol {
				return false
			}
			return true
		})

	// Move right and type.
	ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'e'})
	ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'c'})

	// Simulate prompt change (e.g., continuation prompt "... ")
	// In real shell, you'd update minCol when prompt changes.
	minCol = 4
	ta = ta.SetValue("... echo 'hello'")
	ta = ta.SetCursorPosition(0, 4)

	// Try to move left before new prompt (should be blocked)
	for i := 0; i < 5; i++ {
		ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyLeft})
	}

	row, col := ta.CursorPosition()
	if row != 0 || col < minCol {
		t.Errorf("Expected cursor to stay at or after continuation prompt (0,%d), got (%d,%d)", minCol, row, col)
	}
}

// TestShellBoundaryProtection_RealWorldScenario tests a realistic shell interaction.
func TestShellBoundaryProtection_RealWorldScenario(t *testing.T) {
	// Track events.
	var movements []string
	var boundaries []string

	ta := NewTextArea().
		SetValue("> ").
		SetCursorPosition(0, 2).
		OnMovement(func(_, to CursorPos) bool {
			// Shell boundary: can't go before prompt.
			if to.Row == 0 && to.Col < 2 {
				return false
			}
			return true
		}).
		OnCursorMoved(func(_, _ CursorPos) {
			movements = append(movements, "moved")
		}).
		OnBoundaryHit(func(_ CursorPos, _ string) {
			boundaries = append(boundaries, "blocked")
		})

	// User types command.
	for _, ch := range "git status" {
		ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: ch})
	}

	// User presses Home (Ctrl+A) - should go to after prompt.
	ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'a', Ctrl: true})
	row, col := ta.CursorPosition()
	if row != 0 || col != 2 {
		t.Errorf("Home key should move to (0,2), got (%d,%d)", row, col)
	}

	// User tries to go left (should be blocked)
	ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if len(boundaries) == 0 {
		t.Errorf("Expected at least one boundary hit")
	}

	// User presses End (Ctrl+E) - should go to end of command.
	ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'e', Ctrl: true})
	row, col = ta.CursorPosition()
	expectedCol := 2 + len("git status")
	if row != 0 || col != expectedCol {
		t.Errorf("End key should move to (0,%d), got (%d,%d)", expectedCol, row, col)
	}

	// Verify final command.
	expectedValue := "> git status"
	if ta.Value() != expectedValue {
		t.Errorf("Expected value '%s', got '%s'", expectedValue, ta.Value())
	}
}

// TestShellBoundaryProtection_BackspaceAtPrompt tests that backspace is handled correctly.
func TestShellBoundaryProtection_BackspaceAtPrompt(t *testing.T) {
	ta := NewTextArea().
		SetValue("> test").
		SetCursorPosition(0, 2). // At prompt boundary
		OnMovement(func(_, to CursorPos) bool {
			return to.Col >= 2
		})

	// Press backspace at prompt boundary (should not delete prompt)
	ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	// Cursor should stay at column 2.
	row, col := ta.CursorPosition()
	if row != 0 || col != 2 {
		t.Errorf("Expected cursor to stay at (0,2), got (%d,%d)", row, col)
	}

	// Note: Actual backspace deletion logic is in EditingService.
	// This test just verifies cursor movement protection.
}
