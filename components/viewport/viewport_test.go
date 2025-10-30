package viewport

import (
	"reflect"
	"testing"

	"github.com/phoenix-tui/phoenix/tea"
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

	// Chain multiple appends.
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

	// Original unchanged.
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

	// Original unchanged.
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

	// Original unchanged.
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

// Test fluent API chaining with new methods.
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

// ============================================================================
// Constructor Tests
// ============================================================================

func TestViewport_NewWithContent(t *testing.T) {
	content := "Line 1\nLine 2\nLine 3"
	v := NewWithContent(content, 80, 20)

	lines := v.VisibleLines()
	if len(lines) != 3 {
		t.Errorf("NewWithContent: got %d lines, want 3", len(lines))
	}
	if lines[0] != "Line 1" {
		t.Errorf("NewWithContent: first line = %q, want %q", lines[0], "Line 1")
	}
}

func TestViewport_NewWithContent_EmptyString(t *testing.T) {
	v := NewWithContent("", 80, 20)
	lines := v.VisibleLines()
	if len(lines) != 1 || lines[0] != "" {
		t.Errorf("NewWithContent with empty string: got %v, want ['']", lines)
	}
}

func TestViewport_NewWithContent_SingleLine(t *testing.T) {
	v := NewWithContent("Single line without newline", 80, 20)
	lines := v.VisibleLines()
	if len(lines) != 1 {
		t.Errorf("NewWithContent single line: got %d lines, want 1", len(lines))
	}
	if lines[0] != "Single line without newline" {
		t.Errorf("NewWithContent: line = %q, want %q", lines[0], "Single line without newline")
	}
}

func TestViewport_NewWithLines(t *testing.T) {
	lines := []string{"Line 1", "Line 2", "Line 3"}
	v := NewWithLines(lines, 80, 20)

	result := v.VisibleLines()
	if !reflect.DeepEqual(result, lines) {
		t.Errorf("NewWithLines: got %v, want %v", result, lines)
	}
}

func TestViewport_NewWithLines_EmptySlice(t *testing.T) {
	v := NewWithLines([]string{}, 80, 20)
	lines := v.VisibleLines()
	if len(lines) != 0 {
		t.Errorf("NewWithLines with empty slice: got %d lines, want 0", len(lines))
	}
}

// ============================================================================
// Configuration Method Tests
// ============================================================================

func TestViewport_FollowMode(t *testing.T) {
	v := New(80, 20).FollowMode(true)

	// Add content - should auto-scroll to bottom
	v = v.SetLines(make([]string, 100))

	if !v.IsAtBottom() {
		t.Error("FollowMode(true): viewport should auto-scroll to bottom when content is added")
	}

	// Disable follow mode
	v2 := v.ScrollToTop().FollowMode(false)

	// Add more content - should NOT auto-scroll
	v2 = v2.AppendLine("New line")

	if !v2.IsAtTop() {
		t.Error("FollowMode(false): viewport should stay at top when content is added")
	}
}

func TestViewport_WrapLines(t *testing.T) {
	v := New(10, 5).WrapLines(true)

	// Set content wider than viewport
	v = v.SetLines([]string{"This is a very long line that should wrap"})

	// With wrapping enabled, should see multiple visible lines
	visibleLines := v.VisibleLines()
	if len(visibleLines) < 2 {
		t.Errorf("WrapLines(true): expected multiple lines due to wrapping, got %d", len(visibleLines))
	}
}

func TestViewport_WrapLines_Disabled(t *testing.T) {
	v := New(10, 5).WrapLines(false)

	// Set content wider than viewport
	v = v.SetLines([]string{"This is a very long line that should NOT wrap"})

	// Without wrapping, should see just 1 line (truncated)
	visibleLines := v.VisibleLines()
	if len(visibleLines) != 1 {
		t.Errorf("WrapLines(false): expected 1 line (truncated), got %d", len(visibleLines))
	}
}

func TestViewport_MouseEnabled(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)

	// Mouse wheel should work when enabled
	// We'll test this via Update() later

	v2 := v.MouseEnabled(false)

	// Verify immutability
	if v == v2 {
		t.Error("MouseEnabled should return new instance")
	}
}

func TestViewport_SetContent(t *testing.T) {
	v := New(80, 20)
	v2 := v.SetContent("Line 1\nLine 2\nLine 3")

	lines := v2.VisibleLines()
	if len(lines) != 3 {
		t.Errorf("SetContent: got %d lines, want 3", len(lines))
	}
	if lines[1] != "Line 2" {
		t.Errorf("SetContent: second line = %q, want %q", lines[1], "Line 2")
	}

	// Original unchanged
	originalLines := v.VisibleLines()
	if len(originalLines) != 0 {
		t.Error("Original viewport was mutated")
	}
}

