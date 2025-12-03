// Package table provides a universal table component for Phoenix TUI Framework.
//
// The Table component displays tabular data with support for:
//   - Column definitions (width, alignment, custom rendering)
//   - Sorting (by column, ascending/descending)
//   - Keyboard navigation (arrows, vim keys, home/end)
//   - Scrolling (for tables larger than viewport)
//
// This is a UNIVERSAL component - it works for any application (file managers,.
// data viewers, process lists, etc.). It does NOT include application-specific.
// features like file system operations.
package table

import (
	"fmt"
	"strings"

	model2 "github.com/phoenix-tui/phoenix/components/table/internal/domain/model"
	"github.com/phoenix-tui/phoenix/components/table/internal/domain/service"
	value2 "github.com/phoenix-tui/phoenix/components/table/internal/domain/value"
	"github.com/phoenix-tui/phoenix/components/table/internal/infrastructure"
	"github.com/phoenix-tui/phoenix/style"
	tea "github.com/phoenix-tui/phoenix/tea"
)

// Column defines a table column.
type Column struct {
	Key       string                   // Unique identifier
	Title     string                   // Display title
	Width     int                      // Column width in characters
	Alignment value2.Alignment         // Cell alignment (default: left)
	Sortable  bool                     // Can this column be sorted?
	Renderer  func(interface{}) string // Custom cell renderer (optional)
}

// Row represents a table row as a map of column key to cell value.
type Row map[string]interface{}

// Table is the public API for the table component.
// It implements tea.Model for integration with Phoenix Tea event loop.
type Table struct {
	domain      *model2.Table
	sortService *service.SortService
	keyBindings infrastructure.KeyBindings
	theme       *style.Theme // Optional theme, defaults to DefaultTheme if nil
}

// New creates a new table with the given columns.
func New(columns []Column) *Table {
	domainCols := make([]*model2.Column, len(columns))
	for i, col := range columns {
		domainCol := model2.NewColumn(col.Key, col.Title, col.Width)
		if col.Alignment != value2.AlignmentLeft {
			domainCol = domainCol.WithAlignment(col.Alignment)
		}
		if col.Sortable {
			domainCol = domainCol.WithSortable(true)
		}
		if col.Renderer != nil {
			domainCol = domainCol.WithRenderer(col.Renderer)
		}
		domainCols[i] = domainCol
	}

	return &Table{
		domain:      model2.NewTable(domainCols),
		sortService: service.NewSortService(),
		keyBindings: infrastructure.DefaultKeyBindings(),
	}
}

// NewWithRows creates a new table with columns and initial data.
func NewWithRows(columns []Column, rows []Row) *Table {
	t := New(columns)

	// Convert to domain rows.
	domainRows := make([]model2.Row, len(rows))
	for i, row := range rows {
		domainRows[i] = model2.Row(row)
	}

	t.domain = t.domain.WithRows(domainRows)
	return t
}

// Height returns a new table with the specified visible height.
func (t *Table) Height(height int) *Table {
	return &Table{
		domain:      t.domain.WithHeight(height),
		sortService: t.sortService,
		keyBindings: t.keyBindings,
		theme:       t.theme,
	}
}

// ShowHeader returns a new table with header visibility set.
func (t *Table) ShowHeader(show bool) *Table {
	return &Table{
		domain:      t.domain.WithShowHeader(show),
		sortService: t.sortService,
		keyBindings: t.keyBindings,
		theme:       t.theme,
	}
}

// Theme sets the theme for styling the table component.
// If nil is provided, DefaultTheme will be used during rendering.
func (t *Table) Theme(theme *style.Theme) *Table {
	return &Table{
		domain:      t.domain,
		sortService: t.sortService,
		keyBindings: t.keyBindings,
		theme:       theme,
	}
}

// SetRows returns a new table with updated rows (clears sorting).
func (t *Table) SetRows(rows []Row) *Table {
	domainRows := make([]model2.Row, len(rows))
	for i, row := range rows {
		domainRows[i] = model2.Row(row)
	}
	return &Table{
		domain:      t.domain.WithRows(domainRows),
		sortService: t.sortService,
		keyBindings: t.keyBindings,
	}
}

// KeyBindings returns a new table with custom key bindings.
func (t *Table) KeyBindings(kb infrastructure.KeyBindings) *Table {
	return &Table{
		domain:      t.domain,
		sortService: t.sortService,
		keyBindings: kb,
	}
}

// Init implements tea.Model.
func (t *Table) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model pattern.
func (t *Table) Update(msg tea.Msg) (*Table, tea.Cmd) {
	//nolint:gocritic // Type switch with single case is idiomatic for tea.Model pattern
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return t.handleKeyPress(msg), nil
	}
	return t, nil
}

// handleKeyPress processes keyboard input.
func (t *Table) handleKeyPress(msg tea.KeyMsg) *Table {
	kb := t.keyBindings

	var newDomain *model2.Table

	switch {
	case kb.IsUp(msg):
		newDomain = t.domain.MoveUp()
	case kb.IsDown(msg):
		newDomain = t.domain.MoveDown()
	case kb.IsHome(msg):
		newDomain = t.domain.MoveToStart()
	case kb.IsEnd(msg):
		newDomain = t.domain.MoveToEnd()
	case kb.IsPageUp(msg):
		return t.pageUp()
	case kb.IsPageDown(msg):
		return t.pageDown()
	case kb.IsClearSort(msg):
		newDomain = t.domain.ClearSort()
	default:
		return t
	}

	return &Table{
		domain:      newDomain,
		sortService: t.sortService,
		keyBindings: t.keyBindings,
	}
}

