package value

import "testing"

func TestNewOption(t *testing.T) {
	opt := NewOption("Test", 42)

	if opt.Label() != "Test" {
		t.Errorf("Label() = %q, want %q", opt.Label(), "Test")
	}
	if opt.Value() != 42 {
		t.Errorf("Value() = %d, want %d", opt.Value(), 42)
	}
	if opt.Description() != "" {
		t.Errorf("Description() = %q, want empty", opt.Description())
	}
	if opt.Disabled() {
		t.Error("Disabled() = true, want false")
	}
}

func TestOption_WithDescription(t *testing.T) {
	opt := NewOption("Test", 42).WithDescription("A test option")

	if opt.Description() != "A test option" {
		t.Errorf("Description() = %q, want %q", opt.Description(), "A test option")
	}
	// Original fields preserved
	if opt.Label() != "Test" {
		t.Errorf("Label() = %q, want %q", opt.Label(), "Test")
	}
	if opt.Value() != 42 {
		t.Errorf("Value() = %d, want %d", opt.Value(), 42)
	}
}

func TestOption_WithDisabled(t *testing.T) {
	opt := NewOption("Test", 42).WithDisabled(true)

	if !opt.Disabled() {
		t.Error("Disabled() = false, want true")
	}

	// Can be toggled back
	opt = opt.WithDisabled(false)
	if opt.Disabled() {
		t.Error("Disabled() = true, want false")
	}
}

func TestOption_Immutability(t *testing.T) {
	original := NewOption("Original", 1)
	modified := original.WithDescription("Modified").WithDisabled(true)

	// Original should be unchanged
	if original.Description() != "" {
		t.Errorf("original.Description() = %q, want empty", original.Description())
	}
	if original.Disabled() {
		t.Error("original.Disabled() = true, want false")
	}

	// Modified should have new values
	if modified.Description() != "Modified" {
		t.Errorf("modified.Description() = %q, want %q", modified.Description(), "Modified")
	}
	if !modified.Disabled() {
		t.Error("modified.Disabled() = false, want true")
	}
}

func TestOption_GenericTypes(t *testing.T) {
	// String type
	strOpt := NewOption("String", "value")
	if strOpt.Value() != "value" {
		t.Errorf("strOpt.Value() = %q, want %q", strOpt.Value(), "value")
	}

	// Struct type
	type Custom struct {
		ID   int
		Name string
	}
	customOpt := NewOption("Custom", Custom{ID: 1, Name: "Test"})
	val := customOpt.Value()
	if val.ID != 1 || val.Name != "Test" {
		t.Errorf("customOpt.Value() = %+v, want {ID:1, Name:Test}", val)
	}

	// Pointer type
	ptrOpt := NewOption("Pointer", &Custom{ID: 2, Name: "Ptr"})
	ptrVal := ptrOpt.Value()
	if ptrVal.ID != 2 || ptrVal.Name != "Ptr" {
		t.Errorf("ptrOpt.Value() = %+v, want {ID:2, Name:Ptr}", *ptrVal)
	}
}
