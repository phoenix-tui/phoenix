package keybindings

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/input/textarea/domain/model"
	"github.com/phoenix-tui/phoenix/tea/api"
)

// TestEmacsKeybindings_KeySpace verifies that Space key is handled correctly.
// This is a CRITICAL test - without KeySpace handling, spaces are ignored.
func TestEmacsKeybindings_KeySpace(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert space
	msg := api.KeyMsg{Type: api.KeySpace}
	result, _ := handler.Handle(msg, ta)

	// Verify space was inserted
	if result.Value() != " " {
		t.Errorf("Expected single space, got %q", result.Value())
	}

	// Insert another space
	result, _ = handler.Handle(msg, result)

	// Verify two spaces
	if result.Value() != "  " {
		t.Errorf("Expected two spaces, got %q", result.Value())
	}
}

// TestEmacsKeybindings_SpaceInText verifies space insertion within text.
func TestEmacsKeybindings_SpaceInText(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	for _, r := range "hello" {
		msg := api.KeyMsg{Type: api.KeyRune, Rune: r}
		ta, _ = handler.Handle(msg, ta)
	}

	// Insert space
	spaceMsg := api.KeyMsg{Type: api.KeySpace}
	ta, _ = handler.Handle(spaceMsg, ta)

	// Insert "world"
	for _, r := range "world" {
		msg := api.KeyMsg{Type: api.KeyRune, Rune: r}
		ta, _ = handler.Handle(msg, ta)
	}

	// Verify result
	expected := "hello world"
	if ta.Value() != expected {
		t.Errorf("Expected %q, got %q", expected, ta.Value())
	}
}

// TestEmacsKeybindings_MultipleSpaces verifies multiple consecutive spaces.
func TestEmacsKeybindings_MultipleSpaces(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert 5 spaces
	spaceMsg := api.KeyMsg{Type: api.KeySpace}
	for i := 0; i < 5; i++ {
		ta, _ = handler.Handle(spaceMsg, ta)
	}

	// Verify 5 spaces
	expected := "     "
	if ta.Value() != expected {
		t.Errorf("Expected 5 spaces, got %q (len=%d)", ta.Value(), len(ta.Value()))
	}
}
