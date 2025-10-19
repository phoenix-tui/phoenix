package service

import (
	"reflect"
	"testing"
)

func TestNewScrollService(t *testing.T) {
	s := NewScrollService()
	if s == nil {
		t.Error("NewScrollService() returned nil")
	}
}

func TestScrollService_VisibleLines(t *testing.T) {
	content := []string{"Line 1", "Line 2", "Line 3", "Line 4", "Line 5"}
	s := NewScrollService()

	tests := []struct {
		name    string
		content []string
		offset  int
		height  int
		want    []string
	}{
		{
			name:    "view from start",
			content: content,
			offset:  0,
			height:  3,
			want:    []string{"Line 1", "Line 2", "Line 3"},
		},
		{
			name:    "view from middle",
			content: content,
			offset:  2,
			height:  2,
			want:    []string{"Line 3", "Line 4"},
		},
		{
			name:    "view at end",
			content: content,
			offset:  3,
			height:  5,
			want:    []string{"Line 4", "Line 5"},
		},
		{
			name:    "height larger than content",
			content: content,
			offset:  0,
			height:  10,
			want:    content,
		},
		{
			name:    "empty content",
			content: []string{},
			offset:  0,
			height:  5,
			want:    []string{},
		},
		{
			name:    "zero height",
			content: content,
			offset:  0,
			height:  0,
			want:    []string{},
		},
		{
			name:    "negative height",
			content: content,
			offset:  0,
			height:  -1,
			want:    []string{},
		},
		{
			name:    "negative offset clamped",
			content: content,
			offset:  -5,
			height:  3,
			want:    []string{"Line 1", "Line 2", "Line 3"},
		},
		{
			name:    "offset beyond content",
			content: content,
			offset:  10,
			height:  3,
			want:    []string{"Line 5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.VisibleLines(tt.content, tt.offset, tt.height)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VisibleLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScrollService_MaxScrollOffset(t *testing.T) {
	s := NewScrollService()

	tests := []struct {
		name           string
		totalLines     int
		viewportHeight int
		want           int
	}{
		{"content larger than viewport", 100, 20, 80},
		{"content equal to viewport", 20, 20, 0},
		{"content smaller than viewport", 10, 20, 0},
		{"empty content", 0, 20, 0},
		{"zero viewport height", 100, 0, 100},
		{"both zero", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.MaxScrollOffset(tt.totalLines, tt.viewportHeight)
			if got != tt.want {
				t.Errorf("MaxScrollOffset(%d, %d) = %d, want %d", tt.totalLines, tt.viewportHeight, got, tt.want)
			}
		})
	}
}

func TestScrollService_ScrollUp(t *testing.T) {
	s := NewScrollService()

	tests := []struct {
		name          string
		currentOffset int
		lines         int
		want          int
	}{
		{"scroll up within bounds", 10, 3, 7},
		{"scroll up to zero", 5, 5, 0},
		{"scroll up below zero clamped", 5, 10, 0},
		{"scroll up from zero", 0, 5, 0},
		{"scroll up by zero", 10, 0, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.ScrollUp(tt.currentOffset, tt.lines)
			if got != tt.want {
				t.Errorf("ScrollUp(%d, %d) = %d, want %d", tt.currentOffset, tt.lines, got, tt.want)
			}
		})
	}
}

func TestScrollService_ScrollDown(t *testing.T) {
	s := NewScrollService()

	tests := []struct {
		name          string
		currentOffset int
		lines         int
		maxOffset     int
		want          int
	}{
		{"scroll down within bounds", 5, 3, 20, 8},
		{"scroll down to max", 15, 5, 20, 20},
		{"scroll down beyond max clamped", 15, 10, 20, 20},
		{"scroll down from max", 20, 5, 20, 20},
		{"scroll down by zero", 10, 0, 20, 10},
		{"scroll down with zero max", 0, 5, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.ScrollDown(tt.currentOffset, tt.lines, tt.maxOffset)
			if got != tt.want {
				t.Errorf("ScrollDown(%d, %d, %d) = %d, want %d", tt.currentOffset, tt.lines, tt.maxOffset, got, tt.want)
			}
		})
	}
}

