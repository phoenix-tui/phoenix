package model

import "github.com/phoenix-tui/phoenix/components/table/domain/value"

// Row represents a table row as a map of column key to cell value.
type Row map[string]interface{}

// Table is the aggregate root for the table component.
// It manages columns, rows, sorting, selection, and scrolling.
type Table struct {
	columns       []*Column           // Column definitions
	rows          []Row               // Original data rows
	sortedRows    []Row               // Sorted rows (if sorting active)
	sortColumnKey string              // Currently sorted column key
	sortDirection value.SortDirection // Sort direction
	selectedIndex int                 // Selected row index
	scrollOffset  int                 // Scroll offset for viewport
	height        int                 // Visible height (number of rows)
	showHeader    bool                // Show header row?
}

// NewTable creates a new table with the given columns.
func NewTable(columns []*Column) *Table {
	return &Table{
		columns:       columns,
		rows:          []Row{},
		sortedRows:    nil,
		sortColumnKey: "",
		sortDirection: value.SortDirectionNone,
		selectedIndex: 0,
		scrollOffset:  0,
		height:        10, // Default height
		showHeader:    true,
	}
}

// NewTableWithRows creates a new table with columns and initial data.
func NewTableWithRows(columns []*Column, rows []Row) *Table {
	return &Table{
		columns:       columns,
		rows:          rows,
		sortedRows:    nil,
		sortColumnKey: "",
		sortDirection: value.SortDirectionNone,
		selectedIndex: 0,
		scrollOffset:  0,
		height:        10,
		showHeader:    true,
	}
}

// WithRows returns a new table with updated rows.
func (t *Table) WithRows(rows []Row) *Table {
	return &Table{
		columns:       t.columns,
		rows:          rows,
		sortedRows:    nil, // Clear sort when rows change
		sortColumnKey: "",
		sortDirection: value.SortDirectionNone,
		selectedIndex: 0,
		scrollOffset:  0,
		height:        t.height,
		showHeader:    t.showHeader,
	}
}

// WithHeight returns a new table with the specified visible height.
func (t *Table) WithHeight(height int) *Table {
	return &Table{
		columns:       t.columns,
		rows:          t.rows,
		sortedRows:    t.sortedRows,
		sortColumnKey: t.sortColumnKey,
		sortDirection: t.sortDirection,
		selectedIndex: t.selectedIndex,
		scrollOffset:  t.scrollOffset,
		height:        height,
		showHeader:    t.showHeader,
	}
}

// WithShowHeader returns a new table with header visibility set.
func (t *Table) WithShowHeader(show bool) *Table {
	return &Table{
		columns:       t.columns,
		rows:          t.rows,
		sortedRows:    t.sortedRows,
		sortColumnKey: t.sortColumnKey,
		sortDirection: t.sortDirection,
		selectedIndex: t.selectedIndex,
		scrollOffset:  t.scrollOffset,
		height:        t.height,
		showHeader:    show,
	}
}

// SortBy returns a new table sorted by the specified column and direction.
// Note: Actual sorting is delegated to SortService (domain service).
func (t *Table) SortBy(columnKey string, direction value.SortDirection, sortedRows []Row) *Table {
	return &Table{
		columns:       t.columns,
		rows:          t.rows,
		sortedRows:    sortedRows,
		sortColumnKey: columnKey,
		sortDirection: direction,
		selectedIndex: 0, // Reset selection when sorting
		scrollOffset:  0,
		height:        t.height,
		showHeader:    t.showHeader,
	}
}

// ClearSort returns a new table with sorting removed.
func (t *Table) ClearSort() *Table {
	return &Table{
		columns:       t.columns,
		rows:          t.rows,
		sortedRows:    nil,
		sortColumnKey: "",
		sortDirection: value.SortDirectionNone,
		selectedIndex: t.selectedIndex,
		scrollOffset:  t.scrollOffset,
		height:        t.height,
		showHeader:    t.showHeader,
	}
}

// MoveUp returns a new table with selection moved up one row.
func (t *Table) MoveUp() *Table {
	if t.selectedIndex <= 0 {
		return t // Already at top
	}

	newIndex := t.selectedIndex - 1
	newOffset := t.scrollOffset

	// Scroll up if needed.
	if newIndex < t.scrollOffset {
		newOffset = newIndex
	}

	return &Table{
		columns:       t.columns,
		rows:          t.rows,
		sortedRows:    t.sortedRows,
		sortColumnKey: t.sortColumnKey,
		sortDirection: t.sortDirection,
		selectedIndex: newIndex,
		scrollOffset:  newOffset,
		height:        t.height,
		showHeader:    t.showHeader,
	}
}

