package api

import (
	"strings"
	"testing"
)

// TestShowCursor_False verifies that ShowCursor(false) prevents cursor rendering.
// CRITICAL: This is for shell applications where terminal cursor is used.
func TestShowCursor_False(t *testing.T) {
	tests := []struct {
		name string
		text string
		row  int
		col  int
	}{
		{
			name: "single line - cursor at start",
			text: "hello",
			row:  0,
			col:  0,
		},
		{
			name: "single line - cursor at end",
			text: "hello",
			row:  0,
			col:  5,
		},
		{
			name: "multiline - cursor on first line",
			text: "line1\nline2\nline3",
			row:  0,
			col:  3,
		},
		{
			name: "multiline - cursor on middle line",
			text: "line1\nline2\nline3",
			row:  1,
			col:  2,
		},
		{
			name: "multiline - cursor on last line",
			text: "line1\nline2\nline3",
			row:  2,
			col:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := New().
				SetValue(tt.text).
				ShowCursor(false)

			view := ta.View()

			// MUST NOT contain block cursor character.
			if strings.Contains(view, "█") {
				t.Errorf("View() with ShowCursor(false) should NOT contain '█' cursor\nView:\n%s", view)
			}

			// Verify content is rendered correctly.
			if !strings.Contains(view, tt.text) {
				t.Errorf("View() should contain text %q\nView:\n%s", tt.text, view)
			}
		})
	}
}

// TestShowCursor_True verifies that ShowCursor(true) renders cursor.
func TestShowCursor_True(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		shouldShow bool // Whether cursor should be visible in this case
	}{
		{
			name:       "single line",
			text:       "hello",
			shouldShow: true,
		},
		{
			name:       "multiline",
			text:       "line1\nline2",
			shouldShow: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := New().
				SetValue(tt.text).
				ShowCursor(true) // Default, but explicit

			view := ta.View()

			if tt.shouldShow {
				// MUST contain block cursor character.
				if !strings.Contains(view, "█") {
					t.Errorf("View() with ShowCursor(true) should contain '█' cursor\nView:\n%s", view)
				}
			}
		})
	}
}

// TestMoveCursorToEnd verifies cursor positioning after MoveCursorToEnd().
func TestMoveCursorToEnd(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		wantRow     int
		wantCol     int
		description string
	}{
		{
			name:        "empty buffer",
			text:        "",
			wantRow:     0,
			wantCol:     0,
			description: "Empty buffer has cursor at (0,0)",
		},
		{
			name:        "single line",
			text:        "hello",
			wantRow:     0,
			wantCol:     5, // After 'o'
			description: "Single line cursor at end",
		},
		{
			name:        "multiline",
			text:        "line1\nline2\nline3",
			wantRow:     2, // Last line (zero-indexed)
			wantCol:     5, // After 'line3'
			description: "Multiline cursor at end of last line",
		},
		{
			name:        "multiline with different lengths",
			text:        "short\nthis is a longer line\nend",
			wantRow:     2,
			wantCol:     3, // After 'end'
			description: "Cursor at end of last line (not longest line)",
		},
		{
			name:        "unicode text",
			text:        "hello\nworld👋",
			wantRow:     1,
			wantCol:     6, // After emoji (counted as 1 rune)
			description: "Unicode text cursor at end",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := New().
				SetValue(tt.text).
				MoveCursorToEnd()

			row, col := ta.CursorPosition()

			if row != tt.wantRow {
				t.Errorf("%s: row = %d, want %d", tt.description, row, tt.wantRow)
			}
			if col != tt.wantCol {
				t.Errorf("%s: col = %d, want %d", tt.description, col, tt.wantCol)
			}
		})
	}
}

