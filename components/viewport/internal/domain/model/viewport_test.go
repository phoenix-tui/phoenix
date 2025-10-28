package model

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewViewport(t *testing.T) {
	v := NewViewport(80, 24)

	if v.size.Width() != 80 {
		t.Errorf("Width = %d, want 80", v.size.Width())
	}
	if v.size.Height() != 24 {
		t.Errorf("Height = %d, want 24", v.size.Height())
	}
	if v.scrollOffset.Offset() != 0 {
		t.Errorf("ScrollOffset = %d, want 0", v.scrollOffset.Offset())
	}
	if len(v.content) != 0 {
		t.Errorf("Content length = %d, want 0", len(v.content))
	}
	if v.followMode {
		t.Error("FollowMode should be false by default")
	}
	if v.wrapLines {
		t.Error("WrapLines should be false by default")
	}
}

func TestNewViewportWithContent(t *testing.T) {
	content := []string{"Line 1", "Line 2", "Line 3"}
	v := NewViewportWithContent(content, 80, 24)

	if v.TotalLines() != 3 {
		t.Errorf("TotalLines = %d, want 3", v.TotalLines())
	}
	if !reflect.DeepEqual(v.content, content) {
		t.Errorf("Content = %v, want %v", v.content, content)
	}
}

func TestViewport_WithContent(t *testing.T) {
	v := NewViewport(80, 24)
	content := []string{"Line 1", "Line 2", "Line 3"}

	v2 := v.WithContent(content)

	// Check new viewport has content.
	if v2.TotalLines() != 3 {
		t.Errorf("TotalLines = %d, want 3", v2.TotalLines())
	}

	// Check original is unchanged (immutability)
	if v.TotalLines() != 0 {
		t.Errorf("Original viewport was mutated: TotalLines = %d, want 0", v.TotalLines())
	}
}

func TestViewport_WithContent_FollowMode(t *testing.T) {
	v := NewViewport(80, 5).WithFollowMode(true)

	// Content larger than viewport.
	content := make([]string, 10)
	for i := range content {
		content[i] = "Line"
	}

	v2 := v.WithContent(content)

	// Should scroll to bottom in follow mode.
	if !v2.IsAtBottom() {
		t.Error("WithContent in follow mode should scroll to bottom")
	}
	if v2.ScrollOffset() != 5 { // 10 lines - 5 height = 5 offset
		t.Errorf("ScrollOffset = %d, want 5", v2.ScrollOffset())
	}
}

func TestViewport_WithSize(t *testing.T) {
	v := NewViewport(80, 24)
	v2 := v.WithSize(100, 30)

	if v2.size.Width() != 100 {
		t.Errorf("Width = %d, want 100", v2.size.Width())
	}
	if v2.size.Height() != 30 {
		t.Errorf("Height = %d, want 30", v2.size.Height())
	}

	// Check original is unchanged.
	if v.size.Width() != 80 || v.size.Height() != 24 {
		t.Error("Original viewport was mutated")
	}
}

func TestViewport_WithSize_ClampsOffset(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).
		WithContent(content).
		ScrollToBottom()

	// Offset should be 80 (100 - 20)
	if v.ScrollOffset() != 80 {
		t.Errorf("Initial offset = %d, want 80", v.ScrollOffset())
	}

	// Increase height to 50.
	v2 := v.WithSize(80, 50)

	// New max offset should be 50 (100 - 50)
	if v2.ScrollOffset() > 50 {
		t.Errorf("Offset after resize = %d, should be clamped to <= 50", v2.ScrollOffset())
	}
}

func TestViewport_WithFollowMode(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).WithContent(content)

	// Initially at top.
	if v.IsAtBottom() {
		t.Error("Should not be at bottom initially")
	}

	// Enable follow mode.
	v2 := v.WithFollowMode(true)

	if !v2.FollowMode() {
		t.Error("FollowMode should be true")
	}
	if !v2.IsAtBottom() {
		t.Error("Enabling follow mode should scroll to bottom")
	}

	// Original unchanged.
	if v.FollowMode() {
		t.Error("Original viewport was mutated")
	}
}

