package model

import (
	"testing"

	value2 "github.com/phoenix-tui/phoenix/components/table/internal/domain/value"
)

func createTestColumns() []*Column {
	return []*Column{
		NewColumn("id", "ID", 5),
		NewColumn("name", "Name", 20),
		NewColumn("age", "Age", 5).WithAlignment(value2.AlignmentRight),
	}
}

func createTestRows() []Row {
	return []Row{
		{"id": 1, "name": "Alice", "age": 30},
		{"id": 2, "name": "Bob", "age": 25},
		{"id": 3, "name": "Charlie", "age": 35},
		{"id": 4, "name": "Diana", "age": 28},
		{"id": 5, "name": "Eve", "age": 32},
	}
}

func TestNewTable(t *testing.T) {
	columns := createTestColumns()
	table := NewTable(columns)

	if len(table.Columns()) != 3 {
		t.Errorf("Columns count = %v, want 3", len(table.Columns()))
	}
	if len(table.Rows()) != 0 {
		t.Errorf("Rows count = %v, want 0", len(table.Rows()))
	}
	if table.Height() != 10 {
		t.Errorf("Height = %v, want 10", table.Height())
	}
	if !table.ShowHeader() {
		t.Errorf("ShowHeader should be true by default")
	}
	if table.IsSorted() {
		t.Errorf("Should not be sorted initially")
	}
}

func TestNewTableWithRows(t *testing.T) {
	columns := createTestColumns()
	rows := createTestRows()
	table := NewTableWithRows(columns, rows)

	if len(table.Rows()) != 5 {
		t.Errorf("Rows count = %v, want 5", len(table.Rows()))
	}
	if table.SelectedIndex() != 0 {
		t.Errorf("SelectedIndex = %v, want 0", table.SelectedIndex())
	}
}

func TestTable_WithRows(t *testing.T) {
	columns := createTestColumns()
	table := NewTable(columns)

	rows := createTestRows()
	table2 := table.WithRows(rows)

	// Original unchanged.
	if len(table.Rows()) != 0 {
		t.Errorf("Original table should have 0 rows")
	}

	// New table has rows.
	if len(table2.Rows()) != 5 {
		t.Errorf("New table rows = %v, want 5", len(table2.Rows()))
	}
}

func TestTable_WithHeight(t *testing.T) {
	table := NewTableWithRows(createTestColumns(), createTestRows())
	table2 := table.WithHeight(20)

	// Original unchanged.
	if table.Height() != 10 {
		t.Errorf("Original height should be 10")
	}

	// New table has new height.
	if table2.Height() != 20 {
		t.Errorf("New height = %v, want 20", table2.Height())
	}
}

func TestTable_WithShowHeader(t *testing.T) {
	table := NewTableWithRows(createTestColumns(), createTestRows())
	table2 := table.WithShowHeader(false)

	// Original unchanged.
	if !table.ShowHeader() {
		t.Errorf("Original should show header")
	}

	// New table hides header.
	if table2.ShowHeader() {
		t.Errorf("New table should hide header")
	}
}

func TestTable_MoveDown(t *testing.T) {
	table := NewTableWithRows(createTestColumns(), createTestRows())

	// Move down once.
	table2 := table.MoveDown()
	if table2.SelectedIndex() != 1 {
		t.Errorf("SelectedIndex = %v, want 1", table2.SelectedIndex())
	}

	// Move down to end.
	table3 := table2.MoveDown().MoveDown().MoveDown().MoveDown()
	if table3.SelectedIndex() != 4 {
		t.Errorf("SelectedIndex = %v, want 4", table3.SelectedIndex())
	}

	// Try to move past end.
	table4 := table3.MoveDown()
	if table4.SelectedIndex() != 4 {
		t.Errorf("Should not move past end, got %v", table4.SelectedIndex())
	}
}

func TestTable_MoveUp(t *testing.T) {
	table := NewTableWithRows(createTestColumns(), createTestRows())
	table = table.MoveToEnd()

	// Move up once.
	table2 := table.MoveUp()
	if table2.SelectedIndex() != 3 {
		t.Errorf("SelectedIndex = %v, want 3", table2.SelectedIndex())
	}

	// Move up to start.
	table3 := table2.MoveUp().MoveUp().MoveUp()
	if table3.SelectedIndex() != 0 {
		t.Errorf("SelectedIndex = %v, want 0", table3.SelectedIndex())
	}

	// Try to move before start.
	table4 := table3.MoveUp()
	if table4.SelectedIndex() != 0 {
		t.Errorf("Should not move before start, got %v", table4.SelectedIndex())
	}
}

func TestTable_MoveToStart(t *testing.T) {
	table := NewTableWithRows(createTestColumns(), createTestRows())
	table = table.MoveDown().MoveDown().MoveDown() // Move to index 3

	table2 := table.MoveToStart()
	if table2.SelectedIndex() != 0 {
		t.Errorf("SelectedIndex = %v, want 0", table2.SelectedIndex())
	}
	if table2.ScrollOffset() != 0 {
		t.Errorf("ScrollOffset = %v, want 0", table2.ScrollOffset())
	}
}

func TestTable_MoveToEnd(t *testing.T) {
	table := NewTableWithRows(createTestColumns(), createTestRows())

	table2 := table.MoveToEnd()
	if table2.SelectedIndex() != 4 {
		t.Errorf("SelectedIndex = %v, want 4", table2.SelectedIndex())
	}
}

