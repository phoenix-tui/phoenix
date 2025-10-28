// Package service provides scroll management services.
package service

// ScrollService handles scrolling logic for viewports.
// It provides pure functions for calculating scroll positions and visible content.
type ScrollService struct{}

// NewScrollService creates a new ScrollService instance.
func NewScrollService() *ScrollService {
	return &ScrollService{}
}

// VisibleLines returns the lines that should be visible in the viewport.
// given the content, current offset, and viewport height.
func (s *ScrollService) VisibleLines(content []string, offset, height int) []string {
	if len(content) == 0 || height <= 0 {
		return []string{}
	}

	// Clamp offset to valid range.
	if offset < 0 {
		offset = 0
	}
	if offset >= len(content) {
		offset = len(content) - 1
	}

	// Calculate end index.
	end := offset + height
	if end > len(content) {
		end = len(content)
	}

	return content[offset:end]
}

// MaxScrollOffset calculates the maximum valid scroll offset.
// This is the total number of lines minus the viewport height.
// Returns 0 if content fits entirely in the viewport.
func (s *ScrollService) MaxScrollOffset(totalLines, viewportHeight int) int {
	if totalLines <= viewportHeight {
		return 0
	}
	return totalLines - viewportHeight
}

// ScrollUp calculates the new offset after scrolling up by the given number of lines.
// The result is clamped to 0.
func (s *ScrollService) ScrollUp(currentOffset, lines int) int {
	newOffset := currentOffset - lines
	if newOffset < 0 {
		return 0
	}
	return newOffset
}

// ScrollDown calculates the new offset after scrolling down by the given number of lines.
// The result is clamped to maxOffset.
func (s *ScrollService) ScrollDown(currentOffset, lines, maxOffset int) int {
	newOffset := currentOffset + lines
	if newOffset > maxOffset {
		return maxOffset
	}
	return newOffset
}

// FollowModeOffset calculates the scroll offset for follow mode (tail -f style).
// This keeps the viewport at the bottom of the content.
func (s *ScrollService) FollowModeOffset(totalLines, viewportHeight int) int {
	return s.MaxScrollOffset(totalLines, viewportHeight)
}

// CanScrollUp returns true if the viewport can scroll up from the current offset.
func (s *ScrollService) CanScrollUp(currentOffset int) bool {
	return currentOffset > 0
}

// CanScrollDown returns true if the viewport can scroll down from the current offset.
func (s *ScrollService) CanScrollDown(currentOffset, totalLines, viewportHeight int) bool {
	maxOffset := s.MaxScrollOffset(totalLines, viewportHeight)
	return currentOffset < maxOffset
}

// IsAtTop returns true if the viewport is at the top (offset = 0).
func (s *ScrollService) IsAtTop(currentOffset int) bool {
	return currentOffset == 0
}

// IsAtBottom returns true if the viewport is at the bottom (offset = maxOffset).
func (s *ScrollService) IsAtBottom(currentOffset, totalLines, viewportHeight int) bool {
	maxOffset := s.MaxScrollOffset(totalLines, viewportHeight)
	return currentOffset >= maxOffset
}
