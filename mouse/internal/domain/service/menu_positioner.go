package service

import (
	"github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

// MenuPositioner is a domain service that calculates optimal menu positioning.
// It ensures menus stay fully visible within screen bounds by adjusting position
// when the menu would overflow screen edges.
//
// Positioning priority:
//  1. At cursor position (preferred)
//  2. Shift left if right edge overflow
//  3. Shift up if bottom edge overflow
//  4. Shift both if corner overflow
//  5. Pin to (0,0) if menu larger than screen
type MenuPositioner struct{}

// NewMenuPositioner creates a new MenuPositioner.
func NewMenuPositioner() *MenuPositioner {
	return &MenuPositioner{}
}

// CalculatePosition determines optimal menu position to keep it fully visible.
//
// Parameters:
//   - cursorPos: mouse cursor position where menu should ideally appear
//   - menuWidth: width of the menu in terminal cells
//   - menuHeight: height of the menu in terminal cells
//   - screenWidth: terminal width in cells
//   - screenHeight: terminal height in cells
//
// Returns:
//   - adjusted position that keeps menu fully visible on screen
//
// Behavior:
//   - If menu fits at cursor position: returns cursor position
//   - If menu overflows right edge: shifts left to fit
//   - If menu overflows bottom edge: shifts up to fit
//   - If menu overflows corner: shifts both directions
//   - If menu larger than screen: pins to top-left (0,0)
func (p *MenuPositioner) CalculatePosition(
	cursorPos value.Position,
	menuWidth, menuHeight int,
	screenWidth, screenHeight int,
) value.Position {
	// Ensure non-negative dimensions
	if menuWidth < 0 {
		menuWidth = 0
	}
	if menuHeight < 0 {
		menuHeight = 0
	}
	if screenWidth < 0 {
		screenWidth = 0
	}
	if screenHeight < 0 {
		screenHeight = 0
	}

	// Special case: menu larger than screen
	// Pin to top-left corner (0,0)
	if menuWidth > screenWidth || menuHeight > screenHeight {
		return value.NewPosition(0, 0)
	}

	// Start with cursor position
	x := cursorPos.X()
	y := cursorPos.Y()

	// Check right edge overflow
	if x+menuWidth > screenWidth {
		// Shift left to fit
		x = screenWidth - menuWidth
		// Ensure non-negative
		if x < 0 {
			x = 0
		}
	}

	// Check bottom edge overflow
	if y+menuHeight > screenHeight {
		// Shift up to fit
		y = screenHeight - menuHeight
		// Ensure non-negative
		if y < 0 {
			y = 0
		}
	}

	return value.NewPosition(x, y)
}

// WouldOverflow checks if a menu at the given position would overflow screen bounds.
// Returns (overflowsRight, overflowsBottom).
func (p *MenuPositioner) WouldOverflow(
	position value.Position,
	menuWidth, menuHeight int,
	screenWidth, screenHeight int,
) (overflowsRight, overflowsBottom bool) {
	x := position.X()
	y := position.Y()

	overflowsRight = x+menuWidth > screenWidth
	overflowsBottom = y+menuHeight > screenHeight

	return overflowsRight, overflowsBottom
}

// FitsOnScreen checks if a menu with given dimensions fits on screen.
func (p *MenuPositioner) FitsOnScreen(
	menuWidth, menuHeight int,
	screenWidth, screenHeight int,
) bool {
	return menuWidth <= screenWidth && menuHeight <= screenHeight
}
