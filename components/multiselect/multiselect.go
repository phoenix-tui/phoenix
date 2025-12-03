// Package multiselect provides a multi-choice selection component with fuzzy filtering.
//
// The MultiSelect component allows users to choose multiple options from a list with:
// - Keyboard navigation (j/k or arrows)
// - Toggle selection with Space
// - Select all/none with a/n
// - Optional fuzzy filtering
// - Min/Max selection constraints
// - Virtualization for large lists
//
// Example (basic string multi-selection):
//
//	frameworks := []string{"Phoenix", "Charm", "Bubbletea", "Ink", "Textual"}
//	multi := multiselect.NewStrings("Select frameworks:", frameworks)
//	p := tea.NewProgram(multi)
//	p.Run()
//
// Example (type-safe with constraints):
//
//	type Feature struct {
//	    ID   int
//	    Name string
//	}
//
//	multi := multiselect.New[Feature]("Select features:").
//	    Options(
//	        multiselect.Opt("Auth", Feature{1, "Auth"}),
//	        multiselect.Opt("API", Feature{2, "API"}),
//	        multiselect.Opt("DB", Feature{3, "DB"}),
//	    ).
//	    WithFilterable(true).
//	    Selected(0, 2).  // Pre-select indices 0 and 2
//	    Min(1).          // At least 1 must be selected
//	    Max(3)           // At most 3 can be selected
package multiselect

import (
	"fmt"
	"strings"

	"github.com/phoenix-tui/phoenix/components/multiselect/internal/domain/model"
	"github.com/phoenix-tui/phoenix/components/multiselect/internal/domain/value"
	"github.com/phoenix-tui/phoenix/components/multiselect/internal/infrastructure"
	"github.com/phoenix-tui/phoenix/style"
	"github.com/phoenix-tui/phoenix/tea"
)

// MultiSelect is the public API for the multi-choice selection component.
// It implements tea.Model for use in Elm Architecture applications.
//nolint:unused // theme field will be used for View rendering in future iterations
type MultiSelect[T any] struct {
	theme       *style.Theme  // Optional theme, defaults to DefaultTheme if nil
	title      string
	domain     *model.MultiSelect[T]
	keymap     *infrastructure.KeyBindingMap
	filterable bool
	min        int
	max        int
}

// New creates a new MultiSelect with the given title.
func New[T any](title string) *MultiSelect[T] {
	return &MultiSelect[T]{
		title:      title,
		domain:     model.New[T]([]*value.Option[T]{}, 0, 0),
		keymap:     infrastructure.DefaultKeyBindingMap(),
		filterable: false,
		min:        0,
		max:        0,
	}
}

// Options sets the available options.
func (m *MultiSelect[T]) Options(options ...*value.Option[T]) *MultiSelect[T] {
	return &MultiSelect[T]{
		title:      m.title,
		domain:     model.New(options, m.min, m.max),
		keymap:     m.keymap,
		filterable: m.filterable,
		min:        m.min,
		max:        m.max,
	}
}

// WithHeight sets the visible height of the option list.
func (m *MultiSelect[T]) WithHeight(height int) *MultiSelect[T] {
	return &MultiSelect[T]{
		title:      m.title,
		domain:     m.domain.WithHeight(height),
		keymap:     m.keymap,
		filterable: m.filterable,
		min:        m.min,
		max:        m.max,
	}
}

// WithFilterable enables or disables fuzzy filtering.
func (m *MultiSelect[T]) WithFilterable(enabled bool) *MultiSelect[T] {
	return &MultiSelect[T]{
		title:      m.title,
		domain:     m.domain,
		keymap:     m.keymap,
		filterable: enabled,
		min:        m.min,
		max:        m.max,
	}
}

// WithFilterFunc sets a custom filter function.
func (m *MultiSelect[T]) WithFilterFunc(fn func(opt *value.Option[T], query string) bool) *MultiSelect[T] {
	return &MultiSelect[T]{
		title:      m.title,
		domain:     m.domain.WithFilterFunc(fn),
		keymap:     m.keymap,
		filterable: m.filterable,
		min:        m.min,
		max:        m.max,
	}
}

