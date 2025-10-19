// Package model provides rich domain models for layout system.
package model

import (
	"fmt"
	"strings"

	"github.com/phoenix-tui/phoenix/layout/domain/value"
)

// FlexContainer represents a flexbox container with children.
//
// Design Philosophy:
//   - Simplified Flexbox (NOT full CSS Flexbox)
//   - Row/Column direction only (no wrap in v0.1.0)
//   - Gap support between items
//   - Immutable operations (returns new instances)
//   - Rich domain model (behavior + data)
//
// Flexbox Properties:
//   - Direction: Row (horizontal) or Column (vertical)
//   - JustifyContent: How items distribute along main axis
//   - AlignItems: How items align along cross axis
//   - Gap: Spacing between items
//
// Example:
//
//	container := NewFlexContainer(FlexDirectionRow).
//		WithJustifyContent(value.JustifyContentSpaceBetween).
//		WithAlignItems(value.AlignItemsCenter).
//		WithGap(2).
//		AddItem(NewBox("Item 1")).
//		AddItem(NewBox("Item 2"))
type FlexContainer struct {
	direction      value.FlexDirection  // Main axis direction
	justifyContent value.JustifyContent // Main axis distribution
	alignItems     value.AlignItems     // Cross axis alignment
	gap            int                  // Space between items (in cells)
	items          []*Node              // Child items (nodes wrapping boxes)
	size           value.Size           // Container size constraints
}

// NewFlexContainer creates a FlexContainer with the given direction.
//
// Defaults:
//   - JustifyContent: Start
//   - AlignItems: Stretch
//   - Gap: 0
//   - Items: Empty slice
//   - Size: Unconstrained
//
// Example:
//
//	container := NewFlexContainer(value.FlexDirectionRow)
func NewFlexContainer(direction value.FlexDirection) *FlexContainer {
	if !direction.Validate() {
		panic(fmt.Sprintf("flexbox: invalid direction %d", direction))
	}

	return &FlexContainer{
		direction:      direction,
		justifyContent: value.JustifyContentStart,
		alignItems:     value.AlignItemsStretch,
		gap:            0,
		items:          []*Node{},
		size:           value.NewSizeUnconstrained(),
	}
}

// Direction returns the flex direction (row/column).
func (f *FlexContainer) Direction() value.FlexDirection {
	return f.direction
}

// JustifyContent returns the justify content strategy.
func (f *FlexContainer) JustifyContent() value.JustifyContent {
	return f.justifyContent
}

// AlignItems returns the align items strategy.
func (f *FlexContainer) AlignItems() value.AlignItems {
	return f.alignItems
}

// Gap returns the gap between items (in cells).
func (f *FlexContainer) Gap() int {
	return f.gap
}

// Items returns a copy of the items slice (immutable).
func (f *FlexContainer) Items() []*Node {
	result := make([]*Node, len(f.items))
	copy(result, f.items)
	return result
}

// Size returns the container size constraints.
func (f *FlexContainer) Size() value.Size {
	return f.size
}

// ItemCount returns the number of items in the container.
func (f *FlexContainer) ItemCount() int {
	return len(f.items)
}

// IsEmpty returns true if the container has no items.
func (f *FlexContainer) IsEmpty() bool {
	return len(f.items) == 0
}

// WithDirection returns a new FlexContainer with the given direction.
// Panics if direction is invalid.
func (f *FlexContainer) WithDirection(direction value.FlexDirection) *FlexContainer {
	if !direction.Validate() {
		panic(fmt.Sprintf("flexbox: invalid direction %d", direction))
	}

	copy := *f
	copy.direction = direction
	return &copy
}

// WithJustifyContent returns a new FlexContainer with the given justify content.
// Panics if justify content is invalid.
func (f *FlexContainer) WithJustifyContent(justify value.JustifyContent) *FlexContainer {
	if !justify.Validate() {
		panic(fmt.Sprintf("flexbox: invalid justify content %d", justify))
	}

	copy := *f
	copy.justifyContent = justify
	return &copy
}

// WithAlignItems returns a new FlexContainer with the given align items.
// Panics if align items is invalid.
func (f *FlexContainer) WithAlignItems(align value.AlignItems) *FlexContainer {
	if !align.Validate() {
		panic(fmt.Sprintf("flexbox: invalid align items %d", align))
	}

	copy := *f
	copy.alignItems = align
	return &copy
}

