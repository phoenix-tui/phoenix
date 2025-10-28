package model

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/input/internal/textarea/domain/value"
)

func TestNewSelection(t *testing.T) {
	tests := []struct {
		name       string
		anchor     value.Position
		cursor     value.Position
		wantAnchor value.Position
		wantCursor value.Position
	}{
		{
			name:       "forward selection",
			anchor:     value.NewPosition(0, 0),
			cursor:     value.NewPosition(0, 5),
			wantAnchor: value.NewPosition(0, 0),
			wantCursor: value.NewPosition(0, 5),
		},
		{
			name:       "backward selection",
			anchor:     value.NewPosition(0, 10),
			cursor:     value.NewPosition(0, 3),
			wantAnchor: value.NewPosition(0, 10),
			wantCursor: value.NewPosition(0, 3),
		},
		{
			name:       "multiline forward",
			anchor:     value.NewPosition(1, 5),
			cursor:     value.NewPosition(3, 8),
			wantAnchor: value.NewPosition(1, 5),
			wantCursor: value.NewPosition(3, 8),
		},
		{
			name:       "multiline backward",
			anchor:     value.NewPosition(5, 10),
			cursor:     value.NewPosition(2, 3),
			wantAnchor: value.NewPosition(5, 10),
			wantCursor: value.NewPosition(2, 3),
		},
		{
			name:       "same position",
			anchor:     value.NewPosition(2, 5),
			cursor:     value.NewPosition(2, 5),
			wantAnchor: value.NewPosition(2, 5),
			wantCursor: value.NewPosition(2, 5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel := NewSelection(tt.anchor, tt.cursor)

			if !sel.Anchor().Equals(tt.wantAnchor) {
				t.Errorf("Anchor() = %v, want %v", sel.Anchor(), tt.wantAnchor)
			}
			if !sel.Cursor().Equals(tt.wantCursor) {
				t.Errorf("Cursor() = %v, want %v", sel.Cursor(), tt.wantCursor)
			}
		})
	}
}

func TestSelection_Range(t *testing.T) {
	tests := []struct {
		name      string
		anchor    value.Position
		cursor    value.Position
		wantStart value.Position
		wantEnd   value.Position
	}{
		{
			name:      "forward selection normalized",
			anchor:    value.NewPosition(0, 0),
			cursor:    value.NewPosition(0, 5),
			wantStart: value.NewPosition(0, 0),
			wantEnd:   value.NewPosition(0, 5),
		},
		{
			name:      "backward selection normalized",
			anchor:    value.NewPosition(0, 10),
			cursor:    value.NewPosition(0, 3),
			wantStart: value.NewPosition(0, 3),
			wantEnd:   value.NewPosition(0, 10),
		},
		{
			name:      "multiline forward normalized",
			anchor:    value.NewPosition(1, 5),
			cursor:    value.NewPosition(3, 8),
			wantStart: value.NewPosition(1, 5),
			wantEnd:   value.NewPosition(3, 8),
		},
		{
			name:      "multiline backward normalized",
			anchor:    value.NewPosition(5, 10),
			cursor:    value.NewPosition(2, 3),
			wantStart: value.NewPosition(2, 3),
			wantEnd:   value.NewPosition(5, 10),
		},
		{
			name:      "same position",
			anchor:    value.NewPosition(2, 5),
			cursor:    value.NewPosition(2, 5),
			wantStart: value.NewPosition(2, 5),
			wantEnd:   value.NewPosition(2, 5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel := NewSelection(tt.anchor, tt.cursor)
			r := sel.Range()

			if !r.Start().Equals(tt.wantStart) {
				t.Errorf("Range().Start() = %v, want %v", r.Start(), tt.wantStart)
			}
			if !r.End().Equals(tt.wantEnd) {
				t.Errorf("Range().End() = %v, want %v", r.End(), tt.wantEnd)
			}
		})
	}
}

