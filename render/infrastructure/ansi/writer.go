//nolint:nestif // Style application requires nested conditionals
package ansi

import (
	"bufio"
	"io"
	"sync"

	"github.com/phoenix-tui/phoenix/render/domain/model"
	"github.com/phoenix-tui/phoenix/render/domain/value"
)

// Writer writes ANSI sequences efficiently with buffering.
// Designed for high-performance rendering with minimal allocations.
type Writer struct {
	output       io.Writer
	buf          *bufio.Writer
	currentX     int
	currentY     int
	currentStyle value.Style
	mu           sync.Mutex
}

// NewWriter creates a new ANSI writer.
func NewWriter(output io.Writer) *Writer {
	return &Writer{
		output:       output,
		buf:          bufio.NewWriter(output),
		currentX:     0,
		currentY:     0,
		currentStyle: value.NewStyle(),
	}
}

// NewWriterWithBuffer creates a writer with specified buffer size.
func NewWriterWithBuffer(output io.Writer, bufferSize int) *Writer {
	return &Writer{
		output:       output,
		buf:          bufio.NewWriterSize(output, bufferSize),
		currentX:     0,
		currentY:     0,
		currentStyle: value.NewStyle(),
	}
}

// MoveCursor moves cursor to (x, y) position.
func (w *Writer) MoveCursor(x, y int) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if x == w.currentX && y == w.currentY {
		return nil // Already at position
	}

	_, err := w.buf.WriteString(MoveCursor(x, y))
	if err != nil {
		return err
	}

	w.currentX = x
	w.currentY = y
	return nil
}

// MoveCursorRelative moves cursor relative to current position.
func (w *Writer) MoveCursorRelative(dx, dy int) error {
	return w.MoveCursor(w.currentX+dx, w.currentY+dy)
}

// SetStyle sets the current text style.
func (w *Writer) SetStyle(style value.Style) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if style.Equals(w.currentStyle) {
		return nil // Style already set
	}

	// Write style sequence.
	ansi := style.ToANSI()
	if ansi != "" {
		_, err := w.buf.WriteString(ansi)
		if err != nil {
			return err
		}
	}

	w.currentStyle = style
	return nil
}

// WriteCell writes a single cell at current position.
func (w *Writer) WriteCell(cell model.Cell) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Set style if different.
	if !cell.Style().Equals(w.currentStyle) {
		ansi := cell.Style().ToANSI()
		if ansi != "" {
			if _, err := w.buf.WriteString(ansi); err != nil {
				return err
			}
		}
		w.currentStyle = cell.Style()
	}

	// Write character.
	if _, err := w.buf.WriteRune(cell.Char()); err != nil {
		return err
	}

	// Update position.
	w.currentX += cell.Width()
	return nil
}

// WriteCellAt writes a cell at specific position.
func (w *Writer) WriteCellAt(pos value.Position, cell model.Cell) error {
	if err := w.MoveCursor(pos.X(), pos.Y()); err != nil {
		return err
	}
	return w.WriteCell(cell)
}

// WriteString writes a string with style at current position.
func (w *Writer) WriteString(text string, style value.Style) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Set style.
	if !style.Equals(w.currentStyle) {
		ansi := style.ToANSI()
		if ansi != "" {
			if _, err := w.buf.WriteString(ansi); err != nil {
				return err
			}
		}
		w.currentStyle = style
	}

	// Write text.
	_, err := w.buf.WriteString(text)
	if err != nil {
		return err
	}

	// Update position (approximate - assumes single-width chars).
	w.currentX += len(text)
	return nil
}

// WriteStringAt writes a string at specific position.
func (w *Writer) WriteStringAt(pos value.Position, text string, style value.Style) error {
	if err := w.MoveCursor(pos.X(), pos.Y()); err != nil {
		return err
	}
	return w.WriteString(text, style)
}

// Clear clears the entire screen.
func (w *Writer) Clear() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.buf.WriteString(ClearScreen)
	return err
}

// ClearLine clears the current line.
func (w *Writer) ClearLine() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.buf.WriteString(ClearLine)
	return err
}

// Reset resets all text attributes.
func (w *Writer) Reset() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.buf.WriteString(Reset)
	if err != nil {
		return err
	}

	w.currentStyle = value.NewStyle()
	return nil
}

// HideCursor hides the cursor.
func (w *Writer) HideCursor() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.buf.WriteString(CursorHide)
	return err
}

// ShowCursor shows the cursor.
func (w *Writer) ShowCursor() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.buf.WriteString(CursorShow)
	return err
}

// WriteRaw writes raw bytes without buffering or tracking.
func (w *Writer) WriteRaw(data []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.buf.Write(data)
	return err
}

// WriteEscape writes an ANSI escape sequence.
func (w *Writer) WriteEscape(seq string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.buf.WriteString(seq)
	return err
}

// Flush flushes the buffer to output.
func (w *Writer) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.buf.Flush()
}

// BufferedSize returns the number of bytes currently buffered.
func (w *Writer) BufferedSize() int {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.buf.Buffered()
}

// CurrentPosition returns current cursor position.
func (w *Writer) CurrentPosition() (x, y int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.currentX, w.currentY
}

// CurrentStyle returns current text style.
func (w *Writer) CurrentStyle() value.Style {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.currentStyle
}

// SetPosition sets internal cursor position without moving (for sync).
func (w *Writer) SetPosition(x, y int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.currentX = x
	w.currentY = y
}

// Close flushes and closes the writer.
func (w *Writer) Close() error {
	return w.Flush()
}
