package model

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/layout/domain/value"
)

// TestBox_Creation tests box creation and default values
func TestBox_Creation(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		shouldPanic bool
	}{
		{
			name:        "valid content",
			content:     "Hello",
			shouldPanic: false,
		},
		{
			name:        "empty content panics",
			content:     "",
			shouldPanic: true,
		},
		{
			name:        "multiline content",
			content:     "Line 1\nLine 2",
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.shouldPanic && r == nil {
					t.Error("Expected panic but got none")
				}
				if !tt.shouldPanic && r != nil {
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			box := NewBox(tt.content)

			if tt.shouldPanic {
				return // Test passed if we panicked
			}

			// Check default values
			if box.Content() != tt.content {
				t.Errorf("Expected content %q, got %q", tt.content, box.Content())
			}

			if !box.Padding().IsZero() {
				t.Errorf("Expected zero padding, got %s", box.Padding())
			}

			if !box.Margin().IsZero() {
				t.Errorf("Expected zero margin, got %s", box.Margin())
			}

			if box.HasBorder() {
				t.Error("Expected no border by default")
			}

			if !box.Size().IsUnconstrained() {
				t.Errorf("Expected unconstrained size, got %s", box.Size())
			}

			if !box.Alignment().IsDefault() {
				t.Errorf("Expected default alignment, got %s", box.Alignment())
			}
		})
	}
}

// TestBox_WithContent tests content modification
func TestBox_WithContent(t *testing.T) {
	original := NewBox("Original")

	// Change content
	modified := original.WithContent("Modified")

	// Verify immutability
	if original.Content() != "Original" {
		t.Error("Original box was mutated")
	}

	if modified.Content() != "Modified" {
		t.Errorf("Expected 'Modified', got %q", modified.Content())
	}

	// Test empty content panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for empty content")
		}
	}()
	original.WithContent("")
}

// TestBox_WithPadding tests padding modification
func TestBox_WithPadding(t *testing.T) {
	original := NewBox("Test")
	padding := value.NewSpacingAll(2)

	modified := original.WithPadding(padding)

	// Verify immutability
	if !original.Padding().IsZero() {
		t.Error("Original box was mutated")
	}

	if !modified.Padding().Equals(padding) {
		t.Errorf("Expected padding %s, got %s", padding, modified.Padding())
	}
}

// TestBox_WithMargin tests margin modification
func TestBox_WithMargin(t *testing.T) {
	original := NewBox("Test")
	margin := value.NewSpacingVH(1, 2)

	modified := original.WithMargin(margin)

	// Verify immutability
	if !original.Margin().IsZero() {
		t.Error("Original box was mutated")
	}

	if !modified.Margin().Equals(margin) {
		t.Errorf("Expected margin %s, got %s", margin, modified.Margin())
	}
}

// TestBox_WithBorder tests border enable/disable
func TestBox_WithBorder(t *testing.T) {
	original := NewBox("Test")

	// Enable border
	withBorder := original.WithBorder(true)

	// Verify immutability
	if original.HasBorder() {
		t.Error("Original box was mutated")
	}

	if !withBorder.HasBorder() {
		t.Error("Expected border to be enabled")
	}

	// Disable border
	withoutBorder := withBorder.WithBorder(false)

	if withBorder.HasBorder() == false {
		t.Error("Modified box was mutated")
	}

	if withoutBorder.HasBorder() {
		t.Error("Expected border to be disabled")
	}
}

// TestBox_WithSize tests size constraints
func TestBox_WithSize(t *testing.T) {
	original := NewBox("Test")
	size := value.NewSizeExact(80, 24)

	modified := original.WithSize(size)

	// Verify immutability
	if !original.Size().IsUnconstrained() {
		t.Error("Original box was mutated")
	}

	if modified.Size().Width() != 80 || modified.Size().Height() != 24 {
		t.Errorf("Expected size 80x24, got %s", modified.Size())
	}
}

// TestBox_WithAlignment tests alignment modification
func TestBox_WithAlignment(t *testing.T) {
	original := NewBox("Test")
	alignment := value.NewAlignmentCenter()

	modified := original.WithAlignment(alignment)

	// Verify immutability
	if !original.Alignment().IsDefault() {
		t.Error("Original box was mutated")
	}

	if !modified.Alignment().IsCenter() {
		t.Errorf("Expected center alignment, got %s", modified.Alignment())
	}
}

// TestBox_ContentSize tests content size calculation
func TestBox_ContentSize(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "single line",
			content:        "Hello",
			expectedWidth:  5,
			expectedHeight: 1,
		},
		{
			name:           "multiple lines same width",
			content:        "ABC\nDEF\nGHI",
			expectedWidth:  3,
			expectedHeight: 3,
		},
		{
			name:           "multiple lines different widths",
			content:        "Short\nMedium line\nX",
			expectedWidth:  11, // "Medium line" is longest
			expectedHeight: 3,
		},
		{
			name:           "empty lines",
			content:        "Top\n\nBottom",
			expectedWidth:  6, // "Bottom" is longest
			expectedHeight: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := NewBox(tt.content)
			size := box.ContentSize()

			if size.Width() != tt.expectedWidth {
				t.Errorf("Expected width %d, got %d", tt.expectedWidth, size.Width())
			}

			if size.Height() != tt.expectedHeight {
				t.Errorf("Expected height %d, got %d", tt.expectedHeight, size.Height())
			}
		})
	}
}

