package layout

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/layout/domain/value"
)

func TestRow(t *testing.T) {
	flex := Row()

	if flex == nil {
		t.Fatal("Row() returned nil")
	}

	if flex.domain == nil {
		t.Fatal("domain is nil")
	}

	if !flex.domain.IsHorizontal() {
		t.Error("Row() should create horizontal flex")
	}
}

func TestColumn(t *testing.T) {
	flex := Column()

	if flex == nil {
		t.Fatal("Column() returned nil")
	}

	if flex.domain == nil {
		t.Fatal("domain is nil")
	}

	if !flex.domain.IsVertical() {
		t.Error("Column() should create vertical flex")
	}
}

func TestFlex_Add(t *testing.T) {
	flex := Row().
		Add(NewBox("Item 1")).
		Add(NewBox("Item 2"))

	if flex.domain.ItemCount() != 2 {
		t.Errorf("ItemCount() = %d, want 2", flex.domain.ItemCount())
	}
}

func TestFlex_AddRaw(t *testing.T) {
	flex := Row().
		AddRaw("Item 1").
		AddRaw("Item 2")

	if flex.domain.ItemCount() != 2 {
		t.Errorf("ItemCount() = %d, want 2", flex.domain.ItemCount())
	}
}

func TestFlex_Gap(t *testing.T) {
	flex := Row().Gap(5)

	if flex.domain.Gap() != 5 {
		t.Errorf("Gap() = %d, want 5", flex.domain.Gap())
	}
}

func TestFlex_JustifyMethods(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Flex
		expected value.JustifyContent
	}{
		{
			name:     "JustifyStart",
			setup:    func() *Flex { return Row().JustifyStart() },
			expected: value.JustifyContentStart,
		},
		{
			name:     "JustifyEnd",
			setup:    func() *Flex { return Row().JustifyEnd() },
			expected: value.JustifyContentEnd,
		},
		{
			name:     "JustifyCenter",
			setup:    func() *Flex { return Row().JustifyCenter() },
			expected: value.JustifyContentCenter,
		},
		{
			name:     "JustifySpaceBetween",
			setup:    func() *Flex { return Row().JustifySpaceBetween() },
			expected: value.JustifyContentSpaceBetween,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flex := tt.setup()
			if flex.domain.JustifyContent() != tt.expected {
				t.Errorf("JustifyContent() = %v, want %v", flex.domain.JustifyContent(), tt.expected)
			}
		})
	}
}

func TestFlex_AlignMethods(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Flex
		expected value.AlignItems
	}{
		{
			name:     "AlignStretch",
			setup:    func() *Flex { return Row().AlignStretch() },
			expected: value.AlignItemsStretch,
		},
		{
			name:     "AlignStart",
			setup:    func() *Flex { return Row().AlignStart() },
			expected: value.AlignItemsStart,
		},
		{
			name:     "AlignEnd",
			setup:    func() *Flex { return Row().AlignEnd() },
			expected: value.AlignItemsEnd,
		},
		{
			name:     "AlignCenter",
			setup:    func() *Flex { return Row().AlignCenter() },
			expected: value.AlignItemsCenter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flex := tt.setup()
			if flex.domain.AlignItems() != tt.expected {
				t.Errorf("AlignItems() = %v, want %v", flex.domain.AlignItems(), tt.expected)
			}
		})
	}
}

func TestFlex_SizeMethods(t *testing.T) {
	t.Run("Width", func(t *testing.T) {
		flex := Row().Width(80)
		size := flex.domain.Size()

		if !size.HasWidth() || size.Width() != 80 {
			t.Errorf("Width not set correctly")
		}
	})

	t.Run("Height", func(t *testing.T) {
		flex := Column().Height(24)
		size := flex.domain.Size()

		if !size.HasHeight() || size.Height() != 24 {
			t.Errorf("Height not set correctly")
		}
	})
}

