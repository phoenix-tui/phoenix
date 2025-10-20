// Package command contains application layer command handlers for style rendering.
package command

import (
	"strings"

	"github.com/phoenix-tui/phoenix/core/domain/service"
	"github.com/phoenix-tui/phoenix/style/domain/model"
	"github.com/phoenix-tui/phoenix/style/domain/value"
	"github.com/phoenix-tui/phoenix/style/infrastructure/ansi"
)

// BorderRenderer handles border rendering around content.
// This is a helper for RenderCommand that encapsulates border logic.
type BorderRenderer struct {
	unicodeService *service.UnicodeService
	ansiGenerator  *ansi.ANSICodeGenerator
}

// NewBorderRenderer creates a new BorderRenderer.
func NewBorderRenderer(
	unicodeService *service.UnicodeService,
	ansiGenerator *ansi.ANSICodeGenerator,
) *BorderRenderer {
	return &BorderRenderer{
		unicodeService: unicodeService,
		ansiGenerator:  ansiGenerator,
	}
}

// Render applies border to content based on style.
// Returns the content with border applied.
func (br *BorderRenderer) Render(content string, style model.Style) string {
	border, hasBorder := style.GetBorder()
	if !hasBorder {
		return content
	}

	// Get border sides.
	hasTop := style.GetBorderTop()
	hasBottom := style.GetBorderBottom()
	hasLeft := style.GetBorderLeft()
	hasRight := style.GetBorderRight()

	// If no sides enabled, return content as-is.
	if !hasTop && !hasBottom && !hasLeft && !hasRight {
		return content
	}

	// Split content into lines.
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}

	// Calculate max line width (Unicode-correct).
	maxWidth := br.calculateMaxWidth(lines)

	// Build bordered content.
	result := []string{}

	// Top border.
	if hasTop {
		topLine := br.buildTopBorder(border, maxWidth, hasLeft, hasRight)
		topLine = br.applyBorderColor(topLine, style)
		result = append(result, topLine)
	}

	// Content lines with left/right borders.
	for _, line := range lines {
		borderedLine := br.buildContentLine(line, border, maxWidth, hasLeft, hasRight)
		borderedLine = br.applyBorderColorToSides(borderedLine, border, style, hasLeft, hasRight)
		result = append(result, borderedLine)
	}

	// Bottom border.
	if hasBottom {
		bottomLine := br.buildBottomBorder(border, maxWidth, hasLeft, hasRight)
		bottomLine = br.applyBorderColor(bottomLine, style)
		result = append(result, bottomLine)
	}

	return strings.Join(result, "\n")
}

