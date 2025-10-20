package command

import (
	"fmt"
	"strings"

	coreService "github.com/phoenix-tui/phoenix/core/domain/service"
	"github.com/phoenix-tui/phoenix/style/domain/model"
	"github.com/phoenix-tui/phoenix/style/domain/service"
	"github.com/phoenix-tui/phoenix/style/domain/value"
	"github.com/phoenix-tui/phoenix/style/infrastructure/ansi"
)

// RenderCommand applies a Style to content and generates ANSI-styled output.
// This is the main application service that orchestrates all domain services.
// to execute the complete styling pipeline.
//
// Pipeline:.
//  1. Style validation.
//  2. Size validation (if size constraints set).
//  3. Text alignment (if alignment set).
//  4. Apply padding (if padding set).
//  5. Apply border (if border set).
//  6. Apply margin (if margin set).
//  7. Color adaptation & ANSI generation.
//  8. Text decorations (bold, italic, etc.).
//
// Example:.
//
//	cmd := NewRenderCommand(colorAdapter, spacingCalc, textAligner, unicodeService, ansiGen).
//	style := model.NewStyle().Foreground(value.RGB(255, 0, 0)).Bold(true).
//	output, err := cmd.Execute(style, "Hello, World!").
type RenderCommand struct {
	colorAdapter      service.ColorAdapter
	spacingCalculator service.SpacingCalculator
	textAligner       service.TextAligner
	unicodeService    *coreService.UnicodeService
	ansiGenerator     *ansi.ANSICodeGenerator
	borderRenderer    *BorderRenderer
}

// NewRenderCommand creates a new RenderCommand.
func NewRenderCommand(
	colorAdapter service.ColorAdapter,
	spacingCalculator service.SpacingCalculator,
	textAligner service.TextAligner,
	unicodeService *coreService.UnicodeService,
	ansiGenerator *ansi.ANSICodeGenerator,
) *RenderCommand {
	return &RenderCommand{
		colorAdapter:      colorAdapter,
		spacingCalculator: spacingCalculator,
		textAligner:       textAligner,
		unicodeService:    unicodeService,
		ansiGenerator:     ansiGenerator,
		borderRenderer:    NewBorderRenderer(unicodeService, ansiGenerator),
	}
}

// Execute applies the style to content and returns ANSI-styled string.
//
//nolint:gocognit,gocyclo,cyclop // Complexity justified: comprehensive style application with multiple optional properties
func (rc *RenderCommand) Execute(style model.Style, content string) (string, error) {
	// 1. Validate style.
	if err := style.Validate(); err != nil {
		return "", fmt.Errorf("style validation failed: %w", err)
	}

	// 2. Size validation & content preparation.
	targetWidth := 0
	targetHeight := 0

	//nolint:nestif // Nested checks required for optional style properties with interdependencies
	if size, ok := style.GetSize(); ok {
		// Calculate effective dimensions accounting for spacing.
		contentWidth, contentHeight := rc.calculateContentDimensions(content)

		// Add padding dimensions.
		if padding, hasPadding := style.GetPadding(); hasPadding {
			contentWidth += padding.Left() + padding.Right()
			contentHeight += padding.Top() + padding.Bottom()
		}

		// Add border dimensions.
		if style.GetBorderLeft() || style.GetBorderRight() {
			contentWidth += 2 // Left + right border
		}
		if style.GetBorderTop() || style.GetBorderBottom() {
			contentHeight += 2 // Top + bottom border
		}

		// Add margin dimensions.
		if margin, hasMargin := style.GetMargin(); hasMargin {
			contentWidth += margin.Left() + margin.Right()
			contentHeight += margin.Top() + margin.Bottom()
		}

		// Validate against constraints.
		if err := rc.validateSizeConstraints(size, contentWidth, contentHeight); err != nil {
			return "", err
		}

		// Use size constraints as target dimensions.
		if width, hasWidth := size.Width(); hasWidth {
			targetWidth = width
		}
		if height, hasHeight := size.Height(); hasHeight {
			targetHeight = height
		}
	}

	// 3. Text alignment (before padding/border).
	//nolint:gocritic,nestif // ifElseChain: Sequential checks are clearer; nestif: Optional properties require nested checks
	if alignment, ok := style.GetAlignment(); ok {
		if targetWidth > 0 && targetHeight > 0 {
			content = rc.textAligner.AlignBoth(content, targetWidth, targetHeight, alignment)
		} else if targetWidth > 0 {
			content = rc.textAligner.AlignHorizontal(content, targetWidth, alignment.Horizontal())
		} else if targetHeight > 0 {
			content = rc.textAligner.AlignVertical(content, targetHeight, alignment.Vertical())
		}
	}

	// 4. Apply padding.
	if padding, ok := style.GetPadding(); ok {
		content = rc.spacingCalculator.ApplyPadding(content, padding)
	}

	// 5. Apply border.
	if _, hasBorder := style.GetBorder(); hasBorder {
		content = rc.borderRenderer.Render(content, style)
	}

	// 6. Apply margin.
	if margin, ok := style.GetMargin(); ok {
		content = rc.spacingCalculator.ApplyMargin(content, margin)
	}

	// 7. Color adaptation & ANSI generation.
	content = rc.applyColors(content, style)

	// 8. Text decorations.
	content = rc.applyDecorations(content, style)

	return content, nil
}

