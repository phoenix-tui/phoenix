package service

import (
	"strings"

	"github.com/phoenix-tui/phoenix/core/domain/service"
	"github.com/phoenix-tui/phoenix/style/domain/value"
)

// SpacingCalculator is a domain service that handles spacing calculations for content.
// This includes padding and margin application, as well as total dimension calculations.
// This is pure business logic with no infrastructure dependencies.
type SpacingCalculator interface {
	// CalculateTotalWidth calculates total width including content, padding, and margin.
	CalculateTotalWidth(contentWidth int, padding value.Padding, margin value.Margin) int

	// CalculateTotalHeight calculates total height including content, padding, and margin.
	CalculateTotalHeight(contentHeight int, padding value.Padding, margin value.Margin) int

	// ApplyPadding adds padding around content.
	// Padding adds space INSIDE the content box.
	ApplyPadding(content string, padding value.Padding) string

	// ApplyMargin adds margin around content.
	// Margin adds space OUTSIDE the content box.
	ApplyMargin(content string, margin value.Margin) string

	// ApplyBoth applies both padding and margin (padding first, then margin).
	ApplyBoth(content string, padding value.Padding, margin value.Margin) string
}

// DefaultSpacingCalculator is the default implementation of SpacingCalculator.
type DefaultSpacingCalculator struct {
	unicodeService *service.UnicodeService
}

// NewSpacingCalculator creates a new DefaultSpacingCalculator.
func NewSpacingCalculator(unicodeService *service.UnicodeService) SpacingCalculator {
	return &DefaultSpacingCalculator{
		unicodeService: unicodeService,
	}
}

// CalculateTotalWidth calculates the total width including content, padding, and margin.
// Formula: contentWidth + padding.Left + padding.Right + margin.Left + margin.Right.
func (sc *DefaultSpacingCalculator) CalculateTotalWidth(contentWidth int, padding value.Padding, margin value.Margin) int {
	return contentWidth + padding.Horizontal() + margin.Horizontal()
}

// CalculateTotalHeight calculates the total height including content, padding, and margin.
// Formula: contentHeight + padding.Top + padding.Bottom + margin.Top + margin.Bottom.
func (sc *DefaultSpacingCalculator) CalculateTotalHeight(contentHeight int, padding value.Padding, margin value.Margin) int {
	return contentHeight + padding.Vertical() + margin.Vertical()
}

// ApplyPadding adds padding around content.
// - Top/Bottom padding: Adds empty lines.
// - Left/Right padding: Adds spaces to each line.
//
//nolint:dupl // ApplyPadding and ApplyMargin are similar but semantically different (padding inside, margin outside)
func (sc *DefaultSpacingCalculator) ApplyPadding(content string, padding value.Padding) string {
	if content == "" {
		return ""
	}

	lines := strings.Split(content, "\n")

	// Apply horizontal padding (left and right spaces).
	paddedLines := make([]string, 0, len(lines))
	leftPad := strings.Repeat(" ", padding.Left())
	rightPad := strings.Repeat(" ", padding.Right())

	for _, line := range lines {
		paddedLine := leftPad + line + rightPad
		paddedLines = append(paddedLines, paddedLine)
	}

	// Apply vertical padding (top and bottom empty lines).
	topPadding := make([]string, padding.Top())
	bottomPadding := make([]string, padding.Bottom())

	// Calculate width of padded lines for empty padding lines.
	// Use the first line's width if available.
	emptyLineWidth := padding.Left() + padding.Right()
	if len(paddedLines) > 0 {
		emptyLineWidth = sc.unicodeService.StringWidth(paddedLines[0])
	}
	emptyLine := strings.Repeat(" ", emptyLineWidth)

	for i := range topPadding {
		topPadding[i] = emptyLine
	}
	for i := range bottomPadding {
		bottomPadding[i] = emptyLine
	}

	// Combine all parts.
	//nolint:gocritic,makezero // appendAssign: Pattern is clear; makezero: topPadding slice size is known
	result := append(topPadding, paddedLines...)
	result = append(result, bottomPadding...)

	return strings.Join(result, "\n")
}

// ApplyMargin adds margin around content.
// - Top/Bottom margin: Adds empty lines.
// - Left/Right margin: Adds spaces to each line (outside content box).
//
//nolint:dupl // ApplyMargin and ApplyPadding are similar but semantically different (margin outside, padding inside)
func (sc *DefaultSpacingCalculator) ApplyMargin(content string, margin value.Margin) string {
	if content == "" {
		return ""
	}

	lines := strings.Split(content, "\n")

	// Apply horizontal margin (left and right spaces).
	marginedLines := make([]string, 0, len(lines))
	leftMargin := strings.Repeat(" ", margin.Left())
	rightMargin := strings.Repeat(" ", margin.Right())

	for _, line := range lines {
		marginedLine := leftMargin + line + rightMargin
		marginedLines = append(marginedLines, marginedLine)
	}

	// Apply vertical margin (top and bottom empty lines).
	topMargin := make([]string, margin.Top())
	bottomMargin := make([]string, margin.Bottom())

	// Calculate width of margined lines for empty margin lines.
	emptyLineWidth := margin.Left() + margin.Right()
	if len(marginedLines) > 0 {
		emptyLineWidth = sc.unicodeService.StringWidth(marginedLines[0])
	}
	emptyLine := strings.Repeat(" ", emptyLineWidth)

	for i := range topMargin {
		topMargin[i] = emptyLine
	}
	for i := range bottomMargin {
		bottomMargin[i] = emptyLine
	}

	// Combine all parts.
	//nolint:gocritic,makezero // appendAssign: Pattern is clear; makezero: topMargin slice size is known
	result := append(topMargin, marginedLines...)
	result = append(result, bottomMargin...)

	return strings.Join(result, "\n")
}

// ApplyBoth applies both padding and margin (padding first, then margin).
// This is equivalent to ApplyMargin(ApplyPadding(content, padding), margin).
func (sc *DefaultSpacingCalculator) ApplyBoth(content string, padding value.Padding, margin value.Margin) string {
	padded := sc.ApplyPadding(content, padding)
	return sc.ApplyMargin(padded, margin)
}
