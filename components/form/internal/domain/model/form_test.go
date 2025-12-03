package model_test

import (
	"errors"
	"testing"

	"github.com/phoenix-tui/phoenix/components/form/internal/domain/model"
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

func TestNew(t *testing.T) {
	form := model.New("Test Form")

	if form.Title() != "Test Form" {
		t.Errorf("Title() = %q, want %q", form.Title(), "Test Form")
	}
	if len(form.Fields()) != 0 {
		t.Errorf("Fields() = %d fields, want 0", len(form.Fields()))
	}
	if form.FocusedIndex() != 0 {
		t.Errorf("FocusedIndex() = %d, want 0", form.FocusedIndex())
	}
	if form.Submitted() {
		t.Error("Submitted() = true, want false")
	}
}

func TestWithFields(t *testing.T) {
	form := model.New("Test")

	field1 := value.NewField("name", "Name", &mockModel{content: "field1"})
	field2 := value.NewField("email", "Email", &mockModel{content: "field2"})

	newForm := form.WithFields([]*value.Field{field1, field2})

	if len(newForm.Fields()) != 2 {
		t.Errorf("Fields() = %d fields, want 2", len(newForm.Fields()))
	}

	// Original form should be unchanged
	if len(form.Fields()) != 0 {
		t.Error("Original form was mutated")
	}
}

func TestFocusedField(t *testing.T) {
	field1 := value.NewField("name", "Name", &mockModel{content: "field1"})
	field2 := value.NewField("email", "Email", &mockModel{content: "field2"})

	form := model.New("Test").WithFields([]*value.Field{field1, field2})

	focused := form.FocusedField()
	if focused == nil {
		t.Fatal("FocusedField() returned nil")
	}
	if focused.Name() != "name" {
		t.Errorf("FocusedField().Name() = %q, want %q", focused.Name(), "name")
	}
}

func TestFocusedFieldEmptyForm(t *testing.T) {
	form := model.New("Test")

	focused := form.FocusedField()
	if focused != nil {
		t.Error("FocusedField() should return nil for empty form")
	}
}

func TestMoveFocusNext(t *testing.T) {
	field1 := value.NewField("name", "Name", &mockModel{content: "field1"})
	field2 := value.NewField("email", "Email", &mockModel{content: "field2"})
	field3 := value.NewField("age", "Age", &mockModel{content: "field3"})

	form := model.New("Test").WithFields([]*value.Field{field1, field2, field3})

	// Move to field 2
	form2 := form.MoveFocusNext()
	if form2.FocusedIndex() != 1 {
		t.Errorf("After first MoveFocusNext(), FocusedIndex() = %d, want 1", form2.FocusedIndex())
	}

	// Move to field 3
	form3 := form2.MoveFocusNext()
	if form3.FocusedIndex() != 2 {
		t.Errorf("After second MoveFocusNext(), FocusedIndex() = %d, want 2", form3.FocusedIndex())
	}

	// Wrap around to field 1
	form4 := form3.MoveFocusNext()
	if form4.FocusedIndex() != 0 {
		t.Errorf("After wrap-around, FocusedIndex() = %d, want 0", form4.FocusedIndex())
	}

	// Check that previous field was marked as touched
	if !form2.Fields()[0].Touched() {
		t.Error("Field 0 should be marked as touched after leaving it")
	}
}

func TestMoveFocusPrev(t *testing.T) {
	field1 := value.NewField("name", "Name", &mockModel{content: "field1"})
	field2 := value.NewField("email", "Email", &mockModel{content: "field2"})
	field3 := value.NewField("age", "Age", &mockModel{content: "field3"})

	form := model.New("Test").WithFields([]*value.Field{field1, field2, field3})

	// Move backwards (should wrap to last field)
	form2 := form.MoveFocusPrev()
	if form2.FocusedIndex() != 2 {
		t.Errorf("After MoveFocusPrev() from 0, FocusedIndex() = %d, want 2", form2.FocusedIndex())
	}

	// Move backwards again
	form3 := form2.MoveFocusPrev()
	if form3.FocusedIndex() != 1 {
		t.Errorf("After second MoveFocusPrev(), FocusedIndex() = %d, want 1", form3.FocusedIndex())
	}
}

func TestUpdateField(t *testing.T) {
	field1 := value.NewField("name", "Name", &mockModel{content: "field1"})
	field2 := value.NewField("email", "Email", &mockModel{content: "field2"})

	form := model.New("Test").WithFields([]*value.Field{field1, field2})

	newModel := &mockModel{content: "updated"}
	newForm := form.UpdateField(1, newModel)

	// Check that field was updated
	updatedModel := newForm.Fields()[1].Model()
	if mock, ok := updatedModel.(*mockModel); !ok || mock.content != "updated" {
		t.Error("Field model was not updated correctly")
	}

	// Check that field was marked as dirty
	if !newForm.Fields()[1].Dirty() {
		t.Error("Field should be marked as dirty after update")
	}

	// Original form should be unchanged
	originalModel := form.Fields()[1].Model()
	if mock, ok := originalModel.(*mockModel); !ok || mock.content != "field2" {
		t.Error("Original form was mutated")
	}
}

func TestValidateField(t *testing.T) {
	validator := func(val interface{}) error {
		if str, ok := val.(string); ok && str == "invalid" {
			return errors.New("validation error")
		}
		return nil
	}

	field := value.NewField("name", "Name", &mockModel{content: "field"}).
		WithValidators(validator)

	form := model.New("Test").WithFields([]*value.Field{field})

	// Validate with valid value
	newForm := form.ValidateField(0, "valid")
	if !newForm.Fields()[0].IsValid() {
		t.Error("Field should be valid after validation with valid value")
	}

	// Validate with invalid value
	newForm2 := form.ValidateField(0, "invalid")
	if newForm2.Fields()[0].IsValid() {
		t.Error("Field should be invalid after validation with invalid value")
	}
	if len(newForm2.Fields()[0].Errors()) == 0 {
		t.Error("Field should have errors after validation with invalid value")
	}
}

func TestValidateAll(t *testing.T) {
	validator1 := func(val interface{}) error {
		if str, ok := val.(string); ok && str == "" {
			return errors.New("required")
		}
		return nil
	}

	validator2 := func(val interface{}) error {
		if str, ok := val.(string); ok && len(str) < 3 {
			return errors.New("too short")
		}
		return nil
	}

	field1 := value.NewField("name", "Name", &mockModel{}).WithValidators(validator1)
	field2 := value.NewField("email", "Email", &mockModel{}).WithValidators(validator2)

	form := model.New("Test").WithFields([]*value.Field{field1, field2})

	fieldValues := map[string]interface{}{
		"name":  "",    // Invalid (empty)
		"email": "ab",  // Invalid (too short)
	}

	newForm := form.ValidateAll(fieldValues)

	// Both fields should have errors
	if newForm.Fields()[0].IsValid() {
		t.Error("Field 'name' should be invalid")
	}
	if newForm.Fields()[1].IsValid() {
		t.Error("Field 'email' should be invalid")
	}

	// Form should not be valid
	if newForm.IsValid() {
		t.Error("Form should not be valid when fields have errors")
	}
}

func TestIsValid(t *testing.T) {
	field1 := value.NewField("name", "Name", &mockModel{})
	field2 := value.NewField("email", "Email", &mockModel{})

	form := model.New("Test").WithFields([]*value.Field{field1, field2})

	if !form.IsValid() {
		t.Error("Form with no errors should be valid")
	}

	// Add error to one field
	field1WithError := field1.WithErrors([]error{errors.New("error")})
	formWithError := form.WithFields([]*value.Field{field1WithError, field2})

	if formWithError.IsValid() {
		t.Error("Form with field errors should not be valid")
	}
}

func TestSubmit(t *testing.T) {
	form := model.New("Test")

	newForm := form.Submit()

	if !newForm.Submitted() {
		t.Error("Submit() did not mark form as submitted")
	}
	if form.Submitted() {
		t.Error("Original form was mutated")
	}
}

func TestReset(t *testing.T) {
	validator := func(_ interface{}) error { return errors.New("error") }

	field1 := value.NewField("name", "Name", &mockModel{}).
		WithValidators(validator).
		MarkTouched().
		MarkDirty().
		WithErrors([]error{errors.New("error")})

	field2 := value.NewField("email", "Email", &mockModel{})

	form := model.New("Test").
		WithFields([]*value.Field{field1, field2}).
		Submit()

	newForm := form.Reset()

	// Check that state was reset
	if newForm.Submitted() {
		t.Error("Reset() did not clear submitted flag")
	}
	if newForm.FocusedIndex() != 0 {
		t.Error("Reset() did not reset focus to 0")
	}

	// Check that fields were reset
	if newForm.Fields()[0].Touched() {
		t.Error("Reset() did not clear touched flag on fields")
	}
	if newForm.Fields()[0].Dirty() {
		t.Error("Reset() did not clear dirty flag on fields")
	}
	if len(newForm.Fields()[0].Errors()) > 0 {
		t.Error("Reset() did not clear errors on fields")
	}

	// Check that validators are preserved
	if len(newForm.Fields()[0].Validators()) == 0 {
		t.Error("Reset() removed validators from fields")
	}
}

func TestErrors(t *testing.T) {
	field1 := value.NewField("name", "Name", &mockModel{}).
		WithErrors([]error{
			errors.New("required"),
			errors.New("too short"),
		})

	field2 := value.NewField("email", "Email", &mockModel{}).
		WithErrors([]error{errors.New("invalid format")})

	field3 := value.NewField("age", "Age", &mockModel{}) // No errors

	form := model.New("Test").WithFields([]*value.Field{field1, field2, field3})

	errorMap := form.Errors()

	if len(errorMap) != 2 {
		t.Errorf("Errors() returned %d fields with errors, want 2", len(errorMap))
	}

	nameErrors, ok := errorMap["name"]
	if !ok {
		t.Error("Errors() did not include 'name' field")
	} else if len(nameErrors) != 2 {
		t.Errorf("'name' has %d errors, want 2", len(nameErrors))
	}

	emailErrors, ok := errorMap["email"]
	if !ok {
		t.Error("Errors() did not include 'email' field")
	} else if len(emailErrors) != 1 {
		t.Errorf("'email' has %d errors, want 1", len(emailErrors))
	}

	if _, ok := errorMap["age"]; ok {
		t.Error("Errors() should not include fields without errors")
	}
}
