package model

import (
	"errors"
	"testing"

	"github.com/phoenix-tui/phoenix/components/input/domain/service"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		width     int
		wantWidth int
	}{
		{"normal width", 40, 40},
		{"zero width clamped", 0, 1},
		{"negative width clamped", -10, 1},
		{"large width", 1000, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(tt.width)

			if input.Width() != tt.wantWidth {
				t.Errorf("Width() = %d, want %d", input.Width(), tt.wantWidth)
			}

			// Check initial state
			if input.Content() != "" {
				t.Errorf("Content() = %q, want empty", input.Content())
			}
			if input.CursorPosition() != 0 {
				t.Errorf("CursorPosition() = %d, want 0", input.CursorPosition())
			}
			if input.HasSelection() {
				t.Error("HasSelection() = true, want false")
			}
			if input.Focused() {
				t.Error("Focused() = true, want false")
			}
		})
	}
}

func TestTextInput_WithContent(t *testing.T) {
	tests := []struct {
		name           string
		initial        string
		newContent     string
		initialCursor  int
		expectedCursor int
	}{
		{"empty to text", "", "hello", 0, 0},
		{"text to different", "hello", "world", 3, 3},
		{"text to shorter clamps cursor", "hello world", "hi", 10, 2},
		{"text to empty", "hello", "", 3, 0},
		{"emoji content", "hello", "👋👋👋", 0, 0},
		{"cursor at end stays at end", "hello", "hello world", 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).
				WithContent(tt.initial).
				WithCursor(tt.initialCursor)

			original := input
			result := input.WithContent(tt.newContent)

			// Check immutability
			if original.Content() != tt.initial {
				t.Error("original modified")
			}

			// Check result
			if result.Content() != tt.newContent {
				t.Errorf("Content() = %q, want %q", result.Content(), tt.newContent)
			}
			if result.CursorPosition() != tt.expectedCursor {
				t.Errorf("CursorPosition() = %d, want %d", result.CursorPosition(), tt.expectedCursor)
			}
		})
	}
}

func TestTextInput_SetContent(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		cursorPos      int
		expectedCursor int
	}{
		{"normal position", "hello", 3, 3},
		{"at start", "hello", 0, 0},
		{"at end", "hello", 5, 5},
		{"beyond end clamped", "hello", 10, 5},
		{"negative clamped", "hello", -5, 0},
		{"emoji content", "👋👋👋", 2, 2},
		{"empty content", "", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).SetContent(tt.content, tt.cursorPos)

			if input.Content() != tt.content {
				t.Errorf("Content() = %q, want %q", input.Content(), tt.content)
			}
			if input.CursorPosition() != tt.expectedCursor {
				t.Errorf("CursorPosition() = %d, want %d", input.CursorPosition(), tt.expectedCursor)
			}
			if input.HasSelection() {
				t.Error("SetContent should clear selection")
			}
		})
	}
}

func TestTextInput_ContentParts(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		cursorPos  int
		wantBefore string
		wantAt     string
		wantAfter  string
	}{
		{
			name:       "middle of text",
			content:    "hello",
			cursorPos:  2,
			wantBefore: "he",
			wantAt:     "l",
			wantAfter:  "lo",
		},
		{
			name:       "at start",
			content:    "hello",
			cursorPos:  0,
			wantBefore: "",
			wantAt:     "h",
			wantAfter:  "ello",
		},
		{
			name:       "at end",
			content:    "hello",
			cursorPos:  5,
			wantBefore: "hello",
			wantAt:     "",
			wantAfter:  "",
		},
		{
			name:       "emoji",
			content:    "hello👋world",
			cursorPos:  5,
			wantBefore: "hello",
			wantAt:     "👋",
			wantAfter:  "world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).SetContent(tt.content, tt.cursorPos)
			before, at, after := input.ContentParts()

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

func TestTextInput_MoveLeft(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		initialCursor  int
		expectedCursor int
	}{
		{"from middle", "hello", 3, 2},
		{"from start stays", "hello", 0, 0},
		{"from end", "hello", 5, 4},
		{"emoji", "hello👋world", 6, 5},
		{"single char", "a", 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).
				SetContent(tt.content, tt.initialCursor).
				MoveLeft()

			if input.CursorPosition() != tt.expectedCursor {
				t.Errorf("CursorPosition() = %d, want %d", input.CursorPosition(), tt.expectedCursor)
			}
			if input.HasSelection() {
				t.Error("MoveLeft should clear selection")
			}
		})
	}
}

