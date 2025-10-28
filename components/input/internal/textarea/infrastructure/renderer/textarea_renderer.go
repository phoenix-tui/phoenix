// Package renderer provides rendering infrastructure for textarea.
package renderer

import (
	"fmt"
	"strings"

	"github.com/phoenix-tui/phoenix/components/input/internal/textarea/domain/model"
	"github.com/phoenix-tui/phoenix/style"
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

// renderPlaceholder renders placeholder text with gray styling.
func (r *TextAreaRenderer) renderPlaceholder(ta *model.TextArea) string {
	// Apply gray foreground color (ANSI color 240 is a nice gray)
	grayStyle := style.New().Foreground(style.Color256(240))
	return style.Render(grayStyle, ta.Placeholder())
}

// renderContent renders actual content.
func (r *TextAreaRenderer) renderContent(ta *model.TextArea) string {
	var b strings.Builder

	visibleLines := ta.VisibleLines()
	cursorRow, cursorCol := ta.CursorPosition()

	for i, line := range visibleLines {
		// Calculate actual row number (accounting for scroll offset).
		actualRow := i + ta.ScrollRow()

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

// renderLineWithCursor renders a line with cursor visible using reverse video.
func (r *TextAreaRenderer) renderLineWithCursor(line string, col int) string {
	runes := []rune(line)

	if col >= len(runes) {
		// Cursor at end of line - use reverse video space for better visibility.
		return line + "\x1b[7m \x1b[27m" // Reverse video space
	}

	// Cursor in middle of line - apply reverse video to the character under cursor.
	before := string(runes[:col])
	cursorChar := string(runes[col])
	after := string(runes[col+1:])

	// Apply reverse video to cursor character.
	return before + "\x1b[7m" + cursorChar + "\x1b[27m" + after
}