// TestBox_PaddedSize tests padded size calculation
func TestBox_PaddedSize(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		padding        value.Spacing
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "no padding",
			content:        "Hello",
			padding:        value.NewSpacingZero(),
			expectedWidth:  5,
			expectedHeight: 1,
		},
		{
			name:           "uniform padding",
			content:        "Hi",
			padding:        value.NewSpacingAll(1),
			expectedWidth:  4, // 2 + 1 left + 1 right
			expectedHeight: 3, // 1 + 1 top + 1 bottom
		},
		{
			name:           "asymmetric padding",
			content:        "Test",
			padding:        value.NewSpacing(1, 2, 3, 4),
			expectedWidth:  10, // 4 + 4 left + 2 right
			expectedHeight: 5,  // 1 + 1 top + 3 bottom
		},
		{
			name:           "multiline with padding",
			content:        "A\nB\nC",
			padding:        value.NewSpacingVH(2, 3),
			expectedWidth:  7, // 1 + 3*2
			expectedHeight: 7, // 3 + 2*2
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := NewBox(tt.content).WithPadding(tt.padding)
			size := box.PaddedSize()

			if size.Width() != tt.expectedWidth {
				t.Errorf("Expected width %d, got %d", tt.expectedWidth, size.Width())
			}

			if size.Height() != tt.expectedHeight {
				t.Errorf("Expected height %d, got %d", tt.expectedHeight, size.Height())
			}
		})
	}
}

// TestBox_BorderedSize tests bordered size calculation
func TestBox_BorderedSize(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		padding        value.Spacing
		hasBorder      bool
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "no border",
			content:        "Hello",
			padding:        value.NewSpacingZero(),
			hasBorder:      false,
			expectedWidth:  5,
			expectedHeight: 1,
		},
		{
			name:           "with border no padding",
			content:        "Hello",
			padding:        value.NewSpacingZero(),
			hasBorder:      true,
			expectedWidth:  7, // 5 + 2
			expectedHeight: 3, // 1 + 2
		},
		{
			name:           "with border and padding",
			content:        "Hi",
			padding:        value.NewSpacingAll(1),
			hasBorder:      true,
			expectedWidth:  6, // 2 + 2 (padding) + 2 (border)
			expectedHeight: 5, // 1 + 2 (padding) + 2 (border)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := NewBox(tt.content).
				WithPadding(tt.padding).
				WithBorder(tt.hasBorder)
			size := box.BorderedSize()

			if size.Width() != tt.expectedWidth {
				t.Errorf("Expected width %d, got %d", tt.expectedWidth, size.Width())
			}

			if size.Height() != tt.expectedHeight {
				t.Errorf("Expected height %d, got %d", tt.expectedHeight, size.Height())
			}
		})
	}
}

// TestBox_TotalSize tests total size calculation (full box model)
func TestBox_TotalSize(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		padding        value.Spacing
		margin         value.Spacing
		hasBorder      bool
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "minimal box",
			content:        "X",
			padding:        value.NewSpacingZero(),
			margin:         value.NewSpacingZero(),
			hasBorder:      false,
			expectedWidth:  1,
			expectedHeight: 1,
		},
		{
			name:           "full box model",
			content:        "Hi",
			padding:        value.NewSpacingAll(1),
			margin:         value.NewSpacingAll(2),
			hasBorder:      true,
			expectedWidth:  10, // 2 + 2 (padding) + 2 (border) + 4 (margin)
			expectedHeight: 9,  // 1 + 2 (padding) + 2 (border) + 4 (margin)
		},
		{
			name:           "asymmetric spacing",
			content:        "Test",
			padding:        value.NewSpacing(1, 2, 1, 2),
			margin:         value.NewSpacing(2, 1, 2, 1),
			hasBorder:      true,
			expectedWidth:  12, // 4 + 4 (padding L/R: 2+2) + 2 (border) + 2 (margin L/R: 1+1)
			expectedHeight: 9,  // 1 + 2 (padding T/B: 1+1) + 2 (border) + 4 (margin T/B: 2+2)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := NewBox(tt.content).
				WithPadding(tt.padding).
				WithMargin(tt.margin).
				WithBorder(tt.hasBorder)
			size := box.TotalSize()

			if size.Width() != tt.expectedWidth {
				t.Errorf("Expected width %d, got %d", tt.expectedWidth, size.Width())
			}

			if size.Height() != tt.expectedHeight {
				t.Errorf("Expected height %d, got %d", tt.expectedHeight, size.Height())
			}
		})
	}
}

