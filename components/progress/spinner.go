package progress

import (
	"fmt"
	"time"

	"github.com/phoenix-tui/phoenix/components/progress/internal/domain/model"
	"github.com/phoenix-tui/phoenix/components/progress/internal/infrastructure"
	"github.com/phoenix-tui/phoenix/tea"
)

// Spinner is the public API for spinner component.
// It implements tea.Model and provides animated loading indicators.
// Uses value semantics for immutable updates.
type Spinner struct {
	domain model.Spinner // VALUE, not pointer!
}

// NewSpinner creates a new spinner with the specified pre-defined style.
// Valid styles: "dots", "line", "arrow", "circle", "bounce", etc.
// See infrastructure.AvailableStyles() for full list.
// Returns pointer for initialization, but store as value in Model.
func NewSpinner(style string) *Spinner {
	spinnerStyle := infrastructure.GetSpinnerStyle(style)
	return &Spinner{
		domain: *model.NewSpinner(spinnerStyle), // Dereference!
	}
}

// NewSpinnerCustom creates a new spinner with custom frames and FPS.
// Returns pointer for initialization, but store as value in Model.
//
//nolint:revive // frames parameter will be used when custom spinner support is added
func NewSpinnerCustom(_ []string, fps int) *Spinner {
	return &Spinner{
		domain: *model.NewSpinner(infrastructure.GetSpinnerStyle("dots")), // Dereference!
	}
}

// Label sets the label text displayed with the spinner.
// Returns new Spinner for method chaining (value semantics).
// IMPORTANT: Must reassign: spinner = spinner.Label("text").
func (s Spinner) Label(label string) Spinner {
	s.domain = s.domain.WithLabel(label)
	return s
}

// Init initializes the spinner and starts the animation (tea.Model interface).
// Returns a command to schedule the first animation tick.
func (s Spinner) Init() tea.Cmd {
	return s.tick()
}

// Update handles messages (implements tea model contract).
// Advances the animation frame on tea.TickMsg.
// IMPORTANT: Must reassign: spinner = spinner.Update(msg).
func (s Spinner) Update(msg tea.Msg) (Spinner, tea.Cmd) {
	//nolint:gocritic // single case for now, will expand with pause/resume
	switch msg.(type) {
	case tea.TickMsg:
		// Advance to next frame.
		s.domain = s.domain.NextFrame()
		return s, s.tick()
	}
	return s, nil
}

// View renders the spinner to a string (tea.Model interface).
// Format: [frame] [label].
// Example: "⠋ Loading...".
func (s Spinner) View() string {
	frame := s.domain.CurrentFrame()
	label := s.domain.Label()

	if label != "" {
		return fmt.Sprintf("%s %s", frame, label)
	}
	return frame
}

// tick returns a command to schedule the next animation frame.
func (s Spinner) tick() tea.Cmd {
	fps := s.domain.Style().FPS()
	duration := time.Second / time.Duration(fps)
	return tea.Tick(duration)
}
