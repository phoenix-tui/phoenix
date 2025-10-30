package service

import (
	"testing"

	"github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

func TestNewMenuPositioner(t *testing.T) {
	p := NewMenuPositioner()
	if p == nil {
		t.Fatal("expected non-nil MenuPositioner")
	}
}

func TestCalculatePosition_FitsAtCursor(t *testing.T) {
	p := NewMenuPositioner()

	// Menu fits perfectly at cursor position
	cursorPos := value.NewPosition(10, 5)
	menuWidth, menuHeight := 20, 10
	screenWidth, screenHeight := 80, 24

	result := p.CalculatePosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	// Should return cursor position unchanged
	if result.X() != 10 {
		t.Errorf("expected x=10, got x=%d", result.X())
	}
	if result.Y() != 5 {
		t.Errorf("expected y=5, got y=%d", result.Y())
	}
}

func TestCalculatePosition_RightEdgeOverflow(t *testing.T) {
	p := NewMenuPositioner()

	// Menu would overflow right edge
	cursorPos := value.NewPosition(70, 5)
	menuWidth, menuHeight := 20, 10
	screenWidth, screenHeight := 80, 24

	result := p.CalculatePosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	// Should shift left to fit
	expectedX := 60 // screenWidth (80) - menuWidth (20)
	if result.X() != expectedX {
		t.Errorf("expected x=%d, got x=%d", expectedX, result.X())
	}
	// Y should remain unchanged
	if result.Y() != 5 {
		t.Errorf("expected y=5, got y=%d", result.Y())
	}
}

func TestCalculatePosition_BottomEdgeOverflow(t *testing.T) {
	p := NewMenuPositioner()

	// Menu would overflow bottom edge
	cursorPos := value.NewPosition(10, 20)
	menuWidth, menuHeight := 20, 10
	screenWidth, screenHeight := 80, 24

	result := p.CalculatePosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	// X should remain unchanged
	if result.X() != 10 {
		t.Errorf("expected x=10, got x=%d", result.X())
	}
	// Should shift up to fit
	expectedY := 14 // screenHeight (24) - menuHeight (10)
	if result.Y() != expectedY {
		t.Errorf("expected y=%d, got y=%d", expectedY, result.Y())
	}
}

func TestCalculatePosition_CornerOverflow(t *testing.T) {
	p := NewMenuPositioner()

	// Menu would overflow both right and bottom edges (corner case)
	cursorPos := value.NewPosition(70, 20)
	menuWidth, menuHeight := 20, 10
	screenWidth, screenHeight := 80, 24

	result := p.CalculatePosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	// Should shift both left and up
	expectedX := 60 // screenWidth (80) - menuWidth (20)
	expectedY := 14 // screenHeight (24) - menuHeight (10)

	if result.X() != expectedX {
		t.Errorf("expected x=%d, got x=%d", expectedX, result.X())
	}
	if result.Y() != expectedY {
		t.Errorf("expected y=%d, got y=%d", expectedY, result.Y())
	}
}

func TestCalculatePosition_MenuLargerThanScreen(t *testing.T) {
	p := NewMenuPositioner()

	// Menu is larger than screen
	cursorPos := value.NewPosition(10, 5)
	menuWidth, menuHeight := 100, 30
	screenWidth, screenHeight := 80, 24

	result := p.CalculatePosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	// Should pin to top-left (0,0)
	if result.X() != 0 {
		t.Errorf("expected x=0, got x=%d", result.X())
	}
	if result.Y() != 0 {
		t.Errorf("expected y=0, got y=%d", result.Y())
	}
}

func TestCalculatePosition_MenuWidthLargerThanScreen(t *testing.T) {
	p := NewMenuPositioner()

	// Menu width exceeds screen width (height is OK)
	cursorPos := value.NewPosition(10, 5)
	menuWidth, menuHeight := 100, 10
	screenWidth, screenHeight := 80, 24

	result := p.CalculatePosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	// Should pin to top-left (0,0)
	if result.X() != 0 || result.Y() != 0 {
		t.Errorf("expected (0,0), got (%d,%d)", result.X(), result.Y())
	}
}

func TestCalculatePosition_MenuHeightLargerThanScreen(t *testing.T) {
	p := NewMenuPositioner()

	// Menu height exceeds screen height (width is OK)
	cursorPos := value.NewPosition(10, 5)
	menuWidth, menuHeight := 20, 30
	screenWidth, screenHeight := 80, 24

	result := p.CalculatePosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	// Should pin to top-left (0,0)
	if result.X() != 0 || result.Y() != 0 {
		t.Errorf("expected (0,0), got (%d,%d)", result.X(), result.Y())
	}
}

func TestCalculatePosition_CursorAtOrigin(t *testing.T) {
	p := NewMenuPositioner()

	// Cursor at (0,0)
	cursorPos := value.NewPosition(0, 0)
	menuWidth, menuHeight := 20, 10
	screenWidth, screenHeight := 80, 24

	result := p.CalculatePosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	// Should stay at (0,0)
	if result.X() != 0 || result.Y() != 0 {
		t.Errorf("expected (0,0), got (%d,%d)", result.X(), result.Y())
	}
}

func TestCalculatePosition_CursorAtBottomRight(t *testing.T) {
	p := NewMenuPositioner()

	// Cursor at bottom-right corner
	cursorPos := value.NewPosition(79, 23)
	menuWidth, menuHeight := 20, 10
	screenWidth, screenHeight := 80, 24

	result := p.CalculatePosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	// Should shift to fit
	expectedX := 60 // 80 - 20
	expectedY := 14 // 24 - 10

	if result.X() != expectedX || result.Y() != expectedY {
		t.Errorf("expected (%d,%d), got (%d,%d)", expectedX, expectedY, result.X(), result.Y())
	}
}

func TestCalculatePosition_ExactFit(t *testing.T) {
	p := NewMenuPositioner()

	// Menu exactly fills screen
	cursorPos := value.NewPosition(0, 0)
	menuWidth, menuHeight := 80, 24
	screenWidth, screenHeight := 80, 24

	result := p.CalculatePosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	// Should stay at (0,0)
	if result.X() != 0 || result.Y() != 0 {
		t.Errorf("expected (0,0), got (%d,%d)", result.X(), result.Y())
	}
}

func TestCalculatePosition_NegativeDimensions(t *testing.T) {
	p := NewMenuPositioner()

	// Negative menu dimensions (should be treated as 0)
	cursorPos := value.NewPosition(10, 5)
	menuWidth, menuHeight := -10, -5
	screenWidth, screenHeight := 80, 24

	result := p.CalculatePosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	// Should handle gracefully and return cursor position
	if result.X() != 10 {
		t.Errorf("expected x=10, got x=%d", result.X())
	}
	if result.Y() != 5 {
		t.Errorf("expected y=5, got y=%d", result.Y())
	}
}

func TestCalculatePosition_NegativeScreenDimensions(t *testing.T) {
	p := NewMenuPositioner()

	// Negative screen dimensions (edge case)
	cursorPos := value.NewPosition(10, 5)
	menuWidth, menuHeight := 20, 10
	screenWidth, screenHeight := -80, -24

	result := p.CalculatePosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	// Should handle gracefully (menu larger than "zero" screen)
	if result.X() != 0 || result.Y() != 0 {
		t.Errorf("expected (0,0), got (%d,%d)", result.X(), result.Y())
	}
}

func TestCalculatePosition_ZeroMenuDimensions(t *testing.T) {
	p := NewMenuPositioner()

	// Zero-sized menu
	cursorPos := value.NewPosition(10, 5)
	menuWidth, menuHeight := 0, 0
	screenWidth, screenHeight := 80, 24

	result := p.CalculatePosition(cursorPos, menuWidth, menuHeight, screenWidth, screenHeight)

	// Should return cursor position
	if result.X() != 10 || result.Y() != 5 {
		t.Errorf("expected (10,5), got (%d,%d)", result.X(), result.Y())
	}
}

func TestWouldOverflow_NoOverflow(t *testing.T) {
	p := NewMenuPositioner()

	position := value.NewPosition(10, 5)
	menuWidth, menuHeight := 20, 10
	screenWidth, screenHeight := 80, 24

	overflowsRight, overflowsBottom := p.WouldOverflow(position, menuWidth, menuHeight, screenWidth, screenHeight)

	if overflowsRight {
		t.Error("expected no right overflow")
	}
	if overflowsBottom {
		t.Error("expected no bottom overflow")
	}
}

func TestWouldOverflow_RightOnly(t *testing.T) {
	p := NewMenuPositioner()

	position := value.NewPosition(70, 5)
	menuWidth, menuHeight := 20, 10
	screenWidth, screenHeight := 80, 24

	overflowsRight, overflowsBottom := p.WouldOverflow(position, menuWidth, menuHeight, screenWidth, screenHeight)

	if !overflowsRight {
		t.Error("expected right overflow")
	}
	if overflowsBottom {
		t.Error("expected no bottom overflow")
	}
}

func TestWouldOverflow_BottomOnly(t *testing.T) {
	p := NewMenuPositioner()

	position := value.NewPosition(10, 20)
	menuWidth, menuHeight := 20, 10
	screenWidth, screenHeight := 80, 24

	overflowsRight, overflowsBottom := p.WouldOverflow(position, menuWidth, menuHeight, screenWidth, screenHeight)

	if overflowsRight {
		t.Error("expected no right overflow")
	}
	if !overflowsBottom {
		t.Error("expected bottom overflow")
	}
}

func TestWouldOverflow_Both(t *testing.T) {
	p := NewMenuPositioner()

	position := value.NewPosition(70, 20)
	menuWidth, menuHeight := 20, 10
	screenWidth, screenHeight := 80, 24

	overflowsRight, overflowsBottom := p.WouldOverflow(position, menuWidth, menuHeight, screenWidth, screenHeight)

	if !overflowsRight {
		t.Error("expected right overflow")
	}
	if !overflowsBottom {
		t.Error("expected bottom overflow")
	}
}

func TestFitsOnScreen_Fits(t *testing.T) {
	p := NewMenuPositioner()

	menuWidth, menuHeight := 20, 10
	screenWidth, screenHeight := 80, 24

	if !p.FitsOnScreen(menuWidth, menuHeight, screenWidth, screenHeight) {
		t.Error("expected menu to fit on screen")
	}
}

func TestFitsOnScreen_TooWide(t *testing.T) {
	p := NewMenuPositioner()

	menuWidth, menuHeight := 100, 10
	screenWidth, screenHeight := 80, 24

	if p.FitsOnScreen(menuWidth, menuHeight, screenWidth, screenHeight) {
		t.Error("expected menu not to fit on screen (too wide)")
	}
}

func TestFitsOnScreen_TooTall(t *testing.T) {
	p := NewMenuPositioner()

	menuWidth, menuHeight := 20, 30
	screenWidth, screenHeight := 80, 24

	if p.FitsOnScreen(menuWidth, menuHeight, screenWidth, screenHeight) {
		t.Error("expected menu not to fit on screen (too tall)")
	}
}

func TestFitsOnScreen_TooLarge(t *testing.T) {
	p := NewMenuPositioner()

	menuWidth, menuHeight := 100, 30
	screenWidth, screenHeight := 80, 24

	if p.FitsOnScreen(menuWidth, menuHeight, screenWidth, screenHeight) {
		t.Error("expected menu not to fit on screen (both dimensions)")
	}
}

func TestFitsOnScreen_ExactFit(t *testing.T) {
	p := NewMenuPositioner()

	menuWidth, menuHeight := 80, 24
	screenWidth, screenHeight := 80, 24

	if !p.FitsOnScreen(menuWidth, menuHeight, screenWidth, screenHeight) {
		t.Error("expected menu to fit on screen (exact fit)")
	}
}
