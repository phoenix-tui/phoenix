package ansi_test

import (
	"testing"

	"github.com/phoenix-tui/phoenix/tea/internal/domain/model"
	"github.com/phoenix-tui/phoenix/tea/internal/infrastructure/ansi"
)

func TestParser_ParseKey_RegularKeys(t *testing.T) {
	p := ansi.NewParser()

	tests := []struct {
		name  string
		input []byte
		want  model.KeyMsg
		ok    bool
	}{
		{
			name:  "letter a",
			input: []byte{'a'},
			want:  model.KeyMsg{Type: model.KeyRune, Rune: 'a'},
			ok:    true,
		},
		{
			name:  "letter A",
			input: []byte{'A'},
			want:  model.KeyMsg{Type: model.KeyRune, Rune: 'A'},
			ok:    true,
		},
		{
			name:  "letter z",
			input: []byte{'z'},
			want:  model.KeyMsg{Type: model.KeyRune, Rune: 'z'},
			ok:    true,
		},
		{
			name:  "digit 0",
			input: []byte{'0'},
			want:  model.KeyMsg{Type: model.KeyRune, Rune: '0'},
			ok:    true,
		},
		{
			name:  "digit 9",
			input: []byte{'9'},
			want:  model.KeyMsg{Type: model.KeyRune, Rune: '9'},
			ok:    true,
		},
		{
			name:  "space",
			input: []byte{0x20},
			want:  model.KeyMsg{Type: model.KeySpace},
			ok:    true,
		},
		{
			name:  "enter CR",
			input: []byte{0x0D},
			want:  model.KeyMsg{Type: model.KeyEnter},
			ok:    true,
		},
		{
			name:  "enter LF",
			input: []byte{0x0A},
			want:  model.KeyMsg{Type: model.KeyEnter},
			ok:    true,
		},
		{
			name:  "backspace DEL",
			input: []byte{0x7F},
			want:  model.KeyMsg{Type: model.KeyBackspace},
			ok:    true,
		},
		{
			name:  "backspace BS",
			input: []byte{0x08},
			want:  model.KeyMsg{Type: model.KeyBackspace},
			ok:    true,
		},
		{
			name:  "tab",
			input: []byte{0x09},
			want:  model.KeyMsg{Type: model.KeyTab},
			ok:    true,
		},
		{
			name:  "esc",
			input: []byte{0x1B},
			want:  model.KeyMsg{Type: model.KeyEsc},
			ok:    true,
		},
		{
			name:  "empty input",
			input: []byte{},
			want:  model.KeyMsg{},
			ok:    false,
		},
		{
			name:  "non-printable",
			input: []byte{0x01}, // Ctrl+A will be parsed though
			want:  model.KeyMsg{Type: model.KeyRune, Rune: 'a', Ctrl: true},
			ok:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := p.ParseKey(tt.input)

			if ok != tt.ok {
				t.Errorf("ok = %v, want %v", ok, tt.ok)
			}

			if !ok {
				return
			}

			if got.Type != tt.want.Type {
				t.Errorf("Type = %v, want %v", got.Type, tt.want.Type)
			}

			if got.Type == model.KeyRune && got.Rune != tt.want.Rune {
				t.Errorf("Rune = %c, want %c", got.Rune, tt.want.Rune)
			}

			if got.Ctrl != tt.want.Ctrl {
				t.Errorf("Ctrl = %v, want %v", got.Ctrl, tt.want.Ctrl)
			}
		})
	}
}

func TestParser_ParseKey_ArrowKeys(t *testing.T) {
	p := ansi.NewParser()

	tests := []struct {
		name  string
		input []byte
		want  model.KeyType
	}{
		{"up", []byte{0x1B, '[', 'A'}, model.KeyUp},
		{"down", []byte{0x1B, '[', 'B'}, model.KeyDown},
		{"right", []byte{0x1B, '[', 'C'}, model.KeyRight},
		{"left", []byte{0x1B, '[', 'D'}, model.KeyLeft},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := p.ParseKey(tt.input)

			if !ok {
				t.Error("should parse arrow key")
			}

			if got.Type != tt.want {
				t.Errorf("Type = %v, want %v", got.Type, tt.want)
			}
		})
	}
}

