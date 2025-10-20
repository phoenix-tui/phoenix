package value

import (
	"reflect"
	"testing"
)

func TestNewSpinnerStyle(t *testing.T) {
	tests := []struct {
		name          string
		frames        []string
		fps           int
		expectedFPS   int
		expectedCount int
	}{
		{
			name:          "Valid style",
			frames:        []string{"⠋", "⠙", "⠹"},
			fps:           10,
			expectedFPS:   10,
			expectedCount: 3,
		},
		{
			name:          "Empty frames - defaults to single dot",
			frames:        []string{},
			fps:           8,
			expectedFPS:   8,
			expectedCount: 1,
		},
		{
			name:          "Zero FPS - defaults to 10",
			frames:        []string{"|", "/", "-", "\\"},
			fps:           0,
			expectedFPS:   10,
			expectedCount: 4,
		},
		{
			name:          "Negative FPS - defaults to 10",
			frames:        []string{"•"},
			fps:           -5,
			expectedFPS:   10,
			expectedCount: 1,
		},
		{
			name:          "Single frame",
			frames:        []string{"●"},
			fps:           1,
			expectedFPS:   1,
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := NewSpinnerStyle(tt.frames, tt.fps)
			if style.FPS() != tt.expectedFPS {
				t.Errorf("FPS() = %d, expected %d", style.FPS(), tt.expectedFPS)
			}
			if style.FrameCount() != tt.expectedCount {
				t.Errorf("FrameCount() = %d, expected %d", style.FrameCount(), tt.expectedCount)
			}
		})
	}
}

func TestSpinnerStyleFrames(t *testing.T) {
	original := []string{"⠋", "⠙", "⠹"}
	style := NewSpinnerStyle(original, 10)

	// Get frames.
	frames := style.Frames()

	// Verify content.
	if !reflect.DeepEqual(frames, original) {
		t.Errorf("Frames() = %v, expected %v", frames, original)
	}

	// Verify immutability - modify returned slice.
	frames[0] = "MODIFIED"
	framesAgain := style.Frames()
	if framesAgain[0] == "MODIFIED" {
		t.Errorf("Frames() returned mutable reference")
	}
}

func TestSpinnerStyleGetFrame(t *testing.T) {
	frames := []string{"A", "B", "C", "D"}
	style := NewSpinnerStyle(frames, 10)

	tests := []struct {
		name     string
		index    int
		expected string
	}{
		{"First frame", 0, "A"},
		{"Middle frame", 2, "C"},
		{"Last frame", 3, "D"},
		{"Wrap around - one past", 4, "A"},
		{"Wrap around - two past", 5, "B"},
		{"Large index", 10, "C"}, // 10 % 4 = 2
		{"Negative index - last", -1, "D"},
		{"Negative index - second to last", -2, "C"},
		{"Negative index - wrap", -5, "D"}, // -5 % 4 = -1 → 3
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := style.GetFrame(tt.index)
			if result != tt.expected {
				t.Errorf("GetFrame(%d) = %s, expected %s", tt.index, result, tt.expected)
			}
		})
	}
}

func TestSpinnerStyleGetFrameEmptyFrames(t *testing.T) {
	// Create style with empty frames (will default to single frame)
	style := NewSpinnerStyle([]string{}, 10)

	// GetFrame should return default frame.
	frame := style.GetFrame(0)
	if frame != "•" {
		t.Errorf("GetFrame(0) on default = %s, expected '•'", frame)
	}
}

func TestSpinnerStyleImmutability(t *testing.T) {
	original := []string{"⠋", "⠙", "⠹"}

	// Modify original after creating style.
	style := NewSpinnerStyle(original, 10)
	original[0] = "MODIFIED"

	// Style should not be affected.
	frames := style.Frames()
	if frames[0] == "MODIFIED" {
		t.Errorf("SpinnerStyle was affected by mutation of input slice")
	}
}

func TestSpinnerStyleFPS(t *testing.T) {
	tests := []struct {
		name     string
		fps      int
		expected int
	}{
		{"Low FPS", 1, 1},
		{"Normal FPS", 10, 10},
		{"High FPS", 60, 60},
		{"Very high FPS", 120, 120},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := NewSpinnerStyle([]string{"•"}, tt.fps)
			if style.FPS() != tt.expected {
				t.Errorf("FPS() = %d, expected %d", style.FPS(), tt.expected)
			}
		})
	}
}
