package model

import (
	"github.com/phoenix-tui/phoenix/style/internal/domain/value"
)

// Theme represents a complete UI theme with color palette, borders, spacing, and typography.
// This is a rich domain model in DDD terms - immutable and behavior-focused.
//
// A Theme provides a consistent visual identity across all components.
// Themes can be switched at runtime without restarting the application.
type Theme struct {
	name       string
	colors     ColorPalette
	borders    BorderStyles
	spacing    SpacingScale
	typography Typography
}

// ColorPalette defines all colors used in the theme.
// Organized by semantic purpose rather than specific colors.
type ColorPalette struct {
	// Primary is the main brand color (buttons, links, focus states)
	Primary value.Color

	// Secondary is the secondary brand color (accents, highlights)
	Secondary value.Color

	// Background is the main background color
	Background value.Color

	// Surface is the color for elevated surfaces (cards, modals)
	Surface value.Color

	// Text is the primary text color
	Text value.Color

	// TextMuted is the color for secondary/muted text
	TextMuted value.Color

	// Error is the color for error states
	Error value.Color

	// Warning is the color for warning states
	Warning value.Color

	// Success is the color for success states
	Success value.Color

	// Info is the color for informational messages
	Info value.Color

	// Border is the default border color
	Border value.Color

	// Focus is the color for focused elements
	Focus value.Color

	// Disabled is the color for disabled elements
	Disabled value.Color
}

// BorderStyles defines the border styles for different UI elements.
type BorderStyles struct {
	// Default is the standard border style for most elements
	Default value.Border

	// Input is the border style for input fields
	Input value.Border

	// Modal is the border style for modal dialogs
	Modal value.Border

	// Table is the border style for tables
	Table value.Border

	// Panel is the border style for panels/containers
	Panel value.Border
}

// SpacingScale defines the spacing values used throughout the theme.
// Follows a consistent scale (typically powers of 2 or Fibonacci).
type SpacingScale struct {
	// XS is extra small spacing (2-4 cells)
	XS int

	// SM is small spacing (4-8 cells)
	SM int

	// MD is medium spacing (8-12 cells)
	MD int

	// LG is large spacing (12-16 cells)
	LG int

	// XL is extra large spacing (16-24 cells)
	XL int
}

// Typography defines text-related theme settings.
type Typography struct {
	// Placeholder color for placeholder text
	PlaceholderColor value.Color

	// Code color for code/monospace text
	CodeColor value.Color

	// Link color for links
	LinkColor value.Color

	// Heading color for headings/titles
	HeadingColor value.Color
}

// NewTheme creates a new Theme with the given name and configuration.
//nolint:gocritic // BorderStyles passed by value for immutability - struct is small enough.
// This is the primary constructor - use preset functions for common themes.
func NewTheme(
	name string,
	colors ColorPalette,
	borders BorderStyles,
	spacing SpacingScale,
	typography Typography,
) *Theme {
	return &Theme{
		name:       name,
		colors:     colors,
		borders:    borders,
		spacing:    spacing,
		typography: typography,
	}
}

// Name returns the theme name.
func (t *Theme) Name() string {
	return t.name
}

// Colors returns the color palette.
func (t *Theme) Colors() ColorPalette {
	return t.colors
}

// Borders returns the border styles.
func (t *Theme) Borders() BorderStyles {
	return t.borders
}

// Spacing returns the spacing scale.
func (t *Theme) Spacing() SpacingScale {
	return t.spacing
}

// Typography returns the typography settings.
func (t *Theme) Typography() Typography {
	return t.typography
}

// WithColors returns a new theme with the given color palette.
// Enables theme inheritance with color overrides.
func (t *Theme) WithColors(colors ColorPalette) *Theme {
	return &Theme{
		name:       t.name,
		colors:     colors,
		borders:    t.borders,
		spacing:    t.spacing,
		typography: t.typography,
	}
}

// WithBorders returns a new theme with the given border styles.
//nolint:gocritic // BorderStyles passed by value for immutability - struct is small enough.
// Enables theme inheritance with border overrides.
func (t *Theme) WithBorders(borders BorderStyles) *Theme {
	return &Theme{
		name:       t.name,
		colors:     t.colors,
		borders:    borders,
		spacing:    t.spacing,
		typography: t.typography,
	}
}

