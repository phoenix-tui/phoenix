package value

import (
	"testing"
)

func TestNewSize(t *testing.T) {
	size := NewSize(80, 24)

	if size.Width() != 80 {
		t.Errorf("Expected width 80, got %d", size.Width())
	}
	if size.Height() != 24 {
		t.Errorf("Expected height 24, got %d", size.Height())
	}
}

func TestSizeWithWidth(t *testing.T) {
	original := NewSize(80, 24)
	modified := original.WithWidth(100)

	// Original unchanged (immutability)
	if original.Width() != 80 {
		t.Errorf("Original width should be 80, got %d", original.Width())
	}

	// New size has updated width
	if modified.Width() != 100 {
		t.Errorf("Expected width 100, got %d", modified.Width())
	}
	if modified.Height() != 24 {
		t.Errorf("Expected height 24, got %d", modified.Height())
	}
}

func TestSizeWithHeight(t *testing.T) {
	original := NewSize(80, 24)
	modified := original.WithHeight(40)

	// Original unchanged (immutability)
	if original.Height() != 24 {
		t.Errorf("Original height should be 24, got %d", original.Height())
	}

	// New size has updated height
	if modified.Width() != 80 {
		t.Errorf("Expected width 80, got %d", modified.Width())
	}
	if modified.Height() != 40 {
		t.Errorf("Expected height 40, got %d", modified.Height())
	}
}

func TestSizeChaining(t *testing.T) {
	size := NewSize(50, 20).WithWidth(60).WithHeight(25)

	if size.Width() != 60 {
		t.Errorf("Expected width 60, got %d", size.Width())
	}
	if size.Height() != 25 {
		t.Errorf("Expected height 25, got %d", size.Height())
	}
}

func TestSizeZeroValues(t *testing.T) {
	size := NewSize(0, 0)

	if size.Width() != 0 {
		t.Errorf("Expected width 0, got %d", size.Width())
	}
	if size.Height() != 0 {
		t.Errorf("Expected height 0, got %d", size.Height())
	}
}

func TestSizeNegativeValues(t *testing.T) {
	// Negative values are allowed (rendering logic will clamp)
	size := NewSize(-10, -5)

	if size.Width() != -10 {
		t.Errorf("Expected width -10, got %d", size.Width())
	}
	if size.Height() != -5 {
		t.Errorf("Expected height -5, got %d", size.Height())
	}
}
