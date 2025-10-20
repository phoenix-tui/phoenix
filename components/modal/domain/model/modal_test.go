package model

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/modal/domain/value"
)

func TestNewModal(t *testing.T) {
	modal := NewModal("Test content")

	if modal.Content() != "Test content" {
		t.Errorf("Expected content 'Test content', got '%s'", modal.Content())
	}
	if modal.Title() != "" {
		t.Errorf("Expected empty title, got '%s'", modal.Title())
	}
	if modal.IsVisible() {
		t.Error("Expected modal to be not visible by default")
	}
	if modal.DimBackground() {
		t.Error("Expected background dimming to be disabled by default")
	}
	if len(modal.Buttons()) != 0 {
		t.Errorf("Expected no buttons, got %d", len(modal.Buttons()))
	}
	if !modal.Position().IsCenter() {
		t.Error("Expected position to be centered by default")
	}
	if modal.Size().Width() != 40 || modal.Size().Height() != 10 {
		t.Errorf("Expected default size 40x10, got %dx%d", modal.Size().Width(), modal.Size().Height())
	}
}

func TestNewModalWithTitle(t *testing.T) {
	modal := NewModalWithTitle("Test Title", "Test content")

	if modal.Title() != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", modal.Title())
	}
	if modal.Content() != "Test content" {
		t.Errorf("Expected content 'Test content', got '%s'", modal.Content())
	}
}

func TestModalWithTitle(t *testing.T) {
	original := NewModal("Content")
	modified := original.WithTitle("New Title")

	// Original unchanged.
	if original.Title() != "" {
		t.Errorf("Original title should be empty, got '%s'", original.Title())
	}

	// Modified has new title.
	if modified.Title() != "New Title" {
		t.Errorf("Expected title 'New Title', got '%s'", modified.Title())
	}
	// Other fields preserved.
	if modified.Content() != "Content" {
		t.Error("Content should be preserved")
	}
}

func TestModalWithContent(t *testing.T) {
	original := NewModal("Original")
	modified := original.WithContent("Modified")

	// Original unchanged.
	if original.Content() != "Original" {
		t.Error("Original content should be unchanged")
	}

	// Modified has new content.
	if modified.Content() != "Modified" {
		t.Errorf("Expected content 'Modified', got '%s'", modified.Content())
	}
}

func TestModalWithSize(t *testing.T) {
	original := NewModal("Content")
	modified := original.WithSize(80, 24)

	// Original unchanged.
	if original.Size().Width() != 40 || original.Size().Height() != 10 {
		t.Error("Original size should be unchanged")
	}

	// Modified has new size.
	if modified.Size().Width() != 80 || modified.Size().Height() != 24 {
		t.Errorf("Expected size 80x24, got %dx%d", modified.Size().Width(), modified.Size().Height())
	}
}

func TestModalWithPosition(t *testing.T) {
	original := NewModal("Content")
	customPos := value.NewPositionCustom(10, 20)
	modified := original.WithPosition(customPos)

	// Original unchanged (centered)
	if !original.Position().IsCenter() {
		t.Error("Original position should be centered")
	}

	// Modified has custom position.
	if modified.Position().IsCenter() {
		t.Error("Modified position should not be centered")
	}
	if modified.Position().X() != 10 || modified.Position().Y() != 20 {
		t.Errorf("Expected position (10, 20), got (%d, %d)", modified.Position().X(), modified.Position().Y())
	}
}

func TestModalWithButtons(t *testing.T) {
	original := NewModal("Content")
	buttons := []*Button{
		NewButton("Yes", "y", "confirm"),
		NewButton("No", "n", "cancel"),
	}
	modified := original.WithButtons(buttons)

	// Original unchanged.
	if len(original.Buttons()) != 0 {
		t.Error("Original should have no buttons")
	}

	// Modified has buttons.
	if len(modified.Buttons()) != 2 {
		t.Errorf("Expected 2 buttons, got %d", len(modified.Buttons()))
	}
	if modified.Buttons()[0].Label() != "Yes" {
		t.Error("First button should be Yes")
	}
	if modified.Buttons()[1].Label() != "No" {
		t.Error("Second button should be No")
	}

	// Focused button reset to 0.
	if modified.FocusedButtonIndex() != 0 {
		t.Errorf("Focused button should be 0, got %d", modified.FocusedButtonIndex())
	}
}

func TestModalWithDimBackground(t *testing.T) {
	original := NewModal("Content")
	modified := original.WithDimBackground(true)

	// Original unchanged.
	if original.DimBackground() {
		t.Error("Original should not have dimming")
	}

	// Modified has dimming.
	if !modified.DimBackground() {
		t.Error("Modified should have dimming enabled")
	}
}

func TestModalWithVisible(t *testing.T) {
	original := NewModal("Content")
	modified := original.WithVisible(true)

	// Original unchanged.
	if original.IsVisible() {
		t.Error("Original should not be visible")
	}

	// Modified is visible.
	if !modified.IsVisible() {
		t.Error("Modified should be visible")
	}
}

