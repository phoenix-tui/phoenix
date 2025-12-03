package multiselect

import (
	"reflect"
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/tea"
)

func TestNew(t *testing.T) {
	m := New[string]("Select items:")

	if m.title != "Select items:" {
		t.Errorf("title = %q, want %q", m.title, "Select items:")
	}
	if m.SelectionCount() != 0 {
		t.Errorf("SelectionCount() = %d, want 0", m.SelectionCount())
	}
}

func TestNewStrings(t *testing.T) {
	options := []string{"A", "B", "C"}
	m := NewStrings("Select:", options)

	if m.title != "Select:" {
		t.Errorf("title = %q, want %q", m.title, "Select:")
	}

	// Should have 3 options
	if m.domain.TotalCount() != 3 {
		t.Errorf("TotalCount() = %d, want 3", m.domain.TotalCount())
	}
}

func TestMultiSelect_Options(t *testing.T) {
	m := New[int]("Select numbers:").
		Options(
			Opt("One", 1),
			Opt("Two", 2),
			Opt("Three", 3),
		)

	if m.domain.TotalCount() != 3 {
		t.Errorf("TotalCount() = %d, want 3", m.domain.TotalCount())
	}
}

func TestMultiSelect_Selected(t *testing.T) {
	m := New[string]("Select:").
		Options(
			Opt("A", "a"),
			Opt("B", "b"),
			Opt("C", "c"),
			Opt("D", "d"),
		).
		Selected(0, 2)

	indices := m.SelectedIndices()
	want := []int{0, 2}
	if !reflect.DeepEqual(indices, want) {
		t.Errorf("SelectedIndices() = %v, want %v", indices, want)
	}

	items := m.SelectedItems()
	wantItems := []string{"a", "c"}
	if !reflect.DeepEqual(items, wantItems) {
		t.Errorf("SelectedItems() = %v, want %v", items, wantItems)
	}
}

func TestMultiSelect_MinMax(t *testing.T) {
	m := New[string]("Select:").
		Options(
			Opt("A", "a"),
			Opt("B", "b"),
			Opt("C", "c"),
		).
		Min(1).
		Max(2)

	if m.min != 1 {
		t.Errorf("min = %d, want 1", m.min)
	}
	if m.max != 2 {
		t.Errorf("max = %d, want 2", m.max)
	}

	// Cannot confirm with 0 selected (min=1)
	if m.domain.CanConfirm() {
		t.Error("CanConfirm() = true, want false (min not met)")
	}
}

func TestMultiSelect_WithHeight(t *testing.T) {
	m := New[string]("Select:").
		Options(
			Opt("A", "a"),
			Opt("B", "b"),
			Opt("C", "c"),
		).
		WithHeight(2)

	rendered := m.domain.RenderVisibleOptions()
	if len(rendered) != 2 {
		t.Errorf("len(rendered) = %d, want 2 (height limit)", len(rendered))
	}
}

func TestMultiSelect_WithFilterable(t *testing.T) {
	m := New[string]("Select:").
		Options(Opt("A", "a")).
		WithFilterable(true)

	if !m.filterable {
		t.Error("filterable = false, want true")
	}

	m = m.WithFilterable(false)
	if m.filterable {
		t.Error("filterable = true, want false")
	}
}

// func TestMultiSelect_Update_Navigation(t *testing.T) {
// 	m := New[string]("Select:").
// 		Options(
// 			Opt("A", "a"),
// 			Opt("B", "b"),
// 			Opt("C", "c"),
// 		)
// 
// 	tests := []struct {
// 		name      string
// 		key       tea.KeyMsg
// 		wantIndex int
// 	}{
// 		{"down arrow", tea.KeyMsg{Type: tea.KeyDown}, 1},
// 		{"j", tea.KeyMsg{Type: tea.KeyRune, Rune: 'j'}, 1},
// 		{"up arrow", tea.KeyMsg{Type: tea.KeyUp}, 0},
// 		{"k", tea.KeyMsg{Type: tea.KeyRune, Rune: 'k'}, 0},
// 	}
// 
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			newM, _ := m.Update(tt.key)
// 			}
// 		})
// 	}
// }

