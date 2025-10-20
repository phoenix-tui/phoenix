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
	from := model.NewCursorPos(row, col)

	var newRow, newCol int
	//nolint:gocritic // if-else is more readable than switch for position logic
	if col > 0 {
		// Move left within line.
		newRow, newCol = row, col-1
	} else if row > 0 {
		// Move to end of previous line.
		newRow = row - 1
		newCol = len([]rune(ta.Lines()[newRow]))
	} else {
		// Already at start - no movement.
		return ta
	}

	to := model.NewCursorPos(newRow, newCol)

	// Check validator (if set)
	//nolint:nestif // domain validation logic requires nested conditions
	if validator := ta.GetMovementValidator(); validator != nil {
		if !validator(from, to) {
			// Movement blocked - fire boundary hit.
			if handler := ta.GetBoundaryHitHandler(); handler != nil {
				handler(to, "movement blocked by validator")
			}
			return ta // Return unchanged
		}
	}

	// Apply movement.
	result := ta.WithCursor(model.NewCursor(newRow, newCol))

	// Fire cursor moved event.
	if handler := ta.GetCursorMovedHandler(); handler != nil {
		handler(from, to)
	}

	return result
}

// MoveRight moves cursor right by one character.
func (s *NavigationService) MoveRight(ta *model.TextArea) *model.TextArea {
	row, col := ta.CursorPosition()
	from := model.NewCursorPos(row, col)
	currentLine := ta.CurrentLine()
	currentLineLen := len([]rune(currentLine))

	var newRow, newCol int
	//nolint:gocritic // if-else is more readable than switch for position logic
	if col < currentLineLen {
		// Move right within line.
		newRow, newCol = row, col+1
	} else if row < ta.LineCount()-1 {
		// Move to start of next line.
		newRow, newCol = row+1, 0
	} else {
		// Already at end - no movement.
		return ta
	}

	to := model.NewCursorPos(newRow, newCol)

	// Check validator (if set)
	//nolint:nestif // domain validation logic requires nested conditions
	if validator := ta.GetMovementValidator(); validator != nil {
		if !validator(from, to) {
			// Movement blocked - fire boundary hit.
			if handler := ta.GetBoundaryHitHandler(); handler != nil {
				handler(to, "movement blocked by validator")
			}
			return ta // Return unchanged
		}
	}

	// Apply movement.
	result := ta.WithCursor(model.NewCursor(newRow, newCol))

	// Fire cursor moved event.
	if handler := ta.GetCursorMovedHandler(); handler != nil {
		handler(from, to)
	}

	return result
}

// MoveUp moves cursor up one line.
func (s *NavigationService) MoveUp(ta *model.TextArea) *model.TextArea {
	row, col := ta.CursorPosition()
	from := model.NewCursorPos(row, col)

	//nolint:nestif // domain boundary validation requires nested conditions
	if row == 0 {
		// Already at top - check if validator wants to handle this.
		if validator := ta.GetMovementValidator(); validator != nil {
			to := model.NewCursorPos(row-1, col) // Attempted position
			if !validator(from, to) {
				if handler := ta.GetBoundaryHitHandler(); handler != nil {
					handler(to, "already at top")
				}
			}
		}
		return ta
	}

	// Try to maintain column position.
	newRow := row - 1
	newLineLen := len([]rune(ta.Lines()[newRow]))
	newCol := col
	if newCol > newLineLen {
		newCol = newLineLen
	}

	to := model.NewCursorPos(newRow, newCol)

	// Check validator (if set)
	//nolint:nestif // domain validation logic requires nested conditions
	if validator := ta.GetMovementValidator(); validator != nil {
		if !validator(from, to) {
			// Movement blocked - fire boundary hit.
			if handler := ta.GetBoundaryHitHandler(); handler != nil {
				handler(to, "movement blocked by validator")
			}
			return ta // Return unchanged
		}
	}

	// Apply movement.
	result := ta.WithCursor(model.NewCursor(newRow, newCol))

	// Fire cursor moved event.
	if handler := ta.GetCursorMovedHandler(); handler != nil {
		handler(from, to)
	}

	return result
}

// MoveDown moves cursor down one line.
func (s *NavigationService) MoveDown(ta *model.TextArea) *model.TextArea {
	row, col := ta.CursorPosition()
	from := model.NewCursorPos(row, col)

	//nolint:nestif // domain boundary validation requires nested conditions
	if row >= ta.LineCount()-1 {
		// Already at bottom - check if validator wants to handle this.
		if validator := ta.GetMovementValidator(); validator != nil {
			to := model.NewCursorPos(row+1, col) // Attempted position
			if !validator(from, to) {
				if handler := ta.GetBoundaryHitHandler(); handler != nil {
					handler(to, "already at bottom")
				}
			}
		}
		return ta
	}

	// Try to maintain column position.
	newRow := row + 1
	newLineLen := len([]rune(ta.Lines()[newRow]))
	newCol := col
	if newCol > newLineLen {
		newCol = newLineLen
	}

	to := model.NewCursorPos(newRow, newCol)

	// Check validator (if set)
	//nolint:nestif // domain validation logic requires nested conditions
	if validator := ta.GetMovementValidator(); validator != nil {
		if !validator(from, to) {
			// Movement blocked - fire boundary hit.
			if handler := ta.GetBoundaryHitHandler(); handler != nil {
				handler(to, "movement blocked by validator")
			}
			return ta // Return unchanged
		}
	}

	// Apply movement.
	result := ta.WithCursor(model.NewCursor(newRow, newCol))

	// Fire cursor moved event.
	if handler := ta.GetCursorMovedHandler(); handler != nil {
		handler(from, to)
	}

	return result
}