func TestSelection_Anchor(t *testing.T) {
	anchor := value.NewPosition(5, 10)
	cursor := value.NewPosition(3, 7)
	sel := NewSelection(anchor, cursor)

	result := sel.Anchor()
	if !result.Equals(anchor) {
		t.Errorf("Anchor() = %v, want %v", result, anchor)
	}
}

func TestSelection_Cursor(t *testing.T) {
	anchor := value.NewPosition(5, 10)
	cursor := value.NewPosition(3, 7)
	sel := NewSelection(anchor, cursor)

	result := sel.Cursor()
	if !result.Equals(cursor) {
		t.Errorf("Cursor() = %v, want %v", result, cursor)
	}
}

func TestSelection_WithCursor(t *testing.T) {
	tests := []struct {
		name          string
		anchor        value.Position
		initialCursor value.Position
		newCursor     value.Position
		wantAnchor    value.Position
		wantCursor    value.Position
	}{
		{
			name:          "update cursor forward",
			anchor:        value.NewPosition(0, 0),
			initialCursor: value.NewPosition(0, 5),
			newCursor:     value.NewPosition(0, 10),
			wantAnchor:    value.NewPosition(0, 0),
			wantCursor:    value.NewPosition(0, 10),
		},
		{
			name:          "update cursor backward",
			anchor:        value.NewPosition(0, 10),
			initialCursor: value.NewPosition(0, 5),
			newCursor:     value.NewPosition(0, 3),
			wantAnchor:    value.NewPosition(0, 10),
			wantCursor:    value.NewPosition(0, 3),
		},
		{
			name:          "update cursor to anchor",
			anchor:        value.NewPosition(2, 5),
			initialCursor: value.NewPosition(3, 8),
			newCursor:     value.NewPosition(2, 5),
			wantAnchor:    value.NewPosition(2, 5),
			wantCursor:    value.NewPosition(2, 5),
		},
		{
			name:          "update cursor multiline",
			anchor:        value.NewPosition(1, 0),
			initialCursor: value.NewPosition(1, 10),
			newCursor:     value.NewPosition(5, 20),
			wantAnchor:    value.NewPosition(1, 0),
			wantCursor:    value.NewPosition(5, 20),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel := NewSelection(tt.anchor, tt.initialCursor)
			result := sel.WithCursor(tt.newCursor)

			if !result.Anchor().Equals(tt.wantAnchor) {
				t.Errorf("WithCursor().Anchor() = %v, want %v", result.Anchor(), tt.wantAnchor)
			}
			if !result.Cursor().Equals(tt.wantCursor) {
				t.Errorf("WithCursor().Cursor() = %v, want %v", result.Cursor(), tt.wantCursor)
			}

			// Verify immutability.
			if !sel.Cursor().Equals(tt.initialCursor) {
				t.Errorf("Original selection cursor was modified: %v, want %v", sel.Cursor(), tt.initialCursor)
			}
			if !sel.Anchor().Equals(tt.anchor) {
				t.Errorf("Original selection anchor was modified: %v, want %v", sel.Anchor(), tt.anchor)
			}
		})
	}
}

func TestSelection_Copy(t *testing.T) {
	tests := []struct {
		name  string
		sel   *Selection
		isNil bool
	}{
		{
			name:  "copy non-nil selection",
			sel:   NewSelection(value.NewPosition(1, 5), value.NewPosition(3, 8)),
			isNil: false,
		},
		{
			name:  "copy nil selection",
			sel:   nil,
			isNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sel.Copy()

			if tt.isNil {
				if result != nil {
					t.Error("Copy() of nil selection should return nil")
				}
				return
			}

			if result == nil {
				t.Error("Copy() should not return nil for non-nil selection")
				return
			}

			// Verify copy has same values.
			if !result.Anchor().Equals(tt.sel.Anchor()) {
				t.Errorf("Copy().Anchor() = %v, want %v", result.Anchor(), tt.sel.Anchor())
			}
			if !result.Cursor().Equals(tt.sel.Cursor()) {
				t.Errorf("Copy().Cursor() = %v, want %v", result.Cursor(), tt.sel.Cursor())
			}

			// Verify it's a new instance.
			if result == tt.sel {
				t.Error("Copy() returned same instance, want new instance")
			}

			// Verify modifying copy doesn't affect original.
			newCursor := value.NewPosition(999, 888)
			newCopy := result.WithCursor(newCursor)

			if !tt.sel.Cursor().Equals(value.NewPosition(3, 8)) {
				t.Error("Original selection was modified")
			}
			if !newCopy.Cursor().Equals(newCursor) {
				t.Error("Modified copy has incorrect cursor")
			}
		})
	}
}

