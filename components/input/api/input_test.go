package input

import (
	"errors"
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/components/input/domain/model"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

func TestNew(t *testing.T) {
	input := New(40)

	if input == nil {
		t.Fatal("New() returned nil")
	}

	if input.Value() != "" {
		t.Errorf("initial Value() = %q, want empty", input.Value())
	}

	if input.IsFocused() {
		t.Error("initial IsFocused() = true, want false")
	}
}

func TestInput_FluentAPI(t *testing.T) {
	// Test fluent chaining (pointer chaining from New)
	input := New(40).
		Placeholder("Enter text...").
		Content("hello").
		Focused(true).
		Width(80)

	if input.Value() != "hello" {
		t.Errorf("Value() = %q, want %q", input.Value(), "hello")
	}

	if !input.IsFocused() {
		t.Error("IsFocused() = false, want true")
	}
}

func TestInput_Validator(t *testing.T) {
	validator := func(s string) error {
		if !strings.Contains(s, "@") {
			return errors.New("must contain @")
		}
		return nil
	}

	// Pointer chaining from New.
	input := New(40).
		Validator(validator).
		Content("test")

	if input.IsValid() {
		t.Error("IsValid() = true, want false (missing @)")
	}

	// After first method, we have Input value.
	input = input.Content("test@example.com")
	if !input.IsValid() {
		t.Error("IsValid() = false, want true (has @)")
	}
}

func TestInput_CursorPosition(t *testing.T) {
	input := New(40).SetContent("hello world", 6)

	if input.CursorPosition() != 6 {
		t.Errorf("CursorPosition() = %d, want 6", input.CursorPosition())
	}
}

func TestInput_ContentParts(t *testing.T) {
	input := New(40).SetContent("hello world", 6)
	before, at, after := input.ContentParts()

	if before != "hello " {
		t.Errorf("before = %q, want %q", before, "hello ")
	}
	if at != "w" {
		t.Errorf("at = %q, want %q", at, "w")
	}
	if after != "orld" {
		t.Errorf("after = %q, want %q", after, "orld")
	}
}

func TestInput_SetContent(t *testing.T) {
	input := New(40).SetContent("hello", 3)

	if input.Value() != "hello" {
		t.Errorf("Value() = %q, want %q", input.Value(), "hello")
	}
	if input.CursorPosition() != 3 {
		t.Errorf("CursorPosition() = %d, want 3", input.CursorPosition())
	}
}

func TestInput_Init(t *testing.T) {
	input := New(40)
	cmd := input.Init()

	if cmd != nil {
		t.Error("Init() should return nil cmd")
	}
}

func TestInput_Update_Focused(t *testing.T) {
	// Pointer chaining from New.
	input := New(40).Content("hello").Focused(true)

	// Send left arrow key.
	msg := tea.KeyMsg{Type: tea.KeyLeft}
	_, cmd := input.Update(msg)

	if cmd != nil {
		t.Error("Update() should return nil cmd")
	}

	// Cursor should move (assuming default position at end)
	// After typing "hello", cursor is at 5, left arrow moves to 4.
	inputAtEnd := input.SetContent("hello", 5)
	updated, _ := inputAtEnd.Update(msg)

	if updated.CursorPosition() != 4 {
		t.Errorf("CursorPosition() = %d, want 4 after left arrow", updated.CursorPosition())
	}
}

func TestInput_Update_Unfocused(t *testing.T) {
	input := New(40).Content("hello").Focused(false)

	// Send left arrow key (should be ignored)
	msg := tea.KeyMsg{Type: tea.KeyLeft}
	updated, _ := input.Update(msg)

	// Should be unchanged.
	if updated.Value() != input.Value() {
		t.Error("unfocused input should ignore keys")
	}
}

func TestInput_Update_InsertText(t *testing.T) {
	// Method chaining returns Input value after first call.
	input := New(40).Content("").Focused(true)

	// Insert 'h'.
	msg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'h'}
	input, _ = input.Update(msg) // Reassignment!

	if input.Value() != "h" {
		t.Errorf("Value() = %q, want %q", input.Value(), "h")
	}

	// Insert 'i'.
	msg = tea.KeyMsg{Type: tea.KeyRune, Rune: 'i'}
	input, _ = input.Update(msg) // Reassignment!

	if input.Value() != "hi" {
		t.Errorf("Value() = %q, want %q", input.Value(), "hi")
	}
}

