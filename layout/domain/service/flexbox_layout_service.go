// Package service provides domain services for layout calculations.
package service

import (
	"github.com/phoenix-tui/phoenix/layout/domain/model"
	"github.com/phoenix-tui/phoenix/layout/domain/value"
)

// FlexboxLayoutService calculates layout for flexbox containers.
//
// Design Philosophy:
//   - Simplified Flexbox (NOT full CSS Flexbox)
//   - Two-pass algorithm: measure items, then position them
//   - Supports row/column direction
//   - Handles justify-content and align-items
//   - Gap support between items
//
// Algorithm Overview:
//  1. Measure all items (get natural sizes)
//  2. Calculate main axis distribution (justify-content)
//  3. Calculate cross axis positioning (align-items)
//  4. Apply gap spacing
//  5. Set final positions on nodes
//
// Example:
//
//	service := NewFlexboxLayoutService(measureService)
//	laidOut := service.Layout(container, 80, 24)
type FlexboxLayoutService struct {
	measureService *MeasureService
}

// NewFlexboxLayoutService creates a FlexboxLayoutService.
func NewFlexboxLayoutService(measureService *MeasureService) *FlexboxLayoutService {
	return &FlexboxLayoutService{
		measureService: measureService,
	}
}

// Layout calculates positions for all items in a flexbox container.
// Returns a new container with updated item positions.
//
// Parameters:
//   - container: The flexbox container to layout
//   - containerWidth: Available width for the container
//   - containerHeight: Available height for the container
//
// Returns:
//   - New FlexContainer with positioned items
func (f *FlexboxLayoutService) Layout(
	container *model.FlexContainer,
	containerWidth, containerHeight int,
) *model.FlexContainer {
	if container.IsEmpty() {
		return container // Nothing to layout
	}

	// 1. Measure all items
	itemSizes := f.measureItems(container)

	// 2. Calculate main axis positions
	mainAxisPositions := f.calculateMainAxisPositions(
		container,
		itemSizes,
		containerWidth,
		containerHeight,
	)

	// 3. Calculate cross axis positions
	crossAxisPositions := f.calculateCrossAxisPositions(
		container,
		itemSizes,
		containerWidth,
		containerHeight,
	)

	// 4. Apply positions to items
	newItems := make([]*model.Node, len(container.Items()))
	items := container.Items()

	for i := 0; i < len(items); i++ {
		var x, y int

		if container.IsHorizontal() {
			x = mainAxisPositions[i]
			y = crossAxisPositions[i]
		} else {
			x = crossAxisPositions[i]
			y = mainAxisPositions[i]
		}

		newItems[i] = items[i].SetPosition(value.NewPosition(x, y))
	}

	// 5. Return new container with positioned items
	result := container.ClearItems()
	for _, item := range newItems {
		result = result.AddItem(item)
	}

	return result
}

// measureItems measures all items and returns their sizes.
func (f *FlexboxLayoutService) measureItems(container *model.FlexContainer) []value.Size {
	items := container.Items()
	sizes := make([]value.Size, len(items))

	for i, item := range items {
		sizes[i] = f.measureService.Measure(item.Box())
	}

	return sizes
}

