package model

import (
	"reflect"
	"testing"

	"github.com/phoenix-tui/phoenix/components/multiselect/internal/domain/value"
)

func makeOptions() []*value.Option[string] {
	return []*value.Option[string]{
		value.NewOption("Option 1", "opt1"),
		value.NewOption("Option 2", "opt2"),
		value.NewOption("Option 3", "opt3"),
		value.NewOption("Option 4", "opt4"),
		value.NewOption("Option 5", "opt5"),
	}
}

func TestNew(t *testing.T) {
	opts := makeOptions()
	m := New(opts, 1, 3)

	if m.TotalCount() != 5 {
		t.Errorf("TotalCount() = %d, want 5", m.TotalCount())
	}
	if m.SelectionCount() != 0 {
		t.Errorf("SelectionCount() = %d, want 0", m.SelectionCount())
	}
	if m.FilteredCount() != 5 {
		t.Errorf("FilteredCount() = %d, want 5", m.FilteredCount())
	}
	if m.CanConfirm() {
		// Should not be able to confirm with min=1 and 0 selected
		t.Error("CanConfirm() = true, want false (min=1, selected=0)")
	}
}

func TestNew_EmptyOptions(t *testing.T) {
	m := New([]*value.Option[string]{}, 0, 0)

	if m.TotalCount() != 0 {
		t.Errorf("TotalCount() = %d, want 0", m.TotalCount())
	}
}

func TestMultiSelect_WithHeight(t *testing.T) {
	m := New(makeOptions(), 0, 0).WithHeight(3)

	// Height affects visible rendering (tested in RenderVisibleOptions)
	if m.height != 3 {
		t.Errorf("height = %d, want 3", m.height)
	}

	// Negative/zero height should be clamped to 1
	m = m.WithHeight(0)
	if m.height != 1 {
		t.Errorf("height = %d, want 1 (clamped)", m.height)
	}
}

func TestMultiSelect_WithSelected(t *testing.T) {
	m := New(makeOptions(), 0, 0).WithSelected(0, 2, 4)

	indices := m.SelectedIndices()
	want := []int{0, 2, 4}
	if !reflect.DeepEqual(indices, want) {
		t.Errorf("SelectedIndices() = %v, want %v", indices, want)
	}

	values := m.SelectedValues()
	wantVals := []string{"opt1", "opt3", "opt5"}
	if !reflect.DeepEqual(values, wantVals) {
		t.Errorf("SelectedValues() = %v, want %v", values, wantVals)
	}
}

func TestMultiSelect_MoveUp(t *testing.T) {
	m := New(makeOptions(), 0, 0)

	// Start at 0
	if m.cursor != 0 {
		t.Errorf("initial cursor = %d, want 0", m.cursor)
	}

	// Move up (should stay at 0)
	m = m.MoveUp()
	if m.cursor != 0 {
		t.Errorf("cursor after up = %d, want 0", m.cursor)
	}

	// Move down to 2, then up to 1
	m = m.MoveDown().MoveDown().MoveUp()
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1", m.cursor)
	}
}

func TestMultiSelect_MoveDown(t *testing.T) {
	m := New(makeOptions(), 0, 0)

	m = m.MoveDown()
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1", m.cursor)
	}

	m = m.MoveDown().MoveDown()
	if m.cursor != 3 {
		t.Errorf("cursor = %d, want 3", m.cursor)
	}

	// Move to end (maxIndex = 4)
	m = m.MoveDown().MoveDown()
	if m.cursor != 4 {
		t.Errorf("cursor = %d, want 4", m.cursor)
	}

	// Try to move past end (should stay at 4)
	m = m.MoveDown()
	if m.cursor != 4 {
		t.Errorf("cursor after exceeding = %d, want 4", m.cursor)
	}
}

func TestMultiSelect_MoveToStart(t *testing.T) {
	m := New(makeOptions(), 0, 0).MoveDown().MoveDown().MoveDown()
	if m.cursor != 3 {
		t.Errorf("setup: cursor = %d, want 3", m.cursor)
	}

	m = m.MoveToStart()
	if m.cursor != 0 {
		t.Errorf("cursor after MoveToStart = %d, want 0", m.cursor)
	}
}

func TestMultiSelect_MoveToEnd(t *testing.T) {
	m := New(makeOptions(), 0, 0).MoveToEnd()

	if m.cursor != 4 {
		t.Errorf("cursor = %d, want 4 (maxIndex)", m.cursor)
	}
}

