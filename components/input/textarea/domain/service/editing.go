// Package service provides domain services for textarea.
package service

import "github.com/phoenix-tui/phoenix/components/input/textarea/domain/model"

// EditingService handles text editing operations.
// This is a domain service that operates on the TextArea aggregate.
type EditingService struct{}

// NewEditingService creates editing service.
func NewEditingService() *EditingService {
	return &EditingService{}
}

// InsertChar inserts character at cursor position.
func (s *EditingService) InsertChar(ta *model.TextArea, ch rune) *model.TextArea {
	if ta.IsReadOnly() {
		return ta // Read-only
	}

	row, col := ta.CursorPosition()
	buffer := ta.GetBuffer()

	// Insert character in buffer.
	newBuffer := buffer.InsertChar(row, col, ch)

	// Move cursor right.
	newCursor := model.NewCursor(row, col+1)

	// Return updated TextArea.
	return ta.WithBuffer(newBuffer).WithCursor(newCursor)
}

// DeleteCharBackward deletes character before cursor (Backspace).
// Checks movement validator BEFORE deleting to prevent deletion if cursor can't move back.
func (s *EditingService) DeleteCharBackward(ta *model.TextArea) *model.TextArea {
	if ta.IsReadOnly() {
		return ta
	}

	row, col := ta.CursorPosition()

	if col == 0 && row == 0 {
		return ta // Nothing to delete
	}

	// Calculate where cursor would move after deletion.
	var targetRow, targetCol int
	if col > 0 {
		// Would move left on same line.
		targetRow, targetCol = row, col-1
	} else {
		// Would move to end of previous line (col == 0, row > 0)
		targetRow = row - 1
		targetCol = len([]rune(ta.Lines()[row-1]))
	}

	// Check validator BEFORE deleting.
	//nolint:nestif // domain validation logic requires nested conditions
	if validator := ta.GetMovementValidator(); validator != nil {
		from := model.NewCursorPos(row, col)
		to := model.NewCursorPos(targetRow, targetCol)
		if !validator(from, to) {
			// Movement blocked - can't delete because cursor can't move back.
			if handler := ta.GetBoundaryHitHandler(); handler != nil {
				handler(to, "backspace blocked by validator")
			}
			return ta
		}
	}

	// Validator allows movement - proceed with deletion.
	buffer := ta.GetBuffer()

	if col > 0 {
		// Delete character on current line.
		newBuffer := buffer.DeleteChar(row, col-1)
		newCursor := model.NewCursor(row, col-1)
		return ta.WithBuffer(newBuffer).WithCursor(newCursor)
	}

	// Join with previous line (col == 0, row > 0)
	prevLineLen := len([]rune(ta.Lines()[row-1]))

	// Get content of current line.
	currentLineContent := ta.CurrentLine()

	// Set previous line to previous + current.
	newBuffer := buffer.SetLine(row-1, buffer.Line(row-1)+currentLineContent)

	// Delete current line.
	newBuffer, _ = newBuffer.DeleteLine(row)

	// Move cursor to end of previous line.
	newCursor := model.NewCursor(row-1, prevLineLen)

	return ta.WithBuffer(newBuffer).WithCursor(newCursor)
}

// DeleteCharForward deletes character at cursor (Delete key).
func (s *EditingService) DeleteCharForward(ta *model.TextArea) *model.TextArea {
	if ta.IsReadOnly() {
		return ta
	}

	row, col := ta.CursorPosition()
	currentLine := ta.CurrentLine()
	buffer := ta.GetBuffer()

	if col >= len([]rune(currentLine)) {
		// At end of line, join with next line.
		if row < ta.LineCount()-1 {
			nextLineContent := ta.Lines()[row+1]
			newBuffer := buffer.SetLine(row, currentLine+nextLineContent)
			newBuffer, _ = newBuffer.DeleteLine(row + 1)
			// Preserve cursor position when joining lines.
			return ta.WithBuffer(newBuffer).WithCursor(model.NewCursor(row, col))
		}
		return ta
	}

	// Delete character at cursor (cursor stays in same position)
	newBuffer := buffer.DeleteChar(row, col)
	return ta.WithBuffer(newBuffer).WithCursor(model.NewCursor(row, col))
}

// InsertNewline inserts newline at cursor (Enter key).
func (s *EditingService) InsertNewline(ta *model.TextArea) *model.TextArea {
	if ta.IsReadOnly() {
		return ta
	}

	// Check max lines (if limit is set)
	maxLines := ta.MaxLines()
	if maxLines > 0 && ta.LineCount() >= maxLines {
		// Max lines reached.
		return ta
	}

	row, col := ta.CursorPosition()
	buffer := ta.GetBuffer()

	// Split current line at cursor.
	newBuffer := buffer.InsertNewline(row, col)

	// Move cursor to start of new line.
	newCursor := model.NewCursor(row+1, 0)

	return ta.WithBuffer(newBuffer).WithCursor(newCursor)
}

