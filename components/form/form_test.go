package form_test

import (
	"errors"
	"testing"

	"github.com/phoenix-tui/phoenix/components/form"
	"github.com/phoenix-tui/phoenix/components/form/internal/domain/value"
	tea "github.com/phoenix-tui/phoenix/tea"
)

// mockModel is a simple mock tea.Model for testing.
type mockModel struct {
	content string
}

func (m *mockModel) Init() tea.Cmd { return nil }
func (m *mockModel) Update(_ tea.Msg) (*mockModel, tea.Cmd) {
	return m, nil
}
func (m *mockModel) View() string { return m.content }

// mockValuer is a mock model that implements Value() method.
type mockValuer struct {
	mockModel
	val string
}

func (m *mockValuer) Value() string {
	return m.val
}

func TestNew(t *testing.T) {
	f := form.New("Test Form")

	values := f.Values()
	if len(values) != 0 {
		t.Errorf("New form has %d fields, want 0", len(values))
	}

	if f.IsValid() != true {
		t.Error("New form should be valid")
	}

	if f.IsDirty() {
		t.Error("New form should not be dirty")
	}
}

func TestField(t *testing.T) {
	model1 := &mockModel{content: "field1"}
	model2 := &mockModel{content: "field2"}

	f := form.New("Test").
		Field("name", "Name", model1).
		Field("email", "Email", model2)

	values := f.Values()
	if len(values) != 2 {
		t.Errorf("Form has %d fields, want 2", len(values))
	}

	if values["name"] != model1 {
		t.Error("Field 'name' not stored correctly")
	}

	if values["email"] != model2 {
		t.Error("Field 'email' not stored correctly")
	}
}

func TestValue(t *testing.T) {
	model := &mockModel{content: "test"}

	f := form.New("Test").
		Field("name", "Name", model)

	val := f.Value("name")
	if val != model {
		t.Error("Value('name') did not return correct model")
	}

	val = f.Value("nonexistent")
	if val != nil {
		t.Error("Value() for nonexistent field should return nil")
	}
}

func TestIsValidWithValidators(t *testing.T) {
	model := &mockValuer{val: "test@example.com"}

	validator := func(val interface{}) error {
		if str, ok := val.(string); ok && str == "" {
			return errors.New("required")
		}
		return nil
	}

	f := form.New("Test").
		Field("email", "Email", model, validator)

	// Initially valid (no validation run yet)
	if !f.IsValid() {
		t.Error("Form should be valid initially")
	}
}

func TestIsDirty(t *testing.T) {
	model := &mockModel{content: "test"}

	f := form.New("Test").
		Field("name", "Name", model)

	if f.IsDirty() {
		t.Error("New form should not be dirty")
	}
}

func TestErrors(t *testing.T) {
	model := &mockModel{content: "test"}

	f := form.New("Test").
		Field("name", "Name", model)

	errors := f.Errors()
	if len(errors) != 0 {
		t.Errorf("Form has %d error entries, want 0", len(errors))
	}
}

func TestInit(t *testing.T) {
	model := &mockModel{content: "test"}

	f := form.New("Test").
		Field("name", "Name", model)

	cmd := f.Init()
	// Mock model returns nil for Init, so form should return nil
	if cmd != nil {
		t.Error("Init() should return nil when field models have no init commands")
	}
}

func TestUpdateTabNavigation(t *testing.T) {
	model1 := &mockModel{content: "field1"}
	model2 := &mockModel{content: "field2"}

	f := form.New("Test").
		Field("name", "Name", model1).
		Field("email", "Email", model2)

	// Press Tab to move to next field
	msg := tea.KeyMsg{Type: tea.KeyTab}
	newForm, _ := f.Update(msg)

	// Focus should have moved (we can't directly check focus, but form should be updated)
	if newForm == f {
		t.Error("Update(Tab) should return new form instance")
	}
}

