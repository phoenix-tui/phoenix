// Package layout provides a user-friendly API for creating terminal UI layouts.
// It implements the CSS box model with support for padding, margins, borders, and alignment.
//
// Basic Usage:
//
//	box := layout.NewBox("Hello World").
//		PaddingAll(1).
//		Border().
//		AlignCenter().
//		Render()
//
// The API hides the underlying DDD architecture complexity and provides
// a fluent, chainable interface for building layouts.
package layout

import (
	"strings"

	model2 "github.com/phoenix-tui/phoenix/layout/internal/domain/model"
	service2 "github.com/phoenix-tui/phoenix/layout/internal/domain/service"
	value2 "github.com/phoenix-tui/phoenix/layout/internal/domain/value"
)

// Box is the public API for creating layout boxes.
// It wraps the domain model with a fluent builder API.
//
// Box implements the CSS box model:
//   - Content: The actual text content
//   - Padding: Space between content and border
//   - Border: Visual boundary around the box
//   - Margin: Space outside the border
//
// Example:
//
//	box := layout.NewBox("Hello").
//		Width(20).
//		PaddingAll(1).
//		Border().
//		MarginAll(2).
//		AlignCenter().
//		Render()
type Box struct {
	domain *model2.Box
}

// NewBox creates a new Box with the given content.
// All styling defaults to zero/unconstrained.
//
// Example:
//
//	box := layout.NewBox("Hello World")
func NewBox(content string) *Box {
	return &Box{
		domain: model2.NewBox(content),
	}
}

// ============================================================================
// Size Constraints
// ============================================================================

// Width sets an exact width constraint.
// The box will be exactly this width (in cells).
//
// Example:
//
//	box := layout.NewBox("Hi").Width(20)  // Forces width to 20 cells
func (b *Box) Width(width int) *Box {
	b.domain = b.domain.WithSize(b.domain.Size().WithWidth(width))
	return b
}

// Height sets an exact height constraint.
// The box will be exactly this height (in lines).
//
// Example:
//
//	box := layout.NewBox("Hi").Height(5)  // Forces height to 5 lines
func (b *Box) Height(height int) *Box {
	b.domain = b.domain.WithSize(b.domain.Size().WithHeight(height))
	return b
}

// MinWidth sets a minimum width constraint.
// The box will be at least this width.
//
// Example:
//
//	box := layout.NewBox("Hi").MinWidth(10)  // At least 10 cells wide
func (b *Box) MinWidth(minWidth int) *Box {
	b.domain = b.domain.WithSize(b.domain.Size().WithMinWidth(minWidth))
	return b
}

// MinHeight sets a minimum height constraint.
// The box will be at least this height.
//
// Example:
//
//	box := layout.NewBox("Hi").MinHeight(3)  // At least 3 lines tall
func (b *Box) MinHeight(minHeight int) *Box {
	b.domain = b.domain.WithSize(b.domain.Size().WithMinHeight(minHeight))
	return b
}

// MaxWidth sets a maximum width constraint.
// The box will be at most this width.
//
// Example:
//
//	box := layout.NewBox("Very long text...").MaxWidth(20)  // Max 20 cells
func (b *Box) MaxWidth(maxWidth int) *Box {
	b.domain = b.domain.WithSize(b.domain.Size().WithMaxWidth(maxWidth))
	return b
}

// MaxHeight sets a maximum height constraint.
// The box will be at most this height.
//
// Example:
//
//	box := layout.NewBox("Line1\nLine2\nLine3").MaxHeight(2)  // Max 2 lines
func (b *Box) MaxHeight(maxHeight int) *Box {
	b.domain = b.domain.WithSize(b.domain.Size().WithMaxHeight(maxHeight))
	return b
}

// ============================================================================
// Padding (space inside border)
// ============================================================================

// Padding sets padding for each side individually.
// Padding adds space between content and border.
//
// Parameters: top, right, bottom, left (clockwise from top, like CSS)
//
// Example:
//
//	box := layout.NewBox("Hi").Padding(1, 2, 1, 2)  // Vertical: 1, Horizontal: 2
func (b *Box) Padding(top, right, bottom, left int) *Box {
	b.domain = b.domain.WithPadding(value2.NewSpacing(top, right, bottom, left))
	return b
}

