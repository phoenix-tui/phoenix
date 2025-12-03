// Package value contains value objects for the form component.
package value

// FieldModel is a minimal interface for field models (anything with a View method).
type FieldModel interface {
	View() string
}

// Field represents a form field that can be any FieldModel.
type Field struct {
	name       string
	label      string
	model      FieldModel
	validators []Validator
	touched    bool
	dirty      bool
	errors     []error
}

// NewField creates a new field with the given name, label, and model.
func NewField(name, label string, model FieldModel) *Field {
	return &Field{
		name:       name,
		label:      label,
		model:      model,
		validators: []Validator{},
		touched:    false,
		dirty:      false,
		errors:     []error{},
	}
}

// Name returns the field name.
func (f *Field) Name() string {
	return f.name
}

// Label returns the field label.
func (f *Field) Label() string {
	return f.label
}

// Model returns the underlying FieldModel.
func (f *Field) Model() FieldModel {
	return f.model
}

// Validators returns the field validators.
func (f *Field) Validators() []Validator {
	return f.validators
}

// Touched returns whether the field has been focused.
func (f *Field) Touched() bool {
	return f.touched
}

// Dirty returns whether the field value has changed.
func (f *Field) Dirty() bool {
	return f.dirty
}

// Errors returns the validation errors.
func (f *Field) Errors() []error {
	return f.errors
}

// WithValidators returns a new field with the specified validators.
func (f *Field) WithValidators(validators ...Validator) *Field {
	return &Field{
		name:       f.name,
		label:      f.label,
		model:      f.model,
		validators: validators,
		touched:    f.touched,
		dirty:      f.dirty,
		errors:     f.errors,
	}
}

// WithModel returns a new field with the updated model.
func (f *Field) WithModel(model FieldModel) *Field {
	return &Field{
		name:       f.name,
		label:      f.label,
		model:      model,
		validators: f.validators,
		touched:    f.touched,
		dirty:      f.dirty,
		errors:     f.errors,
	}
}

// MarkTouched returns a new field marked as touched.
func (f *Field) MarkTouched() *Field {
	return &Field{
		name:       f.name,
		label:      f.label,
		model:      f.model,
		validators: f.validators,
		touched:    true,
		dirty:      f.dirty,
		errors:     f.errors,
	}
}

// MarkDirty returns a new field marked as dirty.
func (f *Field) MarkDirty() *Field {
	return &Field{
		name:       f.name,
		label:      f.label,
		model:      f.model,
		validators: f.validators,
		touched:    f.touched,
		dirty:      true,
		errors:     f.errors,
	}
}

// WithErrors returns a new field with the specified errors.
func (f *Field) WithErrors(errors []error) *Field {
	return &Field{
		name:       f.name,
		label:      f.label,
		model:      f.model,
		validators: f.validators,
		touched:    f.touched,
		dirty:      f.dirty,
		errors:     errors,
	}
}

// ClearErrors returns a new field with errors cleared.
func (f *Field) ClearErrors() *Field {
	return &Field{
		name:       f.name,
		label:      f.label,
		model:      f.model,
		validators: f.validators,
		touched:    f.touched,
		dirty:      f.dirty,
		errors:     []error{},
	}
}

// IsValid returns whether the field has no errors.
func (f *Field) IsValid() bool {
	return len(f.errors) == 0
}
