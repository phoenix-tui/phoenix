// Package model contains the domain model for the list component.
package model

import (
	"github.com/phoenix-tui/phoenix/components/list/internal/domain/service"
	"github.com/phoenix-tui/phoenix/components/list/internal/domain/value"
)

// List is the aggregate root for the list component.
type List struct {
	items           []*value.Item       // All items
	filteredItems   []*value.Item       // Filtered items (if filter active)
	selectedIndices map[int]bool        // Selected item indices (in filtered list)
	focusedIndex    int                 // Currently focused item (in filtered list)
	selectionMode   value.SelectionMode // Single or Multi
	itemRenderer    func(item *value.Item, index int, selected, focused bool) string
	filterFunc      func(item *value.Item, query string) bool
	filterQuery     string // Current filter query
	height          int    // Visible height (for scrolling)
	scrollOffset    int    // Scroll offset

	// Services.
	navService    *service.NavigationService
	filterService *service.FilterService
}

// NewList creates a new list with the given selection mode.
func NewList(selectionMode value.SelectionMode) *List {
	return &List{
		items:           []*value.Item{},
		filteredItems:   []*value.Item{},
		selectedIndices: make(map[int]bool),
		focusedIndex:    0,
		selectionMode:   selectionMode,
		itemRenderer:    defaultItemRenderer,
		filterFunc:      nil,
		filterQuery:     "",
		height:          10, // Default height
		scrollOffset:    0,
		navService:      service.NewNavigationService(),
		filterService:   service.NewFilterService(),
	}
}

// NewListWithItems creates a new list with the given items and selection mode.
func NewListWithItems(items []*value.Item, mode value.SelectionMode) *List {
	l := NewList(mode)
	return l.WithItems(items)
}

// defaultItemRenderer provides a simple default rendering for items.
func defaultItemRenderer(item *value.Item, _ int, selected, focused bool) string {
	prefix := "  "
	if selected {
		prefix = "✓ "
	}
	if focused {
		prefix = "> "
	}
	return prefix + item.Label()
}

// WithItems returns a new List with the given items.
func (l *List) WithItems(items []*value.Item) *List {
	newList := l.clone()
	newList.items = make([]*value.Item, len(items))
	copy(newList.items, items)
	newList.applyFilter()
	newList.resetFocus()
	return newList
}

// WithItemRenderer returns a new List with a custom item renderer.
func (l *List) WithItemRenderer(renderer func(*value.Item, int, bool, bool) string) *List {
	newList := l.clone()
	newList.itemRenderer = renderer
	return newList
}

// WithFilter returns a new List with a custom filter function.
func (l *List) WithFilter(filterFunc func(*value.Item, string) bool) *List {
	newList := l.clone()
	newList.filterFunc = filterFunc
	newList.applyFilter()
	return newList
}

// WithHeight returns a new List with the specified visible height.
func (l *List) WithHeight(height int) *List {
	newList := l.clone()
	newList.height = height
	newList.updateScrollOffset()
	return newList
}

// MoveUp moves the focus up by one item.
func (l *List) MoveUp() *List {
	if len(l.filteredItems) == 0 {
		return l
	}
	newList := l.clone()
	newList.focusedIndex = newList.navService.MoveUp(newList.focusedIndex, len(newList.filteredItems))
	newList.updateScrollOffset()
	return newList
}

// MoveDown moves the focus down by one item.
func (l *List) MoveDown() *List {
	if len(l.filteredItems) == 0 {
		return l
	}
	newList := l.clone()
	newList.focusedIndex = newList.navService.MoveDown(newList.focusedIndex, len(newList.filteredItems))
	newList.updateScrollOffset()
	return newList
}

// MovePageUp moves the focus up by one page.
func (l *List) MovePageUp() *List {
	if len(l.filteredItems) == 0 {
		return l
	}
	newList := l.clone()
	newList.focusedIndex = newList.navService.MovePageUp(newList.focusedIndex, newList.height, len(newList.filteredItems))
	newList.updateScrollOffset()
	return newList
}

// MovePageDown moves the focus down by one page.
func (l *List) MovePageDown() *List {
	if len(l.filteredItems) == 0 {
		return l
	}
	newList := l.clone()
	newList.focusedIndex = newList.navService.MovePageDown(newList.focusedIndex, newList.height, len(newList.filteredItems))
	newList.updateScrollOffset()
	return newList
}

// MoveToStart moves the focus to the first item.
func (l *List) MoveToStart() *List {
	if len(l.filteredItems) == 0 {
		return l
	}
	newList := l.clone()
	newList.focusedIndex = newList.navService.MoveToStart()
	newList.updateScrollOffset()
	return newList
}

// MoveToEnd moves the focus to the last item.
func (l *List) MoveToEnd() *List {
	if len(l.filteredItems) == 0 {
		return l
	}
	newList := l.clone()
	newList.focusedIndex = newList.navService.MoveToEnd(len(newList.filteredItems))
	newList.updateScrollOffset()
	return newList
}

// ToggleSelection toggles the selection of the currently focused item.
func (l *List) ToggleSelection() *List {
	if len(l.filteredItems) == 0 {
		return l
	}

	newList := l.clone()

	//nolint:nestif // Selection logic is clear: single vs multi mode with toggle behavior
	if newList.selectionMode.IsSingle() {
		// Single selection: clear all and select focused.
		newList.selectedIndices = make(map[int]bool)
		newList.selectedIndices[newList.focusedIndex] = true
	} else {
		// Multi selection: toggle focused item.
		if newList.selectedIndices[newList.focusedIndex] {
			delete(newList.selectedIndices, newList.focusedIndex)
		} else {
			newList.selectedIndices[newList.focusedIndex] = true
		}
	}

	return newList
}

