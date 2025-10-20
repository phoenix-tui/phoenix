// Package model provides rich domain models for textarea.
package model

import (
	"strings"

	"github.com/phoenix-tui/phoenix/components/input/textarea/domain/value"
)

// Buffer manages the text content as array of lines.
// This is a rich domain model that encapsulates behavior.
// All operations return new instances (immutable).
type Buffer struct {
	lines []string
}

// NewBuffer creates empty buffer with one empty line.
func NewBuffer() *Buffer {
	return &Buffer{
		lines: []string{""},
	}
}

// NewBufferFromString creates buffer from string (splits on newlines).
func NewBufferFromString(text string) *Buffer {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}
	return &Buffer{lines: lines}
}

// Lines returns all lines (defensive copy to prevent mutation).
func (b *Buffer) Lines() []string {
	result := make([]string, len(b.lines))
	copy(result, b.lines)
	return result
}

// Line returns specific line (bounds-checked).
func (b *Buffer) Line(row int) string {
	if row < 0 || row >= len(b.lines) {
		return ""
	}
	return b.lines[row]
}

// LineCount returns number of lines.
func (b *Buffer) LineCount() int {
	return len(b.lines)
}

// String returns all text as single string (lines joined with \n).
func (b *Buffer) String() string {
	return strings.Join(b.lines, "\n")
}

// IsEmpty returns true if buffer has no content.
func (b *Buffer) IsEmpty() bool {
	return len(b.lines) == 1 && b.lines[0] == ""
}

// InsertChar inserts character at position (returns new buffer).
func (b *Buffer) InsertChar(row, col int, ch rune) *Buffer {
	updated := b.Copy()

	if row >= len(updated.lines) {
		return updated
	}

	line := []rune(updated.lines[row])

	// Clamp col to valid range.
	if col > len(line) {
		col = len(line)
	}

	// Insert at position.
	newLine := make([]rune, len(line)+1)
	copied := 0
	//nolint:ineffassign // copied variable tracks offset for clarity
	copied += copy2(newLine[copied:], line[:col])
	newLine[col] = ch
	copy2(newLine[col+1:], line[col:])

	updated.lines[row] = string(newLine)
	return updated
}

// DeleteChar deletes character at position (returns new buffer).
func (b *Buffer) DeleteChar(row, col int) *Buffer {
	updated := b.Copy()

	if row >= len(updated.lines) {
		return updated
	}

	line := []rune(updated.lines[row])

	if col >= len(line) {
		return updated
	}

	newLine := make([]rune, len(line)-1)
	copied := 0
	//nolint:ineffassign // copied variable tracks offset for clarity
	copied += copy2(newLine[copied:], line[:col])
	copy2(newLine[col:], line[col+1:])

	updated.lines[row] = string(newLine)
	return updated
}

// InsertNewline splits line at position (returns new buffer).
func (b *Buffer) InsertNewline(row, col int) *Buffer {
	updated := b.Copy()

	if row >= len(updated.lines) {
		return updated
	}

	line := []rune(updated.lines[row])

	// Clamp col to valid range.
	if col > len(line) {
		col = len(line)
	}

	// Split line at cursor.
	before := string(line[:col])
	after := string(line[col:])

	// Replace current line with "before".
	updated.lines[row] = before

	// Insert "after" as new line.
	newLines := make([]string, len(updated.lines)+1)
	copied := 0
	//nolint:ineffassign // copied variable tracks offset for clarity
	copied += copy2(newLines[copied:], updated.lines[:row+1])
	newLines[row+1] = after
	copy2(newLines[row+2:], updated.lines[row+1:])

	updated.lines = newLines
	return updated
}

// DeleteLine removes entire line (returns new buffer and deleted text).
func (b *Buffer) DeleteLine(row int) (*Buffer, string) {
	updated := b.Copy()

	if row >= len(updated.lines) {
		return updated, ""
	}

	deletedLine := updated.lines[row]

	// Keep at least one empty line.
	if len(updated.lines) == 1 {
		updated.lines[0] = ""
		return updated, deletedLine
	}

	newLines := make([]string, len(updated.lines)-1)
	copied := 0
	//nolint:ineffassign // copied variable tracks offset for clarity
	copied += copy2(newLines[copied:], updated.lines[:row])
	copy2(newLines[row:], updated.lines[row+1:])

	updated.lines = newLines
	return updated, deletedLine
}

// DeleteToLineEnd deletes from position to end of line (returns new buffer and deleted text).
func (b *Buffer) DeleteToLineEnd(row, col int) (*Buffer, string) {
	updated := b.Copy()

	if row >= len(updated.lines) {
		return updated, ""
	}

	line := []rune(updated.lines[row])

	if col >= len(line) {
		return updated, ""
	}

	deleted := string(line[col:])
	updated.lines[row] = string(line[:col])

	return updated, deleted
}

// SetLine replaces entire line (returns new buffer).
func (b *Buffer) SetLine(row int, text string) *Buffer {
	updated := b.Copy()

	if row >= len(updated.lines) {
		return updated
	}

	updated.lines[row] = text
	return updated
}

// JoinWithNextLine joins current line with next line (returns new buffer).
func (b *Buffer) JoinWithNextLine(row int) *Buffer {
	updated := b.Copy()

	if row >= len(updated.lines)-1 {
		return updated
	}

	// Join lines.
	//nolint:gocritic // assignOp less readable for string concatenation in domain logic
	updated.lines[row] = updated.lines[row] + updated.lines[row+1]

	// Remove next line.
	newLines := make([]string, len(updated.lines)-1)
	copied := 0
	//nolint:ineffassign // copied variable tracks offset for clarity
	copied += copy2(newLines[copied:], updated.lines[:row+1])
	copy2(newLines[row+1:], updated.lines[row+2:])

	updated.lines = newLines
	return updated
}

// TextInRange returns text in range.
func (b *Buffer) TextInRange(r value.Range) string {
	startRow, startCol := r.StartRowCol()
	endRow, endCol := r.EndRowCol()

	if startRow == endRow {
		// Single line selection.
		line := []rune(b.Line(startRow))
		if startCol >= len(line) || endCol > len(line) {
			return ""
		}
		return string(line[startCol:endCol])
	}

	// Multi-line selection.
	var result strings.Builder

	// First line (from startCol to end)
	firstLine := []rune(b.Line(startRow))
	if startCol < len(firstLine) {
		result.WriteString(string(firstLine[startCol:]))
	}
	result.WriteRune('\n')

	// Middle lines (full lines)
	for row := startRow + 1; row < endRow; row++ {
		result.WriteString(b.Line(row))
		result.WriteRune('\n')
	}

	// Last line (from start to endCol)
	lastLine := []rune(b.Line(endRow))
	if endCol <= len(lastLine) {
		result.WriteString(string(lastLine[:endCol]))
	}

	return result.String()
}

// Copy returns deep copy of buffer.
func (b *Buffer) Copy() *Buffer {
	linesCopy := make([]string, len(b.lines))
	copy(linesCopy, b.lines)
	return &Buffer{lines: linesCopy}
}

// copy2 is a helper that copies a slice and returns the number of elements copied.
func copy2[T any](dst, src []T) int {
	return copy(dst, src)
}