// pageUp moves up by one page (visible height).
func (t *Table) pageUp() *Table {
	pageSize := t.domain.Height()
	if t.domain.ShowHeader() {
		pageSize--
	}
	newDomain := t.domain
	for i := 0; i < pageSize && newDomain.SelectedIndex() > 0; i++ {
		newDomain = newDomain.MoveUp()
	}
	return &Table{
		domain:      newDomain,
		sortService: t.sortService,
		keyBindings: t.keyBindings,
	}
}

// pageDown moves down by one page (visible height).
func (t *Table) pageDown() *Table {
	pageSize := t.domain.Height()
	if t.domain.ShowHeader() {
		pageSize--
	}
	maxIndex := len(t.domain.Rows()) - 1
	newDomain := t.domain
	for i := 0; i < pageSize && newDomain.SelectedIndex() < maxIndex; i++ {
		newDomain = newDomain.MoveDown()
	}
	return &Table{
		domain:      newDomain,
		sortService: t.sortService,
		keyBindings: t.keyBindings,
	}
}

// View implements tea.Model.
//
//nolint:gocognit,gocyclo,cyclop,funlen // View method is complex by nature: renders header, rows, cells with sorting/styling/padding
func (t *Table) View() string {
	var b strings.Builder

	columns := t.domain.Columns()
	visibleRows := t.domain.VisibleRows()

	// Calculate total width.
	totalWidth := 0
	for _, col := range columns {
		totalWidth += col.Width() + 1 // +1 for separator
	}

	// Render header.
	//nolint:nestif // Header rendering is complex: sort indicators, padding, alignment per column
	if t.domain.ShowHeader() {
		for i, col := range columns {
			title := col.Title()

			// Add sort indicator if this column is sorted.
			if t.domain.IsSorted() && t.domain.SortColumn() == col.Key() {
				if t.domain.SortDirection().IsAscending() {
					title += " ▲"
				} else {
					title += " ▼"
				}
			}

			cell := t.formatCell(title, col.Width(), col.Alignment())
			b.WriteString(cell)

			if i < len(columns)-1 {
				b.WriteString("│")
			}
		}
		b.WriteString("\n")

		// Header separator.
		for i, col := range columns {
			b.WriteString(strings.Repeat("─", col.Width()))
			if i < len(columns)-1 {
				b.WriteString("┼")
			}
		}
		b.WriteString("\n")
	}

	// Render rows.
	selectedIndex := t.domain.SelectedIndex()
	scrollOffset := t.domain.ScrollOffset()

	for rowIdx, row := range visibleRows {
		absoluteIdx := scrollOffset + rowIdx
		isSelected := absoluteIdx == selectedIndex

		for colIdx, col := range columns {
			value := row[col.Key()]

			// Use custom renderer if available.
			var cellText string
			if col.Renderer() != nil {
				cellText = col.Renderer()(value)
			} else {
				cellText = fmt.Sprintf("%v", value)
			}

			cell := t.formatCell(cellText, col.Width(), col.Alignment())

			// Add selection indicator.
			if isSelected && colIdx == 0 {
				if cell != "" {
					cell = ">" + cell[1:]
				}
			}

			b.WriteString(cell)

			if colIdx < len(columns)-1 {
				b.WriteString("│")
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

// formatCell formats a cell with alignment and width.
func (t *Table) formatCell(text string, width int, alignment value2.Alignment) string {
	// Truncate if too long.
	if len(text) > width {
		if width > 3 {
			text = text[:width-3] + "..."
		} else {
			text = text[:width]
		}
	}

	// Pad based on alignment.
	padding := width - len(text)
	if padding <= 0 {
		return text
	}

	switch {
	case alignment.IsLeft():
		return text + strings.Repeat(" ", padding)
	case alignment.IsRight():
		return strings.Repeat(" ", padding) + text
	case alignment.IsCenter():
		leftPad := padding / 2
		rightPad := padding - leftPad
		return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
	default:
		return text
	}
}

// SelectedRow returns the currently selected row.
func (t *Table) SelectedRow() Row {
	return Row(t.domain.SelectedRow())
}

// SelectedIndex returns the index of the selected row.
func (t *Table) SelectedIndex() int {
	return t.domain.SelectedIndex()
}

// Rows returns all rows (unsorted).
func (t *Table) Rows() []Row {
	domainRows := t.domain.Rows()
	rows := make([]Row, len(domainRows))
	for i, row := range domainRows {
		rows[i] = Row(row)
	}
	return rows
}

// SortByColumn returns a new table sorted by the specified column key.
// If already sorted by this column, toggles the direction.
func (t *Table) SortByColumn(columnKey string) *Table {
	// Check if column is sortable.
	var targetCol *model2.Column
	for _, col := range t.domain.Columns() {
		if col.Key() == columnKey {
			targetCol = col
			break
		}
	}

	if targetCol == nil || !targetCol.IsSortable() {
		return t // Column not found or not sortable
	}

	// Determine direction.
	direction := value2.SortDirectionAsc
	if t.domain.IsSorted() && t.domain.SortColumn() == columnKey {
		direction = t.domain.SortDirection().Toggle()
	}

	// Perform sort.
	sortedRows := t.sortService.Sort(t.domain.Rows(), columnKey, direction)
	newDomain := t.domain.SortBy(columnKey, direction, sortedRows)

	return &Table{
		domain:      newDomain,
		sortService: t.sortService,
		keyBindings: t.keyBindings,
	}
}

// ClearSort returns a new table with sorting removed.
func (t *Table) ClearSort() *Table {
	return &Table{
		domain:      t.domain.ClearSort(),
		sortService: t.sortService,
		keyBindings: t.keyBindings,
	}
}