func TestModalFocusNextButton(t *testing.T) {
	buttons := []*Button{
		NewButton("First", "1", "first"),
		NewButton("Second", "2", "second"),
		NewButton("Third", "3", "third"),
	}
	modal := NewModal("Content").WithButtons(buttons)

	// Initially focused on first button (index 0)
	if modal.FocusedButtonIndex() != 0 {
		t.Errorf("Expected initial focus 0, got %d", modal.FocusedButtonIndex())
	}

	// Focus next (index 1)
	modal = modal.FocusNextButton()
	if modal.FocusedButtonIndex() != 1 {
		t.Errorf("Expected focus 1, got %d", modal.FocusedButtonIndex())
	}

	// Focus next (index 2)
	modal = modal.FocusNextButton()
	if modal.FocusedButtonIndex() != 2 {
		t.Errorf("Expected focus 2, got %d", modal.FocusedButtonIndex())
	}

	// Focus next (wraps to index 0)
	modal = modal.FocusNextButton()
	if modal.FocusedButtonIndex() != 0 {
		t.Errorf("Expected focus to wrap to 0, got %d", modal.FocusedButtonIndex())
	}
}

func TestModalFocusPreviousButton(t *testing.T) {
	buttons := []*Button{
		NewButton("First", "1", "first"),
		NewButton("Second", "2", "second"),
		NewButton("Third", "3", "third"),
	}
	modal := NewModal("Content").WithButtons(buttons)

	// Initially focused on first button (index 0)
	if modal.FocusedButtonIndex() != 0 {
		t.Errorf("Expected initial focus 0, got %d", modal.FocusedButtonIndex())
	}

	// Focus previous (wraps to index 2)
	modal = modal.FocusPreviousButton()
	if modal.FocusedButtonIndex() != 2 {
		t.Errorf("Expected focus to wrap to 2, got %d", modal.FocusedButtonIndex())
	}

	// Focus previous (index 1)
	modal = modal.FocusPreviousButton()
	if modal.FocusedButtonIndex() != 1 {
		t.Errorf("Expected focus 1, got %d", modal.FocusedButtonIndex())
	}

	// Focus previous (index 0)
	modal = modal.FocusPreviousButton()
	if modal.FocusedButtonIndex() != 0 {
		t.Errorf("Expected focus 0, got %d", modal.FocusedButtonIndex())
	}
}

func TestModalFocusButtonsNoButtons(t *testing.T) {
	modal := NewModal("Content")

	// FocusNextButton with no buttons should return unchanged modal.
	next := modal.FocusNextButton()
	if next.FocusedButtonIndex() != modal.FocusedButtonIndex() {
		t.Error("Focus should not change when there are no buttons")
	}

	// FocusPreviousButton with no buttons should return unchanged modal.
	prev := modal.FocusPreviousButton()
	if prev.FocusedButtonIndex() != modal.FocusedButtonIndex() {
		t.Error("Focus should not change when there are no buttons")
	}
}

func TestModalFocusedButton(t *testing.T) {
	buttons := []*Button{
		NewButton("Yes", "y", "confirm"),
		NewButton("No", "n", "cancel"),
	}
	modal := NewModal("Content").WithButtons(buttons)

	// Initially focused on first button.
	focused := modal.FocusedButton()
	if focused == nil {
		t.Fatal("Expected focused button, got nil")
	}
	if focused.Label() != "Yes" {
		t.Errorf("Expected focused button 'Yes', got '%s'", focused.Label())
	}

	// Focus next button.
	modal = modal.FocusNextButton()
	focused = modal.FocusedButton()
	if focused == nil {
		t.Fatal("Expected focused button, got nil")
	}
	if focused.Label() != "No" {
		t.Errorf("Expected focused button 'No', got '%s'", focused.Label())
	}
}

func TestModalFocusedButtonNoButtons(t *testing.T) {
	modal := NewModal("Content")

	focused := modal.FocusedButton()
	if focused != nil {
		t.Error("Expected nil focused button when there are no buttons")
	}
}

func TestModalSingleButton(t *testing.T) {
	buttons := []*Button{
		NewButton("OK", "enter", "ok"),
	}
	modal := NewModal("Content").WithButtons(buttons)

	// Focus next with single button (stays at 0)
	next := modal.FocusNextButton()
	if next.FocusedButtonIndex() != 0 {
		t.Error("Single button should stay focused")
	}

	// Focus previous with single button (stays at 0)
	prev := modal.FocusPreviousButton()
	if prev.FocusedButtonIndex() != 0 {
		t.Error("Single button should stay focused")
	}
}

func TestModalImmutability(t *testing.T) {
	original := NewModal("Original Content")

	// Apply multiple changes.
	modified := original.
		WithTitle("Title").
		WithSize(80, 24).
		WithVisible(true).
		WithDimBackground(true)

	// Original should be unchanged.
	if original.Title() != "" {
		t.Error("Original title should be unchanged")
	}
	if original.IsVisible() {
		t.Error("Original visibility should be unchanged")
	}
	if original.DimBackground() {
		t.Error("Original dimming should be unchanged")
	}
	if original.Size().Width() != 40 {
		t.Error("Original size should be unchanged")
	}

	// Modified should have all changes.
	if modified.Title() != "Title" {
		t.Error("Modified title incorrect")
	}
	if !modified.IsVisible() {
		t.Error("Modified visibility incorrect")
	}
	if !modified.DimBackground() {
		t.Error("Modified dimming incorrect")
	}
	if modified.Size().Width() != 80 {
		t.Error("Modified size incorrect")
	}
}
