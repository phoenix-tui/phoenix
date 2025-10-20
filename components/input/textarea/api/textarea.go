// Package api provides the public API for the textarea component.
package api

import (
	"github.com/phoenix-tui/phoenix/components/input/textarea/domain/model"
	"github.com/phoenix-tui/phoenix/components/input/textarea/infrastructure/keybindings"
	"github.com/phoenix-tui/phoenix/components/input/textarea/infrastructure/renderer"
	"github.com/phoenix-tui/phoenix/tea/api"
)

// CursorPos represents a cursor position in the text buffer.
// This is used for cursor movement validation and observation.
type CursorPos struct {
	Row int // Line number (0-based)
	Col int // Column number (0-based, rune offset)
}

// MovementValidator validates cursor movements.
// Return true to allow movement, false to block it.
type MovementValidator func(from, to CursorPos) bool

// CursorMovedHandler is called after successful cursor movement.
// This is an observer pattern - cannot block movement.
type CursorMovedHandler func(from, to CursorPos)

// BoundaryHitHandler provides feedback when movement is blocked.
type BoundaryHitHandler func(attemptedPos CursorPos, reason string)

// TextArea is the public API for multiline text editing.
// This implements the Elm Architecture (Model-View-Update) pattern.
type TextArea struct {
	model       *model.TextArea
	keybindings KeybindingMode
	renderer    *renderer.TextAreaRenderer
}

// KeybindingMode defines keybinding style.
type KeybindingMode int

const (
	// KeybindingsDefault uses Emacs-style keybindings (default).
	KeybindingsDefault KeybindingMode = iota
	// KeybindingsEmacs uses Emacs-style keybindings explicitly.
	KeybindingsEmacs
	// KeybindingsVi uses Vi-style keybindings (future).
	KeybindingsVi
)

// New creates new TextArea with default settings.
func New() TextArea {
	return TextArea{
		model:       model.NewTextArea(),
		keybindings: KeybindingsEmacs, // Default to Emacs
		renderer:    renderer.NewTextAreaRenderer(),
	}
}

// Configuration Methods (Fluent Builder Pattern)

// Width sets display width.
func (t TextArea) Width(width int) TextArea {
	t.model = t.model.WithSize(width, t.model.Height())
	return t
}

// Height sets display height.
func (t TextArea) Height(height int) TextArea {
	t.model = t.model.WithSize(t.model.Width(), height)
	return t
}

// Size sets both width and height.
func (t TextArea) Size(width, height int) TextArea {
	t.model = t.model.WithSize(width, height)
	return t
}

// MaxLines sets maximum line limit (0 = unlimited).
func (t TextArea) MaxLines(maxVal int) TextArea {
	t.model = t.model.WithMaxLines(maxVal)
	return t
}

// MaxChars sets maximum character limit (0 = unlimited).
func (t TextArea) MaxChars(maxVal int) TextArea {
	t.model = t.model.WithMaxChars(maxVal)
	return t
}

// Placeholder sets placeholder text.
func (t TextArea) Placeholder(text string) TextArea {
	t.model = t.model.WithPlaceholder(text)
	return t
}

// Wrap enables/disables word wrap.
func (t TextArea) Wrap(wrap bool) TextArea {
	t.model = t.model.WithWrap(wrap)
	return t
}

// ReadOnly enables/disables read-only mode.
func (t TextArea) ReadOnly(readOnly bool) TextArea {
	t.model = t.model.WithReadOnly(readOnly)
	return t
}

// ShowLineNumbers enables/disables line numbers.
func (t TextArea) ShowLineNumbers(show bool) TextArea {
	t.model = t.model.WithLineNumbers(show)
	return t
}

// ShowCursor enables/disables Phoenix cursor rendering.
// true: Phoenix renders █ cursor (default)
// false: Use terminal cursor (ANSI positioning) - for shell applications.
func (t TextArea) ShowCursor(show bool) TextArea {
	t.model = t.model.WithShowCursor(show)
	return t
}

// Keybindings sets keybinding mode.
func (t TextArea) Keybindings(mode KeybindingMode) TextArea {
	t.keybindings = mode
	return t
}

// State Access Methods (Public Getters - CRITICAL for integration!)

// Value returns all text as single string.
func (t TextArea) Value() string {
	return t.model.Value()
}

// SetValue replaces all text.
func (t TextArea) SetValue(text string) TextArea {
	t.model = t.model.WithBuffer(model.NewBufferFromString(text))
	return t
}

// MoveCursorToEnd moves cursor to end of text.
func (t TextArea) MoveCursorToEnd() TextArea {
	t.model = t.model.MoveCursorToEnd()
	return t
}

// Lines returns all lines.
func (t TextArea) Lines() []string {
	return t.model.Lines()
}

// CursorPosition returns current cursor position (row, col).
func (t TextArea) CursorPosition() (row, col int) {
	return t.model.CursorPosition()
}

// SetCursorPosition sets cursor to specific position with bounds checking.
// Position is clamped to valid range (0 to buffer bounds).
// Returns new instance (immutable).
//
// Example - Jump to specific line/column:
//
//	ta := textarea.New().SetValue("line1\nline2\nline3")
//	ta = ta.SetCursorPosition(1, 2) // Row 1, Col 2.
//	row, col := ta.CursorPosition()
//	// row = 1, col = 2
//
// Example - Clamps to valid bounds:
//
//	ta := textarea.New().SetValue("short")
//	ta = ta.SetCursorPosition(10, 100) // Out of bounds.
//	row, col := ta.CursorPosition()
//	// row = 0, col = 5 (clamped to end of line)
func (t TextArea) SetCursorPosition(row, col int) TextArea {
	t.model = t.model.SetCursorPosition(row, col)
	return t
}

