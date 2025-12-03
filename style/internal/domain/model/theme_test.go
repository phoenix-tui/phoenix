package model

import (
	"testing"

	"github.com/phoenix-tui/phoenix/style/internal/domain/value"
	"github.com/stretchr/testify/assert"
)

func TestNewTheme(t *testing.T) {
	colors := ColorPalette{
		Primary:    value.RGB(59, 130, 246),
		Background: value.RGB(0, 0, 0),
		Text:       value.RGB(255, 255, 255),
	}

	borders := BorderStyles{
		Default: value.RoundedBorder,
		Input:   value.NormalBorder,
	}

	spacing := SpacingScale{
		XS: 2,
		SM: 4,
		MD: 8,
	}

	typography := Typography{
		PlaceholderColor: value.RGB(128, 128, 128),
	}

	theme := NewTheme("TestTheme", colors, borders, spacing, typography)

	assert.NotNil(t, theme)
	assert.Equal(t, "TestTheme", theme.Name())
	assert.Equal(t, colors, theme.Colors())
	assert.Equal(t, borders, theme.Borders())
	assert.Equal(t, spacing, theme.Spacing())
	assert.Equal(t, typography, theme.Typography())
}

func TestTheme_WithColors(t *testing.T) {
	theme := DefaultTheme()
	originalColors := theme.Colors()

	newColors := ColorPalette{
		Primary:    value.RGB(255, 0, 0),
		Background: value.RGB(255, 255, 255),
		Text:       value.RGB(0, 0, 0),
	}

	updated := theme.WithColors(newColors)

	// Original unchanged (immutability)
	assert.Equal(t, originalColors, theme.Colors())

	// Updated has new colors
	assert.Equal(t, newColors, updated.Colors())

	// Other properties unchanged
	assert.Equal(t, theme.Borders(), updated.Borders())
	assert.Equal(t, theme.Spacing(), updated.Spacing())
	assert.Equal(t, theme.Typography(), updated.Typography())
}

func TestTheme_WithBorders(t *testing.T) {
	theme := DefaultTheme()
	originalBorders := theme.Borders()

	newBorders := BorderStyles{
		Default: value.ThickBorder,
		Input:   value.DoubleBorder,
		Modal:   value.ThickBorder,
		Table:   value.NormalBorder,
		Panel:   value.RoundedBorder,
	}

	updated := theme.WithBorders(newBorders)

	// Original unchanged (immutability)
	assert.Equal(t, originalBorders, theme.Borders())

	// Updated has new borders
	assert.Equal(t, newBorders, updated.Borders())

	// Other properties unchanged
	assert.Equal(t, theme.Colors(), updated.Colors())
	assert.Equal(t, theme.Spacing(), updated.Spacing())
	assert.Equal(t, theme.Typography(), updated.Typography())
}

func TestTheme_WithSpacing(t *testing.T) {
	theme := DefaultTheme()
	originalSpacing := theme.Spacing()

	newSpacing := SpacingScale{
		XS: 4,
		SM: 8,
		MD: 16,
		LG: 24,
		XL: 32,
	}

	updated := theme.WithSpacing(newSpacing)

	// Original unchanged (immutability)
	assert.Equal(t, originalSpacing, theme.Spacing())

	// Updated has new spacing
	assert.Equal(t, newSpacing, updated.Spacing())

	// Other properties unchanged
	assert.Equal(t, theme.Colors(), updated.Colors())
	assert.Equal(t, theme.Borders(), updated.Borders())
	assert.Equal(t, theme.Typography(), updated.Typography())
}

func TestTheme_WithTypography(t *testing.T) {
	theme := DefaultTheme()
	originalTypography := theme.Typography()

	newTypography := Typography{
		PlaceholderColor: value.RGB(200, 200, 200),
		CodeColor:        value.RGB(255, 100, 100),
		LinkColor:        value.RGB(100, 100, 255),
		HeadingColor:     value.RGB(255, 255, 255),
	}

	updated := theme.WithTypography(newTypography)

	// Original unchanged (immutability)
	assert.Equal(t, originalTypography, theme.Typography())

	// Updated has new typography
	assert.Equal(t, newTypography, updated.Typography())

	// Other properties unchanged
	assert.Equal(t, theme.Colors(), updated.Colors())
	assert.Equal(t, theme.Borders(), updated.Borders())
	assert.Equal(t, theme.Spacing(), updated.Spacing())
}