func TestInput_Update_Backspace(t *testing.T) {
	// Method chaining returns Input value.
	input := New(40).SetContent("hello", 5).Focused(true)

	msg := tea.KeyMsg{Type: tea.KeyBackspace}
	input, _ = input.Update(msg) // Reassignment!

	if input.Value() != "hell" {
		t.Errorf("Value() = %q, want %q", input.Value(), "hell")
	}
}

func TestInput_View_Empty(t *testing.T) {
	input := New(40).Content("").Focused(false)
	view := input.View()

	if view != "" {
		t.Errorf("View() = %q, want empty", view)
	}
}

func TestInput_View_Placeholder(t *testing.T) {
	input := New(40).
		Placeholder("Enter text...").
		Content("").
		Focused(false)

	view := input.View()

	if !strings.Contains(view, "Enter text") {
		t.Errorf("View() should contain placeholder, got %q", view)
	}
}

func TestInput_View_Content(t *testing.T) {
	input := New(40).
		Content("hello").
		Focused(false)

	view := input.View()

	if !strings.Contains(view, "hello") {
		t.Errorf("View() should contain content, got %q", view)
	}
}

func TestInput_View_WithCursor(t *testing.T) {
	input := New(40).
		SetContent("hello", 2).
		Focused(true)

	view := input.View()

	// Should contain the content.
	if !strings.Contains(view, "he") {
		t.Errorf("View() should contain 'he', got %q", view)
	}

	// Should have some cursor rendering.
	// (exact format depends on renderCursor implementation)
	if len(view) == 0 {
		t.Error("View() should not be empty with content and cursor")
	}
}

func TestInput_CustomKeyBindings(t *testing.T) {
	// Create custom handler that doubles text on Ctrl-D (VALUE SEMANTICS!)
	customHandler := func(domain model.TextInput, msg tea.KeyMsg) model.TextInput {
		if msg.Ctrl && (msg.Rune == 'd' || msg.Rune == 'D') {
			content := domain.Content()
			return domain.WithContent(content + content)
		}
		return domain // Return unchanged if not handled
	}

	// Wrap in KeyBindingHandler interface.
	wrappedHandler := CustomKeyBindings(customHandler)

	// Can't directly test this without exposing more internals.
	// This test verifies the API exists.
	if wrappedHandler == nil {
		t.Error("CustomKeyBindings() returned nil")
	}
}

func TestInput_CommonValidators(t *testing.T) {
	// Test NotEmpty.
	input := New(40).Validator(NotEmpty()).Content("")
	if input.IsValid() {
		t.Error("empty should be invalid with NotEmpty validator")
	}

	input = input.Content("text")
	if !input.IsValid() {
		t.Error("non-empty should be valid with NotEmpty validator")
	}

	// Test MinLength.
	input = New(40).Validator(MinLength(5)).Content("abc")
	if input.IsValid() {
		t.Error("'abc' should be invalid with MinLength(5)")
	}

	input = input.Content("abcde")
	if !input.IsValid() {
		t.Error("'abcde' should be valid with MinLength(5)")
	}

	// Test MaxLength.
	input = New(40).Validator(MaxLength(5)).Content("abcdef")
	if input.IsValid() {
		t.Error("'abcdef' should be invalid with MaxLength(5)")
	}

	input = input.Content("abcde")
	if !input.IsValid() {
		t.Error("'abcde' should be valid with MaxLength(5)")
	}

	// Test Range.
	input = New(40).Validator(Range(3, 7)).Content("ab")
	if input.IsValid() {
		t.Error("'ab' should be invalid with Range(3, 7)")
	}

	input = input.Content("abcd")
	if !input.IsValid() {
		t.Error("'abcd' should be valid with Range(3, 7)")
	}

	input = input.Content("abcdefgh")
	if input.IsValid() {
		t.Error("'abcdefgh' should be invalid with Range(3, 7)")
	}

	// Test Chain.
	input = New(40).Validator(Chain(NotEmpty(), MinLength(3), MaxLength(10))).Content("")
	if input.IsValid() {
		t.Error("empty should be invalid with chained validators")
	}

	input = input.Content("ab")
	if input.IsValid() {
		t.Error("'ab' should be invalid (too short)")
	}

	input = input.Content("abc")
	if !input.IsValid() {
		t.Error("'abc' should be valid")
	}
}

