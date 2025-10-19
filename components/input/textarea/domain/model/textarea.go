// Package model provides rich domain models for textarea.
package model

import (
	"fmt"
)

// TextArea is the rich domain model for multiline text editing.
// This is the aggregate root that coordinates all textarea behavior.
// All operations return new instances (immutable).
type TextArea struct {
	// Buffer state
	buffer    *Buffer    // Lines of text
	cursor    *Cursor    // Current cursor position
	selection *Selection // Current selection (nil if none)
	killRing  *KillRing  // Kill ring for Emacs-style cut/paste

	// Display configuration
	width  int // Display width (for wrapping)
	height int // Display height (visible lines)

	// Scrolling state
	scrollRow int // First visible row (for viewport)
	scrollCol int // First visible column (horizontal scroll)

	// Behavior configuration
	maxLines    int    // Maximum number of lines (0 = unlimited)
	maxChars    int    // Maximum total characters (0 = unlimited)
	placeholder string // Placeholder text
	wrap        bool   // Word wrap (false = horizontal scroll)
	readOnly    bool   // Read-only mode

	// Appearance
	showLineNumbers bool // Show line numbers
	lineNumberWidth int  // Width of line number column
	showCursor      bool // Show cursor (true = Phoenix renders █, false = use terminal cursor)
}

// NewTextArea creates a new TextArea with default settings.
func NewTextArea() *TextArea {
	return &TextArea{
		buffer:          NewBuffer(),
		cursor:          NewCursor(0, 0),
		selection:       nil,
		killRing:        NewKillRing(10), // Keep last 10 kills
		width:           80,
		height:          24,
		maxLines:        0,     // Unlimited
		maxChars:        0,     // Unlimited
		wrap:            false, // No wrap by default
		readOnly:        false, // Editable
		showLineNumbers: false, // No line numbers
		lineNumberWidth: 0,
		showCursor:      true, // Show cursor by default
	}
}

// Configuration Methods (Fluent Builder Pattern)

// WithSize sets display dimensions (returns new instance).
func (t *TextArea) WithSize(width, height int) *TextArea {
	copy := t.copy()
	copy.width = width
	copy.height = height
	return copy
}

// WithMaxLines sets maximum line limit (0 = unlimited).
func (t *TextArea) WithMaxLines(max int) *TextArea {
	copy := t.copy()
	copy.maxLines = max
	return copy
}

// WithMaxChars sets maximum character limit (0 = unlimited).
func (t *TextArea) WithMaxChars(max int) *TextArea {
	copy := t.copy()
	copy.maxChars = max
	return copy
}

// WithWrap enables/disables word wrap.
func (t *TextArea) WithWrap(wrap bool) *TextArea {
	copy := t.copy()
	copy.wrap = wrap
	return copy
}

// WithPlaceholder sets placeholder text.
func (t *TextArea) WithPlaceholder(text string) *TextArea {
	copy := t.copy()
	copy.placeholder = text
	return copy
}

// WithReadOnly enables/disables read-only mode.
func (t *TextArea) WithReadOnly(readOnly bool) *TextArea {
	copy := t.copy()
	copy.readOnly = readOnly
	return copy
}

// WithLineNumbers enables/disables line numbers.
func (t *TextArea) WithLineNumbers(show bool) *TextArea {
	copy := t.copy()
	copy.showLineNumbers = show
	if show {
		copy.lineNumberWidth = len(fmt.Sprintf("%d", t.buffer.LineCount()))
	} else {
		copy.lineNumberWidth = 0
	}
	return copy
}

// WithShowCursor enables/disables Phoenix cursor rendering.
// true: Phoenix renders █ cursor (default)
// false: Use terminal cursor (ANSI positioning) - for shell applications
func (t *TextArea) WithShowCursor(show bool) *TextArea {
	copy := t.copy()
	copy.showCursor = show
	return copy
}

// WithBuffer replaces buffer (returns new instance).
func (t *TextArea) WithBuffer(buffer *Buffer) *TextArea {
	copy := t.copy()
	copy.buffer = buffer
	// Reset cursor to start
	copy.cursor = NewCursor(0, 0)
	copy.selection = nil
	return copy
}

// WithCursor sets cursor position (returns new instance).
// This is used by domain services for cursor movement.
func (t *TextArea) WithCursor(cursor *Cursor) *TextArea {
	return t.withCursor(cursor)
}

// Public Getters (for external integration - CRITICAL!)

// CursorPosition returns current cursor position (row, col).
func (t *TextArea) CursorPosition() (row, col int) {
	return t.cursor.Row(), t.cursor.Col()
}

// Lines returns all lines in buffer.
func (t *TextArea) Lines() []string {
	return t.buffer.Lines()
}

// Value returns all text as single string (lines joined with \n).
func (t *TextArea) Value() string {
	return t.buffer.String()
}

// LineCount returns number of lines.
func (t *TextArea) LineCount() int {
	return t.buffer.LineCount()
}

// CurrentLine returns text of current line.
func (t *TextArea) CurrentLine() string {
	return t.buffer.Line(t.cursor.Row())
}