// calculateContentDimensions calculates the width and height of content.
func (rc *RenderCommand) calculateContentDimensions(content string) (int, int) {
	lines := strings.Split(content, "\n")
	height := len(lines)

	maxWidth := 0
	for _, line := range lines {
		width := rc.unicodeService.StringWidth(line)
		if width > maxWidth {
			maxWidth = width
		}
	}

	return maxWidth, height
}

// validateSizeConstraints validates content dimensions against min/max constraints only.
// Width/Height are TARGET dimensions, not constraints - they're used for alignment.
func (rc *RenderCommand) validateSizeConstraints(size value.Size, width, height int) error {
	// Check minimum width.
	if minWidth, hasMin := size.MinWidth(); hasMin {
		if width < minWidth {
			return fmt.Errorf("content width %d is less than minimum width %d", width, minWidth)
		}
	}

	// Check maximum width.
	if maxWidth, hasMax := size.MaxWidth(); hasMax {
		if width > maxWidth {
			return fmt.Errorf("content width %d exceeds maximum width %d", width, maxWidth)
		}
	}

	// Check minimum height.
	if minHeight, hasMin := size.MinHeight(); hasMin {
		if height < minHeight {
			return fmt.Errorf("content height %d is less than minimum height %d", height, minHeight)
		}
	}

	// Check maximum height.
	if maxHeight, hasMax := size.MaxHeight(); hasMax {
		if height > maxHeight {
			return fmt.Errorf("content height %d exceeds maximum height %d", height, maxHeight)
		}
	}

	return nil
}

// applyColors applies foreground and background colors to content.
func (rc *RenderCommand) applyColors(content string, style model.Style) string {
	termCap := style.GetTerminalCapability()
	var codes []string

	// Foreground color.
	if fg, hasFg := style.GetForeground(); hasFg {
		ansiCode := rc.colorAdapter.ToANSIForeground(fg, termCap)
		if ansiCode != "" {
			codes = append(codes, ansiCode)
		}
	}

	// Background color.
	if bg, hasBg := style.GetBackground(); hasBg {
		ansiCode := rc.colorAdapter.ToANSIBackground(bg, termCap)
		if ansiCode != "" {
			codes = append(codes, ansiCode)
		}
	}

	// No colors to apply.
	if len(codes) == 0 {
		return content
	}

	// Apply colors to entire content.
	prefix := strings.Join(codes, "")
	suffix := rc.ansiGenerator.Reset()

	return prefix + content + suffix
}

// applyDecorations applies text decorations (bold, italic, underline, strikethrough).
func (rc *RenderCommand) applyDecorations(content string, style model.Style) string {
	var codes []string

	if style.GetBold() {
		codes = append(codes, rc.ansiGenerator.Bold())
	}
	if style.GetItalic() {
		codes = append(codes, rc.ansiGenerator.Italic())
	}
	if style.GetUnderline() {
		codes = append(codes, rc.ansiGenerator.Underline())
	}
	if style.GetStrikethrough() {
		codes = append(codes, rc.ansiGenerator.Strikethrough())
	}

	// No decorations to apply.
	if len(codes) == 0 {
		return content
	}

	// Apply decorations to entire content.
	prefix := strings.Join(codes, "")
	suffix := rc.ansiGenerator.Reset()

	return prefix + content + suffix
}
