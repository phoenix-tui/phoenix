// Package model provides the domain model for the multiselect component.
package model

import (
	"github.com/phoenix-tui/phoenix/components/multiselect/internal/domain/value"
)

// MultiSelect represents the domain model for a multi-choice selection list.
// It follows rich domain model pattern with encapsulated behavior.
type MultiSelect[T any] struct {
	options      []*value.Option[T]
	cursor       int
	selection    *value.Selection
	filterQuery  string
	filteredOpts []*value.Option[T]
	height       int
	scrollOffset int
	filterFunc   FilterFunc[T]
	renderFunc   RenderFunc[T]
}

// FilterFunc is a function that determines if an option matches the query.
type FilterFunc[T any] func(opt *value.Option[T], query string) bool

// RenderFunc is a function that renders an option as a string.
type RenderFunc[T any] func(opt *value.Option[T], index int, focused bool, selected bool) string

// New creates a new MultiSelect with the given options.
func New[T any](options []*value.Option[T], minCount, maxCount int) *MultiSelect[T] {
	if len(options) == 0 {
		options = make([]*value.Option[T], 0)
	}

	return &MultiSelect[T]{
		options:      options,
		cursor:       0,
		selection:    value.NewSelection(minCount, maxCount),
		filterQuery:  "",
		filteredOpts: options,
		height:       10,
		scrollOffset: 0,
		filterFunc:   defaultFilter[T],
		renderFunc:   defaultRender[T],
	}
}

// WithHeight returns a new MultiSelect with the specified visible height.
func (m *MultiSelect[T]) WithHeight(height int) *MultiSelect[T] {
	if height < 1 {
		height = 1
	}
	return &MultiSelect[T]{
		options:      m.options,
		cursor:       m.cursor,
		selection:    m.selection,
		filterQuery:  m.filterQuery,
		filteredOpts: m.filteredOpts,
		height:       height,
		scrollOffset: m.scrollOffset,
		filterFunc:   m.filterFunc,
		renderFunc:   m.renderFunc,
	}
}

// WithFilterFunc returns a new MultiSelect with the specified filter function.
func (m *MultiSelect[T]) WithFilterFunc(fn FilterFunc[T]) *MultiSelect[T] {
	return &MultiSelect[T]{
		options:      m.options,
		cursor:       m.cursor,
		selection:    m.selection,
		filterQuery:  m.filterQuery,
		filteredOpts: m.filteredOpts,
		height:       m.height,
		scrollOffset: m.scrollOffset,
		filterFunc:   fn,
		renderFunc:   m.renderFunc,
	}
}

// WithRenderFunc returns a new MultiSelect with the specified render function.
func (m *MultiSelect[T]) WithRenderFunc(fn RenderFunc[T]) *MultiSelect[T] {
	return &MultiSelect[T]{
		options:      m.options,
		cursor:       m.cursor,
		selection:    m.selection,
		filterQuery:  m.filterQuery,
		filteredOpts: m.filteredOpts,
		height:       m.height,
		scrollOffset: m.scrollOffset,
		filterFunc:   m.filterFunc,
		renderFunc:   fn,
	}
}

// WithSelected returns a new MultiSelect with the specified indices pre-selected.
func (m *MultiSelect[T]) WithSelected(indices ...int) *MultiSelect[T] {
	return &MultiSelect[T]{
		options:      m.options,
		cursor:       m.cursor,
		selection:    m.selection.WithSelected(indices...),
		filterQuery:  m.filterQuery,
		filteredOpts: m.filteredOpts,
		height:       m.height,
		scrollOffset: m.scrollOffset,
		filterFunc:   m.filterFunc,
		renderFunc:   m.renderFunc,
	}
}

// MoveUp moves the cursor up one position.
func (m *MultiSelect[T]) MoveUp() *MultiSelect[T] {
	newCursor := m.cursor
	if newCursor > 0 {
		newCursor--
	}
	return m.withCursor(newCursor)
}

