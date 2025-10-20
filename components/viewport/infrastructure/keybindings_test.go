package infrastructure

import (
	"testing"

	tea "github.com/phoenix-tui/phoenix/tea/api"
)

func TestDefaultKeyBindings(t *testing.T) {
	bindings := DefaultKeyBindings()

	expectedActions := []string{"up", "down", "pageup", "pagedown", "home", "end", "halfpageup", "halfpagedown"}

	for _, action := range expectedActions {
		binding, exists := bindings[action]
		if !exists {
			t.Errorf("Missing key binding for action: %s", action)
			continue
		}

		if len(binding.Keys) == 0 {
			t.Errorf("Action %s has no keys defined", action)
		}

		if binding.Help == "" {
			t.Errorf("Action %s has no help text", action)
		}
	}
}

func TestMatchKey(t *testing.T) {
	tests := []struct {
		name     string
		keyMsg   tea.KeyMsg
		keys     []string
		expected bool
	}{
		{
			name:     "matches first key",
			keyMsg:   tea.KeyMsg{Type: tea.KeyUp},
			keys:     []string{"↑", "k"},
			expected: true,
		},
		{
			name:     "matches second key",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRune, Rune: 'k'},
			keys:     []string{"↑", "k"},
			expected: true,
		},
		{
			name:     "no match",
			keyMsg:   tea.KeyMsg{Type: tea.KeyLeft},
			keys:     []string{"↑", "k"},
			expected: false,
		},
		{
			name:     "empty keys",
			keyMsg:   tea.KeyMsg{Type: tea.KeyUp},
			keys:     []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MatchKey(tt.keyMsg, tt.keys)
			if result != tt.expected {
				t.Errorf("MatchKey() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsUpKey(t *testing.T) {
	tests := []struct {
		name     string
		keyMsg   tea.KeyMsg
		expected bool
	}{
		{
			name:     "up arrow",
			keyMsg:   tea.KeyMsg{Type: tea.KeyUp},
			expected: true,
		},
		{
			name:     "k key",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRune, Rune: 'k'},
			expected: true,
		},
		{
			name:     "down arrow",
			keyMsg:   tea.KeyMsg{Type: tea.KeyDown},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsUpKey(tt.keyMsg)
			if result != tt.expected {
				t.Errorf("IsUpKey() = %v, want %v for key %s", result, tt.expected, tt.keyMsg.String())
			}
		})
	}
}

func TestIsDownKey(t *testing.T) {
	tests := []struct {
		name     string
		keyMsg   tea.KeyMsg
		expected bool
	}{
		{
			name:     "down arrow",
			keyMsg:   tea.KeyMsg{Type: tea.KeyDown},
			expected: true,
		},
		{
			name:     "j key",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRune, Rune: 'j'},
			expected: true,
		},
		{
			name:     "up arrow",
			keyMsg:   tea.KeyMsg{Type: tea.KeyUp},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDownKey(tt.keyMsg)
			if result != tt.expected {
				t.Errorf("IsDownKey() = %v, want %v for key %s", result, tt.expected, tt.keyMsg.String())
			}
		})
	}
}

func TestIsPageUpKey(t *testing.T) {
	tests := []struct {
		name     string
		keyMsg   tea.KeyMsg
		expected bool
	}{
		{
			name:     "pgup",
			keyMsg:   tea.KeyMsg{Type: tea.KeyPgUp},
			expected: true,
		},
		{
			name:     "b key",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRune, Rune: 'b'},
			expected: true,
		},
		{
			name:     "ctrl+b",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRune, Rune: 'b', Ctrl: true},
			expected: true,
		},
		{
			name:     "pgdown",
			keyMsg:   tea.KeyMsg{Type: tea.KeyPgDown},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPageUpKey(tt.keyMsg)
			if result != tt.expected {
				t.Errorf("IsPageUpKey() = %v, want %v for key %s", result, tt.expected, tt.keyMsg.String())
			}
		})
	}
}

