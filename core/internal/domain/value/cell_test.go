package value_test

import (
	"testing"

	"github.com/phoenix-tui/phoenix/core/internal/domain/value"
)

func TestNewCell(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		width           int
		expectedContent string
		expectedWidth   int
	}{
		{
			name:            "ASCII character",
			content:         "a",
			width:           1,
			expectedContent: "a",
			expectedWidth:   1,
		},
		{
			name:            "emoji",
			content:         "👋",
			width:           2,
			expectedContent: "👋",
			expectedWidth:   2,
		},
		{
			name:            "emoji with modifier",
			content:         "👋🏻",
			width:           2,
			expectedContent: "👋🏻",
			expectedWidth:   2,
		},
		{
			name:            "zero width",
			content:         "\u200b", // Zero-width space
			width:           0,
			expectedContent: "\u200b",
			expectedWidth:   0,
		},
		{
			name:            "negative width clamped to 0",
			content:         "a",
			width:           -1,
			expectedContent: "a",
			expectedWidth:   0,
		},
		{
			name:            "empty string",
			content:         "",
			width:           0,
			expectedContent: "",
			expectedWidth:   0,
		},
		{
			name:            "space",
			content:         " ",
			width:           1,
			expectedContent: " ",
			expectedWidth:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := value.NewCell(tt.content, tt.width)

			if cell.Content() != tt.expectedContent {
				t.Errorf("expected content %q, got %q", tt.expectedContent, cell.Content())
			}
			if cell.Width() != tt.expectedWidth {
				t.Errorf("expected width %d, got %d", tt.expectedWidth, cell.Width())
			}
		})
	}
}

func TestCell_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		cell     value.Cell
		expected bool
	}{
		{
			name:     "empty string",
			cell:     value.NewCell("", 0),
			expected: true,
		},
		{
			name:     "space",
			cell:     value.NewCell(" ", 1),
			expected: true,
		},
		{
			name:     "zero width",
			cell:     value.NewCell("x", 0),
			expected: true,
		},
		{
			name:     "ASCII character",
			cell:     value.NewCell("a", 1),
			expected: false,
		},
		{
			name:     "emoji",
			cell:     value.NewCell("👋", 2),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cell.IsEmpty()

			if result != tt.expected {
				t.Errorf("expected %v, got %v for cell %q (width %d)",
					tt.expected, result, tt.cell.Content(), tt.cell.Width())
			}
		})
	}
}

func TestCell_Equal(t *testing.T) {
	tests := []struct {
		name     string
		cell1    value.Cell
		cell2    value.Cell
		expected bool
	}{
		{
			name:     "equal cells",
			cell1:    value.NewCell("a", 1),
			cell2:    value.NewCell("a", 1),
			expected: true,
		},
		{
			name:     "different content",
			cell1:    value.NewCell("a", 1),
			cell2:    value.NewCell("b", 1),
			expected: false,
		},
		{
			name:     "different width",
			cell1:    value.NewCell("a", 1),
			cell2:    value.NewCell("a", 2),
			expected: false,
		},
		{
			name:     "both different",
			cell1:    value.NewCell("a", 1),
			cell2:    value.NewCell("b", 2),
			expected: false,
		},
		{
			name:     "equal emoji",
			cell1:    value.NewCell("👋", 2),
			cell2:    value.NewCell("👋", 2),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cell1.Equal(tt.cell2)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCell_Immutability(t *testing.T) {
	// Verify that Cell is truly immutable
	cell := value.NewCell("a", 1)
	originalContent := cell.Content()
	originalWidth := cell.Width()

	// Create another cell - original should be unchanged
	_ = value.NewCell("b", 2)

	if cell.Content() != originalContent {
		t.Error("cell content was mutated")
	}
	if cell.Width() != originalWidth {
		t.Error("cell width was mutated")
	}
}
