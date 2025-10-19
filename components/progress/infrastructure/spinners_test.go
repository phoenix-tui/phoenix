package infrastructure

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/progress/domain/value"
)

func TestPredefinedSpinners(t *testing.T) {
	tests := []struct {
		name        string
		spinner     interface{}
		minFrames   int
		expectedFPS int
	}{
		{"SpinnerDots", SpinnerDots, 5, 10},
		{"SpinnerLine", SpinnerLine, 4, 10},
		{"SpinnerArrow", SpinnerArrow, 8, 8},
		{"SpinnerCircle", SpinnerCircle, 4, 8},
		{"SpinnerBounce", SpinnerBounce, 5, 10},
		{"SpinnerDotPulse", SpinnerDotPulse, 5, 12},
		{"SpinnerGrowVertical", SpinnerGrowVertical, 10, 8},
		{"SpinnerGrowHorizontal", SpinnerGrowHorizontal, 10, 8},
		{"SpinnerBoxBounce", SpinnerBoxBounce, 4, 10},
		{"SpinnerSimpleDots", SpinnerSimpleDots, 3, 6},
		{"SpinnerClock", SpinnerClock, 12, 4},
		{"SpinnerEarth", SpinnerEarth, 3, 6},
		{"SpinnerMoon", SpinnerMoon, 8, 5},
		{"SpinnerToggle", SpinnerToggle, 2, 8},
		{"SpinnerHamburger", SpinnerHamburger, 3, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Type assertion - all should be *value.SpinnerStyle
			style := tt.spinner
			if style == nil {
				t.Fatalf("%s is nil", tt.name)
			}

			// Access through interface (we know it's *value.SpinnerStyle)
			s := style.(*value.SpinnerStyle)

			if s.FrameCount() < tt.minFrames {
				t.Errorf("%s: FrameCount() = %d, expected >= %d",
					tt.name, s.FrameCount(), tt.minFrames)
			}

			if s.FPS() != tt.expectedFPS {
				t.Errorf("%s: FPS() = %d, expected %d",
					tt.name, s.FPS(), tt.expectedFPS)
			}

			// Verify all frames are non-empty
			frames := s.Frames()
			for i, frame := range frames {
				if frame == "" {
					t.Errorf("%s: frame %d is empty", tt.name, i)
				}
			}
		})
	}
}

func TestGetSpinnerStyle(t *testing.T) {
	tests := []struct {
		name          string
		expectedStyle interface{}
	}{
		{"dots", SpinnerDots},
		{"line", SpinnerLine},
		{"arrow", SpinnerArrow},
		{"circle", SpinnerCircle},
		{"bounce", SpinnerBounce},
		{"dot-pulse", SpinnerDotPulse},
		{"grow-vertical", SpinnerGrowVertical},
		{"grow-horizontal", SpinnerGrowHorizontal},
		{"box-bounce", SpinnerBoxBounce},
		{"simple-dots", SpinnerSimpleDots},
		{"clock", SpinnerClock},
		{"earth", SpinnerEarth},
		{"moon", SpinnerMoon},
		{"toggle", SpinnerToggle},
		{"hamburger", SpinnerHamburger},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := GetSpinnerStyle(tt.name)
			if style == nil {
				t.Fatalf("GetSpinnerStyle(%s) returned nil", tt.name)
			}

			// Verify it's the correct style (compare frame count)
			expected := tt.expectedStyle.(*value.SpinnerStyle)
			if style.FrameCount() != expected.FrameCount() {
				t.Errorf("GetSpinnerStyle(%s): FrameCount = %d, expected %d",
					tt.name, style.FrameCount(), expected.FrameCount())
			}
		})
	}
}

func TestGetSpinnerStyleUnknown(t *testing.T) {
	style := GetSpinnerStyle("unknown-style-name")
	if style == nil {
		t.Fatal("GetSpinnerStyle() should return default, not nil")
	}

	// Should return SpinnerDots as default
	if style.FrameCount() != SpinnerDots.FrameCount() {
		t.Errorf("Unknown style should return SpinnerDots")
	}
}

