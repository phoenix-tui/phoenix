package infrastructure

import (
	"testing"

	"github.com/phoenix-tui/phoenix/tea"
)

func TestDefaultKeyBindingMap(t *testing.T) {
	km := DefaultKeyBindingMap()

	tests := []struct {
		name string
		msg  tea.KeyMsg
		want Action
	}{
		// Type-based bindings
		{"up arrow", tea.KeyMsg{Type: tea.KeyUp}, ActionMoveUp},
		{"down arrow", tea.KeyMsg{Type: tea.KeyDown}, ActionMoveDown},
		{"home", tea.KeyMsg{Type: tea.KeyHome}, ActionMoveToStart},
		{"end", tea.KeyMsg{Type: tea.KeyEnd}, ActionMoveToEnd},
		{"space", tea.KeyMsg{Type: tea.KeySpace}, ActionToggle},
		{"enter", tea.KeyMsg{Type: tea.KeyEnter}, ActionConfirm},
		{"esc", tea.KeyMsg{Type: tea.KeyEsc}, ActionClearFilter},
		{"ctrl+c", tea.KeyMsg{Type: tea.KeyCtrlC}, ActionQuit},

		// Rune-based bindings
		{"k", tea.KeyMsg{Type: tea.KeyRune, Rune: 'k'}, ActionMoveUp},
		{"j", tea.KeyMsg{Type: tea.KeyRune, Rune: 'j'}, ActionMoveDown},
		{"g", tea.KeyMsg{Type: tea.KeyRune, Rune: 'g'}, ActionMoveToStart},
		{"G", tea.KeyMsg{Type: tea.KeyRune, Rune: 'G'}, ActionMoveToEnd},
		{"a", tea.KeyMsg{Type: tea.KeyRune, Rune: 'a'}, ActionSelectAll},
		{"n", tea.KeyMsg{Type: tea.KeyRune, Rune: 'n'}, ActionSelectNone},

		// Unmapped keys
		{"x", tea.KeyMsg{Type: tea.KeyRune, Rune: 'x'}, ActionNone},
		{"tab", tea.KeyMsg{Type: tea.KeyTab}, ActionNone},
		{"1", tea.KeyMsg{Type: tea.KeyRune, Rune: '1'}, ActionNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := km.GetAction(tt.msg)
			if got != tt.want {
				t.Errorf("GetAction(%v) = %v, want %v", tt.msg, got, tt.want)
			}
		})
	}
}

func TestKeyBindingMap_GetAction_Priority(t *testing.T) {
	km := DefaultKeyBindingMap()

	// Type-based bindings should take precedence over rune-based
	msg := tea.KeyMsg{Type: tea.KeyUp, Rune: 'k'}
	action := km.GetAction(msg)

	if action != ActionMoveUp {
		t.Errorf("GetAction(up+k) = %v, want %v (type priority)", action, ActionMoveUp)
	}
}

func TestKeyBindingMap_Actions(t *testing.T) {
	// Verify all action constants exist
	actions := []Action{
		ActionMoveUp,
		ActionMoveDown,
		ActionMoveToStart,
		ActionMoveToEnd,
		ActionToggle,
		ActionSelectAll,
		ActionSelectNone,
		ActionConfirm,
		ActionClearFilter,
		ActionQuit,
		ActionNone,
	}

	for _, action := range actions {
		if string(action) == "" {
			t.Errorf("Action %v has empty string", action)
		}
	}
}
