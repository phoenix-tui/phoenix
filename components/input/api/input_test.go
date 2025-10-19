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

	// Pointer chaining from New
	input := New(40).
		Validator(validator).
		Content("test")

	if input.IsValid() {
		t.Error("IsValid() = true, want false (missing @)")
	}

	// After first method, we have Input value
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
	// Pointer chaining from New
	input := New(40).Content("hello").Focused(true)

	// Send left arrow key
	msg := tea.KeyMsg{Type: tea.KeyLeft}
	updated, cmd := input.Update(msg)

	if cmd != nil {
		t.Error("Update() should return nil cmd")
	}

	// Cursor should move (assuming default position at end)
	// After typing "hello", cursor is at 5, left arrow moves to 4
	inputAtEnd := input.SetContent("hello", 5)
	updated, _ = inputAtEnd.Update(msg)

	if updated.CursorPosition() != 4 {
		t.Errorf("CursorPosition() = %d, want 4 after left arrow", updated.CursorPosition())
	}
}

func TestInput_Update_Unfocused(t *testing.T) {
	input := New(40).Content("hello").Focused(false)

	// Send left arrow key (should be ignored)
	msg := tea.KeyMsg{Type: tea.KeyLeft}
	updated, _ := input.Update(msg)

	// Should be unchanged
	if updated.Value() != input.Value() {
		t.Error("unfocused input should ignore keys")
	}
}

func TestInput_Update_InsertText(t *testing.T) {
	// Method chaining returns Input value after first call
	input := New(40).Content("").Focused(true)

	// Insert 'h'
	msg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'h'}
	input, _ = input.Update(msg) // Reassignment!

	if input.Value() != "h" {
		t.Errorf("Value() = %q, want %q", input.Value(), "h")
	}

	// Insert 'i'
	msg = tea.KeyMsg{Type: tea.KeyRune, Rune: 'i'}
	input, _ = input.Update(msg) // Reassignment!

	if input.Value() != "hi" {
		t.Errorf("Value() = %q, want %q", input.Value(), "hi")
	}
}

func TestInput_Update_Backspace(t *testing.T) {
	// Method chaining returns Input value
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

	// Should contain the content
	if !strings.Contains(view, "he") {
		t.Errorf("View() should contain 'he', got %q", view)
	}

	// Should have some cursor rendering
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

	// Wrap in KeyBindingHandler interface
	wrappedHandler := CustomKeyBindings(customHandler)

	// Can't directly test this without exposing more internals
	// This test verifies the API exists
	if wrappedHandler == nil {
		t.Error("CustomKeyBindings() returned nil")
	}
}

func TestInput_CommonValidators(t *testing.T) {
	// Test NotEmpty
	input := New(40).Validator(NotEmpty()).Content("")
	if input.IsValid() {
		t.Error("empty should be invalid with NotEmpty validator")
	}

	input = input.Content("text")
	if !input.IsValid() {
		t.Error("non-empty should be valid with NotEmpty validator")
	}

	// Test MinLength
	input = New(40).Validator(MinLength(5)).Content("abc")
	if input.IsValid() {
		t.Error("'abc' should be invalid with MinLength(5)")
	}

	input = input.Content("abcde")
	if !input.IsValid() {
		t.Error("'abcde' should be valid with MinLength(5)")
	}

	// Test MaxLength
	input = New(40).Validator(MaxLength(5)).Content("abcdef")
	if input.IsValid() {
		t.Error("'abcdef' should be invalid with MaxLength(5)")
	}

	input = input.Content("abcde")
	if !input.IsValid() {
		t.Error("'abcde' should be valid with MaxLength(5)")
	}

	// Test Range
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

	// Test Chain
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
	// Just verify the errors are exported
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
	// Method chaining returns Input value
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
	// Simulate a complete user interaction
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

	// Should be valid now
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