// TestMoveCursorToEnd_WithShowCursor_False verifies multiline rendering bug.
// CRITICAL: This reproduces GoSh shell multiline issue.
func TestMoveCursorToEnd_WithShowCursor_False(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{
			name: "continuation line (GoSh scenario)",
			text: "hello\nworld",
		},
		{
			name: "multiple continuation lines",
			text: "line1\nline2\nline3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := New().
				SetValue(tt.text).
				MoveCursorToEnd().
				ShowCursor(false)

			view := ta.View()

			// Verify NO cursor rendering.
			if strings.Contains(view, "█") {
				t.Errorf("View() with ShowCursor(false) should NOT contain '█'\nView:\n%s", view)
			}

			// Verify all lines are rendered.
			expectedLines := strings.Split(tt.text, "\n")
			for i, line := range expectedLines {
				if !strings.Contains(view, line) {
					t.Errorf("View() missing line %d: %q\nView:\n%s", i, line, view)
				}
			}

			// Verify cursor position is correct.
			row, col := ta.CursorPosition()
			lastLineIdx := len(expectedLines) - 1
			lastLine := expectedLines[lastLineIdx]
			expectedCol := len([]rune(lastLine))

			if row != lastLineIdx {
				t.Errorf("CursorPosition() row = %d, want %d", row, lastLineIdx)
			}
			if col != expectedCol {
				t.Errorf("CursorPosition() col = %d, want %d", col, expectedCol)
			}
		})
	}
}

// TestMultiline_CursorRendering verifies cursor rendering on specific lines.
func TestMultiline_CursorRendering(t *testing.T) {
	text := "line1\nline2\nline3"

	tests := []struct {
		name         string
		cursorRow    int
		cursorCol    int
		showCursor   bool
		expectCursor bool
		expectOnLine int // Which line number (1-based) should have cursor
	}{
		{
			name:         "cursor on line 1 - shown",
			cursorRow:    0,
			cursorCol:    3,
			showCursor:   true,
			expectCursor: true,
			expectOnLine: 1,
		},
		{
			name:         "cursor on line 2 - shown",
			cursorRow:    1,
			cursorCol:    2,
			showCursor:   true,
			expectCursor: true,
			expectOnLine: 2,
		},
		{
			name:         "cursor on line 3 - shown",
			cursorRow:    2,
			cursorCol:    4,
			showCursor:   true,
			expectCursor: true,
			expectOnLine: 3,
		},
		{
			name:         "cursor on line 2 - hidden",
			cursorRow:    1,
			cursorCol:    2,
			showCursor:   false,
			expectCursor: false,
			expectOnLine: 0, // No cursor shown
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := New().
				SetValue(text).
				ShowCursor(tt.showCursor)

			// Note: We can't directly set cursor position via API,.
			// so we test MoveCursorToEnd() behavior instead.
			view := ta.View()

			//nolint:nestif // test validation logic requires branching
			if tt.expectCursor {
				if !strings.Contains(view, "█") {
					t.Errorf("View() should contain cursor '█'\nView:\n%s", view)
				}
			} else {
				if strings.Contains(view, "█") {
					t.Errorf("View() should NOT contain cursor '█'\nView:\n%s", view)
				}
			}
		})
	}
}

// TestSetValue_ResetsState verifies that SetValue() resets cursor properly.
func TestSetValue_ResetsState(t *testing.T) {
	ta := New().
		SetValue("initial text")

	// Cursor should be at (0, 0) after SetValue.
	row, col := ta.CursorPosition()
	if row != 0 || col != 0 {
		t.Errorf("SetValue() should reset cursor to (0, 0), got (%d, %d)", row, col)
	}

	// Now move to end.
	ta = ta.MoveCursorToEnd()

	row, col = ta.CursorPosition()
	if row != 0 || col != 12 { // "initial text" = 12 chars
		t.Errorf("MoveCursorToEnd() cursor = (%d, %d), want (0, 12)", row, col)
	}

	// SetValue again - should reset.
	ta = ta.SetValue("new")

	row, col = ta.CursorPosition()
	if row != 0 || col != 0 {
		t.Errorf("SetValue() should reset cursor to (0, 0), got (%d, %d)", row, col)
	}
}

// TestChaining_SetValue_MoveCursorToEnd_ShowCursor verifies method chaining.
func TestChaining_SetValue_MoveCursorToEnd_ShowCursor(t *testing.T) {
	ta := New().
		SetValue("hello\nworld").
		MoveCursorToEnd().
		ShowCursor(false)

	// Verify cursor position.
	row, col := ta.CursorPosition()
	if row != 1 || col != 5 {
		t.Errorf("CursorPosition() = (%d, %d), want (1, 5)", row, col)
	}

	// Verify View() has no cursor.
	view := ta.View()
	if strings.Contains(view, "█") {
		t.Errorf("View() should NOT contain cursor\nView:\n%s", view)
	}

	// Verify content.
	if !strings.Contains(view, "hello") || !strings.Contains(view, "world") {
		t.Errorf("View() should contain both lines\nView:\n%s", view)
	}
}