// ContentParts returns text before/at/after cursor (for syntax highlighting).
func (t *TextArea) ContentParts() (before, at, after string) {
	line := t.CurrentLine()
	runes := []rune(line)
	col := t.cursor.Col()

	if col >= len(runes) {
		// Cursor at end of line
		return line, " ", ""
	}

	before = string(runes[:col])
	at = string(runes[col : col+1])
	after = string(runes[col+1:])

	return before, at, after
}

// HasSelection returns true if there is active selection.
func (t *TextArea) HasSelection() bool {
	return t.selection != nil
}

// SelectedText returns selected text (empty if no selection).
func (t *TextArea) SelectedText() string {
	if !t.HasSelection() {
		return ""
	}
	return t.buffer.TextInRange(t.selection.Range())
}

// IsEmpty returns true if buffer has no content.
func (t *TextArea) IsEmpty() bool {
	return t.buffer.IsEmpty()
}

// VisibleLines returns slice of visible lines (respecting scroll).
func (t *TextArea) VisibleLines() []string {
	start := t.scrollRow
	end := start + t.height
	if end > t.buffer.LineCount() {
		end = t.buffer.LineCount()
	}

	lines := t.buffer.Lines()
	if start >= len(lines) {
		return []string{}
	}

	return lines[start:end]
}

// Width returns display width.
func (t *TextArea) Width() int {
	return t.width
}

// Height returns display height.
func (t *TextArea) Height() int {
	return t.height
}

// Placeholder returns placeholder text.
func (t *TextArea) Placeholder() string {
	return t.placeholder
}

// IsReadOnly returns true if read-only mode is enabled.
func (t *TextArea) IsReadOnly() bool {
	return t.readOnly
}

// ShowLineNumbers returns true if line numbers are shown.
func (t *TextArea) ShowLineNumbers() bool {
	return t.showLineNumbers
}

// ShowCursor returns whether Phoenix cursor rendering is enabled.
func (t *TextArea) ShowCursor() bool {
	return t.showCursor
}

// MaxLines returns the maximum number of lines (0 = unlimited).
func (t *TextArea) MaxLines() int {
	return t.maxLines
}

// MoveCursorToEnd moves cursor to end of text.
func (t *TextArea) MoveCursorToEnd() *TextArea {
	lastRow := t.buffer.LineCount() - 1
	if lastRow < 0 {
		lastRow = 0
	}
	lastLine := t.buffer.Line(lastRow)
	lastCol := len([]rune(lastLine))

	copy := t.copy()
	copy.cursor = NewCursor(lastRow, lastCol)
	return copy
}

// Public Methods for Services (used by domain services)

// GetBuffer returns buffer (for services).
func (t *TextArea) GetBuffer() *Buffer {
	return t.buffer
}

// GetCursor returns cursor (for services).
func (t *TextArea) GetCursor() *Cursor {
	return t.cursor
}

// GetKillRing returns kill ring (for services).
func (t *TextArea) GetKillRing() *KillRing {
	return t.killRing
}

// WithKillRing sets kill ring (returns new instance).
func (t *TextArea) WithKillRing(killRing *KillRing) *TextArea {
	return t.withKillRing(killRing)
}

// Internal Methods (package-private)

// withCursor returns new TextArea with updated cursor.
func (t *TextArea) withCursor(cursor *Cursor) *TextArea {
	copy := t.copy()
	copy.cursor = cursor
	copy.ensureCursorVisible()
	return copy
}

// withBuffer returns new TextArea with updated buffer.
func (t *TextArea) withBuffer(buffer *Buffer) *TextArea {
	copy := t.copy()
	copy.buffer = buffer
	return copy
}

// withKillRing returns new TextArea with updated kill ring.
func (t *TextArea) withKillRing(killRing *KillRing) *TextArea {
	copy := t.copy()
	copy.killRing = killRing
	return copy
}

// ensureCursorVisible adjusts scroll to keep cursor visible.
func (t *TextArea) ensureCursorVisible() {
	row := t.cursor.Row()

	// Vertical scrolling
	if row < t.scrollRow {
		t.scrollRow = row
	}
	if row >= t.scrollRow+t.height {
		t.scrollRow = row - t.height + 1
	}

	// Horizontal scrolling (if no wrap)
	if !t.wrap {
		col := t.cursor.Col()
		if col < t.scrollCol {
			t.scrollCol = col
		}
		if col >= t.scrollCol+t.width {
			t.scrollCol = col - t.width + 1
		}
	}
}

// Private helper: deep copy all fields.
func (t *TextArea) copy() *TextArea {
	return &TextArea{
		buffer:          t.buffer.Copy(),
		cursor:          t.cursor.Copy(),
		selection:       t.selection.Copy(), // nil-safe
		killRing:        t.killRing.Copy(),
		width:           t.width,
		height:          t.height,
		scrollRow:       t.scrollRow,
		scrollCol:       t.scrollCol,
		maxLines:        t.maxLines,
		maxChars:        t.maxChars,
		placeholder:     t.placeholder,
		wrap:            t.wrap,
		readOnly:        t.readOnly,
		showLineNumbers: t.showLineNumbers,
		lineNumberWidth: t.lineNumberWidth,
		showCursor:      t.showCursor,
	}
}
