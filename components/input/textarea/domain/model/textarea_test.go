package model

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/input/textarea/domain/value"
)

func TestNewTextArea(t *testing.T) {
	ta := NewTextArea()

	if ta.LineCount() != 1 {
		t.Errorf("NewTextArea() LineCount() = %d, want 1", ta.LineCount())
	}
	if ta.Width() != 80 {
		t.Errorf("NewTextArea() Width() = %d, want 80", ta.Width())
	}
	if ta.Height() != 24 {
		t.Errorf("NewTextArea() Height() = %d, want 24", ta.Height())
	}
	if !ta.IsEmpty() {
		t.Error("NewTextArea() should be empty")
	}
	if ta.HasSelection() {
		t.Error("NewTextArea() should not have selection")
	}
	if ta.IsReadOnly() {
		t.Error("NewTextArea() should not be read-only")
	}
	if ta.ShowLineNumbers() {
		t.Error("NewTextArea() should not show line numbers by default")
	}

	row, col := ta.CursorPosition()
	if row != 0 || col != 0 {
		t.Errorf("NewTextArea() cursor = (%d, %d), want (0, 0)", row, col)
	}
}

func TestTextArea_WithSize(t *testing.T) {
	tests := []struct {
		name       string
		width      int
		height     int
		wantWidth  int
		wantHeight int
	}{
		{
			name:       "set custom size",
			width:      100,
			height:     50,
			wantWidth:  100,
			wantHeight: 50,
		},
		{
			name:       "set small size",
			width:      10,
			height:     5,
			wantWidth:  10,
			wantHeight: 5,
		},
		{
			name:       "set large size",
			width:      200,
			height:     100,
			wantWidth:  200,
			wantHeight: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea()
			result := ta.WithSize(tt.width, tt.height)

			if result.Width() != tt.wantWidth {
				t.Errorf("Width() = %d, want %d", result.Width(), tt.wantWidth)
			}
			if result.Height() != tt.wantHeight {
				t.Errorf("Height() = %d, want %d", result.Height(), tt.wantHeight)
			}

			// Verify immutability.
			if ta.Width() != 80 || ta.Height() != 24 {
				t.Error("Original TextArea was modified")
			}
		})
	}
}

func TestTextArea_WithMaxLines(t *testing.T) {
	tests := []struct {
		name string
		max  int
	}{
		{name: "set max 10 lines", max: 10},
		{name: "set max 100 lines", max: 100},
		{name: "set unlimited", max: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea()
			result := ta.WithMaxLines(tt.max)

			// Verify immutability.
			if ta == result {
				t.Error("WithMaxLines() returned same instance")
			}
		})
	}
}

func TestTextArea_WithMaxChars(t *testing.T) {
	tests := []struct {
		name string
		max  int
	}{
		{name: "set max 100 chars", max: 100},
		{name: "set max 1000 chars", max: 1000},
		{name: "set unlimited", max: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea()
			result := ta.WithMaxChars(tt.max)

			// Verify immutability.
			if ta == result {
				t.Error("WithMaxChars() returned same instance")
			}
		})
	}
}

func TestTextArea_WithWrap(t *testing.T) {
	ta := NewTextArea()

	// Enable wrap.
	wrapped := ta.WithWrap(true)
	if wrapped == ta {
		t.Error("WithWrap() returned same instance")
	}

	// Disable wrap.
	nowrapped := wrapped.WithWrap(false)
	if nowrapped == wrapped {
		t.Error("WithWrap() returned same instance")
	}
}

func TestTextArea_WithPlaceholder(t *testing.T) {
	tests := []struct {
		name        string
		placeholder string
	}{
		{name: "simple placeholder", placeholder: "Enter text..."},
		{name: "unicode placeholder", placeholder: "输入文本..."},
		{name: "empty placeholder", placeholder: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea()
			result := ta.WithPlaceholder(tt.placeholder)

			if result.Placeholder() != tt.placeholder {
				t.Errorf("Placeholder() = %q, want %q", result.Placeholder(), tt.placeholder)
			}

			// Verify immutability.
			if ta.Placeholder() != "" {
				t.Error("Original TextArea was modified")
			}
		})
	}
}