func TestMultiSelect_Update_Toggle(t *testing.T) {
	m := New[string]("Select:").
		Options(
			Opt("A", "a"),
			Opt("B", "b"),
		)

	// Toggle first item
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	if m.SelectionCount() != 1 {
		t.Errorf("SelectionCount() = %d, want 1", m.SelectionCount())
	}

	// Toggle again (off)
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	if m.SelectionCount() != 0 {
		t.Errorf("SelectionCount() = %d, want 0", m.SelectionCount())
	}
}

func TestMultiSelect_Update_SelectAll(t *testing.T) {
	m := New[string]("Select:").
		Options(
			Opt("A", "a"),
			Opt("B", "b"),
			Opt("C", "c"),
		)

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'a'})
	if m.SelectionCount() != 3 {
		t.Errorf("SelectionCount() = %d, want 3", m.SelectionCount())
	}
}

func TestMultiSelect_Update_SelectNone(t *testing.T) {
	m := New[string]("Select:").
		Options(
			Opt("A", "a"),
			Opt("B", "b"),
		).
		Selected(0, 1)

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'n'})
	if m.SelectionCount() != 0 {
		t.Errorf("SelectionCount() = %d, want 0", m.SelectionCount())
	}
}

func TestMultiSelect_Update_Confirm(t *testing.T) {
	m := New[string]("Select:").
		Options(
			Opt("A", "a"),
			Opt("B", "b"),
		).
		Selected(0)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("cmd is nil, want ConfirmSelectionCmd")
	}

	msg := cmd()
	confirmMsg, ok := msg.(ConfirmSelectionMsg[string])
	if !ok {
		t.Fatalf("msg type = %T, want ConfirmSelectionMsg[string]", msg)
	}

	want := []string{"a"}
	if !reflect.DeepEqual(confirmMsg.Values, want) {
		t.Errorf("confirmMsg.Values = %v, want %v", confirmMsg.Values, want)
	}
}

func TestMultiSelect_Update_ConfirmBlocked(t *testing.T) {
	// Min=1, selected=0 - cannot confirm
	m := New[string]("Select:").
		Options(Opt("A", "a")).
		Min(1)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Errorf("cmd = %v, want nil (confirm blocked)", cmd)
	}
}

// func TestMultiSelect_Update_FilterInput(t *testing.T) {
// 	m := New[string]("Select:").
// 		Options(
// 			Opt("Apple", "apple"),
// 			Opt("Banana", "banana"),
// 			Opt("Cherry", "cherry"),
// 		).
// 		WithFilterable(true)
// 
// 	// Type 'b'
// 	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'b'})
// 	if m.domain.FilterQuery() != "b" {
// 		t.Errorf("FilterQuery() = %q, want %q", m.domain.FilterQuery(), "b")
// 	}
// 	if m.domain.FilteredCount() != 2 { // Banana
// 		t.Errorf("FilteredCount() = %d, want 1", m.domain.FilteredCount())
// 	}
// 
// 	// Type 'p' (query="ba")
// 	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'p'})
// 	if m.domain.FilterQuery() != "ba" {
// 		t.Errorf("FilterQuery() = %q, want %q", m.domain.FilterQuery(), "ba")
// 	}
// 	if m.domain.FilteredCount() != 1 { // Apple only
// 		t.Errorf("FilteredCount() = %d, want 2", m.domain.FilteredCount())
// 	}
// 
// 	// Backspace
// 	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
// 	if m.domain.FilterQuery() != "b" {
// 		t.Errorf("FilterQuery() = %q, want %q", m.domain.FilterQuery(), "b")
// 	}
// 
// 	// Esc to clear
// 	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
// 	if m.domain.FilterQuery() != "" {
// 		t.Errorf("FilterQuery() = %q, want empty", m.domain.FilterQuery())
// 	}
// }

// func TestMultiSelect_Update_FilterIgnoresBoundKeys(t *testing.T) {
// 	m := New[string]("Select:").
// 		Options(Opt("A", "b"), Opt("B", "b")).
// 		WithFilterable(true)
// 
// 	// Press 'b' (bound to select all)
// 	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'b'})
// 
// 	// Should trigger select all, not filter
// 	if m.SelectionCount() != 2 {
// 		t.Errorf("SelectionCount() = %d, want 1 (select all)", m.SelectionCount())
// 	}
// 	if m.domain.FilterQuery() != "" {
// 		t.Errorf("FilterQuery() = %q, want empty (not filtered)", m.domain.FilterQuery())
// 	}
// }

