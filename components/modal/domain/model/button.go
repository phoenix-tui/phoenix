package model

// Button represents a modal action button.
// Buttons are displayed at the bottom of the modal and can be activated.
// via keyboard shortcuts or navigation (Tab + Enter).
type Button struct {
	label  string // Button text (e.g., "Yes", "No", "OK")
	key    string // Keyboard shortcut (e.g., "y" for Yes, "n" for No)
	action string // Action identifier (e.g., "confirm", "cancel")
}

// NewButton creates a new button with the given label, key, and action.
// - label: Text displayed on the button.
// - key: Keyboard shortcut (single character, case-insensitive)
// - action: Identifier sent in ButtonPressedMsg.
func NewButton(label, key, action string) *Button {
	return &Button{
		label:  label,
		key:    key,
		action: action,
	}
}

// Label returns the button text.
func (b *Button) Label() string {
	return b.label
}

// Key returns the keyboard shortcut.
func (b *Button) Key() string {
	return b.key
}

// Action returns the action identifier.
func (b *Button) Action() string {
	return b.action
}
