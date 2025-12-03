package value_test

import (
	"errors"
	"testing"

	"github.com/phoenix-tui/phoenix/components/form/internal/domain/value"
)

// mockModel is a simple mock tea.Model for testing.
type mockModel struct {
	content string
}

func (m *mockModel) Init() interface{} { return nil }
func (m *mockModel) Update(_ interface{}) (interface{}, interface{}) {
	return m, nil
}
func (m *mockModel) View() string { return m.content }

func TestNewField(t *testing.T) {
	model := &mockModel{content: "test"}
	field := value.NewField("username", "Username", model)

	if field.Name() != "username" {
		t.Errorf("Name() = %q, want %q", field.Name(), "username")
	}
	if field.Label() != "Username" {
		t.Errorf("Label() = %q, want %q", field.Label(), "Username")
	}
	if field.Model() != model {
		t.Error("Model() did not return the same model instance")
	}
	if len(field.Validators()) != 0 {
		t.Errorf("Validators() = %d validators, want 0", len(field.Validators()))
	}
	if field.Touched() {
		t.Error("Touched() = true, want false")
	}
	if field.Dirty() {
		t.Error("Dirty() = true, want false")
	}
	if !field.IsValid() {
		t.Error("IsValid() = false, want true (no errors initially)")
	}
}

func TestFieldWithValidators(t *testing.T) {
	model := &mockModel{content: "test"}
	field := value.NewField("email", "Email", model)

	validator := value.Required()
	newField := field.WithValidators(validator)

	if len(newField.Validators()) != 1 {
		t.Errorf("Validators() = %d validators, want 1", len(newField.Validators()))
	}

	// Original field should be unchanged
	if len(field.Validators()) != 0 {
		t.Error("Original field was mutated")
	}
}

func TestFieldMarkTouched(t *testing.T) {
	model := &mockModel{content: "test"}
	field := value.NewField("name", "Name", model)

	newField := field.MarkTouched()

	if !newField.Touched() {
		t.Error("MarkTouched() did not set Touched to true")
	}
	if field.Touched() {
		t.Error("Original field was mutated")
	}
}

func TestFieldMarkDirty(t *testing.T) {
	model := &mockModel{content: "test"}
	field := value.NewField("name", "Name", model)

	newField := field.MarkDirty()

	if !newField.Dirty() {
		t.Error("MarkDirty() did not set Dirty to true")
	}
	if field.Dirty() {
		t.Error("Original field was mutated")
	}
}

func TestFieldWithErrors(t *testing.T) {
	model := &mockModel{content: "test"}
	field := value.NewField("email", "Email", model)

	errors := []error{
		errors.New("invalid format"),
		errors.New("required"),
	}

	newField := field.WithErrors(errors)

	if len(newField.Errors()) != 2 {
		t.Errorf("Errors() = %d errors, want 2", len(newField.Errors()))
	}
	if newField.IsValid() {
		t.Error("IsValid() = true, want false (field has errors)")
	}

	// Original field should be unchanged
	if len(field.Errors()) != 0 {
		t.Error("Original field was mutated")
	}
}

func TestFieldClearErrors(t *testing.T) {
	model := &mockModel{content: "test"}
	field := value.NewField("email", "Email", model).
		WithErrors([]error{errors.New("test error")})

	newField := field.ClearErrors()

	if len(newField.Errors()) != 0 {
		t.Errorf("ClearErrors() left %d errors, want 0", len(newField.Errors()))
	}
	if !newField.IsValid() {
		t.Error("IsValid() = false after ClearErrors(), want true")
	}
}

func TestFieldWithModel(t *testing.T) {
	model1 := &mockModel{content: "model1"}
	model2 := &mockModel{content: "model2"}

	field := value.NewField("test", "Test", model1)
	newField := field.WithModel(model2)

	if newField.Model() != model2 {
		t.Error("WithModel() did not update the model")
	}
	if field.Model() != model1 {
		t.Error("Original field was mutated")
	}
}

func TestFieldImmutability(t *testing.T) {
	model := &mockModel{content: "test"}
	field := value.NewField("test", "Test", model)

	// Apply multiple transformations
	field2 := field.MarkTouched()
	field3 := field2.MarkDirty()
	field4 := field3.WithErrors([]error{errors.New("error")})

	// Original should be unchanged
	if field.Touched() || field.Dirty() || len(field.Errors()) > 0 {
		t.Error("Original field was mutated during transformations")
	}

	// Each transformation should build on previous
	if !field2.Touched() {
		t.Error("field2 should be touched")
	}
	if !field3.Dirty() {
		t.Error("field3 should be dirty")
	}
	if len(field4.Errors()) == 0 {
		t.Error("field4 should have errors")
	}
}