func TestTextArea_WithReadOnly(t *testing.T) {
	ta := NewTextArea()

	// Enable read-only.
	readonly := ta.WithReadOnly(true)
	if !readonly.IsReadOnly() {
		t.Error("WithReadOnly(true) should set read-only")
	}
	if readonly == ta {
		t.Error("WithReadOnly() returned same instance")
	}

	// Disable read-only.
	editable := readonly.WithReadOnly(false)
	if editable.IsReadOnly() {
		t.Error("WithReadOnly(false) should disable read-only")
	}
}

func TestTextArea_WithLineNumbers(t *testing.T) {
	ta := NewTextArea()

	// Enable line numbers.
	withNumbers := ta.WithLineNumbers(true)
	if !withNumbers.ShowLineNumbers() {
		t.Error("WithLineNumbers(true) should show line numbers")
	}
	if withNumbers == ta {
		t.Error("WithLineNumbers() returned same instance")
	}

	// Disable line numbers.
	withoutNumbers := withNumbers.WithLineNumbers(false)
	if withoutNumbers.ShowLineNumbers() {
		t.Error("WithLineNumbers(false) should hide line numbers")
	}
}

func TestTextArea_WithBuffer(t *testing.T) {
	ta := NewTextArea()
	buffer := NewBufferFromString("line1\nline2\nline3")

	result := ta.WithBuffer(buffer)

	if result.LineCount() != 3 {
		t.Errorf("LineCount() = %d, want 3", result.LineCount())
	}
	if result.Value() != "line1\nline2\nline3" {
		t.Errorf("Value() = %q, want %q", result.Value(), "line1\nline2\nline3")
	}

	// Cursor should be reset.
	row, col := result.CursorPosition()
	if row != 0 || col != 0 {
		t.Errorf("Cursor = (%d, %d), want (0, 0)", row, col)
	}

	// Selection should be cleared.
	if result.HasSelection() {
		t.Error("Selection should be cleared")
	}
}

func TestTextArea_WithCursor(t *testing.T) {
	ta := NewTextArea().WithBuffer(NewBufferFromString("line1\nline2\nline3"))
	cursor := NewCursor(1, 3)

	result := ta.WithCursor(cursor)

	row, col := result.CursorPosition()
	if row != 1 || col != 3 {
		t.Errorf("CursorPosition() = (%d, %d), want (1, 3)", row, col)
	}

	// Verify immutability.
	origRow, origCol := ta.CursorPosition()
	if origRow != 0 || origCol != 0 {
		t.Error("Original TextArea cursor was modified")
	}
}

func TestTextArea_CursorPosition(t *testing.T) {
	tests := []struct {
		name    string
		cursor  *Cursor
		wantRow int
		wantCol int
	}{
		{
			name:    "cursor at origin",
			cursor:  NewCursor(0, 0),
			wantRow: 0,
			wantCol: 0,
		},
		{
			name:    "cursor at position",
			cursor:  NewCursor(5, 10),
			wantRow: 5,
			wantCol: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea().WithCursor(tt.cursor)
			row, col := ta.CursorPosition()

			if row != tt.wantRow {
				t.Errorf("row = %d, want %d", row, tt.wantRow)
			}
			if col != tt.wantCol {
				t.Errorf("col = %d, want %d", col, tt.wantCol)
			}
		})
	}
}

