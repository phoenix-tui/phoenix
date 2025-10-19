package model

import (
	"reflect"
	"testing"

	"github.com/phoenix-tui/phoenix/components/list/domain/value"
)

func createTestItems(count int) []*value.Item {
	items := make([]*value.Item, count)
	for i := 0; i < count; i++ {
		items[i] = value.NewItem(i, "Item "+string(rune('A'+i)))
	}
	return items
}

func TestNewList(t *testing.T) {
	tests := []struct {
		name string
		mode value.SelectionMode
	}{
		{"single selection", value.SelectionModeSingle},
		{"multi selection", value.SelectionModeMulti},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewList(tt.mode)

			if l.SelectionMode() != tt.mode {
				t.Errorf("NewList() mode = %v, want %v", l.SelectionMode(), tt.mode)
			}
			if len(l.Items()) != 0 {
				t.Error("NewList() should start with empty items")
			}
			if l.Height() != 10 {
				t.Errorf("NewList() height = %d, want 10", l.Height())
			}
		})
	}
}

func TestNewListWithItems(t *testing.T) {
	items := createTestItems(5)
	l := NewListWithItems(items, value.SelectionModeSingle)

	if len(l.Items()) != 5 {
		t.Errorf("NewListWithItems() item count = %d, want 5", len(l.Items()))
	}
	if len(l.FilteredItems()) != 5 {
		t.Errorf("NewListWithItems() filtered count = %d, want 5", len(l.FilteredItems()))
	}
}

func TestList_WithItems(t *testing.T) {
	l := NewList(value.SelectionModeSingle)
	items := createTestItems(3)

	l2 := l.WithItems(items)

	// Original unchanged
	if len(l.Items()) != 0 {
		t.Error("WithItems() should not modify original list")
	}

	// New list has items
	if len(l2.Items()) != 3 {
		t.Errorf("WithItems() item count = %d, want 3", len(l2.Items()))
	}
}

func TestList_WithHeight(t *testing.T) {
	l := NewList(value.SelectionModeSingle)
	l2 := l.WithHeight(20)

	if l.Height() != 10 {
		t.Error("WithHeight() should not modify original list")
	}
	if l2.Height() != 20 {
		t.Errorf("WithHeight() height = %d, want 20", l2.Height())
	}
}

func TestList_WithItemRenderer(t *testing.T) {
	l := NewList(value.SelectionModeSingle)

	customRenderer := func(item *value.Item, index int, selected, focused bool) string {
		return "CUSTOM: " + item.Label()
	}

	l2 := l.WithItemRenderer(customRenderer)

	// Renderers should be different (can't directly compare functions, but we can test behavior)
	items := []*value.Item{value.NewItem(1, "Test")}
	l3 := l2.WithItems(items)

	rendered := l3.RenderItem(0)
	if rendered != "CUSTOM: Test" {
		t.Errorf("WithItemRenderer() rendered = %q, want %q", rendered, "CUSTOM: Test")
	}
}

func TestList_WithFilter(t *testing.T) {
	items := []*value.Item{
		value.NewItem(1, "apple"),
		value.NewItem(2, "banana"),
		value.NewItem(3, "apricot"),
	}

	l := NewListWithItems(items, value.SelectionModeSingle)

	// Custom filter: starts with 'a'
	customFilter := func(item *value.Item, query string) bool {
		return len(item.Label()) > 0 && item.Label()[0] == 'a'
	}

	l2 := l.WithFilter(customFilter).SetFilterQuery("a")

	if len(l2.FilteredItems()) != 2 {
		t.Errorf("WithFilter() filtered count = %d, want 2", len(l2.FilteredItems()))
	}
}

func TestList_MoveUp(t *testing.T) {
	items := createTestItems(5)
	l := NewListWithItems(items, value.SelectionModeSingle)

	// Move down twice, then up once
	l = l.MoveDown().MoveDown().MoveUp()

	if l.FocusedIndex() != 1 {
		t.Errorf("MoveUp() focused index = %d, want 1", l.FocusedIndex())
	}
}

func TestList_MoveUp_WrapAround(t *testing.T) {
	items := createTestItems(5)
	l := NewListWithItems(items, value.SelectionModeSingle)

	// At index 0, moving up should wrap to last
	l = l.MoveUp()

	if l.FocusedIndex() != 4 {
		t.Errorf("MoveUp() wrap-around focused index = %d, want 4", l.FocusedIndex())
	}
}

