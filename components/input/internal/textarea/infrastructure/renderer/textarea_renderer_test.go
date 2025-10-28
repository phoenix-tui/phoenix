package renderer

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/components/input/internal/textarea/domain/model"
)

func TestTextAreaRenderer_Render_Empty(t *testing.T) {
	r := NewTextAreaRenderer()
	ta := model.NewTextArea()

	result := r.Render(ta)

	// Empty textarea should render as empty string (no placeholder set)
	if result != "" {
		t.Errorf("Render() empty textarea = %q, want empty string", result)
	}
}

func TestTextAreaRenderer_Render_Placeholder(t *testing.T) {
	r := NewTextAreaRenderer()
	ta := model.NewTextArea().WithPlaceholder("Enter text...")

	result := r.Render(ta)

	// Should show placeholder with gray styling (contains ANSI escape codes).
	// We check that the placeholder text is present and has ANSI codes for gray color.
	if !strings.Contains(result, "Enter text...") {
		t.Errorf("Render() placeholder should contain 'Enter text...', got: %q", result)
	}

	// Check for ANSI color code.
	// Phoenix style system uses TrueColor by default, so color 240 (gray) converts to RGB(88,88,88)
	// This produces: "\x1b[38;2;88;88;88m" (TrueColor) instead of "\x1b[38;5;240m" (256-color)
	hasTrueColor := strings.Contains(result, "\x1b[38;2;88;88;88m")
	has256Color := strings.Contains(result, "\x1b[38;5;240m")
	if !hasTrueColor && !has256Color {
		t.Errorf("Render() placeholder should contain gray ANSI code (TrueColor or 256-color), got: %q", result)
	}

	// Check for reset code
	if !strings.Contains(result, "\x1b[0m") {
		t.Errorf("Render() placeholder should contain reset code, got: %q", result)
	}
}

func TestTextAreaRenderer_Render_SingleLine_NoCursor(t *testing.T) {
	r := NewTextAreaRenderer()
	ta := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("hello")).
		WithShowCursor(false) // Disable cursor

	result := r.Render(ta)

	// Should render text WITHOUT cursor.
	if result != "hello" {
		t.Errorf("Render() = %q, want %q", result, "hello")
	}

	// MUST NOT contain cursor (neither block cursor nor reverse video).
	if strings.Contains(result, "█") {
		t.Errorf("Render() with ShowCursor(false) contains block cursor '█'")
	}
	if strings.Contains(result, "\x1b[7m") {
		t.Errorf("Render() with ShowCursor(false) contains reverse video escape code")
	}
}

func TestTextAreaRenderer_Render_SingleLine_WithCursor(t *testing.T) {
	r := NewTextAreaRenderer()
	ta := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("hello")).
		WithShowCursor(true) // Enable cursor (default)

	result := r.Render(ta)

	// Should render text WITH cursor at position 0 using reverse video.
	// Expected: "\x1b[7mh\x1b[27mello" (cursor on 'h' with reverse video)
	if !strings.Contains(result, "\x1b[7m") {
		t.Errorf("Render() with ShowCursor(true) should contain reverse video escape code, got: %q", result)
	}
	if !strings.Contains(result, "\x1b[27m") {
		t.Errorf("Render() with ShowCursor(true) should contain reverse video off escape code, got: %q", result)
	}
	// Text should still contain "hello"
	if !strings.Contains(result, "hello") && !strings.Contains(result, "ello") {
		t.Errorf("Render() should contain text, got: %q", result)
	}
}

func TestTextAreaRenderer_Render_Multiline_NoCursor(t *testing.T) {
	r := NewTextAreaRenderer()
	ta := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("line1\nline2\nline3")).
		WithShowCursor(false)

	result := r.Render(ta)

	// Should render all lines WITHOUT cursor.
	expected := "line1\nline2\nline3"
	if result != expected {
		t.Errorf("Render() = %q, want %q", result, expected)
	}

	// MUST NOT contain cursor (neither block cursor nor reverse video).
	if strings.Contains(result, "█") {
		t.Errorf("Render() with ShowCursor(false) contains block cursor '█'")
	}
	if strings.Contains(result, "\x1b[7m") {
		t.Errorf("Render() with ShowCursor(false) contains reverse video escape code")
	}
}

func TestTextAreaRenderer_Render_Multiline_WithCursor_FirstLine(t *testing.T) {
	r := NewTextAreaRenderer()
	ta := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("line1\nline2\nline3")).
		WithCursor(model.NewCursor(0, 3)). // Cursor on first line at col 3
		WithShowCursor(true)

	result := r.Render(ta)

	// Should render with cursor on FIRST line using reverse video.
	lines := strings.Split(result, "\n")
	if len(lines) != 3 {
		t.Errorf("Render() lines = %d, want 3", len(lines))
	}

	// First line should have reverse video cursor.
	if !strings.Contains(lines[0], "\x1b[7m") {
		t.Errorf("First line should contain reverse video escape code, got: %q", lines[0])
	}

	// Second and third lines should NOT have cursor.
	if strings.Contains(lines[1], "\x1b[7m") {
		t.Errorf("Second line should NOT contain reverse video, got: %q", lines[1])
	}
	if strings.Contains(lines[2], "\x1b[7m") {
		t.Errorf("Third line should NOT contain reverse video, got: %q", lines[2])
	}
}