func TestViewport_SetSize(t *testing.T) {
	v := New(80, 20)
	v = v.SetLines(make([]string, 100))

	v2 := v.SetSize(40, 10)

	if v2.Height() != 10 {
		t.Errorf("SetSize: height = %d, want 10", v2.Height())
	}

	// Original unchanged
	if v.Height() != 20 {
		t.Error("Original viewport height was mutated")
	}
}

func TestViewport_SetSize_AffectsVisibleLines(t *testing.T) {
	v := New(80, 5)
	v = v.SetLines([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"})

	// Should see first 5 lines
	lines := v.VisibleLines()
	if len(lines) != 5 {
		t.Errorf("Initial: got %d visible lines, want 5", len(lines))
	}

	// Increase height to 8
	v2 := v.SetSize(80, 8)

	// Should now see first 8 lines
	lines2 := v2.VisibleLines()
	if len(lines2) != 8 {
		t.Errorf("After SetSize(80, 8): got %d visible lines, want 8", len(lines2))
	}
}

// ============================================================================
// tea.Model Implementation Tests
// ============================================================================

func TestViewport_Init(t *testing.T) {
	v := New(80, 20)
	cmd := v.Init()

	if cmd != nil {
		t.Error("Init() should return nil Cmd")
	}
}

func TestViewport_Update_KeyMsg_Up(t *testing.T) {
	v := New(80, 5)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Simulate Up arrow key
	msg := tea.KeyMsg{Type: tea.KeyUp}
	v2, cmd := v.Update(msg)

	if cmd != nil {
		t.Error("Update should return nil Cmd for key messages")
	}

	if v2.ScrollOffset() != 49 {
		t.Errorf("After Up key: offset = %d, want 49", v2.ScrollOffset())
	}
}

func TestViewport_Update_KeyMsg_Down(t *testing.T) {
	v := New(80, 5)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Simulate Down arrow key
	msg := tea.KeyMsg{Type: tea.KeyDown}
	v2, cmd := v.Update(msg)

	if cmd != nil {
		t.Error("Update should return nil Cmd")
	}

	if v2.ScrollOffset() != 51 {
		t.Errorf("After Down key: offset = %d, want 51", v2.ScrollOffset())
	}
}

func TestViewport_Update_KeyMsg_PageUp(t *testing.T) {
	v := New(80, 20)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Simulate Page Up key
	msg := tea.KeyMsg{Type: tea.KeyPgUp}
	v2, cmd := v.Update(msg)

	if cmd != nil {
		t.Error("Update should return nil Cmd")
	}

	// Should scroll up by ~height (20 lines)
	if v2.ScrollOffset() >= 50 {
		t.Errorf("After PageUp: offset = %d, should be < 50", v2.ScrollOffset())
	}
}

func TestViewport_Update_KeyMsg_PageDown(t *testing.T) {
	v := New(80, 20)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(10)

	// Simulate Page Down key
	msg := tea.KeyMsg{Type: tea.KeyPgDown}
	v2, cmd := v.Update(msg)

	if cmd != nil {
		t.Error("Update should return nil Cmd")
	}

	// Should scroll down by ~height (20 lines)
	if v2.ScrollOffset() <= 10 {
		t.Errorf("After PageDown: offset = %d, should be > 10", v2.ScrollOffset())
	}
}

func TestViewport_Update_KeyMsg_Home(t *testing.T) {
	v := New(80, 20)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Simulate Home key
	msg := tea.KeyMsg{Type: tea.KeyHome}
	v2, cmd := v.Update(msg)

	if cmd != nil {
		t.Error("Update should return nil Cmd")
	}

	if v2.ScrollOffset() != 0 {
		t.Errorf("After Home key: offset = %d, want 0", v2.ScrollOffset())
	}
}

func TestViewport_Update_KeyMsg_End(t *testing.T) {
	v := New(80, 20)
	v = v.SetLines(make([]string, 100))

	// Simulate End key
	msg := tea.KeyMsg{Type: tea.KeyEnd}
	v2, cmd := v.Update(msg)

	if cmd != nil {
		t.Error("Update should return nil Cmd")
	}

	if !v2.IsAtBottom() {
		t.Error("After End key: viewport should be at bottom")
	}
}

func TestViewport_Update_KeyMsg_HalfPageUp(t *testing.T) {
	v := New(80, 20)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Simulate Ctrl+U (half page up)
	msg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'u', Ctrl: true}
	v2, cmd := v.Update(msg)

	if cmd != nil {
		t.Error("Update should return nil Cmd")
	}

	// Should scroll up by half height (10 lines)
	expectedOffset := 50 - 10
	if v2.ScrollOffset() != expectedOffset {
		t.Errorf("After Ctrl+U: offset = %d, want %d", v2.ScrollOffset(), expectedOffset)
	}
}

func TestViewport_Update_KeyMsg_HalfPageDown(t *testing.T) {
	v := New(80, 20)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(10)

	// Simulate Ctrl+D (half page down)
	msg := tea.KeyMsg{Type: tea.KeyRune, Rune: 'd', Ctrl: true}
	v2, cmd := v.Update(msg)

	if cmd != nil {
		t.Error("Update should return nil Cmd")
	}

	// Should scroll down by half height (10 lines)
	expectedOffset := 10 + 10
	if v2.ScrollOffset() != expectedOffset {
		t.Errorf("After Ctrl+D: offset = %d, want %d", v2.ScrollOffset(), expectedOffset)
	}
}

func TestViewport_Update_MouseMsg_WheelUp(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Simulate mouse wheel up
	msg := tea.MouseMsg{Button: tea.MouseButtonWheelUp}
	v2, cmd := v.Update(msg)

	if cmd != nil {
		t.Error("Update should return nil Cmd")
	}

	// Should scroll up by 3 lines
	expectedOffset := 50 - 3
	if v2.ScrollOffset() != expectedOffset {
		t.Errorf("After mouse wheel up: offset = %d, want %d", v2.ScrollOffset(), expectedOffset)
	}
}

func TestViewport_Update_MouseMsg_WheelDown(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(10)

	// Simulate mouse wheel down
	msg := tea.MouseMsg{Button: tea.MouseButtonWheelDown}
	v2, cmd := v.Update(msg)

	if cmd != nil {
		t.Error("Update should return nil Cmd")
	}

	// Should scroll down by 3 lines (default)
	expectedOffset := 10 + 3
	if v2.ScrollOffset() != expectedOffset {
		t.Errorf("After mouse wheel down: offset = %d, want %d", v2.ScrollOffset(), expectedOffset)
	}
}

// ============================================================================
// Wheel Scrolling Configuration Tests (Week 15 Day 5-6)
// ============================================================================

func TestViewport_SetWheelScrollLines_CustomValue(t *testing.T) {
	v := New(80, 20).MouseEnabled(true).SetWheelScrollLines(5)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Simulate mouse wheel up with custom scroll (5 lines)
	msg := tea.MouseMsg{Button: tea.MouseButtonWheelUp}
	v2, _ := v.Update(msg)

	expectedOffset := 50 - 5
	if v2.ScrollOffset() != expectedOffset {
		t.Errorf("After wheel up (5 lines): offset = %d, want %d", v2.ScrollOffset(), expectedOffset)
	}

	// Simulate mouse wheel down with custom scroll (5 lines)
	msg2 := tea.MouseMsg{Button: tea.MouseButtonWheelDown}
	v3, _ := v2.Update(msg2)

	expectedOffset2 := expectedOffset + 5
	if v3.ScrollOffset() != expectedOffset2 {
		t.Errorf("After wheel down (5 lines): offset = %d, want %d", v3.ScrollOffset(), expectedOffset2)
	}
}

func TestViewport_SetWheelScrollLines_DefaultValue(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Should use default 3 lines
	msg := tea.MouseMsg{Button: tea.MouseButtonWheelUp}
	v2, _ := v.Update(msg)

	expectedOffset := 50 - 3
	if v2.ScrollOffset() != expectedOffset {
		t.Errorf("Default wheel scroll: offset = %d, want %d (3 lines)", v2.ScrollOffset(), expectedOffset)
	}
}

func TestViewport_SetWheelScrollLines_MinimumOne(t *testing.T) {
	// Test that 0 is clamped to 1
	v := New(80, 20).MouseEnabled(true).SetWheelScrollLines(0)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	msg := tea.MouseMsg{Button: tea.MouseButtonWheelUp}
	v2, _ := v.Update(msg)

	expectedOffset := 50 - 1 // Should scroll by 1 (clamped)
	if v2.ScrollOffset() != expectedOffset {
		t.Errorf("SetWheelScrollLines(0) should clamp to 1: offset = %d, want %d", v2.ScrollOffset(), expectedOffset)
	}
}

func TestViewport_SetWheelScrollLines_NegativeValue(t *testing.T) {
	// Test that negative values are clamped to 1
	v := New(80, 20).MouseEnabled(true).SetWheelScrollLines(-5)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	msg := tea.MouseMsg{Button: tea.MouseButtonWheelUp}
	v2, _ := v.Update(msg)

	expectedOffset := 50 - 1 // Should scroll by 1 (clamped)
	if v2.ScrollOffset() != expectedOffset {
		t.Errorf("SetWheelScrollLines(-5) should clamp to 1: offset = %d, want %d", v2.ScrollOffset(), expectedOffset)
	}
}

func TestViewport_SetWheelScrollLines_LargeValue(t *testing.T) {
	v := New(80, 20).MouseEnabled(true).SetWheelScrollLines(10)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Scroll up by 10 lines
	msg := tea.MouseMsg{Button: tea.MouseButtonWheelUp}
	v2, _ := v.Update(msg)

	expectedOffset := 50 - 10
	if v2.ScrollOffset() != expectedOffset {
		t.Errorf("After wheel up (10 lines): offset = %d, want %d", v2.ScrollOffset(), expectedOffset)
	}
}

func TestViewport_SetWheelScrollLines_BoundsTop(t *testing.T) {
	v := New(80, 20).MouseEnabled(true).SetWheelScrollLines(5)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(2) // Near top

	// Scroll up by 5 (would go to -3, should clamp to 0)
	msg := tea.MouseMsg{Button: tea.MouseButtonWheelUp}
	v2, _ := v.Update(msg)

	if v2.ScrollOffset() != 0 {
		t.Errorf("Wheel scroll past top: offset = %d, want 0 (clamped)", v2.ScrollOffset())
	}
}

func TestViewport_SetWheelScrollLines_BoundsBottom(t *testing.T) {
	v := New(80, 20).MouseEnabled(true).SetWheelScrollLines(5)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(77) // Near bottom (max is 80)

	// Scroll down by 5 (would go to 82, should clamp to 80)
	msg := tea.MouseMsg{Button: tea.MouseButtonWheelDown}
	v2, _ := v.Update(msg)

	maxOffset := 100 - 20 // totalLines - height
	if v2.ScrollOffset() != maxOffset {
		t.Errorf("Wheel scroll past bottom: offset = %d, want %d (clamped)", v2.ScrollOffset(), maxOffset)
	}
}

func TestViewport_SetWheelScrollLines_SmallContent(t *testing.T) {
	v := New(80, 20).MouseEnabled(true).SetWheelScrollLines(5)
	v = v.SetLines([]string{"Line 1", "Line 2", "Line 3"}) // Only 3 lines

	// Try to scroll down (content fits entirely, should stay at 0)
	msg := tea.MouseMsg{Button: tea.MouseButtonWheelDown}
	v2, _ := v.Update(msg)

	if v2.ScrollOffset() != 0 {
		t.Errorf("Wheel scroll with small content: offset = %d, want 0", v2.ScrollOffset())
	}

	// Try to scroll up (already at top)
	msg2 := tea.MouseMsg{Button: tea.MouseButtonWheelUp}
	v3, _ := v2.Update(msg2)

	if v3.ScrollOffset() != 0 {
		t.Errorf("Wheel scroll up with small content: offset = %d, want 0", v3.ScrollOffset())
	}
}

func TestViewport_SetWheelScrollLines_MultipleWheels(t *testing.T) {
	v := New(80, 20).MouseEnabled(true).SetWheelScrollLines(3)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// First wheel up
	v, _ = v.Update(tea.MouseMsg{Button: tea.MouseButtonWheelUp})
	if v.ScrollOffset() != 47 {
		t.Errorf("After 1st wheel: offset = %d, want 47", v.ScrollOffset())
	}

	// Second wheel up
	v, _ = v.Update(tea.MouseMsg{Button: tea.MouseButtonWheelUp})
	if v.ScrollOffset() != 44 {
		t.Errorf("After 2nd wheel: offset = %d, want 44", v.ScrollOffset())
	}

	// Third wheel up
	v, _ = v.Update(tea.MouseMsg{Button: tea.MouseButtonWheelUp})
	if v.ScrollOffset() != 41 {
		t.Errorf("After 3rd wheel: offset = %d, want 41", v.ScrollOffset())
	}
}

func TestViewport_SetWheelScrollLines_Immutability(t *testing.T) {
	v1 := New(80, 20).MouseEnabled(true)
	v2 := v1.SetWheelScrollLines(5)

	// Original should still use default (3 lines)
	v1 = v1.SetLines(make([]string, 100)).SetYOffset(50)
	v1, _ = v1.Update(tea.MouseMsg{Button: tea.MouseButtonWheelUp})

	// v1 should scroll by 3 (default)
	if v1.ScrollOffset() != 47 {
		t.Errorf("Original viewport: offset = %d, want 47 (default 3 lines)", v1.ScrollOffset())
	}

	// v2 should scroll by 5 (custom)
	v2 = v2.SetLines(make([]string, 100)).SetYOffset(50)
	v2, _ = v2.Update(tea.MouseMsg{Button: tea.MouseButtonWheelUp})

	if v2.ScrollOffset() != 45 {
		t.Errorf("Modified viewport: offset = %d, want 45 (custom 5 lines)", v2.ScrollOffset())
	}
}

func TestViewport_SetWheelScrollLines_FluentChaining(t *testing.T) {
	v := New(80, 20).
		MouseEnabled(true).
		SetWheelScrollLines(7).
		SetLines(make([]string, 100)).
		SetYOffset(50)

	// Should scroll by 7 lines
	v, _ = v.Update(tea.MouseMsg{Button: tea.MouseButtonWheelUp})

	expectedOffset := 50 - 7
	if v.ScrollOffset() != expectedOffset {
		t.Errorf("Fluent API chaining: offset = %d, want %d", v.ScrollOffset(), expectedOffset)
	}
}

func TestViewport_SetWheelScrollLines_PreservedAfterOtherOperations(t *testing.T) {
	v := New(80, 20).
		MouseEnabled(true).
		SetWheelScrollLines(5).
		SetLines(make([]string, 100))

	// Perform various operations
	v = v.ScrollToBottom().ScrollToTop().SetYOffset(50)

	// Wheel scroll should still use 5 lines
	v, _ = v.Update(tea.MouseMsg{Button: tea.MouseButtonWheelUp})

	expectedOffset := 50 - 5
	if v.ScrollOffset() != expectedOffset {
		t.Errorf("After operations: offset = %d, want %d (should preserve 5 lines)", v.ScrollOffset(), expectedOffset)
	}
}

func TestViewport_Update_MouseMsg_Disabled(t *testing.T) {
	v := New(80, 20).MouseEnabled(false)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Simulate mouse wheel up (should be ignored)
	msg := tea.MouseMsg{Button: tea.MouseButtonWheelUp}
	v2, cmd := v.Update(msg)

	if cmd != nil {
		t.Error("Update should return nil Cmd")
	}

	// Offset should NOT change (mouse disabled)
	if v2.ScrollOffset() != 50 {
		t.Errorf("After mouse wheel (disabled): offset = %d, want 50", v2.ScrollOffset())
	}
}

func TestViewport_Update_WindowSizeMsg(t *testing.T) {
	v := New(80, 20)
	v = v.SetLines(make([]string, 100))

	// Simulate window resize
	msg := tea.WindowSizeMsg{Width: 100, Height: 30}
	v2, cmd := v.Update(msg)

	if cmd != nil {
		t.Error("Update should return nil Cmd")
	}

	if v2.Height() != 30 {
		t.Errorf("After WindowSizeMsg: height = %d, want 30", v2.Height())
	}
}

func TestViewport_Update_UnknownMsg(t *testing.T) {
	v := New(80, 20)

	// Simulate unknown message type
	type UnknownMsg struct{}
	v2, cmd := v.Update(UnknownMsg{})

	if cmd != nil {
		t.Error("Update should return nil Cmd for unknown messages")
	}

	// Viewport should remain unchanged
	if v != v2 {
		t.Error("Unknown message should not modify viewport")
	}
}

// ============================================================================
// View/Rendering Tests
// ============================================================================

func TestViewport_View(t *testing.T) {
	v := New(80, 5)
	v = v.SetLines([]string{"Line 1", "Line 2", "Line 3", "Line 4", "Line 5"})

	view := v.View()

	expectedView := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5"
	if view != expectedView {
		t.Errorf("View() = %q, want %q", view, expectedView)
	}
}

func TestViewport_View_Empty(t *testing.T) {
	v := New(80, 20)

	view := v.View()

	if view != "" {
		t.Errorf("View() for empty viewport = %q, want empty string", view)
	}
}

func TestViewport_View_Scrolled(t *testing.T) {
	v := New(80, 3)
	v = v.SetLines([]string{"Line 1", "Line 2", "Line 3", "Line 4", "Line 5"})
	v = v.SetYOffset(2)

	view := v.View()

	expectedView := "Line 3\nLine 4\nLine 5"
	if view != expectedView {
		t.Errorf("View() scrolled = %q, want %q", view, expectedView)
	}
}

// ============================================================================
// Query Method Tests
// ============================================================================

func TestViewport_CanScrollUp(t *testing.T) {
	v := New(80, 20)
	v = v.SetLines(make([]string, 100))

	// At top - cannot scroll up
	if v.CanScrollUp() {
		t.Error("CanScrollUp() at top should be false")
	}

	// Scroll down - can scroll up
	v2 := v.SetYOffset(10)
	if !v2.CanScrollUp() {
		t.Error("CanScrollUp() after scrolling should be true")
	}
}

func TestViewport_CanScrollDown(t *testing.T) {
	v := New(80, 20)
	v = v.SetLines(make([]string, 100))

	// At top - can scroll down
	if !v.CanScrollDown() {
		t.Error("CanScrollDown() at top should be true")
	}

	// At bottom - cannot scroll down
	v2 := v.ScrollToBottom()
	if v2.CanScrollDown() {
		t.Error("CanScrollDown() at bottom should be false")
	}
}

func TestViewport_TotalLines(t *testing.T) {
	v := New(80, 20)
	v = v.SetLines(make([]string, 100))

	if v.TotalLines() != 100 {
		t.Errorf("TotalLines() = %d, want 100", v.TotalLines())
	}
}

func TestViewport_Height(t *testing.T) {
	v := New(80, 20)

	if v.Height() != 20 {
		t.Errorf("Height() = %d, want 20", v.Height())
	}

	v2 := v.SetSize(80, 30)
	if v2.Height() != 30 {
		t.Errorf("Height() after SetSize = %d, want 30", v2.Height())
	}
}

// ============================================================================
// Drag Scrolling Tests (Week 15 Day 3-4)
// ============================================================================

func TestViewport_DragScroll_Start(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Simulate left mouse button press (start drag)
	msg := tea.MouseMsg{
		X:      10,
		Y:      5,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionPress,
	}
	v2, cmd := v.Update(msg)

	if cmd != nil {
		t.Error("Update should return nil Cmd")
	}

	// Viewport should remain at same offset after press
	if v2.ScrollOffset() != 50 {
		t.Errorf("After drag start: offset = %d, want 50", v2.ScrollOffset())
	}

	// Drag state should be recorded (internal check via motion behavior)
	// We'll verify this works by testing motion next
}

func TestViewport_DragScroll_MotionDown(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Start drag at Y=5
	msg1 := tea.MouseMsg{
		X:      10,
		Y:      5,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionPress,
	}
	v, _ = v.Update(msg1)

	// Drag down to Y=10 (delta = +5)
	// Dragging down should scroll content UP (lower offset)
	msg2 := tea.MouseMsg{
		X:      10,
		Y:      10,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionMotion,
	}
	v2, cmd := v.Update(msg2)

	if cmd != nil {
		t.Error("Update should return nil Cmd")
	}

	// Expected: scroll offset = 50 - 5 = 45
	expectedOffset := 50 - 5
	if v2.ScrollOffset() != expectedOffset {
		t.Errorf("After drag down: offset = %d, want %d", v2.ScrollOffset(), expectedOffset)
	}
}

func TestViewport_DragScroll_MotionUp(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Start drag at Y=10
	msg1 := tea.MouseMsg{
		X:      10,
		Y:      10,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionPress,
	}
	v, _ = v.Update(msg1)

	// Drag up to Y=5 (delta = -5)
	// Dragging up should scroll content DOWN (higher offset)
	msg2 := tea.MouseMsg{
		X:      10,
		Y:      5,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionMotion,
	}
	v2, cmd := v.Update(msg2)

	if cmd != nil {
		t.Error("Update should return nil Cmd")
	}

	// Expected: scroll offset = 50 - (-5) = 55
	expectedOffset := 50 + 5
	if v2.ScrollOffset() != expectedOffset {
		t.Errorf("After drag up: offset = %d, want %d", v2.ScrollOffset(), expectedOffset)
	}
}

func TestViewport_DragScroll_Release(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Start drag at Y=5
	msg1 := tea.MouseMsg{
		X:      10,
		Y:      5,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionPress,
	}
	v, _ = v.Update(msg1)

	// Drag down to Y=10
	msg2 := tea.MouseMsg{
		X:      10,
		Y:      10,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionMotion,
	}
	v, _ = v.Update(msg2)

	currentOffset := v.ScrollOffset()

	// Release mouse button (end drag)
	msg3 := tea.MouseMsg{
		X:      10,
		Y:      10,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionRelease,
	}
	v2, cmd := v.Update(msg3)

	if cmd != nil {
		t.Error("Update should return nil Cmd")
	}

	// Offset should remain at last drag position
	if v2.ScrollOffset() != currentOffset {
		t.Errorf("After drag release: offset = %d, want %d", v2.ScrollOffset(), currentOffset)
	}

	// Further motion without drag should NOT affect scroll
	msg4 := tea.MouseMsg{
		X:      10,
		Y:      15,
		Button: tea.MouseButtonNone,
		Action: tea.MouseActionMotion,
	}
	v3, _ := v2.Update(msg4)

	if v3.ScrollOffset() != currentOffset {
		t.Error("Motion after drag release should not affect scroll")
	}
}

func TestViewport_DragScroll_BoundsTop(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(5)

	// Start drag at Y=10
	msg1 := tea.MouseMsg{
		X:      10,
		Y:      10,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionPress,
	}
	v, _ = v.Update(msg1)

	// Drag down by 20 (would result in offset = 5 - 20 = -15)
	// Should be clamped to 0
	msg2 := tea.MouseMsg{
		X:      10,
		Y:      30,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionMotion,
	}
	v2, _ := v.Update(msg2)

	if v2.ScrollOffset() != 0 {
		t.Errorf("Drag scroll past top: offset = %d, want 0 (clamped)", v2.ScrollOffset())
	}
}

func TestViewport_DragScroll_BoundsBottom(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(75) // Near bottom (max is 80)

	// Start drag at Y=10
	msg1 := tea.MouseMsg{
		X:      10,
		Y:      10,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionPress,
	}
	v, _ = v.Update(msg1)

	// Drag up by 20 (would result in offset = 75 + 20 = 95)
	// Should be clamped to 80 (100 lines - 20 visible)
	msg2 := tea.MouseMsg{
		X:      10,
		Y:      -10, // Negative Y (drag up a lot)
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionMotion,
	}
	v2, _ := v.Update(msg2)

	maxOffset := 100 - 20 // totalLines - height
	if v2.ScrollOffset() != maxOffset {
		t.Errorf("Drag scroll past bottom: offset = %d, want %d (clamped)", v2.ScrollOffset(), maxOffset)
	}
}

func TestViewport_DragScroll_SmallContent(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)
	v = v.SetLines([]string{"Line 1", "Line 2", "Line 3"}) // Only 3 lines

	// Start drag
	msg1 := tea.MouseMsg{
		X:      10,
		Y:      5,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionPress,
	}
	v, _ = v.Update(msg1)

	// Try to drag down
	msg2 := tea.MouseMsg{
		X:      10,
		Y:      10,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionMotion,
	}
	v2, _ := v.Update(msg2)

	// Offset should remain 0 (content fits entirely in viewport)
	if v2.ScrollOffset() != 0 {
		t.Errorf("Drag scroll with small content: offset = %d, want 0", v2.ScrollOffset())
	}
}

