// Package value contains value objects for input domain.
package value

// Cursor represents a grapheme-aware cursor position in text.
// The offset is measured in grapheme clusters, not bytes or runes.
// This ensures proper handling of emoji, combining characters, and CJK text.
type Cursor struct {
	offset int // Grapheme cluster offset (NOT byte offset!)
}

// NewCursor creates a new cursor at the specified grapheme offset.
// Negative offsets are clamped to 0.
func NewCursor(offset int) *Cursor {
	if offset < 0 {
		offset = 0
	}
	return &Cursor{offset: offset}
}

// Offset returns the current grapheme offset.
func (c *Cursor) Offset() int {
	return c.offset
}

// MoveBy moves the cursor by the specified delta.
// The result is clamped to [0, maxOffset].
// Returns a new Cursor instance (immutable).
func (c *Cursor) MoveBy(delta, maxOffset int) *Cursor {
	newOffset := c.offset + delta
	if newOffset < 0 {
		newOffset = 0
	}
	if newOffset > maxOffset {
		newOffset = maxOffset
	}
	return &Cursor{offset: newOffset}
}

// MoveTo moves the cursor to the specified absolute position.
// The result is clamped to [0, maxOffset].
// Returns a new Cursor instance (immutable).
func (c *Cursor) MoveTo(offset, maxOffset int) *Cursor {
	if offset < 0 {
		offset = 0
	}
	if offset > maxOffset {
		offset = maxOffset
	}
	return &Cursor{offset: offset}
}

// Clone creates a copy of this cursor.
func (c *Cursor) Clone() *Cursor {
	return &Cursor{offset: c.offset}
}