// KillLine deletes from cursor to end of line (Ctrl+K).
func (s *EditingService) KillLine(ta *model.TextArea) *model.TextArea {
	if ta.IsReadOnly() {
		return ta
	}

	row, col := ta.CursorPosition()
	currentLine := ta.CurrentLine()
	runes := []rune(currentLine)
	buffer := ta.GetBuffer()

	if col >= len(runes) {
		// At end of line, kill newline (join with next line)
		if row < ta.LineCount()-1 {
			// Kill the newline.
			killed := "\n"
			killRing := ta.GetKillRing().Kill(killed)

			// Join with next line.
			nextLineContent := ta.Lines()[row+1]
			newBuffer := buffer.SetLine(row, currentLine+nextLineContent)
			newBuffer, _ = newBuffer.DeleteLine(row + 1)

			return ta.WithBuffer(newBuffer).WithKillRing(killRing)
		}
		return ta
	}

	// Kill from cursor to end of line.
	killed := string(runes[col:])
	killRing := ta.GetKillRing().Kill(killed)

	// Delete from cursor to end.
	newBuffer, _ := buffer.DeleteToLineEnd(row, col)

	return ta.WithBuffer(newBuffer).WithKillRing(killRing)
}

// KillWord deletes word after cursor (Alt+D).
func (s *EditingService) KillWord(ta *model.TextArea) *model.TextArea {
	if ta.IsReadOnly() {
		return ta
	}

	row, col := ta.CursorPosition()
	line := []rune(ta.CurrentLine())
	buffer := ta.GetBuffer()

	// Find next word boundary.
	startCol := col
	for col < len(line) && !isWordBoundary(line[col]) {
		col++
	}

	if startCol == col {
		return ta // No word to kill
	}

	// Kill from startCol to col.
	killed := string(line[startCol:col])
	killRing := ta.GetKillRing().Kill(killed)

	// Delete from buffer.
	newBuffer := buffer
	for i := startCol; i < col; i++ {
		newBuffer = newBuffer.DeleteChar(row, startCol)
	}

	return ta.WithBuffer(newBuffer).WithKillRing(killRing)
}

// KillWordBackward deletes word before cursor (Ctrl+W, Alt+Backspace).
func (s *EditingService) KillWordBackward(ta *model.TextArea) *model.TextArea {
	if ta.IsReadOnly() {
		return ta
	}

	row, col := ta.CursorPosition()
	line := []rune(ta.CurrentLine())
	buffer := ta.GetBuffer()

	if col == 0 {
		return ta // At start of line, nothing to kill backward
	}

	// Find start of word backward (similar to BackwardWord navigation).
	endCol := col
	newCol := col - 1

	// Skip whitespace backward.
	for newCol > 0 && isWordBoundary(line[newCol]) {
		newCol--
	}

	// Skip current word backward.
	for newCol > 0 && !isWordBoundary(line[newCol-1]) {
		newCol--
	}

	// Note: If newCol is 0 and line[0] is not a word boundary, we delete from start.
	// No special handling needed as the slice operation handles this correctly.

	if newCol == endCol {
		return ta // No word to kill
	}

	// Kill from newCol to endCol.
	killed := string(line[newCol:endCol])
	killRing := ta.GetKillRing().Kill(killed)

	// Delete from buffer (delete backward from cursor).
	newBuffer := buffer
	for i := newCol; i < endCol; i++ {
		newBuffer = newBuffer.DeleteChar(row, newCol)
	}

	// Move cursor to where deletion started.
	newCursor := model.NewCursor(row, newCol)

	return ta.WithBuffer(newBuffer).WithCursor(newCursor).WithKillRing(killRing)
}

// Yank inserts text from kill ring (Ctrl+Y).
func (s *EditingService) Yank(ta *model.TextArea) *model.TextArea {
	if ta.IsReadOnly() {
		return ta
	}

	killRing := ta.GetKillRing()
	text := killRing.Yank()

	if text == "" {
		return ta
	}

	// Insert text at cursor.
	row, col := ta.CursorPosition()
	buffer := ta.GetBuffer()

	// Insert each character.
	for _, ch := range text {
		if ch == '\n' {
			// Insert newline.
			buffer = buffer.InsertNewline(row, col)
			row++
			col = 0
		} else {
			// Insert character.
			buffer = buffer.InsertChar(row, col, ch)
			col++
		}
	}

	newCursor := model.NewCursor(row, col)
	return ta.WithBuffer(buffer).WithCursor(newCursor)
}
