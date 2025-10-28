package service

import (
	model2 "github.com/phoenix-tui/phoenix/layout/internal/domain/model"
	value2 "github.com/phoenix-tui/phoenix/layout/internal/domain/value"
)

// LayoutService calculates positions of boxes within parent containers.
// Applies alignment rules and handles overflow (clamping to bounds).
//
// Design Philosophy:
//   - Domain service (pure business logic)
//   - Uses MeasureService for size calculation
//   - Alignment-based positioning (CSS-like)
//   - Overflow handling via clamping
//   - Recursive tree positioning
//
// Positioning Process:
//  1. Measure box size (via MeasureService)
//  2. Calculate available space in parent
//  3. Apply alignment to determine position
//  4. Clamp position to parent bounds (prevent overflow)
//
// Example:
//
//	ms := NewMeasureService(unicodeService)
//	ls := NewLayoutService(ms)
//
//	box := model.NewBox("Centered Text").
//		WithAlignment(value.NewAlignmentCenter())
//
//	parentSize := value.NewSizeExact(80, 24)
//	position := ls.Layout(box, parentSize)
//	// Position calculated based on center alignment
type LayoutService struct {
	measureService *MeasureService
}

// NewLayoutService creates a new LayoutService.
// Panics if measureService is nil.
func NewLayoutService(ms *MeasureService) *LayoutService {
	if ms == nil {
		panic("layout_service: measure service cannot be nil")
	}
	return &LayoutService{measureService: ms}
}

// Layout calculates the position of a box within parent bounds.
// Uses alignment to determine placement and clamps to prevent overflow.
//
// Algorithm:
//  1. Measure box size
//  2. Calculate alignment offsets
//  3. Clamp to parent bounds
//
// Alignment examples (box width 10, parent width 80):
//   - Left: x = 0
//   - Center: x = (80 - 10) / 2 = 35
//   - Right: x = 80 - 10 = 70
//
// Overflow handling:
//   - If box width > parent width: x = 0 (left-aligned)
//   - If position < 0: clamped to 0
//   - If position + size > parent: clamped to parent - size
//
// Parameters:
//   - box: The box to position
//   - parentSize: Available space in parent container
//
// Returns:
//   - Position within parent (0-based coordinates)
func (ls *LayoutService) Layout(box *model2.Box, parentSize value2.Size) value2.Position {
	// Step 1: Measure box size
	boxSize := ls.measureService.Measure(box)

	// Step 2: Calculate alignment offsets
	alignment := box.Alignment()
	xOffset, yOffset := alignment.CalculateOffsets(
		boxSize.Width(), boxSize.Height(),
		parentSize.Width(), parentSize.Height(),
	)

	// Step 3: Create position (clamping handled by Position constructor)
	return value2.NewPosition(xOffset, yOffset)
}

// LayoutNode positions a node and all its children recursively.
// Creates a new positioned node tree with updated positions.
//
// Algorithm:
//  1. Layout root node's box
//  2. For each child:
//     - Calculate child's available space (parent size - root box margins)
//     - Layout child recursively
//     - Offset child position by root box position
//  3. Return new node tree with positions
//
// Note: This is a simplified version for Day 3.
// Full flexbox/grid positioning will be implemented in Week 9-10 continuation.
//
// Parameters:
//   - node: Root node to position
//   - parentSize: Available space in parent container
//
// Returns:
//   - New node tree with positions calculated
func (ls *LayoutService) LayoutNode(node *model2.Node, parentSize value2.Size) *model2.Node {
	// Step 1: Layout root box
	rootBox := node.Box()
	rootPosition := ls.Layout(rootBox, parentSize)

	// Step 2: Layout children
	// For now, we stack children vertically (simplified)
	// Full flexbox/grid will be implemented later
	children := node.Children()
	result := model2.NewNode(rootBox).SetPosition(rootPosition)

	currentY := 0
	for _, child := range children {
		childBox := child.Box()

		// Calculate available space for child
		// (parent width, remaining height)
		childParentSize := value2.NewSizeExact(
			parentSize.Width(),
			parentSize.Height()-currentY,
		)

		// Layout child recursively
		positionedChild := ls.LayoutNode(child, childParentSize)

		// Offset child position by current Y
		childPosition := positionedChild.Position().Add(0, currentY)

		// Update child with new position
		positionedChild = positionedChild.SetPosition(childPosition)

		// Add positioned child to result
		result = result.AddChild(positionedChild)

		// Update current Y for next child (vertical stacking)
		childSize := ls.measureService.Measure(childBox)
		currentY += childSize.Height()
	}

	return result
}

// CalculatePosition calculates position based on alignment (low-level helper).
// This is exposed for testing and advanced usage.
//
// Parameters:
//   - boxSize: Size of the box to position
//   - parentSize: Size of the parent container
//   - alignment: Alignment within parent
//
// Returns:
//   - Position within parent (0-based coordinates)
func (ls *LayoutService) CalculatePosition(
	boxSize value2.Size,
	parentSize value2.Size,
	alignment value2.Alignment,
) value2.Position {
	xOffset, yOffset := alignment.CalculateOffsets(
		boxSize.Width(), boxSize.Height(),
		parentSize.Width(), parentSize.Height(),
	)
	return value2.NewPosition(xOffset, yOffset)
}

// ClampPosition clamps position to prevent overflow.
// This is a helper for advanced positioning scenarios.
//
// Parameters:
//   - position: Position to clamp
//   - boxSize: Size of the box
//   - parentSize: Size of the parent container
//
// Returns:
//   - Clamped position ensuring box stays within parent bounds
func (ls *LayoutService) ClampPosition(
	position value2.Position,
	boxSize value2.Size,
	parentSize value2.Size,
) value2.Position {
	x := position.X()
	y := position.Y()

	// Clamp X
	if x < 0 {
		x = 0
	}
	maxX := parentSize.Width() - boxSize.Width()
	if maxX < 0 {
		maxX = 0
	}
	if x > maxX {
		x = maxX
	}

	// Clamp Y
	if y < 0 {
		y = 0
	}
	maxY := parentSize.Height() - boxSize.Height()
	if maxY < 0 {
		maxY = 0
	}
	if y > maxY {
		y = maxY
	}

	return value2.NewPosition(x, y)
}