// MoveDown moves the cursor down one position.
func (m *MultiSelect[T]) MoveDown() *MultiSelect[T] {
	newCursor := m.cursor
	maxIndex := len(m.filteredOpts) - 1
	if newCursor < maxIndex {
		newCursor++
	}
	return m.withCursor(newCursor)
}

// MoveToStart moves the cursor to the first option.
func (m *MultiSelect[T]) MoveToStart() *MultiSelect[T] {
	return m.withCursor(0)
}

// MoveToEnd moves the cursor to the last option.
func (m *MultiSelect[T]) MoveToEnd() *MultiSelect[T] {
	maxIndex := len(m.filteredOpts) - 1
	if maxIndex < 0 {
		maxIndex = 0
	}
	return m.withCursor(maxIndex)
}

// Toggle toggles the selection of the currently focused option.
func (m *MultiSelect[T]) Toggle() *MultiSelect[T] {
	if m.cursor >= 0 && m.cursor < len(m.filteredOpts) {
		// Find the original index in options array
		focusedOpt := m.filteredOpts[m.cursor]
		for i, opt := range m.options {
			if opt == focusedOpt {
				return m.withSelection(m.selection.Toggle(i))
			}
		}
	}
	return m
}

// SelectAll selects all filtered options (respecting max constraint).
func (m *MultiSelect[T]) SelectAll() *MultiSelect[T] {
	// Build list of indices for filtered options
	indices := make([]int, 0, len(m.filteredOpts))
	for _, filteredOpt := range m.filteredOpts {
		for i, opt := range m.options {
			if opt == filteredOpt {
				indices = append(indices, i)
				break
			}
		}
	}

	// Start with current selection and add all filtered indices
	newSelection := m.selection
	for _, idx := range indices {
		if !newSelection.IsSelected(idx) && newSelection.CanSelect() {
			newSelection = newSelection.Toggle(idx)
		}
	}

	return m.withSelection(newSelection)
}

// SelectNone clears all selections.
func (m *MultiSelect[T]) SelectNone() *MultiSelect[T] {
	return m.withSelection(m.selection.Clear())
}

// SetFilterQuery sets the filter query and updates filtered options.
func (m *MultiSelect[T]) SetFilterQuery(query string) *MultiSelect[T] {
	// Filter options based on query
	filtered := make([]*value.Option[T], 0)
	for _, opt := range m.options {
		if m.filterFunc(opt, query) {
			filtered = append(filtered, opt)
		}
	}

	newCursor := m.cursor
	if newCursor >= len(filtered) {
		newCursor = len(filtered) - 1
	}
	if newCursor < 0 {
		newCursor = 0
	}

	return &MultiSelect[T]{
		options:      m.options,
		cursor:       newCursor,
		selection:    m.selection,
		filterQuery:  query,
		filteredOpts: filtered,
		height:       m.height,
		scrollOffset: m.scrollOffset,
		filterFunc:   m.filterFunc,
		renderFunc:   m.renderFunc,
	}
}

// ClearFilter clears the filter query and shows all options.
func (m *MultiSelect[T]) ClearFilter() *MultiSelect[T] {
	return m.SetFilterQuery("")
}

// SelectedValues returns the currently selected values.
func (m *MultiSelect[T]) SelectedValues() []T {
	indices := m.selection.Indices()
	result := make([]T, 0, len(indices))
	for _, idx := range indices {
		if idx >= 0 && idx < len(m.options) {
			result = append(result, m.options[idx].Value())
		}
	}
	return result
}

// SelectedIndices returns the currently selected indices.
func (m *MultiSelect[T]) SelectedIndices() []int {
	return m.selection.Indices()
}

// SelectionCount returns the number of selected items.
func (m *MultiSelect[T]) SelectionCount() int {
	return m.selection.Count()
}

// CanConfirm returns true if enough items are selected to confirm.
func (m *MultiSelect[T]) CanConfirm() bool {
	return m.selection.CanConfirm()
}

