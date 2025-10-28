// Package model provides rich domain models for terminal operations.
package model

import (
	value2 "github.com/phoenix-tui/phoenix/core/internal/domain/value"
)

// Terminal represents a terminal instance with its capabilities and state.
// This is the aggregate root for terminal operations.
//
// Invariants:
//   - Capabilities is never nil
//   - Size is always valid (width, height >= 1)
//   - Terminal is immutable after creation
//
// Lifecycle:
//   - Created with detected capabilities
//   - Can be put into raw mode (creates RawMode entity)
//   - Size can be updated on SIGWINCH
type Terminal struct {
	capabilities *value2.Capabilities
	rawMode      *RawMode
	size         value2.Size
}

// NewTerminal creates a new Terminal with detected capabilities.
// Initial size defaults to 80x24 if not specified.
func NewTerminal(caps *value2.Capabilities) *Terminal {
	if caps == nil {
		// Defensive: create minimal capabilities if nil
		caps = value2.NewCapabilities(false, value2.ColorDepthNone, false, false, false)
	}

	return &Terminal{
		capabilities: caps,
		rawMode:      nil,                    // Not in raw mode initially
		size:         value2.NewSize(80, 24), // Standard VT100 default
	}
}

// Capabilities returns terminal capabilities (immutable).
func (t *Terminal) Capabilities() *value2.Capabilities {
	return t.capabilities
}

// Size returns current terminal size.
func (t *Terminal) Size() value2.Size {
	return t.size
}

// IsRawMode returns true if terminal is in raw mode.
func (t *Terminal) IsRawMode() bool {
	return t.rawMode != nil && t.rawMode.IsEnabled()
}

// RawMode returns the raw mode entity (nil if not in raw mode).
func (t *Terminal) RawMode() *RawMode {
	return t.rawMode
}

// WithSize returns new Terminal with updated size (immutable).
// This is called when SIGWINCH is received.
func (t *Terminal) WithSize(size value2.Size) *Terminal {
	newTerm := *t
	newTerm.size = size
	return &newTerm
}

// WithRawMode returns new Terminal in raw mode (immutable).
// The RawMode entity contains the original terminal state for restoration.
func (t *Terminal) WithRawMode(rawMode *RawMode) *Terminal {
	newTerm := *t
	newTerm.rawMode = rawMode
	return &newTerm
}

// WithoutRawMode returns new Terminal not in raw mode (immutable).
// This is called after raw mode is successfully disabled.
func (t *Terminal) WithoutRawMode() *Terminal {
	newTerm := *t
	newTerm.rawMode = nil
	return &newTerm
}

// SupportsANSI is a convenience method for checking ANSI support.
func (t *Terminal) SupportsANSI() bool {
	return t.capabilities.SupportsANSI()
}

// SupportsColor is a convenience method for checking color support.
func (t *Terminal) SupportsColor() bool {
	return t.capabilities.SupportsColor()
}

// ColorDepth is a convenience method for getting color depth.
func (t *Terminal) ColorDepth() value2.ColorDepth {
	return t.capabilities.ColorDepth()
}