// WithGap returns a new FlexContainer with the given gap.
// Gap must be non-negative (panics on negative gap).
func (f *FlexContainer) WithGap(gap int) *FlexContainer {
	if gap < 0 {
		panic(fmt.Sprintf("flexbox: gap must be non-negative, got %d", gap))
	}

	copy := *f
	copy.gap = gap
	return &copy
}

// WithSize returns a new FlexContainer with the given size constraints.
func (f *FlexContainer) WithSize(size value.Size) *FlexContainer {
	copy := *f
	copy.size = size
	return &copy
}

// AddItem returns a new FlexContainer with the given item appended.
// Panics if item is nil.
//
// Example:
//
//	container := container.AddItem(NewNode(NewBox("Item")))
func (f *FlexContainer) AddItem(item *Node) *FlexContainer {
	if item == nil {
		panic("flexbox: item cannot be nil")
	}

	copy := *f
	copy.items = make([]*Node, len(f.items)+1)
	for i, existingItem := range f.items {
		copy.items[i] = existingItem
	}
	copy.items[len(copy.items)-1] = item

	return &copy
}

// AddItems returns a new FlexContainer with multiple items appended.
// This is a convenience method for adding multiple items at once.
//
// Example:
//
//	container := container.AddItems(
//		NewNode(NewBox("Item 1")),
//		NewNode(NewBox("Item 2")),
//		NewNode(NewBox("Item 3")),
//	)
func (f *FlexContainer) AddItems(items ...*Node) *FlexContainer {
	result := f
	for _, item := range items {
		result = result.AddItem(item)
	}
	return result
}

// RemoveItem returns a new FlexContainer with the item at the given index removed.
// Panics if index is out of bounds.
func (f *FlexContainer) RemoveItem(index int) *FlexContainer {
	if index < 0 || index >= len(f.items) {
		panic(fmt.Sprintf("flexbox: index %d out of bounds (0-%d)", index, len(f.items)-1))
	}

	copy := *f
	copy.items = make([]*Node, len(f.items)-1)

	// Copy items before removed index
	for i := 0; i < index; i++ {
		copy.items[i] = f.items[i]
	}

	// Copy items after removed index
	for i := index + 1; i < len(f.items); i++ {
		copy.items[i-1] = f.items[i]
	}

	return &copy
}

// ClearItems returns a new FlexContainer with all items removed.
func (f *FlexContainer) ClearItems() *FlexContainer {
	if len(f.items) == 0 {
		return f // Already empty, return self
	}

	copy := *f
	copy.items = []*Node{}
	return &copy
}

// TotalGap calculates the total gap space between items.
// For N items, there are N-1 gaps.
//
// Example:
//
//	container with 3 items and gap=2 -> total gap = 2 * 2 = 4
func (f *FlexContainer) TotalGap() int {
	if len(f.items) <= 1 {
		return 0
	}
	return f.gap * (len(f.items) - 1)
}

// IsHorizontal returns true if the flex direction is row.
func (f *FlexContainer) IsHorizontal() bool {
	return f.direction.IsHorizontal()
}

// IsVertical returns true if the flex direction is column.
func (f *FlexContainer) IsVertical() bool {
	return f.direction.IsVertical()
}

// String returns a human-readable debug representation.
func (f *FlexContainer) String() string {
	var parts []string

	parts = append(parts, fmt.Sprintf("direction=%s", f.direction))

	if !f.justifyContent.IsDefault() {
		parts = append(parts, fmt.Sprintf("justify=%s", f.justifyContent))
	}

	if !f.alignItems.IsDefault() {
		parts = append(parts, fmt.Sprintf("align=%s", f.alignItems))
	}

	if f.gap > 0 {
		parts = append(parts, fmt.Sprintf("gap=%d", f.gap))
	}

	if !f.size.IsUnconstrained() {
		parts = append(parts, fmt.Sprintf("size=%s", f.size))
	}

	parts = append(parts, fmt.Sprintf("items=%d", len(f.items)))

	return fmt.Sprintf("FlexContainer{%s}", strings.Join(parts, " "))
}
