// Package input provides keyboard input reading and parsing.
package input

import (
	"bufio"
	"io"

	"github.com/phoenix-tui/phoenix/tea/internal/domain/model"
	"github.com/phoenix-tui/phoenix/tea/internal/infrastructure/ansi"
)

// Reader reads input from stdin and parses it into messages.
// Supports cancellation for ExecProcess stdin release.
//
// The zero value is not usable; use NewReader.
type Reader struct {
	reader           *bufio.Reader
	parser           *ansi.Parser
	cancelableReader *CancelableReader // For cancellation support
}

// NewReader creates a new input reader with cancellation support.
//
// The reader wraps the provided io.Reader with CancelableReader,
// allowing Cancel() to immediately unblock any pending Read() calls.
// This is essential for ExecProcess to cleanly release stdin.
func NewReader(r io.Reader) *Reader {
	// Wrap with CancelableReader for cancellation support
	cancelableReader := NewCancelableReader(r)

	return &Reader{
		reader:           bufio.NewReader(cancelableReader),
		parser:           ansi.NewParser(),
		cancelableReader: cancelableReader,
	}
}

// Cancel cancels any pending Read operations.
// After Cancel(), Read() will return io.EOF.
//
// IMPORTANT: Must be called before ExecProcess to release stdin.
// This ensures the inputReader goroutine fully stops before
// the child process attempts to read from stdin.
//
// Safe to call multiple times.
func (ir *Reader) Cancel() {
	if ir.cancelableReader != nil {
		ir.cancelableReader.Cancel()
	}
}

// WaitForShutdown waits for background goroutines to exit after Cancel().
// Call this after Cancel() to ensure the pipe relay and read loop are fully stopped.
func (ir *Reader) WaitForShutdown() {
	if ir.cancelableReader != nil {
		ir.cancelableReader.WaitForShutdown()
	}
}

// IsCanceled returns true if the reader has been canceled.
func (ir *Reader) IsCanceled() bool {
	if ir.cancelableReader != nil {
		return ir.cancelableReader.IsCanceled()
	}
	return false
}

// Read reads input and returns a message.
// Blocks until input is available or Cancel() is called.
//
// Returns:
//   - KeyMsg if keyboard input
//   - nil, io.EOF if canceled or stream ended
//   - nil, error if read fails
//
// Properly handles UTF-8 multi-byte sequences (Russian, Chinese, emoji, etc.).
//
//nolint:gocognit,gocyclo,cyclop,nestif // Input parsing requires sequential byte analysis
func (ir *Reader) Read() (model.Msg, error) {
	// Read one byte to detect input type
	b, err := ir.reader.ReadByte()
	if err != nil {
		return nil, err
	}

	// Build sequence
	seq := []byte{b}

	// If ESC, try to read more bytes for ANSI sequence
	if b == 0x1B {
		// Peek ahead to see if more bytes available
		// (with simple buffered check to avoid blocking forever)

		// Try to read up to 10 more bytes (enough for most sequences)
		for i := 0; i < 10; i++ {
			// Check if byte available (non-blocking peek)
			if ir.reader.Buffered() > 0 {
				nextByte, err := ir.reader.ReadByte()
				if err != nil {
					break
				}
				seq = append(seq, nextByte)

				// Stop if we hit a letter or tilde (end of most sequences)
				if (nextByte >= 'A' && nextByte <= 'Z') ||
					(nextByte >= 'a' && nextByte <= 'z') ||
					nextByte == '~' {
					break
				}
			} else {
				break
			}
		}
	}

	// Parse sequence (handles special keys, ANSI sequences, ASCII)
	keyMsg, ok := ir.parser.ParseKey(seq)
	if ok {
		return keyMsg, nil
	}

	// If not recognized and starts with multi-byte UTF-8 (0x80+),
	// unread the byte and read as rune for proper UTF-8 handling
	if b >= 0x80 {
		// Unread the byte
		if err := ir.reader.UnreadByte(); err != nil {
			return nil, err
		}

		// Read as rune (handles multi-byte UTF-8 automatically)
		r, _, err := ir.reader.ReadRune()
		if err != nil {
			return nil, err
		}

		// Return as KeyRune
		return model.KeyMsg{Type: model.KeyRune, Rune: r}, nil
	}

	// ASCII printable but not recognized (shouldn't happen, but fallback)
	if b >= 32 && b <= 126 {
		return model.KeyMsg{Type: model.KeyRune, Rune: rune(b)}, nil
	}

	// Non-printable unknown byte - ignore and try to read next
	// Return nil to indicate no message (sentinel pattern for skipped input)
	//nolint:nilnil // Intentional: nil msg + nil error = skip this byte, continue reading
	return nil, nil
}
