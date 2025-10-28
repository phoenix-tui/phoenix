// Package service provides rendering services for progress domain.
package service

import (
	"strings"

	"github.com/phoenix-tui/phoenix/components/progress/internal/domain/model"
)

// RenderService handles progress bar rendering logic.
// It provides pure domain logic for converting progress state to visual representation.
type RenderService struct{}

// NewRenderService creates a new RenderService.
func NewRenderService() *RenderService {
	return &RenderService{}
}

// RenderBar renders a progress bar to a string.
// Format: [label] [filled][empty] [percentage].
// Example: "Downloading... ████████░░░░░░░░ 40%".
func (s *RenderService) RenderBar(bar *model.Bar) string {
	if bar == nil {
		return ""
	}

	var parts []string

	// Add label if present.
	if bar.Label() != "" {
		parts = append(parts, bar.Label())
	}

	// Calculate filled and empty widths.
	filledWidth := s.CalculateFilledWidth(bar.Width(), bar.Percentage())
	emptyWidth := bar.Width() - filledWidth

	// Build bar string.
	barStr := strings.Repeat(string(bar.FillChar()), filledWidth) +
		strings.Repeat(string(bar.EmptyChar()), emptyWidth)
	parts = append(parts, barStr)

	// Add percentage if enabled.
	if bar.ShowPercent() {
		parts = append(parts, s.formatPercentage(bar.Percentage()))
	}

	return strings.Join(parts, " ")
}

// CalculateFilledWidth calculates the number of characters to fill based on percentage.
// Returns a value in [0, barWidth].
func (s *RenderService) CalculateFilledWidth(barWidth, percentage int) int {
	if barWidth <= 0 {
		return 0
	}
	if percentage <= 0 {
		return 0
	}
	if percentage >= 100 {
		return barWidth
	}

	// Calculate filled width: (percentage / 100) * barWidth.
	// Use integer arithmetic to avoid floating point.
	filled := (percentage * barWidth) / 100
	return filled
}

// formatPercentage formats a percentage value as a string.
func (s *RenderService) formatPercentage(pct int) string {
	// Simple integer formatting.
	return string(rune('0'+pct/100%10)) +
		string(rune('0'+pct/10%10)) +
		string(rune('0'+pct%10)) + "%"
}
