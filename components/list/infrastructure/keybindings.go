package infrastructure

import (
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// KeyBinding represents a key and its action.
type KeyBinding struct {
	Key    string
	Action string
}

// DefaultKeyBindings returns the default key bindings for list navigation.
func DefaultKeyBindings() []KeyBinding {
	return []KeyBinding{
		// Up movement.
		{Key: "↑", Action: "move_up"}, // Unicode arrow from tea
		{Key: "k", Action: "move_up"},

		// Down movement.
		{Key: "↓", Action: "move_down"}, // Unicode arrow from tea
		{Key: "j", Action: "move_down"},

		// Page movement.
		{Key: "pgup", Action: "page_up"},
		{Key: "ctrl+u", Action: "page_up"},

		{Key: "pgdown", Action: "page_down"},
		{Key: "ctrl+d", Action: "page_down"},

		// Start/End.
		{Key: "home", Action: "move_to_start"},
		{Key: "g", Action: "move_to_start"},

		{Key: "end", Action: "move_to_end"},
		{Key: "G", Action: "move_to_end"},

		// Selection.
		{Key: "space", Action: "toggle_selection"}, // Space key
		{Key: "enter", Action: "confirm"},

		// Multi-select specific.
		{Key: "ctrl+a", Action: "select_all"},
		{Key: "esc", Action: "clear_selection"},

		// Filter.
		{Key: "ctrl+c", Action: "clear_filter"},

		// Quit.
		{Key: "q", Action: "quit"},
		{Key: "ctrl+c", Action: "quit"}, // Also mapped to quit in some contexts
	}
}

// KeyBindingMap creates a map from key to action for fast lookup.
type KeyBindingMap map[string]string

// NewKeyBindingMap creates a key binding map from a slice of key bindings.
func NewKeyBindingMap(bindings []KeyBinding) KeyBindingMap {
	m := make(KeyBindingMap)
	for _, binding := range bindings {
		m[binding.Key] = binding.Action
	}
	return m
}

// GetAction returns the action for a given key message.
func (m KeyBindingMap) GetAction(msg tea.KeyMsg) string {
	key := msg.String()
	if action, ok := m[key]; ok {
		return action
	}
	return ""
}

// DefaultKeyBindingMap returns the default key binding map.
func DefaultKeyBindingMap() KeyBindingMap {
	return NewKeyBindingMap(DefaultKeyBindings())
}
