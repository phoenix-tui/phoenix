package renderer

import (
	"bytes"
	"strings"
	"sync"
	"testing"
)

// ─── Constructor ─────────────────────────────────────────────────────────────

func TestNewInlineRenderer(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 80, 24)

	if r == nil {
		t.Fatal("NewInlineRenderer returned nil")
	}
	if r.out != &buf {
		t.Error("out field not set correctly")
	}
	if r.width != 80 {
		t.Errorf("width: want 80, got %d", r.width)
	}
	if r.height != 24 {
		t.Errorf("height: want 24, got %d", r.height)
	}
	if r.linesRendered != 0 {
		t.Errorf("linesRendered should be 0 initially, got %d", r.linesRendered)
	}
}

// ─── First render ────────────────────────────────────────────────────────────

// TestInlineRenderer_FirstRender verifies the first render writes content
// correctly: no cursor-up, carriage return at start, erase-to-EOL per line,
// carriage return at end.
func TestInlineRenderer_FirstRender(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 80, 24)

	if err := r.Render("Hello\nWorld"); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	out := buf.String()

	// Must contain the text.
	if !strings.Contains(out, "Hello") {
		t.Errorf("output missing 'Hello': %q", out)
	}
	if !strings.Contains(out, "World") {
		t.Errorf("output missing 'World': %q", out)
	}

	// Must contain erase-to-EOL after each line.
	if !strings.Contains(out, eraseLineRight) {
		t.Errorf("output missing eraseLineRight (%q): %q", eraseLineRight, out)
	}

	// Must NOT contain cursor-up (first render, nothing to move up past).
	if strings.Contains(out, "\x1b[") && strings.HasPrefix(out, "\x1b[") {
		// Allow eraseLineRight and eraseScreenBelow but not cursor-up at start.
		if strings.HasPrefix(out, "\x1b[1A") || strings.HasPrefix(out, "\x1b[2A") {
			t.Errorf("first render should not start with cursor-up: %q", out)
		}
	}

	// Must end with carriage return.
	if !strings.HasSuffix(out, carriageReturn) {
		t.Errorf("output should end with carriage return, got: %q", out)
	}

	// linesRendered must equal the number of lines.
	if r.linesRendered != 2 {
		t.Errorf("linesRendered: want 2, got %d", r.linesRendered)
	}
}

// TestInlineRenderer_SingleLine verifies a single-line view does not produce
// a cursor-up sequence (nothing to move past).
func TestInlineRenderer_SingleLine(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 80, 24)

	if err := r.Render("Hello"); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	out := buf.String()

	// No cursor-up sequence.
	if strings.Contains(out, "\x1b[1A") {
		t.Errorf("single line should not produce cursor-up: %q", out)
	}

	if r.linesRendered != 1 {
		t.Errorf("linesRendered: want 1, got %d", r.linesRendered)
	}
}

// ─── Second render (changed) ─────────────────────────────────────────────────

// TestInlineRenderer_SecondRender_Changed verifies that on the second render:
//   - cursor moves up (linesRendered - 1) lines
//   - unchanged lines are skipped (not reprinted)
//   - changed lines are reprinted with eraseLineRight
func TestInlineRenderer_SecondRender_Changed(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 80, 24)

	// First render.
	if err := r.Render("Hello\nWorld"); err != nil {
		t.Fatalf("first Render error: %v", err)
	}
	buf.Reset()

	// Second render: first line unchanged, second line changed.
	if err := r.Render("Hello\nGoBrr"); err != nil {
		t.Fatalf("second Render error: %v", err)
	}

	out := buf.String()

	// Must contain cursor-up 1 line (linesRendered was 2, so 2-1=1).
	if !strings.Contains(out, "\x1b[1A") {
		t.Errorf("expected cursor-up 1, got: %q", out)
	}

	// Changed line must be present.
	if !strings.Contains(out, "GoBrr") {
		t.Errorf("changed line 'GoBrr' missing from output: %q", out)
	}

	// Unchanged line "Hello" should NOT be rewritten — only "\r\n" for advance.
	// We verify by checking "Hello" is absent from the second render output.
	if strings.Contains(out, "Hello") {
		t.Errorf("unchanged line 'Hello' should be skipped, got: %q", out)
	}

	// Changed line must have eraseLineRight.
	if !strings.Contains(out, "GoBrr"+eraseLineRight) {
		t.Errorf("changed line should have eraseLineRight appended: %q", out)
	}
}

