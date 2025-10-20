package modal

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/components/modal/infrastructure"
	tea "github.com/phoenix-tui/phoenix/tea/api"
)

func TestNew(t *testing.T) {
	modal := New("Test content")

	// Modal is not visible by default.
	if modal.IsVisible() {
		t.Error("Modal should not be visible by default")
	}

	// Show modal and check view.
	modal = modal.Show()
	if !strings.Contains(modal.View(), "Test content") {
		t.Error("Modal should contain the content in view when visible")
	}
}

func TestNewWithTitle(t *testing.T) {
	modal := NewWithTitle("Test Title", "Test content")

	view := modal.Show().View()
	if !strings.Contains(view, "Test Title") {
		t.Error("Modal should contain the title in view")
	}
	if !strings.Contains(view, "Test content") {
		t.Error("Modal should contain the content in view")
	}
}

func TestModalSize(t *testing.T) {
	modal := New("Content").Size(60, 15)

	// Size should be applied (verify via domain access or view rendering)
	if modal.domain.Size().Width() != 60 {
		t.Errorf("Expected width 60, got %d", modal.domain.Size().Width())
	}
	if modal.domain.Size().Height() != 15 {
		t.Errorf("Expected height 15, got %d", modal.domain.Size().Height())
	}
}

func TestModalPosition(t *testing.T) {
	modal := New("Content").Position(10, 5)

	if modal.domain.Position().IsCenter() {
		t.Error("Modal should have custom position, not centered")
	}
	if modal.domain.Position().X() != 10 || modal.domain.Position().Y() != 5 {
		t.Error("Modal position not set correctly")
	}
}

func TestModalCentered(t *testing.T) {
	modal := New("Content").Position(10, 5).Centered()

	if !modal.domain.Position().IsCenter() {
		t.Error("Modal should be centered")
	}
}

func TestModalButtons(t *testing.T) {
	buttons := []Button{
		{Label: "Yes", Key: "y", Action: "confirm"},
		{Label: "No", Key: "n", Action: "cancel"},
	}
	modal := New("Content").Buttons(buttons)

	if len(modal.domain.Buttons()) != 2 {
		t.Errorf("Expected 2 buttons, got %d", len(modal.domain.Buttons()))
	}
}

func TestModalDimBackground(t *testing.T) {
	modal := New("Content").DimBackground(true)

	if !modal.domain.DimBackground() {
		t.Error("Modal should have dimming enabled")
	}
}

func TestModalShowHide(t *testing.T) {
	modal := New("Content")

	// Initially hidden.
	if modal.IsVisible() {
		t.Error("Modal should be hidden initially")
	}

	// Show.
	modal = modal.Show()
	if !modal.IsVisible() {
		t.Error("Modal should be visible after Show()")
	}

	// Hide.
	modal = modal.Hide()
	if modal.IsVisible() {
		t.Error("Modal should be hidden after Hide()")
	}
}

func TestModalKeyBindings(t *testing.T) {
	customKB := infrastructure.KeyBindings{
		Close:          []string{"q"},
		NextButton:     []string{"j"},
		PreviousButton: []string{"k"},
		ActivateButton: []string{"enter"},
	}
	modal := New("Content").KeyBindings(customKB)

	if modal.keyBindings.Close[0] != "q" {
		t.Error("Custom key bindings not applied")
	}
}

func TestModalInit(t *testing.T) {
	modal := New("Content")
	cmd := modal.Init()

	if cmd != nil {
		t.Error("Init should return nil cmd")
	}
}

func TestModalUpdateWindowSize(t *testing.T) {
	modal := New("Content").Show()

	msg := tea.WindowSizeMsg{Width: 100, Height: 40}
	updated, _ := modal.Update(msg)

	if updated.terminalWidth != 100 || updated.terminalHeight != 40 {
		t.Error("Terminal size not updated")
	}
}

func TestModalUpdateKeyClose(t *testing.T) {
	modal := New("Content").Show()

	// Press Esc.
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updated, _ := modal.Update(msg)

	if updated.IsVisible() {
		t.Error("Modal should be hidden after Esc")
	}
}

func TestModalUpdateKeyNextButton(t *testing.T) {
	buttons := []Button{
		{Label: "Yes", Key: "y", Action: "confirm"},
		{Label: "No", Key: "n", Action: "cancel"},
	}
	modal := New("Content").Buttons(buttons).Show()

	// Initially focused on first button.
	if modal.FocusedButton() != "confirm" {
		t.Error("Should focus on first button initially")
	}

	// Press Tab to focus next button.
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updated, _ := modal.Update(msg)

	if updated.FocusedButton() != "cancel" {
		t.Error("Should focus on second button after Tab")
	}
}

func TestModalUpdateKeyPreviousButton(t *testing.T) {
	buttons := []Button{
		{Label: "Yes", Key: "y", Action: "confirm"},
		{Label: "No", Key: "n", Action: "cancel"},
	}
	modal := New("Content").Buttons(buttons).Show()

	// Press Shift+Tab to focus previous button (wraps to last)
	msg := tea.KeyMsg{Type: tea.KeyTab, Shift: true}
	updated, _ := modal.Update(msg)

	if updated.FocusedButton() != "cancel" {
		t.Error("Should wrap to last button with Shift+Tab")
	}
}

