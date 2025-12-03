// Package confirm provides a Yes/No/Cancel dialog component for user confirmation.
//
// The Confirm component allows users to confirm or cancel actions with keyboard navigation,
// configurable button labels, and safe defaults (No button focused for dangerous actions).
//
// Example (basic yes/no):
//
//	c := confirm.New("Delete file?")
//	p := tea.NewProgram(c)
//	finalModel, _ := p.Run()
//	if result, ok := finalModel.(*confirm.Confirm); ok {
//	    if result.IsYes() {
//	        // User confirmed
//	    }
//	}
//
// Example (dangerous action with description):
//
//	c := confirm.New("Delete all files?").
//	    Description("This action cannot be undone.").
//	    Affirmative("Delete").
//	    Negative("Cancel").
//	    DefaultNo() // Safe default for dangerous action
//
// Example (three-button with cancel):
//
//	c := confirm.New("Save changes?").
//	    Affirmative("Save").
//	    Negative("Discard").
//	    WithCancel(true)
package confirm

import (
	"strings"

	"github.com/phoenix-tui/phoenix/components/confirm/internal/domain/model"
	"github.com/phoenix-tui/phoenix/components/confirm/internal/infrastructure"
	"github.com/phoenix-tui/phoenix/style"
	"github.com/phoenix-tui/phoenix/tea"
)

// Confirm is the public API for the confirmation dialog component.
// It implements tea.Model for use in Elm Architecture applications.
//nolint:unused // theme field will be used for View rendering in future iterations
type Confirm struct {
	theme  *style.Theme  // Optional theme, defaults to DefaultTheme if nil
	domain  *model.Confirm
	keymap  *infrastructure.KeyBindingMap
	focused bool // Whether this component has input focus
}

// New creates a new Confirm dialog with the given title.
// Default buttons are "Yes" and "No" with "No" focused (safe default).
func New(title string) *Confirm {
	return &Confirm{
		domain:  model.New(title),
		keymap:  infrastructure.DefaultKeyBindingMap(),
		focused: true,
	}
}

// Description sets the description text shown below the title.
func (c *Confirm) Description(desc string) *Confirm {
	return &Confirm{
		domain:  c.domain.WithDescription(desc),
		keymap:  c.keymap,
		focused: c.focused,
	}
}

// Affirmative sets the label for the affirmative button (default: "Yes").
func (c *Confirm) Affirmative(label string) *Confirm {
	buttons := c.domain.Buttons()
	if len(buttons) < 2 {
		buttons = []string{"Yes", "No"}
	}
	buttons[0] = label
	return c.withButtons(buttons...)
}

// Negative sets the label for the negative button (default: "No").
func (c *Confirm) Negative(label string) *Confirm {
	buttons := c.domain.Buttons()
	if len(buttons) < 2 {
		buttons = []string{"Yes", "No"}
	}
	buttons[1] = label
	return c.withButtons(buttons...)
}

// WithCancel adds a third "Cancel" button (for Yes/No/Cancel mode).
func (c *Confirm) WithCancel(enabled bool) *Confirm {
	buttons := c.domain.Buttons()
	if enabled && len(buttons) == 2 {
		buttons = append(buttons, "Cancel")
	} else if !enabled && len(buttons) == 3 {
		buttons = buttons[:2]
	}
	return c.withButtons(buttons...)
}

// DefaultYes focuses the "Yes" button by default.
func (c *Confirm) DefaultYes() *Confirm {
	return &Confirm{
		domain:  c.domain.WithDefaultYes(),
		keymap:  c.keymap,
		focused: c.focused,
	}
}

// DefaultNo focuses the "No" button by default (recommended for dangerous actions).
func (c *Confirm) DefaultNo() *Confirm {
	return &Confirm{
		domain:  c.domain.WithDefaultNo(),
		keymap:  c.keymap,
		focused: c.focused,
	}
}

// IsYes returns true if the user selected the affirmative button.
func (c *Confirm) IsYes() bool {
	return c.domain.Result() == model.ResultYes
}

