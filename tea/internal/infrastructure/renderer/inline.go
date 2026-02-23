// Package renderer provides inline (non-alt-screen) rendering for the tea module.
// InlineRenderer tracks the number of lines rendered in the previous frame and
// uses ANSI cursor-up sequences to overwrite the previous frame on each update,
// preventing the stacked-output bug where each render appends below the last.
package renderer

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"unicode"
)

// ANSI escape sequences used for inline rendering.
const (
	// eraseLineRight erases from the cursor to the end of the current line.
	// Equivalent to \x1b[0K. Used after writing each line to clear leftover chars.
	eraseLineRight = "\x1b[K"

	// eraseScreenBelow erases from the cursor to the end of the screen.
	// Used when the new frame has fewer lines than the previous frame.
	eraseScreenBelow = "\x1b[J"

	// carriageReturn moves the cursor to column 0 of the current line.
	carriageReturn = "\r"
)

// cursorUp returns the ANSI sequence to move the cursor up n lines.
// Returns empty string if n <= 0.
func cursorUp(n int) string {
	if n <= 0 {
		return ""
	}
	return fmt.Sprintf("\x1b[%dA", n)
}

// InlineRenderer handles inline (non-alt-screen) rendering.
//
// It tracks the number of lines from the previous render frame and uses
// ANSI cursor-up sequences to return to the top of the rendered region
// before overwriting. Per-line diffing skips unchanged lines to minimize I/O.
//
// All public methods are safe for concurrent use.
type InlineRenderer struct {
	out           io.Writer
	width         int      // terminal width; 0 means no truncation
	height        int      // terminal height; 0 means no clipping
	lastView      string   // raw string of the previous render (for identical-frame check)
	lastLines     []string // previous render split by "\n" (for per-line diffing)
	linesRendered int      // number of lines written in the previous render
	mu            sync.Mutex
}

// NewInlineRenderer creates a new InlineRenderer writing to out.
//
// width and height are the current terminal dimensions. Pass 0 to disable
// truncation (width) and clipping (height). Call Resize when the terminal
// dimensions change.
func NewInlineRenderer(out io.Writer, width, height int) *InlineRenderer {
	return &InlineRenderer{
		out:    out,
		width:  width,
		height: height,
	}
}

// Render writes view to the output, overwriting the previous render frame.
//
// On the first call it writes the content with ANSI line-erase sequences.
// On subsequent calls it moves the cursor back to the top of the rendered
// region and rewrites only the lines that changed since the last call.
//
// Returns nil on success. Write errors are returned to the caller but do not
// corrupt internal state (a failed render can be retried).
func (r *InlineRenderer) Render(view string) error { //nolint:gocognit,gocyclo,cyclop // render loop with per-line diffing has inherent branching
	r.mu.Lock()
	defer r.mu.Unlock()

	// No-op if the view is identical to the previous render.
	if view == r.lastView && r.linesRendered > 0 {
		return nil
	}

	buf := &bytes.Buffer{}

	// Move cursor back to the top of the previously rendered region.
	// When linesRendered == 0 this is the first render; no cursor-up needed.
	if r.linesRendered > 1 {
		buf.WriteString(cursorUp(r.linesRendered - 1))
	}
	// Return to column 0 of the current line (handles both first render
	// and subsequent renders where cursor sits at end of last line).
	buf.WriteString(carriageReturn)

	newLines := strings.Split(view, "\n")

	// Clip to terminal height — cannot scroll into the scrollback buffer
	// while the inline render region is active.
	if r.height > 0 && len(newLines) > r.height {
		newLines = newLines[len(newLines)-r.height:]
	}

	for i, line := range newLines {
		// Truncate to terminal width to prevent line wrap, which would
		// inflate the physical line count and corrupt cursor positioning.
		if r.width > 0 {
			line = truncateLine(line, r.width)
		}

		// Per-line diff: skip lines that are identical to the previous render.
		// Compare using the original (pre-truncation) content because lastLines
		// also stores originals. If width changed, Resize() clears the cache.
		if i < len(r.lastLines) && r.lastLines[i] == newLines[i] {
			if i < len(newLines)-1 {
				buf.WriteString("\r\n")
			}
			continue
		}

		// Write the (possibly truncated) line content followed by erase-to-EOL.
		// The erase sequence clears any leftover characters from a longer previous line.
		buf.WriteString(line)
		buf.WriteString(eraseLineRight)

		if i < len(newLines)-1 {
			buf.WriteString("\r\n")
		}
	}

	// Erase any extra lines from the previous render that are no longer present.
	if r.linesRendered > len(newLines) {
		buf.WriteString(eraseScreenBelow)
	}

	// Leave cursor at column 0 of the last rendered line for consistent positioning.
	buf.WriteString(carriageReturn)

	_, err := r.out.Write(buf.Bytes())
	if err != nil {
		return err
	}

	// Update state only after a successful write.
	r.linesRendered = len(newLines)
	r.lastView = view
	r.lastLines = newLines

	return nil
}

