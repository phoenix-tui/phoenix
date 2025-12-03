// Package value provides value objects for the select component domain.
package value

// Option represents a selectable option in a list.
// It's a value object containing the label, underlying value, and metadata.
type Option[T any] struct {
	label       string
	value       T
	description string
	disabled    bool
}

// NewOption creates a new option with the given label and value.
func NewOption[T any](label string, value T) *Option[T] {
	return &Option[T]{
		label:       label,
		value:       value,
		description: "",
		disabled:    false,
	}
}

// WithDescription returns a new option with the specified description.
func (o *Option[T]) WithDescription(desc string) *Option[T] {
	return &Option[T]{
		label:       o.label,
		value:       o.value,
		description: desc,
		disabled:    o.disabled,
	}
}

// WithDisabled returns a new option with the specified disabled state.
func (o *Option[T]) WithDisabled(disabled bool) *Option[T] {
	return &Option[T]{
		label:       o.label,
		value:       o.value,
		description: o.description,
		disabled:    disabled,
	}
}

// Label returns the display label for this option.
func (o *Option[T]) Label() string {
	return o.label
}

// Value returns the underlying value for this option.
func (o *Option[T]) Value() T {
	return o.value
}

// Description returns the optional description text.
func (o *Option[T]) Description() string {
	return o.description
}

// Disabled returns true if this option cannot be selected.
func (o *Option[T]) Disabled() bool {
	return o.disabled
}