func TestViewport_DragScroll_EmptyContent(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)

	// Start drag
	msg1 := tea.MouseMsg{
		X:      10,
		Y:      5,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionPress,
	}
	v, _ = v.Update(msg1)

	// Try to drag
	msg2 := tea.MouseMsg{
		X:      10,
		Y:      10,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionMotion,
	}
	v2, _ := v.Update(msg2)

	// Should not panic, offset remains 0
	if v2.ScrollOffset() != 0 {
		t.Errorf("Drag scroll with empty content: offset = %d, want 0", v2.ScrollOffset())
	}
}

func TestViewport_DragScroll_LargeContent(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)
	// Create 10,000 lines
	lines := make([]string, 10000)
	for i := range lines {
		lines[i] = "Line content"
	}
	v = v.SetLines(lines)
	v = v.SetYOffset(5000)

	// Start drag at Y=10
	msg1 := tea.MouseMsg{
		X:      10,
		Y:      10,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionPress,
	}
	v, _ = v.Update(msg1)

	// Drag down by 100 cells
	msg2 := tea.MouseMsg{
		X:      10,
		Y:      110,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionMotion,
	}
	v2, _ := v.Update(msg2)

	// Expected: offset = 5000 - 100 = 4900
	expectedOffset := 5000 - 100
	if v2.ScrollOffset() != expectedOffset {
		t.Errorf("Drag scroll large content: offset = %d, want %d", v2.ScrollOffset(), expectedOffset)
	}
}

