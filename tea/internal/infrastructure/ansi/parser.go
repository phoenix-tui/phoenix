// Package ansi provides ANSI escape sequence parsing for keyboard input.
package ansi

import (
	"bytes"

	"github.com/phoenix-tui/phoenix/tea/internal/domain/model"
)

// Parser parses ANSI escape sequences into messages.
type Parser struct{}

// NewParser creates a new ANSI parser.
func NewParser() *Parser {
	return &Parser{}
}

// ParseKey parses a byte sequence into a KeyMsg.
// Returns KeyMsg and true if parsed, or zero KeyMsg and false if not recognized.
//
// Supports:
// - Regular ASCII keys (a-z, A-Z, 0-9, space, etc.)
// - Enter, Backspace, Tab, Esc.
// - Arrow keys (ESC [ A/B/C/D).
// - Function keys F1-F12 (basic sequences).
// - Ctrl combinations (Ctrl+A through Ctrl+Z).
//
//nolint:gocognit,gocyclo,cyclop,funlen,nestif // ANSI parsing requires sequential checks for all key types
func (p *Parser) ParseKey(data []byte) (model.KeyMsg, bool) {
	if len(data) == 0 {
		return model.KeyMsg{}, false
	}

	// Single byte - regular key or Ctrl combination
	if len(data) == 1 {
		b := data[0]

		// Special keys FIRST (before Ctrl, as they overlap!)
		switch b {
		case 0x0D, 0x0A: // Enter (CR or LF) - also Ctrl+M, Ctrl+J
			return model.KeyMsg{Type: model.KeyEnter}, true
		case 0x7F: // Backspace DEL
			return model.KeyMsg{Type: model.KeyBackspace}, true
		case 0x08: // Backspace BS (also Ctrl+H)
			return model.KeyMsg{Type: model.KeyBackspace}, true
		case 0x09: // Tab (also Ctrl+I)
			return model.KeyMsg{Type: model.KeyTab}, true
		case 0x1B: // Esc
			return model.KeyMsg{Type: model.KeyEsc}, true
		case 0x20: // Space
			return model.KeyMsg{Type: model.KeySpace}, true
		}

		// Ctrl+A through Ctrl+Z (0x01 - 0x1A) - AFTER special keys
		// Skip 0x08 (Ctrl+H/Backspace), 0x09 (Ctrl+I/Tab), 0x0A (Ctrl+J/Enter), 0x0D (Ctrl+M/Enter)
		if b >= 1 && b <= 26 && b != 0x08 && b != 0x09 && b != 0x0A && b != 0x0D {
			return model.KeyMsg{
				Type: model.KeyRune,
				Rune: rune('a' + b - 1),
				Ctrl: true,
			}, true
		}

		// Regular printable character
		if b >= 32 && b <= 126 {
			return model.KeyMsg{
				Type: model.KeyRune,
				Rune: rune(b),
			}, true
		}

		return model.KeyMsg{}, false
	}

	// Multi-byte - ANSI escape sequences
	if data[0] == 0x1B { // ESC
		// Arrow keys: ESC [ A/B/C/D
		if len(data) == 3 && data[1] == '[' {
			switch data[2] {
			case 'A':
				return model.KeyMsg{Type: model.KeyUp}, true
			case 'B':
				return model.KeyMsg{Type: model.KeyDown}, true
			case 'C':
				return model.KeyMsg{Type: model.KeyRight}, true
			case 'D':
				return model.KeyMsg{Type: model.KeyLeft}, true
			}
		}

		// Function keys: ESC O P/Q/R/S (F1-F4)
		if len(data) == 3 && data[1] == 'O' {
			switch data[2] {
			case 'P':
				return model.KeyMsg{Type: model.KeyF1}, true
			case 'Q':
				return model.KeyMsg{Type: model.KeyF2}, true
			case 'R':
				return model.KeyMsg{Type: model.KeyF3}, true
			case 'S':
				return model.KeyMsg{Type: model.KeyF4}, true
			}
		}

		// Special keys with ~ terminator: ESC [ N ~
		// Home (1~), Insert (2~), Delete (3~), End (4~), PageUp (5~), PageDown (6~)
		if len(data) >= 4 && data[1] == '[' && data[len(data)-1] == '~' {
			// Extract the number before ~
			switch {
			case bytes.Equal(data, []byte{0x1B, '[', '1', '~'}):
				return model.KeyMsg{Type: model.KeyHome}, true
			case bytes.Equal(data, []byte{0x1B, '[', '2', '~'}):
				return model.KeyMsg{Type: model.KeyInsert}, true
			case bytes.Equal(data, []byte{0x1B, '[', '3', '~'}):
				return model.KeyMsg{Type: model.KeyDelete}, true // ← CRITICAL FIX!
			case bytes.Equal(data, []byte{0x1B, '[', '4', '~'}):
				return model.KeyMsg{Type: model.KeyEnd}, true
			case bytes.Equal(data, []byte{0x1B, '[', '5', '~'}):
				return model.KeyMsg{Type: model.KeyPgUp}, true
			case bytes.Equal(data, []byte{0x1B, '[', '6', '~'}):
				return model.KeyMsg{Type: model.KeyPgDown}, true

			// Function keys: ESC [ 1 5 ~ (F5), ESC [ 1 7 ~ (F6), etc.
			case bytes.Equal(data, []byte{0x1B, '[', '1', '5', '~'}):
				return model.KeyMsg{Type: model.KeyF5}, true
			case bytes.Equal(data, []byte{0x1B, '[', '1', '7', '~'}):
				return model.KeyMsg{Type: model.KeyF6}, true
			case bytes.Equal(data, []byte{0x1B, '[', '1', '8', '~'}):
				return model.KeyMsg{Type: model.KeyF7}, true
			case bytes.Equal(data, []byte{0x1B, '[', '1', '9', '~'}):
				return model.KeyMsg{Type: model.KeyF8}, true
			case bytes.Equal(data, []byte{0x1B, '[', '2', '0', '~'}):
				return model.KeyMsg{Type: model.KeyF9}, true
			case bytes.Equal(data, []byte{0x1B, '[', '2', '1', '~'}):
				return model.KeyMsg{Type: model.KeyF10}, true
			case bytes.Equal(data, []byte{0x1B, '[', '2', '3', '~'}):
				return model.KeyMsg{Type: model.KeyF11}, true
			case bytes.Equal(data, []byte{0x1B, '[', '2', '4', '~'}):
				return model.KeyMsg{Type: model.KeyF12}, true
			}
		}

		// Alternative Home/End sequences: ESC [ H / ESC [ F
		if len(data) == 3 && data[1] == '[' {
			switch data[2] {
			case 'H':
				return model.KeyMsg{Type: model.KeyHome}, true
			case 'F':
				return model.KeyMsg{Type: model.KeyEnd}, true
			}
		}
	}

	return model.KeyMsg{}, false
}
