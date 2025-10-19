// Package service contains domain services that orchestrate complex domain logic.
// Domain services operate on domain models and value objects.
package service

import (
	"time"

	"github.com/phoenix-tui/phoenix/mouse/domain/model"
	"github.com/phoenix-tui/phoenix/mouse/domain/value"
)

// ClickDetector is a domain service that detects single, double, and triple clicks.
type ClickDetector struct {
	lastClick  *model.MouseEvent
	clickCount int
	timeout    time.Duration
	tolerance  int // Position tolerance for multi-click detection (cells)
}

// NewClickDetector creates a new ClickDetector.
// timeout: maximum time between clicks for multi-click detection (typically 500ms)
// tolerance: maximum distance between clicks for multi-click detection (typically 1 cell)
func NewClickDetector(timeout time.Duration, tolerance int) *ClickDetector {
	if timeout <= 0 {
		timeout = 500 * time.Millisecond
	}
	if tolerance < 0 {
		tolerance = 1
	}
	return &ClickDetector{
		timeout:   timeout,
		tolerance: tolerance,
	}
}

// DetectClick processes a release event and determines if it's a click, double-click, or triple-click.
// Returns a new event with the appropriate click type, or nil if not a click.
func (d *ClickDetector) DetectClick(releaseEvent model.MouseEvent) *model.MouseEvent {
	// Only process release events
	if releaseEvent.Type() != value.EventRelease {
		return nil
	}

	// Check if this is a continuation of previous clicks
	if d.lastClick != nil {
		timeSinceLastClick := releaseEvent.Timestamp().Sub(d.lastClick.Timestamp())
		samePosition := releaseEvent.IsAt(d.lastClick.Position(), d.tolerance)
		sameButton := releaseEvent.Button() == d.lastClick.Button()

		if timeSinceLastClick <= d.timeout && samePosition && sameButton {
			// This is a multi-click
			d.clickCount++
		} else {
			// Too much time passed, different position, or different button - reset
			d.clickCount = 1
		}
	} else {
		// First click
		d.clickCount = 1
	}

	// Determine click type based on count
	var clickType value.EventType
	switch d.clickCount {
	case 1:
		clickType = value.EventClick
	case 2:
		clickType = value.EventDoubleClick
	case 3:
		clickType = value.EventTripleClick
	default:
		// Reset after triple click
		d.clickCount = 1
		clickType = value.EventClick
	}

	// Create click event
	clickEvent := releaseEvent.WithType(clickType)

	// Store this click for next comparison
	d.lastClick = &clickEvent

	return &clickEvent
}

// Reset resets the click detector state.
func (d *ClickDetector) Reset() {
	d.lastClick = nil
	d.clickCount = 0
}

// Timeout returns the timeout duration.
func (d *ClickDetector) Timeout() time.Duration {
	return d.timeout
}

// Tolerance returns the position tolerance.
func (d *ClickDetector) Tolerance() int {
	return d.tolerance
}

// ClickCount returns the current click count.
func (d *ClickDetector) ClickCount() int {
	return d.clickCount
}
