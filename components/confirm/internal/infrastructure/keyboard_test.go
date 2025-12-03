package infrastructure

import (
	"testing"

	"github.com/phoenix-tui/phoenix/tea"
)

func TestKeyBindingMap_GetAction(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.KeyMsg
		want Action
	}{
		{
			name: "Left arrow",
			msg:  tea.KeyMsg{Type: tea.KeyLeft},
			want: ActionMoveLeft,
		},
		{
			name: "Shift+Tab",
			msg:  tea.KeyMsg{Type: tea.KeyTab, Shift: true},
			want: ActionMoveLeft,
		},
		{
			name: "Right arrow",
			msg:  tea.KeyMsg{Type: tea.KeyRight},
			want: ActionMoveRight,
		},
		{
			name: "Tab",
			msg:  tea.KeyMsg{Type: tea.KeyTab},
			want: ActionMoveRight,
		},
		{
			name: "Enter",
			msg:  tea.KeyMsg{Type: tea.KeyEnter},
			want: ActionConfirm,
		},
		{
			name: "Escape",
			msg:  tea.KeyMsg{Type: tea.KeyEsc},
			want: ActionCancel,
		},
		{
			name: "Ctrl+C",
			msg:  tea.KeyMsg{Type: tea.KeyCtrlC},
			want: ActionCancel,
		},
		{
			name: "Rune 'y'",
			msg:  tea.KeyMsg{Type: tea.KeyRune, Rune: 'y'},
			want: ActionShortcut,
		},
		{
			name: "Rune 'n'",
			msg:  tea.KeyMsg{Type: tea.KeyRune, Rune: 'n'},
			want: ActionShortcut,
		},
		{
			name: "Unknown key",
			msg:  tea.KeyMsg{Type: tea.KeyF1},
			want: ActionNone,
		},
	}

	km := DefaultKeyBindingMap()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := km.GetAction(tt.msg); got != tt.want {
				t.Errorf("GetAction() = %v, want %v", got, tt.want)
			}
		})
	}
}