func TestTheme_WithName(t *testing.T) {
	theme := DefaultTheme()

	updated := theme.WithName("CustomDefault")

	assert.Equal(t, "Default", theme.Name())
	assert.Equal(t, "CustomDefault", updated.Name())

	// All other properties unchanged
	assert.Equal(t, theme.Colors(), updated.Colors())
	assert.Equal(t, theme.Borders(), updated.Borders())
	assert.Equal(t, theme.Spacing(), updated.Spacing())
	assert.Equal(t, theme.Typography(), updated.Typography())
}

func TestTheme_Merge(t *testing.T) {
	base := DefaultTheme()

	// Create override with partial changes
	override := NewTheme(
		"Override",
		ColorPalette{
			Primary: value.RGB(255, 0, 0), // Red primary
			// Other colors left as zero (black) - should not override
		},
		BorderStyles{
			Modal: value.ThickBorder,
			// Other borders left as zero - should not override
		},
		SpacingScale{
			MD: 16, // Change MD only
			// Other values left as zero - should not override
		},
		Typography{
			LinkColor: value.RGB(0, 255, 0), // Green links
			// Other typography left as zero - should not override
		},
	)

	merged := base.Merge(override)

	// Name from override
	assert.Equal(t, "Override", merged.Name())

	// Primary color overridden
	assert.Equal(t, value.RGB(255, 0, 0), merged.Colors().Primary)

	// Other colors from base (not black)
	// Background stays from base (DefaultTheme has black background which is 0,0,0)
	assert.Equal(t, base.Colors().Background, merged.Colors().Background)
	assert.NotEqual(t, value.RGB(0, 0, 0), merged.Colors().Text)

	// Modal border overridden
	assert.Equal(t, value.ThickBorder, merged.Borders().Modal)

	// Other borders from base
	assert.Equal(t, base.Borders().Default, merged.Borders().Default)

	// MD spacing overridden
	assert.Equal(t, 16, merged.Spacing().MD)

	// Other spacing from base (not zero)
	assert.Equal(t, base.Spacing().XS, merged.Spacing().XS)
	assert.Equal(t, base.Spacing().SM, merged.Spacing().SM)

	// Link color overridden
	assert.Equal(t, value.RGB(0, 255, 0), merged.Typography().LinkColor)

	// Other typography from base
	assert.Equal(t, base.Typography().PlaceholderColor, merged.Typography().PlaceholderColor)
}

func TestTheme_Merge_FullOverride(t *testing.T) {
	base := DefaultTheme()
	override := DarkTheme()

	merged := base.Merge(override)

	// Should be identical to DarkTheme since all values are non-zero
	assert.Equal(t, override.Name(), merged.Name())
	assert.Equal(t, override.Colors(), merged.Colors())
	assert.Equal(t, override.Borders(), merged.Borders())
	assert.Equal(t, override.Spacing(), merged.Spacing())
	assert.Equal(t, override.Typography(), merged.Typography())
}

func TestTheme_Immutability(t *testing.T) {
	original := DefaultTheme()

	// Create multiple modified versions
	withColors := original.WithColors(ColorPalette{Primary: value.RGB(255, 0, 0)})
	withBorders := original.WithBorders(BorderStyles{Default: value.ThickBorder})
	withSpacing := original.WithSpacing(SpacingScale{MD: 16})

	// Original should be unchanged
	defaultTheme := DefaultTheme()
	assert.Equal(t, defaultTheme.Colors(), original.Colors())
	assert.Equal(t, defaultTheme.Borders(), original.Borders())
	assert.Equal(t, defaultTheme.Spacing(), original.Spacing())

	// Each modified version should only have its changes
	assert.Equal(t, value.RGB(255, 0, 0), withColors.Colors().Primary)
	assert.Equal(t, defaultTheme.Borders(), withColors.Borders())

	assert.Equal(t, value.ThickBorder, withBorders.Borders().Default)
	assert.Equal(t, defaultTheme.Colors(), withBorders.Colors())

	assert.Equal(t, 16, withSpacing.Spacing().MD)
	assert.Equal(t, defaultTheme.Colors(), withSpacing.Colors())
}
