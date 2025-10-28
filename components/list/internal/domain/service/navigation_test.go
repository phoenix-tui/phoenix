package service

import (
	"testing"
)

func TestNavigationService_MoveUp(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name         string
		currentIndex int
		itemCount    int
		want         int
	}{
		{
			name:         "move from middle",
			currentIndex: 5,
			itemCount:    10,
			want:         4,
		},
		{
			name:         "move from start wraps to end",
			currentIndex: 0,
			itemCount:    10,
			want:         9,
		},
		{
			name:         "empty list",
			currentIndex: 0,
			itemCount:    0,
			want:         0,
		},
		{
			name:         "single item wraps",
			currentIndex: 0,
			itemCount:    1,
			want:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.MoveUp(tt.currentIndex, tt.itemCount); got != tt.want {
				t.Errorf("NavigationService.MoveUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigationService_MoveDown(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name         string
		currentIndex int
		itemCount    int
		want         int
	}{
		{
			name:         "move from middle",
			currentIndex: 5,
			itemCount:    10,
			want:         6,
		},
		{
			name:         "move from end wraps to start",
			currentIndex: 9,
			itemCount:    10,
			want:         0,
		},
		{
			name:         "empty list",
			currentIndex: 0,
			itemCount:    0,
			want:         0,
		},
		{
			name:         "single item wraps",
			currentIndex: 0,
			itemCount:    1,
			want:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.MoveDown(tt.currentIndex, tt.itemCount); got != tt.want {
				t.Errorf("NavigationService.MoveDown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigationService_MovePageUp(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name         string
		currentIndex int
		pageSize     int
		itemCount    int
		want         int
	}{
		{
			name:         "normal page up",
			currentIndex: 15,
			pageSize:     10,
			itemCount:    50,
			want:         5,
		},
		{
			name:         "page up from near start",
			currentIndex: 5,
			pageSize:     10,
			itemCount:    50,
			want:         0,
		},
		{
			name:         "page up at start",
			currentIndex: 0,
			pageSize:     10,
			itemCount:    50,
			want:         0,
		},
		{
			name:         "empty list",
			currentIndex: 0,
			pageSize:     10,
			itemCount:    0,
			want:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.MovePageUp(tt.currentIndex, tt.pageSize, tt.itemCount); got != tt.want {
				t.Errorf("NavigationService.MovePageUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigationService_MovePageDown(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name         string
		currentIndex int
		pageSize     int
		itemCount    int
		want         int
	}{
		{
			name:         "normal page down",
			currentIndex: 5,
			pageSize:     10,
			itemCount:    50,
			want:         15,
		},
		{
			name:         "page down near end",
			currentIndex: 45,
			pageSize:     10,
			itemCount:    50,
			want:         49,
		},
		{
			name:         "page down at end",
			currentIndex: 49,
			pageSize:     10,
			itemCount:    50,
			want:         49,
		},
		{
			name:         "empty list",
			currentIndex: 0,
			pageSize:     10,
			itemCount:    0,
			want:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.MovePageDown(tt.currentIndex, tt.pageSize, tt.itemCount); got != tt.want {
				t.Errorf("NavigationService.MovePageDown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigationService_MoveToStart(t *testing.T) {
	svc := NewNavigationService()

	if got := svc.MoveToStart(); got != 0 {
		t.Errorf("NavigationService.MoveToStart() = %v, want 0", got)
	}
}

func TestNavigationService_MoveToEnd(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name      string
		itemCount int
		want      int
	}{
		{
			name:      "normal list",
			itemCount: 10,
			want:      9,
		},
		{
			name:      "single item",
			itemCount: 1,
			want:      0,
		},
		{
			name:      "empty list",
			itemCount: 0,
			want:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.MoveToEnd(tt.itemCount); got != tt.want {
				t.Errorf("NavigationService.MoveToEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigationService_CalculateScrollOffset(t *testing.T) {
	svc := NewNavigationService()

	tests := []struct {
		name          string
		focusedIndex  int
		visibleHeight int
		itemCount     int
		want          int
	}{
		{
			name:          "all items fit - no scroll",
			focusedIndex:  5,
			visibleHeight: 20,
			itemCount:     10,
			want:          0,
		},
		{
			name:          "focused at start",
			focusedIndex:  2,
			visibleHeight: 10,
			itemCount:     50,
			want:          0,
		},
		{
			name:          "focused in middle - centered",
			focusedIndex:  25,
			visibleHeight: 10,
			itemCount:     50,
			want:          20,
		},
		{
			name:          "focused near end",
			focusedIndex:  45,
			visibleHeight: 10,
			itemCount:     50,
			want:          40,
		},
		{
			name:          "empty list",
			focusedIndex:  0,
			visibleHeight: 10,
			itemCount:     0,
			want:          0,
		},
		{
			name:          "zero height",
			focusedIndex:  5,
			visibleHeight: 0,
			itemCount:     10,
			want:          0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.CalculateScrollOffset(tt.focusedIndex, tt.visibleHeight, tt.itemCount); got != tt.want {
				t.Errorf("NavigationService.CalculateScrollOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}