// PaddingAll sets the same padding for all sides.
//
// Example:
//
//	box := layout.NewBox("Hi").PaddingAll(2)  // 2 cells padding on all sides
func (b *Box) PaddingAll(padding int) *Box {
	b.domain = b.domain.WithPadding(value2.NewSpacingAll(padding))
	return b
}

// PaddingVH sets vertical and horizontal padding.
//
// Example:
//
//	box := layout.NewBox("Hi").PaddingVH(1, 3)  // 1 vertical, 3 horizontal
func (b *Box) PaddingVH(vertical, horizontal int) *Box {
	b.domain = b.domain.WithPadding(value2.NewSpacingVH(vertical, horizontal))
	return b
}

// ============================================================================
// Border
// ============================================================================

// Border enables a border around the box.
// The border uses Unicode box drawing characters (┌─┐│└┘).
//
// Note: Borders automatically add 1 cell of aesthetic padding between
// the border and content on all sides (implicit padding).
//
// Example:
//
//	box := layout.NewBox("Hi").Border()  // Draws a border around content
func (b *Box) Border() *Box {
	b.domain = b.domain.WithBorder(true)
	return b
}

// NoBorder explicitly disables the border.
// This is the default, but can be used to override previous settings.
//
// Example:
//
//	box := layout.NewBox("Hi").Border().NoBorder()  // Border disabled
func (b *Box) NoBorder() *Box {
	b.domain = b.domain.WithBorder(false)
	return b
}

// ============================================================================
// Margin (space outside border)
// ============================================================================

// Margin sets margin for each side individually.
// Margin adds space outside the border.
//
// Parameters: top, right, bottom, left (clockwise from top, like CSS)
//
// Example:
//
//	box := layout.NewBox("Hi").Margin(2, 4, 2, 4)  // Vertical: 2, Horizontal: 4
func (b *Box) Margin(top, right, bottom, left int) *Box {
	b.domain = b.domain.WithMargin(value2.NewSpacing(top, right, bottom, left))
	return b
}

// MarginAll sets the same margin for all sides.
//
// Example:
//
//	box := layout.NewBox("Hi").MarginAll(3)  // 3 cells margin on all sides
func (b *Box) MarginAll(margin int) *Box {
	b.domain = b.domain.WithMargin(value2.NewSpacingAll(margin))
	return b
}

// MarginVH sets vertical and horizontal margin.
//
// Example:
//
//	box := layout.NewBox("Hi").MarginVH(2, 5)  // 2 vertical, 5 horizontal
func (b *Box) MarginVH(vertical, horizontal int) *Box {
	b.domain = b.domain.WithMargin(value2.NewSpacingVH(vertical, horizontal))
	return b
}

// ============================================================================
// Alignment
// ============================================================================

// AlignLeft aligns content to the left (default).
//
// Example:
//
//	box := layout.NewBox("Hi").AlignLeft()
func (b *Box) AlignLeft() *Box {
	b.domain = b.domain.WithAlignment(value2.NewAlignmentDefault())
	return b
}

// AlignCenter centers content horizontally and vertically.
//
// Example:
//
//	box := layout.NewBox("Hi").AlignCenter()
func (b *Box) AlignCenter() *Box {
	b.domain = b.domain.WithAlignment(value2.NewAlignmentCenter())
	return b
}

// AlignRight aligns content to the right.
//
// Example:
//
//	box := layout.NewBox("Hi").AlignRight()
func (b *Box) AlignRight() *Box {
	align := value2.NewAlignment(value2.AlignRight, value2.AlignTop)
	b.domain = b.domain.WithAlignment(align)
	return b
}

// AlignTop aligns content to the top (default).
//
// Example:
//
//	box := layout.NewBox("Hi").AlignTop()
func (b *Box) AlignTop() *Box {
	align := value2.NewAlignment(value2.AlignLeft, value2.AlignTop)
	b.domain = b.domain.WithAlignment(align)
	return b
}