func TestTextInput_MoveRight(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		initialCursor  int
		expectedCursor int
	}{
		{"from start", "hello", 0, 1},
		{"from middle", "hello", 2, 3},
		{"from before end", "hello", 4, 5},
		{"from end stays", "hello", 5, 5},
		{"emoji", "hello👋world", 5, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).
				SetContent(tt.content, tt.initialCursor).
				MoveRight()

			if input.CursorPosition() != tt.expectedCursor {
				t.Errorf("CursorPosition() = %d, want %d", input.CursorPosition(), tt.expectedCursor)
			}
			if input.HasSelection() {
				t.Error("MoveRight should clear selection")
			}
		})
	}
}

func TestTextInput_MoveHome(t *testing.T) {
	input := New(40).
		SetContent("hello world", 7).
		MoveHome()

	if input.CursorPosition() != 0 {
		t.Errorf("CursorPosition() = %d, want 0", input.CursorPosition())
	}
	if input.HasSelection() {
		t.Error("MoveHome should clear selection")
	}
}

func TestTextInput_MoveEnd(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantPos int
	}{
		{"ascii", "hello", 5},
		{"emoji", "hello👋", 6},
		{"empty", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).
				SetContent(tt.content, 0).
				MoveEnd()

			if input.CursorPosition() != tt.wantPos {
				t.Errorf("CursorPosition() = %d, want %d", input.CursorPosition(), tt.wantPos)
			}
		})
	}
}

func TestTextInput_InsertRune(t *testing.T) {
	tests := []struct {
		name        string
		initial     string
		cursorPos   int
		insertRune  rune
		wantContent string
		wantCursor  int
	}{
		{"insert at start", "hello", 0, 'X', "Xhello", 1},
		{"insert in middle", "hello", 2, 'X', "heXllo", 3},
		{"insert at end", "hello", 5, 'X', "helloX", 6},
		{"insert emoji", "hello", 5, '👋', "hello👋", 6},
		{"insert into empty", "", 0, 'X', "X", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).
				SetContent(tt.initial, tt.cursorPos).
				InsertRune(tt.insertRune)

			if input.Content() != tt.wantContent {
				t.Errorf("Content() = %q, want %q", input.Content(), tt.wantContent)
			}
			if input.CursorPosition() != tt.wantCursor {
				t.Errorf("CursorPosition() = %d, want %d", input.CursorPosition(), tt.wantCursor)
			}
		})
	}
}

func TestTextInput_InsertRune_WithSelection(t *testing.T) {
	input := New(40).
		SetContent("hello world", 0).
		WithSelection(0, 5). // Select "hello"
		InsertRune('X')

	if input.Content() != "X world" {
		t.Errorf("Content() = %q, want %q", input.Content(), "X world")
	}
	if input.CursorPosition() != 1 {
		t.Errorf("CursorPosition() = %d, want 1", input.CursorPosition())
	}
	if input.HasSelection() {
		t.Error("selection should be cleared")
	}
}

func TestTextInput_DeleteBackward(t *testing.T) {
	tests := []struct {
		name        string
		initial     string
		cursorPos   int
		wantContent string
		wantCursor  int
	}{
		{"from middle", "hello", 3, "helo", 2},
		{"from end", "hello", 5, "hell", 4},
		{"from start no-op", "hello", 0, "hello", 0},
		{"emoji", "hello👋world", 6, "helloworld", 5},
		{"single char", "a", 1, "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).
				SetContent(tt.initial, tt.cursorPos).
				DeleteBackward()

			if input.Content() != tt.wantContent {
				t.Errorf("Content() = %q, want %q", input.Content(), tt.wantContent)
			}
			if input.CursorPosition() != tt.wantCursor {
				t.Errorf("CursorPosition() = %d, want %d", input.CursorPosition(), tt.wantCursor)
			}
		})
	}
}