func TestList_MoveDown(t *testing.T) {
	items := createTestItems(5)
	l := NewListWithItems(items, value.SelectionModeSingle)

	l = l.MoveDown().MoveDown()

	if l.FocusedIndex() != 2 {
		t.Errorf("MoveDown() focused index = %d, want 2", l.FocusedIndex())
	}
}

func TestList_MoveDown_WrapAround(t *testing.T) {
	items := createTestItems(3)
	l := NewListWithItems(items, value.SelectionModeSingle)

	// Move to last item and then down should wrap to first
	l = l.MoveToEnd().MoveDown()

	if l.FocusedIndex() != 0 {
		t.Errorf("MoveDown() wrap-around focused index = %d, want 0", l.FocusedIndex())
	}
}

func TestList_MovePageUp(t *testing.T) {
	items := createTestItems(30)
	l := NewListWithItems(items, value.SelectionModeSingle).WithHeight(10)

	// Move to middle
	l = l.MoveToEnd().MovePageUp()

	// Should be at index 19 (29 - 10)
	if l.FocusedIndex() != 19 {
		t.Errorf("MovePageUp() focused index = %d, want 19", l.FocusedIndex())
	}
}

func TestList_MovePageDown(t *testing.T) {
	items := createTestItems(30)
	l := NewListWithItems(items, value.SelectionModeSingle).WithHeight(10)

	l = l.MovePageDown()

	// Should be at index 10
	if l.FocusedIndex() != 10 {
		t.Errorf("MovePageDown() focused index = %d, want 10", l.FocusedIndex())
	}
}

func TestList_MoveToStart(t *testing.T) {
	items := createTestItems(5)
	l := NewListWithItems(items, value.SelectionModeSingle)

	l = l.MoveToEnd().MoveToStart()

	if l.FocusedIndex() != 0 {
		t.Errorf("MoveToStart() focused index = %d, want 0", l.FocusedIndex())
	}
}

func TestList_MoveToEnd(t *testing.T) {
	items := createTestItems(5)
	l := NewListWithItems(items, value.SelectionModeSingle)

	l = l.MoveToEnd()

	if l.FocusedIndex() != 4 {
		t.Errorf("MoveToEnd() focused index = %d, want 4", l.FocusedIndex())
	}
}

func TestList_ToggleSelection_Single(t *testing.T) {
	items := createTestItems(5)
	l := NewListWithItems(items, value.SelectionModeSingle)

	// Select first item
	l = l.ToggleSelection()

	selected := l.SelectedItems()
	if len(selected) != 1 {
		t.Fatalf("ToggleSelection() single mode selected count = %d, want 1", len(selected))
	}
	if selected[0].Label() != "Item A" {
		t.Errorf("ToggleSelection() selected item = %s, want Item A", selected[0].Label())
	}

	// Move to second item and select (should clear first)
	l = l.MoveDown().ToggleSelection()

	selected = l.SelectedItems()
	if len(selected) != 1 {
		t.Fatalf("ToggleSelection() single mode should only select one = %d, want 1", len(selected))
	}
	if selected[0].Label() != "Item B" {
		t.Errorf("ToggleSelection() selected item = %s, want Item B", selected[0].Label())
	}
}

func TestList_ToggleSelection_Multi(t *testing.T) {
	items := createTestItems(5)
	l := NewListWithItems(items, value.SelectionModeMulti)

	// Select first item
	l = l.ToggleSelection()

	// Select third item
	l = l.MoveDown().MoveDown().ToggleSelection()

	selected := l.SelectedItems()
	if len(selected) != 2 {
		t.Fatalf("ToggleSelection() multi mode selected count = %d, want 2", len(selected))
	}

	// Toggle first item again (should deselect)
	l = l.MoveToStart().ToggleSelection()

	selected = l.SelectedItems()
	if len(selected) != 1 {
		t.Errorf("ToggleSelection() after deselect count = %d, want 1", len(selected))
	}
}

func TestList_SelectAll(t *testing.T) {
	items := createTestItems(5)
	l := NewListWithItems(items, value.SelectionModeMulti)

	l = l.SelectAll()

	selected := l.SelectedItems()
	if len(selected) != 5 {
		t.Errorf("SelectAll() selected count = %d, want 5", len(selected))
	}
}