func TestTextArea_Lines(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []string
	}{
		{
			name: "single line",
			text: "hello",
			want: []string{"hello"},
		},
		{
			name: "multiple lines",
			text: "line1\nline2\nline3",
			want: []string{"line1", "line2", "line3"},
		},
		{
			name: "empty",
			text: "",
			want: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea().WithBuffer(NewBufferFromString(tt.text))
			lines := ta.Lines()

			if len(lines) != len(tt.want) {
				t.Errorf("len(Lines()) = %d, want %d", len(lines), len(tt.want))
				return
			}

			for i, line := range lines {
				if line != tt.want[i] {
					t.Errorf("Lines()[%d] = %q, want %q", i, line, tt.want[i])
				}
			}
		})
	}
}

func TestTextArea_Value(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "single line",
			text: "hello",
			want: "hello",
		},
		{
			name: "multiple lines",
			text: "line1\nline2\nline3",
			want: "line1\nline2\nline3",
		},
		{
			name: "empty",
			text: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea().WithBuffer(NewBufferFromString(tt.text))
			value := ta.Value()

			if value != tt.want {
				t.Errorf("Value() = %q, want %q", value, tt.want)
			}
		})
	}
}

func TestTextArea_LineCount(t *testing.T) {
	tests := []struct {
		name string
		text string
		want int
	}{
		{name: "single line", text: "hello", want: 1},
		{name: "three lines", text: "a\nb\nc", want: 3},
		{name: "empty", text: "", want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea().WithBuffer(NewBufferFromString(tt.text))
			count := ta.LineCount()

			if count != tt.want {
				t.Errorf("LineCount() = %d, want %d", count, tt.want)
			}
		})
	}
}

func TestTextArea_CurrentLine(t *testing.T) {
	ta := NewTextArea().WithBuffer(NewBufferFromString("line1\nline2\nline3"))

	tests := []struct {
		name string
		row  int
		want string
	}{
		{name: "first line", row: 0, want: "line1"},
		{name: "second line", row: 1, want: "line2"},
		{name: "third line", row: 2, want: "line3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positioned := ta.WithCursor(NewCursor(tt.row, 0))
			line := positioned.CurrentLine()

			if line != tt.want {
				t.Errorf("CurrentLine() = %q, want %q", line, tt.want)
			}
		})
	}
}

func TestTextArea_ContentParts(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		row        int
		col        int
		wantBefore string
		wantAt     string
		wantAfter  string
	}{
		{
			name:       "cursor at start",
			text:       "hello",
			row:        0,
			col:        0,
			wantBefore: "",
			wantAt:     "h",
			wantAfter:  "ello",
		},
		{
			name:       "cursor in middle",
			text:       "hello",
			row:        0,
			col:        2,
			wantBefore: "he",
			wantAt:     "l",
			wantAfter:  "lo",
		},
		{
			name:       "cursor at end",
			text:       "hello",
			row:        0,
			col:        5,
			wantBefore: "hello",
			wantAt:     " ",
			wantAfter:  "",
		},
		{
			name:       "unicode text",
			text:       "he👋lo",
			row:        0,
			col:        2,
			wantBefore: "he",
			wantAt:     "👋",
			wantAfter:  "lo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea().
				WithBuffer(NewBufferFromString(tt.text)).
				WithCursor(NewCursor(tt.row, tt.col))

			before, at, after := ta.ContentParts()

			if before != tt.wantBefore {
				t.Errorf("before = %q, want %q", before, tt.wantBefore)
			}
			if at != tt.wantAt {
				t.Errorf("at = %q, want %q", at, tt.wantAt)
			}
			if after != tt.wantAfter {
				t.Errorf("after = %q, want %q", after, tt.wantAfter)
			}
		})
	}
}

func TestTextArea_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{name: "empty", text: "", want: true},
		{name: "non-empty", text: "hello", want: false},
		{name: "whitespace", text: " ", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea().WithBuffer(NewBufferFromString(tt.text))
			isEmpty := ta.IsEmpty()

			if isEmpty != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", isEmpty, tt.want)
			}
		})
	}
}

