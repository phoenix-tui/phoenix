package viewport

import (
	"reflect"
	"testing"
)

// Test AppendLine functionality (Issue #1 fix)
func TestViewport_AppendLine(t *testing.T) {
	v := New(80, 24)
	v = v.SetLines([]string{"Line 1", "Line 2"})

	v2 := v.AppendLine("Line 3")

	content := v2.VisibleLines()
	want := []string{"Line 1", "Line 2", "Line 3"}

	if !reflect.DeepEqual(content, want) {
		t.Errorf("AppendLine() content = %v, want %v", content, want)
	}

	// Original unchanged (immutability)
	originalContent := v.VisibleLines()
	if len(originalContent) != 2 {
		t.Error("Original viewport was mutated")
	}
}

func TestViewport_AppendLine_ToEmptyViewport(t *testing.T) {
	v := New(80, 24)

	v2 := v.AppendLine("First line")

	content := v2.VisibleLines()
	if len(content) != 1 {
		t.Errorf("AppendLine to empty viewport: got %d lines, want 1", len(content))
	}
	if content[0] != "First line" {
		t.Errorf("AppendLine content = %q, want %q", content[0], "First line")
	}
}

func TestViewport_AppendLine_Multiple(t *testing.T) {
	v := New(80, 24)

	// Chain multiple appends
	v = v.AppendLine("Line 1").
		AppendLine("Line 2").
		AppendLine("Line 3")

	content := v.VisibleLines()
	want := []string{"Line 1", "Line 2", "Line 3"}

	if !reflect.DeepEqual(content, want) {
		t.Errorf("Multiple AppendLine() = %v, want %v", content, want)
	}
}

// Test AppendLines functionality (Issue #1 fix - batch append)
func TestViewport_AppendLines(t *testing.T) {
	v := New(80, 24)
	v = v.SetLines([]string{"Line 1", "Line 2"})

	newLines := []string{"Line 3", "Line 4", "Line 5"}
	v2 := v.AppendLines(newLines)

	content := v2.VisibleLines()
	want := []string{"Line 1", "Line 2", "Line 3", "Line 4", "Line 5"}

	if !reflect.DeepEqual(content, want) {
		t.Errorf("AppendLines() content = %v, want %v", content, want)
	}

	// Original unchanged
	originalContent := v.VisibleLines()
	if len(originalContent) != 2 {
		t.Error("Original viewport was mutated")
	}
}

func TestViewport_AppendLines_ToEmptyViewport(t *testing.T) {
	v := New(80, 24)

	lines := []string{"First", "Second", "Third"}
	v2 := v.AppendLines(lines)

	content := v2.VisibleLines()
	if !reflect.DeepEqual(content, lines) {
		t.Errorf("AppendLines to empty viewport = %v, want %v", content, lines)
	}
}

func TestViewport_AppendLines_EmptySlice(t *testing.T) {
	v := New(80, 24)
	v = v.SetLines([]string{"Line 1", "Line 2"})

	v2 := v.AppendLines([]string{})

	content := v2.VisibleLines()
	want := []string{"Line 1", "Line 2"}

	if !reflect.DeepEqual(content, want) {
		t.Errorf("AppendLines with empty slice = %v, want %v", content, want)
	}
}

// Test ScrollToBottom/ScrollToTop public API (Issue #2 - expose existing methods)
func TestViewport_ScrollToBottom_API(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := New(80, 20).SetLines(content)

	v2 := v.ScrollToBottom()

	if !v2.IsAtBottom() {
		t.Error("ScrollToBottom() should scroll to bottom")
	}

	// Original unchanged
	if v.IsAtBottom() {
		t.Error("Original viewport was mutated")
	}
}

func TestViewport_ScrollToTop_API(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := New(80, 20).
		SetLines(content).
		ScrollToBottom() // Start at bottom

	v2 := v.ScrollToTop()

	if !v2.IsAtTop() {
		t.Error("ScrollToTop() should scroll to top")
	}
	if v2.ScrollOffset() != 0 {
		t.Errorf("ScrollToTop() offset = %d, want 0", v2.ScrollOffset())
	}
}

// Test SetYOffset public API (Issue #2 - precise scroll control)
func TestViewport_SetYOffset(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := New(80, 20).SetLines(content)

	v2 := v.SetYOffset(30)

	if v2.ScrollOffset() != 30 {
		t.Errorf("SetYOffset(30) offset = %d, want 30", v2.ScrollOffset())
	}

	// Original unchanged
	if v.ScrollOffset() != 0 {
		t.Error("Original viewport was mutated")
	}
}

func TestViewport_SetYOffset_ClampNegative(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := New(80, 20).SetLines(content)

	v2 := v.SetYOffset(-10)

	if v2.ScrollOffset() != 0 {
		t.Errorf("SetYOffset(-10) should clamp to 0, got %d", v2.ScrollOffset())
	}
}

func TestViewport_SetYOffset_ClampTooLarge(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := New(80, 20).SetLines(content)

	// Max offset is 80 (100 - 20)
	v2 := v.SetYOffset(200)

	if v2.ScrollOffset() != 80 {
		t.Errorf("SetYOffset(200) should clamp to 80, got %d", v2.ScrollOffset())
	}
}

// Test fluent API chaining with new methods
func TestViewport_FluentAPIChaining(t *testing.T) {
	v := New(80, 20).
		SetLines([]string{"Initial"}).
		AppendLine("Line 2").
		AppendLines([]string{"Line 3", "Line 4"}).
		ScrollToBottom().
		SetYOffset(1).
		ScrollToTop()

	if !v.IsAtTop() {
		t.Error("Fluent API chain should result in viewport at top")
	}

	content := v.VisibleLines()
	if len(content) != 4 {
		t.Errorf("Fluent API chain: got %d lines, want 4", len(content))
	}
}