func TestFlex_Render_HorizontalLayout(t *testing.T) {
	flex := Row().
		AddRaw("A").
		AddRaw("B").
		AddRaw("C")

	output := flex.Render(10, 3)

	// Check output is not empty
	if output == "" {
		t.Error("Render() returned empty string")
	}

	// Check output has correct dimensions
	lines := strings.Split(output, "\n")
	if len(lines) != 3 {
		t.Errorf("Render() returned %d lines, want 3", len(lines))
	}

	// Check first line contains all items
	firstLine := lines[0]
	if !strings.Contains(firstLine, "A") {
		t.Error("First line should contain 'A'")
	}
	if !strings.Contains(firstLine, "B") {
		t.Error("First line should contain 'B'")
	}
	if !strings.Contains(firstLine, "C") {
		t.Error("First line should contain 'C'")
	}
}

func TestFlex_Render_VerticalLayout(t *testing.T) {
	flex := Column().
		AddRaw("A").
		AddRaw("B").
		AddRaw("C")

	output := flex.Render(10, 5)

	// Check output is not empty
	if output == "" {
		t.Error("Render() returned empty string")
	}

	// Check items appear on different lines
	lines := strings.Split(output, "\n")

	foundA := false
	foundB := false
	foundC := false

	for _, line := range lines {
		if strings.Contains(line, "A") {
			foundA = true
		}
		if strings.Contains(line, "B") {
			foundB = true
		}
		if strings.Contains(line, "C") {
			foundC = true
		}
	}

	if !foundA || !foundB || !foundC {
		t.Error("All items should appear in vertical layout")
	}
}

func TestFlex_Render_WithGap(t *testing.T) {
	flex := Row().
		Gap(3).
		AddRaw("A").
		AddRaw("B")

	output := flex.Render(20, 3)

	// Check items are separated
	lines := strings.Split(output, "\n")
	firstLine := lines[0]

	indexA := strings.Index(firstLine, "A")
	indexB := strings.Index(firstLine, "B")

	if indexA == -1 || indexB == -1 {
		t.Error("Both items should appear in output")
	}

	// Check gap (B should be at least 4 positions after A: 1 for A + 3 gap)
	if indexB-indexA < 4 {
		t.Errorf("Gap not applied correctly: A at %d, B at %d", indexA, indexB)
	}
}

func TestFlex_String(t *testing.T) {
	flex := Row().AddRaw("Test")

	str := flex.String()

	if str == "" {
		t.Error("String() returned empty string")
	}

	if !strings.Contains(str, "Test") {
		t.Error("String() should contain 'Test'")
	}
}

func TestFlex_Domain(t *testing.T) {
	flex := Row()

	domain := flex.Domain()

	if domain == nil {
		t.Error("Domain() returned nil")
	}

	if domain != flex.domain {
		t.Error("Domain() should return the internal domain model")
	}
}

func TestFlex_Chaining(t *testing.T) {
	// Test that all methods return *Flex for chaining
	flex := Row().
		Gap(2).
		JustifyCenter().
		AlignStart().
		Width(80).
		Height(24).
		AddRaw("Item 1").
		AddRaw("Item 2")

	if flex == nil {
		t.Error("Chaining returned nil")
		return // Prevent nil dereference
	}

	if flex.domain.ItemCount() != 2 {
		t.Error("Chaining didn't preserve items")
	}
}

// Example tests for documentation
func ExampleRow() {
	flex := Row().
		AddRaw("Left").
		AddRaw("Right")

	output := flex.Render(20, 3)
	_ = output // Use output
}

func ExampleColumn() {
	flex := Column().
		AddRaw("Top").
		AddRaw("Middle").
		AddRaw("Bottom")

	output := flex.Render(20, 10)
	_ = output // Use output
}

func ExampleFlex_Gap() {
	flex := Row().
		Gap(5).
		AddRaw("A").
		AddRaw("B").
		AddRaw("C")

	output := flex.Render(30, 3)
	_ = output // Use output
}
