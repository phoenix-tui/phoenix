// Package form provides a form container component with field management and validation.
//
// The Form component allows users to:
// - Navigate between fields with Tab/Shift+Tab
// - Validate fields individually or all at once
// - Track dirty and touched state
// - Submit or reset the form
//
// Example (basic form with validation):
//
//	import (
//	    "github.com/phoenix-tui/phoenix/components/form"
//	    "github.com/phoenix-tui/phoenix/components/form/internal/domain/value"
//	    "github.com/phoenix-tui/phoenix/components/textinput"
//	)
//
//	// Create text input models
//	nameInput := textinput.New("John Doe", "Enter your name")
//	emailInput := textinput.New("", "Enter your email")
//
//	// Create form with fields
//	f := form.New("User Registration").
//	    Field("name", "Name", nameInput, value.Required(), value.MinLength(2)).
//	    Field("email", "Email", emailInput, value.Required(), value.Email())
//
//	// Run the form
//	p := tea.NewProgram(f)
//	p.Run()
//
// Example (retrieving form values after submission):
//
//	// In your parent model's Update function
//	case form.SubmitMsg:
//	    // Access field values through your stored models
//	    name := nameInput.Value()
//	    email := emailInput.Value()
//	    fmt.Printf("Name: %s, Email: %s\n", name, email)
package form

import (
	"strings"

	"github.com/phoenix-tui/phoenix/components/form/internal/domain/model"
	"github.com/phoenix-tui/phoenix/components/form/internal/domain/value"
	"github.com/phoenix-tui/phoenix/components/form/internal/infrastructure"
	"github.com/phoenix-tui/phoenix/style"
	"github.com/phoenix-tui/phoenix/tea"
)

// Form is the public API for the form container component.
// It implements tea.Model for use in Elm Architecture applications.
//nolint:unused // theme field will be used for View rendering in future iterations
type Form struct {
	theme  *style.Theme  // Optional theme, defaults to DefaultTheme if nil
	domain     *model.Form
	keymap     *infrastructure.KeyBindingMap
	fieldNames map[string]int // Maps field name to index
}

// New creates a new Form with the given title.
func New(title string) *Form {
	return &Form{
		domain:     model.New(title),
		keymap:     infrastructure.DefaultKeyBindingMap(),
		fieldNames: make(map[string]int),
	}
}

// Field adds a field to the form.
// The model parameter must implement View() method (any Phoenix component works).
// Validators are applied in order when the field is validated.
func (f *Form) Field(name, label string, fieldModel value.FieldModel, validators ...value.Validator) *Form {
	field := value.NewField(name, label, fieldModel).WithValidators(validators...)

	fields := append(f.domain.Fields(), field)

	// Update field name mapping
	newFieldNames := make(map[string]int, len(f.fieldNames)+1)
	for k, v := range f.fieldNames {
		newFieldNames[k] = v
	}
	newFieldNames[name] = len(fields) - 1

	return &Form{
		domain:     f.domain.WithFields(fields),
		keymap:     f.keymap,
		fieldNames: newFieldNames,
	}
}

// Values returns a map of field name → field model.
// Use type assertion to extract values from the models.
func (f *Form) Values() map[string]value.FieldModel {
	result := make(map[string]value.FieldModel)

	for _, field := range f.domain.Fields() {
		result[field.Name()] = field.Model()
	}

	return result
}

// Value returns the model for a specific field by name.
// Returns nil if the field doesn't exist.
func (f *Form) Value(name string) value.FieldModel {
	index, ok := f.fieldNames[name]
	if !ok || index >= len(f.domain.Fields()) {
		return nil
	}

	return f.domain.Fields()[index].Model()
}

// IsValid returns whether all fields pass validation.
func (f *Form) IsValid() bool {
	return f.domain.IsValid()
}

// IsDirty returns whether any field has been modified.
func (f *Form) IsDirty() bool {
	for _, field := range f.domain.Fields() {
		if field.Dirty() {
			return true
		}
	}
	return false
}

// Errors returns a map of field name → error messages.
func (f *Form) Errors() map[string][]string {
	return f.domain.Errors()
}

