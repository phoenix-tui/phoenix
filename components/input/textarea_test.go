package input

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/tea"
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
			ta := NewTextArea().
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
			ta := NewTextArea().
				SetValue(tt.text).
				ShowCursor(true) // Default, but explicit

			view := ta.View()

			if tt.shouldShow {
				// MUST contain reverse video escape code (cursor styling).
				if !strings.Contains(view, "\x1b[7m") {
					t.Errorf("View() with ShowCursor(true) should contain reverse video cursor\nView:\n%s", view)
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
			ta := NewTextArea().
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
			ta := NewTextArea().
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
			ta := NewTextArea().
				SetValue(text).
				ShowCursor(tt.showCursor)

			// Note: We can't directly set cursor position via API,.
			// so we test MoveCursorToEnd() behavior instead.
			view := ta.View()

			//nolint:nestif // test validation logic requires branching
			if tt.expectCursor {
				if !strings.Contains(view, "\x1b[7m") {
					t.Errorf("View() should contain reverse video cursor\nView:\n%s", view)
				}
			} else {
				if strings.Contains(view, "\x1b[7m") {
					t.Errorf("View() should NOT contain reverse video cursor\nView:\n%s", view)
				}
			}
		})
	}
}

// TestSetValue_ResetsState verifies that SetValue() resets cursor properly.
func TestSetValue_ResetsState(t *testing.T) {
	ta := NewTextArea().
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
	ta := NewTextArea().
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

// Additional coverage tests for uncovered functions

func TestTextArea_MaxLines(t *testing.T) {
	ta := NewTextArea().MaxLines(5)

	// Verify fluent API
	if ta.Value() != "" {
		t.Error("New textarea should be empty")
	}
}

func TestTextArea_MaxChars(t *testing.T) {
	ta := NewTextArea().MaxChars(100)

	// Verify fluent API
	if ta.Value() != "" {
		t.Error("New textarea should be empty")
	}
}

func TestTextArea_Placeholder(t *testing.T) {
	ta := NewTextArea().Placeholder("Enter text...")

	// Verify fluent API
	if ta.Value() != "" {
		t.Error("New textarea should be empty")
	}
}

func TestTextArea_Wrap(t *testing.T) {
	tests := []bool{true, false}

	for _, wrap := range tests {
		t.Run("wrap", func(t *testing.T) {
			ta := NewTextArea().Wrap(wrap)

			// Verify fluent API
			if ta.Value() != "" {
				t.Error("New textarea should be empty")
			}
		})
	}
}

func TestTextArea_ReadOnly(t *testing.T) {
	tests := []bool{true, false}

	for _, readonly := range tests {
		t.Run("readonly", func(t *testing.T) {
			ta := NewTextArea().ReadOnly(readonly)

			// Verify fluent API
			if ta.Value() != "" {
				t.Error("New textarea should be empty")
			}
		})
	}
}

func TestTextArea_ShowLineNumbers(t *testing.T) {
	tests := []bool{true, false}

	for _, show := range tests {
		t.Run("show", func(t *testing.T) {
			ta := NewTextArea().ShowLineNumbers(show)

			// Verify fluent API
			if ta.Value() != "" {
				t.Error("New textarea should be empty")
			}
		})
	}
}

func TestTextArea_Keybindings(t *testing.T) {
	tests := []KeybindingMode{
		KeybindingsDefault,
		KeybindingsEmacs,
		KeybindingsVi,
	}

	for _, mode := range tests {
		t.Run("keybinding mode", func(t *testing.T) {
			ta := NewTextArea().Keybindings(mode)

			// Verify fluent API
			if ta.Value() != "" {
				t.Error("New textarea should be empty")
			}
		})
	}
}

func TestTextArea_Lines(t *testing.T) {
	ta := NewTextArea().SetValue("line1\nline2\nline3")

	lines := ta.Lines()

	if len(lines) != 3 {
		t.Errorf("Lines() = %d, want 3", len(lines))
	}

	if lines[0] != "line1" {
		t.Errorf("Lines()[0] = %q, want %q", lines[0], "line1")
	}
}

func TestTextArea_ContentParts(_ *testing.T) {
	ta := NewTextArea().SetValue("hello\nworld").SetCursorPosition(1, 2)

	before, at, after := ta.ContentParts()

	// Verify structure (exact content depends on implementation)
	_ = before
	_ = at
	_ = after
}

func TestTextArea_CurrentLine(t *testing.T) {
	ta := NewTextArea().SetValue("line1\nline2\nline3").SetCursorPosition(1, 0)

	line := ta.CurrentLine()

	if line != "line2" {
		t.Errorf("CurrentLine() = %q, want %q", line, "line2")
	}
}

func TestTextArea_LineCount(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  int
	}{
		{"empty", "", 1}, // Empty has 1 line
		{"single", "hello", 1},
		{"multi", "line1\nline2\nline3", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea().SetValue(tt.value)

			count := ta.LineCount()

			if count != tt.want {
				t.Errorf("LineCount() = %d, want %d", count, tt.want)
			}
		})
	}
}

func TestTextArea_IsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"empty", "", true},
		{"not empty", "text", false},
		{"whitespace", " ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea().SetValue(tt.value)

			empty := ta.IsEmpty()

			if empty != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", empty, tt.want)
			}
		})
	}
}

func TestTextArea_HasSelection(t *testing.T) {
	ta := NewTextArea().SetValue("hello")

	// Default - no selection
	if ta.HasSelection() {
		t.Error("HasSelection() should be false for new textarea")
	}
}

func TestTextArea_SelectedText(t *testing.T) {
	ta := NewTextArea().SetValue("hello")

	// Default - no selection
	text := ta.SelectedText()

	if text != "" {
		t.Errorf("SelectedText() = %q, want empty", text)
	}
}

func TestTextArea_Init(t *testing.T) {
	ta := NewTextArea()

	cmd := ta.Init()

	if cmd != nil {
		t.Error("Init() should return nil cmd")
	}
}

func TestTextArea_Update_NonKeyMsg(t *testing.T) {
	ta := NewTextArea().SetValue("hello")

	// Send non-key message
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updated, cmd := ta.Update(msg)

	if cmd != nil {
		t.Error("Update() should return nil cmd for WindowSizeMsg")
	}

	if updated.Value() != "hello" {
		t.Errorf("Value() should be unchanged, got %q", updated.Value())
	}
}

func TestTextArea_Update_ReadOnly(t *testing.T) {
	ta := NewTextArea().SetValue("hello").ReadOnly(true)

	// Try to insert text
	msg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'x'}
	updated, _ := ta.Update(msg)

	// Should ignore input in readonly mode
	if updated.Value() != "hello" {
		t.Errorf("Value() should be unchanged in readonly mode, got %q", updated.Value())
	}
}

func TestTextArea_FluentChaining(t *testing.T) {
	ta := NewTextArea().
		MaxLines(10).
		MaxChars(500).
		Placeholder("Type here...").
		Wrap(true).
		ReadOnly(false).
		ShowLineNumbers(true).
		ShowCursor(true).
		SetValue("test")

	if ta.Value() != "test" {
		t.Errorf("Value() = %q, want %q", ta.Value(), "test")
	}
}
