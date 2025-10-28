package table

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/components/table/internal/domain/value"
	tea "github.com/phoenix-tui/phoenix/tea"
)

func createTestTable() *Table {
	columns := []Column{
		{Key: "id", Title: "ID", Width: 5, Alignment: value.AlignmentRight, Sortable: true},
		{Key: "name", Title: "Name", Width: 15, Sortable: true},
		{Key: "age", Title: "Age", Width: 5, Alignment: value.AlignmentRight, Sortable: true},
	}

	rows := []Row{
		{"id": 1, "name": "Alice", "age": 30},
		{"id": 2, "name": "Bob", "age": 25},
		{"id": 3, "name": "Charlie", "age": 35},
	}

	return NewWithRows(columns, rows)
}

func TestNew(t *testing.T) {
	columns := []Column{
		{Key: "id", Title: "ID", Width: 10},
	}

	table := New(columns)

	if table == nil {
		t.Fatal("Table should not be nil")
	}
	if len(table.Rows()) != 0 {
		t.Errorf("New table should have 0 rows")
	}
}

func TestNewWithRows(t *testing.T) {
	table := createTestTable()

	if len(table.Rows()) != 3 {
		t.Errorf("Rows count = %v, want 3", len(table.Rows()))
	}
	if table.SelectedIndex() != 0 {
		t.Errorf("SelectedIndex = %v, want 0", table.SelectedIndex())
	}
}

func TestTable_Height(t *testing.T) {
	table := createTestTable()
	table = table.Height(20)

	if table.domain.Height() != 20 {
		t.Errorf("Height = %v, want 20", table.domain.Height())
	}
}

func TestTable_ShowHeader(t *testing.T) {
	table := createTestTable()
	table = table.ShowHeader(false)

	if table.domain.ShowHeader() {
		t.Errorf("ShowHeader should be false")
	}
}

func TestTable_SetRows(t *testing.T) {
	table := createTestTable()

	newRows := []Row{
		{"id": 10, "name": "NewUser", "age": 40},
	}

	table = table.SetRows(newRows)

	if len(table.Rows()) != 1 {
		t.Errorf("Rows count = %v, want 1", len(table.Rows()))
	}
	if table.Rows()[0]["name"] != "NewUser" {
		t.Errorf("First row name = %v, want NewUser", table.Rows()[0]["name"])
	}
}

func TestTable_Init(t *testing.T) {
	table := createTestTable()
	cmd := table.Init()

	if cmd != nil {
		t.Errorf("Init should return nil cmd")
	}
}

func TestTable_Update_MoveDown(t *testing.T) {
	table := createTestTable()

	msg := tea.KeyMsg{Type: tea.KeyDown}
	updated, _ := table.Update(msg)

	table = updated
	if table.SelectedIndex() != 1 {
		t.Errorf("SelectedIndex = %v, want 1", table.SelectedIndex())
	}
}

func TestTable_Update_MoveUp(t *testing.T) {
	table := createTestTable()
	table.domain = table.domain.MoveDown().MoveDown() // Move to index 2

	msg := tea.KeyMsg{Type: tea.KeyUp}
	updated, _ := table.Update(msg)

	table = updated
	if table.SelectedIndex() != 1 {
		t.Errorf("SelectedIndex = %v, want 1", table.SelectedIndex())
	}
}

func TestTable_Update_VimKeys(t *testing.T) {
	table := createTestTable()

	// Test 'j' (down)
	msg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'j'}
	updated, _ := table.Update(msg)
	table = updated

	if table.SelectedIndex() != 1 {
		t.Errorf("After 'j', SelectedIndex = %v, want 1", table.SelectedIndex())
	}

	// Test 'k' (up)
	msg = tea.KeyMsg{Type: tea.KeyRune, Rune: 'k'}
	updated, _ = table.Update(msg)
	table = updated

	if table.SelectedIndex() != 0 {
		t.Errorf("After 'k', SelectedIndex = %v, want 0", table.SelectedIndex())
	}
}

func TestTable_Update_Home(t *testing.T) {
	table := createTestTable()
	table.domain = table.domain.MoveDown().MoveDown() // Move to index 2

	msg := tea.KeyMsg{Type: tea.KeyHome}
	updated, _ := table.Update(msg)
	table = updated

	if table.SelectedIndex() != 0 {
		t.Errorf("SelectedIndex = %v, want 0", table.SelectedIndex())
	}
}