func TestMultiSelect_Toggle(t *testing.T) {
	m := New(makeOptions(), 0, 0)

	// Toggle on (cursor at 0)
	m = m.Toggle()
	if m.SelectionCount() != 1 {
		t.Errorf("SelectionCount() = %d, want 1", m.SelectionCount())
	}
	if !m.selection.IsSelected(0) {
		t.Error("IsSelected(0) = false, want true")
	}

	// Toggle off
	m = m.Toggle()
	if m.SelectionCount() != 0 {
		t.Errorf("SelectionCount() = %d, want 0", m.SelectionCount())
	}

	// Toggle multiple items
	m = m.Toggle()            // index 0
	m = m.MoveDown().Toggle() // index 1
	m = m.MoveDown().Toggle() // index 2
	if m.SelectionCount() != 3 {
		t.Errorf("SelectionCount() = %d, want 3", m.SelectionCount())
	}
}

func TestMultiSelect_SelectAll(t *testing.T) {
	// No constraints
	m := New(makeOptions(), 0, 0).SelectAll()
	if m.SelectionCount() != 5 {
		t.Errorf("SelectionCount() = %d, want 5", m.SelectionCount())
	}

	// With max constraint
	m = New(makeOptions(), 0, 3).SelectAll()
	if m.SelectionCount() != 3 {
		t.Errorf("SelectionCount() = %d, want 3 (respects max)", m.SelectionCount())
	}
}

func TestMultiSelect_SelectNone(t *testing.T) {
	m := New(makeOptions(), 0, 0).
		WithSelected(0, 1, 2, 3, 4).
		SelectNone()

	if m.SelectionCount() != 0 {
		t.Errorf("SelectionCount() = %d, want 0", m.SelectionCount())
	}

	indices := m.SelectedIndices()
	if len(indices) != 0 {
		t.Errorf("SelectedIndices() = %v, want []", indices)
	}
}

func TestMultiSelect_SetFilterQuery(t *testing.T) {
	m := New(makeOptions(), 0, 0)

	// Filter to "1" (should match "Option 1")
	m = m.SetFilterQuery("1")
	if m.FilteredCount() != 1 {
		t.Errorf("FilteredCount() = %d, want 1", m.FilteredCount())
	}
	if !m.IsFiltered() {
		t.Error("IsFiltered() = false, want true")
	}

	// Filter to "option" (should match all)
	m = m.SetFilterQuery("option")
	if m.FilteredCount() != 5 {
		t.Errorf("FilteredCount() = %d, want 5", m.FilteredCount())
	}

	// Filter to non-matching
	m = m.SetFilterQuery("xyz")
	if m.FilteredCount() != 0 {
		t.Errorf("FilteredCount() = %d, want 0", m.FilteredCount())
	}

	// Cursor should be clamped
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0 (clamped to empty)", m.cursor)
	}
}

func TestMultiSelect_ClearFilter(t *testing.T) {
	m := New(makeOptions(), 0, 0).
		SetFilterQuery("1").
		ClearFilter()

	if m.IsFiltered() {
		t.Error("IsFiltered() = true, want false")
	}
	if m.FilteredCount() != 5 {
		t.Errorf("FilteredCount() = %d, want 5", m.FilteredCount())
	}
	if m.FilterQuery() != "" {
		t.Errorf("FilterQuery() = %q, want empty", m.FilterQuery())
	}
}

