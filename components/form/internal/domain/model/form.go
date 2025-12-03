// Package model contains the domain model for the form component.
package model

import (
	"github.com/phoenix-tui/phoenix/components/form/internal/domain/value"
)

// Form represents the form domain model with fields and validation state.
type Form struct {
	title     string
	fields    []*value.Field
	focused   int
	submitted bool
}

// New creates a new Form with the given title.
func New(title string) *Form {
	return &Form{
		title:     title,
		fields:    []*value.Field{},
		focused:   0,
		submitted: false,
	}
}

// Title returns the form title.
func (f *Form) Title() string {
	return f.title
}

// Fields returns all fields.
func (f *Form) Fields() []*value.Field {
	return f.fields
}

// FocusedIndex returns the index of the focused field.
func (f *Form) FocusedIndex() int {
	return f.focused
}

// Submitted returns whether the form has been submitted.
func (f *Form) Submitted() bool {
	return f.submitted
}

// WithFields returns a new form with the specified fields.
func (f *Form) WithFields(fields []*value.Field) *Form {
	focused := f.focused
	if focused >= len(fields) {
		focused = 0
	}

	return &Form{
		title:     f.title,
		fields:    fields,
		focused:   focused,
		submitted: f.submitted,
	}
}

// FocusedField returns the currently focused field, or nil if no fields.
func (f *Form) FocusedField() *value.Field {
	if len(f.fields) == 0 || f.focused >= len(f.fields) {
		return nil
	}
	return f.fields[f.focused]
}

// MoveFocusNext moves focus to the next field (wraps around).
func (f *Form) MoveFocusNext() *Form {
	if len(f.fields) == 0 {
		return f
	}

	newFocused := (f.focused + 1) % len(f.fields)

	// Mark current field as touched when leaving it
	newFields := make([]*value.Field, len(f.fields))
	copy(newFields, f.fields)
	if f.focused < len(newFields) {
		newFields[f.focused] = newFields[f.focused].MarkTouched()
	}

	return &Form{
		title:     f.title,
		fields:    newFields,
		focused:   newFocused,
		submitted: f.submitted,
	}
}

// MoveFocusPrev moves focus to the previous field (wraps around).
func (f *Form) MoveFocusPrev() *Form {
	if len(f.fields) == 0 {
		return f
	}

	newFocused := f.focused - 1
	if newFocused < 0 {
		newFocused = len(f.fields) - 1
	}

	// Mark current field as touched when leaving it
	newFields := make([]*value.Field, len(f.fields))
	copy(newFields, f.fields)
	if f.focused < len(newFields) {
		newFields[f.focused] = newFields[f.focused].MarkTouched()
	}

	return &Form{
		title:     f.title,
		fields:    newFields,
		focused:   newFocused,
		submitted: f.submitted,
	}
}

// UpdateField updates the field at the given index with a new model.
func (f *Form) UpdateField(index int, model interface{}) *Form {
	if index < 0 || index >= len(f.fields) {
		return f
	}

	// Mark field as dirty
	newFields := make([]*value.Field, len(f.fields))
	copy(newFields, f.fields)

	// Type assert to tea.Model if possible
	if teaModel, ok := model.(interface{ View() string }); ok {
		newFields[index] = newFields[index].WithModel(teaModel).MarkDirty()
	}

	return &Form{
		title:     f.title,
		fields:    newFields,
		focused:   f.focused,
		submitted: f.submitted,
	}
}

// ValidateField validates a single field by index.
func (f *Form) ValidateField(index int, fieldValue interface{}) *Form {
	if index < 0 || index >= len(f.fields) {
		return f
	}

	field := f.fields[index]
	var errors []error

	// Run all validators
	for _, validator := range field.Validators() {
		if err := validator(fieldValue); err != nil {
			errors = append(errors, err)
		}
	}

	// Update field with errors
	newFields := make([]*value.Field, len(f.fields))
	copy(newFields, f.fields)
	newFields[index] = newFields[index].WithErrors(errors)

	return &Form{
		title:     f.title,
		fields:    newFields,
		focused:   f.focused,
		submitted: f.submitted,
	}
}

// ValidateAll validates all fields.
func (f *Form) ValidateAll(fieldValues map[string]interface{}) *Form {
	newFields := make([]*value.Field, len(f.fields))
	copy(newFields, f.fields)

	for i, field := range newFields {
		var errors []error

		// Get value for this field
		fieldValue, ok := fieldValues[field.Name()]
		if !ok {
			continue
		}

		// Run all validators
		for _, validator := range field.Validators() {
			if err := validator(fieldValue); err != nil {
				errors = append(errors, err)
			}
		}

		newFields[i] = newFields[i].WithErrors(errors)
	}

	return &Form{
		title:     f.title,
		fields:    newFields,
		focused:   f.focused,
		submitted: f.submitted,
	}
}

// IsValid returns whether all fields are valid.
func (f *Form) IsValid() bool {
	for _, field := range f.fields {
		if !field.IsValid() {
			return false
		}
	}
	return true
}

// Submit marks the form as submitted.
func (f *Form) Submit() *Form {
	return &Form{
		title:     f.title,
		fields:    f.fields,
		focused:   f.focused,
		submitted: true,
	}
}

// Reset resets the form to initial state.
func (f *Form) Reset() *Form {
	// Clear all errors and reset touched/dirty state
	newFields := make([]*value.Field, len(f.fields))
	for i, field := range f.fields {
		newFields[i] = value.NewField(field.Name(), field.Label(), field.Model()).
			WithValidators(field.Validators()...)
	}

	return &Form{
		title:     f.title,
		fields:    newFields,
		focused:   0,
		submitted: false,
	}
}

// Errors returns a map of field name → error messages.
func (f *Form) Errors() map[string][]string {
	result := make(map[string][]string)

	for _, field := range f.fields {
		if len(field.Errors()) > 0 {
			messages := make([]string, len(field.Errors()))
			for i, err := range field.Errors() {
				messages[i] = err.Error()
			}
			result[field.Name()] = messages
		}
	}

	return result
}
