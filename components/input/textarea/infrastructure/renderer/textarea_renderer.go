// Package renderer provides rendering infrastructure for textarea.
package renderer

import (
	"fmt"
	"strings"

	"github.com/phoenix-tui/phoenix/components/input/textarea/domain/model"
)

// TextAreaRenderer renders a textarea to a string.
// This is an infrastructure component that handles presentation logic.
type TextAreaRenderer struct{}

// NewTextAreaRenderer creates a new textarea renderer.
func NewTextAreaRenderer() *TextAreaRenderer {
	return &TextAreaRenderer{}
}

// Render renders the textarea to a string.
func (r *TextAreaRenderer) Render(ta *model.TextArea) string {
	if ta.IsEmpty() {
		if ta.Placeholder() != "" {
			return r.renderPlaceholder(ta)
		}
		// Empty with no placeholder - return empty string (not cursor)
		return ""
	}

	return r.renderContent(ta)
}

// renderPlaceholder renders placeholder text.
func (r *TextAreaRenderer) renderPlaceholder(ta *model.TextArea) string {
	// TODO: Style placeholder (gray text)
	return ta.Placeholder()
}

// renderContent renders actual content.
func (r *TextAreaRenderer) renderContent(ta *model.TextArea) string {
	var b strings.Builder

	visibleLines := ta.VisibleLines()
	cursorRow, cursorCol := ta.CursorPosition()

	for i, line := range visibleLines {
		actualRow := i // TODO: adjust for scroll offset

		// Render line number if enabled.
		if ta.ShowLineNumbers() {
			lineNum := fmt.Sprintf("%4d ", actualRow+1)
			b.WriteString(lineNum)
		}

		// Render line content.
		if actualRow == cursorRow && ta.ShowCursor() {
			// Render line with cursor (only if ShowCursor enabled)
			b.WriteString(r.renderLineWithCursor(line, cursorCol))
		} else {
			// Render line without cursor.
			b.WriteString(line)
		}

		// Add newline (except for last line)
		if i < len(visibleLines)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

// renderLineWithCursor renders a line with cursor visible.
func (r *TextAreaRenderer) renderLineWithCursor(line string, col int) string {
	runes := []rune(line)

	if col >= len(runes) {
		// Cursor at end of line.
		return line + "█" // Block cursor
	}

	// Cursor in middle of line.
	before := string(runes[:col])
	// TODO: Style cursor as reverse video.
	after := string(runes[col+1:])

	// For now, just show cursor as block.
	return before + "█" + after
}