func TestList_SelectAll_SingleMode(t *testing.T) {
	items := createTestItems(5)
	l := NewListWithItems(items, value.SelectionModeSingle)

	l = l.SelectAll()

	// Should not select all in single mode
	selected := l.SelectedItems()
	if len(selected) != 0 {
		t.Errorf("SelectAll() in single mode should not select, got %d items", len(selected))
	}
}

func TestList_ClearSelection(t *testing.T) {
	items := createTestItems(5)
	l := NewListWithItems(items, value.SelectionModeMulti)

	l = l.SelectAll().ClearSelection()

	selected := l.SelectedItems()
	if len(selected) != 0 {
		t.Errorf("ClearSelection() selected count = %d, want 0", len(selected))
	}
}

func TestList_SetFilterQuery(t *testing.T) {
	items := []*value.Item{
		value.NewItem(1, "apple"),
		value.NewItem(2, "banana"),
		value.NewItem(3, "apricot"),
	}

	l := NewListWithItems(items, value.SelectionModeSingle)
	l = l.SetFilterQuery("ap")

	filtered := l.FilteredItems()
	if len(filtered) != 2 {
		t.Errorf("SetFilterQuery() filtered count = %d, want 2", len(filtered))
	}

	// Check filtered items
	expectedLabels := []string{"apple", "apricot"}
	for i, item := range filtered {
		if item.Label() != expectedLabels[i] {
			t.Errorf("SetFilterQuery() item %d = %s, want %s", i, item.Label(), expectedLabels[i])
		}
	}
}

func TestList_ClearFilter(t *testing.T) {
	items := createTestItems(5)
	l := NewListWithItems(items, value.SelectionModeSingle)

	l = l.SetFilterQuery("Item A").ClearFilter()

	if l.IsFiltered() {
		t.Error("ClearFilter() should clear filter state")
	}
	if len(l.FilteredItems()) != 5 {
		t.Errorf("ClearFilter() filtered count = %d, want 5", len(l.FilteredItems()))
	}
}

func TestList_FocusedItem(t *testing.T) {
	items := createTestItems(5)
	l := NewListWithItems(items, value.SelectionModeSingle)

	l = l.MoveDown().MoveDown()

	focused := l.FocusedItem()
	if focused == nil {
		t.Fatal("FocusedItem() should not be nil")
	}
	if focused.Label() != "Item C" {
		t.Errorf("FocusedItem() label = %s, want Item C", focused.Label())
	}
}

func TestList_FocusedItem_EmptyList(t *testing.T) {
	l := NewList(value.SelectionModeSingle)

	focused := l.FocusedItem()
	if focused != nil {
		t.Error("FocusedItem() should be nil for empty list")
	}
}

func TestList_RenderItem(t *testing.T) {
	items := []*value.Item{value.NewItem(1, "Test Item")}
	l := NewListWithItems(items, value.SelectionModeSingle)

	rendered := l.RenderItem(0)

	// Default renderer should include "> " for focused item
	if rendered != "> Test Item" {
		t.Errorf("RenderItem() = %q, want %q", rendered, "> Test Item")
	}
}

func TestList_RenderVisibleItems(t *testing.T) {
	items := createTestItems(20)
	l := NewListWithItems(items, value.SelectionModeSingle).WithHeight(5)

	// At start
	rendered := l.RenderVisibleItems()
	if len(rendered) != 5 {
		t.Errorf("RenderVisibleItems() count = %d, want 5", len(rendered))
	}

	// Move down to trigger scrolling
	l = l.MoveToEnd()
	rendered = l.RenderVisibleItems()
	if len(rendered) != 5 {
		t.Errorf("RenderVisibleItems() after scroll count = %d, want 5", len(rendered))
	}
}

func TestList_RenderVisibleItems_EmptyList(t *testing.T) {
	l := NewList(value.SelectionModeSingle)

	rendered := l.RenderVisibleItems()
	if len(rendered) != 0 {
		t.Errorf("RenderVisibleItems() empty list should return empty slice, got %d items", len(rendered))
	}
}

func TestList_ScrollOffset(t *testing.T) {
	items := createTestItems(30)
	l := NewListWithItems(items, value.SelectionModeSingle).WithHeight(10)

	// Move to middle
	l = l.MovePageDown().MovePageDown()

	offset := l.ScrollOffset()
	if offset == 0 {
		t.Error("ScrollOffset() should be non-zero after scrolling")
	}
}