func TestTextAreaRenderer_Render_Multiline_WithCursor_MiddleLine(t *testing.T) {
	r := NewTextAreaRenderer()
	ta := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("line1\nline2\nline3")).
		WithCursor(model.NewCursor(1, 2)). // Cursor on SECOND line (row 1)
		WithShowCursor(true)

	result := r.Render(ta)

	// Should render with cursor on SECOND line using reverse video.
	lines := strings.Split(result, "\n")
	if len(lines) != 3 {
		t.Errorf("Render() lines = %d, want 3", len(lines))
	}

	// CRITICAL TEST: This tests the actualRow vs cursorRow bug!
	// The bug is: actualRow (loop index) is compared with cursorRow (buffer position)
	// When cursor is at row 1, actualRow should also be 1 (second visible line)

	// First line should NOT have cursor.
	if strings.Contains(lines[0], "\x1b[7m") {
		t.Errorf("First line should NOT contain reverse video, got: %q", lines[0])
	}

	// Second line SHOULD have cursor with reverse video.
	if !strings.Contains(lines[1], "\x1b[7m") {
		t.Errorf("Second line SHOULD contain reverse video, got: %q\nFull render:\n%s", lines[1], result)
	}

	// Third line should NOT have cursor.
	if strings.Contains(lines[2], "\x1b[7m") {
		t.Errorf("Third line should NOT contain reverse video, got: %q", lines[2])
	}
}

func TestTextAreaRenderer_Render_Multiline_WithCursor_LastLine(t *testing.T) {
	r := NewTextAreaRenderer()
	ta := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("line1\nline2\nline3")).
		WithCursor(model.NewCursor(2, 4)). // Cursor on THIRD line (row 2)
		WithShowCursor(true)

	result := r.Render(ta)

	// Should render with cursor on THIRD line using reverse video.
	lines := strings.Split(result, "\n")
	if len(lines) != 3 {
		t.Errorf("Render() lines = %d, want 3", len(lines))
	}

	// CRITICAL TEST: actualRow=2 should match cursorRow=2.

	// First and second lines should NOT have cursor.
	if strings.Contains(lines[0], "\x1b[7m") {
		t.Errorf("First line should NOT contain reverse video, got: %q", lines[0])
	}
	if strings.Contains(lines[1], "\x1b[7m") {
		t.Errorf("Second line should NOT contain reverse video, got: %q", lines[1])
	}

	// Third line SHOULD have cursor with reverse video.
	if !strings.Contains(lines[2], "\x1b[7m") {
		t.Errorf("Third line SHOULD contain reverse video, got: %q\nFull render:\n%s", lines[2], result)
	}
}

func TestTextAreaRenderer_Render_Multiline_CursorAtEndOfLine(t *testing.T) {
	r := NewTextAreaRenderer()
	ta := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("hello\nworld")).
		WithCursor(model.NewCursor(1, 5)). // Cursor at end of "world" (row 1, col 5)
		WithShowCursor(true)

	result := r.Render(ta)

	lines := strings.Split(result, "\n")
	if len(lines) != 2 {
		t.Errorf("Render() lines = %d, want 2", len(lines))
	}

	// First line should NOT have cursor.
	if strings.Contains(lines[0], "\x1b[7m") {
		t.Errorf("First line should NOT contain reverse video, got: %q", lines[0])
	}

	// Second line SHOULD have cursor at END with reverse video.
	if !strings.Contains(lines[1], "\x1b[7m") {
		t.Errorf("Second line SHOULD contain reverse video, got: %q", lines[1])
	}

	// Cursor should be AFTER "world" (reverse video space).
	if !strings.Contains(lines[1], "world") {
		t.Errorf("Second line should contain 'world': %q", lines[1])
	}
	if !strings.Contains(lines[1], "\x1b[7m \x1b[27m") {
		t.Errorf("Cursor at end should be reverse video space, got: %q", lines[1])
	}
}

func TestTextAreaRenderer_Render_WithLineNumbers(t *testing.T) {
	r := NewTextAreaRenderer()
	ta := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("line1\nline2")).
		WithLineNumbers(true).
		WithShowCursor(false)

	result := r.Render(ta)

	// Should have line numbers.
	if !strings.Contains(result, "1 ") {
		t.Errorf("Render() should contain line number '1 ', got: %q", result)
	}
	if !strings.Contains(result, "2 ") {
		t.Errorf("Render() should contain line number '2 ', got: %q", result)
	}
}

