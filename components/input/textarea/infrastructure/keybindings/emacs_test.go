package keybindings

import (
	"testing"

	"github.com/phoenix-tui/phoenix/components/input/textarea/domain/model"
	"github.com/phoenix-tui/phoenix/tea/api"
)

// TestEmacsKeybindings_KeySpace verifies that Space key is handled correctly.
// This is a CRITICAL test - without KeySpace handling, spaces are ignored.
func TestEmacsKeybindings_KeySpace(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert space.
	msg := api.KeyMsg{Type: api.KeySpace}
	result, _ := handler.Handle(msg, ta)

	// Verify space was inserted.
	if result.Value() != " " {
		t.Errorf("Expected single space, got %q", result.Value())
	}

	// Insert another space.
	result, _ = handler.Handle(msg, result)

	// Verify two spaces.
	if result.Value() != "  " {
		t.Errorf("Expected two spaces, got %q", result.Value())
	}
}

// TestEmacsKeybindings_SpaceInText verifies space insertion within text.
func TestEmacsKeybindings_SpaceInText(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello".
	for _, r := range "hello" {
		msg := api.KeyMsg{Type: api.KeyRune, Rune: r}
		ta, _ = handler.Handle(msg, ta)
	}

	// Insert space.
	spaceMsg := api.KeyMsg{Type: api.KeySpace}
	ta, _ = handler.Handle(spaceMsg, ta)

	// Insert "world".
	for _, r := range "world" {
		msg := api.KeyMsg{Type: api.KeyRune, Rune: r}
		ta, _ = handler.Handle(msg, ta)
	}

	// Verify result.
	expected := "hello world"
	if ta.Value() != expected {
		t.Errorf("Expected %q, got %q", expected, ta.Value())
	}
}

// TestEmacsKeybindings_MultipleSpaces verifies multiple consecutive spaces.
func TestEmacsKeybindings_MultipleSpaces(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert 5 spaces.
	spaceMsg := api.KeyMsg{Type: api.KeySpace}
	for i := 0; i < 5; i++ {
		ta, _ = handler.Handle(spaceMsg, ta)
	}

	// Verify 5 spaces.
	expected := "     "
	if ta.Value() != expected {
		t.Errorf("Expected 5 spaces, got %q (len=%d)", ta.Value(), len(ta.Value()))
	}
}

// ============================================================================
// Ctrl Key Combinations - Navigation
// ============================================================================

func TestEmacsKeybindings_Ctrl_A_MoveToLineStart(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello world"
	insertText(t, handler, &ta, "hello world")

	// Ctrl+A should move to line start
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'a', Ctrl: true}
	result, _ := handler.Handle(msg, ta)

	// Cursor should be at position 0
	_, col := result.CursorPosition()
	if col != 0 {
		t.Errorf("Expected cursor at col 0, got %d", col)
	}
}

func TestEmacsKeybindings_Ctrl_E_MoveToLineEnd(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	insertText(t, handler, &ta, "hello")

	// Move to start
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'a', Ctrl: true}
	ta, _ = handler.Handle(msg, ta)

	// Ctrl+E should move to line end
	msg = api.KeyMsg{Type: api.KeyRune, Rune: 'e', Ctrl: true}
	result, _ := handler.Handle(msg, ta)

	// Cursor should be at end (col 5)
	_, col := result.CursorPosition()
	if col != 5 {
		t.Errorf("Expected cursor at col 5, got %d", col)
	}
}

func TestEmacsKeybindings_Ctrl_P_MoveUp(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert multiline text
	insertText(t, handler, &ta, "line1\nline2\nline3")

	// Ctrl+P should move up
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'p', Ctrl: true}
	result, _ := handler.Handle(msg, ta)

	// Should be on line 1 now (row 1)
	row, _ := result.CursorPosition()
	if row != 1 {
		t.Errorf("Expected row 1, got %d", row)
	}
}

func TestEmacsKeybindings_Ctrl_N_MoveDown(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert multiline text
	insertText(t, handler, &ta, "line1\nline2")

	// Move to start
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'a', Ctrl: true}
	ta, _ = handler.Handle(msg, ta)

	// Ctrl+N should move down
	msg = api.KeyMsg{Type: api.KeyRune, Rune: 'n', Ctrl: true}
	result, _ := handler.Handle(msg, ta)

	// Should be on line 1 now (row 1)
	row, _ := result.CursorPosition()
	if row != 1 {
		t.Errorf("Expected row 1, got %d", row)
	}
}

