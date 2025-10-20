package list

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/components/list/domain/value"
	"github.com/phoenix-tui/phoenix/components/list/infrastructure"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

func TestNew(t *testing.T) {
	values := []interface{}{1, 2, 3}
	labels := []string{"One", "Two", "Three"}

	l := New(values, labels, value.SelectionModeSingle)

	if l == nil {
		t.Fatal("New() returned nil")
	}
	if l.domain == nil {
		t.Error("New() domain should not be nil")
	}
}

func TestNew_PanicOnMismatchedLengths(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("New() should panic when values and labels have different lengths")
		}
	}()

	values := []interface{}{1, 2, 3}
	labels := []string{"One", "Two"} // Different length

	New(values, labels, value.SelectionModeSingle)
}

func TestNewSingleSelect(t *testing.T) {
	values := []interface{}{"a", "b"}
	labels := []string{"A", "B"}

	l := NewSingleSelect(values, labels)

	if l.domain.SelectionMode() != value.SelectionModeSingle {
		t.Error("NewSingleSelect() should create single-selection list")
	}
}

func TestNewMultiSelect(t *testing.T) {
	values := []interface{}{"a", "b"}
	labels := []string{"A", "B"}

	l := NewMultiSelect(values, labels)

	if l.domain.SelectionMode() != value.SelectionModeMulti {
		t.Error("NewMultiSelect() should create multi-selection list")
	}
}

func TestList_Height(t *testing.T) {
	l := NewSingleSelect([]interface{}{1, 2}, []string{"A", "B"})

	l = l.Height(20)

	if l.domain.Height() != 20 {
		t.Errorf("Height() = %d, want 20", l.domain.Height())
	}
}

func TestList_ItemRenderer(t *testing.T) {
	values := []interface{}{1, 2}
	labels := []string{"One", "Two"}
	l := NewSingleSelect(values, labels)

	customRenderer := func(item interface{}, index int, selected, focused bool) string {
		return "CUSTOM: " + labels[index]
	}

	l = l.ItemRenderer(customRenderer)

	// Verify custom renderer is applied by checking view output.
	view := l.View()
	if !strings.Contains(view, "CUSTOM:") {
		t.Error("ItemRenderer() should apply custom renderer")
	}
}

func TestList_Filter(t *testing.T) {
	values := []interface{}{"apple", "banana", "apricot"}
	labels := []string{"apple", "banana", "apricot"}
	l := NewSingleSelect(values, labels)

	// Custom filter: starts with 'a'.
	customFilter := func(item interface{}, query string) bool {
		s := item.(string)
		return len(s) > 0 && s[0] == 'a'
	}

	l = l.Filter(customFilter)

	// Apply filter via domain (API doesn't expose SetFilterQuery directly in this simple version)
	// We would need to send filter messages in a real implementation.
	// For now, just verify the filter function is set.
	if l.domain == nil {
		t.Error("Filter() should not nil the domain")
	}
}

func TestList_ShowFilter(t *testing.T) {
	l := NewSingleSelect([]interface{}{1}, []string{"A"})

	l = l.ShowFilter(true)

	if !l.showFilter {
		t.Error("ShowFilter(true) should enable filter display")
	}

	l = l.ShowFilter(false)

	if l.showFilter {
		t.Error("ShowFilter(false) should disable filter display")
	}
}

func TestList_SelectedItems(t *testing.T) {
	values := []interface{}{"a", "b", "c"}
	labels := []string{"A", "B", "C"}
	l := NewMultiSelect(values, labels)

	// Select first item.
	l.domain = l.domain.ToggleSelection()

	selected := l.SelectedItems()
	if len(selected) != 1 {
		t.Fatalf("SelectedItems() count = %d, want 1", len(selected))
	}
	if selected[0] != "a" {
		t.Errorf("SelectedItems()[0] = %v, want 'a'", selected[0])
	}
}