func TestModalUpdateKeyActivateButton(t *testing.T) {
	buttons := []Button{
		{Label: "Yes", Key: "y", Action: "confirm"},
		{Label: "No", Key: "n", Action: "cancel"},
	}
	modal := New("Content").Buttons(buttons).Show()

	// Press Enter to activate focused button.
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := modal.Update(msg)

	if cmd == nil {
		t.Fatal("Should return command with ButtonPressedMsg")
	}

	// Execute command and verify message.
	result := cmd()
	btnMsg, ok := result.(ButtonPressedMsg)
	if !ok {
		t.Fatal("Command should return ButtonPressedMsg")
	}
	if btnMsg.Action != "confirm" {
		t.Errorf("Expected action 'confirm', got '%s'", btnMsg.Action)
	}
}

func TestModalUpdateKeyShortcut(t *testing.T) {
	buttons := []Button{
		{Label: "Yes", Key: "y", Action: "confirm"},
		{Label: "No", Key: "n", Action: "cancel"},
	}
	modal := New("Content").Buttons(buttons).Show()

	// Press 'n' shortcut.
	msg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'n'}
	_, cmd := modal.Update(msg)

	if cmd == nil {
		t.Fatal("Should return command with ButtonPressedMsg")
	}

	// Execute command and verify message.
	result := cmd()
	btnMsg, ok := result.(ButtonPressedMsg)
	if !ok {
		t.Fatal("Command should return ButtonPressedMsg")
	}
	if btnMsg.Action != "cancel" {
		t.Errorf("Expected action 'cancel', got '%s'", btnMsg.Action)
	}
}

func TestModalUpdateKeyShortcutCaseInsensitive(t *testing.T) {
	buttons := []Button{
		{Label: "Yes", Key: "y", Action: "confirm"},
	}
	modal := New("Content").Buttons(buttons).Show()

	// Press 'Y' (uppercase) - should match 'y' shortcut.
	msg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'Y'}
	_, cmd := modal.Update(msg)

	if cmd == nil {
		t.Fatal("Should return command with ButtonPressedMsg (case-insensitive)")
	}

	result := cmd()
	btnMsg, ok := result.(ButtonPressedMsg)
	if !ok {
		t.Fatal("Command should return ButtonPressedMsg")
	}
	if btnMsg.Action != "confirm" {
		t.Errorf("Expected action 'confirm', got '%s'", btnMsg.Action)
	}
}

func TestModalUpdateHidden(t *testing.T) {
	modal := New("Content") // Not visible

	// Press keys - should have no effect.
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updated, cmd := modal.Update(msg)

	if updated.IsVisible() {
		t.Error("Modal should remain hidden")
	}
	if cmd != nil {
		t.Error("Hidden modal should not process input")
	}
}

func TestModalViewHidden(t *testing.T) {
	modal := New("Content") // Not visible

	view := modal.View()
	if view != "" {
		t.Error("Hidden modal should render empty string")
	}
}

func TestModalViewVisible(t *testing.T) {
	modal := New("Content").Show()

	view := modal.View()
	if view == "" {
		t.Error("Visible modal should render content")
	}
	if !strings.Contains(view, "Content") {
		t.Error("Modal view should contain content")
	}
}

func TestModalViewWithTitle(t *testing.T) {
	modal := NewWithTitle("Title", "Content").Show()

	view := modal.View()
	if !strings.Contains(view, "Title") {
		t.Error("Modal view should contain title")
	}
	if !strings.Contains(view, "Content") {
		t.Error("Modal view should contain content")
	}
}

func TestModalViewWithButtons(t *testing.T) {
	buttons := []Button{
		{Label: "Yes", Key: "y", Action: "confirm"},
		{Label: "No", Key: "n", Action: "cancel"},
	}
	modal := New("Content").Buttons(buttons).Show()

	view := modal.View()
	if !strings.Contains(view, "Yes") {
		t.Error("Modal view should contain button labels")
	}
	if !strings.Contains(view, "No") {
		t.Error("Modal view should contain button labels")
	}
}

func TestModalViewDimmedBackground(t *testing.T) {
	modal := New("Content").DimBackground(true).Show()

	view := modal.View()
	if !strings.Contains(view, "░") {
		t.Error("Modal view should contain dimmed background character")
	}
}

func TestModalFluentAPI(t *testing.T) {
	// Test method chaining.
	modal := New("Content").
		Size(60, 15).
		DimBackground(true).
		Buttons([]Button{
			{Label: "OK", Key: "enter", Action: "ok"},
		}).
		Show()

	if !modal.IsVisible() {
		t.Error("Fluent API should preserve all operations")
	}
	if !modal.domain.DimBackground() {
		t.Error("Fluent API should preserve dimming")
	}
	if len(modal.domain.Buttons()) != 1 {
		t.Error("Fluent API should preserve buttons")
	}
}

func TestModalImmutability(t *testing.T) {
	original := New("Original")
	modified := original.Show()

	// Original should be unchanged.
	if original.IsVisible() {
		t.Error("Original modal should remain hidden")
	}

	// Modified should have changes.
	if !modified.IsVisible() {
		t.Error("Modified modal should be visible")
	}
}

func TestModalFocusedButtonNoButtons(t *testing.T) {
	modal := New("Content")

	focused := modal.FocusedButton()
	if focused != "" {
		t.Error("Should return empty string when no buttons")
	}
}

func TestModalFocusedButton(t *testing.T) {
	buttons := []Button{
		{Label: "Yes", Key: "y", Action: "confirm"},
		{Label: "No", Key: "n", Action: "cancel"},
	}
	modal := New("Content").Buttons(buttons)

	focused := modal.FocusedButton()
	if focused != "confirm" {
		t.Errorf("Expected focused button 'confirm', got '%s'", focused)
	}
}
