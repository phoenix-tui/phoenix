// Package model provides the domain model for the select component.
package model

import (
	"github.com/phoenix-tui/phoenix/components/select/internal/domain/value"
)

// Select represents the domain model for a single-choice selection list.
// It follows rich domain model pattern with encapsulated behavior.
type Select[T any] struct {
	options       []*value.Option[T]
	cursor        int
	selectedIndex int
	filterQuery   string
	filteredOpts  []*value.Option[T]
	height        int
	scrollOffset  int
	filterFunc    FilterFunc[T]
	renderFunc    RenderFunc[T]
}

// FilterFunc is a function that determines if an option matches the query.
type FilterFunc[T any] func(opt *value.Option[T], query string) bool

// RenderFunc is a function that renders an option as a string.
type RenderFunc[T any] func(opt *value.Option[T], index int, focused bool) string

// New creates a new Select with the given options.
func New[T any](options []*value.Option[T]) *Select[T] {
	if len(options) == 0 {
		options = make([]*value.Option[T], 0)
	}

	return &Select[T]{
		options:       options,
		cursor:        0,
		selectedIndex: -1,
		filterQuery:   "",
		filteredOpts:  options,
		height:        10,
		scrollOffset:  0,
		filterFunc:    defaultFilter[T],
		renderFunc:    defaultRender[T],
	}
}

// WithHeight returns a new Select with the specified visible height.
func (s *Select[T]) WithHeight(height int) *Select[T] {
	if height < 1 {
		height = 1
	}
	return &Select[T]{
		options:       s.options,
		cursor:        s.cursor,
		selectedIndex: s.selectedIndex,
		filterQuery:   s.filterQuery,
		filteredOpts:  s.filteredOpts,
		height:        height,
		scrollOffset:  s.scrollOffset,
		filterFunc:    s.filterFunc,
		renderFunc:    s.renderFunc,
	}
}

// WithFilterFunc returns a new Select with the specified filter function.
func (s *Select[T]) WithFilterFunc(fn FilterFunc[T]) *Select[T] {
	return &Select[T]{
		options:       s.options,
		cursor:        s.cursor,
		selectedIndex: s.selectedIndex,
		filterQuery:   s.filterQuery,
		filteredOpts:  s.filteredOpts,
		height:        s.height,
		scrollOffset:  s.scrollOffset,
		filterFunc:    fn,
		renderFunc:    s.renderFunc,
	}
}

// WithRenderFunc returns a new Select with the specified render function.
func (s *Select[T]) WithRenderFunc(fn RenderFunc[T]) *Select[T] {
	return &Select[T]{
		options:       s.options,
		cursor:        s.cursor,
		selectedIndex: s.selectedIndex,
		filterQuery:   s.filterQuery,
		filteredOpts:  s.filteredOpts,
		height:        s.height,
		scrollOffset:  s.scrollOffset,
		filterFunc:    s.filterFunc,
		renderFunc:    fn,
	}
}

// MoveUp moves the cursor up one position.
func (s *Select[T]) MoveUp() *Select[T] {
	newCursor := s.cursor
	if newCursor > 0 {
		newCursor--
	}
	return s.withCursor(newCursor)
}

// MoveDown moves the cursor down one position.
func (s *Select[T]) MoveDown() *Select[T] {
	newCursor := s.cursor
	maxIndex := len(s.filteredOpts) - 1
	if newCursor < maxIndex {
		newCursor++
	}
	return s.withCursor(newCursor)
}

// MoveToStart moves the cursor to the first option.
func (s *Select[T]) MoveToStart() *Select[T] {
	return s.withCursor(0)
}

// MoveToEnd moves the cursor to the last option.
func (s *Select[T]) MoveToEnd() *Select[T] {
	maxIndex := len(s.filteredOpts) - 1
	if maxIndex < 0 {
		maxIndex = 0
	}
	return s.withCursor(maxIndex)
}

// Select marks the currently focused option as selected.
func (s *Select[T]) Select() *Select[T] {
	if s.cursor >= 0 && s.cursor < len(s.filteredOpts) {
		// Find the original index in options array
		selectedOpt := s.filteredOpts[s.cursor]
		for i, opt := range s.options {
			if opt == selectedOpt {
				return s.withSelectedIndex(i)
			}
		}
	}
	return s
}

