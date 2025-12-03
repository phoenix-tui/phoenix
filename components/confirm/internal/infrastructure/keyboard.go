// Package infrastructure provides keyboard handling for the confirm component.
package infrastructure

import (
	"github.com/phoenix-tui/phoenix/tea"
)

// Action represents an action that can be performed on the confirm dialog.
type Action int

const (
	// ActionNone means no action.
	ActionNone Action = iota
	// ActionMoveLeft moves focus to the previous button.
	ActionMoveLeft
	// ActionMoveRight moves focus to the next button.
	ActionMoveRight
	// ActionConfirm confirms the current selection.
	ActionConfirm
	// ActionCancel cancels the dialog.
	ActionCancel
	// ActionShortcut triggers a button by its keyboard shortcut.
	ActionShortcut
)

// KeyBindingMap maps keyboard input to actions.
type KeyBindingMap struct {
	// No configuration needed for now - hard-coded bindings
}

// DefaultKeyBindingMap returns the default key bindings.
func DefaultKeyBindingMap() *KeyBindingMap {
	return &KeyBindingMap{}
}

// GetAction returns the action for the given key message.
func (k *KeyBindingMap) GetAction(msg tea.KeyMsg) Action {
	switch msg.Type {
	case tea.KeyLeft:
		return ActionMoveLeft
	case tea.KeyRight:
		return ActionMoveRight
	case tea.KeyTab:
		// Shift+Tab goes left, Tab goes right
		if msg.Shift {
			return ActionMoveLeft
		}
		return ActionMoveRight
	case tea.KeyEnter:
		return ActionConfirm
	case tea.KeyEsc, tea.KeyCtrlC:
		return ActionCancel
	case tea.KeyRune:
		// Any rune is a potential shortcut key
		return ActionShortcut
	}
	return ActionNone
}
