// Package model provides rich domain models for textarea.
package model

import "github.com/phoenix-tui/phoenix/components/input/textarea/domain/value"

// Selection represents selected text range.
// This is a rich domain model that encapsulates selection behavior.
// All operations return new instances (immutable).
type Selection struct {
	anchor value.Position // Selection start (where selection began)
	cursor value.Position // Selection end (current cursor position)
}

// NewSelection creates selection from anchor to cursor.
func NewSelection(anchor, cursor value.Position) *Selection {
	return &Selection{
		anchor: anchor,
		cursor: cursor,
	}
}

// Range returns selection range (normalized so start <= end).
func (s *Selection) Range() value.Range {
	// Normalize so start <= end.
	if s.anchor.IsBefore(s.cursor) {
		return value.NewRange(s.anchor, s.cursor)
	}
	return value.NewRange(s.cursor, s.anchor)
}

// Anchor returns anchor position.
func (s *Selection) Anchor() value.Position {
	return s.anchor
}

// Cursor returns cursor position.
func (s *Selection) Cursor() value.Position {
	return s.cursor
}

// WithCursor returns new selection with updated cursor position.
func (s *Selection) WithCursor(cursor value.Position) *Selection {
	return &Selection{
		anchor: s.anchor,
		cursor: cursor,
	}
}

// Copy returns copy of selection (nil-safe).
func (s *Selection) Copy() *Selection {
	if s == nil {
		return nil
	}
	return &Selection{
		anchor: s.anchor,
		cursor: s.cursor,
	}
}
