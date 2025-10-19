// Package service provides domain services for textarea.
package service

import "github.com/phoenix-tui/phoenix/components/input/textarea/domain/model"

// NavigationService handles cursor movement logic.
// This is a domain service that operates on the TextArea aggregate.
type NavigationService struct{}

// NewNavigationService creates navigation service.
func NewNavigationService() *NavigationService {
	return &NavigationService{}
}

// MoveLeft moves cursor left by one character.
func (s *NavigationService) MoveLeft(ta *model.TextArea) *model.TextArea {
	row, col := ta.CursorPosition()

	if col > 0 {
		// Move left within line
		return ta.WithCursor(model.NewCursor(row, col-1))
	}

	if row > 0 {
		// Move to end of previous line
		prevLineLen := len([]rune(ta.Lines()[row-1]))
		return ta.WithCursor(model.NewCursor(row-1, prevLineLen))
	}

	// Already at start
	return ta
}

// MoveRight moves cursor right by one character.
func (s *NavigationService) MoveRight(ta *model.TextArea) *model.TextArea {
	row, col := ta.CursorPosition()
	currentLine := ta.CurrentLine()
	currentLineLen := len([]rune(currentLine))

	if col < currentLineLen {
		// Move right within line
		return ta.WithCursor(model.NewCursor(row, col+1))
	}

	if row < ta.LineCount()-1 {
		// Move to start of next line
		return ta.WithCursor(model.NewCursor(row+1, 0))
	}

	// Already at end
	return ta
}

// MoveUp moves cursor up one line.
func (s *NavigationService) MoveUp(ta *model.TextArea) *model.TextArea {
	row, col := ta.CursorPosition()

	if row == 0 {
		return ta
	}

	// Try to maintain column position
	newRow := row - 1
	newLineLen := len([]rune(ta.Lines()[newRow]))

	newCol := col
	if newCol > newLineLen {
		newCol = newLineLen
	}

	return ta.WithCursor(model.NewCursor(newRow, newCol))
}

// MoveDown moves cursor down one line.
func (s *NavigationService) MoveDown(ta *model.TextArea) *model.TextArea {
	row, col := ta.CursorPosition()

	if row >= ta.LineCount()-1 {
		return ta
	}

	// Try to maintain column position
	newRow := row + 1
	newLineLen := len([]rune(ta.Lines()[newRow]))

	newCol := col
	if newCol > newLineLen {
		newCol = newLineLen
	}

	return ta.WithCursor(model.NewCursor(newRow, newCol))
}

// MoveToLineStart moves cursor to start of current line (Ctrl+A / Home).
func (s *NavigationService) MoveToLineStart(ta *model.TextArea) *model.TextArea {
	row, _ := ta.CursorPosition()
	return ta.WithCursor(model.NewCursor(row, 0))
}

// MoveToLineEnd moves cursor to end of current line (Ctrl+E / End).
func (s *NavigationService) MoveToLineEnd(ta *model.TextArea) *model.TextArea {
	row, _ := ta.CursorPosition()
	lineLen := len([]rune(ta.CurrentLine()))
	return ta.WithCursor(model.NewCursor(row, lineLen))
}

// MoveToBufferStart moves cursor to start of buffer (Alt+<).
func (s *NavigationService) MoveToBufferStart(ta *model.TextArea) *model.TextArea {
	return ta.WithCursor(model.NewCursor(0, 0))
}

// MoveToBufferEnd moves cursor to end of buffer (Alt+>).
func (s *NavigationService) MoveToBufferEnd(ta *model.TextArea) *model.TextArea {
	lastRow := ta.LineCount() - 1
	lastLineLen := len([]rune(ta.Lines()[lastRow]))
	return ta.WithCursor(model.NewCursor(lastRow, lastLineLen))
}

// ForwardWord moves cursor forward by one word (Alt+F).
func (s *NavigationService) ForwardWord(ta *model.TextArea) *model.TextArea {
	row, col := ta.CursorPosition()
	line := []rune(ta.CurrentLine())

	// Skip current word
	for col < len(line) && !isWordBoundary(line[col]) {
		col++
	}

	// Skip whitespace
	for col < len(line) && isWordBoundary(line[col]) {
		col++
	}

	return ta.WithCursor(model.NewCursor(row, col))
}

// BackwardWord moves cursor backward by one word (Alt+B).
func (s *NavigationService) BackwardWord(ta *model.TextArea) *model.TextArea {
	row, col := ta.CursorPosition()
	line := []rune(ta.CurrentLine())

	if col == 0 {
		return ta
	}

	col-- // Move back one

	// Skip whitespace
	for col > 0 && isWordBoundary(line[col]) {
		col--
	}

	// Skip current word
	for col > 0 && !isWordBoundary(line[col-1]) {
		col--
	}

	return ta.WithCursor(model.NewCursor(row, col))
}

// isWordBoundary returns true if rune is word boundary (space, punctuation).
func isWordBoundary(r rune) bool {
	return r == ' ' || r == '\t' || r == '.' || r == ',' || r == ';' || r == ':' ||
		r == '!' || r == '?' || r == '(' || r == ')' || r == '[' || r == ']' ||
		r == '{' || r == '}' || r == '/' || r == '\\' || r == '|' || r == '-' ||
		r == '_' || r == '=' || r == '+' || r == '*' || r == '&' || r == '%' ||
		r == '$' || r == '#' || r == '@' || r == '~' || r == '`' || r == '\'' ||
		r == '"' || r == '<' || r == '>'
}
