// Package progress provides progress bar and spinner components for TUI applications.
package progress

import (
	"github.com/phoenix-tui/phoenix/components/progress/domain/model"
	"github.com/phoenix-tui/phoenix/components/progress/domain/service"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// Bar is the public API for progress bar component.
// It implements tea.Model and provides a fluent interface for configuration.
// Uses value semantics for immutable updates.
type Bar struct {
	domain  model.Bar // VALUE, not pointer!
	service *service.RenderService
}

// NewBar creates a new progress bar with the specified width.
// The bar starts at 0% progress with default styling.
// Returns pointer for initialization, but store as value in Model.
func NewBar(width int) *Bar {
	return &Bar{
		domain:  *model.NewBar(width), // Dereference!
		service: service.NewRenderService(),
	}
}

// NewBarWithProgress creates a new progress bar with initial percentage.
// Returns pointer for initialization, but store as value in Model.
func NewBarWithProgress(width, percentage int) *Bar {
	return &Bar{
		domain:  *model.NewBarWithPercentage(width, percentage), // Dereference!
		service: service.NewRenderService(),
	}
}

// FillChar sets the character used for the filled portion of the bar.
// Returns new Bar for method chaining (value semantics).
// IMPORTANT: Must reassign: bar = bar.FillChar('█').
func (b Bar) FillChar(char rune) Bar {
	b.domain = b.domain.WithFillChar(char)
	return b
}

// EmptyChar sets the character used for the empty portion of the bar.
// Returns new Bar for method chaining (value semantics).
// IMPORTANT: Must reassign: bar = bar.EmptyChar('░').
func (b Bar) EmptyChar(char rune) Bar {
	b.domain = b.domain.WithEmptyChar(char)
	return b
}

// ShowPercent toggles whether to display the percentage text.
// Returns new Bar for method chaining (value semantics).
// IMPORTANT: Must reassign: bar = bar.ShowPercent(true).
func (b Bar) ShowPercent(show bool) Bar {
	b.domain = b.domain.WithShowPercent(show)
	return b
}

// Label sets the label text displayed before the bar.
// Returns new Bar for method chaining (value semantics).
// IMPORTANT: Must reassign: bar = bar.Label("Loading").
func (b Bar) Label(label string) Bar {
	b.domain = b.domain.WithLabel(label)
	return b
}

// SetProgress sets the progress percentage (0-100).
// Values are automatically clamped to valid range.
// Returns new Bar for method chaining (value semantics).
// IMPORTANT: Must reassign: bar = bar.SetProgress(50).
func (b Bar) SetProgress(pct int) Bar {
	b.domain = b.domain.WithPercentage(pct)
	return b
}

// Increment increases the progress by the specified delta.
// Result is clamped to [0, 100].
// Returns new Bar for method chaining (value semantics).
// IMPORTANT: Must reassign: bar = bar.Increment(10).
func (b Bar) Increment(delta int) Bar {
	b.domain = b.domain.Increment(delta)
	return b
}

// Decrement decreases the progress by the specified delta.
// Result is clamped to [0, 100].
// Returns new Bar for method chaining (value semantics).
// IMPORTANT: Must reassign: bar = bar.Decrement(10).
func (b Bar) Decrement(delta int) Bar {
	b.domain = b.domain.Decrement(delta)
	return b
}

// Progress returns the current progress percentage (0-100).
func (b Bar) Progress() int {
	return b.domain.Percentage()
}

// IsComplete returns true if the progress is 100%.
func (b Bar) IsComplete() bool {
	return b.domain.IsComplete()
}

// Init initializes the progress bar (tea.Model interface).
// Returns nil as bars don't need initialization commands.
func (b Bar) Init() tea.Cmd {
	return nil
}

// Update handles messages (implements tea model contract).
// Progress bars don't respond to standard messages - use SetProgress() instead.
// IMPORTANT: Must reassign: bar = bar.Update(msg).
func (b Bar) Update(_ tea.Msg) (Bar, tea.Cmd) {
	// Progress bars are controlled programmatically, not by messages.
	// Application code should call SetProgress() to update.
	return b, nil
}

// View renders the progress bar to a string (tea.Model interface).
func (b Bar) View() string {
	return b.service.RenderBar(&b.domain)
}
