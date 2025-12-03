package model

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/select/internal/domain/value"
)

func TestNew(t *testing.T) {
	t.Run("creates select with options", func(t *testing.T) {
		opts := []*value.Option[string]{
			value.NewOption("Option 1", "opt1"),
			value.NewOption("Option 2", "opt2"),
		}
		sel := New(opts)

		if sel == nil {
			t.Fatal("expected non-nil Select")
		}
		if len(sel.options) != 2 {
			t.Errorf("expected 2 options, got %d", len(sel.options))
		}
		if sel.cursor != 0 {
			t.Errorf("expected cursor at 0, got %d", sel.cursor)
		}
		if sel.selectedIndex != -1 {
			t.Errorf("expected selectedIndex -1, got %d", sel.selectedIndex)
		}
	})

	t.Run("handles nil options", func(t *testing.T) {
		sel := New[string](nil)
		if sel == nil {
			t.Fatal("expected non-nil Select")
		}
		if len(sel.options) != 0 {
			t.Errorf("expected 0 options, got %d", len(sel.options))
		}
	})
}

func TestMoveUp(t *testing.T) {
	opts := []*value.Option[string]{
		value.NewOption("Option 1", "opt1"),
		value.NewOption("Option 2", "opt2"),
		value.NewOption("Option 3", "opt3"),
	}
	sel := New(opts)

	t.Run("moves cursor up from position 2", func(t *testing.T) {
		sel = sel.withCursor(2)
		sel = sel.MoveUp()
		if sel.cursor != 1 {
			t.Errorf("expected cursor at 1, got %d", sel.cursor)
		}
	})

	t.Run("does not move up from position 0", func(t *testing.T) {
		sel = sel.withCursor(0)
		sel = sel.MoveUp()
		if sel.cursor != 0 {
			t.Errorf("expected cursor at 0, got %d", sel.cursor)
		}
	})
}

func TestMoveDown(t *testing.T) {
	opts := []*value.Option[string]{
		value.NewOption("Option 1", "opt1"),
		value.NewOption("Option 2", "opt2"),
		value.NewOption("Option 3", "opt3"),
	}
	sel := New(opts)

	t.Run("moves cursor down from position 0", func(t *testing.T) {
		sel = sel.MoveDown()
		if sel.cursor != 1 {
			t.Errorf("expected cursor at 1, got %d", sel.cursor)
		}
	})

	t.Run("does not move down from last position", func(t *testing.T) {
		sel = sel.withCursor(2)
		sel = sel.MoveDown()
		if sel.cursor != 2 {
			t.Errorf("expected cursor at 2, got %d", sel.cursor)
		}
	})
}

func TestMoveToStart(t *testing.T) {
	opts := []*value.Option[string]{
		value.NewOption("Option 1", "opt1"),
		value.NewOption("Option 2", "opt2"),
		value.NewOption("Option 3", "opt3"),
	}
	sel := New(opts).withCursor(2)

	sel = sel.MoveToStart()
	if sel.cursor != 0 {
		t.Errorf("expected cursor at 0, got %d", sel.cursor)
	}
}

func TestMoveToEnd(t *testing.T) {
	opts := []*value.Option[string]{
		value.NewOption("Option 1", "opt1"),
		value.NewOption("Option 2", "opt2"),
		value.NewOption("Option 3", "opt3"),
	}
	sel := New(opts)

	sel = sel.MoveToEnd()
	if sel.cursor != 2 {
		t.Errorf("expected cursor at 2, got %d", sel.cursor)
	}
}

func TestSelect(t *testing.T) {
	opts := []*value.Option[string]{
		value.NewOption("Option 1", "opt1"),
		value.NewOption("Option 2", "opt2"),
		value.NewOption("Option 3", "opt3"),
	}
	sel := New(opts).withCursor(1)

	t.Run("selects focused option", func(t *testing.T) {
		sel = sel.Select()
		if sel.selectedIndex != 1 {
			t.Errorf("expected selectedIndex 1, got %d", sel.selectedIndex)
		}
	})

	t.Run("returns selected value", func(t *testing.T) {
		val, ok := sel.SelectedValue()
		if !ok {
			t.Fatal("expected selected value")
		}
		if val != "opt2" {
			t.Errorf("expected 'opt2', got %q", val)
		}
	})
}

func TestSetFilterQuery(t *testing.T) {
	opts := []*value.Option[string]{
		value.NewOption("Phoenix", "phoenix"),
		value.NewOption("Charm", "charm"),
		value.NewOption("Bubbletea", "bubbletea"),
	}
	sel := New(opts)

	t.Run("filters options by query", func(t *testing.T) {
		sel = sel.SetFilterQuery("ph")
		if len(sel.filteredOpts) != 1 {
			t.Errorf("expected 1 filtered option, got %d", len(sel.filteredOpts))
		}
		if sel.filteredOpts[0].Label() != "Phoenix" {
			t.Errorf("expected 'Phoenix', got %q", sel.filteredOpts[0].Label())
		}
	})

	t.Run("adjusts cursor when filtered list is shorter", func(t *testing.T) {
		sel = sel.withCursor(2)
		sel = sel.SetFilterQuery("ph")
		if sel.cursor != 0 {
			t.Errorf("expected cursor adjusted to 0, got %d", sel.cursor)
		}
	})

	t.Run("empty query shows all options", func(t *testing.T) {
		sel = sel.SetFilterQuery("ph")
		sel = sel.SetFilterQuery("")
		if len(sel.filteredOpts) != 3 {
			t.Errorf("expected 3 options, got %d", len(sel.filteredOpts))
		}
	})
}

