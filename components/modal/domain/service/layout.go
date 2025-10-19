package service

import (
	"github.com/phoenix-tui/phoenix/components/modal/domain/value"
)

// LayoutService handles modal positioning and centering calculations.
// This is a domain service because positioning logic is part of the business rules
// for how modals should be displayed.
type LayoutService struct{}

// NewLayoutService creates a new layout service.
func NewLayoutService() *LayoutService {
	return &LayoutService{}
}

// CenterPosition calculates the centered position for a modal.
// Returns (x, y) coordinates for the top-left corner of the modal.
//
// Parameters:
//   - terminalWidth: Width of the terminal in characters
//   - terminalHeight: Height of the terminal in rows
//   - modalWidth: Width of the modal in characters
//   - modalHeight: Height of the modal in rows
//
// Returns the x, y coordinates for the top-left corner to center the modal.
func (s *LayoutService) CenterPosition(terminalWidth, terminalHeight, modalWidth, modalHeight int) (x, y int) {
	x = (terminalWidth - modalWidth) / 2
	y = (terminalHeight - modalHeight) / 2

	// Ensure non-negative coordinates
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	return x, y
}

// CalculatePosition calculates the actual position for a modal based on its position value.
// If the position is centered, calculates the center coordinates.
// If the position is custom, returns the custom coordinates as-is.
//
// Parameters:
//   - position: The position value object (center or custom)
//   - terminalWidth: Width of the terminal in characters
//   - terminalHeight: Height of the terminal in rows
//   - modalWidth: Width of the modal in characters
//   - modalHeight: Height of the modal in rows
//
// Returns the x, y coordinates for the top-left corner of the modal.
func (s *LayoutService) CalculatePosition(position *value.Position, terminalWidth, terminalHeight, modalWidth, modalHeight int) (x, y int) {
	if position.IsCenter() {
		return s.CenterPosition(terminalWidth, terminalHeight, modalWidth, modalHeight)
	}

	// Custom position - return as-is (may be negative or out of bounds)
	return position.X(), position.Y()
}