// WithRenderFunc sets a custom option rendering function.
func (m *MultiSelect[T]) WithRenderFunc(fn func(opt *value.Option[T], index int, focused bool, selected bool) string) *MultiSelect[T] {
	return &MultiSelect[T]{
		title:      m.title,
		domain:     m.domain.WithRenderFunc(fn),
		keymap:     m.keymap,
		filterable: m.filterable,
		min:        m.min,
		max:        m.max,
	}
}

// Selected sets the initially selected indices.
func (m *MultiSelect[T]) Selected(indices ...int) *MultiSelect[T] {
	return &MultiSelect[T]{
		title:      m.title,
		domain:     m.domain.WithSelected(indices...),
		keymap:     m.keymap,
		filterable: m.filterable,
		min:        m.min,
		max:        m.max,
	}
}

// Min sets the minimum required selections.
func (m *MultiSelect[T]) Min(minCount int) *MultiSelect[T] {
	if minCount < 0 {
		minCount = 0
	}
	return &MultiSelect[T]{
		title:      m.title,
		domain:     m.domain.WithSelectionConstraints(minCount, m.max), // Recreate with new constraints
		keymap:     m.keymap,
		filterable: m.filterable,
		min:        minCount,
		max:        m.max,
	}
}

// Max sets the maximum allowed selections (0 = unlimited).
func (m *MultiSelect[T]) Max(maxCount int) *MultiSelect[T] {
	if maxCount < 0 {
		maxCount = 0
	}
	return &MultiSelect[T]{
		title:      m.title,
		domain:     m.domain.WithSelectionConstraints(m.min, maxCount), // Recreate with new constraints
		keymap:     m.keymap,
		filterable: m.filterable,
		min:        m.min,
		max:        maxCount,
	}
}

// SelectedItems returns the currently selected values.
func (m *MultiSelect[T]) SelectedItems() []T {
	return m.domain.SelectedValues()
}

// SelectedIndices returns the currently selected indices.
func (m *MultiSelect[T]) SelectedIndices() []int {
	return m.domain.SelectedIndices()
}

// SelectionCount returns the number of selected items.
func (m *MultiSelect[T]) SelectionCount() int {
	return m.domain.SelectionCount()
}

// Init implements tea.Model.
func (m *MultiSelect[T]) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *MultiSelect[T]) Update(msg tea.Msg) (*MultiSelect[T], tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		return m.handleKey(keyMsg)
	}
	return m, nil
}
//nolint:cyclop // Switch statement for keyboard actions is necessary

// handleKey processes keyboard input.
func (m *MultiSelect[T]) handleKey(msg tea.KeyMsg) (*MultiSelect[T], tea.Cmd) {
	action := m.keymap.GetAction(msg)

	newM := &MultiSelect[T]{
		title:      m.title,
		domain:     m.domain,
		keymap:     m.keymap,
		filterable: m.filterable,
		min:        m.min,
		max:        m.max,
	}

	switch action {
	case infrastructure.ActionMoveUp:
		newM.domain = newM.domain.MoveUp()
	case infrastructure.ActionMoveDown:
		newM.domain = newM.domain.MoveDown()
	case infrastructure.ActionMoveToStart:
		newM.domain = newM.domain.MoveToStart()
	case infrastructure.ActionMoveToEnd:
		newM.domain = newM.domain.MoveToEnd()
	case infrastructure.ActionToggle:
		newM.domain = newM.domain.Toggle()
	case infrastructure.ActionSelectAll:
		newM.domain = newM.domain.SelectAll()
	case infrastructure.ActionSelectNone:
		newM.domain = newM.domain.SelectNone()
	case infrastructure.ActionConfirm:
		if newM.domain.CanConfirm() {
			return newM, ConfirmSelectionCmd[T](newM.SelectedItems)
		}
		// Don't confirm if min constraint not met
		return newM, nil
	case infrastructure.ActionClearFilter:
		newM.domain = newM.domain.ClearFilter()
	case infrastructure.ActionQuit:
		return newM, tea.Quit()
	default:
		newM = m.handleFilterInput(newM, msg)
	}

	return newM, nil
}