// Init implements tea.Model.
func (f *Form) Init() tea.Cmd {
	// Initialize all field models
	var cmds []tea.Cmd

	for _, field := range f.domain.Fields() {
		if initializer, ok := field.Model().(interface{ Init() tea.Cmd }); ok {
			if cmd := initializer.Init(); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	if len(cmds) == 0 {
		return nil
	}

	return tea.Batch(cmds...)
}

// Update implements tea.Model.
func (f *Form) Update(msg tea.Msg) (*Form, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		return f.handleKey(keyMsg)
	}

	// Delegate message to focused field
	return f.delegateToFocusedField(msg)
}

// handleKey processes keyboard input.
func (f *Form) handleKey(msg tea.KeyMsg) (*Form, tea.Cmd) {
	action := f.keymap.GetAction(msg)

	switch action {
	case infrastructure.ActionNextField:
		newForm := &Form{
			domain:     f.domain.MoveFocusNext(),
			keymap:     f.keymap,
			fieldNames: f.fieldNames,
		}
		return newForm, nil

	case infrastructure.ActionPrevField:
		newForm := &Form{
			domain:     f.domain.MoveFocusPrev(),
			keymap:     f.keymap,
			fieldNames: f.fieldNames,
		}
		return newForm, nil

	case infrastructure.ActionSubmit:
		// Validate all fields before submit
		fieldValues := f.extractFieldValues()
		newDomain := f.domain.ValidateAll(fieldValues)

		newForm := &Form{
			domain:     newDomain,
			keymap:     f.keymap,
			fieldNames: f.fieldNames,
		}

		if newDomain.IsValid() {
			newForm.domain = newForm.domain.Submit()
			return newForm, SubmitCmd()
		}

		return newForm, nil

	case infrastructure.ActionReset:
		newForm := &Form{
			domain:     f.domain.Reset(),
			keymap:     f.keymap,
			fieldNames: f.fieldNames,
		}
		return newForm, ResetCmd()

	case infrastructure.ActionQuit:
		return f, tea.Quit()
	}

	// Not a form-level action, delegate to focused field
	return f.delegateToFocusedField(msg)
}

// delegateToFocusedField forwards the message to the currently focused field.
func (f *Form) delegateToFocusedField(msg tea.Msg) (*Form, tea.Cmd) {
	focusedIndex := f.domain.FocusedIndex()
	if focusedIndex >= len(f.domain.Fields()) {
		return f, nil
	}

	focusedField := f.domain.Fields()[focusedIndex]

	// Update the focused field's model
	if updater, ok := focusedField.Model().(interface {
		Update(tea.Msg) (value.FieldModel, tea.Cmd)
	}); ok {
		newModel, cmd := updater.Update(msg)
		newDomain := f.domain.UpdateField(focusedIndex, newModel)

		newForm := &Form{
			domain:     newDomain,
			keymap:     f.keymap,
			fieldNames: f.fieldNames,
		}

		return newForm, cmd
	}

	return f, nil
}

// extractFieldValues extracts values from field models.
// This assumes field models have a Value() method (like TextInput).
func (f *Form) extractFieldValues() map[string]interface{} {
	result := make(map[string]interface{})

	for _, field := range f.domain.Fields() {
		// Try to extract value using common patterns
		switch valuer := field.Model().(type) {
		case interface{ Value() string }:
			result[field.Name()] = valuer.Value()
		case interface{ SelectedValue() (interface{}, bool) }:
			if val, ok := valuer.SelectedValue(); ok {
				result[field.Name()] = val
			}
		default:
			// Default: use the model itself as value
			result[field.Name()] = field.Model()
		}
	}

	return result
}

// View implements tea.Model.
func (f *Form) View() string {
	var b strings.Builder

	// Render title
	if f.domain.Title() != "" {
		b.WriteString(f.domain.Title())
		_ = b.WriteByte('\n')
		b.WriteString(strings.Repeat("─", len(f.domain.Title())))
		_ = b.WriteByte('\n')
		_ = b.WriteByte('\n')
	}

	// Render fields
	f.renderFields(&b)

	// Render help text
	_ = b.WriteByte('\n')
	b.WriteString(strings.Repeat("─", 40))
	_ = b.WriteByte('\n')
	b.WriteString("Tab: next field  Shift+Tab: prev  Enter: submit")

	return b.String()
}

// renderFields renders all form fields.
func (f *Form) renderFields(b *strings.Builder) {
	fields := f.domain.Fields()
	focusedIndex := f.domain.FocusedIndex()

	for i, field := range fields {
		if i > 0 {
			_ = b.WriteByte('\n')
			_ = b.WriteByte('\n')
		}

		isFocused := (i == focusedIndex)
		f.renderField(b, field, isFocused)
	}
}

// renderField renders a single field with its label, model view, and errors.
func (f *Form) renderField(b *strings.Builder, field *value.Field, focused bool) {
	// Render label
	b.WriteString(field.Label())
	b.WriteString(":")

	// Render focus indicator
	if focused {
		b.WriteString(" ← focused")
	}

	_ = b.WriteByte('\n')

	// Render field model view
	b.WriteString(field.Model().View())

	// Render errors
	if len(field.Errors()) > 0 {
		_ = b.WriteByte('\n')
		for _, err := range field.Errors() {
			b.WriteString("✗ ")
			b.WriteString(err.Error())
		}
	} else if field.Touched() && field.IsValid() {
		_ = b.WriteByte('\n')
		b.WriteString("✓ Valid")
	}
}

// SubmitCmd returns a command that sends a SubmitMsg.
func SubmitCmd() tea.Cmd {
	return func() tea.Msg {
		return SubmitMsg{}
	}
}

// SubmitMsg is sent when the form is submitted and all validations pass.
type SubmitMsg struct{}

// ResetCmd returns a command that sends a ResetMsg.
func ResetCmd() tea.Cmd {
	return func() tea.Msg {
		return ResetMsg{}
	}
}

// ResetMsg is sent when the form is reset.
type ResetMsg struct{}
