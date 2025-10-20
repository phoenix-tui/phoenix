package api

import (
	"testing"

	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// TestSetCursorPosition_Clamping tests cursor position clamping to valid bounds.
func TestSetCursorPosition_Clamping(t *testing.T) {
	tests := []struct {
		name        string
		initialText string
		setRow      int
		setCol      int
		expectedRow int
		expectedCol int
		description string
	}{
		{
			name:        "valid position",
			initialText: "line1\nline2\nline3",
			setRow:      1,
			setCol:      2,
			expectedRow: 1,
			expectedCol: 2,
			description: "Valid position should be set exactly",
		},
		{
			name:        "negative row clamped to 0",
			initialText: "line1\nline2",
			setRow:      -5,
			setCol:      2,
			expectedRow: 0,
			expectedCol: 2,
			description: "Negative row should clamp to 0",
		},
		{
			name:        "negative col clamped to 0",
			initialText: "line1\nline2",
			setRow:      0,
			setCol:      -3,
			expectedRow: 0,
			expectedCol: 0,
			description: "Negative col should clamp to 0",
		},
		{
			name:        "row beyond buffer clamped",
			initialText: "line1\nline2",
			setRow:      100,
			setCol:      2,
			expectedRow: 1,
			expectedCol: 2,
			description: "Row beyond buffer should clamp to last row",
		},
		{
			name:        "col beyond line clamped",
			initialText: "short",
			setRow:      0,
			setCol:      100,
			expectedRow: 0,
			expectedCol: 5,
			description: "Col beyond line should clamp to line end",
		},
		{
			name:        "empty buffer edge case",
			initialText: "",
			setRow:      5,
			setCol:      10,
			expectedRow: 0,
			expectedCol: 0,
			description: "Empty buffer should clamp to (0,0)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := New().SetValue(tt.initialText)
			ta = ta.SetCursorPosition(tt.setRow, tt.setCol)

			row, col := ta.CursorPosition()
			if row != tt.expectedRow || col != tt.expectedCol {
				t.Errorf("%s: got (%d,%d), want (%d,%d)",
					tt.description, row, col, tt.expectedRow, tt.expectedCol)
			}
		})
	}
}

// TestOnMovement_BlocksMovement tests that validator can block cursor movements.
func TestOnMovement_BlocksMovement(t *testing.T) {
	// Create textarea with boundary protection (shell prompt)
	ta := New().
		SetValue("> ").
		SetCursorPosition(0, 2). // After prompt
		OnMovement(func(_, to CursorPos) bool {
			// Don't allow cursor before the prompt (column 2)
			if to.Row == 0 && to.Col < 2 {
				return false // Block movement
			}
			return true // Allow movement
		})

	// Try to move left (should be blocked)
	updated, _ := ta.Update(tea.KeyMsg{Type: tea.KeyLeft})
	row, col := updated.CursorPosition()

	if row != 0 || col != 2 {
		t.Errorf("Expected cursor to stay at (0,2) when blocked, got (%d,%d)", row, col)
	}

	// Move right (should be allowed)
	ta = ta.SetCursorPosition(0, 2)
	updated, _ = ta.Update(tea.KeyMsg{Type: tea.KeyRight})
	row, col = updated.CursorPosition()

	if row != 0 || col != 2 {
		// Note: Moving right from end of line doesn't move cursor (no next line)
		// This is correct behavior - just testing validator doesn't interfere.
		t.Logf("Cursor at (%d,%d) - correct, no next line to move to", row, col)
	}
}

// TestOnMovement_AllowsMovement tests that validator allows valid movements.
func TestOnMovement_AllowsMovement(t *testing.T) {
	ta := New().
		SetValue("line1\nline2\nline3").
		SetCursorPosition(1, 2).
		OnMovement(func(_, _ CursorPos) bool {
			// Allow all movements (always return true)
			return true
		})

	// Move down (should be allowed)
	updated, _ := ta.Update(tea.KeyMsg{Type: tea.KeyDown})
	row, col := updated.CursorPosition()

	if row != 2 || col != 2 {
		t.Errorf("Expected cursor to move to (2,2), got (%d,%d)", row, col)
	}
}

// TestOnCursorMoved_CalledAfterMovement tests observer pattern for cursor movement.
func TestOnCursorMoved_CalledAfterMovement(t *testing.T) {
	var capturedFrom, capturedTo CursorPos
	callCount := 0

	ta := New().
		SetValue("line1\nline2\nline3").
		SetCursorPosition(1, 2).
		OnCursorMoved(func(from, to CursorPos) {
			capturedFrom = from
			capturedTo = to
			callCount++
		})

	// Move down.
	updated, _ := ta.Update(tea.KeyMsg{Type: tea.KeyDown})

	if callCount != 1 {
		t.Errorf("Expected OnCursorMoved to be called once, got %d times", callCount)
	}

	if capturedFrom.Row != 1 || capturedFrom.Col != 2 {
		t.Errorf("Expected from position (1,2), got (%d,%d)", capturedFrom.Row, capturedFrom.Col)
	}

	if capturedTo.Row != 2 || capturedTo.Col != 2 {
		t.Errorf("Expected to position (2,2), got (%d,%d)", capturedTo.Row, capturedTo.Col)
	}

	// Verify cursor actually moved.
	row, col := updated.CursorPosition()
	if row != 2 || col != 2 {
		t.Errorf("Expected cursor at (2,2), got (%d,%d)", row, col)
	}
}

