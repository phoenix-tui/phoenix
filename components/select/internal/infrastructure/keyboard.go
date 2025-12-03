// Package infrastructure provides technical implementations for the select component.
package infrastructure

import (
	"github.com/phoenix-tui/phoenix/tea"
)

// Action represents a keyboard action.
type Action string

// Keyboard actions for Select component navigation and selection.
const (
	ActionMoveUp      Action = "move_up"       // Move cursor up one position
	ActionMoveDown    Action = "move_down"     // Move cursor down one position
	ActionMoveToStart Action = "move_to_start" // Move cursor to first option
	ActionMoveToEnd   Action = "move_to_end"   // Move cursor to last option
	ActionSelect      Action = "select"        // Confirm selection
	ActionClearFilter Action = "clear_filter"  // Clear filter query
	ActionQuit        Action = "quit"          // Quit the application
	ActionNone        Action = "none"          // No action (unmapped key)
)

// KeyBindingMap maps key messages to actions.
type KeyBindingMap struct {
	bindings map[tea.KeyType]Action
	runeMap  map[rune]Action
}

// DefaultKeyBindingMap returns the default keyboard bindings for Select.
func DefaultKeyBindingMap() *KeyBindingMap {
	return &KeyBindingMap{
		bindings: map[tea.KeyType]Action{
			tea.KeyUp:    ActionMoveUp,
			tea.KeyDown:  ActionMoveDown,
			tea.KeyHome:  ActionMoveToStart,
			tea.KeyEnd:   ActionMoveToEnd,
			tea.KeyEnter: ActionSelect,
			tea.KeyEsc:   ActionClearFilter,
			tea.KeyCtrlC: ActionQuit,
		},
		runeMap: map[rune]Action{
			'k': ActionMoveUp,
			'j': ActionMoveDown,
			'g': ActionMoveToStart,
			'G': ActionMoveToEnd,
		},
	}
}

// GetAction returns the action for the given key message.
func (k *KeyBindingMap) GetAction(msg tea.KeyMsg) Action {
	// Check type-based bindings first
	if action, ok := k.bindings[msg.Type]; ok {
		return action
	}

	// Check rune-based bindings
	if msg.Type == tea.KeyRune {
		if action, ok := k.runeMap[msg.Rune]; ok {
			return action
		}
	}

	return ActionNone
}