// TestBox_SizeProgression tests that size methods build on each other correctly
func TestBox_SizeProgression(t *testing.T) {
	box := NewBox("Test").
		WithPadding(value.NewSpacingAll(1)).
		WithBorder(true).
		WithMargin(value.NewSpacingAll(1))

	content := box.ContentSize()
	padded := box.PaddedSize()
	bordered := box.BorderedSize()
	total := box.TotalSize()

	// Each size should be larger than the previous
	if padded.Width() <= content.Width() {
		t.Error("Padded width should be larger than content width")
	}
	if bordered.Width() <= padded.Width() {
		t.Error("Bordered width should be larger than padded width")
	}
	if total.Width() <= bordered.Width() {
		t.Error("Total width should be larger than bordered width")
	}

	// Verify exact progression
	expectedPaddedW := content.Width() + 2   // +1 left, +1 right
	expectedBorderedW := expectedPaddedW + 2 // +1 left, +1 right
	expectedTotalW := expectedBorderedW + 2  // +1 left, +1 right

	if padded.Width() != expectedPaddedW {
		t.Errorf("Expected padded width %d, got %d", expectedPaddedW, padded.Width())
	}
	if bordered.Width() != expectedBorderedW {
		t.Errorf("Expected bordered width %d, got %d", expectedBorderedW, bordered.Width())
	}
	if total.Width() != expectedTotalW {
		t.Errorf("Expected total width %d, got %d", expectedTotalW, total.Width())
	}
}

// TestBox_FluentAPI tests method chaining
func TestBox_FluentAPI(t *testing.T) {
	box := NewBox("Test").
		WithPadding(value.NewSpacingAll(1)).
		WithMargin(value.NewSpacingVH(1, 2)).
		WithBorder(true).
		WithSize(value.NewSizeExact(80, 24)).
		WithAlignment(value.NewAlignmentCenter())

	// Verify all properties were set
	if box.Content() != "Test" {
		t.Error("Content not preserved in chain")
	}
	if box.Padding().Top() != 1 {
		t.Error("Padding not set in chain")
	}
	if box.Margin().Top() != 1 {
		t.Error("Margin not set in chain")
	}
	if !box.HasBorder() {
		t.Error("Border not set in chain")
	}
	if box.Size().Width() != 80 {
		t.Error("Size not set in chain")
	}
	if !box.Alignment().IsCenter() {
		t.Error("Alignment not set in chain")
	}
}

// TestBox_String tests debug representation
func TestBox_String(t *testing.T) {
	tests := []struct {
		name     string
		box      *Box
		contains []string
	}{
		{
			name: "minimal box",
			box:  NewBox("Hello"),
			contains: []string{
				"Box{",
				"content=\"Hello\"",
				"total=5x1",
			},
		},
		{
			name: "full box",
			box: NewBox("Test").
				WithPadding(value.NewSpacingAll(1)).
				WithMargin(value.NewSpacingAll(1)).
				WithBorder(true),
			contains: []string{
				"Box{",
				"content=\"Test\"",
				"padding=",
				"margin=",
				"border=true",
				"total=",
			},
		},
		{
			name: "long content truncated",
			box:  NewBox("This is a very long content string that should be truncated in the output"),
			contains: []string{
				"Box{",
				"content=\"This is a very long content...",
			},
		},
		{
			name: "multiline content",
			box:  NewBox("Line1\nLine2"),
			contains: []string{
				"Box{",
				"content=",
				"Line1\\\\nLine2", // Double-escaped: \\ becomes \ in output, then \n
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := tt.box.String()

			for _, substr := range tt.contains {
				if !strings.Contains(str, substr) {
					t.Errorf("Expected string to contain %q, got: %s", substr, str)
				}
			}
		})
	}
}

// TestBox_Immutability tests that all operations return new instances
func TestBox_Immutability(t *testing.T) {
	original := NewBox("Original").
		WithPadding(value.NewSpacingAll(1)).
		WithMargin(value.NewSpacingAll(1)).
		WithBorder(true).
		WithSize(value.NewSizeExact(80, 24)).
		WithAlignment(value.NewAlignmentCenter())

	// Modify every property
	_ = original.WithContent("Modified")
	_ = original.WithPadding(value.NewSpacingAll(2))
	_ = original.WithMargin(value.NewSpacingAll(2))
	_ = original.WithBorder(false)
	_ = original.WithSize(value.NewSizeExact(100, 30))
	_ = original.WithAlignment(value.NewAlignmentDefault())

	// Original should be unchanged
	if original.Content() != "Original" {
		t.Error("Content was mutated")
	}
	if original.Padding().Top() != 1 {
		t.Error("Padding was mutated")
	}
	if original.Margin().Top() != 1 {
		t.Error("Margin was mutated")
	}
	if !original.HasBorder() {
		t.Error("Border was mutated")
	}
	if original.Size().Width() != 80 {
		t.Error("Size was mutated")
	}
	if !original.Alignment().IsCenter() {
		t.Error("Alignment was mutated")
	}
}
