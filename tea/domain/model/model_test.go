package model

import (
	"fmt"
	"testing"
)

// TestModel is a simple test implementation of the Model interface.
// It increments on "+" key, decrements on "-" key, and quits on "q".
type TestModel struct {
	value int
	ready bool
}

// Init initializes the model with a ready flag.
func (m TestModel) Init() Cmd {
	return func() Msg {
		return testReadyMsg{}
	}
}

// Update handles messages for the test model.
func (m TestModel) Update(msg Msg) (Model[TestModel], Cmd) {
	switch msg := msg.(type) {
	case testReadyMsg:
		m.ready = true
		return m, nil
	case KeyMsg:
		switch msg.String() {
		case "+":
			m.value++
			return m, testIncrementedCmd()
		case "-":
			m.value--
			return m, nil
		case "q":
			return m, testQuitCmd()
		}
	}
	return m, nil
}

// View renders the model to a string.
func (m TestModel) View() string {
	return fmt.Sprintf("Value: %d", m.value)
}

// testReadyMsg is sent by Init to signal ready state.
type testReadyMsg struct{}

func (t testReadyMsg) String() string { return "ready" }

// testIncrementedMsg is sent when value is incremented.
type testIncrementedMsg struct{}

func (t testIncrementedMsg) String() string { return "incremented" }

// testIncrementedCmd returns a command that sends incrementedMsg.
func testIncrementedCmd() Cmd {
	return func() Msg {
		return testIncrementedMsg{}
	}
}

// testQuitCmd returns a command that sends QuitMsg.
func testQuitCmd() Cmd {
	return func() Msg {
		return QuitMsg{}
	}
}

// TestModel_Interface verifies TestModel implements Model interface.
func TestModel_Interface(t *testing.T) {
	var _ Model[TestModel] = TestModel{}
}

// TestModel_Init verifies Init method works correctly.
func TestModel_Init(t *testing.T) {
	m := TestModel{}

	cmd := m.Init()
	if cmd == nil {
		t.Fatal("Init() returned nil, expected command")
	}

	msg := cmd()
	if _, ok := msg.(testReadyMsg); !ok {
		t.Errorf("Init command sent %T, expected testReadyMsg", msg)
	}
}

// TestModel_Update_KeyMessages verifies Update handles key messages.
func TestModel_Update_KeyMessages(t *testing.T) {
	tests := []struct {
		name          string
		initialValue  int
		key           KeyMsg
		expectedValue int
		expectCmd     bool
		expectedMsg   string
	}{
		{
			name:          "increment on +",
			initialValue:  0,
			key:           KeyMsg{Type: KeyRune, Rune: '+'},
			expectedValue: 1,
			expectCmd:     true,
			expectedMsg:   "incremented",
		},
		{
			name:          "increment from 5",
			initialValue:  5,
			key:           KeyMsg{Type: KeyRune, Rune: '+'},
			expectedValue: 6,
			expectCmd:     true,
			expectedMsg:   "incremented",
		},
		{
			name:          "decrement on -",
			initialValue:  10,
			key:           KeyMsg{Type: KeyRune, Rune: '-'},
			expectedValue: 9,
			expectCmd:     false,
		},
		{
			name:          "decrement to negative",
			initialValue:  0,
			key:           KeyMsg{Type: KeyRune, Rune: '-'},
			expectedValue: -1,
			expectCmd:     false,
		},
		{
			name:          "quit on q",
			initialValue:  5,
			key:           KeyMsg{Type: KeyRune, Rune: 'q'},
			expectedValue: 5,
			expectCmd:     true,
			expectedMsg:   "quit",
		},
		{
			name:          "ignore other keys",
			initialValue:  5,
			key:           KeyMsg{Type: KeyRune, Rune: 'x'},
			expectedValue: 5,
			expectCmd:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := TestModel{value: tt.initialValue}

			newModel, cmd := m.Update(tt.key)

			// Check value
			concrete, ok := newModel.(TestModel)
			if !ok {
				t.Fatalf("Update returned %T, expected TestModel", newModel)
			}
			if concrete.value != tt.expectedValue {
				t.Errorf("value = %d, want %d", concrete.value, tt.expectedValue)
			}

			// Check command
			if tt.expectCmd {
				if cmd == nil {
					t.Error("expected command, got nil")
				} else {
					msg := cmd()
					msgStr := fmt.Sprintf("%v", msg)
					if msgStr != tt.expectedMsg {
						t.Errorf("command sent %q, want %q", msgStr, tt.expectedMsg)
					}
				}
			} else {
				if cmd != nil {
					t.Errorf("expected no command, got %T", cmd)
				}
			}
		})
	}
}

