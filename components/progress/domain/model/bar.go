package model

import "github.com/phoenix-tui/phoenix/components/progress/domain/value"

// Bar is the domain model for progress bar.
// It provides rich domain logic with immutability and encapsulated behavior.
type Bar struct {
	percentage  *value.Percentage // Current progress (0-100)
	width       int               // Bar width in characters
	fillChar    rune              // Character for filled portion (e.g., '█')
	emptyChar   rune              // Character for empty portion (e.g., '░')
	showPercent bool              // Show percentage text?
	label       string            // Optional label
}

// NewBar creates a new Bar with default settings.
// width specifies the bar width in characters (minimum 1).
func NewBar(width int) *Bar {
	if width < 1 {
		width = 1
	}

	return &Bar{
		percentage:  value.NewPercentage(0),
		width:       width,
		fillChar:    '█',
		emptyChar:   '░',
		showPercent: false,
		label:       "",
	}
}

// NewBarWithPercentage creates a new Bar with initial percentage.
func NewBarWithPercentage(width int, percentage int) *Bar {
	bar := NewBar(width)
	bar.percentage = value.NewPercentage(percentage)
	return bar
}

// WithPercentage returns a new Bar with the specified percentage.
func (b Bar) WithPercentage(pct int) Bar {
	b.percentage = value.NewPercentage(pct)
	return b
}

// WithFillChar returns a new Bar with the specified fill character.
func (b Bar) WithFillChar(char rune) Bar {
	b.fillChar = char
	return b
}

// WithEmptyChar returns a new Bar with the specified empty character.
func (b Bar) WithEmptyChar(char rune) Bar {
	b.emptyChar = char
	return b
}

// WithShowPercent returns a new Bar with percentage display toggled.
func (b Bar) WithShowPercent(show bool) Bar {
	b.showPercent = show
	return b
}

// WithLabel returns a new Bar with the specified label.
func (b Bar) WithLabel(label string) Bar {
	b.label = label
	return b
}

// Increment returns a new Bar with percentage incremented by delta.
// Result is clamped to [0, 100].
func (b Bar) Increment(delta int) Bar {
	b.percentage = b.percentage.Add(delta)
	return b
}

// Decrement returns a new Bar with percentage decremented by delta.
// Result is clamped to [0, 100].
func (b Bar) Decrement(delta int) Bar {
	b.percentage = b.percentage.Subtract(delta)
	return b
}

// SetComplete returns a new Bar with percentage set to 100%.
func (b Bar) SetComplete() Bar {
	return b.WithPercentage(100)
}

// Reset returns a new Bar with percentage set to 0%.
func (b Bar) Reset() Bar {
	return b.WithPercentage(0)
}

// Percentage returns the current percentage value (0-100).
func (b Bar) Percentage() int {
	return b.percentage.Value()
}

// Width returns the bar width in characters.
func (b Bar) Width() int {
	return b.width
}

// FillChar returns the fill character.
func (b Bar) FillChar() rune {
	return b.fillChar
}

// EmptyChar returns the empty character.
func (b Bar) EmptyChar() rune {
	return b.emptyChar
}

// ShowPercent returns whether percentage display is enabled.
func (b Bar) ShowPercent() bool {
	return b.showPercent
}

// Label returns the label.
func (b Bar) Label() string {
	return b.label
}

// IsComplete returns true if percentage is 100%.
func (b Bar) IsComplete() bool {
	return b.percentage.IsComplete()
}