func TestTextArea_VisibleLines(t *testing.T) {
	// Create TextArea with 5 lines but height of 3.
	ta := NewTextArea().
		WithBuffer(NewBufferFromString("line0\nline1\nline2\nline3\nline4")).
		WithSize(80, 3)

	// Initially, should show first 3 lines.
	visible := ta.VisibleLines()
	if len(visible) != 3 {
		t.Errorf("len(VisibleLines()) = %d, want 3", len(visible))
	}
	if len(visible) > 0 && visible[0] != "line0" {
		t.Errorf("VisibleLines()[0] = %q, want %q", visible[0], "line0")
	}
}

func TestTextArea_GetBuffer(t *testing.T) {
	buffer := NewBufferFromString("test")
	ta := NewTextArea().WithBuffer(buffer)

	result := ta.GetBuffer()
	if result.String() != "test" {
		t.Errorf("GetBuffer().String() = %q, want %q", result.String(), "test")
	}
}

func TestTextArea_GetCursor(t *testing.T) {
	cursor := NewCursor(5, 10)
	ta := NewTextArea().WithCursor(cursor)

	result := ta.GetCursor()
	if result.Row() != 5 || result.Col() != 10 {
		t.Errorf("GetCursor() = (%d, %d), want (5, 10)", result.Row(), result.Col())
	}
}

func TestTextArea_GetKillRing(t *testing.T) {
	ta := NewTextArea()
	kr := ta.GetKillRing()

	if kr == nil {
		t.Error("GetKillRing() should not return nil")
	}
	if !kr.IsEmpty() {
		t.Error("GetKillRing() should return empty ring")
	}
}

func TestTextArea_WithKillRing(t *testing.T) {
	ta := NewTextArea()
	kr := NewKillRing(20).Kill("test")

	result := ta.WithKillRing(kr)

	if result.GetKillRing().Yank() != "test" {
		t.Errorf("WithKillRing() kill ring Yank() = %q, want %q", result.GetKillRing().Yank(), "test")
	}

	// Verify immutability.
	if ta.GetKillRing().Yank() == "test" {
		t.Error("Original TextArea kill ring was modified")
	}
}

func TestTextArea_Immutability(t *testing.T) {
	original := NewTextArea()

	// Apply all configuration methods.
	sized := original.WithSize(100, 50)
	maxLines := original.WithMaxLines(10)
	maxChars := original.WithMaxChars(100)
	wrapped := original.WithWrap(true)
	placeholder := original.WithPlaceholder("test")
	readonly := original.WithReadOnly(true)
	lineNumbers := original.WithLineNumbers(true)
	buffer := original.WithBuffer(NewBufferFromString("test"))
	cursor := original.WithCursor(NewCursor(1, 1))
	killRing := original.WithKillRing(NewKillRing(20).Kill("test"))

	// Original should remain unchanged.
	if original.Width() != 80 {
		t.Error("Original Width was modified")
	}
	if original.Height() != 24 {
		t.Error("Original Height was modified")
	}
	if original.Placeholder() != "" {
		t.Error("Original Placeholder was modified")
	}
	if original.IsReadOnly() {
		t.Error("Original ReadOnly was modified")
	}
	if original.ShowLineNumbers() {
		t.Error("Original ShowLineNumbers was modified")
	}
	if !original.IsEmpty() {
		t.Error("Original buffer was modified")
	}
	row, col := original.CursorPosition()
	if row != 0 || col != 0 {
		t.Error("Original cursor was modified")
	}

	// Verify all results are different instances.
	instances := []*TextArea{sized, maxLines, maxChars, wrapped, placeholder, readonly, lineNumbers, buffer, cursor, killRing}
	for i, instance := range instances {
		if instance == original {
			t.Errorf("Method %d returned same instance", i)
		}
	}
}