func TestMultiSelect_View(t *testing.T) {
	m := New[string]("Choose fruits:").
		Options(
			Opt("Apple", "apple"),
			Opt("Banana", "banana"),
		)

	view := m.View()

	// Should contain title with summary
	if !strings.Contains(view, "Choose fruits:") {
		t.Error("View() should contain title")
	}
	if !strings.Contains(view, "(0 of 2 selected)") {
		t.Error("View() should contain selection summary")
	}

	// Should contain help text
	if !strings.Contains(view, "a: all") {
		t.Error("View() should contain help text")
	}
	if !strings.Contains(view, "Space: toggle") {
		t.Error("View() should contain toggle help")
	}
}

func TestMultiSelect_View_WithSelection(t *testing.T) {
	m := New[string]("Select:").
		Options(Opt("A", "b"), Opt("B", "b")).
		Selected(0)

	view := m.View()

	if !strings.Contains(view, "(1 of 2 selected)") {
		t.Errorf("View() should show selection count, got: %s", view)
	}
}

func TestMultiSelect_View_WithFilter(t *testing.T) {
	m := New[string]("Select:").
		Options(Opt("Apple", "apple")).
		WithFilterable(true)

	view := m.View()

	if !strings.Contains(view, "Type to filter...") {
		t.Error("View() should show filter prompt")
	}

	// Type something
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'b'})
	view = m.View()

	if !strings.Contains(view, "Filter: b") {
		t.Errorf("View() should show filter query, got: %s", view)
	}
}

func TestMultiSelect_View_EmptyOptions(t *testing.T) {
	m := New[string]("Select:")

	view := m.View()

	if !strings.Contains(view, "(no options)") {
		t.Error("View() should show empty state")
	}
}

func TestMultiSelect_View_NoMatches(t *testing.T) {
	m := New[string]("Select:").
		Options(Opt("Apple", "apple")).
		WithFilterable(true)

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'z'})
	view := m.View()

	if !strings.Contains(view, "(no matches)") {
		t.Error("View() should show no matches state")
	}
}

func TestOpt(t *testing.T) {
	// Without description
	opt := Opt("Label", 42)
	if opt.Label() != "Label" {
		t.Errorf("Label() = %q, want %q", opt.Label(), "Label")
	}
	if opt.Value() != 42 {
		t.Errorf("Value() = %d, want %d", opt.Value(), 42)
	}
	if opt.Description() != "" {
		t.Errorf("Description() = %q, want empty", opt.Description())
	}

	// With description
	opt = Opt("Label", 42, "A description")
	if opt.Description() != "A description" {
		t.Errorf("Description() = %q, want %q", opt.Description(), "A description")
	}

	// With empty description (should not set it)
	opt = Opt("Label", 42, "")
	if opt.Description() != "" {
		t.Errorf("Description() = %q, want empty", opt.Description())
	}
}

func TestConfirmSelectionCmd(t *testing.T) {
	values := []string{"b", "b", "c"}
	getValues := func() []string { return values }

	cmd := ConfirmSelectionCmd(getValues)
	msg := cmd()

	confirmMsg, ok := msg.(ConfirmSelectionMsg[string])
	if !ok {
		t.Fatalf("msg type = %T, want ConfirmSelectionMsg[string]", msg)
	}

	if !reflect.DeepEqual(confirmMsg.Values, values) {
		t.Errorf("Values = %v, want %v", confirmMsg.Values, values)
	}
}

func TestMultiSelect_Immutability(t *testing.T) {
	original := New[string]("Select:").
		Options(Opt("A", "b"), Opt("B", "b"))

	modified, _ := original.Update(tea.KeyMsg{Type: tea.KeySpace})

	// Original should be unchanged
	if original.SelectionCount() != 0 {
		t.Errorf("original.SelectionCount() = %d, want 0", original.SelectionCount())
	}

	// Modified should have selection
	if modified.SelectionCount() != 1 {
		t.Errorf("modified.SelectionCount() = %d, want 1", modified.SelectionCount())
	}
}
