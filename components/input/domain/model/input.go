// Package model contains rich domain models for text input components.
package model

import (
	"github.com/phoenix-tui/phoenix/components/input/domain/service"
	"github.com/phoenix-tui/phoenix/components/input/domain/value"
)

// TextInput is the aggregate root for the text input component.
// It manages content, cursor position, selection, and validation.
// All operations are immutable - they return new instances.
type TextInput struct {
	content        string                 // Current text content
	cursor         *value.Cursor          // Cursor position (grapheme-aware)
	selection      *value.Selection       // Selection range (nil if no selection)
	validator      service.ValidationFunc // Validation hook (nil if no validation)
	width          int                    // Visible width (for scrolling)
	scrollOffset   int                    // Horizontal scroll offset
	placeholder    string                 // Placeholder text (when empty)
	focused        bool                   // Focus state
	showCursor     bool                   // Whether to render cursor (default: true)
	cursorMovement *service.CursorMovementService
	validationSvc  *service.ValidationService
}

// New creates a new TextInput with the specified width.
func New(width int) *TextInput {
	if width < 1 {
		width = 1
	}

	return &TextInput{
		content:        "",
		cursor:         value.NewCursor(0),
		selection:      nil,
		validator:      nil,
		width:          width,
		scrollOffset:   0,
		placeholder:    "",
		focused:        false,
		showCursor:     true, // Default: show cursor
		cursorMovement: service.NewCursorMovementService(),
		validationSvc:  service.NewValidationService(),
	}
}

// Content returns the current text content.
func (t *TextInput) Content() string {
	return t.content
}

// CursorPosition returns the current cursor position (grapheme offset).
func (t *TextInput) CursorPosition() int {
	return t.cursor.Offset()
}

// Width returns the visible width.
func (t *TextInput) Width() int {
	return t.width
}

// ScrollOffset returns the current horizontal scroll offset.
func (t *TextInput) ScrollOffset() int {
	return t.scrollOffset
}

// Placeholder returns the placeholder text.
func (t *TextInput) Placeholder() string {
	return t.placeholder
}

// Focused returns the focus state.
func (t *TextInput) Focused() bool {
	return t.focused
}

// ShowCursor returns whether cursor rendering is enabled.
func (t *TextInput) ShowCursor() bool {
	return t.showCursor
}

// HasSelection returns true if there's an active selection.
func (t *TextInput) HasSelection() bool {
	return t.selection != nil && !t.selection.IsEmpty()
}

// Selection returns the current selection (nil if no selection).
func (t *TextInput) Selection() *value.Selection {
	if t.selection == nil {
		return nil
	}
	return t.selection.Clone()
}

// ContentParts splits content around cursor: (before, at, after).
// This is a KEY DIFFERENTIATOR - enables apps to customize cursor rendering.
func (t *TextInput) ContentParts() (before, at, after string) {
	return t.cursorMovement.SplitAtCursor(t.content, t.cursor.Offset())
}

// WithContent sets the content (immutable).
func (t TextInput) WithContent(content string) TextInput {
	t.content = content

	// Clamp cursor to new content length.
	maxOffset := t.cursorMovement.GraphemeCount(content)
	t.cursor = t.cursor.MoveTo(t.cursor.Offset(), maxOffset)

	// Clear selection if it's now invalid.
	if t.selection != nil {
		t.selection = t.selection.Clamp(maxOffset)
		if t.selection.IsEmpty() {
			t.selection = nil
		}
	}

	return t
}

// SetContent sets both content and cursor position atomically.
// This is a KEY DIFFERENTIATOR - prevents race conditions.
func (t TextInput) SetContent(content string, cursorPos int) TextInput {
	t.content = content

	// Clamp cursor to content length.
	maxOffset := t.cursorMovement.GraphemeCount(content)
	if cursorPos < 0 {
		cursorPos = 0
	}
	if cursorPos > maxOffset {
		cursorPos = maxOffset
	}
	t.cursor = value.NewCursor(cursorPos)

	// Clear selection.
	t.selection = nil

	return t
}

// WithCursor sets the cursor position (immutable).
func (t TextInput) WithCursor(pos int) TextInput {
	maxOffset := t.cursorMovement.GraphemeCount(t.content)
	t.cursor = value.NewCursor(pos).MoveTo(pos, maxOffset)
	return t
}

// WithSelection sets the selection range (immutable).
func (t TextInput) WithSelection(start, end int) TextInput {
	maxOffset := t.cursorMovement.GraphemeCount(t.content)
	t.selection = value.NewSelection(start, end).Clamp(maxOffset)

	// Move cursor to end of selection.
	t.cursor = value.NewCursor(t.selection.End())

	return t
}

// WithValidator sets the validation function (immutable).
func (t TextInput) WithValidator(validator service.ValidationFunc) TextInput {
	t.validator = validator
	return t
}

// WithWidth sets the visible width (immutable).
func (t TextInput) WithWidth(width int) TextInput {
	if width < 1 {
		width = 1
	}
	t.width = width
	return t
}

// WithPlaceholder sets the placeholder text (immutable).
func (t TextInput) WithPlaceholder(text string) TextInput {
	t.placeholder = text
	return t
}

// WithFocus sets the focus state (immutable).
func (t TextInput) WithFocus(focused bool) TextInput {
	t.focused = focused
	return t
}