func TestAvailableStyles(t *testing.T) {
	styles := AvailableStyles()

	// Check minimum expected styles
	if len(styles) < 15 {
		t.Errorf("AvailableStyles() = %d styles, expected at least 15", len(styles))
	}

	// Verify some known styles are present
	knownStyles := []string{"dots", "line", "arrow", "circle", "bounce"}
	for _, known := range knownStyles {
		found := false
		for _, style := range styles {
			if style == known {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("AvailableStyles() missing %s", known)
		}
	}
}

func TestSpinnerStylesAreUnique(t *testing.T) {
	// Verify each style has unique frames (not shared references)
	styles := []struct {
		name  string
		style interface{}
	}{
		{"dots", SpinnerDots},
		{"line", SpinnerLine},
		{"arrow", SpinnerArrow},
	}

	for i := 0; i < len(styles); i++ {
		for j := i + 1; j < len(styles); j++ {
			s1 := styles[i].style.(*value.SpinnerStyle)
			s2 := styles[j].style.(*value.SpinnerStyle)

			// Different styles should have different frame counts or content
			frames1 := s1.Frames()
			frames2 := s2.Frames()

			if len(frames1) == len(frames2) {
				// Same count - check if frames are different
				allSame := true
				for k := 0; k < len(frames1); k++ {
					if frames1[k] != frames2[k] {
						allSame = false
						break
					}
				}
				if allSame {
					t.Errorf("%s and %s have identical frames",
						styles[i].name, styles[j].name)
				}
			}
		}
	}
}

func TestSpinnerDotsFrames(t *testing.T) {
	frames := SpinnerDots.Frames()
	if len(frames) != 10 {
		t.Errorf("SpinnerDots: expected 10 frames, got %d", len(frames))
	}

	// Verify specific frames
	expected := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	for i, exp := range expected {
		if frames[i] != exp {
			t.Errorf("SpinnerDots frame %d = %s, expected %s", i, frames[i], exp)
		}
	}
}

func TestSpinnerLineFrames(t *testing.T) {
	frames := SpinnerLine.Frames()
	if len(frames) != 4 {
		t.Errorf("SpinnerLine: expected 4 frames, got %d", len(frames))
	}

	expected := []string{"|", "/", "-", "\\"}
	for i, exp := range expected {
		if frames[i] != exp {
			t.Errorf("SpinnerLine frame %d = %s, expected %s", i, frames[i], exp)
		}
	}
}

func TestSpinnerCircleFrames(t *testing.T) {
	frames := SpinnerCircle.Frames()
	if len(frames) != 4 {
		t.Errorf("SpinnerCircle: expected 4 frames, got %d", len(frames))
	}

	expected := []string{"◐", "◓", "◑", "◒"}
	for i, exp := range expected {
		if frames[i] != exp {
			t.Errorf("SpinnerCircle frame %d = %s, expected %s", i, frames[i], exp)
		}
	}
}

func TestStyleRegistryCompleteness(t *testing.T) {
	// Verify all exported spinners are in registry
	allStyles := []struct {
		name  string
		style interface{}
	}{
		{"dots", SpinnerDots},
		{"line", SpinnerLine},
		{"arrow", SpinnerArrow},
		{"circle", SpinnerCircle},
		{"bounce", SpinnerBounce},
		{"dot-pulse", SpinnerDotPulse},
		{"grow-vertical", SpinnerGrowVertical},
		{"grow-horizontal", SpinnerGrowHorizontal},
		{"box-bounce", SpinnerBoxBounce},
		{"simple-dots", SpinnerSimpleDots},
		{"clock", SpinnerClock},
		{"earth", SpinnerEarth},
		{"moon", SpinnerMoon},
		{"toggle", SpinnerToggle},
		{"hamburger", SpinnerHamburger},
	}

	for _, tt := range allStyles {
		style := GetSpinnerStyle(tt.name)
		if style == nil {
			t.Errorf("Style %s not in registry", tt.name)
		}
	}
}
