package infrastructure

import (
	"testing"

	tea "github.com/phoenix-tui/phoenix/tea/api"
)

func TestDefaultKeyBindings(t *testing.T) {
	bindings := DefaultKeyBindings()

	if len(bindings) == 0 {
		t.Error("DefaultKeyBindings() should return non-empty bindings")
	}

	// Check that essential bindings exist.
	essentialActions := map[string]bool{
		"move_up":          false,
		"move_down":        false,
		"toggle_selection": false,
	}

	for _, binding := range bindings {
		if _, ok := essentialActions[binding.Action]; ok {
			essentialActions[binding.Action] = true
		}
	}

	for action, found := range essentialActions {
		if !found {
			t.Errorf("DefaultKeyBindings() missing essential action: %s", action)
		}
	}
}

func TestNewKeyBindingMap(t *testing.T) {
	bindings := []KeyBinding{
		{Key: "up", Action: "move_up"},
		{Key: "down", Action: "move_down"},
		{Key: "k", Action: "move_up"},
	}

	m := NewKeyBindingMap(bindings)

	if len(m) != 3 {
		t.Errorf("NewKeyBindingMap() map size = %d, want 3", len(m))
	}

	if m["up"] != "move_up" {
		t.Errorf("NewKeyBindingMap()[up] = %s, want move_up", m["up"])
	}
	if m["down"] != "move_down" {
		t.Errorf("NewKeyBindingMap()[down] = %s, want move_down", m["down"])
	}
	if m["k"] != "move_up" {
		t.Errorf("NewKeyBindingMap()[k] = %s, want move_up", m["k"])
	}
}

func TestKeyBindingMap_GetAction(t *testing.T) {
	m := KeyBindingMap{
		"↑":     "move_up",   // Unicode arrow (as returned by KeyMsg.String())
		"↓":     "move_down", // Unicode arrow (as returned by KeyMsg.String())
		"enter": "confirm",
	}

	tests := []struct {
		name string
		key  tea.KeyMsg
		want string
	}{
		{
			name: "existing key up",
			key: tea.KeyMsg{
				Type: tea.KeyUp,
			},
			want: "move_up",
		},
		{
			name: "existing key down",
			key: tea.KeyMsg{
				Type: tea.KeyDown,
			},
			want: "move_down",
		},
		{
			name: "existing key enter",
			key: tea.KeyMsg{
				Type: tea.KeyEnter,
			},
			want: "confirm",
		},
		{
			name: "non-existing key",
			key: tea.KeyMsg{
				Type: tea.KeyRune,
				Rune: 'x',
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.GetAction(tt.key); got != tt.want {
				t.Errorf("KeyBindingMap.GetAction() = %v, want %v (key: %s)", got, tt.want, tt.key.String())
			}
		})
	}
}

func TestDefaultKeyBindingMap(t *testing.T) {
	m := DefaultKeyBindingMap()

	if len(m) == 0 {
		t.Error("DefaultKeyBindingMap() should return non-empty map")
	}

	// Test some standard key mappings.
	tests := []struct {
		key    string
		action string
	}{
		{"↑", "move_up"},
		{"↓", "move_down"},
		{"k", "move_up"},
		{"j", "move_down"},
		{"space", "toggle_selection"},
	}

	for _, tt := range tests {
		if got := m[tt.key]; got != tt.action {
			t.Errorf("DefaultKeyBindingMap()[%s] = %s, want %s", tt.key, got, tt.action)
		}
	}
}

func TestKeyBindingMap_GetAction_IntegrationWithDefaultBindings(t *testing.T) {
	m := DefaultKeyBindingMap()

	// Test common navigation keys.
	upKey := tea.KeyMsg{Type: tea.KeyUp}
	if action := m.GetAction(upKey); action != "move_up" {
		t.Errorf("GetAction(up) = %s, want move_up", action)
	}

	downKey := tea.KeyMsg{Type: tea.KeyDown}
	if action := m.GetAction(downKey); action != "move_down" {
		t.Errorf("GetAction(down) = %s, want move_down", action)
	}

	// Test vim keys.
	jKey := tea.KeyMsg{Type: tea.KeyRune, Rune: 'j'}
	if action := m.GetAction(jKey); action != "move_down" {
		t.Errorf("GetAction(j) = %s, want move_down", action)
	}

	kKey := tea.KeyMsg{Type: tea.KeyRune, Rune: 'k'}
	if action := m.GetAction(kKey); action != "move_up" {
		t.Errorf("GetAction(k) = %s, want move_up", action)
	}
}

func TestKeyBindingMap_GetAction_EmptyMap(t *testing.T) {
	m := KeyBindingMap{}

	key := tea.KeyMsg{Type: tea.KeyUp}
	if action := m.GetAction(key); action != "" {
		t.Errorf("GetAction() on empty map should return empty string, got %s", action)
	}
}

func TestKeyBinding_Duplicates(t *testing.T) {
	// Test that duplicate keys (last one wins) work as expected.
	bindings := []KeyBinding{
		{Key: "k", Action: "move_up"},
		{Key: "k", Action: "kill"}, // Duplicate - should overwrite
	}

	m := NewKeyBindingMap(bindings)

	if m["k"] != "kill" {
		t.Errorf("NewKeyBindingMap() with duplicate keys should use last value, got %s", m["k"])
	}
}
