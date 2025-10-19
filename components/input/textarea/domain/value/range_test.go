package value

import (
	"testing"
)

func TestNewRange(t *testing.T) {
	tests := []struct {
		name       string
		start      Position
		end        Position
		wantStart  Position
		wantEnd    Position
		normalized bool
	}{
		{
			name:       "forward range",
			start:      NewPosition(0, 0),
			end:        NewPosition(0, 5),
			wantStart:  NewPosition(0, 0),
			wantEnd:    NewPosition(0, 5),
			normalized: false,
		},
		{
			name:       "backward range gets normalized",
			start:      NewPosition(0, 10),
			end:        NewPosition(0, 3),
			wantStart:  NewPosition(0, 3),
			wantEnd:    NewPosition(0, 10),
			normalized: true,
		},
		{
			name:       "multiline forward",
			start:      NewPosition(1, 5),
			end:        NewPosition(3, 8),
			wantStart:  NewPosition(1, 5),
			wantEnd:    NewPosition(3, 8),
			normalized: false,
		},
		{
			name:       "multiline backward gets normalized",
			start:      NewPosition(5, 10),
			end:        NewPosition(2, 3),
			wantStart:  NewPosition(2, 3),
			wantEnd:    NewPosition(5, 10),
			normalized: true,
		},
		{
			name:       "same position",
			start:      NewPosition(2, 5),
			end:        NewPosition(2, 5),
			wantStart:  NewPosition(2, 5),
			wantEnd:    NewPosition(2, 5),
			normalized: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRange(tt.start, tt.end)

			if !r.Start().Equals(tt.wantStart) {
				t.Errorf("Start() = %v, want %v", r.Start(), tt.wantStart)
			}
			if !r.End().Equals(tt.wantEnd) {
				t.Errorf("End() = %v, want %v", r.End(), tt.wantEnd)
			}

			// Verify start is always before or equal to end
			if r.Start().IsAfter(r.End()) {
				t.Error("Range not normalized: start is after end")
			}
		})
	}
}

func TestRange_Start(t *testing.T) {
	tests := []struct {
		name  string
		start Position
		end   Position
		want  Position
	}{
		{
			name:  "forward range",
			start: NewPosition(1, 5),
			end:   NewPosition(3, 8),
			want:  NewPosition(1, 5),
		},
		{
			name:  "backward range normalized",
			start: NewPosition(5, 10),
			end:   NewPosition(2, 3),
			want:  NewPosition(2, 3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRange(tt.start, tt.end)
			if !r.Start().Equals(tt.want) {
				t.Errorf("Start() = %v, want %v", r.Start(), tt.want)
			}
		})
	}
}

func TestRange_End(t *testing.T) {
	tests := []struct {
		name  string
		start Position
		end   Position
		want  Position
	}{
		{
			name:  "forward range",
			start: NewPosition(1, 5),
			end:   NewPosition(3, 8),
			want:  NewPosition(3, 8),
		},
		{
			name:  "backward range normalized",
			start: NewPosition(5, 10),
			end:   NewPosition(2, 3),
			want:  NewPosition(5, 10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRange(tt.start, tt.end)
			if !r.End().Equals(tt.want) {
				t.Errorf("End() = %v, want %v", r.End(), tt.want)
			}
		})
	}
}

func TestRange_StartRowCol(t *testing.T) {
	tests := []struct {
		name    string
		start   Position
		end     Position
		wantRow int
		wantCol int
	}{
		{
			name:    "simple range",
			start:   NewPosition(1, 5),
			end:     NewPosition(3, 8),
			wantRow: 1,
			wantCol: 5,
		},
		{
			name:    "backward range normalized",
			start:   NewPosition(5, 10),
			end:     NewPosition(2, 3),
			wantRow: 2,
			wantCol: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRange(tt.start, tt.end)
			row, col := r.StartRowCol()

			if row != tt.wantRow {
				t.Errorf("StartRowCol() row = %d, want %d", row, tt.wantRow)
			}
			if col != tt.wantCol {
				t.Errorf("StartRowCol() col = %d, want %d", col, tt.wantCol)
			}
		})
	}
}

