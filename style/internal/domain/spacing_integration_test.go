package domain

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/core"
	service2 "github.com/phoenix-tui/phoenix/style/internal/domain/service"
	value2 "github.com/phoenix-tui/phoenix/style/internal/domain/value"
)

// --- Integration Tests for Spacing Components ---
// These tests verify that Padding, Margin, Size, Alignment, SpacingCalculator,
// and TextAligner work together correctly.

func TestPaddingAndMarginTogether(t *testing.T) {
	calc := service2.NewSpacingCalculator()

	padding := value2.UniformPadding(1)
	margin := value2.UniformMargin(2)

	content := "X"

	// Apply both padding and margin
	result := calc.ApplyBoth(content, padding, margin)

	// Verify structure:
	// - 2 margin lines (top)
	// - 1 padding line (top) + margin
	// - 1 content line with padding + margin
	// - 1 padding line (bottom) + margin
	// - 2 margin lines (bottom)
	// Total: 7 lines

	lines := strings.Split(result, "\n")
	expectedLines := 7 // 2 + 1 + 1 + 1 + 2

	if len(lines) != expectedLines {
		t.Errorf("Expected %d lines, got %d", expectedLines, len(lines))
	}

	// Middle line should have X with padding and margin
	middleLine := lines[3] // 0-indexed: line 3 is the 4th line
	if !strings.Contains(middleLine, "X") {
		t.Errorf("Middle line should contain X, got %q", middleLine)
	}
}

func TestSizeConstraintsWithPaddingAndMargin(t *testing.T) {
	calc := service2.NewSpacingCalculator()

	padding := value2.UniformPadding(2)
	margin := value2.UniformMargin(1)

	content := "Hello"
	contentWidth := core.StringWidth(content)

	// Calculate total width with padding and margin
	totalWidth := calc.CalculateTotalWidth(contentWidth, padding, margin)

	// Verify total width: content(5) + padding(2*2) + margin(1*2) = 11
	expectedWidth := 11
	if totalWidth != expectedWidth {
		t.Errorf("Total width = %d, want %d", totalWidth, expectedWidth)
	}

	// Verify size constraints can validate this width
	// Size.Width() sets an exact width constraint, but we're testing max constraint
	// The totalWidth (11) is larger than the exact width (10), so it should NOT exceed
	sizeWithMax := value2.NewSize().SetMaxWidth(10)
	validated := sizeWithMax.ValidateWidth(totalWidth)

	// ValidateWidth should clamp to max (10)
	if validated != 10 {
		t.Errorf("ValidateWidth should clamp to max: got %d, want 10", validated)
	}
}

func TestAlignmentWithSpacing(t *testing.T) {
	aligner := service2.NewTextAligner()
	calc := service2.NewSpacingCalculator()

	// Create content with alignment
	content := "Hi"
	width := 10
	alignment := value2.CenterMiddle()

	// First align the content
	aligned := aligner.AlignBoth(content, width, 3, alignment)

	// Then apply padding
	padding := value2.UniformPadding(1)
	padded := calc.ApplyPadding(aligned, padding)

	// Verify result is properly structured
	lines := strings.Split(padded, "\n")
	if len(lines) < 3 {
		t.Errorf("Expected at least 3 lines after padding, got %d", len(lines))
	}
}

func TestCompleteUseCase(t *testing.T) {
	// Complete use case: Create a styled box with:
	// - Content: "Hello"
	// - Centered horizontally and vertically
	// - Padding: 1
	// - Margin: 2
	// - Size constraints: min 10x5, max 20x10

	aligner := service2.NewTextAligner()
	calc := service2.NewSpacingCalculator()

	content := "Hello"
	size := value2.NewSize().SetMinWidth(10).SetMaxWidth(20).SetMinHeight(5).SetMaxHeight(10)
	padding := value2.UniformPadding(1)
	margin := value2.UniformMargin(2)
	alignment := value2.CenterMiddle()

	// Step 1: Determine target dimensions (respecting constraints)
	contentWidth := core.StringWidth(content)
	contentHeight := 1 // Single line

	// Add padding to get inner box size
	innerWidth := contentWidth + padding.Horizontal()
	innerHeight := contentHeight + padding.Vertical()

	// Validate against size constraints
	validatedWidth := size.ValidateWidth(innerWidth)
	validatedHeight := size.ValidateHeight(innerHeight)

	// Step 2: Align content
	aligned := aligner.AlignBoth(content, validatedWidth-padding.Horizontal(),
		validatedHeight-padding.Vertical(), alignment)

	// Step 3: Apply padding
	padded := calc.ApplyPadding(aligned, padding)

	// Step 4: Apply margin
	final := calc.ApplyMargin(padded, margin)

	// Verify final result
	lines := strings.Split(final, "\n")

	// Should have: margin.Top + padding.Top + content lines + padding.Bottom + margin.Bottom
	// At minimum: 2 + 1 + 1 + 1 + 2 = 7 lines
	if len(lines) < 7 {
		t.Errorf("Expected at least 7 lines, got %d", len(lines))
	}

	// Verify content is somewhere in the middle
	contentFound := false
	for _, line := range lines {
		if strings.Contains(line, "Hello") {
			contentFound = true
			break
		}
	}
	if !contentFound {
		t.Errorf("Content 'Hello' not found in result:\n%s", final)
	}
}