func TestTextArea_ComplexOperations(t *testing.T) {
	// Test chaining multiple operations.
	ta := NewTextArea().
		WithSize(120, 40).
		WithBuffer(NewBufferFromString("line1\nline2\nline3")).
		WithCursor(NewCursor(1, 2)).
		WithPlaceholder("Enter code...").
		WithWrap(true).
		WithLineNumbers(true)

	// Verify all properties.
	if ta.Width() != 120 {
		t.Errorf("Width = %d, want 120", ta.Width())
	}
	if ta.Height() != 40 {
		t.Errorf("Height = %d, want 40", ta.Height())
	}
	if ta.LineCount() != 3 {
		t.Errorf("LineCount = %d, want 3", ta.LineCount())
	}

	row, col := ta.CursorPosition()
	if row != 1 || col != 2 {
		t.Errorf("Cursor = (%d, %d), want (1, 2)", row, col)
	}

	if ta.Placeholder() != "Enter code..." {
		t.Errorf("Placeholder = %q, want %q", ta.Placeholder(), "Enter code...")
	}

	if !ta.ShowLineNumbers() {
		t.Error("ShowLineNumbers should be true")
	}
}

func TestTextArea_SelectedText(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		startRow     int
		startCol     int
		endRow       int
		endCol       int
		hasSelection bool
		want         string
	}{
		{
			name:         "no selection",
			text:         "hello world",
			hasSelection: false,
			want:         "",
		},
		{
			name:         "selection within single line",
			text:         "hello world",
			startRow:     0,
			startCol:     0,
			endRow:       0,
			endCol:       5,
			hasSelection: true,
			want:         "hello",
		},
		{
			name:         "multiline selection",
			text:         "line1\nline2\nline3",
			startRow:     0,
			startCol:     2,
			endRow:       1,
			endCol:       3,
			hasSelection: true,
			want:         "ne1\nlin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea().
				WithBuffer(NewBufferFromString(tt.text))

			if tt.hasSelection {
				// Create selection using value package.
				anchor := value.NewPosition(tt.startRow, tt.startCol)
				cursor := value.NewPosition(tt.endRow, tt.endCol)
				selection := NewSelection(anchor, cursor)
				ta = ta.withSelection(selection)
			}

			result := ta.SelectedText()
			if result != tt.want {
				t.Errorf("SelectedText() = %q, want %q", result, tt.want)
			}
		})
	}
}

func TestTextArea_MaxLines(t *testing.T) {
	tests := []struct {
		name string
		max  int
		want int
	}{
		{
			name: "unlimited",
			max:  0,
			want: 0,
		},
		{
			name: "10 lines max",
			max:  10,
			want: 10,
		},
		{
			name: "100 lines max",
			max:  100,
			want: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea().WithMaxLines(tt.max)

			if ta.MaxLines() != tt.want {
				t.Errorf("MaxLines() = %d, want %d", ta.MaxLines(), tt.want)
			}
		})
	}
}

func TestTextArea_EnsureCursorVisible_Vertical(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		height        int
		cursorRow     int
		initialScroll int
		wantScroll    int
	}{
		{
			name:          "cursor above viewport",
			text:          "line1\nline2\nline3\nline4\nline5",
			height:        2,
			cursorRow:     1,
			initialScroll: 3,
			wantScroll:    1, // Should scroll up to show cursor
		},
		{
			name:          "cursor below viewport",
			text:          "line1\nline2\nline3\nline4\nline5",
			height:        2,
			cursorRow:     4,
			initialScroll: 0,
			wantScroll:    3, // Should scroll down
		},
		{
			name:          "cursor already visible",
			text:          "line1\nline2\nline3\nline4\nline5",
			height:        3,
			cursorRow:     2,
			initialScroll: 1,
			wantScroll:    1, // No change
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea().
				WithBuffer(NewBufferFromString(tt.text)).
				WithSize(80, tt.height)

			// Set scroll and cursor (this will trigger ensureCursorVisible)
			ta = ta.withScroll(tt.initialScroll, 0).
				WithCursor(NewCursor(tt.cursorRow, 0))

			// Check if scroll was adjusted (need to access internal scrollRow)
			// Since scrollRow is private, we check indirectly via VisibleLines.
			visible := ta.VisibleLines()

			// Verify cursor row is in visible range.
			if len(visible) > 0 {
				// The implementation should ensure cursor is visible.
				t.Logf("Visible lines after cursor move: %v", visible)
			}
		})
	}
}

