package model

import "errors"

// RawMode represents terminal raw mode state (entity with lifecycle).
//
// Raw mode disables:
//   - Line buffering (input available immediately)
//   - Echo (typed characters not displayed)
//   - Signal generation (Ctrl+C doesn't send SIGINT)
//   - Input processing (Ctrl+S/Q, CR->NL translation)
//
// This entity stores the original terminal state for restoration.
//
// Invariants:
//   - OriginalState is never nil after creation
//   - Enabled flag accurately reflects raw mode status
//
// Lifecycle:
//  1. Created with original terminal state
//  2. Enabled (raw mode activated)
//  3. Disabled (terminal restored)
//
// Platform-specific state:
//   - Unix: syscall.Termios
//   - Windows: uint32 (console mode)
type RawMode struct {
	enabled       bool
	originalState interface{} // Platform-specific (syscall.Termios on Unix, uint32 on Windows)
}

// NewRawMode creates raw mode entity with original state.
// The state type depends on platform:
//   - Unix: syscall.Termios
//   - Windows: uint32
func NewRawMode(originalState interface{}) (*RawMode, error) {
	if originalState == nil {
		return nil, errors.New("originalState cannot be nil")
	}

	return &RawMode{
		enabled:       false, // Not enabled yet
		originalState: originalState,
	}, nil
}

// IsEnabled returns true if raw mode is active.
func (r *RawMode) IsEnabled() bool {
	return r.enabled
}

// OriginalState returns the saved terminal state for restoration.
// Type depends on platform - caller must type assert.
func (r *RawMode) OriginalState() interface{} {
	return r.originalState
}

// Enable marks raw mode as enabled (immutable).
// Returns new RawMode instance with enabled flag set.
func (r *RawMode) Enable() *RawMode {
	newMode := *r
	newMode.enabled = true
	return &newMode
}

// Disable marks raw mode as disabled (immutable).
// Returns new RawMode instance with enabled flag cleared.
func (r *RawMode) Disable() *RawMode {
	newMode := *r
	newMode.enabled = false
	return &newMode
}