// IsNo returns true if the user selected the negative button.
func (c *Confirm) IsNo() bool {
	return c.domain.Result() == model.ResultNo
}

// IsCanceled returns true if the user canceled (Esc or Cancel button).
func (c *Confirm) IsCanceled() bool {
	return c.domain.Result() == model.ResultCanceled
}

// Done returns true if the user has made a selection.
func (c *Confirm) Done() bool {
	return c.domain.Done()
}

// Result returns the raw result value.
func (c *Confirm) Result() model.Result {
	return c.domain.Result()
}

// Init implements tea.Model.
func (c *Confirm) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (c *Confirm) Update(msg tea.Msg) (*Confirm, tea.Cmd) {
	if !c.focused || c.domain.Done() {
		return c, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		return c.handleKey(keyMsg)
	}

	return c, nil
}

// handleKey processes keyboard input.
func (c *Confirm) handleKey(msg tea.KeyMsg) (*Confirm, tea.Cmd) {
	action := c.keymap.GetAction(msg)

	newC := &Confirm{
		domain:  c.domain,
		keymap:  c.keymap,
		focused: c.focused,
	}

	switch action {
	case infrastructure.ActionMoveLeft:
		newC.domain = newC.domain.MoveFocusLeft()
	case infrastructure.ActionMoveRight:
		newC.domain = newC.domain.MoveFocusRight()
	case infrastructure.ActionConfirm:
		newC.domain = newC.domain.Confirm()
		return newC, ConfirmResultCmd(newC.domain.Result())
	case infrastructure.ActionCancel:
		newC.domain = newC.domain.Cancel()
		return newC, ConfirmResultCmd(model.ResultCanceled)
	case infrastructure.ActionShortcut:
		// Try keyboard shortcut (y/n/c)
		newC.domain = newC.domain.ConfirmKey(msg.Rune)
		if newC.domain.Done() {
			return newC, ConfirmResultCmd(newC.domain.Result())
		}
	}

	return newC, nil
}

// View implements tea.Model.
func (c *Confirm) View() string {
	var b strings.Builder

	// Render title
	b.WriteString(c.domain.Title())
	_ = b.WriteByte('\n')

	// Render description if present
	if c.domain.Description() != "" {
		b.WriteString(c.domain.Description())
		_ = b.WriteByte('\n')
	}

	_ = b.WriteByte('\n')

	// Render buttons
	c.renderButtons(&b)

	return b.String()
}

// renderButtons renders the button row.
func (c *Confirm) renderButtons(b *strings.Builder) {
	buttons := c.domain.Buttons()
	focused := c.domain.FocusedIndex()

	for i, label := range buttons {
		if i > 0 {
			b.WriteString("  ")
		}

		isFocused := (i == focused)
		c.renderButton(b, label, isFocused)
	}
}

// renderButton renders a single button with focus indicator.
func (c *Confirm) renderButton(b *strings.Builder, label string, focused bool) {
	if focused {
		b.WriteString("[ ")
		b.WriteString(label)
		b.WriteString(" ]")
	} else {
		b.WriteString("( ")
		b.WriteString(label)
		b.WriteString(" )")
	}
}

// withButtons returns a new Confirm with the specified button labels.
func (c *Confirm) withButtons(labels ...string) *Confirm {
	return &Confirm{
		domain:  c.domain.WithButtons(labels...),
		keymap:  c.keymap,
		focused: c.focused,
	}
}

// ConfirmResultCmd returns a command that sends a ConfirmResultMsg.
// Name intentionally stutters for consistency with Select component (ConfirmSelectionMsg).
//
//nolint:revive // Intentional stuttering for API consistency
func ConfirmResultCmd(result model.Result) tea.Cmd {
	return func() tea.Msg {
		return ConfirmResultMsg{Result: result}
	}
}

// ConfirmResultMsg is sent when the user makes a selection.
// Name intentionally stutters for consistency with Select component (ConfirmSelectionMsg).
//
//nolint:revive // Intentional stuttering for API consistency
type ConfirmResultMsg struct {
	Result model.Result
}