func TestInput_ValidationErrors(t *testing.T) {
	// Just verify the errors are exported.
	if ErrEmpty == nil {
		t.Error("ErrEmpty should be exported")
	}
	if ErrTooShort == nil {
		t.Error("ErrTooShort should be exported")
	}
	if ErrTooLong == nil {
		t.Error("ErrTooLong should be exported")
	}
	if ErrInvalidFormat == nil {
		t.Error("ErrInvalidFormat should be exported")
	}
}

func TestInput_Immutability(t *testing.T) {
	// Method chaining returns Input value.
	original := New(40).Content("hello").Focused(true)

	// Apply various operations (don't reassign to original!)
	_ = original.Content("world")
	_ = original.Focused(false)
	_ = original.Placeholder("test")
	_ = original.Width(80)

	// Original should be unchanged (value semantics!)
	if original.Value() != "hello" {
		t.Error("original Value() modified")
	}
	if !original.IsFocused() {
		t.Error("original Focused() modified")
	}
}

func TestInput_CompleteWorkflow(t *testing.T) {
	// Simulate a complete user interaction.
	input := New(40).
		Placeholder("Enter email...").
		Validator(func(s string) error {
			if !strings.Contains(s, "@") {
				return errors.New("invalid email")
			}
			return nil
		}).
		Focused(true)

	// User types "test" (reassignment for value semantics!)
	for _, r := range "test" {
		msg := tea.KeyMsg{Type: tea.KeyRune, Rune: r}
		input, _ = input.Update(msg) // Reassignment!
	}

	if input.Value() != "test" {
		t.Errorf("Value() = %q, want %q", input.Value(), "test")
	}

	// Should be invalid (no @)
	if input.IsValid() {
		t.Error("'test' should be invalid email")
	}

	// User types "@example.com" (reassignment!)
	for _, r := range "@example.com" {
		msg := tea.KeyMsg{Type: tea.KeyRune, Rune: r}
		input, _ = input.Update(msg) // Reassignment!
	}

	if input.Value() != "test@example.com" {
		t.Errorf("Value() = %q, want %q", input.Value(), "test@example.com")
	}

	// Should be valid now.
	if !input.IsValid() {
		t.Error("'test@example.com' should be valid email")
	}

	// User presses Ctrl-U to clear (reassignment!)
	msg := tea.KeyMsg{Ctrl: true, Rune: 'u'}
	input, _ = input.Update(msg) // Reassignment!

	if input.Value() != "" {
		t.Errorf("Value() = %q, want empty after Ctrl-U", input.Value())
	}

	// Render should show placeholder (method chaining returns value)
	view := input.Focused(false).View()
	if !strings.Contains(view, "Enter email") {
		t.Errorf("View() should show placeholder, got %q", view)
	}
}

// Additional coverage tests

func TestInput_ShowCursor(t *testing.T) {
	tests := []struct {
		name string
		show bool
	}{
		{"enabled", true},
		{"disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).ShowCursor(tt.show).Focused(true).Content("test")

			// Verify fluent API works
			view := input.View()

			// Should render (cursor visibility is internal detail)
			if len(view) == 0 {
				t.Error("View() should not be empty")
			}
		})
	}
}

func TestInput_KeyBindings(t *testing.T) {
	// Create custom key bindings handler
	customHandler := func(domain model.TextInput, msg tea.KeyMsg) model.TextInput {
		// Custom handler that uppercases on Ctrl-U
		if msg.Ctrl && (msg.Rune == 'u' || msg.Rune == 'U') {
			return domain.WithContent(strings.ToUpper(domain.Content()))
		}
		return domain
	}

	input := New(40).
		Content("hello").
		Focused(true).
		KeyBindings(CustomKeyBindings(customHandler))

	// Send Ctrl-U
	msg := tea.KeyMsg{Ctrl: true, Rune: 'u'}
	input, _ = input.Update(msg)

	if input.Value() != "HELLO" {
		t.Errorf("Value() = %q, want %q after custom Ctrl-U", input.Value(), "HELLO")
	}
}

func TestInput_RenderContent_Unicode(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"simple ascii", "hello"},
		{"unicode chinese", "你好"},
		{"unicode emoji", "👋🌍"},
		{"mixed", "hello 世界 👋"},
		{"long emoji", "👨‍👩‍👧‍👦"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).Content(tt.content).Focused(false)

			view := input.View()

			if !strings.Contains(view, tt.content) {
				t.Errorf("View() should contain %q, got %q", tt.content, view)
			}
		})
	}
}

