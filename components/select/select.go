// Package selectcomponent provides a single-choice selection component with fuzzy filtering.
//
// The Select component allows users to choose one option from a list with keyboard navigation,
// optional fuzzy filtering, and virtualization for large lists (10,000+ items).
//
// Example (basic string selection):
//
//	frameworks := []string{"Phoenix", "Charm", "Bubbletea", "Ink", "Textual"}
//	sel := selectcomponent.NewString("Choose framework:", frameworks)
//	p := tea.NewProgram(sel)
//	p.Run()
//
// Example (type-safe enum selection):
//
//	type Status int
//	const (
//	    Active Status = iota
//	    Inactive
//	    Pending
//	)
//
//	sel := selectcomponent.New[Status]("Choose status:").
//	    Options(
//	        selectcomponent.Opt("Active", Active, "User is active"),
//	        selectcomponent.Opt("Inactive", Inactive, "User is inactive"),
//	        selectcomponent.Opt("Pending", Pending, "Pending approval"),
//	    ).
//	    WithDefault(Active)
package selectcomponent

import (
	"strings"

	"github.com/phoenix-tui/phoenix/components/select/internal/domain/model"
	"github.com/phoenix-tui/phoenix/components/select/internal/domain/value"
	"github.com/phoenix-tui/phoenix/components/select/internal/infrastructure"
	"github.com/phoenix-tui/phoenix/tea"
)

// Select is the public API for the single-choice selection component.
// It implements tea.Model for use in Elm Architecture applications.
type Select[T any] struct {
	title      string
	domain     *model.Select[T]
	keymap     *infrastructure.KeyBindingMap
	filterable bool
}

// New creates a new Select with the given title.
func New[T any](title string) *Select[T] {
	return &Select[T]{
		title:      title,
		domain:     model.New[T]([]*value.Option[T]{}),
		keymap:     infrastructure.DefaultKeyBindingMap(),
		filterable: false,
	}
}

// Options sets the available options.
func (s *Select[T]) Options(options ...*value.Option[T]) *Select[T] {
	return &Select[T]{
		title:      s.title,
		domain:     model.New(options),
		keymap:     s.keymap,
		filterable: s.filterable,
	}
}

// WithHeight sets the visible height of the option list.
func (s *Select[T]) WithHeight(height int) *Select[T] {
	return &Select[T]{
		title:      s.title,
		domain:     s.domain.WithHeight(height),
		keymap:     s.keymap,
		filterable: s.filterable,
	}
}

// WithFilterable enables or disables fuzzy filtering.
func (s *Select[T]) WithFilterable(enabled bool) *Select[T] {
	return &Select[T]{
		title:      s.title,
		domain:     s.domain,
		keymap:     s.keymap,
		filterable: enabled,
	}
}

// WithFilterFunc sets a custom filter function.
func (s *Select[T]) WithFilterFunc(fn func(opt *value.Option[T], query string) bool) *Select[T] {
	return &Select[T]{
		title:      s.title,
		domain:     s.domain.WithFilterFunc(fn),
		keymap:     s.keymap,
		filterable: s.filterable,
	}
}

// WithRenderFunc sets a custom option rendering function.
func (s *Select[T]) WithRenderFunc(fn func(opt *value.Option[T], index int, focused bool) string) *Select[T] {
	return &Select[T]{
		title:      s.title,
		domain:     s.domain.WithRenderFunc(fn),
		keymap:     s.keymap,
		filterable: s.filterable,
	}
}

// WithDefault sets the default selected value.
func (s *Select[T]) WithDefault(_ T) *Select[T] {
	// TODO: Find matching option and select it
	// This is a simplified implementation for now
	return s
}

// SelectedValue returns the currently selected value, or zero value if nothing selected.
func (s *Select[T]) SelectedValue() (T, bool) {
	return s.domain.SelectedValue()
}

// FocusedValue returns the currently focused value, or zero value if no options.
func (s *Select[T]) FocusedValue() (T, bool) {
	return s.domain.FocusedValue()
}

// Init implements tea.Model.
func (s *Select[T]) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (s *Select[T]) Update(msg tea.Msg) (*Select[T], tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		return s.handleKey(keyMsg)
	}
	return s, nil
}

