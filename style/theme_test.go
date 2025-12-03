package style

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultTheme_PublicAPI(t *testing.T) {
	theme := DefaultTheme()

	assert.NotNil(t, theme)
	assert.Equal(t, "Default", theme.Name())
}

func TestDarkTheme_PublicAPI(t *testing.T) {
	theme := DarkTheme()

	assert.NotNil(t, theme)
	assert.Equal(t, "Dark", theme.Name())
}

func TestLightTheme_PublicAPI(t *testing.T) {
	theme := LightTheme()

	assert.NotNil(t, theme)
	assert.Equal(t, "Light", theme.Name())
}

func TestHighContrastTheme_PublicAPI(t *testing.T) {
	theme := HighContrastTheme()

	assert.NotNil(t, theme)
	assert.Equal(t, "HighContrast", theme.Name())
}

func TestAllThemes(t *testing.T) {
	themes := AllThemes()

	assert.Len(t, themes, 4)
}

func TestThemeByName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"default", "Default"},
		{"dark", "Dark"},
		{"light", "Light"},
		{"highcontrast", "HighContrast"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := ThemeByName(tt.name)
			assert.NotNil(t, theme)
			assert.Equal(t, tt.expected, theme.Name())
		})
	}
}

func TestThemeByName_NotFound(t *testing.T) {
	theme := ThemeByName("unknown")
	assert.Nil(t, theme)
}

func TestNewThemeManager_PublicAPI(t *testing.T) {
	tm := NewThemeManager(nil)

	assert.NotNil(t, tm)
	assert.Equal(t, "Default", tm.Current().Name())
}

func TestNewThemeManager_WithTheme(t *testing.T) {
	darkTheme := DarkTheme()
	tm := NewThemeManager(darkTheme)

	assert.NotNil(t, tm)
	assert.Equal(t, "Dark", tm.Current().Name())
}

func TestThemeManager_SetTheme_PublicAPI(t *testing.T) {
	tm := NewThemeManager(DefaultTheme())

	lightTheme := LightTheme()
	previous := tm.SetTheme(lightTheme)

	assert.Equal(t, "Default", previous.Name())
	assert.Equal(t, "Light", tm.Current().Name())
}

func TestThemeManager_SetPreset_PublicAPI(t *testing.T) {
	tm := NewThemeManager(DefaultTheme())

	ok := tm.SetPreset("dark")
	assert.True(t, ok)
	assert.Equal(t, "Dark", tm.Current().Name())
}

func TestThemeManager_Reset_PublicAPI(t *testing.T) {
	tm := NewThemeManager(DarkTheme())

	tm.Reset()

	assert.Equal(t, "Default", tm.Current().Name())
}

func TestNewColorPalette(t *testing.T) {
	palette := NewColorPalette(
		RGB(255, 0, 0),   // primary
		RGB(0, 255, 0),   // secondary
		RGB(0, 0, 0),     // background
		RGB(10, 10, 10),  // surface
		RGB(255, 255, 255), // text
		RGB(128, 128, 128), // textMuted
		RGB(255, 0, 0),   // error
		RGB(255, 255, 0), // warning
		RGB(0, 255, 0),   // success
		RGB(0, 0, 255),   // info
		RGB(100, 100, 100), // border
		RGB(0, 255, 255), // focus
		RGB(80, 80, 80),  // disabled
	)

	assert.Equal(t, RGB(255, 0, 0), palette.Primary)
	assert.Equal(t, RGB(0, 255, 0), palette.Secondary)
}

func TestNewBorderStyles(t *testing.T) {
	borders := NewBorderStyles(
		RoundedBorder,
		NormalBorder,
		ThickBorder,
		DoubleBorder,
		RoundedBorder,
	)

	assert.Equal(t, RoundedBorder, borders.Default)
	assert.Equal(t, NormalBorder, borders.Input)
	assert.Equal(t, ThickBorder, borders.Modal)
}

func TestNewSpacingScale(t *testing.T) {
	spacing := NewSpacingScale(2, 4, 8, 12, 16)

	assert.Equal(t, 2, spacing.XS)
	assert.Equal(t, 4, spacing.SM)
	assert.Equal(t, 8, spacing.MD)
	assert.Equal(t, 12, spacing.LG)
	assert.Equal(t, 16, spacing.XL)
}

func TestNewTypography(t *testing.T) {
	typography := NewTypography(
		RGB(128, 128, 128),
		RGB(255, 100, 100),
		RGB(100, 100, 255),
		RGB(255, 255, 255),
	)

	assert.Equal(t, RGB(128, 128, 128), typography.PlaceholderColor)
	assert.Equal(t, RGB(255, 100, 100), typography.CodeColor)
	assert.Equal(t, RGB(100, 100, 255), typography.LinkColor)
	assert.Equal(t, RGB(255, 255, 255), typography.HeadingColor)
}

func TestNewTheme_PublicAPI(t *testing.T) {
	colors := NewColorPalette(
		RGB(59, 130, 246),
		RGB(139, 92, 246),
		RGB(0, 0, 0),
		RGB(30, 30, 30),
		RGB(255, 255, 255),
		RGB(156, 163, 175),
		RGB(239, 68, 68),
		RGB(245, 158, 11),
		RGB(34, 197, 94),
		RGB(59, 130, 246),
		RGB(75, 85, 99),
		RGB(96, 165, 250),
		RGB(107, 114, 128),
	)

	borders := NewBorderStyles(
		RoundedBorder,
		RoundedBorder,
		ThickBorder,
		NormalBorder,
		RoundedBorder,
	)

	spacing := NewSpacingScale(2, 4, 8, 12, 16)

	typography := NewTypography(
		RGB(107, 114, 128),
		RGB(236, 72, 153),
		RGB(96, 165, 250),
		RGB(255, 255, 255),
	)

	theme := NewTheme("CustomTheme", colors, borders, spacing, typography)

	assert.NotNil(t, theme)
	assert.Equal(t, "CustomTheme", theme.Name())
	assert.Equal(t, colors, theme.Colors())
	assert.Equal(t, borders, theme.Borders())
	assert.Equal(t, spacing, theme.Spacing())
	assert.Equal(t, typography, theme.Typography())
}