// SetFilterQuery sets the filter query and updates filtered options.
func (s *Select[T]) SetFilterQuery(query string) *Select[T] {
	// Filter options based on query
	filtered := make([]*value.Option[T], 0)
	for _, opt := range s.options {
		if s.filterFunc(opt, query) {
			filtered = append(filtered, opt)
		}
	}

	newCursor := s.cursor
	if newCursor >= len(filtered) {
		newCursor = len(filtered) - 1
	}
	if newCursor < 0 {
		newCursor = 0
	}

	return &Select[T]{
		options:       s.options,
		cursor:        newCursor,
		selectedIndex: s.selectedIndex,
		filterQuery:   query,
		filteredOpts:  filtered,
		height:        s.height,
		scrollOffset:  s.scrollOffset,
		filterFunc:    s.filterFunc,
		renderFunc:    s.renderFunc,
	}
}

// ClearFilter clears the filter query and shows all options.
func (s *Select[T]) ClearFilter() *Select[T] {
	return s.SetFilterQuery("")
}

// SelectedValue returns the currently selected value, or zero value if nothing selected.
func (s *Select[T]) SelectedValue() (T, bool) {
	if s.selectedIndex >= 0 && s.selectedIndex < len(s.options) {
		return s.options[s.selectedIndex].Value(), true
	}
	var zero T
	return zero, false
}

// FocusedValue returns the currently focused value, or zero value if no options.
func (s *Select[T]) FocusedValue() (T, bool) {
	if s.cursor >= 0 && s.cursor < len(s.filteredOpts) {
		return s.filteredOpts[s.cursor].Value(), true
	}
	var zero T
	return zero, false
}

// RenderVisibleOptions returns the rendered strings for visible options.
func (s *Select[T]) RenderVisibleOptions() []string {
	if len(s.filteredOpts) == 0 {
		return []string{}
	}

	// Adjust scroll offset to keep cursor visible
	newS := s.adjustScrollOffset()

	// Calculate visible range
	start := newS.scrollOffset
	end := start + newS.height
	if end > len(newS.filteredOpts) {
		end = len(newS.filteredOpts)
	}

	// Render visible options
	result := make([]string, end-start)
	for i := start; i < end; i++ {
		focused := (i == newS.cursor)
		result[i-start] = newS.renderFunc(newS.filteredOpts[i], i, focused)
	}

	return result
}

// FilterQuery returns the current filter query.
func (s *Select[T]) FilterQuery() string {
	return s.filterQuery
}

// IsFiltered returns true if a filter is currently active.
func (s *Select[T]) IsFiltered() bool {
	return s.filterQuery != ""
}

// withCursor returns a new Select with the specified cursor position.
func (s *Select[T]) withCursor(cursor int) *Select[T] {
	return &Select[T]{
		options:       s.options,
		cursor:        cursor,
		selectedIndex: s.selectedIndex,
		filterQuery:   s.filterQuery,
		filteredOpts:  s.filteredOpts,
		height:        s.height,
		scrollOffset:  s.scrollOffset,
		filterFunc:    s.filterFunc,
		renderFunc:    s.renderFunc,
	}
}

// withSelectedIndex returns a new Select with the specified selected index.
func (s *Select[T]) withSelectedIndex(index int) *Select[T] {
	return &Select[T]{
		options:       s.options,
		cursor:        s.cursor,
		selectedIndex: index,
		filterQuery:   s.filterQuery,
		filteredOpts:  s.filteredOpts,
		height:        s.height,
		scrollOffset:  s.scrollOffset,
		filterFunc:    s.filterFunc,
		renderFunc:    s.renderFunc,
	}
}

// adjustScrollOffset returns a new Select with scroll offset adjusted to keep cursor visible.
func (s *Select[T]) adjustScrollOffset() *Select[T] {
	newOffset := s.scrollOffset

	// Cursor above visible area
	if s.cursor < newOffset {
		newOffset = s.cursor
	}

	// Cursor below visible area
	if s.cursor >= newOffset+s.height {
		newOffset = s.cursor - s.height + 1
	}

	if newOffset == s.scrollOffset {
		return s
	}

	return &Select[T]{
		options:       s.options,
		cursor:        s.cursor,
		selectedIndex: s.selectedIndex,
		filterQuery:   s.filterQuery,
		filteredOpts:  s.filteredOpts,
		height:        s.height,
		scrollOffset:  newOffset,
		filterFunc:    s.filterFunc,
		renderFunc:    s.renderFunc,
	}
}
