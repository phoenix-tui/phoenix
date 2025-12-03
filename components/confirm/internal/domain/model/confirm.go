// Package model provides the domain model for the confirm component.
package model

import (
	"github.com/phoenix-tui/phoenix/components/confirm/internal/domain/value"
)

// Result represents the outcome of the confirm dialog.
type Result int

const (
	// ResultNone means no selection has been made yet.
	ResultNone Result = iota
	// ResultYes means the affirmative button was selected.
	ResultYes
	// ResultNo means the negative button was selected.
	ResultNo
	// ResultCanceled means the cancel button was selected or Esc was pressed.
	ResultCanceled
)

// Confirm represents the domain model for a confirmation dialog.
// It follows rich domain model pattern with encapsulated behavior.
type Confirm struct {
	title       string
	description string
	buttons     []*value.Button
	focused     int
	result      Result
	done        bool
}

// New creates a new Confirm with the given title.
// Default buttons are "Yes" and "No" with "No" focused (safe default).
func New(title string) *Confirm {
	return &Confirm{
		title:       title,
		description: "",
		buttons: []*value.Button{
			value.NewButton("Yes"),
			value.NewButton("No"),
		},
		focused: 1, // Default to "No" (safe choice)
		result:  ResultNone,
		done:    false,
	}
}

// WithDescription returns a new Confirm with the specified description.
func (c *Confirm) WithDescription(desc string) *Confirm {
	return &Confirm{
		title:       c.title,
		description: desc,
		buttons:     c.buttons,
		focused:     c.focused,
		result:      c.result,
		done:        c.done,
	}
}

// WithButtons returns a new Confirm with custom button labels.
func (c *Confirm) WithButtons(labels ...string) *Confirm {
	buttons := make([]*value.Button, len(labels))
	for i, label := range labels {
		buttons[i] = value.NewButton(label)
	}

	// Keep focused index within bounds
	focused := c.focused
	if focused >= len(buttons) {
		focused = len(buttons) - 1
	}
	if focused < 0 {
		focused = 0
	}

	return &Confirm{
		title:       c.title,
		description: c.description,
		buttons:     buttons,
		focused:     focused,
		result:      c.result,
		done:        c.done,
	}
}

// WithDefaultYes returns a new Confirm with "Yes" focused by default.
func (c *Confirm) WithDefaultYes() *Confirm {
	return c.withFocusedButton(0)
}

// WithDefaultNo returns a new Confirm with "No" focused by default.
func (c *Confirm) WithDefaultNo() *Confirm {
	// Find "No" button index
	for i, btn := range c.buttons {
		if btn.Label() == "No" {
			return c.withFocusedButton(i)
		}
	}
	// Fallback to last button if "No" not found
	return c.withFocusedButton(len(c.buttons) - 1)
}

// MoveFocusLeft moves focus to the previous button.
func (c *Confirm) MoveFocusLeft() *Confirm {
	newFocused := c.focused - 1
	if newFocused < 0 {
		newFocused = len(c.buttons) - 1 // Wrap around
	}
	return c.withFocusedButton(newFocused)
}

// MoveFocusRight moves focus to the next button.
func (c *Confirm) MoveFocusRight() *Confirm {
	newFocused := c.focused + 1
	if newFocused >= len(c.buttons) {
		newFocused = 0 // Wrap around
	}
	return c.withFocusedButton(newFocused)
}

// Confirm confirms the currently focused button.
func (c *Confirm) Confirm() *Confirm {
	if c.focused < 0 || c.focused >= len(c.buttons) {
		return c
	}

	label := c.buttons[c.focused].Label()
	result := c.labelToResult(label)

	return &Confirm{
		title:       c.title,
		description: c.description,
		buttons:     c.buttons,
		focused:     c.focused,
		result:      result,
		done:        true,
	}
}

// ConfirmKey confirms the button matching the given key.
func (c *Confirm) ConfirmKey(key rune) *Confirm {
	for i, btn := range c.buttons {
		if btn.MatchesKey(key) {
			// Focus and confirm this button
			focused := c.withFocusedButton(i)
			return focused.Confirm()
		}
	}
	return c
}

// Cancel cancels the dialog (Esc key).
func (c *Confirm) Cancel() *Confirm {
	return &Confirm{
		title:       c.title,
		description: c.description,
		buttons:     c.buttons,
		focused:     c.focused,
		result:      ResultCanceled,
		done:        true,
	}
}

// Title returns the confirm dialog title.
func (c *Confirm) Title() string {
	return c.title
}

// Description returns the confirm dialog description.
func (c *Confirm) Description() string {
	return c.description
}

// Buttons returns the button labels.
func (c *Confirm) Buttons() []string {
	labels := make([]string, len(c.buttons))
	for i, btn := range c.buttons {
		labels[i] = btn.Label()
	}
	return labels
}

// FocusedIndex returns the index of the focused button.
func (c *Confirm) FocusedIndex() int {
	return c.focused
}

// Result returns the result of the confirmation.
func (c *Confirm) Result() Result {
	return c.result
}

// Done returns true if the user has made a selection.
func (c *Confirm) Done() bool {
	return c.done
}

// withFocusedButton returns a new Confirm with the specified focused button.
func (c *Confirm) withFocusedButton(index int) *Confirm {
	return &Confirm{
		title:       c.title,
		description: c.description,
		buttons:     c.buttons,
		focused:     index,
		result:      c.result,
		done:        c.done,
	}
}

// labelToResult converts a button label to a Result.
func (c *Confirm) labelToResult(label string) Result {
	switch label {
	case "Yes":
		return ResultYes
	case "No":
		return ResultNo
	case "Cancel":
		return ResultCanceled
	default:
		// For custom labels, assume first button = Yes, others = No
		if c.focused == 0 {
			return ResultYes
		}
		return ResultNo
	}
}
