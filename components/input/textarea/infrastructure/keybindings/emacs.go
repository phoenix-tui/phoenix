// Package keybindings provides keybinding handlers for textarea.
package keybindings

import (
	"github.com/phoenix-tui/phoenix/components/input/textarea/domain/model"
	"github.com/phoenix-tui/phoenix/components/input/textarea/domain/service"
	"github.com/phoenix-tui/phoenix/tea/api"
)

// EmacsKeybindings implements Emacs-style keybindings.
// This is an infrastructure component that translates key events to domain operations.
type EmacsKeybindings struct {
	navigation *service.NavigationService
	editing    *service.EditingService
}

// NewEmacsKeybindings creates Emacs keybindings handler.
func NewEmacsKeybindings() *EmacsKeybindings {
	return &EmacsKeybindings{
		navigation: service.NewNavigationService(),
		editing:    service.NewEditingService(),
	}
}

// Handle processes key message and returns updated TextArea.
//
//nolint:gocognit,gocyclo,cyclop,funlen // keybindings require state machine logic
func (e *EmacsKeybindings) Handle(msg api.KeyMsg, ta *model.TextArea) (*model.TextArea, api.Cmd) {
	// Handle Ctrl key combinations.
	if msg.Ctrl {
		//nolint:gocritic // switch with single case is intentional for Emacs bindings structure
		switch msg.Type {
		case api.KeyRune:
			switch msg.Rune {
			// Navigation.
			case 'a', 'A':
				return e.navigation.MoveToLineStart(ta), nil

			case 'e', 'E':
				return e.navigation.MoveToLineEnd(ta), nil

			case 'p', 'P':
				return e.navigation.MoveUp(ta), nil

			case 'n', 'N':
				return e.navigation.MoveDown(ta), nil

			case 'f', 'F':
				return e.navigation.MoveRight(ta), nil

			case 'b', 'B':
				return e.navigation.MoveLeft(ta), nil

			// Editing.
			case 'k', 'K':
				return e.editing.KillLine(ta), nil

			case 'u', 'U':
				// Kill from start of line to cursor.
				ta = e.navigation.MoveToLineStart(ta)
				return e.editing.KillLine(ta), nil

			case 'w', 'W':
				return e.editing.KillWordBackward(ta), nil

			case 'y', 'Y':
				return e.editing.Yank(ta), nil

			case 'd', 'D':
				return e.editing.DeleteCharForward(ta), nil

			case 'h', 'H':
				return e.editing.DeleteCharBackward(ta), nil

			case 'm', 'M':
				return e.editing.InsertNewline(ta), nil
			}
		}
	}

	// Handle Alt key combinations.
	if msg.Alt {
		switch msg.Type {
		case api.KeyRune:
			switch msg.Rune {
			case 'f', 'F':
				return e.navigation.ForwardWord(ta), nil

			case 'b', 'B':
				return e.navigation.BackwardWord(ta), nil

			case '<':
				return e.navigation.MoveToBufferStart(ta), nil

			case '>':
				return e.navigation.MoveToBufferEnd(ta), nil

			case 'd', 'D':
				return e.editing.KillWord(ta), nil
			}

		case api.KeyBackspace:
			return e.editing.KillWordBackward(ta), nil
		}
	}

	// Handle special keys (without modifiers)
	if !msg.Ctrl && !msg.Alt {
		switch msg.Type {
		case api.KeyUp:
			return e.navigation.MoveUp(ta), nil

		case api.KeyDown:
			return e.navigation.MoveDown(ta), nil

		case api.KeyLeft:
			return e.navigation.MoveLeft(ta), nil

		case api.KeyRight:
			return e.navigation.MoveRight(ta), nil

		case api.KeyHome:
			return e.navigation.MoveToLineStart(ta), nil

		case api.KeyEnd:
			return e.navigation.MoveToLineEnd(ta), nil

		case api.KeyBackspace:
			return e.editing.DeleteCharBackward(ta), nil

		case api.KeyDelete:
			return e.editing.DeleteCharForward(ta), nil

		case api.KeyEnter:
			return e.editing.InsertNewline(ta), nil

		case api.KeySpace:
			// Insert space character (0x20 is parsed as KeySpace, not KeyRune)
			// CRITICAL FIX: Without this, spaces are ignored until next character.
			return e.editing.InsertChar(ta, ' '), nil

		case api.KeyRune:
			// Insert character.
			return e.editing.InsertChar(ta, msg.Rune), nil
		}
	}

	// Unhandled key.
	return ta, nil
}