// handleFilterInput handles filter input for text-based filtering.
func (m *MultiSelect[T]) handleFilterInput(newM *MultiSelect[T], msg tea.KeyMsg) *MultiSelect[T] {
	if !m.filterable {
		return newM
	}

	if msg.Type == tea.KeyRune {
		// Skip if it's a bound rune (a, n, j, k, g, G)
		action := m.keymap.GetAction(msg)
		if action != infrastructure.ActionNone {
			return newM
		}

		currentQuery := m.domain.FilterQuery()
		newQuery := currentQuery + string(msg.Rune)
		newM.domain = newM.domain.SetFilterQuery(newQuery)
		return newM
	}

	if msg.Type == tea.KeyBackspace {
		currentQuery := m.domain.FilterQuery()
		if currentQuery != "" {
			newQuery := currentQuery[:len(currentQuery)-1]
			newM.domain = newM.domain.SetFilterQuery(newQuery)
		}
		return newM
	}

	return newM
}

// View implements tea.Model.
func (m *MultiSelect[T]) View() string {
	var b strings.Builder

	// Render title with selection summary
	if m.title != "" {
		b.WriteString(m.title)
		if m.domain.TotalCount() > 0 {
			b.WriteString(fmt.Sprintf(" (%d of %d selected)",
				m.domain.SelectionCount(),
				m.domain.TotalCount()))
		}
		b.WriteRune('\n')
		b.WriteRune('\n')
	}

	// Render options
	m.renderOptions(&b)

	// Show filter query if active
	if m.filterable {
		m.renderFilterStatus(&b)
	}

	// Show help text
	m.renderHelp(&b)

	return b.String()
}

// renderOptions renders the option list or empty state message.
func (m *MultiSelect[T]) renderOptions(b *strings.Builder) {
	options := m.domain.RenderVisibleOptions()

	if len(options) == 0 {
		m.renderEmptyState(b)
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
func (m *MultiSelect[T]) renderEmptyState(b *strings.Builder) {
	if m.domain.IsFiltered() {
		b.WriteString("(no matches)")
	} else {
		b.WriteString("(no options)")
	}
}

// renderFilterStatus renders the filter input status.
func (m *MultiSelect[T]) renderFilterStatus(b *strings.Builder) {
	_ = b.WriteByte('\n')
	_ = b.WriteByte('\n')
	filterQuery := m.domain.FilterQuery()
	if filterQuery != "" {
		b.WriteString("Filter: ")
		b.WriteString(filterQuery)
	} else {
		b.WriteString("Type to filter...")
	}
}

// renderHelp renders the help text at the bottom.
func (m *MultiSelect[T]) renderHelp(b *strings.Builder) {
	_ = b.WriteByte('\n')
	_ = b.WriteByte('\n')
	b.WriteString("  a: all  n: none  Space: toggle  Enter: confirm")
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
func ConfirmSelectionCmd[T any](getValuesFunc func() []T) tea.Cmd {
	return func() tea.Msg {
		values := getValuesFunc()
		return ConfirmSelectionMsg[T]{Values: values}
	}
}

// ConfirmSelectionMsg is sent when the user confirms their selection.
type ConfirmSelectionMsg[T any] struct {
	Values []T
}

// NewStrings creates a new string-based MultiSelect (convenience constructor).
func NewStrings(title string, options []string) *MultiSelect[string] {
	opts := make([]*value.Option[string], len(options))
	for i, opt := range options {
		opts[i] = value.NewOption(opt, opt)
	}
	return New[string](title).Options(opts...)
}
