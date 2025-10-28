// Package input provides keyboard input reading and parsing.
package input

import (
	"bufio"
	"io"

	"github.com/phoenix-tui/phoenix/tea/internal/domain/model"
	"github.com/phoenix-tui/phoenix/tea/internal/infrastructure/ansi"
)

// Reader reads input from stdin and parses it into messages.
type Reader struct {
	reader *bufio.Reader
	parser *ansi.Parser
}

// NewReader creates a new input reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		reader: bufio.NewReader(r),
		parser: ansi.NewParser(),
	}
}

// Read reads input and returns a message.
// Blocks until input is available.
//
// Returns:
// - KeyMsg if keyboard input.
// - error if read fails.
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
