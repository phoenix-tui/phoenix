// Package platform manages platform-specific terminal mode operations.
// It handles writing ANSI sequences to stdout and tracking terminal state.
package platform

import (
	"fmt"
	"os"

	"github.com/phoenix-tui/phoenix/mouse/infrastructure/ansi"
)

// TerminalMode manages terminal mouse mode state.
type TerminalMode struct {
	enabled bool
	mode    ansi.MouseMode
}

// NewTerminalMode creates a new TerminalMode.
func NewTerminalMode() *TerminalMode {
	return &TerminalMode{
		enabled: false,
		mode:    ansi.ModeSGR, // Default to modern SGR protocol
	}
}

// Enable enables mouse tracking with the specified mode.
func (t *TerminalMode) Enable(mode ansi.MouseMode) error {
	if t.enabled {
		// Already enabled, switch mode
		if err := t.Disable(); err != nil {
			return fmt.Errorf("failed to disable previous mode: %w", err)
		}
	}

	// Write enable sequence
	seq := ansi.EnableMouseTracking(mode)
	if _, err := os.Stdout.WriteString(seq); err != nil {
		return fmt.Errorf("failed to enable mouse tracking: %w", err)
	}

	// Also enable button-motion tracking for drag support
	if mode == ansi.ModeSGR {
		seq = ansi.EnableMouseTracking(ansi.ModeButtonMotion)
		if _, err := os.Stdout.WriteString(seq); err != nil {
			return fmt.Errorf("failed to enable button motion tracking: %w", err)
		}
	}

	t.enabled = true
	t.mode = mode
	return nil
}

// Disable disables mouse tracking.
func (t *TerminalMode) Disable() error {
	if !t.enabled {
		return nil
	}

	// Write disable sequence
	seq := ansi.DisableMouseTracking(t.mode)
	if _, err := os.Stdout.WriteString(seq); err != nil {
		return fmt.Errorf("failed to disable mouse tracking: %w", err)
	}

	// Also disable button-motion tracking
	if t.mode == ansi.ModeSGR {
		seq = ansi.DisableMouseTracking(ansi.ModeButtonMotion)
		if _, err := os.Stdout.WriteString(seq); err != nil {
			return fmt.Errorf("failed to disable button motion tracking: %w", err)
		}
	}

	t.enabled = false
	return nil
}

// IsEnabled returns true if mouse tracking is enabled.
func (t *TerminalMode) IsEnabled() bool {
	return t.enabled
}

// Mode returns the current mouse tracking mode.
func (t *TerminalMode) Mode() ansi.MouseMode {
	return t.mode
}

// EnableAll enables comprehensive mouse tracking (SGR + button motion).
func (t *TerminalMode) EnableAll() error {
	if t.enabled {
		if err := t.Disable(); err != nil {
			return err
		}
	}

	// Write enable sequence
	seq := ansi.EnableMouseAll()
	if _, err := os.Stdout.WriteString(seq); err != nil {
		return fmt.Errorf("failed to enable mouse tracking: %w", err)
	}

	t.enabled = true
	t.mode = ansi.ModeSGR
	return nil
}

// DisableAll disables comprehensive mouse tracking.
func (t *TerminalMode) DisableAll() error {
	if !t.enabled {
		return nil
	}

	// Write disable sequence
	seq := ansi.DisableMouseAll()
	if _, err := os.Stdout.WriteString(seq); err != nil {
		return fmt.Errorf("failed to disable mouse tracking: %w", err)
	}

	t.enabled = false
	return nil
}
