package application

import (
	"sync"
	"testing"

	"github.com/phoenix-tui/phoenix/style/internal/domain/model"
	"github.com/phoenix-tui/phoenix/style/internal/domain/value"
	"github.com/stretchr/testify/assert"
)

func TestNewThemeManager_WithTheme(t *testing.T) {
	theme := model.DarkTheme()
	tm := NewThemeManager(theme)

	assert.NotNil(t, tm)
	assert.Equal(t, theme, tm.Current())
}

func TestNewThemeManager_WithNil(t *testing.T) {
	tm := NewThemeManager(nil)

	assert.NotNil(t, tm)
	assert.Equal(t, model.DefaultTheme().Name(), tm.Current().Name())
}

func TestThemeManager_Current(t *testing.T) {
	theme := model.LightTheme()
	tm := NewThemeManager(theme)

	current := tm.Current()
	assert.Equal(t, theme, current)
	assert.Equal(t, "Light", current.Name())
}

func TestThemeManager_SetTheme(t *testing.T) {
	tm := NewThemeManager(model.DefaultTheme())

	// Set new theme
	darkTheme := model.DarkTheme()
	previous := tm.SetTheme(darkTheme)

	// Verify previous theme returned
	assert.Equal(t, "Default", previous.Name())

	// Verify current theme changed
	assert.Equal(t, darkTheme, tm.Current())
	assert.Equal(t, "Dark", tm.Current().Name())
}

func TestThemeManager_SetTheme_WithNil(t *testing.T) {
	defaultTheme := model.DefaultTheme()
	tm := NewThemeManager(defaultTheme)

	// Setting nil should not change theme
	previous := tm.SetTheme(nil)

	assert.Equal(t, defaultTheme, previous)
	assert.Equal(t, defaultTheme, tm.Current())
}

func TestThemeManager_SetPreset_Found(t *testing.T) {
	tm := NewThemeManager(model.DefaultTheme())

	tests := []struct {
		name     string
		expected string
	}{
		{"dark", "Dark"},
		{"Dark", "Dark"},
		{"light", "Light"},
		{"highcontrast", "HighContrast"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := tm.SetPreset(tt.name)
			assert.True(t, ok)
			assert.Equal(t, tt.expected, tm.Current().Name())
		})
	}
}

func TestThemeManager_SetPreset_NotFound(t *testing.T) {
	tm := NewThemeManager(model.DefaultTheme())
	originalTheme := tm.Current()

	ok := tm.SetPreset("unknown")

	assert.False(t, ok)
	assert.Equal(t, originalTheme, tm.Current()) // Theme unchanged
}

func TestThemeManager_MergeTheme(t *testing.T) {
	tm := NewThemeManager(model.DefaultTheme())
	original := tm.Current()

	// Create override with only primary color changed
	override := model.NewTheme(
		"Override",
		model.ColorPalette{
			Primary: value.RGB(255, 0, 0), // Red primary
		},
		model.BorderStyles{},
		model.SpacingScale{},
		model.Typography{},
	)

	tm.MergeTheme(override)

	current := tm.Current()

	// Primary color should be overridden
	assert.Equal(t, value.RGB(255, 0, 0), current.Colors().Primary)

	// Other colors should remain from original
	assert.Equal(t, original.Colors().Background, current.Colors().Background)
	assert.Equal(t, original.Colors().Text, current.Colors().Text)

	// Other properties should remain unchanged
	assert.Equal(t, original.Borders(), current.Borders())
	assert.Equal(t, original.Spacing(), current.Spacing())
}

func TestThemeManager_Reset(t *testing.T) {
	tm := NewThemeManager(model.DarkTheme())

	// Verify starting with dark theme
	assert.Equal(t, "Dark", tm.Current().Name())

	// Reset to default
	tm.Reset()

	// Verify reset to default theme
	assert.Equal(t, "Default", tm.Current().Name())
}

func TestThemeManager_AvailablePresets(t *testing.T) {
	tm := NewThemeManager(nil)

	presets := tm.AvailablePresets()

	assert.Len(t, presets, 4)
	assert.Contains(t, presets, "Default")
	assert.Contains(t, presets, "Dark")
	assert.Contains(t, presets, "Light")
	assert.Contains(t, presets, "HighContrast")
}

func TestThemeManager_ConcurrentReads(t *testing.T) {
	tm := NewThemeManager(model.DefaultTheme())

	var wg sync.WaitGroup
	numReaders := 100

	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			// Just read multiple times
			for j := 0; j < 10; j++ {
				theme := tm.Current()
				assert.NotNil(t, theme)
			}
		}()
	}

	wg.Wait()
}

func TestThemeManager_ConcurrentWrites(t *testing.T) {
	tm := NewThemeManager(model.DefaultTheme())

	themes := []*model.Theme{
		model.DefaultTheme(),
		model.DarkTheme(),
		model.LightTheme(),
		model.HighContrastTheme(),
	}

	var wg sync.WaitGroup
	numWriters := 100

	wg.Add(numWriters)
	for i := 0; i < numWriters; i++ {
		go func(index int) {
			defer wg.Done()
			theme := themes[index%len(themes)]
			tm.SetTheme(theme)
		}(i)
	}

	wg.Wait()

	// Final theme should be one of the valid themes
	current := tm.Current()
	found := false
	for _, theme := range themes {
		if current.Name() == theme.Name() {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestThemeManager_ConcurrentReadWrite(t *testing.T) {
	tm := NewThemeManager(model.DefaultTheme())

	themes := []*model.Theme{
		model.DefaultTheme(),
		model.DarkTheme(),
		model.LightTheme(),
	}

	var wg sync.WaitGroup

	// Readers
	numReaders := 50
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				theme := tm.Current()
				assert.NotNil(t, theme)
				assert.NotEmpty(t, theme.Name())
			}
		}()
	}

	// Writers
	numWriters := 10
	wg.Add(numWriters)
	for i := 0; i < numWriters; i++ {
		go func(index int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				theme := themes[index%len(themes)]
				tm.SetTheme(theme)
			}
		}(i)
	}

	wg.Wait()

	// Should not panic and should have valid final state
	current := tm.Current()
	assert.NotNil(t, current)
	assert.NotEmpty(t, current.Name())
}

func TestThemeManager_ConcurrentSetPreset(t *testing.T) {
	tm := NewThemeManager(model.DefaultTheme())

	presets := []string{"default", "dark", "light", "highcontrast"}

	var wg sync.WaitGroup
	numGoroutines := 100

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			preset := presets[index%len(presets)]
			tm.SetPreset(preset)
		}(i)
	}

	wg.Wait()

	// Should have one of the valid themes
	current := tm.Current()
	assert.NotNil(t, current)
	assert.Contains(t, []string{"Default", "Dark", "Light", "HighContrast"}, current.Name())
}

func TestThemeManager_ConcurrentMerge(t *testing.T) {
	tm := NewThemeManager(model.DefaultTheme())

	var wg sync.WaitGroup
	numGoroutines := 50

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()

			override := model.NewTheme(
				"Override",
				model.ColorPalette{
					Primary: value.RGB(uint8(index), 0, 0),
				},
				model.BorderStyles{},
				model.SpacingScale{},
				model.Typography{},
			)

			tm.MergeTheme(override)
		}(i)
	}

	wg.Wait()

	// Should have valid final state
	current := tm.Current()
	assert.NotNil(t, current)

	// Primary color should have been set by one of the goroutines
	r, _, _ := current.Colors().Primary.RGB()
	assert.GreaterOrEqual(t, r, uint8(0))
	assert.Less(t, r, uint8(50)) // Should be less than numGoroutines
}