func TestTable_Update_End(t *testing.T) {
	table := createTestTable()

	msg := tea.KeyMsg{Type: tea.KeyEnd}
	updated, _ := table.Update(msg)
	table = updated

	if table.SelectedIndex() != 2 {
		t.Errorf("SelectedIndex = %v, want 2", table.SelectedIndex())
	}
}

func TestTable_Update_VimHomeEnd(t *testing.T) {
	table := createTestTable()

	// Test 'g' (home)
	msg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'g'}
	table.domain = table.domain.MoveToEnd()
	updated, _ := table.Update(msg)
	table = updated

	if table.SelectedIndex() != 0 {
		t.Errorf("After 'g', SelectedIndex = %v, want 0", table.SelectedIndex())
	}

	// Test 'G' (end)
	msg = tea.KeyMsg{Type: tea.KeyRune, Rune: 'G'}
	updated, _ = table.Update(msg)
	table = updated

	if table.SelectedIndex() != 2 {
		t.Errorf("After 'G', SelectedIndex = %v, want 2", table.SelectedIndex())
	}
}

func TestTable_SelectedRow(t *testing.T) {
	table := createTestTable()

	row := table.SelectedRow()
	if row["name"] != "Alice" {
		t.Errorf("Selected row name = %v, want Alice", row["name"])
	}

	table.domain = table.domain.MoveDown()
	row = table.SelectedRow()
	if row["name"] != "Bob" {
		t.Errorf("Selected row name = %v, want Bob", row["name"])
	}
}

func TestTable_SortByColumn(t *testing.T) {
	table := createTestTable()

	// Sort by age ascending.
	table = table.SortByColumn("age")

	// First row should be Bob (age 25)
	visibleRows := table.domain.VisibleRows()
	if visibleRows[0]["name"] != "Bob" {
		t.Errorf("After sort asc, first row = %v, want Bob", visibleRows[0]["name"])
	}

	// Toggle to descending.
	table = table.SortByColumn("age")

	visibleRows = table.domain.VisibleRows()
	if visibleRows[0]["name"] != "Charlie" {
		t.Errorf("After sort desc, first row = %v, want Charlie", visibleRows[0]["name"])
	}
}

func TestTable_SortByColumn_Unsortable(t *testing.T) {
	columns := []Column{
		{Key: "id", Title: "ID", Width: 5, Sortable: false}, // Not sortable
		{Key: "name", Title: "Name", Width: 15},
	}

	rows := []Row{
		{"id": 3, "name": "Charlie"},
		{"id": 1, "name": "Alice"},
	}

	table := NewWithRows(columns, rows)
	table = table.SortByColumn("id")

	// Should not be sorted (column not sortable)
	if table.domain.IsSorted() {
		t.Errorf("Table should not be sorted (column not sortable)")
	}
}

func TestTable_ClearSort(t *testing.T) {
	table := createTestTable()
	table = table.SortByColumn("name")

	if !table.domain.IsSorted() {
		t.Errorf("Table should be sorted")
	}

	table = table.ClearSort()

	if table.domain.IsSorted() {
		t.Errorf("Table should not be sorted after ClearSort")
	}
}

func TestTable_View_Header(t *testing.T) {
	table := createTestTable()
	view := table.View()

	// Should contain header titles.
	if !strings.Contains(view, "ID") {
		t.Errorf("View should contain 'ID' header")
	}
	if !strings.Contains(view, "Name") {
		t.Errorf("View should contain 'Name' header")
	}
	if !strings.Contains(view, "Age") {
		t.Errorf("View should contain 'Age' header")
	}

	// Should contain separator.
	if !strings.Contains(view, "─") {
		t.Errorf("View should contain header separator")
	}
}

func TestTable_View_NoHeader(t *testing.T) {
	table := createTestTable().ShowHeader(false)
	view := table.View()

	// Should not contain separator.
	if strings.Contains(view, "─") {
		t.Errorf("View should not contain header separator when header hidden")
	}
}

func TestTable_View_Data(t *testing.T) {
	table := createTestTable()
	view := table.View()

	// Should contain data.
	if !strings.Contains(view, "Alice") {
		t.Errorf("View should contain 'Alice'")
	}
	if !strings.Contains(view, "Bob") {
		t.Errorf("View should contain 'Bob'")
	}
	if !strings.Contains(view, "Charlie") {
		t.Errorf("View should contain 'Charlie'")
	}
}

