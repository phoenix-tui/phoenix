package model

import (
	"testing"
)

func TestNew(t *testing.T) {
	c := New("Delete file?")

	if c.Title() != "Delete file?" {
		t.Errorf("Title() = %v, want %v", c.Title(), "Delete file?")
	}

	if c.FocusedIndex() != 1 {
		t.Errorf("FocusedIndex() = %v, want 1 (No button)", c.FocusedIndex())
	}

	buttons := c.Buttons()
	if len(buttons) != 2 {
		t.Errorf("Buttons() length = %v, want 2", len(buttons))
	}

	if buttons[0] != "Yes" {
		t.Errorf("Buttons()[0] = %v, want Yes", buttons[0])
	}

	if buttons[1] != "No" {
		t.Errorf("Buttons()[1] = %v, want No", buttons[1])
	}

	if c.Done() {
		t.Error("Done() = true, want false")
	}

	if c.Result() != ResultNone {
		t.Errorf("Result() = %v, want ResultNone", c.Result())
	}
}

func TestConfirm_WithDescription(t *testing.T) {
	c := New("Delete?").WithDescription("This action cannot be undone.")

	if c.Description() != "This action cannot be undone." {
		t.Errorf("Description() = %v, want 'This action cannot be undone.'", c.Description())
	}
}

func TestConfirm_WithButtons(t *testing.T) {
	c := New("Save changes?").WithButtons("Save", "Discard", "Cancel")

	buttons := c.Buttons()
	if len(buttons) != 3 {
		t.Errorf("Buttons() length = %v, want 3", len(buttons))
	}

	if buttons[0] != "Save" {
		t.Errorf("Buttons()[0] = %v, want Save", buttons[0])
	}

	if buttons[1] != "Discard" {
		t.Errorf("Buttons()[1] = %v, want Discard", buttons[1])
	}

	if buttons[2] != "Cancel" {
		t.Errorf("Buttons()[2] = %v, want Cancel", buttons[2])
	}
}

func TestConfirm_WithDefaultYes(t *testing.T) {
	c := New("Confirm?").WithDefaultYes()

	if c.FocusedIndex() != 0 {
		t.Errorf("FocusedIndex() = %v, want 0 (Yes button)", c.FocusedIndex())
	}
}

func TestConfirm_WithDefaultNo(t *testing.T) {
	c := New("Confirm?").WithDefaultNo()

	if c.FocusedIndex() != 1 {
		t.Errorf("FocusedIndex() = %v, want 1 (No button)", c.FocusedIndex())
	}
}

func TestConfirm_MoveFocusLeft(t *testing.T) {
	c := New("Confirm?") // Starts at index 1 (No)

	// Move left to Yes
	c = c.MoveFocusLeft()
	if c.FocusedIndex() != 0 {
		t.Errorf("FocusedIndex() after left = %v, want 0", c.FocusedIndex())
	}

	// Move left again (should wrap to last button)
	c = c.MoveFocusLeft()
	if c.FocusedIndex() != 1 {
		t.Errorf("FocusedIndex() after wrap = %v, want 1", c.FocusedIndex())
	}
}

func TestConfirm_MoveFocusRight(t *testing.T) {
	c := New("Confirm?").WithDefaultYes() // Start at index 0 (Yes)

	// Move right to No
	c = c.MoveFocusRight()
	if c.FocusedIndex() != 1 {
		t.Errorf("FocusedIndex() after right = %v, want 1", c.FocusedIndex())
	}

	// Move right again (should wrap to first button)
	c = c.MoveFocusRight()
	if c.FocusedIndex() != 0 {
		t.Errorf("FocusedIndex() after wrap = %v, want 0", c.FocusedIndex())
	}
}

func TestConfirm_Confirm(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *Confirm
		wantResult Result
		wantDone   bool
	}{
		{
			name: "Confirm Yes",
			setup: func() *Confirm {
				return New("Confirm?").WithDefaultYes()
			},
			wantResult: ResultYes,
			wantDone:   true,
		},
		{
			name: "Confirm No",
			setup: func() *Confirm {
				return New("Confirm?").WithDefaultNo()
			},
			wantResult: ResultNo,
			wantDone:   true,
		},
		{
			name: "Confirm Cancel",
			setup: func() *Confirm {
				return New("Save?").WithButtons("Yes", "No", "Cancel").
					WithDefaultYes().MoveFocusRight().MoveFocusRight()
			},
			wantResult: ResultCanceled,
			wantDone:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setup().Confirm()

			if c.Result() != tt.wantResult {
				t.Errorf("Result() = %v, want %v", c.Result(), tt.wantResult)
			}

			if c.Done() != tt.wantDone {
				t.Errorf("Done() = %v, want %v", c.Done(), tt.wantDone)
			}
		})
	}
}

func TestConfirm_ConfirmKey(t *testing.T) {
	tests := []struct {
		name       string
		key        rune
		wantResult Result
		wantDone   bool
	}{
		{
			name:       "Press y for Yes",
			key:        'y',
			wantResult: ResultYes,
			wantDone:   true,
		},
		{
			name:       "Press Y for Yes (uppercase)",
			key:        'Y',
			wantResult: ResultYes,
			wantDone:   true,
		},
		{
			name:       "Press n for No",
			key:        'n',
			wantResult: ResultNo,
			wantDone:   true,
		},
		{
			name:       "Press N for No (uppercase)",
			key:        'N',
			wantResult: ResultNo,
			wantDone:   true,
		},
		{
			name:       "Press invalid key",
			key:        'x',
			wantResult: ResultNone,
			wantDone:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New("Confirm?").ConfirmKey(tt.key)

			if c.Result() != tt.wantResult {
				t.Errorf("Result() = %v, want %v", c.Result(), tt.wantResult)
			}

			if c.Done() != tt.wantDone {
				t.Errorf("Done() = %v, want %v", c.Done(), tt.wantDone)
			}
		})
	}
}

func TestConfirm_Cancel(t *testing.T) {
	c := New("Confirm?").Cancel()

	if c.Result() != ResultCanceled {
		t.Errorf("Result() = %v, want ResultCanceled", c.Result())
	}

	if !c.Done() {
		t.Error("Done() = false, want true")
	}
}

func TestConfirm_CustomLabels(t *testing.T) {
	c := New("Proceed?").WithButtons("Delete", "Keep")

	buttons := c.Buttons()
	if len(buttons) != 2 {
		t.Errorf("Buttons() length = %v, want 2", len(buttons))
	}

	if buttons[0] != "Delete" {
		t.Errorf("Buttons()[0] = %v, want Delete", buttons[0])
	}

	if buttons[1] != "Keep" {
		t.Errorf("Buttons()[1] = %v, want Keep", buttons[1])
	}

	// Confirm first button (Delete) - should map to ResultYes
	c = c.WithDefaultYes().Confirm()
	if c.Result() != ResultYes {
		t.Errorf("Result() for custom 'Delete' button = %v, want ResultYes", c.Result())
	}
}