func TestRange_EndRowCol(t *testing.T) {
	tests := []struct {
		name    string
		start   Position
		end     Position
		wantRow int
		wantCol int
	}{
		{
			name:    "simple range",
			start:   NewPosition(1, 5),
			end:     NewPosition(3, 8),
			wantRow: 3,
			wantCol: 8,
		},
		{
			name:    "backward range normalized",
			start:   NewPosition(5, 10),
			end:     NewPosition(2, 3),
			wantRow: 5,
			wantCol: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRange(tt.start, tt.end)
			row, col := r.EndRowCol()

			if row != tt.wantRow {
				t.Errorf("EndRowCol() row = %d, want %d", row, tt.wantRow)
			}
			if col != tt.wantCol {
				t.Errorf("EndRowCol() col = %d, want %d", col, tt.wantCol)
			}
		})
	}
}

func TestRange_Contains(t *testing.T) {
	r := NewRange(NewPosition(2, 5), NewPosition(5, 10))

	tests := []struct {
		name string
		pos  Position
		want bool
	}{
		{
			name: "position before range",
			pos:  NewPosition(1, 0),
			want: false,
		},
		{
			name: "position at start",
			pos:  NewPosition(2, 5),
			want: true,
		},
		{
			name: "position inside range",
			pos:  NewPosition(3, 7),
			want: true,
		},
		{
			name: "position at end",
			pos:  NewPosition(5, 10),
			want: true,
		},
		{
			name: "position after range",
			pos:  NewPosition(6, 0),
			want: false,
		},
		{
			name: "position before range same row as start",
			pos:  NewPosition(2, 3),
			want: false,
		},
		{
			name: "position after range same row as end",
			pos:  NewPosition(5, 15),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Contains(tt.pos)
			if result != tt.want {
				t.Errorf("Contains(%v) = %v, want %v", tt.pos, result, tt.want)
			}
		})
	}
}