func TestViewport_WithWrapLines(t *testing.T) {
	v := NewViewport(80, 24)

	if v.WrapLines() {
		t.Error("WrapLines should be false by default")
	}

	v2 := v.WithWrapLines(true)

	if !v2.WrapLines() {
		t.Error("WrapLines should be true")
	}
	if v.WrapLines() {
		t.Error("Original viewport was mutated")
	}
}

func TestViewport_ScrollUp(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).
		WithContent(content).
		ScrollDown(50) // Go to offset 50

	v2 := v.ScrollUp(10)

	if v2.ScrollOffset() != 40 {
		t.Errorf("ScrollOffset = %d, want 40", v2.ScrollOffset())
	}

	// Original unchanged.
	if v.ScrollOffset() != 50 {
		t.Error("Original viewport was mutated")
	}
}

func TestViewport_ScrollUp_DisablesFollowMode(t *testing.T) {
	v := NewViewport(80, 20).
		WithFollowMode(true).
		ScrollUp(5)

	if v.FollowMode() {
		t.Error("ScrollUp should disable follow mode")
	}
}

func TestViewport_ScrollDown(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).
		WithContent(content)

	v2 := v.ScrollDown(10)

	if v2.ScrollOffset() != 10 {
		t.Errorf("ScrollOffset = %d, want 10", v2.ScrollOffset())
	}
}

func TestViewport_ScrollDown_EnablesFollowModeAtBottom(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).
		WithContent(content).
		ScrollDown(70) // Not quite at bottom

	if v.FollowMode() {
		t.Error("Should not be in follow mode yet")
	}

	v2 := v.ScrollDown(20) // Scroll to bottom

	if !v2.FollowMode() {
		t.Error("Scrolling to bottom should enable follow mode")
	}
}

func TestViewport_ScrollToTop(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).
		WithContent(content).
		ScrollDown(50).
		ScrollToTop()

	if v.ScrollOffset() != 0 {
		t.Errorf("ScrollOffset = %d, want 0", v.ScrollOffset())
	}
	if !v.IsAtTop() {
		t.Error("Should be at top")
	}
	if v.FollowMode() {
		t.Error("Follow mode should be disabled")
	}
}

func TestViewport_ScrollToBottom(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).
		WithContent(content).
		ScrollToBottom()

	if !v.IsAtBottom() {
		t.Error("Should be at bottom")
	}
	if !v.FollowMode() {
		t.Error("Follow mode should be enabled")
	}
	if v.ScrollOffset() != 80 { // 100 - 20
		t.Errorf("ScrollOffset = %d, want 80", v.ScrollOffset())
	}
}

func TestViewport_PageUp(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).
		WithContent(content).
		ScrollDown(50).
		PageUp()

	// Should scroll up by viewport height (20)
	if v.ScrollOffset() != 30 {
		t.Errorf("ScrollOffset = %d, want 30", v.ScrollOffset())
	}
}

func TestViewport_PageDown(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).
		WithContent(content).
		PageDown()

	// Should scroll down by viewport height (20)
	if v.ScrollOffset() != 20 {
		t.Errorf("ScrollOffset = %d, want 20", v.ScrollOffset())
	}
}

func TestViewport_VisibleLines(t *testing.T) {
	content := []string{"Line 1", "Line 2", "Line 3", "Line 4", "Line 5"}
	v := NewViewport(80, 3).WithContent(content)

	visible := v.VisibleLines()

	want := []string{"Line 1", "Line 2", "Line 3"}
	if !reflect.DeepEqual(visible, want) {
		t.Errorf("VisibleLines = %v, want %v", visible, want)
	}

	// Scroll down.
	v2 := v.ScrollDown(2)
	visible2 := v2.VisibleLines()

	want2 := []string{"Line 3", "Line 4", "Line 5"}
	if !reflect.DeepEqual(visible2, want2) {
		t.Errorf("VisibleLines after scroll = %v, want %v", visible2, want2)
	}
}