// TestInlineRenderer_SecondRender_AllChanged verifies that when all lines
// change, all are reprinted.
func TestInlineRenderer_SecondRender_AllChanged(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 80, 24)

	if err := r.Render("Line1\nLine2"); err != nil {
		t.Fatalf("first Render error: %v", err)
	}
	buf.Reset()

	if err := r.Render("New1\nNew2"); err != nil {
		t.Fatalf("second Render error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "New1") {
		t.Errorf("expected 'New1' in output: %q", out)
	}
	if !strings.Contains(out, "New2") {
		t.Errorf("expected 'New2' in output: %q", out)
	}
}

// ─── Identical render (no-op) ─────────────────────────────────────────────────

// TestInlineRenderer_IdenticalRender verifies that rendering the exact same
// view twice results in no output on the second call.
func TestInlineRenderer_IdenticalRender(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 80, 24)

	view := "Count: 5\nStatus: ok"

	if err := r.Render(view); err != nil {
		t.Fatalf("first Render error: %v", err)
	}
	buf.Reset()

	if err := r.Render(view); err != nil {
		t.Fatalf("second Render error: %v", err)
	}

	// No output for identical render.
	if buf.Len() != 0 {
		t.Errorf("identical render should produce no output, got: %q", buf.String())
	}
}

// ─── View shrinks ────────────────────────────────────────────────────────────

// TestInlineRenderer_ViewShrinks verifies that when the new view has fewer
// lines than the previous render, eraseScreenBelow is emitted to clear the
// leftover lines.
func TestInlineRenderer_ViewShrinks(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 80, 24)

	// First render: 3 lines.
	if err := r.Render("Line1\nLine2\nLine3"); err != nil {
		t.Fatalf("first Render error: %v", err)
	}
	buf.Reset()

	// Second render: 1 line.
	if err := r.Render("Short"); err != nil {
		t.Fatalf("second Render error: %v", err)
	}

	out := buf.String()

	// Must erase extra lines below.
	if !strings.Contains(out, eraseScreenBelow) {
		t.Errorf("expected eraseScreenBelow (%q) when view shrinks: %q", eraseScreenBelow, out)
	}

	// linesRendered updated to new count.
	if r.linesRendered != 1 {
		t.Errorf("linesRendered: want 1, got %d", r.linesRendered)
	}
}

// ─── View grows ──────────────────────────────────────────────────────────────

// TestInlineRenderer_ViewGrows verifies that when the new view has more lines
// than the previous render, all new lines are rendered without eraseScreenBelow.
func TestInlineRenderer_ViewGrows(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 80, 24)

	// First render: 1 line.
	if err := r.Render("Short"); err != nil {
		t.Fatalf("first Render error: %v", err)
	}
	buf.Reset()

	// Second render: 3 lines.
	if err := r.Render("Line1\nLine2\nLine3"); err != nil {
		t.Fatalf("second Render error: %v", err)
	}

	out := buf.String()

	// All three lines present.
	if !strings.Contains(out, "Line2") {
		t.Errorf("expected 'Line2' in output: %q", out)
	}
	if !strings.Contains(out, "Line3") {
		t.Errorf("expected 'Line3' in output: %q", out)
	}

	// No eraseScreenBelow when view grows.
	if strings.Contains(out, eraseScreenBelow) {
		t.Errorf("should not erase screen below when view grows: %q", out)
	}

	if r.linesRendered != 3 {
		t.Errorf("linesRendered: want 3, got %d", r.linesRendered)
	}
}

// ─── Empty view ──────────────────────────────────────────────────────────────

