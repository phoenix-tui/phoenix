// Package infrastructure provides keybindings implementations.
package infrastructure

import (
	"github.com/phoenix-tui/phoenix/components/input/domain/model"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// KeyHandler is a function that handles a key message and returns an updated input.
// VALUE SEMANTICS - takes value, returns value!
type KeyHandler func(model.TextInput, tea.KeyMsg) model.TextInput

// DefaultKeyBindings provides the default key bindings for TextInput.
// Applications can override these by providing custom handlers.
type DefaultKeyBindings struct{}

// NewDefaultKeyBindings creates a new default key bindings handler.
func NewDefaultKeyBindings() *DefaultKeyBindings {
	return &DefaultKeyBindings{}
}

// Handle processes a key message and returns the updated input.
// VALUE SEMANTICS - takes value, returns value!
//nolint:gocyclo,cyclop // keybindings require state machine logic
func (kb *DefaultKeyBindings) Handle(input model.TextInput, msg tea.KeyMsg) model.TextInput {
	// Handle Ctrl key combinations.
	if msg.Ctrl {
		switch msg.Rune {
		case 'a', 'A':
			// Ctrl-A selects all.
			return input.SelectAll()
		case 'u', 'U':
			// Ctrl-U clears input.
			return input.Clear()
		case 'e', 'E':
			// Ctrl-E moves to end.
			return input.MoveEnd()
		}
	}

	switch msg.Type {
	// Navigation.
	case tea.KeyLeft:
		return input.MoveLeft()

	case tea.KeyRight:
		return input.MoveRight()

	case tea.KeyHome:
		return input.MoveHome()

	case tea.KeyEnd:
		return input.MoveEnd()

	// Editing.
	case tea.KeyBackspace:
		return input.DeleteBackward()

	case tea.KeyDelete:
		return input.DeleteForward()

	// Character input.
	case tea.KeySpace:
		// Insert space character (0x20 is parsed as KeySpace, not KeyRune)
		return input.InsertRune(' ')

	case tea.KeyRune:
		// Insert the rune.
		return input.InsertRune(msg.Rune)

	default:
		// Key not handled by default bindings - return input unchanged.
		return input
	}
}

// IsNavigationKey returns true if the key is a navigation key (doesn't modify content).
func IsNavigationKey(msg tea.KeyMsg) bool {
	// Ctrl+A (select all) is navigation.
	if msg.Ctrl && (msg.Rune == 'a' || msg.Rune == 'A') {
		return true
	}

	switch msg.Type {
	case tea.KeyLeft, tea.KeyRight, tea.KeyHome, tea.KeyEnd:
		return true
	default:
		return false
	}
}

// IsEditingKey returns true if the key modifies content.
func IsEditingKey(msg tea.KeyMsg) bool {
	// Ctrl+U (clear) is editing.
	if msg.Ctrl && (msg.Rune == 'u' || msg.Rune == 'U') {
		return true
	}

	switch msg.Type {
	case tea.KeyBackspace, tea.KeyDelete, tea.KeySpace, tea.KeyRune:
		return true
	default:
		return false
	}
}

// CustomKeyBindings allows applications to provide custom key handlers.
// Handlers are checked in order; first non-nil result wins.
type CustomKeyBindings struct {
	handlers []KeyHandler
	fallback *DefaultKeyBindings
}

// NewCustomKeyBindings creates a new custom key bindings handler.
func NewCustomKeyBindings(handlers ...KeyHandler) *CustomKeyBindings {
	return &CustomKeyBindings{
		handlers: handlers,
		fallback: NewDefaultKeyBindings(),
	}
}

// Handle processes a key message through custom handlers first,.
// then falls back to default bindings if no custom handler handles it.
// VALUE SEMANTICS - takes value, returns value!
func (kb *CustomKeyBindings) Handle(input model.TextInput, msg tea.KeyMsg) model.TextInput {
	// Try custom handlers first.
	for _, handler := range kb.handlers {
		result := handler(input, msg)
		// Check if handler modified the input.
		if result.Content() != input.Content() || result.CursorPosition() != input.CursorPosition() {
			return result
		}
	}

	// Fall back to default bindings.
	return kb.fallback.Handle(input, msg)
}

// AddHandler adds a custom key handler.
func (kb *CustomKeyBindings) AddHandler(handler KeyHandler) {
	kb.handlers = append(kb.handlers, handler)
}