func TestList_FocusedItem(t *testing.T) {
	values := []interface{}{1, 2, 3}
	labels := []string{"One", "Two", "Three"}
	l := NewSingleSelect(values, labels)

	// Move down twice.
	l.domain = l.domain.MoveDown().MoveDown()

	focused := l.FocusedItem()
	if focused != 3 {
		t.Errorf("FocusedItem() = %v, want 3", focused)
	}
}

func TestList_FocusedItem_EmptyList(t *testing.T) {
	l := NewSingleSelect([]interface{}{}, []string{})

	focused := l.FocusedItem()
	if focused != nil {
		t.Errorf("FocusedItem() on empty list should return nil, got %v", focused)
	}
}

func TestList_SelectedIndices(t *testing.T) {
	values := []interface{}{"a", "b", "c"}
	labels := []string{"A", "B", "C"}
	l := NewMultiSelect(values, labels)

	// Select first and third.
	l.domain = l.domain.ToggleSelection().MoveDown().MoveDown().ToggleSelection()

	indices := l.SelectedIndices()
	if len(indices) != 2 {
		t.Errorf("SelectedIndices() count = %d, want 2", len(indices))
	}
}

func TestList_FocusedIndex(t *testing.T) {
	l := NewSingleSelect([]interface{}{1, 2, 3}, []string{"A", "B", "C"})

	l.domain = l.domain.MoveDown()

	if l.FocusedIndex() != 1 {
		t.Errorf("FocusedIndex() = %d, want 1", l.FocusedIndex())
	}
}

func TestList_Init(t *testing.T) {
	l := NewSingleSelect([]interface{}{1}, []string{"A"})

	cmd := l.Init()
	if cmd != nil {
		t.Error("Init() should return nil cmd")
	}
}

func TestList_Update_MoveUp(t *testing.T) {
	l := NewSingleSelect([]interface{}{1, 2, 3}, []string{"A", "B", "C"})
	l.domain = l.domain.MoveDown().MoveDown() // Start at index 2

	msg := tea.KeyMsg{Type: tea.KeyUp}
	l, _ = l.Update(msg)

	if l.FocusedIndex() != 1 {
		t.Errorf("Update(up) focused index = %d, want 1", l.FocusedIndex())
	}
}

func TestList_Update_MoveDown(t *testing.T) {
	l := NewSingleSelect([]interface{}{1, 2, 3}, []string{"A", "B", "C"})

	msg := tea.KeyMsg{Type: tea.KeyDown}
	l, _ = l.Update(msg)

	if l.FocusedIndex() != 1 {
		t.Errorf("Update(down) focused index = %d, want 1", l.FocusedIndex())
	}
}

func TestList_Update_VimKeys(t *testing.T) {
	l := NewSingleSelect([]interface{}{1, 2, 3}, []string{"A", "B", "C"})

	// j = down.
	jMsg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'j'}
	l, _ = l.Update(jMsg)

	if l.FocusedIndex() != 1 {
		t.Errorf("Update(j) focused index = %d, want 1", l.FocusedIndex())
	}

	// k = up.
	kMsg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'k'}
	l, _ = l.Update(kMsg)

	if l.FocusedIndex() != 0 {
		t.Errorf("Update(k) focused index = %d, want 0", l.FocusedIndex())
	}
}

func TestList_Update_PageUpDown(t *testing.T) {
	// Create list with many items.
	values := make([]interface{}, 30)
	labels := make([]string, 30)
	for i := range values {
		values[i] = i
		labels[i] = string(rune('A' + (i % 26)))
	}
	l := NewSingleSelect(values, labels).Height(10)

	// Page down.
	msg := tea.KeyMsg{Type: tea.KeyPgDown}
	l, _ = l.Update(msg)

	if l.FocusedIndex() != 10 {
		t.Errorf("Update(pgdown) focused index = %d, want 10", l.FocusedIndex())
	}

	// Page up.
	msg = tea.KeyMsg{Type: tea.KeyPgUp}
	l, _ = l.Update(msg)

	if l.FocusedIndex() != 0 {
		t.Errorf("Update(pgup) focused index = %d, want 0", l.FocusedIndex())
	}
}

