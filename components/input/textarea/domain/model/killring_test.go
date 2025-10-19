package model

import (
	"testing"
)

func TestNewKillRing(t *testing.T) {
	tests := []struct {
		name        string
		maxSize     int
		wantMaxSize int
	}{
		{
			name:        "positive max size",
			maxSize:     10,
			wantMaxSize: 10,
		},
		{
			name:        "large max size",
			maxSize:     100,
			wantMaxSize: 100,
		},
		{
			name:        "zero max size defaults to 10",
			maxSize:     0,
			wantMaxSize: 10,
		},
		{
			name:        "negative max size defaults to 10",
			maxSize:     -5,
			wantMaxSize: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kr := NewKillRing(tt.maxSize)

			if !kr.IsEmpty() {
				t.Error("NewKillRing() should be empty")
			}

			// Test that ring accepts items up to maxSize
			testKr := kr
			for i := 0; i < tt.wantMaxSize; i++ {
				testKr = testKr.Kill("item")
			}

			// Ring should now be at max capacity
			// Adding one more should remove oldest
			testKr = testKr.Kill("newest")
			yanked := testKr.Yank()
			if yanked != "newest" {
				t.Errorf("Expected newest item, got %q", yanked)
			}
		})
	}
}

func TestKillRing_Kill(t *testing.T) {
	tests := []struct {
		name     string
		initial  []string
		killText string
		wantLast string
		wantSize int
	}{
		{
			name:     "kill first item",
			initial:  []string{},
			killText: "first",
			wantLast: "first",
			wantSize: 1,
		},
		{
			name:     "kill multiple items",
			initial:  []string{"first"},
			killText: "second",
			wantLast: "second",
			wantSize: 2,
		},
		{
			name:     "kill empty string",
			initial:  []string{"first"},
			killText: "",
			wantLast: "first",
			wantSize: 1,
		},
		{
			name:     "kill unicode text",
			initial:  []string{},
			killText: "hello 世界 👋",
			wantLast: "hello 世界 👋",
			wantSize: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kr := NewKillRing(10)
			for _, item := range tt.initial {
				kr = kr.Kill(item)
			}

			result := kr.Kill(tt.killText)

			if tt.killText != "" {
				yanked := result.Yank()
				if yanked != tt.wantLast {
					t.Errorf("Yank() = %q, want %q", yanked, tt.wantLast)
				}
			}

			// Verify immutability
			originalYank := kr.Yank()
			if len(tt.initial) > 0 && originalYank != tt.initial[len(tt.initial)-1] {
				t.Errorf("Original kill ring was modified")
			}
		})
	}
}

func TestKillRing_Kill_MaxSize(t *testing.T) {
	kr := NewKillRing(3)

	// Fill ring to capacity
	kr = kr.Kill("first")
	kr = kr.Kill("second")
	kr = kr.Kill("third")

	// Verify we have 3 items
	if kr.IsEmpty() {
		t.Error("Ring should not be empty")
	}

	// Add fourth item (should evict first)
	kr = kr.Kill("fourth")

	// Current item should be fourth
	if kr.Yank() != "fourth" {
		t.Errorf("Yank() = %q, want %q", kr.Yank(), "fourth")
	}

	// Rotate back through ring
	kr = kr.YankPop() // third
	if kr.Yank() != "third" {
		t.Errorf("After YankPop(), Yank() = %q, want %q", kr.Yank(), "third")
	}

	kr = kr.YankPop() // second
	if kr.Yank() != "second" {
		t.Errorf("After YankPop(), Yank() = %q, want %q", kr.Yank(), "second")
	}

	kr = kr.YankPop() // fourth (wrapped around)
	if kr.Yank() != "fourth" {
		t.Errorf("After YankPop() wrap, Yank() = %q, want %q", kr.Yank(), "fourth")
	}
}

func TestKillRing_Yank(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		want  string
	}{
		{
			name:  "yank from empty ring",
			items: []string{},
			want:  "",
		},
		{
			name:  "yank single item",
			items: []string{"only"},
			want:  "only",
		},
		{
			name:  "yank latest item",
			items: []string{"first", "second", "third"},
			want:  "third",
		},
		{
			name:  "yank unicode",
			items: []string{"hello", "世界", "👋"},
			want:  "👋",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kr := NewKillRing(10)
			for _, item := range tt.items {
				kr = kr.Kill(item)
			}

			result := kr.Yank()
			if result != tt.want {
				t.Errorf("Yank() = %q, want %q", result, tt.want)
			}
		})
	}
}