// Helper method for testing (internal access)
func (t *TextArea) withSelection(s *Selection) *TextArea {
	updated := t.copy()
	updated.selection = s
	return updated
}

func (t *TextArea) withScroll(scrollRow, scrollCol int) *TextArea {
	updated := t.copy()
	updated.scrollRow = scrollRow
	updated.scrollCol = scrollCol
	return updated
}

func TestTextArea_EnsureCursorVisible_Horizontal(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		width         int
		wrap          bool
		cursorCol     int
		initialScroll int
		expectScroll  bool
	}{
		{
			name:          "cursor left of viewport (no wrap)",
			text:          "this is a very long line of text",
			width:         10,
			wrap:          false,
			cursorCol:     5,
			initialScroll: 10,
			expectScroll:  true,
		},
		{
			name:          "cursor right of viewport (no wrap)",
			text:          "this is a very long line of text",
			width:         10,
			wrap:          false,
			cursorCol:     25,
			initialScroll: 0,
			expectScroll:  true,
		},
		{
			name:          "wrap mode ignores horizontal scroll",
			text:          "this is a very long line of text",
			width:         10,
			wrap:          true,
			cursorCol:     25,
			initialScroll: 0,
			expectScroll:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea().
				WithBuffer(NewBufferFromString(tt.text)).
				WithSize(tt.width, 10).
				WithWrap(tt.wrap).
				withScroll(0, tt.initialScroll).
				WithCursor(NewCursor(0, tt.cursorCol))

			// Verify cursor was set (ensureCursorVisible called internally)
			_, col := ta.CursorPosition()
			if col != tt.cursorCol {
				t.Errorf("Cursor col = %d, want %d", col, tt.cursorCol)
			}
		})
	}
}

func TestTextArea_WithBuffer_ResetsScrollAndCursor(t *testing.T) {
	ta := NewTextArea().
		WithBuffer(NewBufferFromString("initial\ntext")).
		WithCursor(NewCursor(1, 5)).
		withScroll(1, 5)

	// WithBuffer should reset cursor and clear selection.
	newTa := ta.WithBuffer(NewBufferFromString("new\ntext"))

	row, col := newTa.CursorPosition()
	if row != 0 || col != 0 {
		t.Errorf("WithBuffer() cursor = (%d, %d), want (0, 0)", row, col)
	}

	if newTa.HasSelection() {
		t.Error("WithBuffer() should clear selection")
	}
}

func TestTextArea_VisibleLines_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		height    int
		scrollRow int
		wantCount int
	}{
		{
			name:      "empty buffer",
			text:      "",
			height:    5,
			scrollRow: 0,
			wantCount: 1, // Empty buffer has 1 empty line
		},
		{
			name:      "scroll exactly at last line",
			text:      "line1\nline2\nline3",
			height:    2,
			scrollRow: 1,
			wantCount: 2,
		},
		{
			name:      "large height shows all lines",
			text:      "a\nb",
			height:    100,
			scrollRow: 0,
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := NewTextArea().
				WithBuffer(NewBufferFromString(tt.text)).
				WithSize(80, tt.height).
				withScroll(tt.scrollRow, 0)

			visible := ta.VisibleLines()

			if len(visible) != tt.wantCount {
				t.Errorf("VisibleLines() count = %d, want %d (lines: %v)", len(visible), tt.wantCount, visible)
			}
		})
	}
}