// handleKey processes keyboard input.
func (s *Select[T]) handleKey(msg tea.KeyMsg) (*Select[T], tea.Cmd) {
	action := s.keymap.GetAction(msg)

	newS := &Select[T]{
		title:      s.title,
		domain:     s.domain,
		keymap:     s.keymap,
		filterable: s.filterable,
	}

	switch action {
	case infrastructure.ActionMoveUp:
		newS.domain = newS.domain.MoveUp()
	case infrastructure.ActionMoveDown:
		newS.domain = newS.domain.MoveDown()
	case infrastructure.ActionMoveToStart:
		newS.domain = newS.domain.MoveToStart()
	case infrastructure.ActionMoveToEnd:
		newS.domain = newS.domain.MoveToEnd()
	case infrastructure.ActionSelect:
		newS.domain = newS.domain.Select()
		return newS, ConfirmSelectionCmd[T](newS.SelectedValue)
	case infrastructure.ActionClearFilter:
		newS.domain = newS.domain.ClearFilter()
	case infrastructure.ActionQuit:
		return newS, tea.Quit()
	default:
		newS = s.handleFilterInput(newS, msg)
	}

	return newS, nil
}

// handleFilterInput handles filter input for text-based filtering.
func (s *Select[T]) handleFilterInput(newS *Select[T], msg tea.KeyMsg) *Select[T] {
	if !s.filterable {
		return newS
	}

	if msg.Type == tea.KeyRune {
		currentQuery := s.domain.FilterQuery()
		newQuery := currentQuery + string(msg.Rune)
		newS.domain = newS.domain.SetFilterQuery(newQuery)
		return newS
	}

	if msg.Type == tea.KeyBackspace {
		currentQuery := s.domain.FilterQuery()
		if currentQuery != "" {
			newQuery := currentQuery[:len(currentQuery)-1]
			newS.domain = newS.domain.SetFilterQuery(newQuery)
		}
		return newS
	}

	return newS
}

// View implements tea.Model.
func (s *Select[T]) View() string {
	var b strings.Builder

	// Render title
	if s.title != "" {
		b.WriteString(s.title)
		b.WriteRune('\n')
		b.WriteRune('\n')
	}

	// Render options
	s.renderOptions(&b)

	// Show filter query if active
	if s.filterable {
		s.renderFilterStatus(&b)
	}

	return b.String()
}

// renderOptions renders the option list or empty state message.
func (s *Select[T]) renderOptions(b *strings.Builder) {
	options := s.domain.RenderVisibleOptions()

	if len(options) == 0 {
		s.renderEmptyState(b)
		return
	}

	for i, opt := range options {
		if i > 0 {
			_ = b.WriteByte('\n')
		}
		b.WriteString(opt)
	}
}

// renderEmptyState renders the empty list message.
func (s *Select[T]) renderEmptyState(b *strings.Builder) {
	if s.domain.IsFiltered() {
		b.WriteString("(no matches)")
	} else {
		b.WriteString("(no options)")
	}
}

// renderFilterStatus renders the filter input status at the bottom.
func (s *Select[T]) renderFilterStatus(b *strings.Builder) {
	_ = b.WriteByte('\n')
	_ = b.WriteByte('\n')
	filterQuery := s.domain.FilterQuery()
	if filterQuery != "" {
		b.WriteString("Filter: ")
		b.WriteString(filterQuery)
	} else {
		b.WriteString("Type to filter...")
	}
}

// Opt creates a new option with label, value, and optional description.
func Opt[T any](label string, val T, description ...string) *value.Option[T] {
	opt := value.NewOption(label, val)
	if len(description) > 0 && description[0] != "" {
		opt = opt.WithDescription(description[0])
	}
	return opt
}

// ConfirmSelectionCmd returns a command that sends a ConfirmSelectionMsg.
func ConfirmSelectionCmd[T any](getValueFunc func() (T, bool)) tea.Cmd {
	return func() tea.Msg {
		val, ok := getValueFunc()
		return ConfirmSelectionMsg[T]{Value: val, OK: ok}
	}
}

// ConfirmSelectionMsg is sent when the user confirms their selection.
type ConfirmSelectionMsg[T any] struct {
	Value T
	OK    bool
}

// NewString creates a new string-based Select (convenience constructor).
func NewString(title string, options []string) *Select[string] {
	opts := make([]*value.Option[string], len(options))
	for i, opt := range options {
		opts[i] = value.NewOption(opt, opt)
	}
	return New[string](title).Options(opts...)
}
