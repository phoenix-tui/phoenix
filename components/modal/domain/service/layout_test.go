package service

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/modal/domain/value"
)

func TestCenterPosition(t *testing.T) {
	service := NewLayoutService()

	tests := []struct {
		name           string
		terminalWidth  int
		terminalHeight int
		modalWidth     int
		modalHeight    int
		expectedX      int
		expectedY      int
	}{
		{
			name:           "center in 80x24 terminal, 40x10 modal",
			terminalWidth:  80,
			terminalHeight: 24,
			modalWidth:     40,
			modalHeight:    10,
			expectedX:      20,
			expectedY:      7,
		},
		{
			name:           "center in 100x30 terminal, 50x15 modal",
			terminalWidth:  100,
			terminalHeight: 30,
			modalWidth:     50,
			modalHeight:    15,
			expectedX:      25,
			expectedY:      7,
		},
		{
			name:           "center in 60x20 terminal, 20x8 modal",
			terminalWidth:  60,
			terminalHeight: 20,
			modalWidth:     20,
			modalHeight:    8,
			expectedX:      20,
			expectedY:      6,
		},
		{
			name:           "modal larger than terminal (width)",
			terminalWidth:  40,
			terminalHeight: 24,
			modalWidth:     80,
			modalHeight:    10,
			expectedX:      0, // Clamped to 0
			expectedY:      7,
		},
		{
			name:           "modal larger than terminal (height)",
			terminalWidth:  80,
			terminalHeight: 10,
			modalWidth:     40,
			modalHeight:    24,
			expectedX:      20,
			expectedY:      0, // Clamped to 0
		},
		{
			name:           "modal larger than terminal (both)",
			terminalWidth:  40,
			terminalHeight: 10,
			modalWidth:     80,
			modalHeight:    24,
			expectedX:      0, // Clamped to 0
			expectedY:      0, // Clamped to 0
		},
		{
			name:           "odd-sized terminal",
			terminalWidth:  81,
			terminalHeight: 25,
			modalWidth:     40,
			modalHeight:    10,
			expectedX:      20, // (81 - 40) / 2 = 20
			expectedY:      7,  // (25 - 10) / 2 = 7
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := service.CenterPosition(tt.terminalWidth, tt.terminalHeight, tt.modalWidth, tt.modalHeight)

			if x != tt.expectedX {
				t.Errorf("Expected x=%d, got x=%d", tt.expectedX, x)
			}
			if y != tt.expectedY {
				t.Errorf("Expected y=%d, got y=%d", tt.expectedY, y)
			}
		})
	}
}

func TestCalculatePositionCenter(t *testing.T) {
	service := NewLayoutService()
	position := value.NewPositionCenter()

	x, y := service.CalculatePosition(position, 80, 24, 40, 10)

	expectedX, expectedY := 20, 7
	if x != expectedX || y != expectedY {
		t.Errorf("Expected position (%d, %d), got (%d, %d)", expectedX, expectedY, x, y)
	}
}

func TestCalculatePositionCustom(t *testing.T) {
	service := NewLayoutService()
	position := value.NewPositionCustom(10, 5)

	x, y := service.CalculatePosition(position, 80, 24, 40, 10)

	// Custom position should be returned as-is.
	if x != 10 || y != 5 {
		t.Errorf("Expected position (10, 5), got (%d, %d)", x, y)
	}
}

func TestCalculatePositionCustomNegative(t *testing.T) {
	service := NewLayoutService()
	position := value.NewPositionCustom(-5, -3)

	x, y := service.CalculatePosition(position, 80, 24, 40, 10)

	// Negative custom positions allowed (caller decides how to handle)
	if x != -5 || y != -3 {
		t.Errorf("Expected position (-5, -3), got (%d, %d)", x, y)
	}
}

func TestCalculatePositionCustomOutOfBounds(t *testing.T) {
	service := NewLayoutService()
	position := value.NewPositionCustom(100, 50)

	x, y := service.CalculatePosition(position, 80, 24, 40, 10)

	// Out-of-bounds custom positions allowed (caller decides how to handle)
	if x != 100 || y != 50 {
		t.Errorf("Expected position (100, 50), got (%d, %d)", x, y)
	}
}