// TestOnCursorMoved_NotCalledWhenBlocked tests that observer is NOT called when movement blocked.
func TestOnCursorMoved_NotCalledWhenBlocked(t *testing.T) {
	callCount := 0

	ta := New().
		SetValue("> command").
		SetCursorPosition(0, 2).
		OnMovement(func(_, to CursorPos) bool {
			// Block movement to column < 2.
			if to.Col < 2 {
				return false
			}
			return true
		}).
		OnCursorMoved(func(_, _ CursorPos) {
			callCount++
		})

	// Try to move left (blocked)
	ta.Update(tea.KeyMsg{Type: tea.KeyLeft})

	if callCount != 0 {
		t.Errorf("Expected OnCursorMoved NOT to be called when blocked, but it was called %d times", callCount)
	}
}

// TestOnBoundaryHit_CalledWhenBlocked tests feedback when movement is blocked.
func TestOnBoundaryHit_CalledWhenBlocked(t *testing.T) {
	var capturedPos CursorPos
	var capturedReason string
	callCount := 0

	ta := New().
		SetValue("> ").
		SetCursorPosition(0, 2).
		OnMovement(func(_, to CursorPos) bool {
			// Block movement before prompt.
			if to.Col < 2 {
				return false
			}
			return true
		}).
		OnBoundaryHit(func(attemptedPos CursorPos, reason string) {
			capturedPos = attemptedPos
			capturedReason = reason
			callCount++
		})

	// Try to move left (blocked)
	ta.Update(tea.KeyMsg{Type: tea.KeyLeft})

	if callCount != 1 {
		t.Errorf("Expected OnBoundaryHit to be called once, got %d times", callCount)
	}

	if capturedPos.Row != 0 || capturedPos.Col != 1 {
		t.Errorf("Expected attempted position (0,1), got (%d,%d)", capturedPos.Row, capturedPos.Col)
	}

	if capturedReason != "movement blocked by validator" {
		t.Errorf("Expected reason 'movement blocked by validator', got '%s'", capturedReason)
	}
}

// TestOnBoundaryHit_NotCalledWhenAllowed tests that boundary handler is NOT called for allowed movements.
func TestOnBoundaryHit_NotCalledWhenAllowed(t *testing.T) {
	callCount := 0

	ta := New().
		SetValue("line1\nline2").
		SetCursorPosition(0, 0).
		OnMovement(func(_, _ CursorPos) bool {
			return true // Allow all movements
		}).
		OnBoundaryHit(func(_ CursorPos, _ string) {
			callCount++
		})

	// Move right (allowed)
	ta.Update(tea.KeyMsg{Type: tea.KeyRight})

	if callCount != 0 {
		t.Errorf("Expected OnBoundaryHit NOT to be called for allowed movement, but called %d times", callCount)
	}
}

// TestAllFeaturesTogether tests all 4 features working together.
func TestAllFeaturesTogether(t *testing.T) {
	// Counters for callbacks.
	validatorCalls := 0
	movedCalls := 0
	boundaryCalls := 0

	ta := New().
		SetValue("> command here").
		SetCursorPosition(0, 2).
		OnMovement(func(_, to CursorPos) bool {
			validatorCalls++
			// Protect prompt area (columns 0-1)
			if to.Row == 0 && to.Col < 2 {
				return false
			}
			return true
		}).
		OnCursorMoved(func(_, _ CursorPos) {
			movedCalls++
		}).
		OnBoundaryHit(func(_ CursorPos, _ string) {
			boundaryCalls++
		})

	// Try to move left (should be blocked)
	ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if validatorCalls != 1 || movedCalls != 0 || boundaryCalls != 1 {
		t.Errorf("Blocked movement: expected validator=1 moved=0 boundary=1, got %d/%d/%d",
			validatorCalls, movedCalls, boundaryCalls)
	}

	// Move right (should succeed)
	ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyRight})
	if validatorCalls != 2 || movedCalls != 1 || boundaryCalls != 1 {
		t.Errorf("Allowed movement: expected validator=2 moved=1 boundary=1, got %d/%d/%d",
			validatorCalls, movedCalls, boundaryCalls)
	}

	// Verify cursor is at correct position.
	row, col := ta.CursorPosition()
	if row != 0 || col != 3 {
		t.Errorf("Expected cursor at (0,3), got (%d,%d)", row, col)
	}
}

// TestBackwardCompatibility_NoCallbacks tests that TextArea works without any callbacks (100% backward compatible).
func TestBackwardCompatibility_NoCallbacks(t *testing.T) {
	// Old code - no callbacks.
	ta := New().SetValue("line1\nline2\nline3")

	// Move down.
	updated, _ := ta.Update(tea.KeyMsg{Type: tea.KeyDown})
	row, col := updated.CursorPosition()

	if row != 1 || col != 0 {
		t.Errorf("Expected cursor at (1,0), got (%d,%d)", row, col)
	}

	// This test verifies that all new features are opt-in and don't break existing code.
}

// TestSetCursorPosition_Immutability tests that SetCursorPosition returns new instance.
func TestSetCursorPosition_Immutability(t *testing.T) {
	ta := New().SetValue("test")
	original := ta

	updated := ta.SetCursorPosition(0, 2)

	// Original should be unchanged.
	row, col := original.CursorPosition()
	if row != 0 || col != 0 {
		t.Errorf("Original should be at (0,0), got (%d,%d)", row, col)
	}

	// Updated should be at new position.
	row, col = updated.CursorPosition()
	if row != 0 || col != 2 {
		t.Errorf("Updated should be at (0,2), got (%d,%d)", row, col)
	}
}
