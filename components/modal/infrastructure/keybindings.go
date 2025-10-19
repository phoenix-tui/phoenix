package infrastructure

import (
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// KeyBindings defines key bindings for modal navigation.
type KeyBindings struct {
	Close          []string // Keys to close modal (default: Esc)
	NextButton     []string // Keys to focus next button (default: Tab, Right)
	PreviousButton []string // Keys to focus previous button (default: Shift+Tab, Left)
	ActivateButton []string // Keys to activate focused button (default: Enter)
}

// DefaultKeyBindings returns the default key bindings for modal navigation.
func DefaultKeyBindings() KeyBindings {
	return KeyBindings{
		Close:          []string{"esc"},
		NextButton:     []string{"tab", "→"},
		PreviousButton: []string{"shift+tab", "←"},
		ActivateButton: []string{"enter"},
	}
}

// IsClose checks if the key message is a close key.
func (kb KeyBindings) IsClose(msg tea.KeyMsg) bool {
	return kb.matchesKey(msg, kb.Close)
}

// IsNextButton checks if the key message is a next button key.
func (kb KeyBindings) IsNextButton(msg tea.KeyMsg) bool {
	return kb.matchesKey(msg, kb.NextButton)
}

// IsPreviousButton checks if the key message is a previous button key.
func (kb KeyBindings) IsPreviousButton(msg tea.KeyMsg) bool {
	return kb.matchesKey(msg, kb.PreviousButton)
}

// IsActivateButton checks if the key message is an activate button key.
func (kb KeyBindings) IsActivateButton(msg tea.KeyMsg) bool {
	return kb.matchesKey(msg, kb.ActivateButton)
}

// matchesKey checks if the key message matches any of the given keys.
func (kb KeyBindings) matchesKey(msg tea.KeyMsg, keys []string) bool {
	keyStr := msg.String()
	for _, k := range keys {
		if keyStr == k {
			return true
		}
	}
	return false
}
