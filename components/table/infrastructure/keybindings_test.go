package infrastructure

import (
	"testing"

	tea "github.com/phoenix-tui/phoenix/tea/api"
)

func TestDefaultKeyBindings(t *testing.T) {
	kb := DefaultKeyBindings()

	if len(kb.Up) == 0 {
		t.Errorf("Up bindings should not be empty")
	}
	if len(kb.Down) == 0 {
		t.Errorf("Down bindings should not be empty")
	}
	if len(kb.Home) == 0 {
		t.Errorf("Home bindings should not be empty")
	}
	if len(kb.End) == 0 {
		t.Errorf("End bindings should not be empty")
	}
}

func TestKeyBindings_IsUp(t *testing.T) {
	kb := DefaultKeyBindings()

	tests := []struct {
		name string
		msg  tea.KeyMsg
		want bool
	}{
		{
			name: "UpArrow",
			msg:  tea.KeyMsg{Type: tea.KeyUp},
			want: true,
		},
		{
			name: "K",
			msg:  tea.KeyMsg{Type: tea.KeyRune, Rune: 'k'},
			want: true,
		},
		{
			name: "Down",
			msg:  tea.KeyMsg{Type: tea.KeyDown},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := kb.IsUp(tt.msg)
			if got != tt.want {
				t.Errorf("IsUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyBindings_IsDown(t *testing.T) {
	kb := DefaultKeyBindings()

	tests := []struct {
		name string
		msg  tea.KeyMsg
		want bool
	}{
		{
			name: "DownArrow",
			msg:  tea.KeyMsg{Type: tea.KeyDown},
			want: true,
		},
		{
			name: "J",
			msg:  tea.KeyMsg{Type: tea.KeyRune, Rune: 'j'},
			want: true,
		},
		{
			name: "Up",
			msg:  tea.KeyMsg{Type: tea.KeyUp},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := kb.IsDown(tt.msg)
			if got != tt.want {
				t.Errorf("IsDown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyBindings_IsHome(t *testing.T) {
	kb := DefaultKeyBindings()

	tests := []struct {
		name string
		msg  tea.KeyMsg
		want bool
	}{
		{
			name: "Home",
			msg:  tea.KeyMsg{Type: tea.KeyHome},
			want: true,
		},
		{
			name: "G",
			msg:  tea.KeyMsg{Type: tea.KeyRune, Rune: 'g'},
			want: true,
		},
		{
			name: "Other",
			msg:  tea.KeyMsg{Type: tea.KeyRune, Rune: 'x'},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := kb.IsHome(tt.msg)
			if got != tt.want {
				t.Errorf("IsHome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyBindings_IsEnd(t *testing.T) {
	kb := DefaultKeyBindings()

	tests := []struct {
		name string
		msg  tea.KeyMsg
		want bool
	}{
		{
			name: "End",
			msg:  tea.KeyMsg{Type: tea.KeyEnd},
			want: true,
		},
		{
			name: "ShiftG",
			msg:  tea.KeyMsg{Type: tea.KeyRune, Rune: 'G'},
			want: true,
		},
		{
			name: "Other",
			msg:  tea.KeyMsg{Type: tea.KeyRune, Rune: 'x'},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := kb.IsEnd(tt.msg)
			if got != tt.want {
				t.Errorf("IsEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyBindings_IsSort(t *testing.T) {
	kb := DefaultKeyBindings()

	if !kb.IsSort(tea.KeyMsg{Type: tea.KeyRune, Rune: 's'}) {
		t.Errorf("IsSort('s') should be true")
	}

	if !kb.IsSort(tea.KeyMsg{Type: tea.KeyEnter}) {
		t.Errorf("IsSort('enter') should be true")
	}
}

func TestKeyBindings_IsClearSort(t *testing.T) {
	kb := DefaultKeyBindings()

	if !kb.IsClearSort(tea.KeyMsg{Type: tea.KeyRune, Rune: 'c'}) {
		t.Errorf("IsClearSort('c') should be true")
	}
}

func TestKeyBindings_Custom(t *testing.T) {
	// Test custom key bindings
	kb := KeyBindings{
		Up:   []string{"w"},
		Down: []string{"s"},
		Home: []string{"q"},
		End:  []string{"e"},
	}

	if !kb.IsUp(tea.KeyMsg{Type: tea.KeyRune, Rune: 'w'}) {
		t.Errorf("Custom up binding 'w' should work")
	}

	if !kb.IsDown(tea.KeyMsg{Type: tea.KeyRune, Rune: 's'}) {
		t.Errorf("Custom down binding 's' should work")
	}
}
