// Package service provides domain services for layout calculations.
// Domain services contain business logic that doesn't naturally fit in entities/value objects.
package service

import (
	"strings"

	"github.com/phoenix-tui/phoenix/core"
	"github.com/phoenix-tui/phoenix/layout/internal/domain/model"
	"github.com/phoenix-tui/phoenix/layout/internal/domain/value"
)

// MeasureService calculates natural sizes of boxes using Unicode-aware width calculation.
// This service integrates with phoenix/core.UnicodeService for correct visual width.
//
// Design Philosophy:
//   - Domain service (pure business logic, no side effects)
//   - Unicode-aware via phoenix/core integration
//   - Follows CSS box model layering
//   - Returns constrained sizes (respects min/max)
//
// Measurement Process:
//  1. Measure content width using UnicodeService (correct for CJK/emoji)
//  2. Add padding (left + right, top + bottom)
//  3. Add border if enabled (2 cells horizontally, 2 vertically)
//  4. Add margin (left + right, top + bottom)
//  5. Apply size constraints (min/max clamping)
//
// Example:
//
//	us := coreService.NewUnicodeService()
//	ms := NewMeasureService(us)
//
//	box := model.NewBox("Hello 世界").
//		WithPadding(value.NewSpacingAll(1)).
//		WithBorder(true)
//
//	size := ms.Measure(box)
//	// Content: "Hello 世界" = 5 + 4 = 9 cells wide
//	// + padding: 1 left + 1 right = 2
//	// + border: 1 left + 1 right = 2
//	// Total width: 9 + 2 + 2 = 13 cells
type MeasureService struct {
}

// NewMeasureService creates a new MeasureService.
// Panics if unicodeService is nil.
func NewMeasureService() *MeasureService {
	return &MeasureService{}
}

// Measure calculates the total size of a box (content + padding + border + margin).
// The size is constrained by the box's size constraints (min/max).
//
// Calculation steps:
//  1. Content size (Unicode-aware width calculation)
//  2. + Padding (both sides)
//  3. + Border (1 cell per side if enabled)
//  4. + Margin (both sides)
//  5. Apply size constraints
//
// Multi-line content:
//   - Width: Maximum line width (Unicode-aware)
//   - Height: Number of lines
//
// Example:
//
//	box := model.NewBox("Hi\n世界")  // 2 lines, max width 4 (CJK chars)
//	size := ms.Measure(box)          // Size{width=4, height=2}
//
// Returns:
//   - Size with width and height including all box model layers
func (ms *MeasureService) Measure(box *model.Box) value.Size {
	// Step 1: Measure content (Unicode-aware)
	contentWidth, contentHeight := ms.measureContent(box.Content())

	// Step 2: Add padding (explicit + implicit for borders)
	// When border is enabled, add 1-space aesthetic padding between border and content
	padding := box.Padding()
	totalPaddingHorizontal := padding.Horizontal()
	totalPaddingVertical := padding.Vertical()

	if box.HasBorder() {
		// Add implicit aesthetic spacing (1 cell per side)
		totalPaddingHorizontal += 2 // 1 left + 1 right
		totalPaddingVertical += 2   // 1 top + 1 bottom
	}

	width := contentWidth + totalPaddingHorizontal
	height := contentHeight + totalPaddingVertical

	// Step 3: Add border characters (1 cell per side if enabled)
	if box.HasBorder() {
		width += 2  // Border characters
		height += 2 // Border characters
	}

	// Step 4: Add margin
	margin := box.Margin()
	width += margin.Horizontal()
	height += margin.Vertical()

	// Step 5: Apply size constraints
	sizeConstraints := box.Size()
	finalWidth, finalHeight := sizeConstraints.Constrain(width, height)

	return value.NewSizeExact(finalWidth, finalHeight)
}

// measureContent calculates the size of content text (Unicode-aware).
// Multi-line content is split by \n and max width is calculated.
//
// Returns:
//   - width: Maximum line width (Unicode-aware via UnicodeService)
//   - height: Number of lines
func (ms *MeasureService) measureContent(content string) (width, height int) {
	if content == "" {
		return 0, 0
	}

	// Split by newlines
	lines := strings.Split(content, "\n")
	height = len(lines)

	// Find maximum width (Unicode-aware)
	maxWidth := 0
	for _, line := range lines {
		lineWidth := core.StringWidth(line)
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}

	return maxWidth, height
}

// MeasureContent is a public helper for measuring raw content text.
// This is useful for measuring text before creating a box.
//
// Example:
//
//	ms := NewMeasureService(unicodeService)
//	width, height := ms.MeasureContent("Hello 世界")
//	// width = 9 (5 ASCII + 2 CJK chars = 5 + 4 = 9)
//	// height = 1
func (ms *MeasureService) MeasureContent(content string) (width, height int) {
	return ms.measureContent(content)
}