// MoveDown returns a new table with selection moved down one row.
func (t *Table) MoveDown() *Table {
	maxIndex := len(t.effectiveRows()) - 1
	if t.selectedIndex >= maxIndex {
		return t // Already at bottom
	}

	newIndex := t.selectedIndex + 1
	newOffset := t.scrollOffset

	// Scroll down if needed.
	visibleRows := t.height
	if t.showHeader {
		visibleRows-- // Header takes one row
	}
	if newIndex >= t.scrollOffset+visibleRows {
		newOffset = newIndex - visibleRows + 1
	}

	return &Table{
		columns:       t.columns,
		rows:          t.rows,
		sortedRows:    t.sortedRows,
		sortColumnKey: t.sortColumnKey,
		sortDirection: t.sortDirection,
		selectedIndex: newIndex,
		scrollOffset:  newOffset,
		height:        t.height,
		showHeader:    t.showHeader,
	}
}

// MoveToStart returns a new table with selection at the first row.
func (t *Table) MoveToStart() *Table {
	return &Table{
		columns:       t.columns,
		rows:          t.rows,
		sortedRows:    t.sortedRows,
		sortColumnKey: t.sortColumnKey,
		sortDirection: t.sortDirection,
		selectedIndex: 0,
		scrollOffset:  0,
		height:        t.height,
		showHeader:    t.showHeader,
	}
}

// MoveToEnd returns a new table with selection at the last row.
func (t *Table) MoveToEnd() *Table {
	maxIndex := len(t.effectiveRows()) - 1
	if maxIndex < 0 {
		maxIndex = 0
	}

	visibleRows := t.height
	if t.showHeader {
		visibleRows--
	}

	newOffset := maxIndex - visibleRows + 1
	if newOffset < 0 {
		newOffset = 0
	}

	return &Table{
		columns:       t.columns,
		rows:          t.rows,
		sortedRows:    t.sortedRows,
		sortColumnKey: t.sortColumnKey,
		sortDirection: t.sortDirection,
		selectedIndex: maxIndex,
		scrollOffset:  newOffset,
		height:        t.height,
		showHeader:    t.showHeader,
	}
}

// effectiveRows returns the rows to display (sorted if sorting is active).
func (t *Table) effectiveRows() []Row {
	if t.sortedRows != nil {
		return t.sortedRows
	}
	return t.rows
}

// SelectedRow returns the currently selected row.
func (t *Table) SelectedRow() Row {
	rows := t.effectiveRows()
	if t.selectedIndex >= 0 && t.selectedIndex < len(rows) {
		return rows[t.selectedIndex]
	}
	return nil
}

// SelectedIndex returns the index of the selected row.
func (t *Table) SelectedIndex() int {
	return t.selectedIndex
}

// Columns returns the column definitions.
func (t *Table) Columns() []*Column {
	return t.columns
}

// Rows returns the original rows (unsorted).
func (t *Table) Rows() []Row {
	return t.rows
}

// VisibleRows returns the rows currently visible in the viewport.
func (t *Table) VisibleRows() []Row {
	rows := t.effectiveRows()
	visibleRows := t.height
	if t.showHeader {
		visibleRows--
	}

	start := t.scrollOffset
	end := start + visibleRows

	if start >= len(rows) {
		return []Row{}
	}
	if end > len(rows) {
		end = len(rows)
	}

	return rows[start:end]
}

// IsSorted returns true if sorting is currently active.
func (t *Table) IsSorted() bool {
	return !t.sortDirection.IsNone()
}

// SortColumn returns the key of the currently sorted column.
func (t *Table) SortColumn() string {
	return t.sortColumnKey
}

// SortDirection returns the current sort direction.
func (t *Table) SortDirection() value.SortDirection {
	return t.sortDirection
}

// Height returns the visible height.
func (t *Table) Height() int {
	return t.height
}

// ShowHeader returns whether the header is visible.
func (t *Table) ShowHeader() bool {
	return t.showHeader
}

// ScrollOffset returns the current scroll offset.
func (t *Table) ScrollOffset() int {
	return t.scrollOffset
}