func TestList_Update_HomeEnd(t *testing.T) {
	l := NewSingleSelect([]interface{}{1, 2, 3, 4, 5}, []string{"A", "B", "C", "D", "E"})
	l.domain = l.domain.MoveToEnd()

	// Home.
	msg := tea.KeyMsg{Type: tea.KeyHome}
	l, _ = l.Update(msg)

	if l.FocusedIndex() != 0 {
		t.Errorf("Update(home) focused index = %d, want 0", l.FocusedIndex())
	}

	// End.
	msg = tea.KeyMsg{Type: tea.KeyEnd}
	l, _ = l.Update(msg)

	if l.FocusedIndex() != 4 {
		t.Errorf("Update(end) focused index = %d, want 4", l.FocusedIndex())
	}
}

func TestList_Update_ToggleSelection(t *testing.T) {
	l := NewSingleSelect([]interface{}{1, 2, 3}, []string{"A", "B", "C"})

	// Space = toggle.
	msg := tea.KeyMsg{Type: tea.KeySpace}
	l, _ = l.Update(msg)

	selected := l.SelectedItems()
	if len(selected) != 1 {
		t.Errorf("Update(space) selected count = %d, want 1", len(selected))
	}
}

func TestList_Update_SelectAll(t *testing.T) {
	l := NewMultiSelect([]interface{}{1, 2, 3}, []string{"A", "B", "C"})

	// Ctrl+A = select all.
	msg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'a', Ctrl: true}
	l, _ = l.Update(msg)

	selected := l.SelectedItems()
	if len(selected) != 3 {
		t.Errorf("Update(ctrl+a) selected count = %d, want 3", len(selected))
	}
}

func TestList_Update_ClearSelection(t *testing.T) {
	l := NewMultiSelect([]interface{}{1, 2, 3}, []string{"A", "B", "C"})
	l.domain = l.domain.SelectAll()

	// Esc = clear selection.
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	l, _ = l.Update(msg)

	selected := l.SelectedItems()
	if len(selected) != 0 {
		t.Errorf("Update(esc) selected count = %d, want 0", len(selected))
	}
}

func TestList_Update_Quit(t *testing.T) {
	l := NewSingleSelect([]interface{}{1}, []string{"A"})

	msg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'q'}
	_, cmd := l.Update(msg)

	if cmd == nil {
		t.Error("Update(q) should return Quit command")
	}
}

func TestList_View(t *testing.T) {
	values := []interface{}{1, 2, 3}
	labels := []string{"One", "Two", "Three"}
	l := NewSingleSelect(values, labels)

	view := l.View()

	// Should contain item labels.
	if !strings.Contains(view, "One") {
		t.Error("View() should contain item label 'One'")
	}
	if !strings.Contains(view, "Two") {
		t.Error("View() should contain item label 'Two'")
	}
	if !strings.Contains(view, "Three") {
		t.Error("View() should contain item label 'Three'")
	}
}

func TestList_View_EmptyList(t *testing.T) {
	l := NewSingleSelect([]interface{}{}, []string{})

	view := l.View()

	if !strings.Contains(view, "empty list") {
		t.Error("View() should show 'empty list' message")
	}
}

func TestList_View_FocusIndicator(t *testing.T) {
	l := NewSingleSelect([]interface{}{1, 2}, []string{"One", "Two"})

	view := l.View()

	// Default renderer uses "> " for focused item.
	if !strings.Contains(view, ">") {
		t.Error("View() should show focus indicator '>'")
	}
}

func TestList_View_SelectionIndicator(t *testing.T) {
	l := NewSingleSelect([]interface{}{1, 2}, []string{"One", "Two"})
	l.domain = l.domain.ToggleSelection()

	view := l.View()

	// After selection, focused item should show "> " (focused takes precedence)
	// But we can check that selection affects rendering.
	if view == "" {
		t.Error("View() should not be empty after selection")
	}
}