// AlignMiddle centers content vertically.
//
// Example:
//
//	box := layout.NewBox("Hi").AlignMiddle()
func (b *Box) AlignMiddle() *Box {
	align := value2.NewAlignment(value2.AlignLeft, value2.AlignMiddle)
	b.domain = b.domain.WithAlignment(align)
	return b
}

// AlignBottom aligns content to the bottom.
//
// Example:
//
//	box := layout.NewBox("Hi").AlignBottom()
func (b *Box) AlignBottom() *Box {
	align := value2.NewAlignment(value2.AlignLeft, value2.AlignBottom)
	b.domain = b.domain.WithAlignment(align)
	return b
}

// Align sets both horizontal and vertical alignment.
//
// Example:
//
//	box := layout.NewBox("Hi").Align(layout.AlignCenter, layout.AlignMiddle)
func (b *Box) Align(horizontal value2.HorizontalAlignment, vertical value2.VerticalAlignment) *Box {
	b.domain = b.domain.WithAlignment(value2.NewAlignment(horizontal, vertical))
	return b
}

// ============================================================================
// Rendering
// ============================================================================

// Render generates the final string output for this box.
// This applies all styling (padding, border, margin, alignment) and
// returns the rendered string.
//
// Example:
//
//	output := layout.NewBox("Hello World").
//		PaddingAll(1).
//		Border().
//		AlignCenter().
//		Render()
//	fmt.Println(output)
func (b *Box) Render() string {
	// Create services (these are lightweight)
	measureService := service2.NewMeasureService()
	renderService := service2.NewRenderService()

	// Measure box to get natural size
	size := measureService.Measure(b.domain)

	// Update domain box with measured size (for rendering)
	b.domain = b.domain.WithSize(value2.NewSizeExact(size.Width(), size.Height()))

	// Render the box
	return renderService.Render(b.domain)
}

// String implements fmt.Stringer for convenient printing.
// Equivalent to calling Render().
func (b *Box) String() string {
	return b.Render()
}

// ============================================================================
// Layout (positioning within parent)
// ============================================================================

// Layout positions this box within a parent container.
// Returns the calculated position (x, y) where the box should be rendered.
//
// Parameters:
//   - parentWidth: Width of parent container (in cells)
//   - parentHeight: Height of parent container (in lines)
//
// The position is calculated based on the box's alignment settings.
//
// Example:
//
//	box := layout.NewBox("Hi").AlignCenter()
//	pos := box.Layout(80, 24)  // Position in 80x24 terminal
//	fmt.Printf("Render at (%d, %d)\n", pos.X(), pos.Y())
func (b *Box) Layout(parentWidth, parentHeight int) value2.Position {
	// Create services
	measureService := service2.NewMeasureService()
	layoutService := service2.NewLayoutService(measureService)

	// Layout the box
	parentSize := value2.NewSizeExact(parentWidth, parentHeight)
	return layoutService.Layout(b.domain, parentSize)
}

// ============================================================================
// Advanced: Direct access to domain model (for advanced use cases)
// ============================================================================

// Domain returns the underlying domain model.
// This is provided for advanced use cases where you need direct access
// to the domain layer (e.g., for custom layout algorithms).
//
// Most users should not need this method.
func (b *Box) Domain() *model2.Box {
	return b.domain
}

// ============================================================================
// Flexbox Layout API
// ============================================================================

// Flex is the public API for creating flexbox layouts.
// It wraps the domain FlexContainer with a fluent builder API.
//
// Flexbox provides row/column layouts with flexible sizing and gap support.
//
// Simplified from CSS Flexbox (v0.1.0):
//   - Row/Column direction
//   - Justify content (start, end, center, space-between)
//   - Align items (start, end, center, stretch)
//   - Gap between items
//
// Example (horizontal split):
//
//	flex := layout.Row().
//		Gap(2).
//		JustifyStart().
//		AlignStretch().
//		Add(layout.NewBox("Left Panel")).
//		Add(layout.NewBox("Right Panel")).
//		Render(80, 24)
//
// Example (vertical stack):
//
//	flex := layout.Column().
//		Gap(1).
//		JustifyCenter().
//		Add(layout.NewBox("Header")).
//		Add(layout.NewBox("Content")).
//		Add(layout.NewBox("Footer")).
//		Render(80, 24)
type Flex struct {
	domain *model2.FlexContainer
}

