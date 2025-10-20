package api

import (
	"strings"
	"testing"
)

// TestMultiline_GoSh_Scenario reproduces the exact GoSh multiline behavior.
//
// GoSh scenario:
// 1. User types "hello" → Enter (incomplete command)
// 2. Continuation prompt appears: ">>    ".
// 3. User types "world" → Enter.
// 4. Command executes: "hello\nworld".
//
// This test verifies that Phoenix TextArea correctly handles:
// - SetValue("hello\n") → MoveCursorToEnd() → cursor at (1, 0)
// - User types "world" → cursor at (1, 5)
// - View() with ShowCursor(false) shows "hello\nworld" without "█".
func TestMultiline_GoSh_Scenario(t *testing.T) {
	// Step 1: User types "hello" and presses Enter.
	ta := New().
		SetValue("hello").
		ShowCursor(false)

	// Verify initial state.
	if ta.Value() != "hello" {
		t.Errorf("Step 1: Value() = %q, want %q", ta.Value(), "hello")
	}

	row, col := ta.CursorPosition()
	if row != 0 || col != 0 {
		t.Errorf("Step 1: CursorPosition() = (%d, %d), want (0, 0)", row, col)
	}

	// Step 2: Continuation line started (shell adds newline)
	// CRITICAL: This is what GoSh does - it adds "\n" to trigger multiline.
	ta = ta.SetValue("hello\n").MoveCursorToEnd()

	// Verify cursor is at end of first line (after newline)
	row, col = ta.CursorPosition()
	if row != 1 || col != 0 {
		t.Errorf("Step 2: CursorPosition() = (%d, %d), want (1, 0)", row, col)
	}

	// Verify View() has NO cursor.
	view := ta.View()
	if strings.Contains(view, "█") {
		t.Errorf("Step 2: View() should NOT contain cursor '█'\nView:\n%s", view)
	}

	// Verify content rendering.
	expectedLines := []string{"hello", ""}
	actualLines := strings.Split(view, "\n")
	if len(actualLines) != len(expectedLines) {
		t.Errorf("Step 2: View() lines = %d, want %d\nView:\n%s", len(actualLines), len(expectedLines), view)
	}

	// Step 3: User types "world" on continuation line.
	// In real usage, this would be done character-by-character via Update(KeyMsg)
	// But we can simulate by directly setting value.
	ta = ta.SetValue("hello\nworld").MoveCursorToEnd()

	// Verify cursor is at end of second line.
	row, col = ta.CursorPosition()
	if row != 1 || col != 5 {
		t.Errorf("Step 3: CursorPosition() = (%d, %d), want (1, 5)", row, col)
	}

	// Verify View() has NO cursor.
	view = ta.View()
	if strings.Contains(view, "█") {
		t.Errorf("Step 3: View() should NOT contain cursor '█'\nView:\n%s", view)
	}

	// Verify final content.
	if ta.Value() != "hello\nworld" {
		t.Errorf("Step 3: Value() = %q, want %q", ta.Value(), "hello\nworld")
	}

	// Print final View for debugging (when test passes)
	t.Logf("Final View (ShowCursor=false):\n%s\n---\nCursor at: (%d, %d)", view, row, col)
}