func TestViewport_DragScroll_Disabled(t *testing.T) {
	v := New(80, 20).MouseEnabled(false) // Mouse disabled
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Try to start drag
	msg1 := tea.MouseMsg{
		X:      10,
		Y:      5,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionPress,
	}
	v, _ = v.Update(msg1)

	// Try to drag
	msg2 := tea.MouseMsg{
		X:      10,
		Y:      10,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionMotion,
	}
	v2, _ := v.Update(msg2)

	// Offset should NOT change (mouse disabled)
	if v2.ScrollOffset() != 50 {
		t.Errorf("Drag scroll when disabled: offset = %d, want 50", v2.ScrollOffset())
	}
}

func TestViewport_DragScroll_MultipleDrags(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// First drag: Y=5 → Y=10 (scroll up by 5)
	v, _ = v.Update(tea.MouseMsg{X: 10, Y: 5, Button: tea.MouseButtonLeft, Action: tea.MouseActionPress})
	v, _ = v.Update(tea.MouseMsg{X: 10, Y: 10, Button: tea.MouseButtonLeft, Action: tea.MouseActionMotion})
	v, _ = v.Update(tea.MouseMsg{X: 10, Y: 10, Button: tea.MouseButtonLeft, Action: tea.MouseActionRelease})

	if v.ScrollOffset() != 45 {
		t.Errorf("After first drag: offset = %d, want 45", v.ScrollOffset())
	}

	// Second drag: Y=15 → Y=10 (scroll down by 5)
	v, _ = v.Update(tea.MouseMsg{X: 10, Y: 15, Button: tea.MouseButtonLeft, Action: tea.MouseActionPress})
	v, _ = v.Update(tea.MouseMsg{X: 10, Y: 10, Button: tea.MouseButtonLeft, Action: tea.MouseActionMotion})
	v, _ = v.Update(tea.MouseMsg{X: 10, Y: 10, Button: tea.MouseButtonLeft, Action: tea.MouseActionRelease})

	if v.ScrollOffset() != 50 {
		t.Errorf("After second drag: offset = %d, want 50", v.ScrollOffset())
	}
}