// Repaint clears the per-line diff cache so the next Render call performs a
// full repaint of all lines. Call this after Resume() to restore the TUI view
// after an external process has written to the terminal.
//
// NOTE: linesRendered is intentionally NOT reset here. The cursor is still at
// the end of the last rendered region; the next Render will move it back up
// the correct number of lines.
func (r *InlineRenderer) Repaint() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastView = ""
	r.lastLines = nil
}

// Resize updates the terminal dimensions and forces a full repaint on the next
// Render call. Call this when a WindowSizeMsg is received.
func (r *InlineRenderer) Resize(width, height int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.width = width
	r.height = height
	r.lastView = ""
	r.lastLines = nil
}

// SetOutput changes the destination writer. Useful for redirecting output
// without recreating the renderer (e.g., after a terminal handoff).
func (r *InlineRenderer) SetOutput(w io.Writer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.out = w
	r.lastView = ""
	r.lastLines = nil
	r.linesRendered = 0
}

// ─── Internal helpers ────────────────────────────────────────────────────────

// stripANSI removes ANSI escape sequences from s and returns the plain text.
// This is used to compute the visual display width of a line that may contain
// color codes, bold, underline, etc.
//
// Supported sequence types:
//   - CSI sequences: ESC [ <params> <final>  (e.g. \x1b[31m, \x1b[2A, \x1b[K)
//   - OSC sequences: ESC ] ... BEL
//   - Other two-byte sequences: ESC <single-char>
func stripANSI(s string) string { //nolint:gocognit,gocyclo,cyclop // byte-level ANSI escape parser with CSI/OSC/two-byte branches
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		c := s[i]
		if c != '\x1b' {
			b.WriteByte(c)
			i++
			continue
		}
		// ESC found. Peek at the next byte to determine sequence type.
		i++ // consume ESC
		if i >= len(s) {
			break
		}
		next := s[i]
		switch next {
		case '[': // CSI: consume until final byte in range 0x40-0x7E
			i++ // consume '['
			for i < len(s) {
				fc := s[i]
				i++
				if fc >= 0x40 && fc <= 0x7E {
					break // final byte consumed
				}
			}
		case ']': // OSC: consume until BEL (0x07) or ESC (start of ST)
			i++ // consume ']'
			for i < len(s) {
				fc := s[i]
				i++
				if fc == 0x07 {
					break
				}
				if fc == '\x1b' {
					// ST = ESC \ — consume the '\'
					if i < len(s) && s[i] == '\\' {
						i++
					}
					break
				}
			}
		default:
			// Two-byte sequence: ESC <char> — skip the next char
			i++
		}
	}
	return b.String()
}

// runeDisplayWidth returns the display width of a single rune.
// Wide characters (CJK, fullwidth) count as 2 columns; most others count as 1.
// Control characters and combining marks count as 0.
func runeDisplayWidth(r rune) int { //nolint:gocognit,gocyclo,cyclop,funlen // Unicode width lookup table across 18 codepoint ranges
	// Control characters occupy 0 visible columns.
	if r < 0x20 || (r >= 0x7F && r < 0xA0) {
		return 0
	}
	// East Asian Wide and Fullwidth ranges that occupy 2 columns.
	switch {
	case r >= 0x1100 && r <= 0x115F: // Hangul Jamo
		return 2
	case r == 0x2329 || r == 0x232A: // Angle brackets
		return 2
	case r >= 0x2E80 && r <= 0x303E: // CJK Radicals, etc.
		return 2
	case r >= 0x3040 && r <= 0x33FF: // Japanese, Katakana, Bopomofo, etc.
		return 2
	case r >= 0x3400 && r <= 0x4DBF: // CJK Extension A
		return 2
	case r >= 0x4E00 && r <= 0xA4CF: // CJK Unified Ideographs
		return 2
	case r >= 0xA960 && r <= 0xA97F: // Hangul Jamo Extended-A
		return 2
	case r >= 0xAC00 && r <= 0xD7FF: // Hangul Syllables + Jamo Extended-B
		return 2
	case r >= 0xF900 && r <= 0xFAFF: // CJK Compatibility Ideographs
		return 2
	case r >= 0xFE10 && r <= 0xFE1F: // Vertical forms
		return 2
	case r >= 0xFE30 && r <= 0xFE6F: // CJK Compatibility Forms
		return 2
	case r >= 0xFF00 && r <= 0xFF60: // Fullwidth Forms
		return 2
	case r >= 0xFFE0 && r <= 0xFFE6: // Fullwidth Signs
		return 2
	case r >= 0x1B000 && r <= 0x1B0FF: // Kana Supplement
		return 2
	case r >= 0x1F004 && r <= 0x1F0CF: // Playing cards, Mahjong
		return 2
	case r >= 0x1F300 && r <= 0x1F9FF: // Miscellaneous Symbols + Emoji
		return 2
	case r >= 0x20000 && r <= 0x2FFFD: // CJK Extension B-F
		return 2
	case r >= 0x30000 && r <= 0x3FFFD: // CJK Extension G+
		return 2
	}
	// Combining / zero-width characters.
	if unicode.In(r, unicode.Mn, unicode.Me, unicode.Cf) {
		return 0
	}
	return 1
}