// Row creates a new horizontal flexbox container (direction: row).
// Items are arranged left-to-right.
//
// Example:
//
//	flex := layout.Row().Add(box1).Add(box2)
func Row() *Flex {
	return &Flex{
		domain: model2.NewFlexContainer(value2.FlexDirectionRow),
	}
}

// Column creates a new vertical flexbox container (direction: column).
// Items are arranged top-to-bottom.
//
// Example:
//
//	flex := layout.Column().Add(box1).Add(box2)
func Column() *Flex {
	return &Flex{
		domain: model2.NewFlexContainer(value2.FlexDirectionColumn),
	}
}

// ============================================================================
// Item Management
// ============================================================================

// Add adds a box to the flexbox container.
//
// Example:
//
//	flex := layout.Row().
//		Add(layout.NewBox("Item 1")).
//		Add(layout.NewBox("Item 2"))
func (f *Flex) Add(box *Box) *Flex {
	node := model2.NewNode(box.domain)
	f.domain = f.domain.AddItem(node)
	return f
}

// AddRaw adds a raw string as a box to the container.
// This is a convenience method equivalent to Add(NewBox(content)).
//
// Example:
//
//	flex := layout.Row().
//		AddRaw("Item 1").
//		AddRaw("Item 2")
func (f *Flex) AddRaw(content string) *Flex {
	return f.Add(NewBox(content))
}

// ============================================================================
// Gap Spacing
// ============================================================================

// Gap sets the spacing between items (in cells).
//
// Example:
//
//	flex := layout.Row().Gap(3)  // 3 cells between each item
func (f *Flex) Gap(gap int) *Flex {
	f.domain = f.domain.WithGap(gap)
	return f
}

// ============================================================================
// Justify Content (Main Axis Distribution)
// ============================================================================

// JustifyStart packs items at the start of the container (default).
//
// Visual (Row):
//
//	[1][2][3]         (remaining space)
//
// Example:
//
//	flex := layout.Row().JustifyStart()
func (f *Flex) JustifyStart() *Flex {
	f.domain = f.domain.WithJustifyContent(value2.JustifyContentStart)
	return f
}

// JustifyEnd packs items at the end of the container.
//
// Visual (Row):
//
//	(remaining space)         [1][2][3]
//
// Example:
//
//	flex := layout.Row().JustifyEnd()
func (f *Flex) JustifyEnd() *Flex {
	f.domain = f.domain.WithJustifyContent(value2.JustifyContentEnd)
	return f
}

// JustifyCenter centers items in the container.
//
// Visual (Row):
//
//	(space)    [1][2][3]    (space)
//
// Example:
//
//	flex := layout.Row().JustifyCenter()
func (f *Flex) JustifyCenter() *Flex {
	f.domain = f.domain.WithJustifyContent(value2.JustifyContentCenter)
	return f
}

// JustifySpaceBetween distributes items with equal spacing.
// First item at start, last item at end, equal gaps between.
//
// Visual (Row):
//
//	[1]    (gap)    [2]    (gap)    [3]
//
// Example:
//
//	flex := layout.Row().JustifySpaceBetween()
func (f *Flex) JustifySpaceBetween() *Flex {
	f.domain = f.domain.WithJustifyContent(value2.JustifyContentSpaceBetween)
	return f
}

// ============================================================================
// Align Items (Cross Axis Alignment)
// ============================================================================

// AlignStretch stretches items to fill the cross axis (default).
//
// Visual (Row):
//
//	┌───┐ ┌───┐ ┌───┐
//	│ 1 │ │ 2 │ │ 3 │  ← All same height
//	└───┘ └───┘ └───┘
//
// Example:
//
//	flex := layout.Row().AlignStretch()
func (f *Flex) AlignStretch() *Flex {
	f.domain = f.domain.WithAlignItems(value2.AlignItemsStretch)
	return f
}

