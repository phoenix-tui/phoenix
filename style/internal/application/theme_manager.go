// Package application provides application services for the style module.
package application

import (
	"sync"

	"github.com/phoenix-tui/phoenix/style/internal/domain/model"
)

// ThemeManager manages the current theme and provides thread-safe theme switching.
// This is an application service in DDD terms - orchestrates domain logic.
//
// ThemeManager is safe for concurrent use from multiple goroutines.
// Theme changes are atomic and immediately visible to all readers.
type ThemeManager struct {
	mu      sync.RWMutex
	current *model.Theme
}

// NewThemeManager creates a new ThemeManager with the given initial theme.
// If theme is nil, DefaultTheme is used.
func NewThemeManager(theme *model.Theme) *ThemeManager {
	if theme == nil {
		theme = model.DefaultTheme()
	}

	return &ThemeManager{
		current: theme,
	}
}

// Current returns the current theme.
// This is safe for concurrent access.
func (tm *ThemeManager) Current() *model.Theme {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.current
}

// SetTheme sets the current theme.
// This is safe for concurrent access.
// Returns the previous theme for potential undo functionality.
func (tm *ThemeManager) SetTheme(theme *model.Theme) *model.Theme {
	if theme == nil {
		return tm.Current()
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	previous := tm.current
	tm.current = theme
	return previous
}

// SetPreset sets the current theme by preset name.
// Returns true if the preset was found and applied, false otherwise.
// This is safe for concurrent access.
func (tm *ThemeManager) SetPreset(name string) bool {
	preset := model.PresetByName(name)
	if preset == nil {
		return false
	}

	tm.SetTheme(preset)
	return true
}

// MergeTheme merges the given theme with the current theme.
// The given theme takes precedence for non-zero values.
// This enables runtime theme customization without replacing the entire theme.
func (tm *ThemeManager) MergeTheme(override *model.Theme) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.current = tm.current.Merge(override)
}

// Reset resets the theme to the default theme.
func (tm *ThemeManager) Reset() {
	tm.SetTheme(model.DefaultTheme())
}

// AvailablePresets returns all available preset theme names.
func (tm *ThemeManager) AvailablePresets() []string {
	return []string{"Default", "Dark", "Light", "HighContrast"}
}
