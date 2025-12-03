package confirm

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/tea"
)

func TestNew(t *testing.T) {
	c := New("Delete file?")

	if !strings.Contains(c.View(), "Delete file?") {
		t.Error("View() should contain title")
	}

	if c.Done() {
		t.Error("Done() should be false initially")
	}
}

func TestConfirm_Description(t *testing.T) {
	c := New("Delete?").Description("This cannot be undone.")

	view := c.View()
	if !strings.Contains(view, "This cannot be undone.") {
		t.Error("View() should contain description")
	}
}

func TestConfirm_Affirmative(t *testing.T) {
	c := New("Proceed?").Affirmative("Delete")

	view := c.View()
	if !strings.Contains(view, "Delete") {
		t.Error("View() should contain custom affirmative label")
	}
}

func TestConfirm_Negative(t *testing.T) {
	c := New("Proceed?").Negative("Keep")

	view := c.View()
	if !strings.Contains(view, "Keep") {
		t.Error("View() should contain custom negative label")
	}
}

func TestConfirm_WithCancel(t *testing.T) {
	c := New("Save?").WithCancel(true)

	view := c.View()
	if !strings.Contains(view, "Cancel") {
		t.Error("View() should contain Cancel button")
	}
}

func TestConfirm_DefaultYes(t *testing.T) {
	c := New("Confirm?").DefaultYes()

	view := c.View()
	// First button should have [ ] (focused) instead of ( )
	if !strings.Contains(view, "[ Yes ]") {
		t.Error("View() should show Yes button as focused")
	}
}

func TestConfirm_DefaultNo(t *testing.T) {
	c := New("Confirm?").DefaultNo()

	view := c.View()
	// Second button should have [ ] (focused)
	if !strings.Contains(view, "[ No ]") {
		t.Error("View() should show No button as focused")
	}
}

func TestConfirm_Update_MoveLeft(t *testing.T) {
	c := New("Confirm?").DefaultNo() // Start at No

	newC, _ := c.Update(tea.KeyMsg{Type: tea.KeyLeft})

	view := newC.View()
	if !strings.Contains(view, "[ Yes ]") {
		t.Error("After left arrow, Yes should be focused")
	}
}

func TestConfirm_Update_MoveRight(t *testing.T) {
	c := New("Confirm?").DefaultYes() // Start at Yes

	newC, _ := c.Update(tea.KeyMsg{Type: tea.KeyRight})

	view := newC.View()
	if !strings.Contains(view, "[ No ]") {
		t.Error("After right arrow, No should be focused")
	}
}

func TestConfirm_Update_Tab(t *testing.T) {
	c := New("Confirm?").DefaultYes()

	newC, _ := c.Update(tea.KeyMsg{Type: tea.KeyTab})

	view := newC.View()
	if !strings.Contains(view, "[ No ]") {
		t.Error("After Tab, No should be focused")
	}
}

func TestConfirm_Update_ShiftTab(t *testing.T) {
	c := New("Confirm?").DefaultNo()

	newC, _ := c.Update(tea.KeyMsg{Type: tea.KeyTab, Shift: true})

	view := newC.View()
	if !strings.Contains(view, "[ Yes ]") {
		t.Error("After Shift+Tab, Yes should be focused")
	}
}

func TestConfirm_Update_EnterYes(t *testing.T) {
	c := New("Confirm?").DefaultYes()

	newC, cmd := c.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !newC.Done() {
		t.Error("Done() should be true after Enter")
	}

	if !newC.IsYes() {
		t.Error("IsYes() should be true after confirming Yes")
	}

	if newC.IsNo() {
		t.Error("IsNo() should be false")
	}

	if newC.IsCanceled() {
		t.Error("IsCanceled() should be false")
	}

	// Check command sends message
	if cmd == nil {
		t.Fatal("Command should not be nil")
	}

	msg := cmd()
	if _, ok := msg.(ConfirmResultMsg); !ok {
		t.Error("Command should return ConfirmResultMsg")
	}
}

func TestConfirm_Update_EnterNo(t *testing.T) {
	c := New("Confirm?").DefaultNo()

	newC, _ := c.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !newC.Done() {
		t.Error("Done() should be true after Enter")
	}

	if !newC.IsNo() {
		t.Error("IsNo() should be true after confirming No")
	}

	if newC.IsYes() {
		t.Error("IsYes() should be false")
	}
}