func TestParser_ParseKey_CtrlKeys(t *testing.T) {
	p := ansi.NewParser()

	tests := []struct {
		name string
		byte byte
		want rune
	}{
		{"Ctrl+A", 0x01, 'a'},
		{"Ctrl+B", 0x02, 'b'},
		{"Ctrl+C", 0x03, 'c'},
		{"Ctrl+D", 0x04, 'd'},
		{"Ctrl+Z", 0x1A, 'z'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := p.ParseKey([]byte{tt.byte})

			if !ok {
				t.Errorf("should parse %s", tt.name)
			}

			if got.Type != model.KeyRune {
				t.Errorf("Type = %v, want KeyRune", got.Type)
			}

			if got.Rune != tt.want {
				t.Errorf("Rune = %c, want %c", got.Rune, tt.want)
			}

			if !got.Ctrl {
				t.Error("Ctrl should be true")
			}
		})
	}
}

func TestParser_ParseKey_FunctionKeys(t *testing.T) {
	p := ansi.NewParser()

	tests := []struct {
		name  string
		input []byte
		want  model.KeyType
	}{
		{"F1", []byte{0x1B, 'O', 'P'}, model.KeyF1},
		{"F2", []byte{0x1B, 'O', 'Q'}, model.KeyF2},
		{"F3", []byte{0x1B, 'O', 'R'}, model.KeyF3},
		{"F4", []byte{0x1B, 'O', 'S'}, model.KeyF4},
		{"F5", []byte{0x1B, '[', '1', '5', '~'}, model.KeyF5},
		{"F6", []byte{0x1B, '[', '1', '7', '~'}, model.KeyF6},
		{"F7", []byte{0x1B, '[', '1', '8', '~'}, model.KeyF7},
		{"F8", []byte{0x1B, '[', '1', '9', '~'}, model.KeyF8},
		{"F9", []byte{0x1B, '[', '2', '0', '~'}, model.KeyF9},
		{"F10", []byte{0x1B, '[', '2', '1', '~'}, model.KeyF10},
		{"F11", []byte{0x1B, '[', '2', '3', '~'}, model.KeyF11},
		{"F12", []byte{0x1B, '[', '2', '4', '~'}, model.KeyF12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := p.ParseKey(tt.input)

			if !ok {
				t.Errorf("should parse %s", tt.name)
			}

			if got.Type != tt.want {
				t.Errorf("Type = %v, want %v", got.Type, tt.want)
			}
		})
	}
}

func TestParser_ParseKey_SpecialKeys(t *testing.T) {
	p := ansi.NewParser()

	tests := []struct {
		name  string
		input []byte
		want  model.KeyType
	}{
		// ESC[N~ format
		{"Home (ESC[1~)", []byte{0x1B, '[', '1', '~'}, model.KeyHome},
		{"Insert (ESC[2~)", []byte{0x1B, '[', '2', '~'}, model.KeyInsert},
		{"Delete (ESC[3~)", []byte{0x1B, '[', '3', '~'}, model.KeyDelete}, // ← CRITICAL!
		{"End (ESC[4~)", []byte{0x1B, '[', '4', '~'}, model.KeyEnd},
		{"PageUp (ESC[5~)", []byte{0x1B, '[', '5', '~'}, model.KeyPgUp},
		{"PageDown (ESC[6~)", []byte{0x1B, '[', '6', '~'}, model.KeyPgDown},

		// Alternative Home/End sequences
		{"Home (ESC[H)", []byte{0x1B, '[', 'H'}, model.KeyHome},
		{"End (ESC[F)", []byte{0x1B, '[', 'F'}, model.KeyEnd},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := p.ParseKey(tt.input)

			if !ok {
				t.Errorf("should parse %s", tt.name)
			}

			if got.Type != tt.want {
				t.Errorf("Type = %v, want %v", got.Type, tt.want)
			}
		})
	}
}

func TestParser_ParseKey_InvalidSequences(t *testing.T) {
	p := ansi.NewParser()

	tests := []struct {
		name  string
		input []byte
	}{
		{"invalid ESC sequence", []byte{0x1B, 'X'}},
		{"incomplete arrow", []byte{0x1B, '['}},
		{"unknown function key", []byte{0x1B, '[', '9', '9', '~'}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := p.ParseKey(tt.input)

			if ok {
				t.Error("should not parse invalid sequence")
			}
		})
	}
}
