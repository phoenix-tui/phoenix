package model

import (
	"testing"
	"time"

	value2 "github.com/phoenix-tui/phoenix/mouse/internal/domain/value"
)

// TestNewHoverState tests the constructor.
func TestNewHoverState(t *testing.T) {
	state := NewHoverState()

	if state == nil {
		t.Fatal("NewHoverState() returned nil")
	}

	if state.ComponentID() != "" {
		t.Errorf("Expected empty componentID, got %s", state.ComponentID())
	}

	if state.IsActive() {
		t.Error("Expected IsActive() to be false for new state")
	}

	if state.IsHovering() {
		t.Error("Expected IsHovering() to be false for new state")
	}

	// Position should be initialized to 0,0
	pos := state.Position()
	if pos.X() != 0 || pos.Y() != 0 {
		t.Errorf("Expected position (0,0), got (%d,%d)", pos.X(), pos.Y())
	}
}

// TestHoverStateEnter tests the Enter method.
func TestHoverStateEnter(t *testing.T) {
	tests := []struct {
		name        string
		componentID string
		position    value2.Position
	}{
		{
			name:        "enter component at origin",
			componentID: "button1",
			position:    value2.NewPosition(0, 0),
		},
		{
			name:        "enter component with position",
			componentID: "button2",
			position:    value2.NewPosition(10, 5),
		},
		{
			name:        "enter component with large coordinates",
			componentID: "panel",
			position:    value2.NewPosition(100, 50),
		},
		{
			name:        "enter component with empty ID",
			componentID: "",
			position:    value2.NewPosition(5, 5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := NewHoverState()
			before := time.Now()

			state.Enter(tt.componentID, tt.position)

			// Check state is now active and hovering
			if !state.IsActive() {
				t.Error("Expected IsActive() to be true after Enter()")
			}

			// IsHovering() should match whether componentID is non-empty
			expectedHovering := tt.componentID != ""
			if state.IsHovering() != expectedHovering {
				t.Errorf("Expected IsHovering() to be %v after Enter(), got %v",
					expectedHovering, state.IsHovering())
			}

			// Check component ID
			if state.ComponentID() != tt.componentID {
				t.Errorf("Expected componentID %s, got %s", tt.componentID, state.ComponentID())
			}

			// Check position
			if !state.Position().Equals(tt.position) {
				t.Errorf("Expected position %v, got %v", tt.position, state.Position())
			}

			// Check timestamp is updated
			if state.LastUpdate().Before(before) {
				t.Error("LastUpdate() should be updated to current time")
			}
		})
	}
}

// TestHoverStateMove tests the Move method.
func TestHoverStateMove(t *testing.T) {
	tests := []struct {
		name            string
		setupActive     bool
		setupComponent  string
		setupPosition   value2.Position
		movePosition    value2.Position
		expectUpdate    bool
		expectComponent string
	}{
		{
			name:            "move within active hover",
			setupActive:     true,
			setupComponent:  "button1",
			setupPosition:   value2.NewPosition(5, 5),
			movePosition:    value2.NewPosition(6, 6),
			expectUpdate:    true,
			expectComponent: "button1",
		},
		{
			name:            "move to same position",
			setupActive:     true,
			setupComponent:  "button1",
			setupPosition:   value2.NewPosition(5, 5),
			movePosition:    value2.NewPosition(5, 5),
			expectUpdate:    true,
			expectComponent: "button1",
		},
		{
			name:            "move when not active (should be ignored)",
			setupActive:     false,
			setupComponent:  "",
			setupPosition:   value2.NewPosition(0, 0),
			movePosition:    value2.NewPosition(10, 10),
			expectUpdate:    false,
			expectComponent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := NewHoverState()

			// Setup state
			if tt.setupActive {
				state.Enter(tt.setupComponent, tt.setupPosition)
			}

			// Record last update time
			time.Sleep(1 * time.Millisecond) // Ensure time difference
			before := state.LastUpdate()

			// Execute move
			state.Move(tt.movePosition)

			// Verify component ID unchanged
			if state.ComponentID() != tt.expectComponent {
				t.Errorf("Expected componentID %s, got %s", tt.expectComponent, state.ComponentID())
			}

			// Verify position updated (or not)
			if tt.expectUpdate {
				if !state.Position().Equals(tt.movePosition) {
					t.Errorf("Expected position %v, got %v", tt.movePosition, state.Position())
				}

				// Timestamp should be updated
				if !state.LastUpdate().After(before) {
					t.Error("LastUpdate() should be updated after Move()")
				}
			} else {
				// Position should not change if not active
				if !state.Position().Equals(tt.setupPosition) {
					t.Errorf("Position should not change when inactive, expected %v, got %v",
						tt.setupPosition, state.Position())
				}
			}

			// Active state should not change
			if state.IsActive() != tt.setupActive {
				t.Errorf("IsActive() should not change after Move(), expected %v, got %v",
					tt.setupActive, state.IsActive())
			}
		})
	}
}

