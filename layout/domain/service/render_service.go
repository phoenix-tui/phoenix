package service

import (
	"strings"

	"github.com/phoenix-tui/phoenix/layout/domain/model"
)

// RenderService converts positioned box trees into final text output.
// Produces multi-line strings with borders, padding, and margin.
//
// Design Philosophy:
//   - Domain service (pure business logic)
//   - ASCII box drawing characters (Unicode support future)
//   - Follows CSS box model rendering order
//   - Returns multi-line strings (joined with \n)
//
// Rendering Process (inside-out):
//  1. Content lines (text)
//  2. Padding (spaces around content)
//  3. Border (box drawing characters)
//  4. Margin (empty lines/spaces)
//
// Border Characters:
//
//	┌──────┐  Top border
//	│ Text │  Side borders with content
//	└──────┘  Bottom border
//
// Example:
//
//	rs := NewRenderService()
//	box := model.NewBox("Hello").
//		WithPadding(value.NewSpacingAll(1)).
//		WithBorder(true)
//
//	output := rs.Render(box)
//	// Output:
//	// ┌───────┐
//	// │       │
//	// │ Hello │
//	// │       │
//	// └───────┘
type RenderService struct{}

// NewRenderService creates a new RenderService.
func NewRenderService() *RenderService {
	return &RenderService{}
}

// Render produces the final text output for a box.
// Returns a multi-line string with borders, padding, and margin.
//
// Algorithm:
//  1. Add margin top (empty lines)
//  2. Add top border (if enabled)
//  3. Add padding top (empty lines with side borders)
//  4. Add content lines (with padding and borders)
//  5. Add padding bottom (empty lines with side borders)
//  6. Add bottom border (if enabled)
//  7. Add margin bottom (empty lines)
//
// Example:
//
//	box := model.NewBox("Hi").WithBorder(true)
//	output := rs.Render(box)
//	// ┌────┐
//	// │ Hi │
//	// └────┘
//
// Returns:
//   - Multi-line string (lines joined with \n)
func (rs *RenderService) Render(box *model.Box) string {
	var lines []string

	// Get box properties
	content := box.Content()
	padding := box.Padding()
	margin := box.Margin()
	hasBorder := box.HasBorder()

	// Split content into lines
	contentLines := strings.Split(content, "\n")

	// Calculate widths
	contentWidth := rs.calculateContentWidth(contentLines)

	// Calculate total padding (explicit + implicit for borders)
	// When border is enabled, add 1-space aesthetic padding between border and content
	// This is added to any explicit padding (│  Hi  │ with padding=1 gives 2 spaces)
	totalPaddingLeft := padding.Left()
	totalPaddingRight := padding.Right()
	if hasBorder {
		totalPaddingLeft += 1 // Aesthetic spacing between border and content
		totalPaddingRight += 1
	}

	// innerWidth includes content + total padding
	innerWidth := contentWidth + totalPaddingLeft + totalPaddingRight

	totalWidth := innerWidth
	if hasBorder {
		totalWidth += 2 // Border characters
	}

	// Step 1: Margin top
	for i := 0; i < margin.Top(); i++ {
		lines = append(lines, strings.Repeat(" ", totalWidth+margin.Horizontal()))
	}

	// Step 2: Top border
	if hasBorder {
		borderLine := rs.renderMarginLeft(margin) + "┌" + strings.Repeat("─", innerWidth) + "┐" + rs.renderMarginRight(margin)
		lines = append(lines, borderLine)
	}

	// Step 3: Padding top
	for i := 0; i < padding.Top(); i++ {
		line := rs.renderMarginLeft(margin)
		if hasBorder {
			line += "│"
		}
		line += strings.Repeat(" ", innerWidth)
		if hasBorder {
			line += "│"
		}
		line += rs.renderMarginRight(margin)
		lines = append(lines, line)
	}

	// Step 4: Content lines
	for _, contentLine := range contentLines {
		line := rs.renderMarginLeft(margin)
		if hasBorder {
			line += "│"
		}
		// Add left padding (total = explicit + implicit for borders)
		line += strings.Repeat(" ", totalPaddingLeft)

		// Add content
		line += contentLine

		// Add right padding to align right border
		lineWidth := len(contentLine)
		spacesNeeded := contentWidth - lineWidth

		// Only pad to contentWidth if we have a border (to align it)
		// Without border, no alignment padding needed
		if hasBorder {
			line += strings.Repeat(" ", spacesNeeded) + strings.Repeat(" ", totalPaddingRight)
		} else {
			// Without border, just add explicit padding (if any)
			line += strings.Repeat(" ", totalPaddingRight)
		}

		if hasBorder {
			line += "│"
		}
		line += rs.renderMarginRight(margin)
		lines = append(lines, line)
	}

	// Step 5: Padding bottom
	for i := 0; i < padding.Bottom(); i++ {
		line := rs.renderMarginLeft(margin)
		if hasBorder {
			line += "│"
		}
		line += strings.Repeat(" ", innerWidth)
		if hasBorder {
			line += "│"
		}
		line += rs.renderMarginRight(margin)
		lines = append(lines, line)
	}

	// Step 6: Bottom border
	if hasBorder {
		borderLine := rs.renderMarginLeft(margin) + "└" + strings.Repeat("─", innerWidth) + "┘" + rs.renderMarginRight(margin)
		lines = append(lines, borderLine)
	}

	// Step 7: Margin bottom
	for i := 0; i < margin.Bottom(); i++ {
		lines = append(lines, strings.Repeat(" ", totalWidth+margin.Horizontal()))
	}

	return strings.Join(lines, "\n")
}

// RenderNode produces output for entire node tree.
// For Day 3, this is simplified to render root box only.
// Full tree rendering with absolute positioning will come later.
//
// Parameters:
//   - node: Root node to render
//
// Returns:
//   - Multi-line string of rendered output
func (rs *RenderService) RenderNode(node *model.Node) string {
	// For now, just render the root box
	// Full tree rendering with child positioning comes in later iterations
	return rs.Render(node.Box())
}

// calculateContentWidth finds the maximum line width in content.
// This is a simple character count (Unicode width handled by caller if needed).
func (rs *RenderService) calculateContentWidth(lines []string) int {
	maxWidth := 0
	for _, line := range lines {
		width := len(line)
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}

// renderMarginLeft renders left margin spaces.
func (rs *RenderService) renderMarginLeft(margin interface{ Left() int }) string {
	return strings.Repeat(" ", margin.Left())
}

// renderMarginRight renders right margin spaces.
func (rs *RenderService) renderMarginRight(margin interface{ Right() int }) string {
	return strings.Repeat(" ", margin.Right())
}

// renderPaddingLeft renders left padding spaces.
func (rs *RenderService) renderPaddingLeft(padding interface{ Left() int }) string {
	return strings.Repeat(" ", padding.Left())
}

// renderPaddingRight renders right padding spaces.
// Pads to contentWidth to align right border.
func (rs *RenderService) renderPaddingRight(padding interface{ Right() int }, contentLine string, contentWidth int) string {
	// Calculate spaces needed to reach contentWidth
	lineWidth := len(contentLine)
	spacesNeeded := contentWidth - lineWidth
	return strings.Repeat(" ", spacesNeeded) + strings.Repeat(" ", padding.Right())
}
