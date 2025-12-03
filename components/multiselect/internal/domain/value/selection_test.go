package value

import (
	"reflect"
	"testing"
)

func TestNewSelection(t *testing.T) {
	tests := []struct {
		name     string
		min      int
		max      int
		wantMin  int
		wantMax  int
		wantLen  int
	}{
		{"no constraints", 0, 0, 0, 0, 0},
		{"min only", 2, 0, 2, 0, 0},
		{"max only", 0, 5, 0, 5, 0},
		{"min and max", 1, 3, 1, 3, 0},
		{"negative min", -5, 10, 0, 10, 0},
		{"negative max", 2, -5, 2, 0, 0},
		{"min > max (should clamp)", 10, 5, 5, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSelection(tt.min, tt.max)
			if s.Min() != tt.wantMin {
				t.Errorf("Min() = %d, want %d", s.Min(), tt.wantMin)
			}
			if s.Max() != tt.wantMax {
				t.Errorf("Max() = %d, want %d", s.Max(), tt.wantMax)
			}
			if s.Count() != tt.wantLen {
				t.Errorf("Count() = %d, want %d", s.Count(), tt.wantLen)
			}
		})
	}
}

func TestSelection_WithSelected(t *testing.T) {
	tests := []struct {
		name        string
		initial     *Selection
		indices     []int
		wantIndices []int
	}{
		{
			name:        "add to empty",
			initial:     NewSelection(0, 0),
			indices:     []int{0, 2, 4},
			wantIndices: []int{0, 2, 4},
		},
		{
			name:        "add to existing",
			initial:     NewSelection(0, 0).WithSelected(1, 3),
			indices:     []int{0, 2},
			wantIndices: []int{0, 1, 2, 3},
		},
		{
			name:        "ignore negative",
			initial:     NewSelection(0, 0),
			indices:     []int{-1, 0, -5, 2},
			wantIndices: []int{0, 2},
		},
		{
			name:        "respect max constraint",
			initial:     NewSelection(0, 3),
			indices:     []int{0, 1, 2, 3, 4},
			wantIndices: []int{0, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.initial.WithSelected(tt.indices...)
			got := s.Indices()
			if !reflect.DeepEqual(got, tt.wantIndices) {
				t.Errorf("Indices() = %v, want %v", got, tt.wantIndices)
			}
		})
	}
}

func TestSelection_Toggle(t *testing.T) {
	tests := []struct {
		name        string
		initial     *Selection
		toggle      int
		wantIndices []int
	}{
		{
			name:        "toggle on empty",
			initial:     NewSelection(0, 0),
			toggle:      2,
			wantIndices: []int{2},
		},
		{
			name:        "toggle off",
			initial:     NewSelection(0, 0).WithSelected(1, 2, 3),
			toggle:      2,
			wantIndices: []int{1, 3},
		},
		{
			name:        "toggle on with existing",
			initial:     NewSelection(0, 0).WithSelected(1, 3),
			toggle:      2,
			wantIndices: []int{1, 2, 3},
		},
		{
			name:        "toggle negative (no change)",
			initial:     NewSelection(0, 0).WithSelected(1, 2),
			toggle:      -1,
			wantIndices: []int{1, 2},
		},
		{
			name:        "toggle on respects max",
			initial:     NewSelection(0, 2).WithSelected(1, 3),
			toggle:      5,
			wantIndices: []int{1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.initial.Toggle(tt.toggle)
			got := s.Indices()
			if !reflect.DeepEqual(got, tt.wantIndices) {
				t.Errorf("Indices() = %v, want %v", got, tt.wantIndices)
			}
		})
	}
}

func TestSelection_SelectAll(t *testing.T) {
	tests := []struct {
		name        string
		initial     *Selection
		maxIndex    int
		wantIndices []int
	}{
		{
			name:        "select all no constraints",
			initial:     NewSelection(0, 0),
			maxIndex:    4,
			wantIndices: []int{0, 1, 2, 3, 4},
		},
		{
			name:        "select all with max constraint",
			initial:     NewSelection(0, 3),
			maxIndex:    10,
			wantIndices: []int{0, 1, 2},
		},
		{
			name:        "select all negative maxIndex",
			initial:     NewSelection(0, 0),
			maxIndex:    -1,
			wantIndices: []int{},
		},
		{
			name:        "select all zero maxIndex",
			initial:     NewSelection(0, 0),
			maxIndex:    0,
			wantIndices: []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.initial.SelectAll(tt.maxIndex)
			got := s.Indices()
			if !reflect.DeepEqual(got, tt.wantIndices) {
				t.Errorf("Indices() = %v, want %v", got, tt.wantIndices)
			}
		})
	}
}

