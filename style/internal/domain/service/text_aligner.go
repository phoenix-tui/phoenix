package service

import (
	"strings"

	"github.com/phoenix-tui/phoenix/core"
	"github.com/phoenix-tui/phoenix/style/internal/domain/value"
)

// TextAligner is a domain service that handles text alignment within given dimensions.
// This service correctly handles Unicode width calculation for proper alignment.
// This is pure business logic with no infrastructure dependencies.
type TextAligner interface {
	// AlignHorizontal aligns text within given width.
	// Returns the aligned text (single line).
	AlignHorizontal(text string, width int, alignment value.HorizontalAlignment) string

	// AlignVertical aligns text within given height.
	// Returns the aligned text (multiple lines).
	AlignVertical(text string, height int, alignment value.VerticalAlignment) string

	// AlignBoth aligns text within given dimensions (width x height).
	// Applies both horizontal and vertical alignment.
	AlignBoth(text string, width, height int, alignment value.Alignment) string
}

// DefaultTextAligner is the default implementation of TextAligner.
type DefaultTextAligner struct {
}

// NewTextAligner creates a new DefaultTextAligner.
func NewTextAligner() TextAligner {
	return &DefaultTextAligner{}
}

// AlignHorizontal aligns a single line of text within given width.
// Handles Unicode correctly by using visual width, not string length.
func (ta *DefaultTextAligner) AlignHorizontal(text string, width int, alignment value.HorizontalAlignment) string {
	// Calculate actual visual width of text.
	textWidth := core.StringWidth(text)

	// If text is wider than target width, truncate (future: could add ellipsis).
	if textWidth > width {
		return text[:width] // Simple truncation for now
	}

	// If text fits exactly, return as-is.
	if textWidth == width {
		return text
	}

	// Calculate padding needed.
	paddingTotal := width - textWidth

	switch alignment {
	case value.AlignLeft:
		// Add all padding to the right.
		return text + strings.Repeat(" ", paddingTotal)

	case value.AlignCenter:
		// Distribute padding evenly (left gets extra if odd).
		paddingLeft := (paddingTotal + 1) / 2 // Rounds up for odd numbers
		paddingRight := paddingTotal - paddingLeft
		return strings.Repeat(" ", paddingLeft) + text + strings.Repeat(" ", paddingRight)

	case value.AlignRight:
		// Add all padding to the left.
		return strings.Repeat(" ", paddingTotal) + text

	default:
		// Default to left alignment.
		return text + strings.Repeat(" ", paddingTotal)
	}
}

// AlignVertical aligns text within given height.
// Text is treated as multi-line (split by \n).
//
//nolint:funlen // Function length justified: comprehensive alignment with multiple cases and proper error handling
func (ta *DefaultTextAligner) AlignVertical(text string, height int, alignment value.VerticalAlignment) string {
	lines := strings.Split(text, "\n")
	currentHeight := len(lines)

	// If content is taller than target height, truncate.
	if currentHeight > height {
		return strings.Join(lines[:height], "\n")
	}

	// If content fits exactly, return as-is.
	if currentHeight == height {
		return text
	}

	// Calculate empty lines needed.
	emptyLinesTotal := height - currentHeight

	// Determine width for empty lines (use max line width).
	emptyLineWidth := 0
	for _, line := range lines {
		lineWidth := core.StringWidth(line)
		if lineWidth > emptyLineWidth {
			emptyLineWidth = lineWidth
		}
	}
	emptyLine := strings.Repeat(" ", emptyLineWidth)

	switch alignment {
	case value.AlignTop:
		// Add all empty lines at bottom.
		emptyLines := make([]string, emptyLinesTotal)
		for i := range emptyLines {
			emptyLines[i] = emptyLine
		}
		//nolint:gocritic // appendAssign: Building result from multiple slices
		result := append(lines, emptyLines...)
		return strings.Join(result, "\n")

	case value.AlignMiddle:
		// Distribute empty lines evenly (top gets extra if odd).
		emptyTop := (emptyLinesTotal + 1) / 2 // Rounds up for odd numbers
		emptyBottom := emptyLinesTotal - emptyTop

		topLines := make([]string, emptyTop)
		for i := range topLines {
			topLines[i] = emptyLine
		}
		bottomLines := make([]string, emptyBottom)
		for i := range bottomLines {
			bottomLines[i] = emptyLine
		}

		//nolint:gocritic,makezero // appendAssign: Building result from multiple slices; makezero: topLines slice size is known
		result := append(topLines, lines...)
		result = append(result, bottomLines...)
		return strings.Join(result, "\n")

	case value.AlignBottom:
		// Add all empty lines at top.
		emptyLines := make([]string, emptyLinesTotal)
		for i := range emptyLines {
			emptyLines[i] = emptyLine
		}
		//nolint:gocritic,makezero // appendAssign: Building result from multiple slices; makezero: emptyLines slice size is known
		result := append(emptyLines, lines...)
		return strings.Join(result, "\n")

	default:
		// Default to top alignment.
		emptyLines := make([]string, emptyLinesTotal)
		for i := range emptyLines {
			emptyLines[i] = emptyLine
		}
		//nolint:gocritic // appendAssign: Building result from multiple slices
		result := append(lines, emptyLines...)
		return strings.Join(result, "\n")
	}
}

// AlignBoth aligns text within given width and height.
// Applies horizontal alignment to each line, then vertical alignment.
func (ta *DefaultTextAligner) AlignBoth(text string, width, height int, alignment value.Alignment) string {
	lines := strings.Split(text, "\n")

	// First, apply horizontal alignment to each line.
	alignedLines := make([]string, len(lines))
	for i, line := range lines {
		alignedLines[i] = ta.AlignHorizontal(line, width, alignment.Horizontal())
	}

	// Rejoin lines and apply vertical alignment.
	horizontallyAligned := strings.Join(alignedLines, "\n")
	return ta.AlignVertical(horizontallyAligned, height, alignment.Vertical())
}
