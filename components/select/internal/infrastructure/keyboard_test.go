package infrastructure

import (
	"testing"

	"github.com/phoenix-tui/phoenix/tea"
)

func TestDefaultKeyBindingMap(t *testing.T) {
	km := DefaultKeyBindingMap()

	tests := []struct {
		name     string
		key      tea.KeyMsg
		expected Action
	}{
		{"up arrow", tea.KeyMsg{Type: tea.KeyUp}, ActionMoveUp},
		{"down arrow", tea.KeyMsg{Type: tea.KeyDown}, ActionMoveDown},
		{"home key", tea.KeyMsg{Type: tea.KeyHome}, ActionMoveToStart},
		{"end key", tea.KeyMsg{Type: tea.KeyEnd}, ActionMoveToEnd},
		{"enter key", tea.KeyMsg{Type: tea.KeyEnter}, ActionSelect},
		{"escape key", tea.KeyMsg{Type: tea.KeyEsc}, ActionClearFilter},
		{"ctrl+c", tea.KeyMsg{Type: tea.KeyCtrlC}, ActionQuit},
		{"k key", tea.KeyMsg{Type: tea.KeyRune, Rune: 'k'}, ActionMoveUp},
		{"j key", tea.KeyMsg{Type: tea.KeyRune, Rune: 'j'}, ActionMoveDown},
		{"g key", tea.KeyMsg{Type: tea.KeyRune, Rune: 'g'}, ActionMoveToStart},
		{"G key", tea.KeyMsg{Type: tea.KeyRune, Rune: 'G'}, ActionMoveToEnd},
		{"unmapped key", tea.KeyMsg{Type: tea.KeyRune, Rune: 'x'}, ActionNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := km.GetAction(tt.key)
			if action != tt.expected {
				t.Errorf("expected action %q, got %q", tt.expected, action)
			}
		})
	}
}

func TestKeyBindingMapGetAction(t *testing.T) {
	t.Run("returns ActionNone for unmapped keys", func(t *testing.T) {
		km := DefaultKeyBindingMap()
		action := km.GetAction(tea.KeyMsg{Type: tea.KeyF1})
		if action != ActionNone {
			t.Errorf("expected ActionNone, got %q", action)
		}
	})

	t.Run("prefers type-based bindings over rune bindings", func(t *testing.T) {
		km := DefaultKeyBindingMap()
		action := km.GetAction(tea.KeyMsg{Type: tea.KeyUp})
		if action != ActionMoveUp {
			t.Errorf("expected ActionMoveUp, got %q", action)
		}
	})
}
