// Package list provides a universal selectable list component for Phoenix TUI.
//
// The List component is a foundation for building selection interfaces like file pickers,
// menus, and searchable lists. It supports both single and multi-selection modes,
// custom item rendering, filtering, keyboard navigation, and scrolling for long lists.
//
// Example (basic file picker):
//
//	files := []string{"file1.txt", "file2.go", "file3.md"}
//	l := list.New(files, files, value.SelectionModeSingle).Height(10)
//	p := tea.New(l)
//	p.Run()
//
// Example (custom rendering):
//
//	type Person struct { Name string; Age int }
//	people := []Person{{"Alice", 30}, {"Bob", 25}}
//
//	l := list.New(people, getLabels(people), value.SelectionModeSingle).
//		ItemRenderer(func(item interface{}, idx int, selected, focused bool) string {
//			p := item.(Person)
//			prefix := "  "
//			if selected { prefix = "✓ " }
//			if focused { prefix = "> " }
//			return fmt.Sprintf("%s%s (age %d)", prefix, p.Name, p.Age)
//		})
package list

import (
	"strings"

	"github.com/phoenix-tui/phoenix/components/list/domain/model"
	"github.com/phoenix-tui/phoenix/components/list/domain/value"
	"github.com/phoenix-tui/phoenix/components/list/infrastructure"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

// List is the public API for the List component.
// It implements the tea.Model pattern (Init, Update, View).
type List struct {
	domain     *model.List
	keymap     infrastructure.KeyBindingMap
	showFilter bool // Whether to show filter input at bottom
}

// New creates a new list with the given items and selection mode.
// Values and labels must be the same length.
func New(values []interface{}, labels []string, mode value.SelectionMode) *List {
	if len(values) != len(labels) {
		panic("list.New: values and labels must have the same length")
	}

	items := make([]*value.Item, len(values))
	for i := range values {
		items[i] = value.NewItem(values[i], labels[i])
	}

	return &List{
		domain:     model.NewListWithItems(items, mode),
		keymap:     infrastructure.DefaultKeyBindingMap(),
		showFilter: false,
	}
}

// NewSingleSelect is a convenience constructor for single-selection lists.
func NewSingleSelect(values []interface{}, labels []string) *List {
	return New(values, labels, value.SelectionModeSingle)
}

// NewMultiSelect is a convenience constructor for multi-selection lists.
func NewMultiSelect(values []interface{}, labels []string) *List {
	return New(values, labels, value.SelectionModeMulti)
}

// Height sets the visible height of the list.
func (l *List) Height(height int) *List {
	newList := l.clone()
	newList.domain = newList.domain.WithHeight(height)
	return newList
}

// ItemRenderer sets a custom item renderer function.
// The function receives the item, index, selected state, and focused state.
func (l *List) ItemRenderer(renderer func(item interface{}, index int, selected, focused bool) string) *List {
	newList := l.clone()
	// Wrap the user's renderer to work with domain Item
	wrappedRenderer := func(domainItem *value.Item, index int, selected, focused bool) string {
		return renderer(domainItem.Value(), index, selected, focused)
	}
	newList.domain = newList.domain.WithItemRenderer(wrappedRenderer)
	return newList
}

// Filter sets a custom filter function.
// The function receives the item and query, and returns true if the item matches.
func (l *List) Filter(filterFunc func(item interface{}, query string) bool) *List {
	newList := l.clone()
	// Wrap the user's filter to work with domain Item
	wrappedFilter := func(domainItem *value.Item, query string) bool {
		return filterFunc(domainItem.Value(), query)
	}
	newList.domain = newList.domain.WithFilter(wrappedFilter)
	return newList
}

// ShowFilter enables the filter input display at the bottom of the list.
func (l *List) ShowFilter(show bool) *List {
	newList := l.clone()
	newList.showFilter = show
	return newList
}

// KeyBindings sets custom key bindings.
func (l *List) KeyBindings(bindings []infrastructure.KeyBinding) *List {
	newList := l.clone()
	newList.keymap = infrastructure.NewKeyBindingMap(bindings)
	return newList
}

// clone creates a shallow copy of the list for immutability
func (l *List) clone() *List {
	return &List{
		domain:     l.domain,
		keymap:     l.keymap,
		showFilter: l.showFilter,
	}
}

// SelectedItems returns the currently selected item values.
func (l *List) SelectedItems() []interface{} {
	domainItems := l.domain.SelectedItems()
	result := make([]interface{}, len(domainItems))
	for i, item := range domainItems {
		result[i] = item.Value()
	}
	return result
}

// FocusedItem returns the currently focused item value, or nil if no items.
func (l *List) FocusedItem() interface{} {
	item := l.domain.FocusedItem()
	if item == nil {
		return nil
	}
	return item.Value()
}

// SelectedIndices returns the indices of selected items.
func (l *List) SelectedIndices() []int {
	return l.domain.SelectedIndices()
}

// FocusedIndex returns the index of the focused item.
func (l *List) FocusedIndex() int {
	return l.domain.FocusedIndex()
}

// Init implements tea.Model.
func (l *List) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (l *List) Update(msg tea.Msg) (*List, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return l.handleKey(msg)
	}
	return l, nil
}

// handleKey processes keyboard input.
func (l *List) handleKey(msg tea.KeyMsg) (*List, tea.Cmd) {
	action := l.keymap.GetAction(msg)
	newList := l.clone()

	switch action {
	case "move_up":
		newList.domain = newList.domain.MoveUp()
	case "move_down":
		newList.domain = newList.domain.MoveDown()
	case "page_up":
		newList.domain = newList.domain.MovePageUp()
	case "page_down":
		newList.domain = newList.domain.MovePageDown()
	case "move_to_start":
		newList.domain = newList.domain.MoveToStart()
	case "move_to_end":
		newList.domain = newList.domain.MoveToEnd()
	case "toggle_selection":
		newList.domain = newList.domain.ToggleSelection()
	case "select_all":
		newList.domain = newList.domain.SelectAll()
	case "clear_selection":
		newList.domain = newList.domain.ClearSelection()
	case "clear_filter":
		newList.domain = newList.domain.ClearFilter()
	case "confirm":
		// Return selected items (application can handle this)
		return newList, nil
	case "quit":
		return newList, tea.Quit()
	default:
		// If no action matched, check if it's a printable character for filtering
		if newList.showFilter && msg.Type == tea.KeyRune {
			// Add character to filter query
			currentQuery := string(msg.Rune)
			newList.domain = newList.domain.SetFilterQuery(currentQuery)
		} else if newList.showFilter && msg.Type == tea.KeyBackspace {
			// Remove last character from filter query
			newList.domain = newList.domain.ClearFilter()
		}
	}

	return newList, nil
}

// View implements tea.Model.
func (l *List) View() string {
	var b strings.Builder

	// Render visible items
	items := l.domain.RenderVisibleItems()
	if len(items) == 0 {
		if l.domain.IsFiltered() {
			b.WriteString("(no matches)")
		} else {
			b.WriteString("(empty list)")
		}
	} else {
		for i, item := range items {
			if i > 0 {
				b.WriteRune('\n')
			}
			b.WriteString(item)
		}
	}

	// Show filter status if enabled
	if l.showFilter && l.domain.IsFiltered() {
		b.WriteRune('\n')
		b.WriteString("Filter: (active)")
	}

	return b.String()
}