// TestModel_Update_ReadyMessage verifies Update handles ready message.
func TestModel_Update_ReadyMessage(t *testing.T) {
	m := TestModel{value: 42}

	newModel, cmd := m.Update(testReadyMsg{})

	concrete := newModel.(TestModel)
	if !concrete.ready {
		t.Error("ready flag not set after testReadyMsg")
	}
	if concrete.value != 42 {
		t.Errorf("value changed to %d, expected unchanged 42", concrete.value)
	}
	if cmd != nil {
		t.Errorf("expected no command, got %T", cmd)
	}
}

// TestModel_Update_UnknownMessage verifies Update ignores unknown messages.
func TestModel_Update_UnknownMessage(t *testing.T) {
	m := TestModel{value: 10}

	type unknownMsg struct{}

	newModel, cmd := m.Update(unknownMsg{})

	concrete := newModel.(TestModel)
	if concrete.value != 10 {
		t.Errorf("value = %d, expected unchanged 10", concrete.value)
	}
	if cmd != nil {
		t.Errorf("expected no command, got %T", cmd)
	}
}

// TestModel_View verifies View renders correctly.
func TestModel_View(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected string
	}{
		{
			name:     "zero value",
			value:    0,
			expected: "Value: 0",
		},
		{
			name:     "positive value",
			value:    42,
			expected: "Value: 42",
		},
		{
			name:     "negative value",
			value:    -10,
			expected: "Value: -10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := TestModel{value: tt.value}

			view := m.View()

			if view != tt.expected {
				t.Errorf("View() = %q, want %q", view, tt.expected)
			}
		})
	}
}

// TestModel_Update_ReturnsCorrectType verifies Update returns Model[TestModel].
func TestModel_Update_ReturnsCorrectType(t *testing.T) {
	m := TestModel{value: 5}

	newModel, _ := m.Update(KeyMsg{Type: KeyRune, Rune: '+'})

	// Should be able to assign to Model[TestModel]
	var _ Model[TestModel] = newModel

	// Should also be able to type assert to concrete type
	concrete, ok := newModel.(TestModel)
	if !ok {
		t.Errorf("Update returned %T, expected TestModel", newModel)
	}
	if concrete.value != 6 {
		t.Errorf("value = %d, want 6", concrete.value)
	}
}

// TestModel_Immutability verifies Update returns new instance (immutability).
func TestModel_Immutability(t *testing.T) {
	m := TestModel{value: 5}
	original := m

	newModel, _ := m.Update(KeyMsg{Type: KeyRune, Rune: '+'})

	// Original should be unchanged
	if original.value != 5 {
		t.Errorf("original.value = %d, want 5 (mutated!)", original.value)
	}

	// New model should have updated value
	concrete := newModel.(TestModel)
	if concrete.value != 6 {
		t.Errorf("newModel.value = %d, want 6", concrete.value)
	}
}

// TestModel_CommandExecution verifies commands execute correctly.
func TestModel_CommandExecution(t *testing.T) {
	m := TestModel{value: 0}

	// Update with "+" should return increment command
	newModel, cmd := m.Update(KeyMsg{Type: KeyRune, Rune: '+'})

	if cmd == nil {
		t.Fatal("expected command, got nil")
	}

	// Execute command
	msg := cmd()

	// Should receive incrementedMsg
	if _, ok := msg.(testIncrementedMsg); !ok {
		t.Errorf("command sent %T, expected testIncrementedMsg", msg)
	}

	// Update with that message
	finalModel, _ := newModel.Update(msg)

	// Model should still have incremented value
	concrete := finalModel.(TestModel)
	if concrete.value != 1 {
		t.Errorf("value = %d, want 1", concrete.value)
	}
}
