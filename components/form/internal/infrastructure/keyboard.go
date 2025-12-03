// Package infrastructure contains infrastructure implementations for the form component.
package infrastructure

import (
	"github.com/phoenix-tui/phoenix/tea"
)

// Action represents a keyboard action.
type Action int

const (
	// ActionNone represents no action.
	ActionNone Action = iota
	// ActionNextField moves focus to the next field.
	ActionNextField
	// ActionPrevField moves focus to the previous field.
	ActionPrevField
	// ActionSubmit submits the form.
	ActionSubmit
	// ActionReset resets the form.
	ActionReset
	// ActionQuit quits the application.
	ActionQuit
)

// KeyBindingMap maps keyboard inputs to actions.
type KeyBindingMap struct {
	bindings map[tea.KeyType]Action
	runes    map[rune]Action
}

// DefaultKeyBindingMap returns the default keyboard bindings for forms.
func DefaultKeyBindingMap() *KeyBindingMap {
	return &KeyBindingMap{
		bindings: map[tea.KeyType]Action{
			tea.KeyTab:   ActionNextField,
			tea.KeyEnter: ActionSubmit,
			tea.KeyCtrlC: ActionQuit,
			tea.KeyEsc:   ActionQuit,
		},
		runes: map[rune]Action{},
	}
}

// GetAction returns the action for the given key message.
func (k *KeyBindingMap) GetAction(msg tea.KeyMsg) Action {
	// Check Shift+Tab for previous field
	if msg.Type == tea.KeyTab && msg.Shift {
		return ActionPrevField
	}

	// Check key type bindings
	if action, ok := k.bindings[msg.Type]; ok {
		return action
	}

	// Check rune bindings
	if msg.Type == tea.KeyRune {
		if action, ok := k.runes[msg.Rune]; ok {
			return action
		}
	}

	return ActionNone
}