func TestMultiSelect_CanConfirm(t *testing.T) {
	tests := []struct {
		name string
		min  int
		max  int
		sel  []int
		want bool
	}{
		{"no constraints, no selection", 0, 0, []int{}, true},
		{"no constraints, with selection", 0, 0, []int{0, 1}, true},
		{"min=2, selected=2", 2, 0, []int{0, 1}, true},
		{"min=2, selected=1", 2, 0, []int{0}, false},
		{"min=2, selected=3", 2, 0, []int{0, 1, 2}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(makeOptions(), tt.min, tt.max).WithSelected(tt.sel...)
			if got := m.CanConfirm(); got != tt.want {
				t.Errorf("CanConfirm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiSelect_RenderVisibleOptions(t *testing.T) {
	m := New(makeOptions(), 0, 0).WithHeight(3)

	rendered := m.RenderVisibleOptions()
	if len(rendered) != 3 {
		t.Errorf("len(rendered) = %d, want 3 (height limit)", len(rendered))
	}

	// Check focus indicator on first item
	if rendered[0][:1] != ">" {
		t.Errorf("rendered[0] = %q, want to start with '>'", rendered[0])
	}

	// Check checkbox format
	if rendered[0][2:5] != "[ ]" {
		t.Errorf("rendered[0] checkbox = %q, want '[ ]'", rendered[0][2:5])
	}

	// Select item 0 and re-render
	m = m.Toggle()
	rendered = m.RenderVisibleOptions()
	if rendered[0][2:5] != "[x]" {
		t.Errorf("rendered[0] checkbox = %q, want '[x]'", rendered[0][2:5])
	}
}

func TestMultiSelect_RenderVisibleOptions_Empty(t *testing.T) {
	m := New([]*value.Option[string]{}, 0, 0)

	rendered := m.RenderVisibleOptions()
	if len(rendered) != 0 {
		t.Errorf("len(rendered) = %d, want 0", len(rendered))
	}
}

func TestMultiSelect_RenderVisibleOptions_Scrolling(t *testing.T) {
	m := New(makeOptions(), 0, 0).WithHeight(2)

	// Initially shows items 0-1
	rendered := m.RenderVisibleOptions()
	if len(rendered) != 2 {
		t.Errorf("len(rendered) = %d, want 2", len(rendered))
	}

	// Move cursor to item 2 (should scroll)
	m = m.MoveDown().MoveDown()
	rendered = m.RenderVisibleOptions()
	if len(rendered) != 2 {
		t.Errorf("len(rendered) = %d, want 2", len(rendered))
	}

	// Focus should be on second visible item
	if rendered[1][:1] != ">" {
		t.Errorf("rendered[1] should be focused, got: %q", rendered[1])
	}
}

func TestMultiSelect_SelectAll_WithFilter(t *testing.T) {
	m := New(makeOptions(), 0, 0).
		SetFilterQuery("1"). // Filters to "Option 1" only
		SelectAll()

	if m.SelectionCount() != 1 {
		t.Errorf("SelectionCount() = %d, want 1 (only filtered item)", m.SelectionCount())
	}

	// Clear filter - should still have 1 selected
	m = m.ClearFilter()
	if m.SelectionCount() != 1 {
		t.Errorf("SelectionCount() = %d, want 1 (persists after filter clear)", m.SelectionCount())
	}
}

func TestMultiSelect_Immutability(t *testing.T) {
	original := New(makeOptions(), 0, 0)
	modified := original.MoveDown().Toggle().SetFilterQuery("test")

	// Original should be unchanged
	if original.SelectionCount() != 0 {
		t.Errorf("original.SelectionCount() = %d, want 0", original.SelectionCount())
	}
	if original.FilterQuery() != "" {
		t.Errorf("original.FilterQuery() = %q, want empty", original.FilterQuery())
	}

	// Modified should have changes
	if modified.SelectionCount() == 0 {
		t.Error("modified.SelectionCount() should not be 0")
	}
	if modified.FilterQuery() == "" {
		t.Error("modified.FilterQuery() should not be empty")
	}
}
func TestMultiSelect_CustomFilterFunc(t *testing.T) {
	m := New(makeOptions(), 0, 0)

	// Custom filter: only odd numbers
	customFilter := func(opt *value.Option[string], _ string) bool {
		label := opt.Label()
		return label == "Option 1" || label == "Option 3" || label == "Option 5"
	}

	m = m.WithFilterFunc(customFilter).SetFilterQuery("anything")

	if m.FilteredCount() != 3 {
		t.Errorf("FilteredCount() = %d, want 3 (custom filter)", m.FilteredCount())
	}
}

func TestMultiSelect_CustomRenderFunc(t *testing.T) {
	m := New(makeOptions(), 0, 0)

	// Custom render: just show label
	customRender := func(opt *value.Option[string], _ int, _, _ bool) string {
		return opt.Label()
	}

	m = m.WithRenderFunc(customRender)
	rendered := m.RenderVisibleOptions()

	// Should just be the label (no checkbox, no focus indicator)
	if rendered[0] != "Option 1" {
		t.Errorf("rendered[0] = %q, want %q", rendered[0], "Option 1")
	}
}