// TestHoverStateLeave tests the Leave method.
func TestHoverStateLeave(t *testing.T) {
	tests := []struct {
		name           string
		setupActive    bool
		setupComponent string
		setupPosition  value2.Position
		leavePosition  value2.Position
		expectUpdate   bool
	}{
		{
			name:           "leave from active hover",
			setupActive:    true,
			setupComponent: "button1",
			setupPosition:  value2.NewPosition(5, 5),
			leavePosition:  value2.NewPosition(10, 10),
			expectUpdate:   true,
		},
		{
			name:           "leave when not active (should be ignored)",
			setupActive:    false,
			setupComponent: "",
			setupPosition:  value2.NewPosition(0, 0),
			leavePosition:  value2.NewPosition(10, 10),
			expectUpdate:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := NewHoverState()

			// Setup state
			if tt.setupActive {
				state.Enter(tt.setupComponent, tt.setupPosition)
			}

			// Record last update time
			time.Sleep(1 * time.Millisecond) // Ensure time difference
			before := state.LastUpdate()

			// Execute leave
			state.Leave(tt.leavePosition)

			if tt.expectUpdate {
				// Component ID should be cleared
				if state.ComponentID() != "" {
					t.Errorf("Expected empty componentID after Leave(), got %s", state.ComponentID())
				}

				// Should no longer be active or hovering
				if state.IsActive() {
					t.Error("Expected IsActive() to be false after Leave()")
				}

				if state.IsHovering() {
					t.Error("Expected IsHovering() to be false after Leave()")
				}

				// Position should be updated
				if !state.Position().Equals(tt.leavePosition) {
					t.Errorf("Expected position %v, got %v", tt.leavePosition, state.Position())
				}

				// Timestamp should be updated
				if !state.LastUpdate().After(before) {
					t.Error("LastUpdate() should be updated after Leave()")
				}
			} else {
				// Nothing should change if not active
				if state.IsActive() != tt.setupActive {
					t.Errorf("IsActive() should not change when leaving inactive state")
				}
			}
		})
	}
}

// TestHoverStateReset tests the Reset method.
func TestHoverStateReset(t *testing.T) {
	state := NewHoverState()

	// Setup some state
	state.Enter("button1", value2.NewPosition(10, 5))

	// Verify setup
	if !state.IsActive() {
		t.Error("Expected active state before reset")
	}

	// Reset
	before := time.Now()
	state.Reset()

	// Verify reset
	if state.ComponentID() != "" {
		t.Errorf("Expected empty componentID after Reset(), got %s", state.ComponentID())
	}

	if state.IsActive() {
		t.Error("Expected IsActive() to be false after Reset()")
	}

	if state.IsHovering() {
		t.Error("Expected IsHovering() to be false after Reset()")
	}

	// Position should be reset to 0,0
	pos := state.Position()
	if pos.X() != 0 || pos.Y() != 0 {
		t.Errorf("Expected position (0,0) after Reset(), got (%d,%d)", pos.X(), pos.Y())
	}

	// Timestamp should be updated
	if state.LastUpdate().Before(before) {
		t.Error("LastUpdate() should be updated after Reset()")
	}
}

// TestHoverStateIsHovering tests the IsHovering method logic.
func TestHoverStateIsHovering(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*HoverState)
		expectHover bool
	}{
		{
			name: "hovering when active with component",
			setup: func(s *HoverState) {
				s.Enter("button1", value2.NewPosition(5, 5))
			},
			expectHover: true,
		},
		{
			name: "not hovering when not active",
			setup: func(s *HoverState) {
				// Leave default state (not active)
			},
			expectHover: false,
		},
		{
			name: "not hovering after leave",
			setup: func(s *HoverState) {
				s.Enter("button1", value2.NewPosition(5, 5))
				s.Leave(value2.NewPosition(10, 10))
			},
			expectHover: false,
		},
		{
			name: "not hovering with empty component ID",
			setup: func(s *HoverState) {
				s.Enter("", value2.NewPosition(5, 5))
			},
			expectHover: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := NewHoverState()
			tt.setup(state)

			if state.IsHovering() != tt.expectHover {
				t.Errorf("Expected IsHovering() to be %v, got %v", tt.expectHover, state.IsHovering())
			}
		})
	}
}