// calculateMainAxisPositions calculates positions along the main axis (row/column direction).
// This implements justify-content (start, end, center, space-between).
func (f *FlexboxLayoutService) calculateMainAxisPositions(
	container *model.FlexContainer,
	itemSizes []value.Size,
	containerWidth, containerHeight int,
) []int {
	positions := make([]int, len(itemSizes))
	if len(itemSizes) == 0 {
		return positions
	}

	// Calculate total size of items along main axis
	var totalItemSize int
	for _, size := range itemSizes {
		if container.IsHorizontal() {
			totalItemSize += size.Width()
		} else {
			totalItemSize += size.Height()
		}
	}

	// Add gap spacing
	totalGap := container.TotalGap()
	totalItemSize += totalGap

	// Get container size along main axis
	var containerSize int
	if container.IsHorizontal() {
		containerSize = containerWidth
	} else {
		containerSize = containerHeight
	}

	// Calculate remaining space
	remainingSpace := containerSize - totalItemSize
	if remainingSpace < 0 {
		remainingSpace = 0
	}

	// Apply justify-content strategy
	justify := container.JustifyContent()
	gap := container.Gap()

	switch justify {
	case value.JustifyContentStart:
		// Pack items at start with gap spacing
		pos := 0
		for i, size := range itemSizes {
			positions[i] = pos
			if container.IsHorizontal() {
				pos += size.Width() + gap
			} else {
				pos += size.Height() + gap
			}
		}

	case value.JustifyContentEnd:
		// Pack items at end with gap spacing
		pos := remainingSpace
		for i, size := range itemSizes {
			positions[i] = pos
			if container.IsHorizontal() {
				pos += size.Width() + gap
			} else {
				pos += size.Height() + gap
			}
		}

	case value.JustifyContentCenter:
		// Center items with gap spacing
		startPos := remainingSpace / 2
		pos := startPos
		for i, size := range itemSizes {
			positions[i] = pos
			if container.IsHorizontal() {
				pos += size.Width() + gap
			} else {
				pos += size.Height() + gap
			}
		}

	case value.JustifyContentSpaceBetween:
		// Equal spacing between items
		if len(itemSizes) == 1 {
			// Single item: position at start
			positions[0] = 0
		} else {
			// Calculate gap between items (excluding original gap)
			gapBetween := remainingSpace / (len(itemSizes) - 1)

			pos := 0
			for i, size := range itemSizes {
				positions[i] = pos
				if container.IsHorizontal() {
					pos += size.Width() + gap + gapBetween
				} else {
					pos += size.Height() + gap + gapBetween
				}
			}
		}
	}

	return positions
}

// calculateCrossAxisPositions calculates positions along the cross axis.
// This implements align-items (start, end, center, stretch).
func (f *FlexboxLayoutService) calculateCrossAxisPositions(
	container *model.FlexContainer,
	itemSizes []value.Size,
	containerWidth, containerHeight int,
) []int {
	positions := make([]int, len(itemSizes))
	if len(itemSizes) == 0 {
		return positions
	}

	// Get container size along cross axis
	var containerSize int
	if container.IsHorizontal() {
		containerSize = containerHeight
	} else {
		containerSize = containerWidth
	}

	// Apply align-items strategy
	align := container.AlignItems()

	for i, size := range itemSizes {
		var itemSize int
		if container.IsHorizontal() {
			itemSize = size.Height()
		} else {
			itemSize = size.Width()
		}

		switch align {
		case value.AlignItemsStart:
			// Align at start (top for row, left for column)
			positions[i] = 0

		case value.AlignItemsEnd:
			// Align at end (bottom for row, right for column)
			positions[i] = containerSize - itemSize

		case value.AlignItemsCenter:
			// Center along cross axis
			positions[i] = (containerSize - itemSize) / 2

		case value.AlignItemsStretch:
			// Stretch to fill (position at 0, size adjusted elsewhere)
			positions[i] = 0
		}

		// Clamp to non-negative
		if positions[i] < 0 {
			positions[i] = 0
		}
	}

	return positions
}

// LayoutResult represents the result of a flexbox layout calculation.
// This is useful for debugging and testing.
type LayoutResult struct {
	Container     *model.FlexContainer
	ItemPositions []value.Position
	ItemSizes     []value.Size
}

// LayoutWithDetails performs layout and returns detailed results.
// This is useful for debugging and testing.
func (f *FlexboxLayoutService) LayoutWithDetails(
	container *model.FlexContainer,
	containerWidth, containerHeight int,
) LayoutResult {
	laidOutContainer := f.Layout(container, containerWidth, containerHeight)

	// Extract positions and sizes
	items := laidOutContainer.Items()
	positions := make([]value.Position, len(items))
	sizes := make([]value.Size, len(items))

	for i, item := range items {
		positions[i] = item.Position()
		totalSize := item.Box().TotalSize()
		sizes[i] = totalSize
	}

	return LayoutResult{
		Container:     laidOutContainer,
		ItemPositions: positions,
		ItemSizes:     sizes,
	}
}
