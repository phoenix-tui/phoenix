package infrastructure

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/input/internal/input/domain/model"
	"github.com/phoenix-tui/phoenix/tea"
)

func TestDefaultKeyBindings_Navigation(t *testing.T) {
	kb := NewDefaultKeyBindings()
	input := model.New(40).SetContent("hello world", 6)

	tests := []struct {
		name           string
		key            tea.KeyMsg
		expectedCursor int
	}{
		{"left arrow", tea.KeyMsg{Type: tea.KeyLeft}, 5},
		{"right arrow", tea.KeyMsg{Type: tea.KeyRight}, 7},
		{"home", tea.KeyMsg{Type: tea.KeyHome}, 0},
		{"end", tea.KeyMsg{Type: tea.KeyEnd}, 11},
		{"ctrl-e (end)", tea.KeyMsg{Ctrl: true, Rune: 'e'}, 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kb.Handle(input, tt.key)
			if result.CursorPosition() != tt.expectedCursor {
				t.Errorf("CursorPosition() = %d, want %d", result.CursorPosition(), tt.expectedCursor)
			}
		})
	}
}

func TestDefaultKeyBindings_Editing(t *testing.T) {
	kb := NewDefaultKeyBindings()

	tests := []struct {
		name        string
		initial     string
		cursorPos   int
		key         tea.KeyMsg
		wantContent string
		wantCursor  int
	}{
		{
			name:        "backspace",
			initial:     "hello",
			cursorPos:   5,
			key:         tea.KeyMsg{Type: tea.KeyBackspace},
			wantContent: "hell",
			wantCursor:  4,
		},
		{
			name:        "delete",
			initial:     "hello",
			cursorPos:   0,
			key:         tea.KeyMsg{Type: tea.KeyDelete},
			wantContent: "ello",
			wantCursor:  0,
		},
		{
			name:        "ctrl-u clear",
			initial:     "hello",
			cursorPos:   3,
			key:         tea.KeyMsg{Ctrl: true, Rune: 'u'},
			wantContent: "",
			wantCursor:  0,
		},
		{
			name:        "insert rune",
			initial:     "hello",
			cursorPos:   5,
			key:         tea.KeyMsg{Type: tea.KeyRune, Rune: '!'},
			wantContent: "hello!",
			wantCursor:  6,
		},
		{
			name:        "insert space (KeySpace)",
			initial:     "hello",
			cursorPos:   5,
			key:         tea.KeyMsg{Type: tea.KeySpace}, // Space parsed as KeySpace, not KeyRune
			wantContent: "hello ",
			wantCursor:  6,
		},
		{
			name:        "insert space in middle",
			initial:     "helloworld",
			cursorPos:   5,
			key:         tea.KeyMsg{Type: tea.KeySpace},
			wantContent: "hello world",
			wantCursor:  6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := model.New(40).SetContent(tt.initial, tt.cursorPos)
			result := kb.Handle(input, tt.key)

			if result.Content() != tt.wantContent {
				t.Errorf("Content() = %q, want %q", result.Content(), tt.wantContent)
			}
			if result.CursorPosition() != tt.wantCursor {
				t.Errorf("CursorPosition() = %d, want %d", result.CursorPosition(), tt.wantCursor)
			}
		})
	}
}

func TestDefaultKeyBindings_InsertMultipleRunes(t *testing.T) {
	kb := NewDefaultKeyBindings()
	input := model.New(40).SetContent("", 0)

	// Insert "hello" one rune at a time.
	runes := []rune{'h', 'e', 'l', 'l', 'o'}
	result := input
	for _, r := range runes {
		msg := tea.KeyMsg{Type: tea.KeyRune, Rune: r}
		updated := kb.Handle(result, msg)
		result = updated
	}

	if result.Content() != "hello" {
		t.Errorf("Content() = %q, want %q", result.Content(), "hello")
	}
	if result.CursorPosition() != 5 {
		t.Errorf("CursorPosition() = %d, want 5", result.CursorPosition())
	}
}

func TestDefaultKeyBindings_SelectAll(t *testing.T) {
	kb := NewDefaultKeyBindings()
	input := model.New(40).SetContent("hello world", 5)

	// Ctrl-A.
	msg := tea.KeyMsg{Ctrl: true, Rune: 'a'}

	result := kb.Handle(input, msg)

	if !result.HasSelection() {
		t.Error("SelectAll should create selection")
	}

	sel := result.Selection()
	if sel.Start() != 0 || sel.End() != 11 {
		t.Errorf("Selection = (%d, %d), want (0, 11)", sel.Start(), sel.End())
	}
}

func TestDefaultKeyBindings_UnhandledKey(t *testing.T) {
	kb := NewDefaultKeyBindings()
	input := model.New(40).SetContent("hello", 5)

	// Unhandled key should return input unchanged.
	msg := tea.KeyMsg{Type: tea.KeyF1}
	result := kb.Handle(input, msg)

	// Check that input is unchanged.
	if result.Content() != input.Content() {
		t.Errorf("Content changed for unhandled key: got %q, want %q", result.Content(), input.Content())
	}
	if result.CursorPosition() != input.CursorPosition() {
		t.Errorf("Cursor changed for unhandled key: got %d, want %d", result.CursorPosition(), input.CursorPosition())
	}
}

