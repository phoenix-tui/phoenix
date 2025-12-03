package model

import (
	"testing"

	"github.com/phoenix-tui/phoenix/style/internal/domain/value"
	"github.com/stretchr/testify/assert"
)

func TestDefaultTheme(t *testing.T) {
	theme := DefaultTheme()

	assert.NotNil(t, theme)
	assert.Equal(t, "Default", theme.Name())

	// Verify color palette
	colors := theme.Colors()
	assert.Equal(t, value.RGB(59, 130, 246), colors.Primary)
	assert.Equal(t, value.RGB(0, 0, 0), colors.Background)
	assert.Equal(t, value.RGB(255, 255, 255), colors.Text)

	// Verify borders
	borders := theme.Borders()
	assert.Equal(t, value.RoundedBorder, borders.Default)
	assert.Equal(t, value.RoundedBorder, borders.Input)
	assert.Equal(t, value.ThickBorder, borders.Modal)

	// Verify spacing
	spacing := theme.Spacing()
	assert.Equal(t, 2, spacing.XS)
	assert.Equal(t, 4, spacing.SM)
	assert.Equal(t, 8, spacing.MD)
	assert.Equal(t, 12, spacing.LG)
	assert.Equal(t, 16, spacing.XL)

	// Verify typography
	typography := theme.Typography()
	assert.NotEqual(t, value.RGB(0, 0, 0), typography.PlaceholderColor)
}

func TestDarkTheme(t *testing.T) {
	theme := DarkTheme()

	assert.NotNil(t, theme)
	assert.Equal(t, "Dark", theme.Name())

	// Dark theme should have black background
	colors := theme.Colors()
	assert.Equal(t, value.RGB(0, 0, 0), colors.Background)

	// Text should be light
	r, g, b := colors.Text.RGB()
	assert.Greater(t, int(r)+int(g)+int(b), 600) // Sum > 600 means light color

	// Verify spacing matches Default
	spacing := theme.Spacing()
	assert.Equal(t, 2, spacing.XS)
	assert.Equal(t, 4, spacing.SM)
	assert.Equal(t, 8, spacing.MD)
}

func TestLightTheme(t *testing.T) {
	theme := LightTheme()

	assert.NotNil(t, theme)
	assert.Equal(t, "Light", theme.Name())

	// Light theme should have white background
	colors := theme.Colors()
	assert.Equal(t, value.RGB(255, 255, 255), colors.Background)

	// Text should be dark
	r, g, b := colors.Text.RGB()
	assert.Less(t, int(r)+int(g)+int(b), 100) // Sum < 100 means dark color

	// Verify spacing matches Default
	spacing := theme.Spacing()
	assert.Equal(t, 2, spacing.XS)
	assert.Equal(t, 4, spacing.SM)
	assert.Equal(t, 8, spacing.MD)
}

func TestHighContrastTheme(t *testing.T) {
	theme := HighContrastTheme()

	assert.NotNil(t, theme)
	assert.Equal(t, "HighContrast", theme.Name())

	colors := theme.Colors()

	// Background should be pure black
	assert.Equal(t, value.RGB(0, 0, 0), colors.Background)

	// Text should be pure white
	assert.Equal(t, value.RGB(255, 255, 255), colors.Text)

	// Primary colors should be pure (high contrast)
	assert.Equal(t, value.RGB(255, 0, 0), colors.Error)
	assert.Equal(t, value.RGB(0, 255, 0), colors.Success)
	assert.Equal(t, value.RGB(0, 0, 255), colors.Primary)

	// Borders should be thick/double for visibility
	borders := theme.Borders()
	assert.Equal(t, value.DoubleBorder, borders.Default)
	assert.Equal(t, value.ThickBorder, borders.Modal)

	// Spacing should be slightly larger
	spacing := theme.Spacing()
	assert.Greater(t, spacing.SM, 4)
	assert.Greater(t, spacing.MD, 8)
}

func TestAllPresets(t *testing.T) {
	presets := AllPresets()

	assert.Len(t, presets, 4)

	// Verify all themes are present
	names := make(map[string]bool)
	for _, theme := range presets {
		names[theme.Name()] = true
	}

	assert.True(t, names["Default"])
	assert.True(t, names["Dark"])
	assert.True(t, names["Light"])
	assert.True(t, names["HighContrast"])

	// Verify all themes are unique
	assert.Len(t, names, 4)
}