// SelectAll selects all items (multi-selection only).
func (l *List) SelectAll() *List {
	if l.selectionMode.IsSingle() || len(l.filteredItems) == 0 {
		return l
	}

	newList := l.clone()
	newList.selectedIndices = make(map[int]bool)
	for i := range newList.filteredItems {
		newList.selectedIndices[i] = true
	}

	return newList
}

// ClearSelection clears all selections.
func (l *List) ClearSelection() *List {
	newList := l.clone()
	newList.selectedIndices = make(map[int]bool)
	return newList
}

// SelectedItems returns the selected items.
func (l *List) SelectedItems() []*value.Item {
	result := make([]*value.Item, 0, len(l.selectedIndices))
	for idx := range l.selectedIndices {
		if idx < len(l.filteredItems) {
			result = append(result, l.filteredItems[idx])
		}
	}
	return result
}

// FocusedItem returns the currently focused item, or nil if no items.
func (l *List) FocusedItem() *value.Item {
	if len(l.filteredItems) == 0 || l.focusedIndex >= len(l.filteredItems) {
		return nil
	}
	return l.filteredItems[l.focusedIndex]
}

// SetFilterQuery updates the filter query and reapplies filtering.
func (l *List) SetFilterQuery(query string) *List {
	newList := l.clone()
	newList.filterQuery = query
	newList.applyFilter()
	newList.resetFocus()
	return newList
}

// ClearFilter clears the filter query.
func (l *List) ClearFilter() *List {
	return l.SetFilterQuery("")
}

// Items returns all items (not filtered).
func (l *List) Items() []*value.Item {
	result := make([]*value.Item, len(l.items))
	copy(result, l.items)
	return result
}

// FilteredItems returns the currently filtered items.
func (l *List) FilteredItems() []*value.Item {
	result := make([]*value.Item, len(l.filteredItems))
	copy(result, l.filteredItems)
	return result
}

// IsFiltered returns true if a filter is currently active.
func (l *List) IsFiltered() bool {
	return l.filterQuery != ""
}

// FocusedIndex returns the currently focused index (in filtered list).
func (l *List) FocusedIndex() int {
	return l.focusedIndex
}

// SelectedIndices returns the indices of selected items (in filtered list).
func (l *List) SelectedIndices() []int {
	result := make([]int, 0, len(l.selectedIndices))
	for idx := range l.selectedIndices {
		result = append(result, idx)
	}
	return result
}

// Height returns the visible height.
func (l *List) Height() int {
	return l.height
}

// ScrollOffset returns the current scroll offset.
func (l *List) ScrollOffset() int {
	return l.scrollOffset
}

// SelectionMode returns the selection mode.
func (l *List) SelectionMode() value.SelectionMode {
	return l.selectionMode
}

// RenderItem renders a specific item using the item renderer.
func (l *List) RenderItem(index int) string {
	if index < 0 || index >= len(l.filteredItems) {
		return ""
	}

	item := l.filteredItems[index]
	selected := l.selectedIndices[index]
	focused := index == l.focusedIndex

	return l.itemRenderer(item, index, selected, focused)
}

// RenderVisibleItems renders all visible items based on scroll offset.
func (l *List) RenderVisibleItems() []string {
	if len(l.filteredItems) == 0 {
		return []string{}
	}

	start := l.scrollOffset
	end := l.scrollOffset + l.height
	if end > len(l.filteredItems) {
		end = len(l.filteredItems)
	}
	if start >= len(l.filteredItems) {
		start = 0
	}

	result := make([]string, 0, end-start)
	for i := start; i < end; i++ {
		result = append(result, l.RenderItem(i))
	}

	return result
}

// clone creates a shallow copy of the list for immutability.
func (l *List) clone() *List {
	newSelectedIndices := make(map[int]bool, len(l.selectedIndices))
	for k, v := range l.selectedIndices {
		newSelectedIndices[k] = v
	}

	return &List{
		items:           l.items,
		filteredItems:   l.filteredItems,
		selectedIndices: newSelectedIndices,
		focusedIndex:    l.focusedIndex,
		selectionMode:   l.selectionMode,
		itemRenderer:    l.itemRenderer,
		filterFunc:      l.filterFunc,
		filterQuery:     l.filterQuery,
		height:          l.height,
		scrollOffset:    l.scrollOffset,
		navService:      l.navService,
		filterService:   l.filterService,
	}
}

// applyFilter applies the current filter to items.
func (l *List) applyFilter() {
	l.filteredItems = l.filterService.Filter(l.items, l.filterQuery, l.filterFunc)
	l.selectedIndices = make(map[int]bool) // Clear selections after filter
}

// resetFocus resets the focus to the first item if current focus is out of bounds.
func (l *List) resetFocus() {
	if l.focusedIndex >= len(l.filteredItems) {
		l.focusedIndex = 0
	}
	l.updateScrollOffset()
}

// updateScrollOffset updates the scroll offset based on the focused item.
func (l *List) updateScrollOffset() {
	l.scrollOffset = l.navService.CalculateScrollOffset(
		l.focusedIndex,
		l.height,
		len(l.filteredItems),
	)
}
