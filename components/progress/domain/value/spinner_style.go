package value

// SpinnerStyle defines spinner animation frames and timing.
// It provides value object semantics for spinner configurations.
type SpinnerStyle struct {
	frames []string // Animation frames (e.g., ["⠋", "⠙", "⠹", ...])
	fps    int      // Frames per second (for timing)
}

// NewSpinnerStyle creates a new SpinnerStyle value object.
// frames must contain at least one frame.
// fps must be positive (defaults to 10 if invalid).
func NewSpinnerStyle(frames []string, fps int) *SpinnerStyle {
	if len(frames) == 0 {
		frames = []string{"•"} // Default single frame
	}
	if fps <= 0 {
		fps = 10 // Default FPS
	}

	// Copy frames to ensure immutability
	framesCopy := make([]string, len(frames))
	copy(framesCopy, frames)

	return &SpinnerStyle{
		frames: framesCopy,
		fps:    fps,
	}
}

// Frames returns a copy of the animation frames.
func (s *SpinnerStyle) Frames() []string {
	framesCopy := make([]string, len(s.frames))
	copy(framesCopy, s.frames)
	return framesCopy
}

// FPS returns the frames per second.
func (s *SpinnerStyle) FPS() int {
	return s.fps
}

// FrameCount returns the number of frames in the animation.
func (s *SpinnerStyle) FrameCount() int {
	return len(s.frames)
}

// GetFrame returns the frame at the given index, with wrap-around.
// Negative indices are supported (wrap from end).
func (s *SpinnerStyle) GetFrame(index int) string {
	if len(s.frames) == 0 {
		return ""
	}

	// Normalize negative indices
	if index < 0 {
		index = len(s.frames) + (index % len(s.frames))
	}

	// Wrap around
	index = index % len(s.frames)
	return s.frames[index]
}