func TestViewport_DragScroll_RightButton(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Try to drag with right button (should not work, only left button)
	msg1 := tea.MouseMsg{
		X:      10,
		Y:      5,
		Button: tea.MouseButtonRight,
		Action: tea.MouseActionPress,
	}
	v, _ = v.Update(msg1)

	msg2 := tea.MouseMsg{
		X:      10,
		Y:      10,
		Button: tea.MouseButtonRight,
		Action: tea.MouseActionMotion,
	}
	v2, _ := v.Update(msg2)

	// Offset should NOT change (right button doesn't drag)
	if v2.ScrollOffset() != 50 {
		t.Errorf("Drag with right button: offset = %d, want 50", v2.ScrollOffset())
	}
}

func TestViewport_DragScroll_MiddleButton(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	// Try to drag with middle button (should not work, only left button)
	msg1 := tea.MouseMsg{
		X:      10,
		Y:      5,
		Button: tea.MouseButtonMiddle,
		Action: tea.MouseActionPress,
	}
	v, _ = v.Update(msg1)

	msg2 := tea.MouseMsg{
		X:      10,
		Y:      10,
		Button: tea.MouseButtonMiddle,
		Action: tea.MouseActionMotion,
	}
	v2, _ := v.Update(msg2)

	// Offset should NOT change (middle button doesn't drag)
	if v2.ScrollOffset() != 50 {
		t.Errorf("Drag with middle button: offset = %d, want 50", v2.ScrollOffset())
	}
}

func TestViewport_DragScroll_Immutability(t *testing.T) {
	v := New(80, 20).MouseEnabled(true)
	v = v.SetLines(make([]string, 100))
	v = v.SetYOffset(50)

	originalOffset := v.ScrollOffset()

	// Start drag on original
	msg1 := tea.MouseMsg{X: 10, Y: 5, Button: tea.MouseButtonLeft, Action: tea.MouseActionPress}
	v2, _ := v.Update(msg1)

	// Drag on v2
	msg2 := tea.MouseMsg{X: 10, Y: 10, Button: tea.MouseButtonLeft, Action: tea.MouseActionMotion}
	v3, _ := v2.Update(msg2)

	// Original should remain unchanged
	if v.ScrollOffset() != originalOffset {
		t.Error("Original viewport was mutated by drag operations")
	}

	// v3 should have changed offset
	if v3.ScrollOffset() == originalOffset {
		t.Error("Drag did not create new viewport with updated offset")
	}
}
