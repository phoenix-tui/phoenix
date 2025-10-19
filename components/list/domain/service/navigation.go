package service

// NavigationService handles list navigation logic
type NavigationService struct{}

// NewNavigationService creates a new navigation service
func NewNavigationService() *NavigationService {
	return &NavigationService{}
}

// MoveUp moves the focus up by one item with wrap-around
func (s *NavigationService) MoveUp(currentIndex, itemCount int) int {
	if itemCount == 0 {
		return 0
	}
	if currentIndex <= 0 {
		return itemCount - 1 // Wrap to end
	}
	return currentIndex - 1
}

// MoveDown moves the focus down by one item with wrap-around
func (s *NavigationService) MoveDown(currentIndex, itemCount int) int {
	if itemCount == 0 {
		return 0
	}
	if currentIndex >= itemCount-1 {
		return 0 // Wrap to start
	}
	return currentIndex + 1
}

// MovePageUp moves the focus up by page size
func (s *NavigationService) MovePageUp(currentIndex, pageSize, itemCount int) int {
	if itemCount == 0 {
		return 0
	}
	newIndex := currentIndex - pageSize
	if newIndex < 0 {
		return 0
	}
	return newIndex
}

// MovePageDown moves the focus down by page size
func (s *NavigationService) MovePageDown(currentIndex, pageSize, itemCount int) int {
	if itemCount == 0 {
		return 0
	}
	newIndex := currentIndex + pageSize
	if newIndex >= itemCount {
		return itemCount - 1
	}
	return newIndex
}

// MoveToStart moves the focus to the first item
func (s *NavigationService) MoveToStart() int {
	return 0
}

// MoveToEnd moves the focus to the last item
func (s *NavigationService) MoveToEnd(itemCount int) int {
	if itemCount == 0 {
		return 0
	}
	return itemCount - 1
}

// CalculateScrollOffset calculates the scroll offset to keep the focused item visible
// Returns the offset from the start of the list
func (s *NavigationService) CalculateScrollOffset(focusedIndex, visibleHeight, itemCount int) int {
	if itemCount == 0 || visibleHeight <= 0 {
		return 0
	}

	// If all items fit in the visible area, no scrolling needed
	if itemCount <= visibleHeight {
		return 0
	}

	// Calculate the current scroll offset to keep focused item visible
	// We want to keep the focused item in the middle of the viewport when possible

	// Try to center the focused item
	centerOffset := focusedIndex - visibleHeight/2
	if centerOffset < 0 {
		return 0
	}
	if centerOffset > itemCount-visibleHeight {
		return itemCount - visibleHeight
	}

	return centerOffset
}
