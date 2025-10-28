package program

import (
	"bytes"
	"strings"
	"testing"
)

// TestWithInput verifies custom input option works.
func TestWithInput(t *testing.T) {
	customInput := strings.NewReader("test input")

	p := New(
		TestModel{},
		WithInput[TestModel](customInput),
	)

	// Verify p.input is custom reader
	if p.input != customInput {
		t.Error("WithInput should set custom reader")
	}
}

// TestWithOutput verifies custom output option works.
func TestWithOutput(t *testing.T) {
	var customOutput bytes.Buffer

	p := New(
		TestModel{},
		WithOutput[TestModel](&customOutput),
	)

	// Verify p.output is custom writer
	if p.output != &customOutput {
		t.Error("WithOutput should set custom writer")
	}
}

// TestWithAltScreen verifies alt screen flag option works.
func TestWithAltScreen(t *testing.T) {
	p := New(
		TestModel{},
		WithAltScreen[TestModel](),
	)

	// Verify p.altScreen = true
	if !p.altScreen {
		t.Error("WithAltScreen should set altScreen flag to true")
	}
}

// TestWithMouseAllMotion verifies mouse flag option works.
func TestWithMouseAllMotion(t *testing.T) {
	p := New(
		TestModel{},
		WithMouseAllMotion[TestModel](),
	)

	// Verify p.mouseAllMotion = true
	if !p.mouseAllMotion {
		t.Error("WithMouseAllMotion should set mouseAllMotion flag to true")
	}
}

// TestOptions_Combination verifies multiple options work together.
func TestOptions_Combination(t *testing.T) {
	customInput := strings.NewReader("test input")
	var customOutput bytes.Buffer

	p := New(
		TestModel{},
		WithInput[TestModel](customInput),
		WithOutput[TestModel](&customOutput),
		WithAltScreen[TestModel](),
		WithMouseAllMotion[TestModel](),
	)

	// Verify all flags/fields set correctly
	if p.input != customInput {
		t.Error("input should be custom reader")
	}
	if p.output != &customOutput {
		t.Error("output should be custom writer")
	}
	if !p.altScreen {
		t.Error("altScreen should be true")
	}
	if !p.mouseAllMotion {
		t.Error("mouseAllMotion should be true")
	}
}

// TestOptions_Order verifies options order independence.
func TestOptions_Order(t *testing.T) {
	customInput := strings.NewReader("test input")
	var customOutput bytes.Buffer

	// Order 1: Input then Output
	p1 := New(
		TestModel{},
		WithInput[TestModel](customInput),
		WithOutput[TestModel](&customOutput),
	)

	if p1.input != customInput {
		t.Error("input should be set (order 1)")
	}
	if p1.output != &customOutput {
		t.Error("output should be set (order 1)")
	}

	// Order 2: Output then Input
	customInput2 := strings.NewReader("test input 2")
	var customOutput2 bytes.Buffer

	p2 := New(
		TestModel{},
		WithOutput[TestModel](&customOutput2),
		WithInput[TestModel](customInput2),
	)

	if p2.input != customInput2 {
		t.Error("input should be set (order 2)")
	}
	if p2.output != &customOutput2 {
		t.Error("output should be set (order 2)")
	}

	// Both should work (order doesn't matter)
}

// TestOptions_Default verifies default values without options.
func TestOptions_Default(t *testing.T) {
	p := New(TestModel{})

	// Default flags should be false
	if p.altScreen {
		t.Error("altScreen should default to false")
	}
	if p.mouseAllMotion {
		t.Error("mouseAllMotion should default to false")
	}
}

// TestOptions_Overwrite verifies later options overwrite earlier ones.
func TestOptions_Overwrite(t *testing.T) {
	input1 := strings.NewReader("input 1")
	input2 := strings.NewReader("input 2")

	p := New(
		TestModel{},
		WithInput[TestModel](input1),
		WithInput[TestModel](input2), // Should overwrite
	)

	// Later option should win
	if p.input != input2 {
		t.Error("later WithInput should overwrite earlier one")
	}
}