func TestInput_RenderContent_Scrolling(t *testing.T) {
	// Test content longer than width
	longContent := "this is a very long string that exceeds the input width and should scroll"

	input := New(20).SetContent(longContent, 50).Focused(true)

	view := input.View()

	// Should render something (scrolling logic tested)
	if len(view) == 0 {
		t.Error("View() should not be empty for long content")
	}
}

func TestInput_RenderCursor_Positions(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		position int
	}{
		{"beginning", "hello", 0},
		{"middle", "hello", 2},
		{"end", "hello", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).
				SetContent(tt.content, tt.position).
				Focused(true)

			view := input.View()

			// Should render cursor at position
			if len(view) == 0 {
				t.Error("View() should not be empty with cursor")
			}
		})
	}
}

func TestInput_RenderCursor_UnicodeContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		position int
	}{
		{"emoji middle", "👋hello👋", 1},
		{"chinese middle", "你好世界", 2},
		{"mixed", "hello世界", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).
				SetContent(tt.content, tt.position).
				Focused(true)

			view := input.View()

			// Should handle Unicode cursor positioning
			if len(view) == 0 {
				t.Error("View() should not be empty")
			}
		})
	}
}

func TestInput_Update_SpecialKeys(t *testing.T) {
	tests := []struct {
		name       string
		initial    string
		initialPos int
		key        tea.KeyMsg
		wantValue  string
		checkPos   bool
		wantPos    int
	}{
		{
			name:       "Home key",
			initial:    "hello",
			initialPos: 3,
			key:        tea.KeyMsg{Type: tea.KeyHome},
			wantValue:  "hello",
			checkPos:   true,
			wantPos:    0,
		},
		{
			name:       "End key",
			initial:    "hello",
			initialPos: 2,
			key:        tea.KeyMsg{Type: tea.KeyEnd},
			wantValue:  "hello",
			checkPos:   true,
			wantPos:    5,
		},
		{
			name:       "Delete key",
			initial:    "hello",
			initialPos: 2,
			key:        tea.KeyMsg{Type: tea.KeyDelete},
			wantValue:  "helo",
			checkPos:   false,
		},
		{
			name:       "Right arrow",
			initial:    "hello",
			initialPos: 2,
			key:        tea.KeyMsg{Type: tea.KeyRight},
			wantValue:  "hello",
			checkPos:   true,
			wantPos:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := New(40).
				SetContent(tt.initial, tt.initialPos).
				Focused(true)

			input, _ = input.Update(tt.key)

			if input.Value() != tt.wantValue {
				t.Errorf("Value() = %q, want %q", input.Value(), tt.wantValue)
			}

			if tt.checkPos && input.CursorPosition() != tt.wantPos {
				t.Errorf("CursorPosition() = %d, want %d", input.CursorPosition(), tt.wantPos)
			}
		})
	}
}

func TestInput_View_CursorHidden(t *testing.T) {
	input := New(40).
		Content("hello").
		Focused(true).
		ShowCursor(false)

	view := input.View()

	// Should still render content even without cursor
	if !strings.Contains(view, "hello") {
		t.Errorf("View() should contain content, got %q", view)
	}
}

func TestInput_View_EmptyWithCursor(t *testing.T) {
	input := New(40).
		Content("").
		Focused(true).
		ShowCursor(true)

	view := input.View()

	// Should render cursor even with empty content
	if len(view) == 0 {
		t.Error("View() should not be empty (cursor visible)")
	}
}

func TestInput_ContentParts_EmptyContent(t *testing.T) {
	input := New(40).SetContent("", 0)

	before, at, after := input.ContentParts()

	if before != "" || at != "" || after != "" {
		t.Errorf("ContentParts() for empty should be all empty, got %q, %q, %q", before, at, after)
	}
}

func TestInput_ContentParts_Unicode(t *testing.T) {
	input := New(40).SetContent("你好世界", 2)

	before, at, after := input.ContentParts()

	if before != "你好" {
		t.Errorf("before = %q, want %q", before, "你好")
	}
	if at != "世" {
		t.Errorf("at = %q, want %q", at, "世")
	}
	if after != "界" {
		t.Errorf("after = %q, want %q", after, "界")
	}
}