// visualWidth returns the number of terminal columns that s occupies,
// ignoring any ANSI escape sequences.
func visualWidth(s string) int {
	plain := stripANSI(s)
	w := 0
	for _, r := range plain {
		w += runeDisplayWidth(r)
	}
	return w
}

// truncateLine truncates s so that its visual display width does not exceed
// maxWidth columns. ANSI escape sequences are preserved in the output but
// do not count toward the width budget. If the line fits within maxWidth
// the original string is returned unchanged.
func truncateLine(s string, maxWidth int) string { //nolint:gocognit,gocyclo,cyclop,funlen // ANSI-preserving truncation with CSI/OSC/two-byte parsing
	if maxWidth <= 0 || visualWidth(s) <= maxWidth {
		return s
	}

	var b strings.Builder
	b.Grow(len(s))
	width := 0
	i := 0

	for i < len(s) {
		c := s[i]

		// Pass ANSI escape sequences through without counting columns.
		if c == '\x1b' { //nolint:nestif // ANSI escape parsing requires nested CSI/OSC/ST branches
			b.WriteByte(c)
			i++ // consume ESC
			if i >= len(s) {
				break
			}
			next := s[i]
			b.WriteByte(next)
			i++ // consume byte after ESC
			switch next {
			case '[': // CSI: consume until final byte in 0x40-0x7E
				for i < len(s) {
					fc := s[i]
					b.WriteByte(fc)
					i++
					if fc >= 0x40 && fc <= 0x7E {
						break
					}
				}
			case ']': // OSC: consume until BEL or ESC
				for i < len(s) {
					fc := s[i]
					b.WriteByte(fc)
					i++
					if fc == 0x07 {
						break
					}
					if fc == '\x1b' {
						if i < len(s) && s[i] == '\\' {
							b.WriteByte(s[i])
							i++
						}
						break
					}
				}
				// default: two-byte sequence already consumed above
			}
			continue
		}

		// Decode rune and check remaining width budget.
		r, size := decodeRune(s, i)
		rw := runeDisplayWidth(r)
		if width+rw > maxWidth {
			break
		}
		b.WriteRune(r)
		width += rw
		i += size
	}

	return b.String()
}

// decodeRune decodes a UTF-8 rune from s starting at byte offset i.
// Returns the rune and the number of bytes consumed.
func decodeRune(s string, i int) (rune, int) {
	b := s[i]
	// Single-byte ASCII.
	if b < 0x80 {
		return rune(b), 1
	}
	// Multi-byte sequence: determine length from leading byte.
	var size int
	switch {
	case b&0xE0 == 0xC0:
		size = 2
	case b&0xF0 == 0xE0:
		size = 3
	case b&0xF8 == 0xF0:
		size = 4
	default:
		// Invalid lead byte — treat as replacement character.
		return '\uFFFD', 1
	}
	if i+size > len(s) {
		return '\uFFFD', 1
	}
	r := rune(b & (0xFF >> (size + 1)))
	for j := 1; j < size; j++ {
		cb := s[i+j]
		if cb&0xC0 != 0x80 {
			return '\uFFFD', 1
		}
		r = (r << 6) | rune(cb&0x3F)
	}
	return r, size
}
