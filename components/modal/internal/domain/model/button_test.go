package model

import (
	"testing"
)

func TestNewButton(t *testing.T) {
	button := NewButton("Yes", "y", "confirm")

	if button.Label() != "Yes" {
		t.Errorf("Expected label 'Yes', got '%s'", button.Label())
	}
	if button.Key() != "y" {
		t.Errorf("Expected key 'y', got '%s'", button.Key())
	}
	if button.Action() != "confirm" {
		t.Errorf("Expected action 'confirm', got '%s'", button.Action())
	}
}

func TestButtonEmptyValues(t *testing.T) {
	button := NewButton("", "", "")

	if button.Label() != "" {
		t.Errorf("Expected empty label, got '%s'", button.Label())
	}
	if button.Key() != "" {
		t.Errorf("Expected empty key, got '%s'", button.Key())
	}
	if button.Action() != "" {
		t.Errorf("Expected empty action, got '%s'", button.Action())
	}
}

func TestButtonMultipleButtons(t *testing.T) {
	yes := NewButton("Yes", "y", "confirm")
	no := NewButton("No", "n", "cancel")
	ok := NewButton("OK", "enter", "ok")

	// Verify each button maintains its own state.
	if yes.Action() != "confirm" {
		t.Error("Yes button action incorrect")
	}
	if no.Action() != "cancel" {
		t.Error("No button action incorrect")
	}
	if ok.Action() != "ok" {
		t.Error("OK button action incorrect")
	}
}
