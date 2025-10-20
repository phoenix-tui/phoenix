package infrastructure

import (
	"testing"

	tea "github.com/phoenix-tui/phoenix/tea/api"
)

func TestDefaultKeyBindings(t *testing.T) {
	kb := DefaultKeyBindings()

	// Verify default close keys.
	if len(kb.Close) != 1 || kb.Close[0] != "esc" {
		t.Errorf("Expected Close to be [esc], got %v", kb.Close)
	}

	// Verify default next button keys.
	if len(kb.NextButton) != 2 || kb.NextButton[0] != "tab" || kb.NextButton[1] != "→" {
		t.Errorf("Expected NextButton to be [tab, →], got %v", kb.NextButton)
	}

	// Verify default previous button keys.
	if len(kb.PreviousButton) != 2 || kb.PreviousButton[0] != "shift+tab" || kb.PreviousButton[1] != "←" {
		t.Errorf("Expected PreviousButton to be [shift+tab, ←], got %v", kb.PreviousButton)
	}

	// Verify default activate button keys.
	if len(kb.ActivateButton) != 1 || kb.ActivateButton[0] != "enter" {
		t.Errorf("Expected ActivateButton to be [enter], got %v", kb.ActivateButton)
	}
}

func TestIsClose(t *testing.T) {
	kb := DefaultKeyBindings()

	// Esc should match.
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	if !kb.IsClose(escMsg) {
		t.Error("Esc key should be recognized as Close")
	}

	// Other keys should not match.
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	if kb.IsClose(enterMsg) {
		t.Error("Enter key should not be recognized as Close")
	}
}

func TestIsNextButton(t *testing.T) {
	kb := DefaultKeyBindings()

	// Tab should match.
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	if !kb.IsNextButton(tabMsg) {
		t.Error("Tab key should be recognized as NextButton")
	}

	// Right arrow should match.
	rightMsg := tea.KeyMsg{Type: tea.KeyRight}
	if !kb.IsNextButton(rightMsg) {
		t.Error("Right arrow should be recognized as NextButton")
	}

	// Other keys should not match.
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	if kb.IsNextButton(escMsg) {
		t.Error("Esc key should not be recognized as NextButton")
	}
}

func TestIsPreviousButton(t *testing.T) {
	kb := DefaultKeyBindings()

	// Shift+Tab should match.
	shiftTabMsg := tea.KeyMsg{Type: tea.KeyTab, Shift: true}
	if !kb.IsPreviousButton(shiftTabMsg) {
		t.Error("Shift+Tab should be recognized as PreviousButton")
	}

	// Left arrow should match.
	leftMsg := tea.KeyMsg{Type: tea.KeyLeft}
	if !kb.IsPreviousButton(leftMsg) {
		t.Error("Left arrow should be recognized as PreviousButton")
	}

	// Other keys should not match.
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	if kb.IsPreviousButton(tabMsg) {
		t.Error("Tab key should not be recognized as PreviousButton")
	}
}

func TestIsActivateButton(t *testing.T) {
	kb := DefaultKeyBindings()

	// Enter should match.
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	if !kb.IsActivateButton(enterMsg) {
		t.Error("Enter key should be recognized as ActivateButton")
	}

	// Other keys should not match.
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	if kb.IsActivateButton(escMsg) {
		t.Error("Esc key should not be recognized as ActivateButton")
	}
}

func TestCustomKeyBindings(t *testing.T) {
	kb := KeyBindings{
		Close:          []string{"q", "esc"},
		NextButton:     []string{"j", "↓"},
		PreviousButton: []string{"k", "↑"},
		ActivateButton: []string{"enter", "space"},
	}

	// Test custom close keys.
	qMsg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'q'}
	if !kb.IsClose(qMsg) {
		t.Error("'q' should be recognized as Close")
	}

	// Test custom next button keys.
	jMsg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'j'}
	if !kb.IsNextButton(jMsg) {
		t.Error("'j' should be recognized as NextButton")
	}

	// Test custom previous button keys.
	kMsg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'k'}
	if !kb.IsPreviousButton(kMsg) {
		t.Error("'k' should be recognized as PreviousButton")
	}

	// Test custom activate button keys.
	spaceMsg := tea.KeyMsg{Type: tea.KeySpace}
	if !kb.IsActivateButton(spaceMsg) {
		t.Error("Space should be recognized as ActivateButton")
	}
}

func TestMultipleKeysForSameAction(t *testing.T) {
	kb := KeyBindings{
		Close: []string{"esc", "q", "ctrl+c"},
	}

	// All keys should match.
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	if !kb.IsClose(escMsg) {
		t.Error("Esc should match Close")
	}

	qMsg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'q'}
	if !kb.IsClose(qMsg) {
		t.Error("'q' should match Close")
	}

	ctrlCMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	if !kb.IsClose(ctrlCMsg) {
		t.Error("Ctrl+C should match Close")
	}
}