func TestTable_SelectedRow(t *testing.T) {
	rows := createTestRows()
	table := NewTableWithRows(createTestColumns(), rows)

	// First row selected by default.
	row := table.SelectedRow()
	if row == nil {
		t.Fatal("SelectedRow should not be nil")
	}
	if row["name"] != "Alice" {
		t.Errorf("Selected row name = %v, want Alice", row["name"])
	}

	// Select second row.
	table2 := table.MoveDown()
	row2 := table2.SelectedRow()
	if row2["name"] != "Bob" {
		t.Errorf("Selected row name = %v, want Bob", row2["name"])
	}
}

func TestTable_SelectedRow_Empty(t *testing.T) {
	table := NewTable(createTestColumns())
	row := table.SelectedRow()
	if row != nil {
		t.Errorf("SelectedRow should be nil for empty table")
	}
}

func TestTable_SortBy(t *testing.T) {
	table := NewTableWithRows(createTestColumns(), createTestRows())

	// Simulate sorting (actual sorting done by service)
	sortedRows := []Row{
		{"id": 2, "name": "Bob", "age": 25},
		{"id": 1, "name": "Alice", "age": 30},
		{"id": 3, "name": "Charlie", "age": 35},
	}

	table2 := table.SortBy("name", value2.SortDirectionAsc, sortedRows)

	if !table2.IsSorted() {
		t.Errorf("Table should be sorted")
	}
	if table2.SortColumn() != "name" {
		t.Errorf("SortColumn = %v, want name", table2.SortColumn())
	}
	if !table2.SortDirection().IsAscending() {
		t.Errorf("SortDirection should be ascending")
	}
	if table2.SelectedIndex() != 0 {
		t.Errorf("Selection should reset to 0 when sorting")
	}
}

func TestTable_ClearSort(t *testing.T) {
	table := NewTableWithRows(createTestColumns(), createTestRows())

	sortedRows := []Row{{"id": 2, "name": "Bob", "age": 25}}
	table2 := table.SortBy("name", value2.SortDirectionAsc, sortedRows)
	table3 := table2.ClearSort()

	if table3.IsSorted() {
		t.Errorf("Table should not be sorted after ClearSort")
	}
	if table3.SortColumn() != "" {
		t.Errorf("SortColumn should be empty")
	}
	if !table3.SortDirection().IsNone() {
		t.Errorf("SortDirection should be None")
	}
}

func TestTable_VisibleRows(t *testing.T) {
	// Create table with 10 rows.
	rows := []Row{}
	for i := 1; i <= 10; i++ {
		rows = append(rows, Row{"id": i, "name": "User", "age": 20 + i})
	}

	table := NewTableWithRows(createTestColumns(), rows).WithHeight(5)

	// First 4 visible (5 height - 1 for header)
	visible := table.VisibleRows()
	if len(visible) != 4 {
		t.Errorf("VisibleRows count = %v, want 4", len(visible))
	}
	if visible[0]["id"] != 1 {
		t.Errorf("First visible row id = %v, want 1", visible[0]["id"])
	}

	// Scroll down.
	table2 := table.MoveDown().MoveDown().MoveDown().MoveDown().MoveDown()
	visible2 := table2.VisibleRows()
	if len(visible2) != 4 {
		t.Errorf("VisibleRows count = %v, want 4", len(visible2))
	}
	// Should show rows starting from scrollOffset.
	if visible2[0]["id"].(int) < 1 {
		t.Errorf("Visible rows should have scrolled")
	}
}

func TestTable_VisibleRows_NoHeader(t *testing.T) {
	rows := []Row{}
	for i := 1; i <= 10; i++ {
		rows = append(rows, Row{"id": i})
	}

	table := NewTableWithRows(createTestColumns(), rows).
		WithHeight(5).
		WithShowHeader(false)

	// All 5 rows visible (no header)
	visible := table.VisibleRows()
	if len(visible) != 5 {
		t.Errorf("VisibleRows count = %v, want 5", len(visible))
	}
}

func TestTable_Scrolling(t *testing.T) {
	// Create table with many rows.
	rows := []Row{}
	for i := 1; i <= 20; i++ {
		rows = append(rows, Row{"id": i})
	}

	table := NewTableWithRows(createTestColumns(), rows).WithHeight(5)

	// Move down several times.
	table2 := table
	for i := 0; i < 10; i++ {
		table2 = table2.MoveDown()
	}

	// Should have scrolled.
	if table2.ScrollOffset() <= 0 {
		t.Errorf("ScrollOffset should be > 0 after moving down")
	}

	// Selected index should be visible.
	if table2.SelectedIndex() < table2.ScrollOffset() {
		t.Errorf("Selected index should be >= scroll offset")
	}
}

func TestTable_Immutability(t *testing.T) {
	table1 := NewTableWithRows(createTestColumns(), createTestRows())

	// Chain multiple operations.
	table2 := table1.
		WithHeight(20).
		WithShowHeader(false).
		MoveDown().
		MoveDown()

	// Original unchanged.
	if table1.Height() != 10 {
		t.Errorf("Original height should be 10")
	}
	if !table1.ShowHeader() {
		t.Errorf("Original should show header")
	}
	if table1.SelectedIndex() != 0 {
		t.Errorf("Original selection should be 0")
	}

	// New table has all changes.
	if table2.Height() != 20 {
		t.Errorf("New height = %v, want 20", table2.Height())
	}
	if table2.ShowHeader() {
		t.Errorf("New table should hide header")
	}
	if table2.SelectedIndex() != 2 {
		t.Errorf("New selection = %v, want 2", table2.SelectedIndex())
	}
}
