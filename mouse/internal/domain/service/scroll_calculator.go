package service

import (
	"github.com/phoenix-tui/phoenix/mouse/internal/domain/model"
	value2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

// ScrollCalculator is a domain service that calculates scroll deltas.
type ScrollCalculator struct {
	linesPerScroll int
}

// NewScrollCalculator creates a new ScrollCalculator.
// linesPerScroll: number of lines to scroll per wheel event (typically 3).
func NewScrollCalculator(linesPerScroll int) *ScrollCalculator {
	if linesPerScroll <= 0 {
		linesPerScroll = 3 // Default
	}
	return &ScrollCalculator{
		linesPerScroll: linesPerScroll,
	}
}

// CalculateDelta calculates the scroll delta for a scroll event.
// Returns the number of lines to scroll (positive = down, negative = up).
func (s *ScrollCalculator) CalculateDelta(scrollEvent model.MouseEvent) int {
	if scrollEvent.Type() != value2.EventScroll {
		return 0
	}

	switch scrollEvent.Button() {
	case value2.ButtonWheelUp:
		return -s.linesPerScroll
	case value2.ButtonWheelDown:
		return s.linesPerScroll
	default:
		return 0
	}
}

// IsScrollUp returns true if the event is a scroll up event.
func (s *ScrollCalculator) IsScrollUp(scrollEvent model.MouseEvent) bool {
	return scrollEvent.Type() == value2.EventScroll && scrollEvent.Button() == value2.ButtonWheelUp
}

// IsScrollDown returns true if the event is a scroll down event.
func (s *ScrollCalculator) IsScrollDown(scrollEvent model.MouseEvent) bool {
	return scrollEvent.Type() == value2.EventScroll && scrollEvent.Button() == value2.ButtonWheelDown
}

// LinesPerScroll returns the configured lines per scroll.
func (s *ScrollCalculator) LinesPerScroll() int {
	return s.linesPerScroll
}
