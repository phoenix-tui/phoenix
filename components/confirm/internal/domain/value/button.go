// Package value provides value objects for the confirm component.
package value

// Button represents a button option in the confirm dialog.
// This is a value object - immutable and defined by its attributes.
type Button struct {
	label string
	key   rune // Keyboard shortcut (e.g., 'y', 'n', 'c')
}

// NewButton creates a new Button with the given label.
// The key is automatically derived from the first letter of the label.
func NewButton(label string) *Button {
	var key rune
	if label != "" {
		key = rune(label[0])
	}
	return &Button{
		label: label,
		key:   key,
	}
}

// Label returns the button's label.
func (b *Button) Label() string {
	return b.label
}

// Key returns the button's keyboard shortcut.
func (b *Button) Key() rune {
	return b.key
}

// MatchesKey returns true if the given rune matches this button's key (case-insensitive).
func (b *Button) MatchesKey(r rune) bool {
	// Case-insensitive comparison
	keyLower := toLower(b.key)
	rLower := toLower(r)
	return keyLower == rLower
}

// toLower converts a rune to lowercase (simple ASCII conversion).
func toLower(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + 32
	}
	return r
}