// WithSpacing returns a new theme with the given spacing scale.
// Enables theme inheritance with spacing overrides.
func (t *Theme) WithSpacing(spacing SpacingScale) *Theme {
	return &Theme{
		name:       t.name,
		colors:     t.colors,
		borders:    t.borders,
		spacing:    spacing,
		typography: t.typography,
	}
}

// WithTypography returns a new theme with the given typography settings.
// Enables theme inheritance with typography overrides.
func (t *Theme) WithTypography(typography Typography) *Theme {
	return &Theme{
		name:       t.name,
		colors:     t.colors,
		borders:    t.borders,
		spacing:    t.spacing,
		typography: typography,
	}
}

// WithName returns a new theme with the given name.
// Useful for creating theme variants.
func (t *Theme) WithName(name string) *Theme {
	return &Theme{
		name:       name,
		colors:     t.colors,
		borders:    t.borders,
		spacing:    t.spacing,
		typography: t.typography,
	}
}

// Merge combines this theme with another theme, with the other theme taking precedence.
// Only non-zero values from the other theme are applied.
// This enables theme composition and inheritance patterns.
func (t *Theme) Merge(other *Theme) *Theme {
	// For now, simple override strategy
	// In future, could implement selective merging of non-zero values
	return &Theme{
		name:       other.name,
		colors:     mergeColors(t.colors, other.colors),
		borders:    mergeBorders(t.borders, other.borders),
		spacing:    mergeSpacing(t.spacing, other.spacing),
		typography: mergeTypography(t.typography, other.typography),
	}
}

// Helper functions for merging theme components

//nolint:gocyclo,cyclop // Color palette has 13 fields - complexity is acceptable for merge logic.
func mergeColors(base, override ColorPalette) ColorPalette {
	result := base

	// Only override if the color is not the zero value (black)
	zero := value.RGB(0, 0, 0)

	if !override.Primary.Equal(zero) {
		result.Primary = override.Primary
	}
	if !override.Secondary.Equal(zero) {
		result.Secondary = override.Secondary
	}
	if !override.Background.Equal(zero) {
		result.Background = override.Background
	}
	if !override.Surface.Equal(zero) {
		result.Surface = override.Surface
	}
	if !override.Text.Equal(zero) {
		result.Text = override.Text
	}
	if !override.TextMuted.Equal(zero) {
		result.TextMuted = override.TextMuted
	}
	if !override.Error.Equal(zero) {
		result.Error = override.Error
	}
	if !override.Warning.Equal(zero) {
		result.Warning = override.Warning
	}
	if !override.Success.Equal(zero) {
		result.Success = override.Success
	}
	if !override.Info.Equal(zero) {
		result.Info = override.Info
	}
	if !override.Border.Equal(zero) {
		result.Border = override.Border
	}
	if !override.Focus.Equal(zero) {
		result.Focus = override.Focus
	}
	if !override.Disabled.Equal(zero) {
		result.Disabled = override.Disabled
	}

	return result
}

//nolint:gocritic // BorderStyles passed by value for immutability - struct is small enough.
func mergeBorders(base, override BorderStyles) BorderStyles {
	result := base

	if !override.Default.IsHidden() {
		result.Default = override.Default
	}
	if !override.Input.IsHidden() {
		result.Input = override.Input
	}
	if !override.Modal.IsHidden() {
		result.Modal = override.Modal
	}
	if !override.Table.IsHidden() {
		result.Table = override.Table
	}
	if !override.Panel.IsHidden() {
		result.Panel = override.Panel
	}

	return result
}

func mergeSpacing(base, override SpacingScale) SpacingScale {
	result := base

	if override.XS > 0 {
		result.XS = override.XS
	}
	if override.SM > 0 {
		result.SM = override.SM
	}
	if override.MD > 0 {
		result.MD = override.MD
	}
	if override.LG > 0 {
		result.LG = override.LG
	}
	if override.XL > 0 {
		result.XL = override.XL
	}

	return result
}

func mergeTypography(base, override Typography) Typography {
	result := base

	zero := value.RGB(0, 0, 0)

	if !override.PlaceholderColor.Equal(zero) {
		result.PlaceholderColor = override.PlaceholderColor
	}
	if !override.CodeColor.Equal(zero) {
		result.CodeColor = override.CodeColor
	}
	if !override.LinkColor.Equal(zero) {
		result.LinkColor = override.LinkColor
	}
	if !override.HeadingColor.Equal(zero) {
		result.HeadingColor = override.HeadingColor
	}

	return result
}