// calculateMaxWidth calculates the maximum width of all lines (Unicode-correct).
func (br *BorderRenderer) calculateMaxWidth(lines []string) int {
	maxWidth := 0
	for _, line := range lines {
		width := br.unicodeService.StringWidth(line)
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}

// buildTopBorder builds the top border line.
func (br *BorderRenderer) buildTopBorder(border value.Border, width int, hasLeft, hasRight bool) string {
	var parts []string

	// Left corner.
	if hasLeft {
		parts = append(parts, border.TopLeft)
	}

	// Top edge.
	parts = append(parts, strings.Repeat(border.Top, width))

	// Right corner.
	if hasRight {
		parts = append(parts, border.TopRight)
	}

	return strings.Join(parts, "")
}

// buildBottomBorder builds the bottom border line.
func (br *BorderRenderer) buildBottomBorder(border value.Border, width int, hasLeft, hasRight bool) string {
	var parts []string

	// Left corner.
	if hasLeft {
		parts = append(parts, border.BottomLeft)
	}

	// Bottom edge.
	parts = append(parts, strings.Repeat(border.Bottom, width))

	// Right corner.
	if hasRight {
		parts = append(parts, border.BottomRight)
	}

	return strings.Join(parts, "")
}

// buildContentLine builds a content line with left/right borders.
func (br *BorderRenderer) buildContentLine(line string, border value.Border, targetWidth int, hasLeft, hasRight bool) string {
	// Calculate padding needed to reach target width.
	currentWidth := br.unicodeService.StringWidth(line)
	paddingWidth := targetWidth - currentWidth
	if paddingWidth < 0 {
		paddingWidth = 0
	}

	// Build line with padding.
	paddedLine := line + strings.Repeat(" ", paddingWidth)

	var parts []string

	// Left border.
	if hasLeft {
		parts = append(parts, border.Left)
	}

	// Content.
	parts = append(parts, paddedLine)

	// Right border.
	if hasRight {
		parts = append(parts, border.Right)
	}

	return strings.Join(parts, "")
}

// applyBorderColor applies border color to entire line.
func (br *BorderRenderer) applyBorderColor(line string, style model.Style) string {
	borderColor, hasBorderColor := style.GetBorderColor()
	if !hasBorderColor {
		return line
	}

	// Adapt color to terminal capability.
	termCap := style.GetTerminalCapability()
	ansiCode := br.colorToANSI(borderColor, termCap, true)

	return ansiCode + line + br.ansiGenerator.Reset()
}

// applyBorderColorToSides applies border color only to left/right border characters.
// This is more complex because we need to color only the border chars, not the content.
func (br *BorderRenderer) applyBorderColorToSides(line string, _ value.Border, style model.Style, hasLeft, hasRight bool) string {
	borderColor, hasBorderColor := style.GetBorderColor()
	if !hasBorderColor {
		return line
	}

	// If neither side has border, return as-is.
	if !hasLeft && !hasRight {
		return line
	}

	termCap := style.GetTerminalCapability()
	ansiCode := br.colorToANSI(borderColor, termCap, true)
	reset := br.ansiGenerator.Reset()

	// Split line to color only border characters.
	runes := []rune(line)
	var result strings.Builder

	// Left border (first character).
	if hasLeft && len(runes) > 0 {
		result.WriteString(ansiCode)
		result.WriteRune(runes[0])
		result.WriteString(reset)
		runes = runes[1:]
	}

	// Middle content (no color).
	if hasRight && len(runes) > 0 {
		// All but last character.
		result.WriteString(string(runes[:len(runes)-1]))

		// Right border (last character).
		result.WriteString(ansiCode)
		result.WriteRune(runes[len(runes)-1])
		result.WriteString(reset)
	} else {
		// No right border, write remaining content.
		result.WriteString(string(runes))
	}

	return result.String()
}

// colorToANSI converts a Color to ANSI code based on terminal capability.
func (br *BorderRenderer) colorToANSI(color value.Color, termCap value.TerminalCapability, isForeground bool) string {
	r, g, b := color.RGB()

	switch termCap {
	case value.TrueColor:
		if isForeground {
			return br.ansiGenerator.Foreground(r, g, b)
		}
		return br.ansiGenerator.Background(r, g, b)

	case value.ANSI256:
		// Convert RGB to 256-color palette.
		code := br.rgbTo256(r, g, b)
		if isForeground {
			return br.ansiGenerator.Foreground256(code)
		}
		return br.ansiGenerator.Background256(code)

	case value.ANSI16:
		// Convert RGB to 16-color palette.
		code := br.rgbTo16(r, g, b)
		if isForeground {
			return br.ansiGenerator.Foreground16(code)
		}
		return br.ansiGenerator.Background16(code)

	default:
		return ""
	}
}

// rgbTo256 converts RGB to 256-color palette index.
// Uses 6x6x6 color cube + grayscale ramp.
func (br *BorderRenderer) rgbTo256(r, g, b uint8) uint8 {
	// Grayscale detection.
	if r == g && g == b {
		// Use grayscale ramp (232-255).
		if r < 8 {
			return 16 // black
		}
		if r > 247 {
			return 231 // white
		}
		//nolint:gosec // G115: Overflow impossible - result range 0-24 from (r-8)/10 where r ∈ [8,247]
		return 232 + uint8((int(r)-8)/10)
	}

	// Color cube (16-231): 16 + 36*R + 6*G + B.
	// Where R, G, B are in range 0-5.
	//nolint:gosec // G115: Overflow impossible - result range 0-5 from (r*6)/256 where r ∈ [0,255]
	r6 := uint8((int(r) * 6) / 256)
	//nolint:gosec // G115: Overflow impossible - result range 0-5 from (g*6)/256
	g6 := uint8((int(g) * 6) / 256)
	//nolint:gosec // G115: Overflow impossible - result range 0-5 from (b*6)/256
	b6 := uint8((int(b) * 6) / 256)

	return 16 + 36*r6 + 6*g6 + b6
}

// rgbTo16 converts RGB to 16-color palette index.
// Maps to basic ANSI colors (0-15).
func (br *BorderRenderer) rgbTo16(r, g, b uint8) uint8 {
	// Calculate brightness.
	brightness := (int(r) + int(g) + int(b)) / 3

	// Determine if color is bright (8-15) or normal (0-7).
	bright := brightness > 128

	// Determine dominant color.
	red := r > 128
	green := g > 128
	blue := b > 128

	// Build color code.
	var code uint8
	if red {
		code |= 1
	}
	if green {
		code |= 2
	}
	if blue {
		code |= 4
	}

	// Add bright bit.
	if bright {
		code |= 8
	}

	return code
}