func TestTextInput_DeleteForward(t *testing.T) {
	tests := []struct {
		name        string
		initial     string
		cursorPos   int
		wantContent string
		wantCursor  int
	}{
		{"from start", "hello", 0, "ello", 0},
		{"from middle", "hello", 2, "helo", 2},
		{"from end no-op", "hello", 5, "hello", 5},
		{"emoji", "hello👋world", 5, "helloworld", 5},
		{"single char", "a", 0, "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).
				SetContent(tt.initial, tt.cursorPos).
				DeleteForward()

			if input.Content() != tt.wantContent {
				t.Errorf("Content() = %q, want %q", input.Content(), tt.wantContent)
			}
			if input.CursorPosition() != tt.wantCursor {
				t.Errorf("CursorPosition() = %d, want %d", input.CursorPosition(), tt.wantCursor)
			}
		})
	}
}

func TestTextInput_Clear(t *testing.T) {
	input := New(40).
		SetContent("hello world", 7).
		WithSelection(0, 5).
		Clear()

	if input.Content() != "" {
		t.Errorf("Content() = %q, want empty", input.Content())
	}
	if input.CursorPosition() != 0 {
		t.Errorf("CursorPosition() = %d, want 0", input.CursorPosition())
	}
	if input.HasSelection() {
		t.Error("selection should be cleared")
	}
	if input.ScrollOffset() != 0 {
		t.Errorf("ScrollOffset() = %d, want 0", input.ScrollOffset())
	}
}

func TestTextInput_SelectAll(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantEnd int
	}{
		{"ascii", "hello", 5},
		{"emoji", "hello👋", 6},
		{"empty", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).
				SetContent(tt.content, 0).
				SelectAll()

			// Empty content should not have selection (0,0 is empty)
			if tt.content == "" {
				if input.HasSelection() {
					t.Error("HasSelection() = true for empty content, want false")
				}
			} else {
				if !input.HasSelection() {
					t.Error("HasSelection() = false, want true")
				}

				sel := input.Selection()
				if sel.Start() != 0 {
					t.Errorf("Selection.Start() = %d, want 0", sel.Start())
				}
				if sel.End() != tt.wantEnd {
					t.Errorf("Selection.End() = %d, want %d", sel.End(), tt.wantEnd)
				}
			}

			if input.CursorPosition() != tt.wantEnd {
				t.Errorf("CursorPosition() = %d, want %d", input.CursorPosition(), tt.wantEnd)
			}
		})
	}
}

func TestTextInput_WithSelection(t *testing.T) {
	input := New(40).
		SetContent("hello world", 0).
		WithSelection(0, 5)

	if !input.HasSelection() {
		t.Error("HasSelection() = false, want true")
	}

	sel := input.Selection()
	if sel.Start() != 0 || sel.End() != 5 {
		t.Errorf("Selection = (%d, %d), want (0, 5)", sel.Start(), sel.End())
	}

	// Cursor should be at end of selection
	if input.CursorPosition() != 5 {
		t.Errorf("CursorPosition() = %d, want 5", input.CursorPosition())
	}
}

func TestTextInput_ClearSelection(t *testing.T) {
	input := New(40).
		SetContent("hello world", 0).
		WithSelection(0, 5).
		ClearSelection()

	if input.HasSelection() {
		t.Error("HasSelection() = true, want false")
	}
}

func TestTextInput_WithValidator(t *testing.T) {
	validator := service.NotEmpty()
	input := New(40).WithValidator(validator)

	// Empty should be invalid
	if input.IsValid() {
		t.Error("IsValid() = true, want false for empty content")
	}

	// Non-empty should be valid
	input = input.WithContent("hello")
	if !input.IsValid() {
		t.Error("IsValid() = false, want true for non-empty content")
	}
}