func TestViewport_VisibleLines_Truncation(t *testing.T) {
	content := []string{
		"Short",
		"This is a very long line that should be truncated to fit the viewport width",
	}
	v := NewViewport(20, 5).WithContent(content)

	visible := v.VisibleLines()

	if len(visible) != 2 {
		t.Errorf("VisibleLines count = %d, want 2", len(visible))
	}

	// Second line should be truncated.
	if len(visible[1]) > 20 {
		t.Errorf("Line was not truncated: length = %d, want <= 20", len(visible[1]))
	}
}

func TestViewport_VisibleLines_Wrapping(t *testing.T) {
	content := []string{"This is a long line that should be wrapped"}
	v := NewViewport(10, 10).
		WithContent(content).
		WithWrapLines(true)

	visible := v.VisibleLines()

	// Line should be wrapped into multiple lines.
	if len(visible) <= 1 {
		t.Errorf("Line was not wrapped: got %d lines", len(visible))
	}

	// Each wrapped line should fit in width.
	for i, line := range visible {
		if len(line) > 10 {
			t.Errorf("Wrapped line %d exceeds width: length = %d", i, len(line))
		}
	}
}

func TestViewport_CanScrollUp(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).WithContent(content)

	if v.CanScrollUp() {
		t.Error("Should not be able to scroll up from top")
	}

	v2 := v.ScrollDown(10)
	if !v2.CanScrollUp() {
		t.Error("Should be able to scroll up from middle")
	}
}

func TestViewport_CanScrollDown(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).WithContent(content)

	if !v.CanScrollDown() {
		t.Error("Should be able to scroll down from top")
	}

	v2 := v.ScrollToBottom()
	if v2.CanScrollDown() {
		t.Error("Should not be able to scroll down from bottom")
	}
}

func TestViewport_IsAtTop(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).WithContent(content)

	if !v.IsAtTop() {
		t.Error("Should be at top initially")
	}

	v2 := v.ScrollDown(10)
	if v2.IsAtTop() {
		t.Error("Should not be at top after scrolling")
	}
}

func TestViewport_IsAtBottom(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).WithContent(content)

	if v.IsAtBottom() {
		t.Error("Should not be at bottom initially")
	}

	v2 := v.ScrollToBottom()
	if !v2.IsAtBottom() {
		t.Error("Should be at bottom after ScrollToBottom")
	}
}

func TestViewport_Immutability(t *testing.T) {
	original := NewViewport(80, 24).
		WithContent([]string{"Line 1", "Line 2", "Line 3"})

	// Perform various operations.
	_ = original.WithContent([]string{"New content"})
	_ = original.WithSize(100, 30)
	_ = original.WithFollowMode(true)
	_ = original.WithWrapLines(true)
	_ = original.ScrollUp(5)
	_ = original.ScrollDown(5)
	_ = original.ScrollToTop()
	_ = original.ScrollToBottom()
	_ = original.PageUp()
	_ = original.PageDown()

	// Original should remain unchanged.
	if original.TotalLines() != 3 {
		t.Error("Content was mutated")
	}
	if original.size.Width() != 80 || original.size.Height() != 24 {
		t.Error("Size was mutated")
	}
	if original.FollowMode() {
		t.Error("FollowMode was mutated")
	}
	if original.WrapLines() {
		t.Error("WrapLines was mutated")
	}
	if original.ScrollOffset() != 0 {
		t.Error("ScrollOffset was mutated")
	}
}

func TestViewport_EmptyContent(t *testing.T) {
	v := NewViewport(80, 24)

	if v.TotalLines() != 0 {
		t.Errorf("TotalLines = %d, want 0", v.TotalLines())
	}
	if len(v.VisibleLines()) != 0 {
		t.Errorf("VisibleLines count = %d, want 0", len(v.VisibleLines()))
	}
	if !v.IsAtTop() {
		t.Error("Empty viewport should be at top")
	}
	if !v.IsAtBottom() {
		t.Error("Empty viewport should be at bottom")
	}
	if v.CanScrollUp() {
		t.Error("Empty viewport should not allow scrolling up")
	}
	if v.CanScrollDown() {
		t.Error("Empty viewport should not allow scrolling down")
	}
}