func TestEmacsKeybindings_Ctrl_F_MoveRight(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	insertText(t, handler, &ta, "hello")

	// Move to start
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'a', Ctrl: true}
	ta, _ = handler.Handle(msg, ta)

	// Ctrl+F should move right
	msg = api.KeyMsg{Type: api.KeyRune, Rune: 'f', Ctrl: true}
	result, _ := handler.Handle(msg, ta)

	// Should be at col 1
	_, col := result.CursorPosition()
	if col != 1 {
		t.Errorf("Expected col 1, got %d", col)
	}
}

func TestEmacsKeybindings_Ctrl_B_MoveLeft(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	insertText(t, handler, &ta, "hello")

	// Ctrl+B should move left
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'b', Ctrl: true}
	result, _ := handler.Handle(msg, ta)

	// Should be at col 4
	_, col := result.CursorPosition()
	if col != 4 {
		t.Errorf("Expected col 4, got %d", col)
	}
}

// ============================================================================
// Ctrl Key Combinations - Editing
// ============================================================================

func TestEmacsKeybindings_Ctrl_K_KillLine(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello world"
	insertText(t, handler, &ta, "hello world")

	// Move to middle (col 5)
	for i := 0; i < 6; i++ {
		msg := api.KeyMsg{Type: api.KeyRune, Rune: 'b', Ctrl: true}
		ta, _ = handler.Handle(msg, ta)
	}

	// Ctrl+K should kill from cursor to end of line
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'k', Ctrl: true}
	result, _ := handler.Handle(msg, ta)

	// Should have "hello" remaining
	if result.Value() != "hello" {
		t.Errorf("Expected \"hello\", got %q", result.Value())
	}
}

func TestEmacsKeybindings_Ctrl_U_KillToStart(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello world"
	insertText(t, handler, &ta, "hello world")

	// Ctrl+U should kill from start to cursor
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'u', Ctrl: true}
	result, _ := handler.Handle(msg, ta)

	// Should be empty (killed entire line)
	if result.Value() != "" {
		t.Errorf("Expected empty, got %q", result.Value())
	}
}

func TestEmacsKeybindings_Ctrl_W_KillWord(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello world"
	insertText(t, handler, &ta, "hello world")

	// Ctrl+W should kill word backward
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'w', Ctrl: true}
	result, _ := handler.Handle(msg, ta)

	// Should have "hello " remaining
	expected := "hello "
	if result.Value() != expected {
		t.Errorf("Expected %q, got %q", expected, result.Value())
	}
}

func TestEmacsKeybindings_Ctrl_Y_Yank(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello world"
	insertText(t, handler, &ta, "hello world")

	// Kill word
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'w', Ctrl: true}
	ta, _ = handler.Handle(msg, ta)

	// Yank (paste)
	msg = api.KeyMsg{Type: api.KeyRune, Rune: 'y', Ctrl: true}
	result, _ := handler.Handle(msg, ta)

	// Should have "hello world" again
	expected := "hello world"
	if result.Value() != expected {
		t.Errorf("Expected %q, got %q", expected, result.Value())
	}
}

func TestEmacsKeybindings_Ctrl_D_DeleteCharForward(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	insertText(t, handler, &ta, "hello")

	// Move to start
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'a', Ctrl: true}
	ta, _ = handler.Handle(msg, ta)

	// Ctrl+D should delete char forward
	msg = api.KeyMsg{Type: api.KeyRune, Rune: 'd', Ctrl: true}
	result, _ := handler.Handle(msg, ta)

	// Should have "ello"
	if result.Value() != "ello" {
		t.Errorf("Expected \"ello\", got %q", result.Value())
	}
}

func TestEmacsKeybindings_Ctrl_H_DeleteCharBackward(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	insertText(t, handler, &ta, "hello")

	// Ctrl+H should delete char backward
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'h', Ctrl: true}
	result, _ := handler.Handle(msg, ta)

	// Should have "hell"
	if result.Value() != "hell" {
		t.Errorf("Expected \"hell\", got %q", result.Value())
	}
}

func TestEmacsKeybindings_Ctrl_M_InsertNewline(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	insertText(t, handler, &ta, "hello")

	// Ctrl+M should insert newline
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'm', Ctrl: true}
	result, _ := handler.Handle(msg, ta)

	// Should have "hello\n"
	expected := "hello\n"
	if result.Value() != expected {
		t.Errorf("Expected %q, got %q", expected, result.Value())
	}
}

// ============================================================================
// Alt Key Combinations - Navigation
// ============================================================================

func TestEmacsKeybindings_Alt_F_ForwardWord(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello world"
	insertText(t, handler, &ta, "hello world")

	// Move to start
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'a', Ctrl: true}
	ta, _ = handler.Handle(msg, ta)

	// Alt+F should move forward word
	msg = api.KeyMsg{Type: api.KeyRune, Rune: 'f', Alt: true}
	result, _ := handler.Handle(msg, ta)

	// Should be after "hello" (col 5)
	_, col := result.CursorPosition()
	if col < 5 {
		t.Errorf("Expected cursor at col >= 5, got %d", col)
	}
}