func TestTextInput_Validate(t *testing.T) {
	validator := func(s string) error {
		if s != "valid" {
			return errors.New("must be 'valid'")
		}
		return nil
	}

	input := New(40).WithValidator(validator)

	// Invalid content
	input = input.WithContent("invalid")
	if err := input.Validate(); err == nil {
		t.Error("Validate() = nil, want error")
	}

	// Valid content
	input = input.WithContent("valid")
	if err := input.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}

	// No validator
	input = New(40).WithContent("anything")
	if err := input.Validate(); err != nil {
		t.Errorf("Validate() with no validator = %v, want nil", err)
	}
}

func TestTextInput_WithFocus(t *testing.T) {
	input := New(40)

	if input.Focused() {
		t.Error("initial Focused() = true, want false")
	}

	focused := (*input).WithFocus(true)
	if !focused.Focused() {
		t.Error("Focused() = false after WithFocus(true)")
	}

	unfocused := (*input).WithFocus(false)
	if unfocused.Focused() {
		t.Error("Focused() = true after WithFocus(false)")
	}
}

func TestTextInput_WithPlaceholder(t *testing.T) {
	placeholder := "Enter text..."
	input := New(40).WithPlaceholder(placeholder)

	if input.Placeholder() != placeholder {
		t.Errorf("Placeholder() = %q, want %q", input.Placeholder(), placeholder)
	}
}

func TestTextInput_WithWidth(t *testing.T) {
	input := New(40).WithWidth(80)

	if input.Width() != 80 {
		t.Errorf("Width() = %d, want 80", input.Width())
	}

	// Zero/negative should be clamped
	input = input.WithWidth(0)
	if input.Width() != 1 {
		t.Errorf("Width() = %d, want 1 (clamped)", input.Width())
	}
}

func TestTextInput_Immutability(t *testing.T) {
	original := New(40).
		SetContent("hello", 3).
		WithPlaceholder("test").
		WithFocus(true)

	// Perform various operations
	_ = original.WithContent("world")
	_ = original.MoveLeft()
	_ = original.InsertRune('x')
	_ = original.DeleteBackward()
	_ = original.Clear()
	_ = original.SelectAll()

	// Original should be unchanged
	if original.Content() != "hello" {
		t.Error("original content modified")
	}
	if original.CursorPosition() != 3 {
		t.Error("original cursor modified")
	}
	if original.Placeholder() != "test" {
		t.Error("original placeholder modified")
	}
	if !original.Focused() {
		t.Error("original focus modified")
	}
}

func TestTextInput_DeleteSelection(t *testing.T) {
	input := New(40).
		SetContent("hello world", 0).
		WithSelection(0, 5). // Select "hello"
		DeleteBackward()     // Should delete selection

	if input.Content() != " world" {
		t.Errorf("Content() = %q, want %q", input.Content(), " world")
	}
	if input.CursorPosition() != 0 {
		t.Errorf("CursorPosition() = %d, want 0", input.CursorPosition())
	}
	if input.HasSelection() {
		t.Error("selection should be cleared")
	}
}

func TestTextInput_ComplexUnicode(t *testing.T) {
	// Test with various Unicode complexities
	tests := []struct {
		name    string
		content string
	}{
		{"family emoji", "👨‍👩‍👧‍👦"},
		{"flag emoji", "🇺🇸"},
		{"skin tone emoji", "👋🏽"},
		{"combining chars", "é"},
		{"mixed", "Hello 世界 👋"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).SetContent(tt.content, 0)

			// Calculate expected max position (grapheme count)
			svc := service.NewCursorMovementService()
			maxPos := svc.GraphemeCount(tt.content)

			// Should be able to navigate through content
			count := 0
			for input.CursorPosition() < maxPos {
				input = input.MoveRight()
				count++
				if count > 100 {
					t.Fatal("infinite loop detected")
				}
			}

			// Should be at end
			if input.CursorPosition() != maxPos {
				t.Errorf("final cursor = %d, want %d", input.CursorPosition(), maxPos)
			}
		})
	}
}