// TestInlineRenderer_EmptyView verifies that an empty view is handled
// gracefully and does not panic or produce corrupt sequences.
func TestInlineRenderer_EmptyView(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 80, 24)

	if err := r.Render(""); err != nil {
		t.Fatalf("Render empty view error: %v", err)
	}

	// linesRendered must be 1 (strings.Split("", "\n") = []string{""}).
	if r.linesRendered != 1 {
		t.Errorf("linesRendered for empty view: want 1, got %d", r.linesRendered)
	}
}

// TestInlineRenderer_EmptyView_AfterContent verifies that rendering an empty
// view after real content erases the previous output.
func TestInlineRenderer_EmptyView_AfterContent(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 80, 24)

	if err := r.Render("Hello\nWorld\nEnd"); err != nil {
		t.Fatalf("first Render error: %v", err)
	}
	buf.Reset()

	if err := r.Render(""); err != nil {
		t.Fatalf("second Render empty error: %v", err)
	}

	out := buf.String()

	// 3 previous lines → 1 new line: must erase below.
	if !strings.Contains(out, eraseScreenBelow) {
		t.Errorf("expected eraseScreenBelow after switching to empty view: %q", out)
	}
}

// ─── Resize ──────────────────────────────────────────────────────────────────

// TestInlineRenderer_Resize verifies that Resize updates dimensions and forces
// a full repaint (diff cache cleared) on the next Render.
func TestInlineRenderer_Resize(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 80, 24)

	view := "Hello\nWorld"
	if err := r.Render(view); err != nil {
		t.Fatalf("first Render error: %v", err)
	}
	buf.Reset()

	r.Resize(120, 40)

	if r.width != 120 {
		t.Errorf("width after Resize: want 120, got %d", r.width)
	}
	if r.height != 40 {
		t.Errorf("height after Resize: want 40, got %d", r.height)
	}

	// Diff cache must be cleared.
	if r.lastView != "" {
		t.Error("lastView should be empty after Resize")
	}
	if r.lastLines != nil {
		t.Error("lastLines should be nil after Resize")
	}

	// Re-render identical view — should produce output because cache is cleared.
	if err := r.Render(view); err != nil {
		t.Fatalf("Render after Resize error: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("Render after Resize should produce output (full repaint)")
	}
}

// ─── Repaint ─────────────────────────────────────────────────────────────────

// TestInlineRenderer_Repaint verifies that Repaint clears the diff cache so
// the next Render performs a full repaint.
func TestInlineRenderer_Repaint(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 80, 24)

	view := "Hello\nWorld"
	if err := r.Render(view); err != nil {
		t.Fatalf("first Render error: %v", err)
	}
	buf.Reset()

	r.Repaint()

	if r.lastView != "" {
		t.Error("lastView should be empty after Repaint")
	}
	if r.lastLines != nil {
		t.Error("lastLines should be nil after Repaint")
	}

	// linesRendered must NOT be reset by Repaint (cursor position still valid).
	if r.linesRendered != 2 {
		t.Errorf("linesRendered should remain 2 after Repaint, got %d", r.linesRendered)
	}

	// Re-render same view — should produce output (full repaint triggered).
	if err := r.Render(view); err != nil {
		t.Fatalf("Render after Repaint error: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("Render after Repaint should produce output")
	}

	// Content must be present after full repaint.
	out := buf.String()
	if !strings.Contains(out, "Hello") {
		t.Errorf("expected 'Hello' after Repaint: %q", out)
	}
}

// TestInlineRenderer_Repaint_DoesNotResetLinesRendered confirms linesRendered
// is preserved so the cursor-up on the next flush is still correct.
func TestInlineRenderer_Repaint_DoesNotResetLinesRendered(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 80, 24)

	if err := r.Render("A\nB\nC"); err != nil {
		t.Fatalf("Render error: %v", err)
	}
	if r.linesRendered != 3 {
		t.Fatalf("linesRendered before Repaint: want 3, got %d", r.linesRendered)
	}

	r.Repaint()

	if r.linesRendered != 3 {
		t.Errorf("linesRendered after Repaint: want 3, got %d", r.linesRendered)
	}
}

// ─── SetOutput ───────────────────────────────────────────────────────────────

