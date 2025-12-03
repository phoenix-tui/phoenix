package value

import "testing"

func TestNewOption(t *testing.T) {
	opt := NewOption("Test Label", 42)

	if opt.Label() != "Test Label" {
		t.Errorf("expected label 'Test Label', got %q", opt.Label())
	}
	if opt.Value() != 42 {
		t.Errorf("expected value 42, got %d", opt.Value())
	}
	if opt.Description() != "" {
		t.Errorf("expected empty description, got %q", opt.Description())
	}
	if opt.Disabled() {
		t.Error("expected option not disabled")
	}
}

func TestWithDescription(t *testing.T) {
	opt := NewOption("Label", "value")
	opt = opt.WithDescription("This is a description")

	if opt.Description() != "This is a description" {
		t.Errorf("expected description, got %q", opt.Description())
	}
}

func TestWithDisabled(t *testing.T) {
	opt := NewOption("Label", "value")

	opt = opt.WithDisabled(true)
	if !opt.Disabled() {
		t.Error("expected option disabled")
	}

	opt = opt.WithDisabled(false)
	if opt.Disabled() {
		t.Error("expected option not disabled")
	}
}

func TestOptionImmutability(t *testing.T) {
	original := NewOption("Original", 1)
	modified := original.WithDescription("Modified")

	if original.Description() != "" {
		t.Error("original option was mutated")
	}
	if modified.Description() != "Modified" {
		t.Error("modified option does not have expected description")
	}
}

func TestOptionGenericTypes(t *testing.T) {
	t.Run("string option", func(t *testing.T) {
		opt := NewOption("Label", "value")
		if opt.Value() != "value" {
			t.Errorf("expected 'value', got %q", opt.Value())
		}
	})

	t.Run("int option", func(t *testing.T) {
		opt := NewOption("Number", 42)
		if opt.Value() != 42 {
			t.Errorf("expected 42, got %d", opt.Value())
		}
	})

	t.Run("struct option", func(t *testing.T) {
		type CustomType struct {
			ID   int
			Name string
		}
		custom := CustomType{ID: 1, Name: "Test"}
		opt := NewOption("Custom", custom)
		if opt.Value().ID != 1 {
			t.Errorf("expected ID 1, got %d", opt.Value().ID)
		}
		if opt.Value().Name != "Test" {
			t.Errorf("expected name 'Test', got %q", opt.Value().Name)
		}
	})
}