func TestSelection_Clear(t *testing.T) {
	s := NewSelection(0, 0).WithSelected(1, 2, 3, 4, 5)
	s = s.Clear()

	if s.Count() != 0 {
		t.Errorf("Count() = %d, want 0", s.Count())
	}

	indices := s.Indices()
	if len(indices) != 0 {
		t.Errorf("Indices() = %v, want []", indices)
	}

	// Constraints should be preserved
	if s.Min() != 0 || s.Max() != 0 {
		t.Errorf("Min/Max = %d/%d, want 0/0", s.Min(), s.Max())
	}
}

func TestSelection_IsSelected(t *testing.T) {
	s := NewSelection(0, 0).WithSelected(1, 3, 5)

	tests := []struct {
		index int
		want  bool
	}{
		{0, false},
		{1, true},
		{2, false},
		{3, true},
		{4, false},
		{5, true},
		{-1, false},
		{100, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := s.IsSelected(tt.index); got != tt.want {
				t.Errorf("IsSelected(%d) = %v, want %v", tt.index, got, tt.want)
			}
		})
	}
}

func TestSelection_Count(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *Selection
		want  int
	}{
		{
			name:  "empty",
			setup: func() *Selection { return NewSelection(0, 0) },
			want:  0,
		},
		{
			name:  "single",
			setup: func() *Selection { return NewSelection(0, 0).WithSelected(5) },
			want:  1,
		},
		{
			name:  "multiple",
			setup: func() *Selection { return NewSelection(0, 0).WithSelected(1, 2, 3, 5, 8) },
			want:  5,
		},
		{
			name: "after toggle off",
			setup: func() *Selection {
				return NewSelection(0, 0).WithSelected(1, 2, 3).Toggle(2)
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			if got := s.Count(); got != tt.want {
				t.Errorf("Count() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSelection_Indices(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Selection
		want    []int
		wantLen int
	}{
		{
			name:    "empty",
			setup:   func() *Selection { return NewSelection(0, 0) },
			want:    []int{},
			wantLen: 0,
		},
		{
			name:    "sorted ascending",
			setup:   func() *Selection { return NewSelection(0, 0).WithSelected(5, 1, 3, 2) },
			want:    []int{1, 2, 3, 5},
			wantLen: 4,
		},
		{
			name:    "single",
			setup:   func() *Selection { return NewSelection(0, 0).WithSelected(42) },
			want:    []int{42},
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			got := s.Indices()
			if len(got) != tt.wantLen {
				t.Errorf("len(Indices()) = %d, want %d", len(got), tt.wantLen)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Indices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelection_CanSelect(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Selection
		want    bool
	}{
		{
			name:    "no max constraint",
			setup:   func() *Selection { return NewSelection(0, 0).WithSelected(1, 2, 3) },
			want:    true,
		},
		{
			name:    "max not reached",
			setup:   func() *Selection { return NewSelection(0, 5).WithSelected(1, 2) },
			want:    true,
		},
		{
			name:    "max reached",
			setup:   func() *Selection { return NewSelection(0, 3).WithSelected(1, 2, 3) },
			want:    false,
		},
		{
			name:    "empty with max",
			setup:   func() *Selection { return NewSelection(0, 5) },
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			if got := s.CanSelect(); got != tt.want {
				t.Errorf("CanSelect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelection_CanConfirm(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *Selection
		want  bool
	}{
		{
			name:  "no min constraint",
			setup: func() *Selection { return NewSelection(0, 0) },
			want:  true,
		},
		{
			name:  "min reached",
			setup: func() *Selection { return NewSelection(2, 0).WithSelected(1, 2) },
			want:  true,
		},
		{
			name:  "min exceeded",
			setup: func() *Selection { return NewSelection(2, 0).WithSelected(1, 2, 3) },
			want:  true,
		},
		{
			name:  "min not reached",
			setup: func() *Selection { return NewSelection(3, 0).WithSelected(1) },
			want:  false,
		},
		{
			name:  "empty with min",
			setup: func() *Selection { return NewSelection(1, 0) },
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			if got := s.CanConfirm(); got != tt.want {
				t.Errorf("CanConfirm() = %v, want %v", got, tt.want)
			}
		})
	}
}