func TestIsNavigationKey(t *testing.T) {
	tests := []struct {
		name string
		key  tea.KeyMsg
		want bool
	}{
		{"left arrow", tea.KeyMsg{Type: tea.KeyLeft}, true},
		{"right arrow", tea.KeyMsg{Type: tea.KeyRight}, true},
		{"home", tea.KeyMsg{Type: tea.KeyHome}, true},
		{"end", tea.KeyMsg{Type: tea.KeyEnd}, true},
		{"ctrl-a", tea.KeyMsg{Ctrl: true, Rune: 'a'}, true},
		{"backspace", tea.KeyMsg{Type: tea.KeyBackspace}, false},
		{"rune", tea.KeyMsg{Type: tea.KeyRune, Rune: 'x'}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNavigationKey(tt.key)
			if got != tt.want {
				t.Errorf("IsNavigationKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEditingKey(t *testing.T) {
	tests := []struct {
		name string
		key  tea.KeyMsg
		want bool
	}{
		{"backspace", tea.KeyMsg{Type: tea.KeyBackspace}, true},
		{"delete", tea.KeyMsg{Type: tea.KeyDelete}, true},
		{"space", tea.KeyMsg{Type: tea.KeySpace}, true}, // Space is editing key
		{"ctrl-u", tea.KeyMsg{Ctrl: true, Rune: 'u'}, true},
		{"rune", tea.KeyMsg{Type: tea.KeyRune, Rune: 'x'}, true},
		{"left arrow", tea.KeyMsg{Type: tea.KeyLeft}, false},
		{"home", tea.KeyMsg{Type: tea.KeyHome}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsEditingKey(tt.key)
			if got != tt.want {
				t.Errorf("IsEditingKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomKeyBindings(t *testing.T) {
	// Create a custom handler that intercepts Ctrl-D to duplicate line.
	customHandler := func(input model.TextInput, msg tea.KeyMsg) model.TextInput {
		if msg.Ctrl && (msg.Rune == 'd' || msg.Rune == 'D') {
			content := input.Content()
			return input.WithContent(content + content)
		}
		return input // Not handled - return unchanged
	}

	kb := NewCustomKeyBindings(customHandler)
	input := model.New(40).SetContent("hello", 5)

	// Test custom handler - Ctrl-D.
	msg := tea.KeyMsg{Ctrl: true, Rune: 'd'}
	result := kb.Handle(input, msg)

	if result.Content() != "hellohello" {
		t.Errorf("Content() = %q, want %q", result.Content(), "hellohello")
	}

	// Test fallback to default.
	msg2 := tea.KeyMsg{Type: tea.KeyLeft}
	result2 := kb.Handle(input, msg2)

	if result2.CursorPosition() != 4 {
		t.Errorf("CursorPosition() = %d, want 4", result2.CursorPosition())
	}
}

func TestCustomKeyBindings_AddHandler(t *testing.T) {
	kb := NewCustomKeyBindings()

	// Add handler dynamically - Ctrl-X.
	handler := func(input model.TextInput, msg tea.KeyMsg) model.TextInput {
		if msg.Ctrl && (msg.Rune == 'x' || msg.Rune == 'X') {
			return input.Clear()
		}
		return input
	}
	kb.AddHandler(handler)

	input := model.New(40).SetContent("hello", 5)
	msg := tea.KeyMsg{Ctrl: true, Rune: 'x'}
	result := kb.Handle(input, msg)

	if result.Content() != "" {
		t.Errorf("Content() = %q, want empty", result.Content())
	}
}

func TestCustomKeyBindings_MultipleHandlers(t *testing.T) {
	// First handler handles ctrl+1.
	handler1 := func(input model.TextInput, msg tea.KeyMsg) model.TextInput {
		if msg.Ctrl && msg.Rune == '1' {
			return input.WithContent("one")
		}
		return input
	}

	// Second handler handles ctrl+2.
	handler2 := func(input model.TextInput, msg tea.KeyMsg) model.TextInput {
		if msg.Ctrl && msg.Rune == '2' {
			return input.WithContent("two")
		}
		return input
	}

	kb := NewCustomKeyBindings(handler1, handler2)
	input := *model.New(40) // Dereference pointer

	// Test first handler - Ctrl-1.
	msg1 := tea.KeyMsg{Ctrl: true, Rune: '1'}
	result1 := kb.Handle(input, msg1)
	if result1.Content() != "one" {
		t.Error("first handler failed")
	}

	// Test second handler - Ctrl-2.
	msg2 := tea.KeyMsg{Ctrl: true, Rune: '2'}
	result2 := kb.Handle(input, msg2)
	if result2.Content() != "two" {
		t.Error("second handler failed")
	}
}
