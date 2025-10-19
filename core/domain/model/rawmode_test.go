package model_test

import (
	"testing"

	"github.com/phoenix-tui/phoenix/core/domain/model"
)

func TestNewRawMode(t *testing.T) {
	tests := []struct {
		name          string
		originalState interface{}
		wantErr       bool
	}{
		{
			name:          "valid state (string)",
			originalState: "terminal-state",
			wantErr:       false,
		},
		{
			name:          "valid state (struct)",
			originalState: struct{ mode int }{mode: 123},
			wantErr:       false,
		},
		{
			name:          "nil state (invalid)",
			originalState: nil,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawMode, err := model.NewRawMode(tt.originalState)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if rawMode != nil {
					t.Error("expected nil raw mode on error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if rawMode == nil {
				t.Fatal("expected raw mode, got nil")
			}

			if rawMode.IsEnabled() {
				t.Error("new raw mode should not be enabled")
			}

			if rawMode.OriginalState() != tt.originalState {
				t.Error("original state not preserved")
			}
		})
	}
}

func TestRawMode_Enable(t *testing.T) {
	rawMode, err := model.NewRawMode("original-state")
	if err != nil {
		t.Fatalf("failed to create raw mode: %v", err)
	}

	// Initially not enabled
	if rawMode.IsEnabled() {
		t.Fatal("new raw mode should not be enabled")
	}

	// Enable
	enabledMode := rawMode.Enable()

	// Check enabled
	if !enabledMode.IsEnabled() {
		t.Error("raw mode should be enabled")
	}

	// Check original unchanged (immutability)
	if rawMode.IsEnabled() {
		t.Error("original raw mode was mutated")
	}

	if enabledMode == rawMode {
		t.Error("Enable should return new instance")
	}

	// Check original state preserved
	if enabledMode.OriginalState() != "original-state" {
		t.Error("original state was lost during enable")
	}
}

func TestRawMode_Disable(t *testing.T) {
	rawMode, err := model.NewRawMode("original-state")
	if err != nil {
		t.Fatalf("failed to create raw mode: %v", err)
	}

	// Enable first
	enabledMode := rawMode.Enable()

	if !enabledMode.IsEnabled() {
		t.Fatal("raw mode should be enabled")
	}

	// Disable
	disabledMode := enabledMode.Disable()

	// Check disabled
	if disabledMode.IsEnabled() {
		t.Error("raw mode should be disabled")
	}

	// Check original unchanged (immutability)
	if !enabledMode.IsEnabled() {
		t.Error("original raw mode was mutated")
	}

	if disabledMode == enabledMode {
		t.Error("Disable should return new instance")
	}

	// Check original state preserved
	if disabledMode.OriginalState() != "original-state" {
		t.Error("original state was lost during disable")
	}
}

func TestRawMode_Lifecycle(t *testing.T) {
	// Test complete lifecycle: create → enable → disable
	rawMode, err := model.NewRawMode("original-state")
	if err != nil {
		t.Fatalf("failed to create raw mode: %v", err)
	}

	// Phase 1: Created (not enabled)
	if rawMode.IsEnabled() {
		t.Error("phase 1: should not be enabled")
	}

	// Phase 2: Enabled
	enabledMode := rawMode.Enable()
	if !enabledMode.IsEnabled() {
		t.Error("phase 2: should be enabled")
	}

	// Phase 3: Disabled
	disabledMode := enabledMode.Disable()
	if disabledMode.IsEnabled() {
		t.Error("phase 3: should be disabled")
	}

	// Verify all phases preserved original state
	if rawMode.OriginalState() != "original-state" {
		t.Error("phase 1: lost original state")
	}
	if enabledMode.OriginalState() != "original-state" {
		t.Error("phase 2: lost original state")
	}
	if disabledMode.OriginalState() != "original-state" {
		t.Error("phase 3: lost original state")
	}

	// Verify all phases are different instances
	if rawMode == enabledMode || enabledMode == disabledMode {
		t.Error("lifecycle phases returned same instance (not immutable)")
	}
}

func TestRawMode_Immutability(t *testing.T) {
	// Verify RawMode is truly immutable
	rawMode1, _ := model.NewRawMode("state1")

	// Store original values
	originalEnabled := rawMode1.IsEnabled()
	originalState := rawMode1.OriginalState()

	// Perform multiple operations
	rawMode2 := rawMode1.Enable()
	rawMode3 := rawMode2.Disable()
	rawMode4 := rawMode3.Enable()

	// Verify rawMode1 unchanged
	if rawMode1.IsEnabled() != originalEnabled {
		t.Error("rawMode1 enabled flag was mutated")
	}

	if rawMode1.OriginalState() != originalState {
		t.Error("rawMode1 original state was mutated")
	}

	// Verify all instances are different
	if rawMode1 == rawMode2 || rawMode1 == rawMode3 || rawMode1 == rawMode4 {
		t.Error("operations returned same instance (not immutable)")
	}
}

func TestRawMode_OriginalStateTypes(t *testing.T) {
	// Test that different platform-specific state types work
	// Note: Using comparable types only (struct, primitives)
	// Maps are uncomparable and would cause panic with !=

	type unixState struct{ Iflag, Oflag, Cflag, Lflag uint32 }

	tests := []struct {
		name      string
		state     interface{}
		checkFunc func(interface{}) bool
	}{
		{
			name:  "Unix (simulated syscall.Termios)",
			state: unixState{123, 456, 789, 012},
			checkFunc: func(s interface{}) bool {
				us, ok := s.(unixState)
				return ok && us.Iflag == 123 && us.Oflag == 456
			},
		},
		{
			name:  "Windows (simulated console mode)",
			state: uint32(0x0001 | 0x0002 | 0x0004),
			checkFunc: func(s interface{}) bool {
				mode, ok := s.(uint32)
				return ok && mode == uint32(0x0001|0x0002|0x0004)
			},
		},
		{
			name:  "String state",
			state: "raw-mode-active",
			checkFunc: func(s interface{}) bool {
				str, ok := s.(string)
				return ok && str == "raw-mode-active"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawMode, err := model.NewRawMode(tt.state)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tt.checkFunc(rawMode.OriginalState()) {
				t.Error("original state type not preserved")
			}

			// Verify state preserved through lifecycle
			enabled := rawMode.Enable()
			if !tt.checkFunc(enabled.OriginalState()) {
				t.Error("state lost during enable")
			}

			disabled := enabled.Disable()
			if !tt.checkFunc(disabled.OriginalState()) {
				t.Error("state lost during disable")
			}
		})
	}
}