func TestTextAreaRenderer_renderLineWithCursor_Middle(t *testing.T) {
	r := NewTextAreaRenderer()
	result := r.renderLineWithCursor("hello", 2)

	// Cursor at position 2 should apply reverse video to 'l'.
	// Expected: "he\x1b[7ml\x1b[27mlo" (cursor on third character with reverse video)
	expected := "he\x1b[7ml\x1b[27mlo"
	if result != expected {
		t.Errorf("renderLineWithCursor() = %q, want %q", result, expected)
	}
}

func TestTextAreaRenderer_renderLineWithCursor_End(t *testing.T) {
	r := NewTextAreaRenderer()
	result := r.renderLineWithCursor("hello", 5)

	// Cursor at end of line should append reverse video space.
	// Expected: "hello\x1b[7m \x1b[27m".
	expected := "hello\x1b[7m \x1b[27m"
	if result != expected {
		t.Errorf("renderLineWithCursor() = %q, want %q", result, expected)
	}
}

func TestTextAreaRenderer_renderLineWithCursor_Start(t *testing.T) {
	r := NewTextAreaRenderer()
	result := r.renderLineWithCursor("hello", 0)

	// Cursor at start should apply reverse video to first character.
	// Expected: "\x1b[7mh\x1b[27mello".
	expected := "\x1b[7mh\x1b[27mello"
	if result != expected {
		t.Errorf("renderLineWithCursor() = %q, want %q", result, expected)
	}
}

// TestTextAreaRenderer_Render_WithScrollOffset tests the actualRow vs cursorRow bug.
//
// BUG SCENARIO:
// When TextArea has scrolled down (scrollRow > 0), the actualRow (visible line index)
// does NOT match cursorRow (absolute buffer position).
//
// Example:
// - Buffer has 10 lines (0-9)
// - Viewport height = 3 (shows 3 lines at a time)
// - Scroll offset = 5 (viewing lines 5, 6, 7)
// - Cursor at row 6, col 2.
//
// VisibleLines() returns: lines[5], lines[6], lines[7].
// Loop: actualRow goes from 0 to 2.
// But cursorRow = 6 (absolute position in buffer)
//
// CURRENT CODE: actualRow == cursorRow → 1 == 6 → FALSE → NO CURSOR RENDERED!
// CORRECT: actualRow should be (i + scrollOffset) to match cursorRow.
func TestTextAreaRenderer_Render_WithScrollOffset(t *testing.T) {
	// Test that renderer correctly handles scroll offset for cursor positioning.
	// When scrolled, cursor should render on the correct visible line.
	r := NewTextAreaRenderer()
	ta := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("line0\nline1\nline2\nline3\nline4\nline5\nline6\nline7")).
		WithSize(80, 3).                  // Height = 3, shows 3 lines.
		WithCursor(model.NewCursor(6, 2)) // Cursor on row 6, ensureCursorVisible() auto-scrolls.

	// WithCursor() triggers ensureCursorVisible() which sets scrollRow = 4
	// (cursor on row 6, height 3, so scrollRow = 6 - 3 + 1 = 4)
	// VisibleLines will be: [line4, line5, line6]
	// Cursor should appear on line6 (index 2 in VisibleLines).

	result := r.Render(ta)
	lines := strings.Split(result, "\n")

	// Verify we have 3 lines.
	if len(lines) != 3 {
		t.Fatalf("Expected 3 lines, got %d", len(lines))
	}

	// Cursor should be on third visible line (line6 in buffer) with reverse video.
	if !strings.Contains(lines[2], "\x1b[7m") {
		t.Errorf("Line at index 2 (buffer row 6) should have reverse video cursor, got: %q", lines[2])
	}

	// First two lines should NOT have cursor.
	if strings.Contains(lines[0], "\x1b[7m") {
		t.Errorf("Line at index 0 should not have reverse video, got: %q", lines[0])
	}
	if strings.Contains(lines[1], "\x1b[7m") {
		t.Errorf("Line at index 1 should not have reverse video, got: %q", lines[1])
	}

	// Verify content is correct (line4, line5, line6).
	if !strings.HasPrefix(lines[0], "line4") {
		t.Errorf("First visible line should be line4, got: %q", lines[0])
	}
	if !strings.HasPrefix(lines[1], "line5") {
		t.Errorf("Second visible line should be line5, got: %q", lines[1])
	}
	// Cursor at col 2 should be: "li\x1b[7mn\x1b[27me6"
	if !strings.Contains(lines[2], "li\x1b[7mn\x1b[27me6") {
		t.Errorf("Third visible line should have cursor on 'n': %q", lines[2])
	}
}
