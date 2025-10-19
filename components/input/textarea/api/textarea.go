// Package api provides the public API for the textarea component.
package api

import (
	"github.com/phoenix-tui/phoenix/components/input/textarea/domain/model"
	"github.com/phoenix-tui/phoenix/components/input/textarea/infrastructure/keybindings"
	"github.com/phoenix-tui/phoenix/components/input/textarea/infrastructure/renderer"
	"github.com/phoenix-tui/phoenix/tea/api"
)

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
func (t TextArea) MaxLines(max int) TextArea {
	t.model = t.model.WithMaxLines(max)
	return t
}

// MaxChars sets maximum character limit (0 = unlimited).
func (t TextArea) MaxChars(max int) TextArea {
	t.model = t.model.WithMaxChars(max)
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
// false: Use terminal cursor (ANSI positioning) - for shell applications
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
		// Delegate to keybindings handler
		switch t.keybindings {
		case KeybindingsEmacs, KeybindingsDefault:
			handler := keybindings.NewEmacsKeybindings()
			newModel, cmd := handler.Handle(msg, t.model)
			t.model = newModel
			return t, cmd

		case KeybindingsVi:
			// Future: Vi keybindings
			return t, nil

		default:
			return t, nil
		}

	case api.WindowSizeMsg:
		// Handle window resize
		// For now, we don't auto-resize the textarea
		// User should explicitly set size if needed
		return t, nil
	}

	return t, nil
}

// View renders the component.
func (t TextArea) View() string {
	return t.renderer.Render(t.model)
}