// TestHoverStateEquals tests the Equals method.
func TestHoverStateEquals(t *testing.T) {
	tests := []struct {
		name        string
		state1      *HoverState
		state2      *HoverState
		expectEqual bool
	}{
		{
			name:        "equal new states",
			state1:      NewHoverState(),
			state2:      NewHoverState(),
			expectEqual: true,
		},
		{
			name: "equal after same enter",
			state1: func() *HoverState {
				s := NewHoverState()
				s.Enter("button1", value2.NewPosition(5, 5))
				return s
			}(),
			state2: func() *HoverState {
				s := NewHoverState()
				s.Enter("button1", value2.NewPosition(5, 5))
				return s
			}(),
			expectEqual: true,
		},
		{
			name: "not equal - different component IDs",
			state1: func() *HoverState {
				s := NewHoverState()
				s.Enter("button1", value2.NewPosition(5, 5))
				return s
			}(),
			state2: func() *HoverState {
				s := NewHoverState()
				s.Enter("button2", value2.NewPosition(5, 5))
				return s
			}(),
			expectEqual: false,
		},
		{
			name: "not equal - different positions",
			state1: func() *HoverState {
				s := NewHoverState()
				s.Enter("button1", value2.NewPosition(5, 5))
				return s
			}(),
			state2: func() *HoverState {
				s := NewHoverState()
				s.Enter("button1", value2.NewPosition(6, 6))
				return s
			}(),
			expectEqual: false,
		},
		{
			name: "not equal - different active state",
			state1: func() *HoverState {
				s := NewHoverState()
				s.Enter("button1", value2.NewPosition(5, 5))
				return s
			}(),
			state2: func() *HoverState {
				s := NewHoverState()
				s.Enter("button1", value2.NewPosition(5, 5))
				s.Leave(value2.NewPosition(5, 5))
				return s
			}(),
			expectEqual: false,
		},
		{
			name:        "not equal - nil comparison",
			state1:      NewHoverState(),
			state2:      nil,
			expectEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.state1.Equals(tt.state2)
			if result != tt.expectEqual {
				t.Errorf("Expected Equals() to be %v, got %v", tt.expectEqual, result)
			}
		})
	}
}

// TestHoverStateSequence tests a complete sequence of hover events.
func TestHoverStateSequence(t *testing.T) {
	state := NewHoverState()

	// Initial state
	if state.IsHovering() {
		t.Error("Should not be hovering initially")
	}

	// Enter button1
	state.Enter("button1", value2.NewPosition(10, 5))
	if !state.IsHovering() {
		t.Error("Should be hovering after Enter()")
	}
	if state.ComponentID() != "button1" {
		t.Errorf("Expected button1, got %s", state.ComponentID())
	}

	// Move within button1
	state.Move(value2.NewPosition(11, 5))
	if !state.IsHovering() {
		t.Error("Should still be hovering after Move()")
	}
	if state.ComponentID() != "button1" {
		t.Errorf("Expected button1, got %s", state.ComponentID())
	}
	if state.Position().X() != 11 {
		t.Errorf("Position should be updated to (11,5), got %v", state.Position())
	}

	// Leave button1
	state.Leave(value2.NewPosition(20, 5))
	if state.IsHovering() {
		t.Error("Should not be hovering after Leave()")
	}
	if state.ComponentID() != "" {
		t.Errorf("ComponentID should be empty, got %s", state.ComponentID())
	}

	// Enter button2
	state.Enter("button2", value2.NewPosition(10, 10))
	if !state.IsHovering() {
		t.Error("Should be hovering after entering button2")
	}
	if state.ComponentID() != "button2" {
		t.Errorf("Expected button2, got %s", state.ComponentID())
	}

	// Reset
	state.Reset()
	if state.IsHovering() {
		t.Error("Should not be hovering after Reset()")
	}
	if state.ComponentID() != "" {
		t.Errorf("ComponentID should be empty after Reset(), got %s", state.ComponentID())
	}
}
