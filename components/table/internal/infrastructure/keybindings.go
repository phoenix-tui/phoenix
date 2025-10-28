// Package infrastructure provides technical implementations for the table component.
package infrastructure

import tea "github.com/phoenix-tui/phoenix/tea"

// KeyBindings defines the default keyboard shortcuts for table navigation.
type KeyBindings struct {
	Up        []string // Move selection up
	Down      []string // Move selection down
	PageUp    []string // Move up by page
	PageDown  []string // Move down by page
	Home      []string // Move to first row
	End       []string // Move to last row
	Sort      []string // Toggle sort on current column
	ClearSort []string // Clear all sorting
}

// DefaultKeyBindings returns the default key bindings for table navigation.
func DefaultKeyBindings() KeyBindings {
	return KeyBindings{
		Up:        []string{"↑", "k"},
		Down:      []string{"↓", "j"},
		PageUp:    []string{"pgup"},
		PageDown:  []string{"pgdown"},
		Home:      []string{"home", "g"},
		End:       []string{"end", "G"},
		Sort:      []string{"s", "enter"},
		ClearSort: []string{"c"},
	}
}

// IsUp returns true if the key message matches an "up" binding.
func (kb KeyBindings) IsUp(msg tea.KeyMsg) bool {
	return kb.matchesAny(msg, kb.Up)
}

// IsDown returns true if the key message matches a "down" binding.
func (kb KeyBindings) IsDown(msg tea.KeyMsg) bool {
	return kb.matchesAny(msg, kb.Down)
}

// IsPageUp returns true if the key message matches a "page up" binding.
func (kb KeyBindings) IsPageUp(msg tea.KeyMsg) bool {
	return kb.matchesAny(msg, kb.PageUp)
}

// IsPageDown returns true if the key message matches a "page down" binding.
func (kb KeyBindings) IsPageDown(msg tea.KeyMsg) bool {
	return kb.matchesAny(msg, kb.PageDown)
}

// IsHome returns true if the key message matches a "home" binding.
func (kb KeyBindings) IsHome(msg tea.KeyMsg) bool {
	return kb.matchesAny(msg, kb.Home)
}

// IsEnd returns true if the key message matches an "end" binding.
func (kb KeyBindings) IsEnd(msg tea.KeyMsg) bool {
	return kb.matchesAny(msg, kb.End)
}

// IsSort returns true if the key message matches a "sort" binding.
func (kb KeyBindings) IsSort(msg tea.KeyMsg) bool {
	return kb.matchesAny(msg, kb.Sort)
}

// IsClearSort returns true if the key message matches a "clear sort" binding.
func (kb KeyBindings) IsClearSort(msg tea.KeyMsg) bool {
	return kb.matchesAny(msg, kb.ClearSort)
}

// matchesAny returns true if the key message matches any of the bindings.
func (kb KeyBindings) matchesAny(msg tea.KeyMsg, bindings []string) bool {
	key := msg.String()
	for _, binding := range bindings {
		if key == binding {
			return true
		}
	}
	return false
}