// MoveToLineStart moves cursor to start of current line (Ctrl+A / Home).
// If validator blocks column 0, tries to find first allowed position (col 1, 2, ...).
func (s *NavigationService) MoveToLineStart(ta *model.TextArea) *model.TextArea {
	row, col := ta.CursorPosition()
	from := model.NewCursorPos(row, col)

	validator := ta.GetMovementValidator()
	if validator == nil {
		// No validator - move to column 0.
		return s.moveCursor(ta, row, col, row, 0)
	}

	// Try column 0 first.
	to := model.NewCursorPos(row, 0)
	if validator(from, to) {
		return s.moveCursor(ta, row, col, row, 0)
	}

	// Column 0 blocked - find first allowed position.
	lineLen := len([]rune(ta.CurrentLine()))
	for tryCol := 1; tryCol <= lineLen; tryCol++ {
		to = model.NewCursorPos(row, tryCol)
		if validator(from, to) {
			return s.moveCursor(ta, row, col, row, tryCol)
		}
	}

	// All positions blocked - fire boundary hit and stay put.
	if handler := ta.GetBoundaryHitHandler(); handler != nil {
		handler(model.NewCursorPos(row, 0), "no valid position found in line")
	}
	return ta
}

// MoveToLineEnd moves cursor to end of current line (Ctrl+E / End).
func (s *NavigationService) MoveToLineEnd(ta *model.TextArea) *model.TextArea {
	row, col := ta.CursorPosition()
	lineLen := len([]rune(ta.CurrentLine()))
	return s.moveCursor(ta, row, col, row, lineLen)
}

// MoveToBufferStart moves cursor to start of buffer (Alt+<).
func (s *NavigationService) MoveToBufferStart(ta *model.TextArea) *model.TextArea {
	row, col := ta.CursorPosition()
	return s.moveCursor(ta, row, col, 0, 0)
}

// MoveToBufferEnd moves cursor to end of buffer (Alt+>).
func (s *NavigationService) MoveToBufferEnd(ta *model.TextArea) *model.TextArea {
	row, col := ta.CursorPosition()
	lastRow := ta.LineCount() - 1
	lastLineLen := len([]rune(ta.Lines()[lastRow]))
	return s.moveCursor(ta, row, col, lastRow, lastLineLen)
}

// ForwardWord moves cursor forward by one word (Alt+F).
func (s *NavigationService) ForwardWord(ta *model.TextArea) *model.TextArea {
	row, col := ta.CursorPosition()
	line := []rune(ta.CurrentLine())

	newCol := col
	// Skip current word.
	for newCol < len(line) && !isWordBoundary(line[newCol]) {
		newCol++
	}

	// Skip whitespace.
	for newCol < len(line) && isWordBoundary(line[newCol]) {
		newCol++
	}

	return s.moveCursor(ta, row, col, row, newCol)
}

// BackwardWord moves cursor backward by one word (Alt+B).
func (s *NavigationService) BackwardWord(ta *model.TextArea) *model.TextArea {
	row, col := ta.CursorPosition()
	line := []rune(ta.CurrentLine())

	if col == 0 {
		return ta
	}

	newCol := col - 1 // Move back one

	// Skip whitespace.
	for newCol > 0 && isWordBoundary(line[newCol]) {
		newCol--
	}

	// Skip current word.
	for newCol > 0 && !isWordBoundary(line[newCol-1]) {
		newCol--
	}

	return s.moveCursor(ta, row, col, row, newCol)
}

// moveCursor is a helper that handles validation and callbacks for cursor movement.
// This reduces code duplication across all movement methods.
func (s *NavigationService) moveCursor(ta *model.TextArea, fromRow, fromCol, toRow, toCol int) *model.TextArea {
	from := model.NewCursorPos(fromRow, fromCol)
	to := model.NewCursorPos(toRow, toCol)

	// No movement - return unchanged.
	if from.Equals(to) {
		return ta
	}

	// Check validator (if set)
	//nolint:nestif // domain validation logic requires nested conditions
	if validator := ta.GetMovementValidator(); validator != nil {
		if !validator(from, to) {
			// Movement blocked - fire boundary hit.
			if handler := ta.GetBoundaryHitHandler(); handler != nil {
				handler(to, "movement blocked by validator")
			}
			return ta // Return unchanged
		}
	}

	// Apply movement.
	result := ta.WithCursor(model.NewCursor(toRow, toCol))

	// Fire cursor moved event.
	if handler := ta.GetCursorMovedHandler(); handler != nil {
		handler(from, to)
	}

	return result
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