func TestViewport_ContentSmallerThanViewport(t *testing.T) {
	content := []string{"Line 1", "Line 2", "Line 3"}
	v := NewViewport(80, 10).WithContent(content)

	if v.CanScrollDown() {
		t.Error("Should not be able to scroll when content fits")
	}
	if !v.IsAtBottom() {
		t.Error("Should be at bottom when content fits")
	}

	visible := v.VisibleLines()
	if !reflect.DeepEqual(visible, content) {
		t.Errorf("VisibleLines = %v, want %v", visible, content)
	}
}

func TestViewport_UnicodeContent(t *testing.T) {
	content := []string{
		"Hello 世界",
		"Emoji: 😀🎉",
		"Mixed: ABC 日本語 123",
	}
	v := NewViewport(20, 5).WithContent(content)

	visible := v.VisibleLines()

	if len(visible) != 3 {
		t.Errorf("VisibleLines count = %d, want 3", len(visible))
	}

	// Lines should be properly truncated respecting Unicode width.
	for i, line := range visible {
		if strings.Contains(line, "�") {
			t.Errorf("Line %d contains replacement character: %s", i, line)
		}
	}
}

func TestViewport_LargeContent(t *testing.T) {
	// Test with large content (10k lines)
	content := make([]string, 10000)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).WithContent(content)

	// Should handle large content efficiently.
	visible := v.VisibleLines()
	if len(visible) != 20 {
		t.Errorf("VisibleLines count = %d, want 20", len(visible))
	}

	// Scroll operations should be fast.
	v2 := v.ScrollToBottom()
	if !v2.IsAtBottom() {
		t.Error("Failed to scroll to bottom of large content")
	}

	v3 := v2.ScrollToTop()
	if !v3.IsAtTop() {
		t.Error("Failed to scroll to top of large content")
	}
}

// New tests for WithScrollOffset (Issue #2 - SetYOffset support)
func TestViewport_WithScrollOffset(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).WithContent(content)

	// Set specific offset.
	v2 := v.WithScrollOffset(30)

	if v2.ScrollOffset() != 30 {
		t.Errorf("ScrollOffset = %d, want 30", v2.ScrollOffset())
	}

	// Follow mode should be disabled.
	if v2.FollowMode() {
		t.Error("WithScrollOffset should disable follow mode")
	}

	// Original unchanged.
	if v.ScrollOffset() != 0 {
		t.Error("Original viewport was mutated")
	}
}

func TestViewport_WithScrollOffset_ClampsToZero(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).WithContent(content)

	// Try negative offset.
	v2 := v.WithScrollOffset(-10)

	if v2.ScrollOffset() != 0 {
		t.Errorf("Negative offset should be clamped to 0, got %d", v2.ScrollOffset())
	}
}

func TestViewport_WithScrollOffset_ClampsToMaxOffset(t *testing.T) {
	content := make([]string, 100)
	for i := range content {
		content[i] = "Line"
	}

	v := NewViewport(80, 20).WithContent(content)

	// Max offset is 80 (100 lines - 20 height)
	// Try offset beyond max.
	v2 := v.WithScrollOffset(200)

	if v2.ScrollOffset() != 80 {
		t.Errorf("Offset should be clamped to 80, got %d", v2.ScrollOffset())
	}
}

// New test for Content() accessor (Issue #1 - AppendLine support)
func TestViewport_Content(t *testing.T) {
	content := []string{"Line 1", "Line 2", "Line 3"}
	v := NewViewport(80, 24).WithContent(content)

	retrievedContent := v.Content()

	// Should return copy of content.
	if !reflect.DeepEqual(retrievedContent, content) {
		t.Errorf("Content() = %v, want %v", retrievedContent, content)
	}

	// Modifying returned slice should not affect viewport (defensive copy)
	retrievedContent[0] = "Modified"

	retrievedContent2 := v.Content()
	if retrievedContent2[0] != "Line 1" {
		t.Error("Viewport content was mutated through Content() return value")
	}
}