func TestEmacsKeybindings_Alt_B_BackwardWord(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello world"
	insertText(t, handler, &ta, "hello world")

	// Alt+B should move backward word
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'b', Alt: true}
	result, _ := handler.Handle(msg, ta)

	// Should be at start of "world" (col 6)
	_, col := result.CursorPosition()
	if col != 6 {
		t.Errorf("Expected cursor at col 6, got %d", col)
	}
}

func TestEmacsKeybindings_Alt_LessThan_MoveToBufferStart(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert multiline text
	insertText(t, handler, &ta, "line1\nline2\nline3")

	// Alt+< should move to buffer start
	msg := api.KeyMsg{Type: api.KeyRune, Rune: '<', Alt: true}
	result, _ := handler.Handle(msg, ta)

	// Should be at row 0, col 0
	row, col := result.CursorPosition()
	if row != 0 || col != 0 {
		t.Errorf("Expected (0,0), got (%d,%d)", row, col)
	}
}

func TestEmacsKeybindings_Alt_GreaterThan_MoveToBufferEnd(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert multiline text
	insertText(t, handler, &ta, "line1\nline2\nline3")

	// Move to start
	msg := api.KeyMsg{Type: api.KeyRune, Rune: '<', Alt: true}
	ta, _ = handler.Handle(msg, ta)

	// Alt+> should move to buffer end
	msg = api.KeyMsg{Type: api.KeyRune, Rune: '>', Alt: true}
	result, _ := handler.Handle(msg, ta)

	// Should be at last line
	row, _ := result.CursorPosition()
	if row != 2 {
		t.Errorf("Expected row 2, got %d", row)
	}
}

func TestEmacsKeybindings_Alt_D_KillWord(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello world"
	insertText(t, handler, &ta, "hello world")

	// Move to start
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'a', Ctrl: true}
	ta, _ = handler.Handle(msg, ta)

	// Alt+D should kill word forward
	msg = api.KeyMsg{Type: api.KeyRune, Rune: 'd', Alt: true}
	result, _ := handler.Handle(msg, ta)

	// Should have " world" remaining
	expected := " world"
	if result.Value() != expected {
		t.Errorf("Expected %q, got %q", expected, result.Value())
	}
}

func TestEmacsKeybindings_Alt_Backspace_KillWord(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello world"
	insertText(t, handler, &ta, "hello world")

	// Alt+Backspace should kill word backward
	msg := api.KeyMsg{Type: api.KeyBackspace, Alt: true}
	result, _ := handler.Handle(msg, ta)

	// Should have "hello " remaining
	expected := "hello "
	if result.Value() != expected {
		t.Errorf("Expected %q, got %q", expected, result.Value())
	}
}

// ============================================================================
// Special Keys (no modifiers)
// ============================================================================

func TestEmacsKeybindings_ArrowUp(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert multiline
	insertText(t, handler, &ta, "line1\nline2")

	// Arrow Up
	msg := api.KeyMsg{Type: api.KeyUp}
	result, _ := handler.Handle(msg, ta)

	// Should be on line 0
	row, _ := result.CursorPosition()
	if row != 0 {
		t.Errorf("Expected row 0, got %d", row)
	}
}

func TestEmacsKeybindings_ArrowDown(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert multiline
	insertText(t, handler, &ta, "line1\nline2")

	// Move to start
	msg := api.KeyMsg{Type: api.KeyRune, Rune: '<', Alt: true}
	ta, _ = handler.Handle(msg, ta)

	// Arrow Down
	msg = api.KeyMsg{Type: api.KeyDown}
	result, _ := handler.Handle(msg, ta)

	// Should be on line 1
	row, _ := result.CursorPosition()
	if row != 1 {
		t.Errorf("Expected row 1, got %d", row)
	}
}

func TestEmacsKeybindings_ArrowLeft(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	insertText(t, handler, &ta, "hello")

	// Arrow Left
	msg := api.KeyMsg{Type: api.KeyLeft}
	result, _ := handler.Handle(msg, ta)

	// Should be at col 4
	_, col := result.CursorPosition()
	if col != 4 {
		t.Errorf("Expected col 4, got %d", col)
	}
}

func TestEmacsKeybindings_ArrowRight(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	insertText(t, handler, &ta, "hello")

	// Move to start
	msg := api.KeyMsg{Type: api.KeyHome}
	ta, _ = handler.Handle(msg, ta)

	// Arrow Right
	msg = api.KeyMsg{Type: api.KeyRight}
	result, _ := handler.Handle(msg, ta)

	// Should be at col 1
	_, col := result.CursorPosition()
	if col != 1 {
		t.Errorf("Expected col 1, got %d", col)
	}
}