func TestScrollService_FollowModeOffset(t *testing.T) {
	s := NewScrollService()

	tests := []struct {
		name           string
		totalLines     int
		viewportHeight int
		want           int
	}{
		{"content larger than viewport", 100, 20, 80},
		{"content equal to viewport", 20, 20, 0},
		{"content smaller than viewport", 10, 20, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.FollowModeOffset(tt.totalLines, tt.viewportHeight)
			if got != tt.want {
				t.Errorf("FollowModeOffset(%d, %d) = %d, want %d", tt.totalLines, tt.viewportHeight, got, tt.want)
			}
		})
	}
}

func TestScrollService_CanScrollUp(t *testing.T) {
	s := NewScrollService()

	tests := []struct {
		name          string
		currentOffset int
		want          bool
	}{
		{"can scroll up from middle", 10, true},
		{"can scroll up from 1", 1, true},
		{"cannot scroll up from 0", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.CanScrollUp(tt.currentOffset)
			if got != tt.want {
				t.Errorf("CanScrollUp(%d) = %v, want %v", tt.currentOffset, got, tt.want)
			}
		})
	}
}

func TestScrollService_CanScrollDown(t *testing.T) {
	s := NewScrollService()

	tests := []struct {
		name           string
		currentOffset  int
		totalLines     int
		viewportHeight int
		want           bool
	}{
		{"can scroll down from top", 0, 100, 20, true},
		{"can scroll down from middle", 40, 100, 20, true},
		{"cannot scroll down at max", 80, 100, 20, false},
		{"cannot scroll down when content fits", 0, 20, 20, false},
		{"cannot scroll down when content smaller", 0, 10, 20, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.CanScrollDown(tt.currentOffset, tt.totalLines, tt.viewportHeight)
			if got != tt.want {
				t.Errorf("CanScrollDown(%d, %d, %d) = %v, want %v", tt.currentOffset, tt.totalLines, tt.viewportHeight, got, tt.want)
			}
		})
	}
}

func TestScrollService_IsAtTop(t *testing.T) {
	s := NewScrollService()

	tests := []struct {
		name          string
		currentOffset int
		want          bool
	}{
		{"is at top when offset is 0", 0, true},
		{"not at top when offset is positive", 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.IsAtTop(tt.currentOffset)
			if got != tt.want {
				t.Errorf("IsAtTop(%d) = %v, want %v", tt.currentOffset, got, tt.want)
			}
		})
	}
}

func TestScrollService_IsAtBottom(t *testing.T) {
	s := NewScrollService()

	tests := []struct {
		name           string
		currentOffset  int
		totalLines     int
		viewportHeight int
		want           bool
	}{
		{"is at bottom at max offset", 80, 100, 20, true},
		{"not at bottom before max", 79, 100, 20, false},
		{"is at bottom when content fits", 0, 20, 20, true},
		{"is at bottom when content smaller", 0, 10, 20, true},
		{"not at bottom at top", 0, 100, 20, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.IsAtBottom(tt.currentOffset, tt.totalLines, tt.viewportHeight)
			if got != tt.want {
				t.Errorf("IsAtBottom(%d, %d, %d) = %v, want %v", tt.currentOffset, tt.totalLines, tt.viewportHeight, got, tt.want)
			}
		})
	}
}

func TestScrollService_EdgeCases(t *testing.T) {
	s := NewScrollService()

	t.Run("single line content", func(t *testing.T) {
		content := []string{"Only line"}
		visible := s.VisibleLines(content, 0, 5)
		if len(visible) != 1 || visible[0] != "Only line" {
			t.Errorf("Single line content failed: got %v", visible)
		}

		maxOffset := s.MaxScrollOffset(1, 5)
		if maxOffset != 0 {
			t.Errorf("MaxScrollOffset for single line = %d, want 0", maxOffset)
		}
	})

	t.Run("exact fit content", func(t *testing.T) {
		content := []string{"Line 1", "Line 2", "Line 3"}
		visible := s.VisibleLines(content, 0, 3)
		if !reflect.DeepEqual(visible, content) {
			t.Errorf("Exact fit content: got %v, want %v", visible, content)
		}

		if s.CanScrollDown(0, 3, 3) {
			t.Error("Should not be able to scroll down when content fits exactly")
		}
	})
}
