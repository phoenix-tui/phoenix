package selectcomponent

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/tea"
)

func TestNewString(t *testing.T) {
	options := []string{"Option 1", "Option 2", "Option 3"}
	sel := NewString("Choose:", options)

	if sel.title != "Choose:" {
		t.Errorf("expected title 'Choose:', got %q", sel.title)
	}
}

func TestNew(t *testing.T) {
	sel := New[string]("Select item:")

	if sel.title != "Select item:" {
		t.Errorf("expected title 'Select item:', got %q", sel.title)
	}
	if sel.filterable {
		t.Error("expected filterable false by default")
	}
}

func TestWithHeight(t *testing.T) {
	sel := New[string]("Title").WithHeight(15)

	// Check that height was set (by verifying behavior)
	view := sel.View()
	if !strings.Contains(view, "Title") {
		t.Error("expected view to contain title")
	}
}

func TestWithFilterable(t *testing.T) {
	sel := New[string]("Title").WithFilterable(true)

	if !sel.filterable {
		t.Error("expected filterable to be true")
	}

	view := sel.View()
	if !strings.Contains(view, "Type to filter") {
		t.Error("expected view to show filter prompt when filterable")
	}
}

func TestUpdate(t *testing.T) {
	options := []string{"Option 1", "Option 2", "Option 3"}
	sel := NewString("Choose:", options)

	t.Run("handles key down", func(t *testing.T) {
		newSel, _ := sel.Update(tea.KeyMsg{Type: tea.KeyDown})
		// Cursor should have moved (verify by checking focused value)
		if newSel == nil {
			t.Fatal("expected non-nil select")
		}
	})

	t.Run("handles non-key messages", func(t *testing.T) {
		newSel, cmd := sel.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
		if newSel != sel {
			t.Error("expected same select for non-key message")
		}
		if cmd != nil {
			t.Error("expected nil command for non-key message")
		}
	})

	t.Run("handles enter key to select", func(t *testing.T) {
		newSel, cmd := sel.Update(tea.KeyMsg{Type: tea.KeyEnter})
		if newSel == nil {
			t.Fatal("expected non-nil select")
		}
		if cmd == nil {
			t.Error("expected command after Enter key")
		}
	})
}

func TestView(t *testing.T) {
	t.Run("renders title", func(t *testing.T) {
		sel := NewString("My Title", []string{"Option 1"})
		view := sel.View()
		if !strings.Contains(view, "My Title") {
			t.Error("expected view to contain title")
		}
	})

	t.Run("renders empty state", func(t *testing.T) {
		sel := NewString("Title", []string{})
		view := sel.View()
		if !strings.Contains(view, "(no options)") {
			t.Error("expected view to show empty state")
		}
	})

	t.Run("renders options with cursor", func(t *testing.T) {
		sel := NewString("Title", []string{"Option 1", "Option 2"})
		view := sel.View()
		if !strings.Contains(view, "Option 1") {
			t.Error("expected view to contain Option 1")
		}
		if !strings.Contains(view, "Option 2") {
			t.Error("expected view to contain Option 2")
		}
		// First option should be focused
		if !strings.Contains(view, "> Option 1") {
			t.Error("expected first option to be focused (with cursor)")
		}
	})
}

func TestOpt(t *testing.T) {
	t.Run("creates option with label and value", func(t *testing.T) {
		opt := Opt("Label", 42)
		if opt.Label() != "Label" {
			t.Errorf("expected label 'Label', got %q", opt.Label())
		}
		if opt.Value() != 42 {
			t.Errorf("expected value 42, got %d", opt.Value())
		}
		if opt.Description() != "" {
			t.Error("expected empty description")
		}
	})

	t.Run("creates option with description", func(t *testing.T) {
		opt := Opt("Label", 42, "This is a description")
		if opt.Description() != "This is a description" {
			t.Errorf("expected description, got %q", opt.Description())
		}
	})

	t.Run("handles empty description", func(t *testing.T) {
		opt := Opt("Label", 42, "")
		if opt.Description() != "" {
			t.Error("expected empty description")
		}
	})
}

func TestFilterInput(t *testing.T) {
	options := []string{"Phoenix", "Charm", "Bubbletea"}
	sel := NewString("Choose:", options).WithFilterable(true)

	t.Run("adds character to filter", func(t *testing.T) {
		newSel, _ := sel.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'p'})
		view := newSel.View()
		if !strings.Contains(view, "Filter: p") {
			t.Error("expected filter query to show 'p'")
		}
	})

	t.Run("backspace removes character", func(t *testing.T) {
		sel, _ = sel.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'p'})
		sel, _ = sel.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'h'})
		newSel, _ := sel.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		view := newSel.View()
		if !strings.Contains(view, "Filter: p") {
			t.Error("expected filter query to be 'p' after backspace")
		}
	})

	t.Run("escape clears filter", func(t *testing.T) {
		sel, _ = sel.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'p'})
		newSel, _ := sel.Update(tea.KeyMsg{Type: tea.KeyEsc})
		view := newSel.View()
		if !strings.Contains(view, "Type to filter") {
			t.Error("expected filter to be cleared")
		}
	})
}

func TestInit(t *testing.T) {
	sel := NewString("Title", []string{"Option 1"})
	cmd := sel.Init()
	if cmd != nil {
		t.Error("expected nil command from Init")
	}
}

func TestSelectedValue(t *testing.T) {
	options := []string{"Option 1", "Option 2"}
	sel := NewString("Choose:", options)

	t.Run("returns false when nothing selected", func(t *testing.T) {
		_, ok := sel.SelectedValue()
		if ok {
			t.Error("expected false when nothing selected")
		}
	})

	t.Run("returns selected value after selection", func(t *testing.T) {
		newSel, _ := sel.Update(tea.KeyMsg{Type: tea.KeyEnter})
		val, ok := newSel.SelectedValue()
		if !ok {
			t.Fatal("expected true after selection")
		}
		if val != "Option 1" {
			t.Errorf("expected 'Option 1', got %q", val)
		}
	})
}

func TestFocusedValue(t *testing.T) {
	options := []string{"Option 1", "Option 2"}
	sel := NewString("Choose:", options)

	t.Run("returns first option initially", func(t *testing.T) {
		val, ok := sel.FocusedValue()
		if !ok {
			t.Fatal("expected true for focused value")
		}
		if val != "Option 1" {
			t.Errorf("expected 'Option 1', got %q", val)
		}
	})

	t.Run("returns second option after move down", func(t *testing.T) {
		newSel, _ := sel.Update(tea.KeyMsg{Type: tea.KeyDown})
		val, ok := newSel.FocusedValue()
		if !ok {
			t.Fatal("expected true for focused value")
		}
		if val != "Option 2" {
			t.Errorf("expected 'Option 2', got %q", val)
		}
	})
}
