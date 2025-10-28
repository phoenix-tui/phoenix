package value

import (
	"testing"
)

func TestNewPositionCenter(t *testing.T) {
	pos := NewPositionCenter()

	if !pos.IsCenter() {
		t.Error("Expected position to be centered")
	}
	if pos.X() != -1 {
		t.Errorf("Expected X to be -1, got %d", pos.X())
	}
	if pos.Y() != -1 {
		t.Errorf("Expected Y to be -1, got %d", pos.Y())
	}
}

func TestNewPositionCustom(t *testing.T) {
	pos := NewPositionCustom(10, 20)

	if pos.IsCenter() {
		t.Error("Expected position to not be centered")
	}
	if pos.X() != 10 {
		t.Errorf("Expected X to be 10, got %d", pos.X())
	}
	if pos.Y() != 20 {
		t.Errorf("Expected Y to be 20, got %d", pos.Y())
	}
}

func TestPositionCustomAtOrigin(t *testing.T) {
	pos := NewPositionCustom(0, 0)

	if pos.IsCenter() {
		t.Error("Expected position to not be centered")
	}
	if pos.X() != 0 {
		t.Errorf("Expected X to be 0, got %d", pos.X())
	}
	if pos.Y() != 0 {
		t.Errorf("Expected Y to be 0, got %d", pos.Y())
	}
}

func TestPositionNegativeCustomCoordinates(t *testing.T) {
	// Negative coordinates are allowed (e.g., for positioning off-screen effects)
	pos := NewPositionCustom(-5, -10)

	if pos.IsCenter() {
		t.Error("Expected position to not be centered")
	}
	if pos.X() != -5 {
		t.Errorf("Expected X to be -5, got %d", pos.X())
	}
	if pos.Y() != -10 {
		t.Errorf("Expected Y to be -10, got %d", pos.Y())
	}
}
