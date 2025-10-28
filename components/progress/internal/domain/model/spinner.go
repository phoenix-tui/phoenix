package model

import "github.com/phoenix-tui/phoenix/components/progress/internal/domain/value"

// Spinner is the domain model for animated spinner.
// It provides rich domain logic with immutability and encapsulated behavior.
type Spinner struct {
	style      *value.SpinnerStyle // Animation frames and timing
	frameIndex int                 // Current frame index
	label      string              // Optional label
}

// NewSpinner creates a new Spinner with the given style.
func NewSpinner(style *value.SpinnerStyle) *Spinner {
	if style == nil {
		// Default to single dot.
		style = value.NewSpinnerStyle([]string{"•"}, 10)
	}

	return &Spinner{
		style:      style,
		frameIndex: 0,
		label:      "",
	}
}

// WithLabel returns a new Spinner with the specified label.
func (s Spinner) WithLabel(label string) Spinner {
	s.label = label
	return s
}

// NextFrame returns a new Spinner advanced to the next animation frame.
// Wraps around to frame 0 after the last frame.
func (s Spinner) NextFrame() Spinner {
	s.frameIndex = (s.frameIndex + 1) % s.style.FrameCount()
	return s
}

// Reset returns a new Spinner with frame index reset to 0.
func (s Spinner) Reset() Spinner {
	s.frameIndex = 0
	return s
}

// CurrentFrame returns the current animation frame string.
func (s Spinner) CurrentFrame() string {
	return s.style.GetFrame(s.frameIndex)
}

// Label returns the label.
func (s Spinner) Label() string {
	return s.label
}

// FrameIndex returns the current frame index.
func (s Spinner) FrameIndex() int {
	return s.frameIndex
}

// Style returns the spinner style.
func (s Spinner) Style() *value.SpinnerStyle {
	return s.style
}