// AlignStart aligns items at the start of the cross axis.
//
// Visual (Row):
//
//	┌───┐ ┌───┐ ┌───┐
//	│ 1 │ │ 2 │ │ 3 │  ← Aligned to top
//	│   │ └───┘ │   │
//	└───┘       └───┘
//
// Example:
//
//	flex := layout.Row().AlignStart()
func (f *Flex) AlignStart() *Flex {
	f.domain = f.domain.WithAlignItems(value2.AlignItemsStart)
	return f
}

// AlignEnd aligns items at the end of the cross axis.
//
// Visual (Row):
//
//	┌───┐       ┌───┐
//	│ 1 │ ┌───┐ │ 3 │
//	│   │ │ 2 │ │   │  ← Aligned to bottom
//	└───┘ └───┘ └───┘
//
// Example:
//
//	flex := layout.Row().AlignEnd()
func (f *Flex) AlignEnd() *Flex {
	f.domain = f.domain.WithAlignItems(value2.AlignItemsEnd)
	return f
}

// AlignCenter centers items along the cross axis.
//
// Visual (Row):
//
//	┌───┐
//	│ 1 │ ┌───┐ ┌───┐
//	│   │ │ 2 │ │ 3 │  ← Centered vertically
//	└───┘ └───┘ │   │
//	            └───┘
//
// Example:
//
//	flex := layout.Row().AlignCenter()
func (f *Flex) AlignCenter() *Flex {
	f.domain = f.domain.WithAlignItems(value2.AlignItemsCenter)
	return f
}

// ============================================================================
// Size Constraints
// ============================================================================

// Width sets an exact width constraint for the container.
//
// Example:
//
//	flex := layout.Row().Width(80)
func (f *Flex) Width(width int) *Flex {
	f.domain = f.domain.WithSize(f.domain.Size().WithWidth(width))
	return f
}

// Height sets an exact height constraint for the container.
//
// Example:
//
//	flex := layout.Column().Height(24)
func (f *Flex) Height(height int) *Flex {
	f.domain = f.domain.WithSize(f.domain.Size().WithHeight(height))
	return f
}

// ============================================================================
// Rendering
// ============================================================================

// Render lays out and renders the flexbox container.
//
// Parameters:
//   - containerWidth: Available width for the container (in cells)
//   - containerHeight: Available height for the container (in lines)
//
// Returns:
//   - Rendered string output with all items positioned and rendered
//
// Example:
//
//	output := layout.Row().
//		Add(layout.NewBox("Left")).
//		Add(layout.NewBox("Right")).
//		Render(80, 24)
//	fmt.Println(output)
//
//nolint:gocognit // Flex rendering orchestrates multiple services
func (f *Flex) Render(containerWidth, containerHeight int) string {
	// Create services
	measureService := service2.NewMeasureService()
	flexService := service2.NewFlexboxLayoutService(measureService)
	renderService := service2.NewRenderService()

	// Layout the flexbox
	laidOut := flexService.Layout(f.domain, containerWidth, containerHeight)

	// Render each item
	var result strings.Builder

	// Create a grid to track occupied positions
	grid := make([][]rune, containerHeight)
	for i := range grid {
		grid[i] = make([]rune, containerWidth)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Render each item into the grid
	for _, item := range laidOut.Items() {
		itemOutput := renderService.Render(item.Box())
		itemLines := strings.Split(itemOutput, "\n")

		pos := item.Position()
		for lineIdx, line := range itemLines {
			y := pos.Y() + lineIdx
			if y >= containerHeight {
				break
			}

			x := pos.X()
			for _, ch := range line {
				if x >= containerWidth {
					break
				}
				grid[y][x] = ch
				x++
			}
		}
	}

	// Convert grid to string
	for i, row := range grid {
		result.WriteString(string(row))
		if i < len(grid)-1 {
			result.WriteRune('\n')
		}
	}

	return result.String()
}

// String implements fmt.Stringer.
// Renders the flexbox with default terminal size (80x24).
func (f *Flex) String() string {
	return f.Render(80, 24)
}

// ============================================================================
// Advanced: Direct access to domain model
// ============================================================================

// Domain returns the underlying domain model.
// This is provided for advanced use cases.
func (f *Flex) Domain() *model2.FlexContainer {
	return f.domain
}