// TestInlineRenderer_SetOutput verifies that SetOutput changes the destination
// writer for future renders.
func TestInlineRenderer_SetOutput(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	r := NewInlineRenderer(&buf1, 80, 24)

	if err := r.Render("First"); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	if buf1.Len() == 0 {
		t.Fatal("buf1 should have received output")
	}

	r.SetOutput(&buf2)
	r.Repaint()

	if err := r.Render("Second"); err != nil {
		t.Fatalf("Render after SetOutput error: %v", err)
	}

	if buf2.Len() == 0 {
		t.Error("buf2 should have received output after SetOutput")
	}
	if !strings.Contains(buf2.String(), "Second") {
		t.Errorf("buf2 should contain 'Second': %q", buf2.String())
	}
}

// ─── Height clipping ─────────────────────────────────────────────────────────

// TestInlineRenderer_HeightClipping verifies that a view taller than the
// terminal height is clipped from the top, preserving the bottom lines.
func TestInlineRenderer_HeightClipping(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 80, 3) // only 3 lines tall

	// 5-line view.
	view := "L1\nL2\nL3\nL4\nL5"
	if err := r.Render(view); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	out := buf.String()

	// Bottom 3 lines should be present.
	if !strings.Contains(out, "L3") {
		t.Errorf("expected 'L3' in clipped output: %q", out)
	}
	if !strings.Contains(out, "L5") {
		t.Errorf("expected 'L5' in clipped output: %q", out)
	}

	// Top 2 lines should be absent.
	if strings.Contains(out, "L1") {
		t.Errorf("'L1' should be clipped: %q", out)
	}
	if strings.Contains(out, "L2") {
		t.Errorf("'L2' should be clipped: %q", out)
	}

	// linesRendered must be 3 (the clipped height), not 5.
	if r.linesRendered != 3 {
		t.Errorf("linesRendered after clipping: want 3, got %d", r.linesRendered)
	}
}

// TestInlineRenderer_ZeroHeight verifies that height=0 disables clipping.
func TestInlineRenderer_ZeroHeight(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 80, 0) // height 0 = no clipping

	view := "L1\nL2\nL3\nL4\nL5"
	if err := r.Render(view); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	out := buf.String()

	// All lines present.
	for _, line := range []string{"L1", "L2", "L3", "L4", "L5"} {
		if !strings.Contains(out, line) {
			t.Errorf("expected %q in output with zero height: %q", line, out)
		}
	}

	if r.linesRendered != 5 {
		t.Errorf("linesRendered: want 5, got %d", r.linesRendered)
	}
}

// ─── Width truncation ─────────────────────────────────────────────────────────

// TestInlineRenderer_WidthTruncation verifies that lines longer than the
// terminal width are truncated.
func TestInlineRenderer_WidthTruncation(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 5, 24) // only 5 columns wide

	if err := r.Render("Hello World"); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	out := buf.String()

	// "Hello" fits in 5 columns; " World" should be cut off.
	if strings.Contains(out, "World") {
		t.Errorf("'World' should be truncated: %q", out)
	}
	if !strings.Contains(out, "Hello") {
		t.Errorf("'Hello' should be present: %q", out)
	}
}

// TestInlineRenderer_ZeroWidth verifies that width=0 disables truncation.
func TestInlineRenderer_ZeroWidth(t *testing.T) {
	var buf bytes.Buffer
	r := NewInlineRenderer(&buf, 0, 24) // width 0 = no truncation

	long := strings.Repeat("X", 200)
	if err := r.Render(long); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// Full line should be present.
	if !strings.Contains(buf.String(), long) {
		t.Errorf("line should not be truncated with width=0")
	}
}

// ─── Write error propagation ──────────────────────────────────────────────────

// TestInlineRenderer_WriteError verifies that a write error is returned from
// Render and that internal state is not updated (failed render is retryable).
func TestInlineRenderer_WriteError(t *testing.T) {
	r := NewInlineRenderer(&errorWriter{}, 80, 24)

	err := r.Render("Hello")
	if err == nil {
		t.Fatal("expected write error, got nil")
	}

	// State must not be updated on error.
	if r.linesRendered != 0 {
		t.Errorf("linesRendered should be 0 after failed write, got %d", r.linesRendered)
	}
	if r.lastView != "" {
		t.Errorf("lastView should be empty after failed write, got %q", r.lastView)
	}
}