func TestList_Immutability(t *testing.T) {
	items := createTestItems(5)
	l1 := NewListWithItems(items, value.SelectionModeSingle)

	// Perform various operations
	l2 := l1.MoveDown()
	l3 := l1.ToggleSelection()
	l4 := l1.WithHeight(20)

	// Original should be unchanged
	if l1.FocusedIndex() != 0 {
		t.Error("Original list focus should not change")
	}
	if len(l1.SelectedItems()) != 0 {
		t.Error("Original list selection should not change")
	}
	if l1.Height() != 10 {
		t.Error("Original list height should not change")
	}

	// New instances should have changes
	if l2.FocusedIndex() != 1 {
		t.Error("New list after MoveDown should have changed focus")
	}
	if len(l3.SelectedItems()) != 1 {
		t.Error("New list after ToggleSelection should have selection")
	}
	if l4.Height() != 20 {
		t.Error("New list after WithHeight should have new height")
	}
}

func TestList_SelectedIndices(t *testing.T) {
	items := createTestItems(5)
	l := NewListWithItems(items, value.SelectionModeMulti)

	// Select first and third items
	l = l.ToggleSelection().MoveDown().MoveDown().ToggleSelection()

	indices := l.SelectedIndices()
	if len(indices) != 2 {
		t.Errorf("SelectedIndices() count = %d, want 2", len(indices))
	}

	// Check that indices contain 0 and 2
	hasZero := false
	hasTwo := false
	for _, idx := range indices {
		if idx == 0 {
			hasZero = true
		}
		if idx == 2 {
			hasTwo = true
		}
	}
	if !hasZero || !hasTwo {
		t.Errorf("SelectedIndices() = %v, want indices 0 and 2", indices)
	}
}

func TestList_FilterResetsSelection(t *testing.T) {
	items := createTestItems(5)
	l := NewListWithItems(items, value.SelectionModeMulti)

	// Select all items
	l = l.SelectAll()
	if len(l.SelectedItems()) != 5 {
		t.Fatal("Setup: SelectAll should select 5 items")
	}

	// Apply filter
	l = l.SetFilterQuery("Item A")

	// Selection should be cleared after filter
	selected := l.SelectedItems()
	if len(selected) != 0 {
		t.Errorf("Filter should clear selection, got %d selected items", len(selected))
	}
}

func TestList_Navigation_EmptyList(t *testing.T) {
	l := NewList(value.SelectionModeSingle)

	// All navigation operations should not panic on empty list
	l = l.MoveUp()
	l = l.MoveDown()
	l = l.MovePageUp()
	l = l.MovePageDown()
	l = l.MoveToStart()
	l = l.MoveToEnd()

	// Focus should remain 0
	if l.FocusedIndex() != 0 {
		t.Errorf("Navigation on empty list should keep focus at 0, got %d", l.FocusedIndex())
	}
}

func TestList_Selection_EmptyList(t *testing.T) {
	l := NewList(value.SelectionModeSingle)

	// Selection operations should not panic on empty list
	l = l.ToggleSelection()
	l = l.SelectAll()
	l = l.ClearSelection()

	if len(l.SelectedItems()) != 0 {
		t.Error("Selection on empty list should not select any items")
	}
}

func TestDefaultItemRenderer(t *testing.T) {
	item := value.NewItem(1, "Test")

	tests := []struct {
		name     string
		selected bool
		focused  bool
		want     string
	}{
		{"neither", false, false, "  Test"},
		{"selected", true, false, "✓ Test"},
		{"focused", false, true, "> Test"},
		{"both (focused takes precedence)", true, true, "> Test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := defaultItemRenderer(item, 0, tt.selected, tt.focused)
			if got != tt.want {
				t.Errorf("defaultItemRenderer() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestList_Items_ReturnsCopy(t *testing.T) {
	items := createTestItems(3)
	l := NewListWithItems(items, value.SelectionModeSingle)

	retrieved := l.Items()
	retrieved[0] = value.NewItem(999, "Modified")

	// Original should be unchanged
	original := l.Items()
	if !reflect.DeepEqual(original[0], items[0]) {
		t.Error("Items() should return a copy, not allow modification of original")
	}
}

func TestList_FilteredItems_ReturnsCopy(t *testing.T) {
	items := createTestItems(3)
	l := NewListWithItems(items, value.SelectionModeSingle)

	retrieved := l.FilteredItems()
	retrieved[0] = value.NewItem(999, "Modified")

	// Original should be unchanged
	original := l.FilteredItems()
	if !reflect.DeepEqual(original[0], items[0]) {
		t.Error("FilteredItems() should return a copy, not allow modification of original")
	}
}