// WithShowCursor sets whether cursor should be rendered (immutable).
// When false, applications can use the terminal's native cursor instead.
func (t TextInput) WithShowCursor(show bool) TextInput {
	t.showCursor = show
	return t
}

// MoveLeft moves cursor left by one grapheme (immutable).
func (t TextInput) MoveLeft() TextInput {
	newPos := t.cursorMovement.MoveLeft(t.content, t.cursor.Offset())
	t.cursor = value.NewCursor(newPos)
	t.selection = nil // Clear selection on cursor movement
	return t
}

// MoveRight moves cursor right by one grapheme (immutable).
func (t TextInput) MoveRight() TextInput {
	newPos := t.cursorMovement.MoveRight(t.content, t.cursor.Offset())
	t.cursor = value.NewCursor(newPos)
	t.selection = nil // Clear selection on cursor movement
	return t
}

// MoveHome moves cursor to start (immutable).
func (t TextInput) MoveHome() TextInput {
	t.cursor = value.NewCursor(0)
	t.selection = nil // Clear selection on cursor movement
	return t
}

// MoveEnd moves cursor to end (immutable).
func (t TextInput) MoveEnd() TextInput {
	maxOffset := t.cursorMovement.GraphemeCount(t.content)
	t.cursor = value.NewCursor(maxOffset)
	t.selection = nil // Clear selection on cursor movement
	return t
}

// InsertRune inserts a rune at cursor position (immutable).
func (t TextInput) InsertRune(r rune) TextInput {
	// Delete selection if present.
	if t.selection != nil && !t.selection.IsEmpty() {
		t = t.deleteSelection()
	}

	// Insert rune at cursor position.
	// SplitAtCursor gives us (before, at, after) where 'at' is the grapheme at cursor.
	// We want to insert BEFORE 'at', so: before + newRune + at + after.
	before, at, after := t.cursorMovement.SplitAtCursor(t.content, t.cursor.Offset())
	t.content = before + string(r) + at + after

	// Move cursor right.
	t.cursor = t.cursor.MoveBy(1, t.cursorMovement.GraphemeCount(t.content))

	return t
}

// DeleteBackward deletes grapheme before cursor (Backspace) (immutable).
func (t TextInput) DeleteBackward() TextInput {
	// If selection exists, delete it.
	if t.selection != nil && !t.selection.IsEmpty() {
		return t.deleteSelection()
	}

	// Can't delete if at start.
	if t.cursor.Offset() == 0 {
		return t
	}

	// Delete grapheme before cursor.
	cursorPos := t.cursor.Offset()
	before, at, after := t.cursorMovement.SplitAtCursor(t.content, cursorPos)

	// Remove last grapheme from before.
	gr := t.cursorMovement
	if before != "" {
		beforePos := gr.GraphemeCount(before) - 1
		beforeByteOffset := gr.GraphemeOffsetToByteOffset(before, beforePos)
		before = before[:beforeByteOffset]
	}

	// Reconstruct: trimmed before + at (grapheme at cursor) + after.
	t.content = before + at + after
	t.cursor = t.cursor.MoveBy(-1, t.cursorMovement.GraphemeCount(t.content))

	return t
}

// DeleteForward deletes grapheme after cursor (Delete) (immutable).
func (t TextInput) DeleteForward() TextInput {
	// If selection exists, delete it.
	if t.selection != nil && !t.selection.IsEmpty() {
		return t.deleteSelection()
	}

	// Can't delete if at end.
	maxOffset := t.cursorMovement.GraphemeCount(t.content)
	if t.cursor.Offset() >= maxOffset {
		return t
	}

	// Delete grapheme at cursor.
	before, _, after := t.cursorMovement.SplitAtCursor(t.content, t.cursor.Offset())
	t.content = before + after

	// Cursor stays in place.
	return t
}

// Clear removes all content (Ctrl-U) (immutable).
func (t TextInput) Clear() TextInput {
	t.content = ""
	t.cursor = value.NewCursor(0)
	t.selection = nil
	t.scrollOffset = 0
	return t
}

// SelectAll selects all content (Ctrl-A) (immutable).
func (t TextInput) SelectAll() TextInput {
	maxOffset := t.cursorMovement.GraphemeCount(t.content)
	t.selection = value.NewSelection(0, maxOffset)
	t.cursor = value.NewCursor(maxOffset)
	return t
}

// ClearSelection removes the selection without deleting content (immutable).
func (t TextInput) ClearSelection() TextInput {
	t.selection = nil
	return t
}

// Validate runs the validator on current content.
func (t *TextInput) Validate() error {
	return t.validationSvc.Validate(t.content, t.validator)
}

// IsValid returns true if content passes validation.
func (t *TextInput) IsValid() bool {
	return t.validationSvc.IsValid(t.content, t.validator)
}

// deleteSelection is a helper that deletes the selected text.
func (t TextInput) deleteSelection() TextInput {
	if t.selection == nil || t.selection.IsEmpty() {
		return t
	}

	// Split at selection boundaries.
	startByteOffset := t.cursorMovement.GraphemeOffsetToByteOffset(t.content, t.selection.Start())
	endByteOffset := t.cursorMovement.GraphemeOffsetToByteOffset(t.content, t.selection.End())

	newContent := t.content[:startByteOffset] + t.content[endByteOffset:]

	t.content = newContent
	t.cursor = value.NewCursor(t.selection.Start())
	t.selection = nil

	return t
}