func TestEmacsKeybindings_Home(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	insertText(t, handler, &ta, "hello")

	// Home
	msg := api.KeyMsg{Type: api.KeyHome}
	result, _ := handler.Handle(msg, ta)

	// Should be at col 0
	_, col := result.CursorPosition()
	if col != 0 {
		t.Errorf("Expected col 0, got %d", col)
	}
}

func TestEmacsKeybindings_End(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	insertText(t, handler, &ta, "hello")

	// Move to start
	msg := api.KeyMsg{Type: api.KeyHome}
	ta, _ = handler.Handle(msg, ta)

	// End
	msg = api.KeyMsg{Type: api.KeyEnd}
	result, _ := handler.Handle(msg, ta)

	// Should be at col 5
	_, col := result.CursorPosition()
	if col != 5 {
		t.Errorf("Expected col 5, got %d", col)
	}
}

func TestEmacsKeybindings_Backspace(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	insertText(t, handler, &ta, "hello")

	// Backspace
	msg := api.KeyMsg{Type: api.KeyBackspace}
	result, _ := handler.Handle(msg, ta)

	// Should have "hell"
	if result.Value() != "hell" {
		t.Errorf("Expected \"hell\", got %q", result.Value())
	}
}

func TestEmacsKeybindings_Delete(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	insertText(t, handler, &ta, "hello")

	// Move to start
	msg := api.KeyMsg{Type: api.KeyHome}
	ta, _ = handler.Handle(msg, ta)

	// Delete
	msg = api.KeyMsg{Type: api.KeyDelete}
	result, _ := handler.Handle(msg, ta)

	// Should have "ello"
	if result.Value() != "ello" {
		t.Errorf("Expected \"ello\", got %q", result.Value())
	}
}

func TestEmacsKeybindings_Enter(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	insertText(t, handler, &ta, "hello")

	// Enter
	msg := api.KeyMsg{Type: api.KeyEnter}
	result, _ := handler.Handle(msg, ta)

	// Should have "hello\n"
	expected := "hello\n"
	if result.Value() != expected {
		t.Errorf("Expected %q, got %q", expected, result.Value())
	}
}

func TestEmacsKeybindings_Rune(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert 'a'
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'a'}
	result, _ := handler.Handle(msg, ta)

	// Should have "a"
	if result.Value() != "a" {
		t.Errorf("Expected \"a\", got %q", result.Value())
	}
}

// ============================================================================
// Unhandled Keys
// ============================================================================

func TestEmacsKeybindings_UnhandledKey_NoChange(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	insertText(t, handler, &ta, "hello")

	// Unhandled key (e.g., F1)
	msg := api.KeyMsg{Type: api.KeyType(999)} // Invalid key type
	result, _ := handler.Handle(msg, ta)

	// Should remain unchanged
	if result.Value() != "hello" {
		t.Errorf("Expected \"hello\", got %q", result.Value())
	}
}

// ============================================================================
// Case Insensitivity Tests (Ctrl+A == Ctrl+Shift+A)
// ============================================================================

func TestEmacsKeybindings_Ctrl_UppercaseWorks(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello"
	insertText(t, handler, &ta, "hello")

	// Ctrl+Shift+A (uppercase) should also move to line start
	msg := api.KeyMsg{Type: api.KeyRune, Rune: 'A', Ctrl: true}
	result, _ := handler.Handle(msg, ta)

	// Should be at col 0
	_, col := result.CursorPosition()
	if col != 0 {
		t.Errorf("Expected col 0, got %d", col)
	}
}

func TestEmacsKeybindings_Alt_UppercaseWorks(t *testing.T) {
	ta := model.NewTextArea()
	handler := NewEmacsKeybindings()

	// Insert "hello world"
	insertText(t, handler, &ta, "hello world")

	// Move to start
	msg := api.KeyMsg{Type: api.KeyHome}
	ta, _ = handler.Handle(msg, ta)

	// Alt+Shift+F (uppercase) should also move forward word
	msg = api.KeyMsg{Type: api.KeyRune, Rune: 'F', Alt: true}
	result, _ := handler.Handle(msg, ta)

	// Should be after "hello"
	_, col := result.CursorPosition()
	if col < 5 {
		t.Errorf("Expected col >= 5, got %d", col)
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

func insertText(t *testing.T, handler *EmacsKeybindings, ta **model.TextArea, text string) {
	t.Helper()
	for _, r := range text {
		var msg api.KeyMsg
		if r == '\n' {
			msg = api.KeyMsg{Type: api.KeyEnter}
		} else {
			msg = api.KeyMsg{Type: api.KeyRune, Rune: r}
		}
		*ta, _ = handler.Handle(msg, *ta)
	}
}