func TestPresetByName_Found(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"default", "Default"},
		{"Default", "Default"},
		{"DEFAULT", "Default"},
		{"dark", "Dark"},
		{"Dark", "Dark"},
		{"light", "Light"},
		{"highcontrast", "HighContrast"},
		{"HighContrast", "HighContrast"},
		{"HIGHCONTRAST", "HighContrast"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := PresetByName(tt.name)
			assert.NotNil(t, theme)
			assert.Equal(t, tt.expected, theme.Name())
		})
	}
}

func TestPresetByName_NotFound(t *testing.T) {
	tests := []string{
		"unknown",
		"",
		"custom",
		"blue",
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			theme := PresetByName(name)
			assert.Nil(t, theme)
		})
	}
}

func TestThemePresets_Consistency(t *testing.T) {
	// All themes should have consistent spacing scales
	themes := AllPresets()

	for _, theme := range themes {
		t.Run(theme.Name(), func(t *testing.T) {
			spacing := theme.Spacing()

			// Verify ascending order
			assert.Less(t, spacing.XS, spacing.SM)
			assert.Less(t, spacing.SM, spacing.MD)
			assert.Less(t, spacing.MD, spacing.LG)
			assert.Less(t, spacing.LG, spacing.XL)

			// Verify positive values
			assert.Greater(t, spacing.XS, 0)
			assert.Greater(t, spacing.SM, 0)
			assert.Greater(t, spacing.MD, 0)
			assert.Greater(t, spacing.LG, 0)
			assert.Greater(t, spacing.XL, 0)
		})
	}
}

func TestThemePresets_ColorContrast(t *testing.T) {
	// Verify text/background contrast for readability

	tests := []struct {
		name  string
		theme *Theme
	}{
		{"Default", DefaultTheme()},
		{"Dark", DarkTheme()},
		{"Light", LightTheme()},
		{"HighContrast", HighContrastTheme()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colors := tt.theme.Colors()

			// Calculate luminance difference (simple approximation)
			bgR, bgG, bgB := colors.Background.RGB()
			textR, textG, textB := colors.Text.RGB()

			bgLum := int(bgR) + int(bgG) + int(bgB)
			textLum := int(textR) + int(textG) + int(textB)

			contrast := textLum - bgLum
			if contrast < 0 {
				contrast = -contrast
			}

			// Contrast should be significant (at least 200 for regular themes, 600 for high contrast)
			if tt.name == "HighContrast" {
				assert.Greater(t, contrast, 600, "High contrast theme should have very high contrast")
			} else {
				assert.Greater(t, contrast, 200, "Theme should have readable contrast")
			}
		})
	}
}

func TestThemePresets_BorderStyles(t *testing.T) {
	// All themes should have valid borders (not hidden)

	themes := AllPresets()

	for _, theme := range themes {
		t.Run(theme.Name(), func(t *testing.T) {
			borders := theme.Borders()

			// None of the borders should be completely hidden
			assert.False(t, borders.Default.IsHidden())
			assert.False(t, borders.Input.IsHidden())
			assert.False(t, borders.Modal.IsHidden())
			assert.False(t, borders.Table.IsHidden())
			assert.False(t, borders.Panel.IsHidden())
		})
	}
}

func TestThemePresets_Typography(t *testing.T) {
	// All themes should have valid typography colors

	themes := AllPresets()
	zero := value.RGB(0, 0, 0)

	for _, theme := range themes {
		t.Run(theme.Name(), func(t *testing.T) {
			typography := theme.Typography()

			// PlaceholderColor should not be zero (except in dark/light themes where it might be black)
			// For most themes, it should be distinct
			if theme.Name() != "Light" {
				assert.NotEqual(t, zero, typography.PlaceholderColor)
			}

			// Other colors should not be zero
			assert.NotEqual(t, zero, typography.CodeColor)
			assert.NotEqual(t, zero, typography.LinkColor)
			assert.NotEqual(t, zero, typography.HeadingColor)
		})
	}
}
