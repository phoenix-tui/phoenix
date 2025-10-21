package viewport

import (
	"reflect"
	"testing"

	tea "github.com/phoenix-tui/phoenix/tea/api"
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

	// Should scroll down by 3 lines
	expectedOffset := 10 + 3
	if v2.ScrollOffset() != expectedOffset {
		t.Errorf("After mouse wheel down: offset = %d, want %d", v2.ScrollOffset(), expectedOffset)
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