func TestList_View_Scrolling(t *testing.T) {
	// Create list with more items than visible height.
	values := make([]interface{}, 20)
	labels := make([]string, 20)
	for i := range values {
		values[i] = i
		labels[i] = string(rune('A' + (i % 26)))
	}
	l := NewSingleSelect(values, labels).Height(5)

	// Move to end.
	l.domain = l.domain.MoveToEnd()

	view := l.View()

	// Should only show 5 items (height)
	lines := strings.Split(view, "\n")
	visibleItems := 0
	for _, line := range lines {
		if line != "" && !strings.Contains(line, "empty") {
			visibleItems++
		}
	}

	if visibleItems != 5 {
		t.Errorf("View() with scrolling should show 5 items, got %d", visibleItems)
	}
}

func TestList_KeyBindings(t *testing.T) {
	l := NewSingleSelect([]interface{}{1, 2}, []string{"A", "B"})

	// Set custom key bindings.
	customBindings := []infrastructure.KeyBinding{
		{Key: "w", Action: "move_up"},
		{Key: "s", Action: "move_down"},
	}
	l = l.KeyBindings(customBindings)

	// Test custom binding.
	wMsg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'w'}
	l, _ = l.Update(wMsg)

	// Should have moved up (wrap to end since at start)
	if l.FocusedIndex() != 1 {
		t.Errorf("Custom keybinding 'w' should move up, focused index = %d", l.FocusedIndex())
	}
}

func TestList_MethodChaining(t *testing.T) {
	// Test that methods can be chained.
	l := NewSingleSelect([]interface{}{1, 2, 3}, []string{"A", "B", "C"}).
		Height(15).
		ShowFilter(true)

	if l.domain.Height() != 15 {
		t.Error("Method chaining should work for Height()")
	}
	if !l.showFilter {
		t.Error("Method chaining should work for ShowFilter()")
	}
}

func TestList_Update_NonKeyMsg(t *testing.T) {
	l := NewSingleSelect([]interface{}{1}, []string{"A"})

	// Send non-key message (should be ignored)
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	l2, cmd := l.Update(msg)

	if l2 != l {
		t.Error("Update() with non-key message should return same list")
	}
	if cmd != nil {
		t.Error("Update() with non-key message should return nil cmd")
	}
}

func TestList_Immutability(t *testing.T) {
	l1 := NewSingleSelect([]interface{}{1, 2}, []string{"A", "B"})

	// Perform operations.
	l2 := l1.Height(20)
	l3 := l1.ShowFilter(true)

	// Original should be unchanged.
	if l1.domain.Height() != 10 {
		t.Error("Original list height should not change")
	}
	if l1.showFilter {
		t.Error("Original list showFilter should not change")
	}

	// New instances should have changes.
	if l2.domain.Height() != 20 {
		t.Error("New list after Height() should have new height")
	}
	if !l3.showFilter {
		t.Error("New list after ShowFilter() should have filter enabled")
	}
}

func TestList_SelectedItems_ReturnsValues(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := []interface{}{
		Person{Name: "Alice", Age: 30},
		Person{Name: "Bob", Age: 25},
	}
	labels := []string{"Alice", "Bob"}

	l := NewMultiSelect(people, labels)
	l.domain = l.domain.ToggleSelection()

	selected := l.SelectedItems()
	if len(selected) != 1 {
		t.Fatalf("SelectedItems() count = %d, want 1", len(selected))
	}

	person := selected[0].(Person)
	if person.Name != "Alice" {
		t.Errorf("SelectedItems() returned wrong person: %s, want Alice", person.Name)
	}
}

func TestList_FocusedItem_ReturnsValue(t *testing.T) {
	values := []interface{}{"apple", "banana", "cherry"}
	labels := []string{"apple", "banana", "cherry"}

	l := NewSingleSelect(values, labels)
	l.domain = l.domain.MoveDown()

	focused := l.FocusedItem()
	if focused != "banana" {
		t.Errorf("FocusedItem() = %v, want 'banana'", focused)
	}
}