func TestIsPageDownKey(t *testing.T) {
	tests := []struct {
		name     string
		keyMsg   tea.KeyMsg
		expected bool
	}{
		{
			name:     "pgdown",
			keyMsg:   tea.KeyMsg{Type: tea.KeyPgDown},
			expected: true,
		},
		{
			name:     "f key",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRune, Rune: 'f'},
			expected: true,
		},
		{
			name:     "ctrl+f",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRune, Rune: 'f', Ctrl: true},
			expected: true,
		},
		{
			name:     "space",
			keyMsg:   tea.KeyMsg{Type: tea.KeySpace},
			expected: true,
		},
		{
			name:     "pgup",
			keyMsg:   tea.KeyMsg{Type: tea.KeyPgUp},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPageDownKey(tt.keyMsg)
			if result != tt.expected {
				t.Errorf("IsPageDownKey() = %v, want %v for key %s", result, tt.expected, tt.keyMsg.String())
			}
		})
	}
}

func TestIsHomeKey(t *testing.T) {
	tests := []struct {
		name     string
		keyMsg   tea.KeyMsg
		expected bool
	}{
		{
			name:     "home",
			keyMsg:   tea.KeyMsg{Type: tea.KeyHome},
			expected: true,
		},
		{
			name:     "g key",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRune, Rune: 'g'},
			expected: true,
		},
		{
			name:     "end",
			keyMsg:   tea.KeyMsg{Type: tea.KeyEnd},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsHomeKey(tt.keyMsg)
			if result != tt.expected {
				t.Errorf("IsHomeKey() = %v, want %v for key %s", result, tt.expected, tt.keyMsg.String())
			}
		})
	}
}

func TestIsEndKey(t *testing.T) {
	tests := []struct {
		name     string
		keyMsg   tea.KeyMsg
		expected bool
	}{
		{
			name:     "end",
			keyMsg:   tea.KeyMsg{Type: tea.KeyEnd},
			expected: true,
		},
		{
			name:     "G key (shift+g)",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRune, Rune: 'G'},
			expected: true,
		},
		{
			name:     "home",
			keyMsg:   tea.KeyMsg{Type: tea.KeyHome},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEndKey(tt.keyMsg)
			if result != tt.expected {
				t.Errorf("IsEndKey() = %v, want %v for key %s", result, tt.expected, tt.keyMsg.String())
			}
		})
	}
}

func TestIsHalfPageUpKey(t *testing.T) {
	tests := []struct {
		name     string
		keyMsg   tea.KeyMsg
		expected bool
	}{
		{
			name:     "ctrl+u",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRune, Rune: 'u', Ctrl: true},
			expected: true,
		},
		{
			name:     "ctrl+d",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRune, Rune: 'd', Ctrl: true},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsHalfPageUpKey(tt.keyMsg)
			if result != tt.expected {
				t.Errorf("IsHalfPageUpKey() = %v, want %v for key %s", result, tt.expected, tt.keyMsg.String())
			}
		})
	}
}

func TestIsHalfPageDownKey(t *testing.T) {
	tests := []struct {
		name     string
		keyMsg   tea.KeyMsg
		expected bool
	}{
		{
			name:     "ctrl+d",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRune, Rune: 'd', Ctrl: true},
			expected: true,
		},
		{
			name:     "ctrl+u",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRune, Rune: 'u', Ctrl: true},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsHalfPageDownKey(tt.keyMsg)
			if result != tt.expected {
				t.Errorf("IsHalfPageDownKey() = %v, want %v for key %s", result, tt.expected, tt.keyMsg.String())
			}
		})
	}
}

func TestAllKeyBindingsCovered(t *testing.T) {
	// Ensure all key binding functions are tested.
	bindings := DefaultKeyBindings()

	testers := map[string]func(tea.KeyMsg) bool{
		"up":           IsUpKey,
		"down":         IsDownKey,
		"pageup":       IsPageUpKey,
		"pagedown":     IsPageDownKey,
		"home":         IsHomeKey,
		"end":          IsEndKey,
		"halfpageup":   IsHalfPageUpKey,
		"halfpagedown": IsHalfPageDownKey,
	}

	for action := range bindings {
		if _, exists := testers[action]; !exists {
			t.Errorf("No test function for key binding: %s", action)
		}
	}
}
