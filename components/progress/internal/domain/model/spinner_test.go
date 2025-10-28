package model

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/progress/internal/domain/value"
)

func TestNewSpinner(t *testing.T) {
	style := value.NewSpinnerStyle([]string{"⠋", "⠙", "⠹"}, 10)
	spinner := NewSpinner(style)

	if spinner.FrameIndex() != 0 {
		t.Errorf("FrameIndex() = %d, expected 0", spinner.FrameIndex())
	}
	if spinner.Label() != "" {
		t.Errorf("Label() = %s, expected empty", spinner.Label())
	}
	if spinner.CurrentFrame() != "⠋" {
		t.Errorf("CurrentFrame() = %s, expected '⠋'", spinner.CurrentFrame())
	}
}

func TestNewSpinnerNilStyle(t *testing.T) {
	spinner := NewSpinner(nil)

	if spinner.FrameIndex() != 0 {
		t.Errorf("FrameIndex() = %d, expected 0", spinner.FrameIndex())
	}
	// Should have default style.
	if spinner.CurrentFrame() == "" {
		t.Errorf("CurrentFrame() should not be empty with nil style")
	}
}

func TestSpinnerWithLabel(t *testing.T) {
	style := value.NewSpinnerStyle([]string{"•"}, 10)
	spinner := NewSpinner(style)

	newSpinner := spinner.WithLabel("Loading...")
	if newSpinner.Label() != "Loading..." {
		t.Errorf("WithLabel() = %s, expected 'Loading...'", newSpinner.Label())
	}

	// Verify immutability.
	if spinner.Label() != "" {
		t.Errorf("WithLabel() mutated original")
	}
}

func TestSpinnerNextFrame(t *testing.T) {
	frames := []string{"A", "B", "C"}
	style := value.NewSpinnerStyle(frames, 10)
	spinner := *NewSpinner(style) // Dereference to get value

	// Frame 0 → A.
	if spinner.CurrentFrame() != "A" {
		t.Errorf("Initial CurrentFrame() = %s, expected 'A'", spinner.CurrentFrame())
	}

	// Frame 1 → B.
	spinner = spinner.NextFrame()
	if spinner.CurrentFrame() != "B" {
		t.Errorf("NextFrame() = %s, expected 'B'", spinner.CurrentFrame())
	}
	if spinner.FrameIndex() != 1 {
		t.Errorf("FrameIndex() = %d, expected 1", spinner.FrameIndex())
	}

	// Frame 2 → C.
	spinner = spinner.NextFrame()
	if spinner.CurrentFrame() != "C" {
		t.Errorf("NextFrame() = %s, expected 'C'", spinner.CurrentFrame())
	}

	// Frame 3 → A (wrap around)
	spinner = spinner.NextFrame()
	if spinner.CurrentFrame() != "A" {
		t.Errorf("NextFrame() wrap = %s, expected 'A'", spinner.CurrentFrame())
	}
	if spinner.FrameIndex() != 0 {
		t.Errorf("FrameIndex() after wrap = %d, expected 0", spinner.FrameIndex())
	}
}

func TestSpinnerReset(t *testing.T) {
	frames := []string{"A", "B", "C"}
	style := value.NewSpinnerStyle(frames, 10)
	spinner := *NewSpinner(style) // Dereference to get value

	// Advance to frame 2.
	spinner = spinner.NextFrame().NextFrame()
	if spinner.FrameIndex() != 2 {
		t.Errorf("Setup: FrameIndex() = %d, expected 2", spinner.FrameIndex())
	}

	// Reset.
	newSpinner := spinner.Reset()
	if newSpinner.FrameIndex() != 0 {
		t.Errorf("Reset() FrameIndex() = %d, expected 0", newSpinner.FrameIndex())
	}
	if newSpinner.CurrentFrame() != "A" {
		t.Errorf("Reset() CurrentFrame() = %s, expected 'A'", newSpinner.CurrentFrame())
	}

	// Verify immutability.
	if spinner.FrameIndex() != 2 {
		t.Errorf("Reset() mutated original")
	}
}

func TestSpinnerImmutability(t *testing.T) {
	style := value.NewSpinnerStyle([]string{"A", "B", "C"}, 10)
	spinner := NewSpinner(style).WithLabel("Test")

	// Apply all mutations.
	_ = spinner.WithLabel("Changed")
	_ = spinner.NextFrame()
	_ = spinner.Reset()

	// Original should be unchanged.
	if spinner.Label() != "Test" {
		t.Errorf("Label mutated: %s != 'Test'", spinner.Label())
	}
	if spinner.FrameIndex() != 0 {
		t.Errorf("FrameIndex mutated: %d != 0", spinner.FrameIndex())
	}
	if spinner.CurrentFrame() != "A" {
		t.Errorf("CurrentFrame mutated: %s != 'A'", spinner.CurrentFrame())
	}
}

func TestSpinnerFluentInterface(t *testing.T) {
	style := value.NewSpinnerStyle([]string{"⠋", "⠙", "⠹"}, 10)
	spinner := NewSpinner(style).
		WithLabel("Loading...").
		NextFrame()

	if spinner.Label() != "Loading..." {
		t.Errorf("Fluent Label() = %s, expected 'Loading...'", spinner.Label())
	}
	if spinner.FrameIndex() != 1 {
		t.Errorf("Fluent FrameIndex() = %d, expected 1", spinner.FrameIndex())
	}
	if spinner.CurrentFrame() != "⠙" {
		t.Errorf("Fluent CurrentFrame() = %s, expected '⠙'", spinner.CurrentFrame())
	}
}

func TestSpinnerSingleFrame(t *testing.T) {
	style := value.NewSpinnerStyle([]string{"●"}, 10)
	spinner := *NewSpinner(style) // Dereference to get value

	// Single frame should always return the same frame.
	for i := 0; i < 5; i++ {
		if spinner.CurrentFrame() != "●" {
			t.Errorf("CurrentFrame() = %s, expected '●'", spinner.CurrentFrame())
		}
		spinner = spinner.NextFrame()
	}
}

func TestSpinnerStyle(t *testing.T) {
	style := value.NewSpinnerStyle([]string{"⠋", "⠙", "⠹"}, 10)
	spinner := NewSpinner(style)

	retrievedStyle := spinner.Style()
	if retrievedStyle.FrameCount() != 3 {
		t.Errorf("Style().FrameCount() = %d, expected 3", retrievedStyle.FrameCount())
	}
	if retrievedStyle.FPS() != 10 {
		t.Errorf("Style().FPS() = %d, expected 10", retrievedStyle.FPS())
	}
}