func TestRange_IsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		start Position
		end   Position
		want  bool
	}{
		{
			name:  "empty range same position",
			start: NewPosition(2, 5),
			end:   NewPosition(2, 5),
			want:  true,
		},
		{
			name:  "non-empty range different col",
			start: NewPosition(2, 5),
			end:   NewPosition(2, 10),
			want:  false,
		},
		{
			name:  "non-empty range different row",
			start: NewPosition(2, 5),
			end:   NewPosition(3, 5),
			want:  false,
		},
		{
			name:  "non-empty multiline",
			start: NewPosition(1, 0),
			end:   NewPosition(5, 10),
			want:  false,
		},
		{
			name:  "empty range at origin",
			start: NewPosition(0, 0),
			end:   NewPosition(0, 0),
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRange(tt.start, tt.end)
			result := r.IsEmpty()

			if result != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestRange_IsSingleLine(t *testing.T) {
	tests := []struct {
		name  string
		start Position
		end   Position
		want  bool
	}{
		{
			name:  "single line range",
			start: NewPosition(2, 5),
			end:   NewPosition(2, 10),
			want:  true,
		},
		{
			name:  "multiline range",
			start: NewPosition(2, 5),
			end:   NewPosition(3, 5),
			want:  false,
		},
		{
			name:  "empty range is single line",
			start: NewPosition(5, 10),
			end:   NewPosition(5, 10),
			want:  true,
		},
		{
			name:  "large multiline range",
			start: NewPosition(1, 0),
			end:   NewPosition(100, 0),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRange(tt.start, tt.end)
			result := r.IsSingleLine()

			if result != tt.want {
				t.Errorf("IsSingleLine() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestRange_Normalization(t *testing.T) {
	tests := []struct {
		name        string
		start       Position
		end         Position
		description string
	}{
		{
			name:        "already normalized forward",
			start:       NewPosition(1, 5),
			end:         NewPosition(3, 8),
			description: "should remain unchanged",
		},
		{
			name:        "backward needs normalization",
			start:       NewPosition(3, 8),
			end:         NewPosition(1, 5),
			description: "should swap positions",
		},
		{
			name:        "same row backward",
			start:       NewPosition(2, 10),
			end:         NewPosition(2, 3),
			description: "should swap columns",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRange(tt.start, tt.end)

			// After normalization, start should always be before or equal to end
			if r.Start().IsAfter(r.End()) {
				t.Errorf("Range not normalized: start=%v > end=%v", r.Start(), r.End())
			}

			// The smaller position should become start
			if tt.start.IsBefore(tt.end) {
				if !r.Start().Equals(tt.start) {
					t.Errorf("Forward range: start=%v, want %v", r.Start(), tt.start)
				}
				if !r.End().Equals(tt.end) {
					t.Errorf("Forward range: end=%v, want %v", r.End(), tt.end)
				}
			} else {
				if !r.Start().Equals(tt.end) {
					t.Errorf("Backward range: start=%v, want %v", r.Start(), tt.end)
				}
				if !r.End().Equals(tt.start) {
					t.Errorf("Backward range: end=%v, want %v", r.End(), tt.start)
				}
			}
		})
	}
}

func TestRange_ValueObject(t *testing.T) {
	// Test that Range behaves as a value object
	r1 := NewRange(NewPosition(1, 5), NewPosition(3, 8))
	r2 := NewRange(NewPosition(1, 5), NewPosition(3, 8))
	r3 := NewRange(NewPosition(1, 5), NewPosition(3, 9))

	// Equal ranges should have equal start and end
	if !r1.Start().Equals(r2.Start()) || !r1.End().Equals(r2.End()) {
		t.Error("Equal ranges should have equal start and end")
	}

	// Different ranges should not be equal
	if r1.Start().Equals(r3.Start()) && r1.End().Equals(r3.End()) {
		t.Error("Different ranges should not be equal")
	}
}

func TestRange_EdgeCases(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "range at origin",
			test: func(t *testing.T) {
				r := NewRange(NewPosition(0, 0), NewPosition(0, 5))
				if !r.Start().Equals(NewPosition(0, 0)) {
					t.Error("Range should start at origin")
				}
			},
		},
		{
			name: "zero-width range at various positions",
			test: func(t *testing.T) {
				positions := []Position{
					NewPosition(0, 0),
					NewPosition(5, 10),
					NewPosition(100, 200),
				}
				for _, pos := range positions {
					r := NewRange(pos, pos)
					if !r.IsEmpty() {
						t.Errorf("Zero-width range at %v should be empty", pos)
					}
					if !r.IsSingleLine() {
						t.Errorf("Zero-width range at %v should be single line", pos)
					}
				}
			},
		},
		{
			name: "contains boundary positions",
			test: func(t *testing.T) {
				r := NewRange(NewPosition(2, 5), NewPosition(5, 10))
				// Should contain start
				if !r.Contains(r.Start()) {
					t.Error("Range should contain its start position")
				}
				// Should contain end
				if !r.Contains(r.End()) {
					t.Error("Range should contain its end position")
				}
			},
		},
		{
			name: "single character range",
			test: func(t *testing.T) {
				r := NewRange(NewPosition(2, 5), NewPosition(2, 6))
				if r.IsEmpty() {
					t.Error("Single character range should not be empty")
				}
				if !r.IsSingleLine() {
					t.Error("Single character range should be single line")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestRange_ContainsEdgeCases(t *testing.T) {
	// Test Contains with various edge cases
	r := NewRange(NewPosition(2, 5), NewPosition(4, 10))

	tests := []struct {
		name string
		pos  Position
		want bool
	}{
		// Boundary tests
		{name: "start position", pos: NewPosition(2, 5), want: true},
		{name: "end position", pos: NewPosition(4, 10), want: true},
		{name: "just before start", pos: NewPosition(2, 4), want: false},
		{name: "just after end", pos: NewPosition(4, 11), want: false},

		// Middle tests
		{name: "middle of range", pos: NewPosition(3, 7), want: true},
		{name: "start of middle row", pos: NewPosition(3, 0), want: true},
		{name: "end of middle row", pos: NewPosition(3, 100), want: true},

		// Row tests
		{name: "row before", pos: NewPosition(1, 100), want: false},
		{name: "row after", pos: NewPosition(5, 0), want: false},
		{name: "start row before col", pos: NewPosition(2, 0), want: false},
		{name: "end row after col", pos: NewPosition(4, 100), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.Contains(tt.pos)
			if result != tt.want {
				t.Errorf("Contains(%v) = %v, want %v", tt.pos, result, tt.want)
			}
		})
	}
}