func TestUnicodeWithCompleteSpacing(t *testing.T) {
	// Test with Unicode content (emoji + CJK)
	aligner := service2.NewTextAligner()
	calc := service2.NewSpacingCalculator()

	content := "👋你好" // emoji(2) + CJK(4) = 6 visual width
	padding := value2.NewPadding(1, 1, 1, 1)
	margin := value2.NewMargin(1, 1, 1, 1)
	alignment := value2.CenterMiddle()

	// Align in 12-width box
	aligned := aligner.AlignHorizontal(content, 12, alignment.Horizontal())

	// Apply spacing
	final := calc.ApplyBoth(aligned, padding, margin)

	// Verify structure is correct
	lines := strings.Split(final, "\n")
	if len(lines) < 5 {
		t.Errorf("Expected at least 5 lines (margin+padding+content+padding+margin), got %d", len(lines))
	}

	// Verify content exists and is centered
	contentFound := false
	for _, line := range lines {
		if strings.Contains(line, "👋") && strings.Contains(line, "你好") {
			contentFound = true
			// Content should be centered with padding and margin
			break
		}
	}
	if !contentFound {
		t.Errorf("Unicode content not found in result:\n%s", final)
	}
}

func TestMultiLineContentWithAllFeatures(t *testing.T) {
	aligner := service2.NewTextAligner()
	calc := service2.NewSpacingCalculator()

	content := "Line1\nLine2\nLine3"
	padding := value2.UniformPadding(1)
	margin := value2.UniformMargin(1)
	alignment := value2.LeftTop()

	// Align multi-line content
	aligned := aligner.AlignBoth(content, 10, 5, alignment)

	// Apply spacing
	final := calc.ApplyBoth(aligned, padding, margin)

	// Verify all original lines are present
	if !strings.Contains(final, "Line1") {
		t.Errorf("Line1 not found in result")
	}
	if !strings.Contains(final, "Line2") {
		t.Errorf("Line2 not found in result")
	}
	if !strings.Contains(final, "Line3") {
		t.Errorf("Line3 not found in result")
	}

	// Verify structure has padding and margin
	lines := strings.Split(final, "\n")
	// Should have margin + padding + 5 aligned lines + padding + margin
	// At minimum: 1 + 1 + 5 + 1 + 1 = 9 lines
	if len(lines) < 9 {
		t.Errorf("Expected at least 9 lines, got %d", len(lines))
	}
}

func TestAsymmetricSpacingIntegration(t *testing.T) {
	calc := service2.NewSpacingCalculator()

	content := "X"
	padding := value2.NewPadding(1, 2, 3, 4) // top, right, bottom, left
	margin := value2.NewMargin(5, 6, 7, 8)   // top, right, bottom, left

	result := calc.ApplyBoth(content, padding, margin)

	// Verify total dimensions
	lines := strings.Split(result, "\n")

	// Height: margin.Top(5) + padding.Top(1) + content(1) + padding.Bottom(3) + margin.Bottom(7)
	expectedHeight := 17
	if len(lines) != expectedHeight {
		t.Errorf("Expected height %d, got %d", expectedHeight, len(lines))
	}

	// Find content line (should be at line index 5+1=6)
	contentLineIndex := margin.Top() + padding.Top()
	contentLine := lines[contentLineIndex]

	// Width: margin.Left(8) + padding.Left(4) + content(1) + padding.Right(2) + margin.Right(6)
	expectedWidth := 21
	actualWidth := core.StringWidth(contentLine)
	if actualWidth != expectedWidth {
		t.Errorf("Expected width %d, got %d", expectedWidth, actualWidth)
	}
}

func TestSizeValidationInPipeline(t *testing.T) {
	calc := service2.NewSpacingCalculator()

	content := "VeryLongContentThatExceedsMaxWidth"
	contentWidth := core.StringWidth(content)
	padding := value2.UniformPadding(2)
	margin := value2.UniformMargin(1)
	size := value2.NewSize().SetMaxWidth(20)

	// Calculate total width
	totalWidth := calc.CalculateTotalWidth(contentWidth, padding, margin)

	// Validate against size constraints
	validatedWidth := size.ValidateWidth(totalWidth)

	// Validated width should be clamped to max
	if maxWidth, isSet := size.MaxWidth(); isSet {
		if validatedWidth != maxWidth {
			t.Errorf("ValidateWidth should clamp to max: got %d, want %d", validatedWidth, maxWidth)
		}
	}
}

func TestEmptyContentWithSpacing(t *testing.T) {
	calc := service2.NewSpacingCalculator()

	padding := value2.UniformPadding(2)
	margin := value2.UniformMargin(1)

	result := calc.ApplyBoth("", padding, margin)

	// Empty content should return empty string
	if result != "" {
		t.Errorf("Empty content should remain empty after spacing, got %q", result)
	}
}