// OnMovement sets a validator that is called BEFORE cursor movements.
// Return false to block the movement.
// This is useful for implementing boundary protection (e.g., shell prompts).
//
// Example - Shell boundary protection:
//
//	ta := textarea.New().
//	    SetValue("> ").
//	    SetCursorPosition(0, 2).
//	    OnMovement(func(from, to textarea.CursorPos) bool {.
//	        // Don't allow cursor before the prompt ("> ")
//	        if to.Row == 0 && to.Col < 2 {.
//	            return false // Block movement.
//	        }
//	        return true // Allow movement.
//	    })
func (t TextArea) OnMovement(validator MovementValidator) TextArea {
	// Convert API CursorPos to domain model CursorPos.
	if validator != nil {
		domainValidator := func(from, to model.CursorPos) bool {
			apiFrom := CursorPos{Row: from.Row, Col: from.Col}
			apiTo := CursorPos{Row: to.Row, Col: to.Col}
			return validator(apiFrom, apiTo)
		}
		// Use reflection to access package-private method.
		// We need to call withMovementValidator which is package-private.
		// For now, we'll directly set it since we're in the same module.
		t.model = setMovementValidator(t.model, domainValidator)
	}
	return t
}

// OnCursorMoved sets an observer that is called AFTER successful cursor movement.
// This cannot block movement - use OnMovement() for validation.
// This is useful for updating UI (e.g., syntax highlighting) when cursor moves.
//
// Example - Refresh syntax highlighting on row change:
//
//	ta := textarea.New().
//	    OnCursorMoved(func(from, to textarea.CursorPos) {.
//	        if from.Row != to.Row {.
//	            // Cursor moved to different line.
//	            refreshSyntaxHighlight(to.Row)
//	        }
//	    })
func (t TextArea) OnCursorMoved(handler CursorMovedHandler) TextArea {
	// Convert API CursorPos to domain model CursorPos.
	if handler != nil {
		domainHandler := func(from, to model.CursorPos) {
			apiFrom := CursorPos{Row: from.Row, Col: from.Col}
			apiTo := CursorPos{Row: to.Row, Col: to.Col}
			handler(apiFrom, apiTo)
		}
		t.model = setCursorMovedHandler(t.model, domainHandler)
	}
	return t
}

// OnBoundaryHit sets a handler that is called when cursor movement is blocked.
// This provides feedback to the user when they try to move beyond allowed boundaries.
// Useful for accessibility and user experience.
//
// Example - Visual feedback for blocked movement:
//
//	ta := textarea.New().
//	    OnMovement(func(from, to textarea.CursorPos) bool {.
//	        return to.Row >= 0 && to.Col >= 0 // Block negative positions.
//	    }).
//	    OnBoundaryHit(func(attemptedPos textarea.CursorPos, reason string) {.
//	        // Flash screen, beep, or show message.
//	        fmt.Println("Cannot move to", attemptedPos, ":", reason)
//	    })
func (t TextArea) OnBoundaryHit(handler BoundaryHitHandler) TextArea {
	// Convert API CursorPos to domain model CursorPos.
	if handler != nil {
		domainHandler := func(attemptedPos model.CursorPos, reason string) {
			apiPos := CursorPos{Row: attemptedPos.Row, Col: attemptedPos.Col}
			handler(apiPos, reason)
		}
		t.model = setBoundaryHitHandler(t.model, domainHandler)
	}
	return t
}

// Helper functions to access package-private setters from domain model.
func setMovementValidator(ta *model.TextArea, validator func(from, to model.CursorPos) bool) *model.TextArea {
	return ta.WithMovementValidator(validator)
}

func setCursorMovedHandler(ta *model.TextArea, handler func(from, to model.CursorPos)) *model.TextArea {
	return ta.WithCursorMovedHandler(handler)
}

func setBoundaryHitHandler(ta *model.TextArea, handler func(attemptedPos model.CursorPos, reason string)) *model.TextArea {
	return ta.WithBoundaryHitHandler(handler)
}

// ContentParts returns text before/at/after cursor (for syntax highlighting).
func (t TextArea) ContentParts() (before, at, after string) {
	return t.model.ContentParts()
}

// CurrentLine returns text of current line.
func (t TextArea) CurrentLine() string {
	return t.model.CurrentLine()
}

// LineCount returns number of lines.
func (t TextArea) LineCount() int {
	return t.model.LineCount()
}

// IsEmpty returns true if buffer has no content.
func (t TextArea) IsEmpty() bool {
	return t.model.IsEmpty()
}

// HasSelection returns true if there is active selection.
func (t TextArea) HasSelection() bool {
	return t.model.HasSelection()
}

// SelectedText returns selected text (empty if no selection).
func (t TextArea) SelectedText() string {
	return t.model.SelectedText()
}

// Bubbletea Integration (Elm Architecture)

// Init initializes the component.
func (t TextArea) Init() api.Cmd {
	return nil
}

// Update handles messages and returns updated component.
func (t TextArea) Update(msg api.Msg) (TextArea, api.Cmd) {
	switch msg := msg.(type) {
	case api.KeyMsg:
		// Delegate to keybindings handler.
		switch t.keybindings {
		case KeybindingsEmacs, KeybindingsDefault:
			handler := keybindings.NewEmacsKeybindings()
			newModel, cmd := handler.Handle(msg, t.model)
			t.model = newModel
			return t, cmd

		case KeybindingsVi:
			// Future: Vi keybindings.
			return t, nil

		default:
			return t, nil
		}

	case api.WindowSizeMsg:
		// Handle window resize.
		// For now, we don't auto-resize the textarea.
		// User should explicitly set size if needed.
		return t, nil
	}

	return t, nil
}

// View renders the component.
func (t TextArea) View() string {
	return t.renderer.Render(t.model)
}