// errorWriter always returns an error on Write.
type errorWriter struct{}

func (errorWriter) Write(_ []byte) (int, error) {
	return 0, &writeError{}
}

type writeError struct{}

func (writeError) Error() string { return "simulated write error" }

// ─── Concurrent access ────────────────────────────────────────────────────────

// TestInlineRenderer_ConcurrentAccess verifies that the mutex protects
// concurrent calls to Render, Repaint, Resize, and SetOutput.
func TestInlineRenderer_ConcurrentAccess(t *testing.T) {
	t.Parallel()
	var buf syncBuffer
	r := NewInlineRenderer(&buf, 80, 24)

	var wg sync.WaitGroup
	const goroutines = 20

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			switch n % 4 {
			case 0:
				_ = r.Render("Hello\nWorld")
			case 1:
				r.Repaint()
			case 2:
				r.Resize(80, 24)
			case 3:
				r.SetOutput(&buf)
			}
		}(i)
	}

	wg.Wait()
	// No data race or panic = pass.
}

// syncBuffer is a bytes.Buffer protected by a mutex, suitable for concurrent writes.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

// ─── Internal helper tests ────────────────────────────────────────────────────

// TestStripANSI verifies ANSI escape sequence removal.
func TestStripANSI(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"Hello", "Hello"},
		{"\x1b[31mRed\x1b[0m", "Red"},
		{"\x1b[1;32mBold Green\x1b[0m", "Bold Green"},
		{"No escapes", "No escapes"},
		{"\x1b[K", ""},
		{"\x1b[J", ""},
		{"\x1b[2A", ""},
		{"Hello\x1b[K World", "Hello World"},
	}

	for _, c := range cases {
		got := stripANSI(c.input)
		if got != c.want {
			t.Errorf("stripANSI(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

// TestVisualWidth verifies column-width calculation.
func TestVisualWidth(t *testing.T) {
	cases := []struct {
		input string
		want  int
	}{
		{"Hello", 5},
		{"", 0},
		{"\x1b[31mHello\x1b[0m", 5}, // ANSI codes ignored
		{"日本語", 6},                  // 3 × 2 columns
		{"A日B", 4},                  // 1 + 2 + 1
	}

	for _, c := range cases {
		got := visualWidth(c.input)
		if got != c.want {
			t.Errorf("visualWidth(%q) = %d, want %d", c.input, got, c.want)
		}
	}
}

// TestTruncateLine verifies line truncation respects display width.
func TestTruncateLine(t *testing.T) {
	cases := []struct {
		input    string
		maxWidth int
		wantLen  int // expected visual width of result
		wantHas  string
		wantNot  string
	}{
		{"Hello World", 5, 5, "Hello", "World"},
		{"Hello", 10, 5, "Hello", ""}, // shorter than max — unchanged
		{"Hello", 5, 5, "Hello", ""},  // exactly max — unchanged
		{"日本語", 4, 4, "日本", "語"},      // wide chars
		{"A日B", 3, 3, "A日", "B"},      // mixed
	}

	for _, c := range cases {
		got := truncateLine(c.input, c.maxWidth)
		gotWidth := visualWidth(got)
		if gotWidth > c.maxWidth {
			t.Errorf("truncateLine(%q, %d) visual width %d > max %d; got %q",
				c.input, c.maxWidth, gotWidth, c.maxWidth, got)
		}
		if gotWidth != c.wantLen {
			t.Errorf("truncateLine(%q, %d) width = %d, want %d; got %q",
				c.input, c.maxWidth, gotWidth, c.wantLen, got)
		}
		if c.wantHas != "" && !strings.Contains(got, c.wantHas) {
			t.Errorf("truncateLine(%q, %d) should contain %q; got %q",
				c.input, c.maxWidth, c.wantHas, got)
		}
		if c.wantNot != "" && strings.Contains(got, c.wantNot) {
			t.Errorf("truncateLine(%q, %d) should not contain %q; got %q",
				c.input, c.maxWidth, c.wantNot, got)
		}
	}
}

// TestTruncateLine_PreservesANSI verifies that ANSI escape sequences are passed
// through without counting toward the width budget.
func TestTruncateLine_PreservesANSI(t *testing.T) {
	// Red "Hello" followed by reset — visual width is 5 but string is longer.
	input := "\x1b[31mHello\x1b[0m World"
	got := truncateLine(input, 5)

	// Visual width must not exceed 5.
	if w := visualWidth(got); w > 5 {
		t.Errorf("visual width after truncation = %d, want <= 5; got %q", w, got)
	}

	// ANSI codes should still be present in the output.
	if !strings.Contains(got, "\x1b[31m") {
		t.Errorf("ANSI opening sequence should be preserved: %q", got)
	}
}

// ─── Sequence correctness ─────────────────────────────────────────────────────

// TestInlineRenderer_CursorUpCorrectness verifies the exact cursor-up count
// for multi-line renders.
func TestInlineRenderer_CursorUpCorrectness(t *testing.T) {
	cases := []struct {
		firstView  string
		secondView string
		wantUp     string // expected cursor-up sequence, empty = none
	}{
		{"A", "B", ""},                    // 1 line → no cursor-up
		{"A\nB", "C\nD", "\x1b[1A"},       // 2 lines → cursor-up 1
		{"A\nB\nC", "D\nE\nF", "\x1b[2A"}, // 3 lines → cursor-up 2
	}

	for _, c := range cases {
		var buf bytes.Buffer
		r := NewInlineRenderer(&buf, 80, 24)

		if err := r.Render(c.firstView); err != nil {
			t.Fatalf("first Render error: %v", err)
		}
		buf.Reset()

		if err := r.Render(c.secondView); err != nil {
			t.Fatalf("second Render error: %v", err)
		}

		out := buf.String()
		if c.wantUp == "" {
			hasCursorUp := strings.Contains(out, "\x1b[1A") || strings.Contains(out, "\x1b[2A")
			if hasCursorUp {
				t.Errorf("firstView=%q → secondView=%q: unexpected cursor-up in %q",
					c.firstView, c.secondView, out)
			}
			continue
		}
		if !strings.Contains(out, c.wantUp) {
			t.Errorf("firstView=%q → secondView=%q: want cursor-up %q, got %q",
				c.firstView, c.secondView, c.wantUp, out)
		}
	}
}

// ─── cursorUp edge cases ──────────────────────────────────────────────────────

// TestCursorUp verifies the helper function for cursor-up sequences.
func TestCursorUp(t *testing.T) {
	cases := []struct {
		n    int
		want string
	}{
		{0, ""},
		{-1, ""},
		{-100, ""},
		{1, "\x1b[1A"},
		{5, "\x1b[5A"},
		{100, "\x1b[100A"},
	}
	for _, c := range cases {
		got := cursorUp(c.n)
		if got != c.want {
			t.Errorf("cursorUp(%d) = %q, want %q", c.n, got, c.want)
		}
	}
}

// ─── stripANSI additional coverage ───────────────────────────────────────────

// TestStripANSI_OSC verifies that OSC sequences (window title etc.) are stripped.
func TestStripANSI_OSC(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		// OSC terminated by BEL
		{"\x1b]0;Title\x07Text", "Text"},
		// OSC terminated by ST (ESC \)
		{"\x1b]0;Title\x1b\\Text", "Text"},
		// Bare ESC at end of string
		{"\x1b", ""},
		// Two-byte ESC sequences (not CSI)
		{"\x1b7Text\x1b8", "Text"},
		// ESC at end of string (truncated)
		{"Hello\x1b", "Hello"},
	}
	for _, c := range cases {
		got := stripANSI(c.input)
		if got != c.want {
			t.Errorf("stripANSI(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

// ─── decodeRune coverage ──────────────────────────────────────────────────────

// TestDecodeRune verifies UTF-8 rune decoding including error cases.
func TestDecodeRune(t *testing.T) {
	cases := []struct {
		s        string
		i        int
		wantRune rune
		wantSize int
	}{
		// ASCII
		{"Hello", 0, 'H', 1},
		{"Hello", 4, 'o', 1},
		// 2-byte UTF-8 (é = U+00E9)
		{"\xc3\xa9", 0, 'é', 2},
		// 3-byte UTF-8 (日 = U+65E5)
		{"\xe6\x97\xa5", 0, '日', 3},
		// 4-byte UTF-8 (𝄞 = U+1D11E)
		{"\xf0\x9d\x84\x9e", 0, '𝄞', 4},
		// Invalid lead byte
		{"\xff", 0, '\uFFFD', 1},
		// Truncated 2-byte sequence
		{"\xc3", 0, '\uFFFD', 1},
		// Invalid continuation byte
		{"\xc3\x00", 0, '\uFFFD', 1},
	}

	for _, c := range cases {
		gotRune, gotSize := decodeRune(c.s, c.i)
		if gotRune != c.wantRune {
			t.Errorf("decodeRune(%q, %d) rune = %q (%d), want %q (%d)",
				c.s, c.i, gotRune, gotRune, c.wantRune, c.wantRune)
		}
		if gotSize != c.wantSize {
			t.Errorf("decodeRune(%q, %d) size = %d, want %d",
				c.s, c.i, gotSize, c.wantSize)
		}
	}
}

// ─── runeDisplayWidth coverage ────────────────────────────────────────────────

// TestRuneDisplayWidth verifies display width for representative Unicode ranges.
func TestRuneDisplayWidth(t *testing.T) {
	cases := []struct {
		r    rune
		want int
		desc string
	}{
		// ASCII printable
		{'A', 1, "ASCII letter"},
		{' ', 1, "ASCII space"},
		// Control characters
		{'\n', 0, "newline"},
		{'\t', 0, "tab"},
		{'\x1b', 0, "ESC"},
		{'\x80', 0, "C1 control start"},
		{'\x9F', 0, "C1 control end"},
		// Hangul Jamo (wide)
		{'\u1100', 2, "Hangul Jamo start"},
		{'\u115F', 2, "Hangul Jamo end"},
		// Angle brackets (wide)
		{'\u2329', 2, "left angle bracket"},
		{'\u232A', 2, "right angle bracket"},
		// CJK Radicals (wide)
		{'\u2E80', 2, "CJK Radicals start"},
		// Japanese Hiragana (wide)
		{'\u3041', 2, "Hiragana"},
		// CJK Extension A (wide)
		{'\u3400', 2, "CJK Ext-A start"},
		// CJK Unified Ideographs (wide)
		{'\u4E00', 2, "CJK start"},
		{'\u9FFF', 2, "CJK end area"},
		// Hangul Extended-A (wide)
		{'\uA960', 2, "Hangul Ext-A"},
		// Hangul Syllables (wide)
		{'\uAC00', 2, "Hangul Syllables start"},
		// CJK Compatibility Ideographs (wide)
		{'\uF900', 2, "CJK Compat start"},
		// Vertical forms (wide)
		{'\uFE10', 2, "Vertical forms start"},
		// CJK Compat Forms (wide)
		{'\uFE30', 2, "CJK Compat Forms start"},
		// Fullwidth Forms (wide)
		{'\uFF01', 2, "Fullwidth exclamation"},
		// Fullwidth Signs (wide)
		{'\uFFE0', 2, "Fullwidth signs start"},
		// Emoji (wide)
		{'\U0001F600', 2, "Emoji grinning face"},
		// Combining mark (zero width)
		{'\u0300', 0, "combining grave accent"},
		// Regular Latin (1 column)
		{'z', 1, "Latin small z"},
	}

	for _, c := range cases {
		got := runeDisplayWidth(c.r)
		if got != c.want {
			t.Errorf("runeDisplayWidth(%q %s U+%04X) = %d, want %d",
				c.r, c.desc, c.r, got, c.want)
		}
	}
}

// TestRuneDisplayWidth_AdditionalRanges covers ranges not hit by the basic test.
func TestRuneDisplayWidth_AdditionalRanges(t *testing.T) {
	cases := []struct {
		r    rune
		want int
	}{
		// Kana Supplement
		{'\U0001B000', 2},
		// Playing cards / Mahjong tiles
		{'\U0001F004', 2},
		// CJK Extension B-F
		{'\U00020000', 2},
		{'\U0002FFFD', 2},
		// CJK Extension G+
		{'\U00030000', 2},
		{'\U0003FFFD', 2},
		// Normal Latin character (should be 1, not wide)
		{'a', 1},
	}

	for _, c := range cases {
		got := runeDisplayWidth(c.r)
		if got != c.want {
			t.Errorf("runeDisplayWidth(U+%04X) = %d, want %d", c.r, got, c.want)
		}
	}
}

// ─── truncateLine OSC coverage ───────────────────────────────────────────────

// TestTruncateLine_OSCSequence verifies that OSC sequences in a line are
// passed through without counting toward the width budget.
func TestTruncateLine_OSCSequence(t *testing.T) {
	// OSC window title followed by visible text
	input := "\x1b]0;Title\x07Hello World"
	got := truncateLine(input, 5)

	// Visual width must not exceed 5.
	if w := visualWidth(got); w > 5 {
		t.Errorf("visual width = %d, want <= 5; got %q", w, got)
	}
	// "Hello" should be present.
	if !strings.Contains(got, "Hello") {
		t.Errorf("expected 'Hello' in output: %q", got)
	}
	// OSC sequence should be preserved.
	if !strings.Contains(got, "\x1b]0;Title\x07") {
		t.Errorf("OSC sequence should be preserved: %q", got)
	}
}

// TestTruncateLine_TwoByteESC verifies two-byte ESC sequences (e.g., DECSC)
// are passed through without counting toward the width budget.
func TestTruncateLine_TwoByteESC(t *testing.T) {
	// DECSC (save cursor) + text
	input := "\x1b7Hello World\x1b8"
	got := truncateLine(input, 5)

	if w := visualWidth(got); w > 5 {
		t.Errorf("visual width = %d, want <= 5; got %q", w, got)
	}
	if !strings.Contains(got, "Hello") {
		t.Errorf("expected 'Hello' in output: %q", got)
	}
}

// TestTruncateLine_ESCAtEndOfString verifies truncateLine handles a trailing
// ESC byte gracefully.
func TestTruncateLine_ESCAtEndOfString(t *testing.T) {
	input := "Hi\x1b"
	// Should not panic; result may or may not include trailing ESC.
	got := truncateLine(input, 80)
	if !strings.HasPrefix(got, "Hi") {
		t.Errorf("expected result to start with 'Hi': %q", got)
	}
}

// TestTruncateLine_OSCSequence_STTerminated verifies that OSC sequences
// terminated by ST (ESC \) are handled correctly in truncateLine.
func TestTruncateLine_OSCSequence_STTerminated(t *testing.T) {
	// OSC terminated by String Terminator (ESC \)
	input := "\x1b]0;Title\x1b\\Hello World"
	got := truncateLine(input, 5)

	if w := visualWidth(got); w > 5 {
		t.Errorf("visual width = %d, want <= 5; got %q", w, got)
	}
	if !strings.Contains(got, "Hello") {
		t.Errorf("expected 'Hello' in output: %q", got)
	}
}

// TestTruncateLine_ESCAtEnd_Short verifies the path where ESC appears
// at the end of the string in the main loop of truncateLine.
func TestTruncateLine_ESCAtEnd_Short(t *testing.T) {
	// String that is shorter than maxWidth but ends with ESC — should return unchanged.
	input := "Hi\x1b"
	got := truncateLine(input, 2) // visual width 2 = "Hi", but ESC is still there in original
	// visualWidth("Hi\x1b") strips ANSI → "Hi" = 2 columns which equals maxWidth.
	// So the early-return branch fires and we get the original string back.
	if got != input {
		t.Errorf("truncateLine(%q, 2): expected original string back (fits in maxWidth), got %q", input, got)
	}
}