// RenderVisibleOptions returns the rendered strings for visible options.
func (m *MultiSelect[T]) RenderVisibleOptions() []string {
	if len(m.filteredOpts) == 0 {
		return []string{}
	}

	// Adjust scroll offset to keep cursor visible
	newM := m.adjustScrollOffset()

	// Calculate visible range
	start := newM.scrollOffset
	end := start + newM.height
	if end > len(newM.filteredOpts) {
		end = len(newM.filteredOpts)
	}

	// Render visible options
	result := make([]string, end-start)
	for i := start; i < end; i++ {
		focused := (i == newM.cursor)
		// Find original index to check selection
		selected := false
		for origIdx, opt := range newM.options {
			if opt == newM.filteredOpts[i] {
				selected = newM.selection.IsSelected(origIdx)
				break
			}
		}
		result[i-start] = newM.renderFunc(newM.filteredOpts[i], i, focused, selected)
	}

	return result
}

// FilterQuery returns the current filter query.
func (m *MultiSelect[T]) FilterQuery() string {
	return m.filterQuery
}

// IsFiltered returns true if a filter is currently active.
func (m *MultiSelect[T]) IsFiltered() bool {
	return m.filterQuery != ""
}

// TotalCount returns the total number of options (before filtering).
func (m *MultiSelect[T]) TotalCount() int {
	return len(m.options)
}

// FilteredCount returns the number of filtered options.
func (m *MultiSelect[T]) FilteredCount() int {
	return len(m.filteredOpts)
}

// withCursor returns a new MultiSelect with the specified cursor position.
func (m *MultiSelect[T]) withCursor(cursor int) *MultiSelect[T] {
	return &MultiSelect[T]{
		options:      m.options,
		cursor:       cursor,
		selection:    m.selection,
		filterQuery:  m.filterQuery,
		filteredOpts: m.filteredOpts,
		height:       m.height,
		scrollOffset: m.scrollOffset,
		filterFunc:   m.filterFunc,
		renderFunc:   m.renderFunc,
	}
}

// withSelection returns a new MultiSelect with the specified selection.
func (m *MultiSelect[T]) withSelection(selection *value.Selection) *MultiSelect[T] {
	return &MultiSelect[T]{
		options:      m.options,
		cursor:       m.cursor,
		selection:    selection,
		filterQuery:  m.filterQuery,
		filteredOpts: m.filteredOpts,
		height:       m.height,
		scrollOffset: m.scrollOffset,
		filterFunc:   m.filterFunc,
		renderFunc:   m.renderFunc,
	}
}

// adjustScrollOffset returns a new MultiSelect with scroll offset adjusted to keep cursor visible.
func (m *MultiSelect[T]) adjustScrollOffset() *MultiSelect[T] {
	newOffset := m.scrollOffset

	// Cursor above visible area
	if m.cursor < newOffset {
		newOffset = m.cursor
	}

	// Cursor below visible area
	if m.cursor >= newOffset+m.height {
		newOffset = m.cursor - m.height + 1
	}

	if newOffset == m.scrollOffset {
		return m
	}

	return &MultiSelect[T]{
		options:      m.options,
		cursor:       m.cursor,
		selection:    m.selection,
		filterQuery:  m.filterQuery,
		filteredOpts: m.filteredOpts,
		height:       m.height,
		scrollOffset: newOffset,
		filterFunc:   m.filterFunc,
		renderFunc:   m.renderFunc,
	}
}

// WithSelectionConstraints returns a new MultiSelect with updated min/max constraints.
func (m *MultiSelect[T]) WithSelectionConstraints(minCount, maxCount int) *MultiSelect[T] {
	// Get current selection indices
	currentIndices := m.selection.Indices()
	
	// Create new selection with new constraints and preserve selections
	newSelection := value.NewSelection(minCount, maxCount).WithSelected(currentIndices...)
	
	return &MultiSelect[T]{
		options:      m.options,
		cursor:       m.cursor,
		selection:    newSelection,
		filterQuery:  m.filterQuery,
		filteredOpts: m.filteredOpts,
		height:       m.height,
		scrollOffset: m.scrollOffset,
		filterFunc:   m.filterFunc,
		renderFunc:   m.renderFunc,
	}
}
