// Package infrastructure provides keybindings for viewport.
package infrastructure

import tea "github.com/phoenix-tui/phoenix/tea/api"

// KeyBinding represents a key binding for viewport actions.
type KeyBinding struct {
	Keys []string
	Help string
}

// DefaultKeyBindings returns the default key bindings for viewport navigation.
//
// Key string format follows tea.KeyMsg.String() representation:
//   - Arrow keys: "↑", "↓", "←", "→" (Unicode arrows)
//   - Special keys: "space", "enter", "backspace", etc.
//   - Modifiers: "ctrl+key", "alt+key", "shift+key".
//   - Regular runes: "k", "j", "f", "G", etc.
func DefaultKeyBindings() map[string]KeyBinding {
	return map[string]KeyBinding{
		"up": {
			Keys: []string{"↑", "k"},
			Help: "Scroll up one line",
		},
		"down": {
			Keys: []string{"↓", "j"},
			Help: "Scroll down one line",
		},
		"pageup": {
			Keys: []string{"pgup", "b", "ctrl+b"},
			Help: "Scroll up one page",
		},
		"pagedown": {
			Keys: []string{"pgdown", "f", "ctrl+f", "space"},
			Help: "Scroll down one page",
		},
		"home": {
			Keys: []string{"home", "g"},
			Help: "Scroll to top",
		},
		"end": {
			Keys: []string{"end", "G"},
			Help: "Scroll to bottom",
		},
		"halfpageup": {
			Keys: []string{"ctrl+u"},
			Help: "Scroll up half page",
		},
		"halfpagedown": {
			Keys: []string{"ctrl+d"},
			Help: "Scroll down half page",
		},
	}
}

// MatchKey checks if a key message matches any of the keys in the binding.
func MatchKey(msg tea.KeyMsg, keys []string) bool {
	keyStr := msg.String()

	for _, k := range keys {
		if keyStr == k {
			return true
		}
	}
	return false
}

// IsUpKey checks if the key message is an "up" key.
func IsUpKey(msg tea.KeyMsg) bool {
	bindings := DefaultKeyBindings()
	return MatchKey(msg, bindings["up"].Keys)
}

// IsDownKey checks if the key message is a "down" key.
func IsDownKey(msg tea.KeyMsg) bool {
	bindings := DefaultKeyBindings()
	return MatchKey(msg, bindings["down"].Keys)
}

// IsPageUpKey checks if the key message is a "page up" key.
func IsPageUpKey(msg tea.KeyMsg) bool {
	bindings := DefaultKeyBindings()
	return MatchKey(msg, bindings["pageup"].Keys)
}

// IsPageDownKey checks if the key message is a "page down" key.
func IsPageDownKey(msg tea.KeyMsg) bool {
	bindings := DefaultKeyBindings()
	return MatchKey(msg, bindings["pagedown"].Keys)
}

// IsHomeKey checks if the key message is a "home" key.
func IsHomeKey(msg tea.KeyMsg) bool {
	bindings := DefaultKeyBindings()
	return MatchKey(msg, bindings["home"].Keys)
}

// IsEndKey checks if the key message is an "end" key.
func IsEndKey(msg tea.KeyMsg) bool {
	bindings := DefaultKeyBindings()
	return MatchKey(msg, bindings["end"].Keys)
}

// IsHalfPageUpKey checks if the key message is a "half page up" key.
func IsHalfPageUpKey(msg tea.KeyMsg) bool {
	bindings := DefaultKeyBindings()
	return MatchKey(msg, bindings["halfpageup"].Keys)
}

// IsHalfPageDownKey checks if the key message is a "half page down" key.
func IsHalfPageDownKey(msg tea.KeyMsg) bool {
	bindings := DefaultKeyBindings()
	return MatchKey(msg, bindings["halfpagedown"].Keys)
}