// TestMultiline_GoSh_ViewWithPrompts simulates ViewWithPrompts() rendering.
//
// This is what GoSh actually does:
// 1. Phoenix TextArea.View() → "hello\nworld".
// 2. Split by "\n" → ["hello", "world"].
// 3. Add prompts → "gosh> hello\n>>    world".
//
// PROBLEM INVESTIGATION:
// When user types each character in "world", does Phoenix correctly.
// render the cursor position WITHOUT inserting "█"?
func TestMultiline_GoSh_ViewWithPrompts(t *testing.T) {
	// Simulate character-by-character typing on continuation line.
	tests := []struct {
		name      string
		text      string
		wantLines []string
	}{
		{
			name:      "empty continuation line",
			text:      "hello\n",
			wantLines: []string{"hello", ""},
		},
		{
			name:      "one character on continuation line",
			text:      "hello\nw",
			wantLines: []string{"hello", "w"},
		},
		{
			name:      "two characters on continuation line",
			text:      "hello\nwo",
			wantLines: []string{"hello", "wo"},
		},
		{
			name:      "full word on continuation line",
			text:      "hello\nworld",
			wantLines: []string{"hello", "world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := New().
				SetValue(tt.text).
				MoveCursorToEnd().
				ShowCursor(false)

			view := ta.View()

			// Verify NO cursor.
			if strings.Contains(view, "█") {
				t.Errorf("View() should NOT contain cursor '█'\nView:\n%s", view)
			}

			// Verify line count.
			lines := strings.Split(view, "\n")
			if len(lines) != len(tt.wantLines) {
				t.Errorf("View() lines = %d, want %d\nView:\n%s", len(lines), len(tt.wantLines), view)
			}

			// Verify each line content.
			for i, wantLine := range tt.wantLines {
				if i >= len(lines) {
					t.Errorf("Missing line %d: %q", i, wantLine)
					continue
				}
				if lines[i] != wantLine {
					t.Errorf("Line %d = %q, want %q", i, lines[i], wantLine)
				}
			}

			// Simulate adding prompts (what GoSh does)
			primaryPrompt := "gosh> "
			continuationPrompt := ">>    "

			var result strings.Builder
			for i, line := range lines {
				if i == 0 {
					result.WriteString(primaryPrompt)
				} else {
					result.WriteString(continuationPrompt)
				}
				result.WriteString(line)
				if i < len(lines)-1 {
					result.WriteString("\n")
				}
			}

			withPrompts := result.String()
			t.Logf("With prompts:\n%s", withPrompts)

			// Verify prompts are added correctly.
			if !strings.Contains(withPrompts, primaryPrompt) {
				t.Errorf("Missing primary prompt in: %s", withPrompts)
			}
			if len(lines) > 1 && !strings.Contains(withPrompts, continuationPrompt) {
				t.Errorf("Missing continuation prompt in: %s", withPrompts)
			}
		})
	}
}

// TestMultiline_Cursor_Position_Bug_Investigation investigates the "jumping" behavior.
//
// HYPOTHESIS:
// When Phoenix TextArea renders with ShowCursor(false), the View() output.
// might not match the expected line structure, causing GoSh's ViewWithPrompts()
// to add prompts incorrectly.
//
// This test checks if View() output structure is stable across character additions.
func TestMultiline_Cursor_Position_Bug_Investigation(t *testing.T) {
	ta := New().
		SetValue("hello\n").
		MoveCursorToEnd().
		ShowCursor(false)

	// Record View() output before each character.
	//nolint:prealloc // slice size is small and dynamic growth is acceptable in tests
	var views []string
	texts := []string{"hello\n", "hello\nw", "hello\nwo", "hello\nwor", "hello\nworl", "hello\nworld"}

	for _, text := range texts {
		ta = ta.SetValue(text).MoveCursorToEnd().ShowCursor(false)
		view := ta.View()
		views = append(views, view)

		t.Logf("Text: %q\nView:\n%s\n---", text, view)

		// Verify View() structure consistency.
		lines := strings.Split(view, "\n")
		if len(lines) != 2 {
			t.Errorf("Text %q: View() lines = %d, want 2\nView:\n%s", text, len(lines), view)
		}

		// Verify NO cursor.
		if strings.Contains(view, "█") {
			t.Errorf("Text %q: View() contains cursor '█'\nView:\n%s", text, view)
		}
	}

	// Check if View() output is consistent (only second line changes)
	for i := 1; i < len(views); i++ {
		prev := strings.Split(views[i-1], "\n")
		curr := strings.Split(views[i], "\n")

		if len(prev) != 2 || len(curr) != 2 {
			t.Errorf("View %d structure inconsistent: prev=%d lines, curr=%d lines", i, len(prev), len(curr))
			continue
		}

		// First line should NEVER change.
		if prev[0] != curr[0] {
			t.Errorf("View %d first line CHANGED: %q → %q", i, prev[0], curr[0])
		}

		t.Logf("View %d consistency: first line=%q (unchanged), second line=%q → %q",
			i, curr[0], prev[1], curr[1])
	}
}
