package renderer

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/components/input/textarea/domain/model"
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

	// Should show placeholder.
	if result != "Enter text..." {
		t.Errorf("Render() placeholder = %q, want %q", result, "Enter text...")
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

	// MUST NOT contain cursor.
	if strings.Contains(result, "█") {
		t.Errorf("Render() with ShowCursor(false) contains cursor '█'")
	}
}

func TestTextAreaRenderer_Render_SingleLine_WithCursor(t *testing.T) {
	r := NewTextAreaRenderer()
	ta := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("hello")).
		WithShowCursor(true) // Enable cursor (default)

	result := r.Render(ta)

	// Should render text WITH cursor at position 0.
	// Expected: "█ello" (cursor replaces 'h')
	if !strings.Contains(result, "█") {
		t.Errorf("Render() with ShowCursor(true) should contain cursor '█', got: %q", result)
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

	// MUST NOT contain cursor.
	if strings.Contains(result, "█") {
		t.Errorf("Render() with ShowCursor(false) contains cursor '█'")
	}
}

func TestTextAreaRenderer_Render_Multiline_WithCursor_FirstLine(t *testing.T) {
	r := NewTextAreaRenderer()
	ta := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("line1\nline2\nline3")).
		WithCursor(model.NewCursor(0, 3)). // Cursor on first line at col 3
		WithShowCursor(true)

	result := r.Render(ta)

	// Should render with cursor on FIRST line.
	lines := strings.Split(result, "\n")
	if len(lines) != 3 {
		t.Errorf("Render() lines = %d, want 3", len(lines))
	}

	// First line should have cursor.
	if !strings.Contains(lines[0], "█") {
		t.Errorf("First line should contain cursor, got: %q", lines[0])
	}

	// Second and third lines should NOT have cursor.
	if strings.Contains(lines[1], "█") {
		t.Errorf("Second line should NOT contain cursor, got: %q", lines[1])
	}
	if strings.Contains(lines[2], "█") {
		t.Errorf("Third line should NOT contain cursor, got: %q", lines[2])
	}
}

func TestTextAreaRenderer_Render_Multiline_WithCursor_MiddleLine(t *testing.T) {
	r := NewTextAreaRenderer()
	ta := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("line1\nline2\nline3")).
		WithCursor(model.NewCursor(1, 2)). // Cursor on SECOND line (row 1)
		WithShowCursor(true)

	result := r.Render(ta)

	// Should render with cursor on SECOND line.
	lines := strings.Split(result, "\n")
	if len(lines) != 3 {
		t.Errorf("Render() lines = %d, want 3", len(lines))
	}

	// CRITICAL TEST: This tests the actualRow vs cursorRow bug!
	// The bug is: actualRow (loop index) is compared with cursorRow (buffer position)
	// When cursor is at row 1, actualRow should also be 1 (second visible line)

	// First line should NOT have cursor.
	if strings.Contains(lines[0], "█") {
		t.Errorf("First line should NOT contain cursor, got: %q", lines[0])
	}

	// Second line SHOULD have cursor.
	if !strings.Contains(lines[1], "█") {
		t.Errorf("Second line SHOULD contain cursor, got: %q\nFull render:\n%s", lines[1], result)
	}

	// Third line should NOT have cursor.
	if strings.Contains(lines[2], "█") {
		t.Errorf("Third line should NOT contain cursor, got: %q", lines[2])
	}
}

func TestTextAreaRenderer_Render_Multiline_WithCursor_LastLine(t *testing.T) {
	r := NewTextAreaRenderer()
	ta := model.NewTextArea().
		WithBuffer(model.NewBufferFromString("line1\nline2\nline3")).
		WithCursor(model.NewCursor(2, 4)). // Cursor on THIRD line (row 2)
		WithShowCursor(true)

	result := r.Render(ta)

	// Should render with cursor on THIRD line.
	lines := strings.Split(result, "\n")
	if len(lines) != 3 {
		t.Errorf("Render() lines = %d, want 3", len(lines))
	}

	// CRITICAL TEST: actualRow=2 should match cursorRow=2.

	// First and second lines should NOT have cursor.
	if strings.Contains(lines[0], "█") {
		t.Errorf("First line should NOT contain cursor, got: %q", lines[0])
	}
	if strings.Contains(lines[1], "█") {
		t.Errorf("Second line should NOT contain cursor, got: %q", lines[1])
	}

	// Third line SHOULD have cursor.
	if !strings.Contains(lines[2], "█") {
		t.Errorf("Third line SHOULD contain cursor, got: %q\nFull render:\n%s", lines[2], result)
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
	if strings.Contains(lines[0], "█") {
		t.Errorf("First line should NOT contain cursor, got: %q", lines[0])
	}

	// Second line SHOULD have cursor at END.
	if !strings.Contains(lines[1], "█") {
		t.Errorf("Second line SHOULD contain cursor, got: %q", lines[1])
	}

	// Cursor should be AFTER "world".
	if !strings.HasSuffix(lines[1], "world█") {
		t.Errorf("Cursor should be at end: %q", lines[1])
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

	// Cursor at position 2 should replace 'l'.
	// Expected: "he█lo" (cursor replaces third character)
	if result != "he█lo" {
		t.Errorf("renderLineWithCursor() = %q, want %q", result, "he█lo")
	}
}

func TestTextAreaRenderer_renderLineWithCursor_End(t *testing.T) {
	r := NewTextAreaRenderer()
	result := r.renderLineWithCursor("hello", 5)

	// Cursor at end of line should append.
	// Expected: "hello█".
	if result != "hello█" {
		t.Errorf("renderLineWithCursor() = %q, want %q", result, "hello█")
	}
}

func TestTextAreaRenderer_renderLineWithCursor_Start(t *testing.T) {
	r := NewTextAreaRenderer()
	result := r.renderLineWithCursor("hello", 0)

	// Cursor at start should replace first character.
	// Expected: "█ello".
	if result != "█ello" {
		t.Errorf("renderLineWithCursor() = %q, want %q", result, "█ello")
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
	// Note: We can't easily create a scrolled TextArea via public API.
	// because scroll state is private and managed by ensureCursorVisible().
	//
	// This test documents the expected behavior if scroll were implemented.
	// The renderer assumes actualRow == cursorRow, which breaks with scroll.
	t.Skip("Scroll offset not yet implemented - renderer has TODO for this")

	// TODO: When scroll is implemented, add this test:
	// r := NewTextAreaRenderer()
	// ta := model.NewTextArea().
	// 	WithBuffer(model.NewBufferFromString("line0\nline1\nline2\nline3\nline4\nline5\nline6\nline7")).
	// 	WithSize(80, 3). // Height = 3, shows 3 lines.
	// 	WithCursor(model.NewCursor(6, 2)). // Cursor on row 6.
	// 	withScrollOffset(5) // Scroll down to show lines 5, 6, 7.
	//
	// result := r.Render(ta)
	// lines := strings.Split(result, "\n")
	//
	// // Second visible line (index 1 in VisibleLines, which is line 6 in buffer)
	// // SHOULD have cursor
	// if !strings.Contains(lines[1], "█") {.
	// 	t.Errorf("Line at index 1 (buffer row 6) should have cursor, got: %q", lines[1])
	// }
}