func TestTable_View_Selection(t *testing.T) {
	table := createTestTable()
	view := table.View()

	// First row should have selection indicator.
	lines := strings.Split(view, "\n")
	var dataLines []string
	for _, line := range lines {
		if !strings.Contains(line, "─") && !strings.Contains(line, "ID") && line != "" {
			dataLines = append(dataLines, line)
		}
	}

	if len(dataLines) > 0 {
		firstDataLine := dataLines[0]
		if !strings.HasPrefix(firstDataLine, ">") {
			t.Errorf("First data row should have '>' indicator, got: %s", firstDataLine)
		}
	}
}

func TestTable_View_SortIndicator(t *testing.T) {
	table := createTestTable()
	table = table.SortByColumn("name")

	view := table.View()

	// Should contain ascending indicator.
	if !strings.Contains(view, "▲") {
		t.Errorf("View should contain ascending sort indicator '▲'")
	}

	// Toggle to descending.
	table = table.SortByColumn("name")
	view = table.View()

	// Should contain descending indicator.
	if !strings.Contains(view, "▼") {
		t.Errorf("View should contain descending sort indicator '▼'")
	}
}

func TestTable_View_CustomRenderer(t *testing.T) {
	columns := []Column{
		{
			Key:   "status",
			Title: "Status",
			Width: 10,
			Renderer: func(v interface{}) string {
				if v.(bool) {
					return "ACTIVE"
				}
				return "INACTIVE"
			},
		},
	}

	rows := []Row{
		{"status": true},
		{"status": false},
	}

	table := NewWithRows(columns, rows)
	view := table.View()

	if !strings.Contains(view, "ACTIVE") {
		t.Errorf("View should contain custom rendered 'ACTIVE', got: %s", view)
	}
	if !strings.Contains(view, "INACTIVE") {
		t.Errorf("View should contain custom rendered 'INACTIVE', got: %s", view)
	}
}

func TestTable_View_Alignment(t *testing.T) {
	columns := []Column{
		{Key: "left", Title: "Left", Width: 10, Alignment: value.AlignmentLeft},
		{Key: "center", Title: "Center", Width: 10, Alignment: value.AlignmentCenter},
		{Key: "right", Title: "Right", Width: 10, Alignment: value.AlignmentRight},
	}

	rows := []Row{
		{"left": "L", "center": "C", "right": "R"},
	}

	table := NewWithRows(columns, rows)
	view := table.View()

	// Just check that view renders without error.
	// Detailed alignment testing would require parsing the output.
	if view == "" {
		t.Errorf("View should not be empty")
	}
}

func TestTable_PageDown(t *testing.T) {
	// Create table with many rows.
	columns := []Column{
		{Key: "id", Title: "ID", Width: 5},
	}

	rows := []Row{}
	for i := 1; i <= 20; i++ {
		rows = append(rows, Row{"id": i})
	}

	table := NewWithRows(columns, rows).Height(5)

	// Page down.
	msg := tea.KeyMsg{Type: tea.KeyPgDown}
	updated, _ := table.Update(msg)
	table = updated

	// Should have moved down by page size (4, because 1 row for header)
	if table.SelectedIndex() < 4 {
		t.Errorf("After PageDown, SelectedIndex = %v, should be >= 4", table.SelectedIndex())
	}
}

func TestTable_PageUp(t *testing.T) {
	// Create table with many rows.
	columns := []Column{
		{Key: "id", Title: "ID", Width: 5},
	}

	rows := []Row{}
	for i := 1; i <= 20; i++ {
		rows = append(rows, Row{"id": i})
	}

	table := NewWithRows(columns, rows).Height(5)

	// Move to end first.
	table.domain = table.domain.MoveToEnd()

	// Page up.
	msg := tea.KeyMsg{Type: tea.KeyPgUp}
	updated, _ := table.Update(msg)
	table = updated

	// Should have moved up.
	if table.SelectedIndex() >= 19 {
		t.Errorf("After PageUp from end, SelectedIndex = %v, should be < 19", table.SelectedIndex())
	}
}

func TestTable_Immutability(t *testing.T) {
	table1 := createTestTable()

	// Perform operations.
	table2 := table1.Height(20).ShowHeader(false)

	// table1 should be unchanged.
	if table1.domain.Height() != 10 {
		t.Errorf("Original table height should be 10")
	}
	if !table1.domain.ShowHeader() {
		t.Errorf("Original table should show header")
	}

	// table2 should have changes.
	if table2.domain.Height() != 20 {
		t.Errorf("New table height should be 20")
	}
	if table2.domain.ShowHeader() {
		t.Errorf("New table should hide header")
	}
}