func TestSelection_Immutability(t *testing.T) {
	anchor := value.NewPosition(1, 5)
	cursor := value.NewPosition(3, 8)
	original := NewSelection(anchor, cursor)

	// Test all mutation operations preserve original.
	withNewCursor := original.WithCursor(value.NewPosition(5, 10))
	copied := original.Copy()

	// Original should remain unchanged.
	if !original.Anchor().Equals(anchor) {
		t.Errorf("Original anchor changed: %v, want %v", original.Anchor(), anchor)
	}
	if !original.Cursor().Equals(cursor) {
		t.Errorf("Original cursor changed: %v, want %v", original.Cursor(), cursor)
	}

	// Results should have correct values.
	if !withNewCursor.Cursor().Equals(value.NewPosition(5, 10)) {
		t.Error("WithCursor() produced incorrect selection")
	}
	if !withNewCursor.Anchor().Equals(anchor) {
		t.Error("WithCursor() modified anchor")
	}
	if !copied.Anchor().Equals(anchor) || !copied.Cursor().Equals(cursor) {
		t.Error("Copy() produced incorrect selection")
	}
}

func TestSelection_RangeNormalization(t *testing.T) {
	tests := []struct {
		name        string
		anchor      value.Position
		cursor      value.Position
		description string
	}{
		{
			name:        "forward selection should maintain order",
			anchor:      value.NewPosition(0, 0),
			cursor:      value.NewPosition(0, 10),
			description: "start=anchor, end=cursor",
		},
		{
			name:        "backward selection should swap order",
			anchor:      value.NewPosition(0, 10),
			cursor:      value.NewPosition(0, 0),
			description: "start=cursor, end=anchor",
		},
		{
			name:        "multiline forward should maintain order",
			anchor:      value.NewPosition(1, 5),
			cursor:      value.NewPosition(5, 10),
			description: "start=anchor, end=cursor",
		},
		{
			name:        "multiline backward should swap order",
			anchor:      value.NewPosition(5, 10),
			cursor:      value.NewPosition(1, 5),
			description: "start=cursor, end=anchor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel := NewSelection(tt.anchor, tt.cursor)
			r := sel.Range()

			// Range should always have start <= end.
			if r.Start().IsAfter(r.End()) {
				t.Errorf("Range not normalized: start=%v > end=%v", r.Start(), r.End())
			}

			// For forward selections, start should be anchor.
			//nolint:nestif // test validation logic requires branching
			if tt.anchor.IsBefore(tt.cursor) {
				if !r.Start().Equals(tt.anchor) {
					t.Errorf("Forward selection: start=%v, want %v (anchor)", r.Start(), tt.anchor)
				}
				if !r.End().Equals(tt.cursor) {
					t.Errorf("Forward selection: end=%v, want %v (cursor)", r.End(), tt.cursor)
				}
			} else {
				// For backward selections, start should be cursor.
				if !r.Start().Equals(tt.cursor) {
					t.Errorf("Backward selection: start=%v, want %v (cursor)", r.Start(), tt.cursor)
				}
				if !r.End().Equals(tt.anchor) {
					t.Errorf("Backward selection: end=%v, want %v (anchor)", r.End(), tt.anchor)
				}
			}
		})
	}
}