func TestKillRing_YankPop(t *testing.T) {
	tests := []struct {
		name    string
		items   []string
		numPops int
		want    string
	}{
		{
			name:    "yank pop from empty ring",
			items:   []string{},
			numPops: 1,
			want:    "",
		},
		{
			name:    "yank pop single item",
			items:   []string{"only"},
			numPops: 1,
			want:    "only",
		},
		{
			name:    "yank pop once",
			items:   []string{"first", "second", "third"},
			numPops: 1,
			want:    "second",
		},
		{
			name:    "yank pop twice",
			items:   []string{"first", "second", "third"},
			numPops: 2,
			want:    "first",
		},
		{
			name:    "yank pop wraps around",
			items:   []string{"first", "second", "third"},
			numPops: 3,
			want:    "third",
		},
		{
			name:    "yank pop multiple wraps",
			items:   []string{"first", "second", "third"},
			numPops: 7, // 3 full rotations + 1
			want:    "second",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kr := NewKillRing(10)
			for _, item := range tt.items {
				kr = kr.Kill(item)
			}

			// Pop numPops times
			for i := 0; i < tt.numPops; i++ {
				kr = kr.YankPop()
			}

			result := kr.Yank()
			if result != tt.want {
				t.Errorf("After %d pops, Yank() = %q, want %q", tt.numPops, result, tt.want)
			}
		})
	}
}

func TestKillRing_IsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		want  bool
	}{
		{
			name:  "empty ring",
			items: []string{},
			want:  true,
		},
		{
			name:  "non-empty ring",
			items: []string{"item"},
			want:  false,
		},
		{
			name:  "ring with multiple items",
			items: []string{"first", "second", "third"},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kr := NewKillRing(10)
			for _, item := range tt.items {
				kr = kr.Kill(item)
			}

			result := kr.IsEmpty()
			if result != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestKillRing_Copy(t *testing.T) {
	kr := NewKillRing(10)
	kr = kr.Kill("first")
	kr = kr.Kill("second")
	kr = kr.Kill("third")

	copy := kr.Copy()

	// Verify copy has same content
	if copy.Yank() != kr.Yank() {
		t.Errorf("Copy() Yank() = %q, want %q", copy.Yank(), kr.Yank())
	}
	if copy.IsEmpty() != kr.IsEmpty() {
		t.Errorf("Copy() IsEmpty() = %v, want %v", copy.IsEmpty(), kr.IsEmpty())
	}

	// Verify it's a new instance
	if copy == kr {
		t.Error("Copy() returned same instance, want new instance")
	}

	// Verify modifying copy doesn't affect original
	copy = copy.Kill("fourth")
	if kr.Yank() != "third" {
		t.Errorf("Original ring was modified: Yank() = %q, want %q", kr.Yank(), "third")
	}
	if copy.Yank() != "fourth" {
		t.Errorf("Copy Yank() = %q, want %q", copy.Yank(), "fourth")
	}

	// Verify modifying index in copy doesn't affect original
	copy = copy.YankPop()
	if kr.Yank() != "third" {
		t.Error("Original ring index was modified")
	}
}

func TestKillRing_Immutability(t *testing.T) {
	original := NewKillRing(10)
	original = original.Kill("first")

	// Test all mutation operations preserve original
	killed := original.Kill("second")
	popped := original.YankPop()
	copied := original.Copy()

	// Original should remain unchanged
	if original.Yank() != "first" {
		t.Errorf("Original ring changed after operations: %q, want %q", original.Yank(), "first")
	}

	// Results should have correct values
	if killed.Yank() != "second" {
		t.Error("Kill() produced incorrect ring")
	}
	if popped.Yank() != "first" {
		t.Error("YankPop() produced incorrect ring")
	}
	if copied.Yank() != "first" {
		t.Error("Copy() produced incorrect ring")
	}
}

func TestKillRing_EmptyStringHandling(t *testing.T) {
	kr := NewKillRing(10)

	// Kill empty string should not add to ring
	kr = kr.Kill("")
	if !kr.IsEmpty() {
		t.Error("Kill(\"\") should not add to ring")
	}

	// Kill non-empty, then try empty
	kr = kr.Kill("valid")
	kr = kr.Kill("")

	// Should still have only "valid"
	if kr.Yank() != "valid" {
		t.Errorf("After Kill(\"\"), Yank() = %q, want %q", kr.Yank(), "valid")
	}
}

func TestKillRing_Yank_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		index int
		want  string
	}{
		{
			name:  "empty ring",
			items: []string{},
			index: 0,
			want:  "",
		},
		{
			name:  "index out of bounds negative",
			items: []string{"a", "b"},
			index: -1,
			want:  "",
		},
		{
			name:  "index out of bounds positive",
			items: []string{"a", "b"},
			index: 10,
			want:  "",
		},
		{
			name:  "valid index 0",
			items: []string{"first", "second"},
			index: 0,
			want:  "first",
		},
		{
			name:  "valid index 1",
			items: []string{"first", "second", "third"},
			index: 1,
			want:  "second",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ring := &KillRing{
				items:   tt.items,
				index:   tt.index,
				maxSize: 10,
			}

			result := ring.Yank()
			if result != tt.want {
				t.Errorf("Yank() = %q, want %q", result, tt.want)
			}
		})
	}
}