func TestClearFilter(t *testing.T) {
	opts := []*value.Option[string]{
		value.NewOption("Phoenix", "phoenix"),
		value.NewOption("Charm", "charm"),
	}
	sel := New(opts).SetFilterQuery("ph")

	sel = sel.ClearFilter()
	if sel.filterQuery != "" {
		t.Errorf("expected empty filter query, got %q", sel.filterQuery)
	}
	if len(sel.filteredOpts) != 2 {
		t.Errorf("expected 2 options after clear, got %d", len(sel.filteredOpts))
	}
}

func TestRenderVisibleOptions(t *testing.T) {
	opts := []*value.Option[string]{
		value.NewOption("Option 1", "opt1"),
		value.NewOption("Option 2", "opt2"),
		value.NewOption("Option 3", "opt3"),
		value.NewOption("Option 4", "opt4"),
		value.NewOption("Option 5", "opt5"),
	}
	sel := New(opts).WithHeight(3)

	t.Run("renders visible portion", func(t *testing.T) {
		rendered := sel.RenderVisibleOptions()
		if len(rendered) != 3 {
			t.Errorf("expected 3 rendered options, got %d", len(rendered))
		}
	})

	t.Run("scrolls when cursor moves down", func(t *testing.T) {
		sel = sel.withCursor(4) // Last option
		rendered := sel.RenderVisibleOptions()
		if len(rendered) != 3 {
			t.Errorf("expected 3 rendered options, got %d", len(rendered))
		}
		// Should render options 2, 3, 4 (indices relative to scroll offset)
	})
}

func TestFocusedValue(t *testing.T) {
	opts := []*value.Option[string]{
		value.NewOption("Option 1", "opt1"),
		value.NewOption("Option 2", "opt2"),
	}
	sel := New(opts).withCursor(1)

	val, ok := sel.FocusedValue()
	if !ok {
		t.Fatal("expected focused value")
	}
	if val != "opt2" {
		t.Errorf("expected 'opt2', got %q", val)
	}
}

func TestIsFiltered(t *testing.T) {
	opts := []*value.Option[string]{
		value.NewOption("Option 1", "opt1"),
	}
	sel := New(opts)

	if sel.IsFiltered() {
		t.Error("expected not filtered initially")
	}

	sel = sel.SetFilterQuery("opt")
	if !sel.IsFiltered() {
		t.Error("expected filtered after setting query")
	}

	sel = sel.ClearFilter()
	if sel.IsFiltered() {
		t.Error("expected not filtered after clear")
	}
}

func TestWithHeight(t *testing.T) {
	opts := []*value.Option[string]{
		value.NewOption("Option 1", "opt1"),
	}
	sel := New(opts)

	sel = sel.WithHeight(20)
	if sel.height != 20 {
		t.Errorf("expected height 20, got %d", sel.height)
	}

	sel = sel.WithHeight(0)
	if sel.height != 1 {
		t.Errorf("expected height clamped to 1, got %d", sel.height)
	}
}

func TestWithFilterFunc(t *testing.T) {
	opts := []*value.Option[string]{
		value.NewOption("Phoenix", "phoenix"),
		value.NewOption("Charm", "charm"),
	}
	sel := New(opts)

	// Custom filter that only matches exact case
	customFilter := func(opt *value.Option[string], query string) bool {
		if query == "" {
			return true
		}
		return opt.Label() == query
	}

	sel = sel.WithFilterFunc(customFilter)
	sel = sel.SetFilterQuery("Phoenix")

	if len(sel.filteredOpts) != 1 {
		t.Errorf("expected 1 match, got %d", len(sel.filteredOpts))
	}

	sel = sel.SetFilterQuery("phoenix") // Wrong case
	if len(sel.filteredOpts) != 0 {
		t.Errorf("expected 0 matches with wrong case, got %d", len(sel.filteredOpts))
	}
}

func TestWithRenderFunc(t *testing.T) {
	opts := []*value.Option[string]{
		value.NewOption("Option 1", "opt1"),
	}
	sel := New(opts)

	customRender := func(opt *value.Option[string], _ int, focused bool) string {
		if focused {
			return "[x] " + opt.Label()
		}
		return "[ ] " + opt.Label()
	}

	sel = sel.WithRenderFunc(customRender)
	rendered := sel.RenderVisibleOptions()

	if len(rendered) != 1 {
		t.Fatalf("expected 1 rendered option, got %d", len(rendered))
	}

	if rendered[0] != "[x] Option 1" {
		t.Errorf("expected '[x] Option 1', got %q", rendered[0])
	}
}