func TestConfirm_Update_Escape(t *testing.T) {
	c := New("Confirm?")

	newC, cmd := c.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if !newC.Done() {
		t.Error("Done() should be true after Escape")
	}

	if !newC.IsCanceled() {
		t.Error("IsCanceled() should be true after Escape")
	}

	if cmd == nil {
		t.Fatal("Command should not be nil")
	}

	msg := cmd()
	if _, ok := msg.(ConfirmResultMsg); !ok {
		t.Error("Command should return ConfirmResultMsg")
	}
}

func TestConfirm_Update_KeyboardShortcut_Y(t *testing.T) {
	c := New("Confirm?")

	newC, cmd := c.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'y'})

	if !newC.Done() {
		t.Error("Done() should be true after pressing 'y'")
	}

	if !newC.IsYes() {
		t.Error("IsYes() should be true after pressing 'y'")
	}

	if cmd == nil {
		t.Error("Command should not be nil")
	}
}

func TestConfirm_Update_KeyboardShortcut_N(t *testing.T) {
	c := New("Confirm?")

	newC, cmd := c.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'n'})

	if !newC.Done() {
		t.Error("Done() should be true after pressing 'n'")
	}

	if !newC.IsNo() {
		t.Error("IsNo() should be true after pressing 'n'")
	}

	if cmd == nil {
		t.Error("Command should not be nil")
	}
}

func TestConfirm_Update_KeyboardShortcut_UppercaseY(t *testing.T) {
	c := New("Confirm?")

	newC, _ := c.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'Y'})

	if !newC.IsYes() {
		t.Error("IsYes() should be true after pressing 'Y' (uppercase)")
	}
}

func TestConfirm_Update_KeyboardShortcut_InvalidKey(t *testing.T) {
	c := New("Confirm?")

	newC, cmd := c.Update(tea.KeyMsg{Type: tea.KeyRune, Rune: 'x'})

	if newC.Done() {
		t.Error("Done() should be false after invalid key")
	}

	if cmd != nil {
		t.Error("Command should be nil for invalid key")
	}
}

func TestConfirm_Update_AfterDone(t *testing.T) {
	c := New("Confirm?")

	// Confirm selection
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !c.Done() {
		t.Fatal("Should be done after Enter")
	}

	// Try to change selection after done
	c2, cmd := c.Update(tea.KeyMsg{Type: tea.KeyLeft})

	// Should not change
	if c2 != c {
		t.Error("Update() should not change confirm after done")
	}

	if cmd != nil {
		t.Error("Command should be nil after done")
	}
}

func TestConfirm_View_ButtonRendering(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *Confirm
		want  []string // Substrings that should appear
	}{
		{
			name: "Yes focused",
			setup: func() *Confirm {
				return New("Test?").DefaultYes()
			},
			want: []string{"[ Yes ]", "( No )"},
		},
		{
			name: "No focused",
			setup: func() *Confirm {
				return New("Test?").DefaultNo()
			},
			want: []string{"( Yes )", "[ No ]"},
		},
		{
			name: "Three buttons with Cancel focused",
			setup: func() *Confirm {
				c := New("Test?").WithCancel(true).DefaultYes()
				c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRight})
				c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRight})
				return c
			},
			want: []string{"( Yes )", "( No )", "[ Cancel ]"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setup()
			view := c.View()

			for _, substr := range tt.want {
				if !strings.Contains(view, substr) {
					t.Errorf("View() should contain %q, got:\n%s", substr, view)
				}
			}
		})
	}
}

func TestConfirm_CustomLabels(t *testing.T) {
	c := New("Proceed?").
		Affirmative("Delete").
		Negative("Keep")

	view := c.View()

	if !strings.Contains(view, "Delete") {
		t.Error("View() should contain custom affirmative label 'Delete'")
	}

	if !strings.Contains(view, "Keep") {
		t.Error("View() should contain custom negative label 'Keep'")
	}

	// Confirm Delete
	c = c.DefaultYes()
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !c.IsYes() {
		t.Error("IsYes() should be true after confirming custom 'Delete' button")
	}
}