func TestUpdateShiftTab(t *testing.T) {
	model1 := &mockModel{content: "field1"}
	model2 := &mockModel{content: "field2"}

	f := form.New("Test").
		Field("name", "Name", model1).
		Field("email", "Email", model2)

	// Press Shift+Tab to move to previous field
	msg := tea.KeyMsg{Type: tea.KeyTab, Shift: true}
	newForm, _ := f.Update(msg)

	if newForm == f {
		t.Error("Update(Shift+Tab) should return new form instance")
	}
}

func TestUpdateSubmit(t *testing.T) {
	model := &mockValuer{val: "test@example.com"}

	f := form.New("Test").
		Field("email", "Email", model, value.Email())

	// Press Enter to submit
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newForm, cmd := f.Update(msg)

	if newForm == f {
		t.Error("Update(Enter) should return new form instance")
	}

	// Should receive submit command if validation passes
	if cmd == nil {
		t.Error("Update(Enter) should return submit command when valid")
	}

	// Execute command to get SubmitMsg
	if cmd != nil {
		resultMsg := cmd()
		if _, ok := resultMsg.(form.SubmitMsg); !ok {
			t.Error("Submit command should return SubmitMsg")
		}
	}
}

func TestUpdateSubmitInvalid(t *testing.T) {
	model := &mockValuer{val: ""} // Empty value

	f := form.New("Test").
		Field("email", "Email", model, value.Required())

	// Press Enter to submit
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newForm, cmd := f.Update(msg)

	if newForm == f {
		t.Error("Update(Enter) should return new form instance")
	}

	// Should NOT receive submit command if validation fails
	if cmd != nil {
		resultMsg := cmd()
		if _, ok := resultMsg.(form.SubmitMsg); ok {
			t.Error("Submit command should not be sent when form is invalid")
		}
	}

	// Check that form has errors
	if len(newForm.Errors()) == 0 {
		t.Error("Form should have validation errors after submit attempt")
	}
}

func TestUpdateQuit(t *testing.T) {
	model := &mockModel{content: "test"}

	f := form.New("Test").
		Field("name", "Name", model)

	// Press Ctrl+C to quit
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := f.Update(msg)

	if cmd == nil {
		t.Error("Update(Ctrl+C) should return quit command")
	}

	// Execute command to verify it's a quit
	if cmd != nil {
		resultMsg := cmd()
		if _, ok := resultMsg.(tea.QuitMsg); !ok {
			t.Error("Quit command should return tea.QuitMsg")
		}
	}
}

func TestView(t *testing.T) {
	model := &mockModel{content: "test"}

	f := form.New("User Form").
		Field("name", "Name", model)

	view := f.View()

	// View should contain title
	if view == "" {
		t.Error("View() returned empty string")
	}

	// Should contain title
	if !containsSubstring(view, "User Form") {
		t.Error("View() should contain form title")
	}

	// Should contain help text
	if !containsSubstring(view, "Tab") {
		t.Error("View() should contain help text with 'Tab'")
	}
}

func TestMultipleValidators(t *testing.T) {
	model := &mockValuer{val: "a"} // Too short

	f := form.New("Test").
		Field("name", "Name", model,
			value.Required(),
			value.MinLength(3))

	// Submit to trigger validation
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newForm, _ := f.Update(msg)

	// Should have validation errors
	errors := newForm.Errors()
	if len(errors["name"]) == 0 {
		t.Error("Field should have MinLength validation error")
	}
}

func TestFormImmutability(t *testing.T) {
	model := &mockModel{content: "test"}

	f := form.New("Test").
		Field("name", "Name", model)

	// Apply multiple updates
	msg := tea.KeyMsg{Type: tea.KeyTab}
	f2, _ := f.Update(msg)
	f3, _ := f2.Update(msg)

	// Each should be a different instance
	if f == f2 || f2 == f3 || f == f3 {
		t.Error("Form updates should return new instances (immutability)")
	}

	// Original should be unchanged
	if len(f.Values()) != 1 {
		t.Error("Original form was mutated")
	}
}

// Helper function to check if a string contains a substring.
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
