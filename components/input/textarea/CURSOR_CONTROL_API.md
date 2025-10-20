# TextArea Cursor Control API

**Status**: ✅ Released in v0.1.0-beta.2 (2025-10-20)
**Priority**: HIGH - Enables GoSh shell REPL development
**Specification**: docs/dev/TEXTAREA_CURSOR_CONTROL_IMPLEMENTATION.md

---

## Overview

Phoenix TextArea now supports **programmatic cursor control** and **observer pattern** for cursor movements. This enables shell-like applications (GoSh) to:

1. Set cursor position programmatically
2. Validate cursor movements (boundary protection)
3. React to cursor movements (observers)
4. Provide user feedback when boundaries hit

**100% Backward Compatible** - All features are opt-in!

---

## Features

### 1. SetCursorPosition() - Programmatic Cursor Control

Set cursor to specific position with automatic bounds checking.

```go
ta := textarea.New().SetValue("line1\nline2\nline3")
ta = ta.SetCursorPosition(1, 2) // Row 1, Col 2

row, col := ta.CursorPosition()
// row = 1, col = 2
```

**Automatic clamping:**
- Row clamped to [0, lineCount-1]
- Col clamped to [0, len(line in runes)]
- Negative values → 0
- Out of bounds → max valid value

---

### 2. OnMovement() - Movement Validator

Validate cursor movements BEFORE they happen. Return `false` to block movement.

```go
ta := textarea.New().
    SetValue("> command").
    SetCursorPosition(0, 2).
    OnMovement(func(from, to textarea.CursorPos) bool {
        // Don't allow cursor before prompt ("> ")
        if to.Row == 0 && to.Col < 2 {
            return false // Block movement
        }
        return true // Allow movement
    })
```

**When validator is called:**
- ALL arrow key movements (←↑→↓)
- Home/End, Ctrl+A/Ctrl+E
- PgUp/PgDn
- Word navigation (Alt+F/Alt+B)
- Buffer start/end (Alt+</Alt+>)

**NOT called for:**
- Typing characters (use EditingService for that)
- SetCursorPosition() (explicit API calls bypass validator)

---

### 3. OnCursorMoved() - Cursor Movement Observer

React to cursor movements AFTER they happen. Cannot block movement.

```go
ta := textarea.New().
    OnCursorMoved(func(from, to textarea.CursorPos) {
        if from.Row != to.Row {
            // Cursor moved to different line
            refreshSyntaxHighlight(to.Row)
        }
    })
```

**When observer is called:**
- After successful cursor movement
- Only when movement actually happened (from != to)
- NOT called when movement is blocked by validator

**Use cases:**
- Update syntax highlighting
- Track cursor history
- Update UI state
- Log movements for debugging

---

### 4. OnBoundaryHit() - Boundary Hit Feedback

Provide feedback when movement is blocked by validator.

```go
ta := textarea.New().
    OnMovement(func(from, to textarea.CursorPos) bool {
        return to.Row >= 0 && to.Col >= 0 // Block negative positions
    }).
    OnBoundaryHit(func(attemptedPos textarea.CursorPos, reason string) {
        // Flash screen, beep, or show message
        fmt.Println("Cannot move to", attemptedPos, ":", reason)
    })
```

**Reasons provided:**
- `"movement blocked by validator"` - when validator returns false
- `"already at top"` - when trying to move up from row 0
- `"already at bottom"` - when trying to move down from last row

**Use cases:**
- Visual feedback (flash, beep)
- Error messages
- Accessibility
- User experience

---

## Types

### CursorPos

```go
type CursorPos struct {
    Row int // Line number (0-based)
    Col int // Column number (0-based, rune offset)
}
```

### Callback Types

```go
// Validator - return true to allow, false to block
type MovementValidator func(from, to CursorPos) bool

// Observer - called after successful movement
type CursorMovedHandler func(from, to CursorPos)

// Feedback - called when movement blocked
type BoundaryHitHandler func(attemptedPos CursorPos, reason string)
```

---

## Examples

### Example 1: Shell Prompt Protection (GoSh Use Case)

```go
ta := textarea.New().
    SetValue("> ").
    SetCursorPosition(0, 2).
    OnMovement(func(from, to textarea.CursorPos) bool {
        // Protect prompt area (columns 0-1)
        if to.Row == 0 && to.Col < 2 {
            return false
        }
        return true
    }).
    OnBoundaryHit(func(attemptedPos textarea.CursorPos, reason string) {
        flash("Cannot edit prompt area")
    })

// User types "ls -la"
for _, ch := range "ls -la" {
    ta, _ = ta.Update(textarea.KeyMsg{Type: textarea.KeyRune, Rune: ch})
}

// User presses Home (Ctrl+A) - moves to column 2, not 0
ta, _ = ta.Update(textarea.KeyMsg{Type: textarea.KeyRune, Rune: 'a', Ctrl: true})

row, col := ta.CursorPosition()
// row = 0, col = 2 (blocked at prompt boundary!)
```

### Example 2: Read-Only History Lines

```go
ta := textarea.New().
    SetValue("> pwd\n/home/user\n> ls\nfile.txt\n> ").
    SetCursorPosition(4, 2). // Current prompt (last line)
    OnMovement(func(from, to textarea.CursorPos) bool {
        // Only allow editing current line (row 4)
        if to.Row < 4 {
            return false // Block access to history
        }
        if to.Row == 4 && to.Col < 2 {
            return false // Block moving before prompt
        }
        return true
    }).
    OnBoundaryHit(func(attemptedPos textarea.CursorPos, reason string) {
        showMessage("Cannot edit command history")
    })

// Arrow up is blocked (can't edit history)
ta, _ = ta.Update(textarea.KeyMsg{Type: textarea.KeyUp})

row, col := ta.CursorPosition()
// row = 4, col = 2 (stayed at current prompt)
```

### Example 3: Syntax Highlighting on Row Change

```go
ta := textarea.New().
    SetValue("package main\n\nfunc main() {\n}\n").
    OnCursorMoved(func(from, to textarea.CursorPos) {
        if from.Row != to.Row {
            // Cursor moved to different line - refresh syntax highlighting
            rehighlightLine(to.Row)
        }
    })
```

### Example 4: All Features Together

```go
ta := textarea.New().
    SetValue("> command here").
    SetCursorPosition(0, 2).

    // Validator - protect prompt
    OnMovement(func(from, to textarea.CursorPos) bool {
        if to.Row == 0 && to.Col < 2 {
            return false
        }
        return true
    }).

    // Observer - track movements
    OnCursorMoved(func(from, to textarea.CursorPos) {
        logMovement(from, to)
    }).

    // Feedback - notify user
    OnBoundaryHit(func(attemptedPos textarea.CursorPos, reason string) {
        flash("Boundary hit!")
    })
```

---

## Architecture

### DDD Layered Implementation

```
API Layer (components/input/textarea/api/)
├── CursorPos (API type)
├── MovementValidator
├── CursorMovedHandler
├── BoundaryHitHandler
├── OnMovement() - converts API → Domain types
├── OnCursorMoved() - converts API → Domain types
├── OnBoundaryHit() - converts API → Domain types
└── SetCursorPosition() - calls domain method

Domain Layer (components/input/textarea/domain/model/)
├── CursorPos (domain type)
├── TextArea struct (holds callbacks)
├── SetCursorPosition() - clamping logic
├── GetMovementValidator()
├── GetCursorMovedHandler()
├── GetBoundaryHitHandler()
├── WithMovementValidator() - immutable setter
├── WithCursorMovedHandler() - immutable setter
└── WithBoundaryHitHandler() - immutable setter

Service Layer (components/input/textarea/domain/service/)
└── NavigationService
    ├── MoveLeft() - checks validator, fires callbacks
    ├── MoveRight() - checks validator, fires callbacks
    ├── MoveUp() - checks validator, fires callbacks
    ├── MoveDown() - checks validator, fires callbacks
    ├── MoveToLineStart() - checks validator, fires callbacks
    ├── MoveToLineEnd() - checks validator, fires callbacks
    ├── ForwardWord() - checks validator, fires callbacks
    ├── BackwardWord() - checks validator, fires callbacks
    └── moveCursor() - helper to reduce duplication
```

### Flow

1. User presses arrow key
2. EmacsKeybindings calls NavigationService.MoveLeft()
3. NavigationService:
   - Calculates new position (to)
   - Checks validator (if set) - BEFORE movement
   - If blocked: fires OnBoundaryHit, returns unchanged
   - If allowed: applies movement, fires OnCursorMoved, returns updated
4. Updated TextArea returned to API layer

---

## Testing

### Unit Tests (textarea_cursor_control_test.go)

- ✅ SetCursorPosition clamping (negative, out of bounds, empty buffer)
- ✅ OnMovement blocks invalid movements
- ✅ OnMovement allows valid movements
- ✅ OnCursorMoved called after successful movement
- ✅ OnCursorMoved NOT called when blocked
- ✅ OnBoundaryHit called when blocked
- ✅ OnBoundaryHit NOT called when allowed
- ✅ All 4 features working together
- ✅ Backward compatibility (no callbacks)
- ✅ Immutability (SetCursorPosition returns new instance)

### Integration Tests (textarea_shell_integration_test.go)

- ✅ Shell boundary protection (prevent cursor before prompt)
- ✅ Multiple lines (command history)
- ✅ User feedback on boundary hit
- ✅ Typing still allowed
- ✅ Dynamic prompt (continuation lines)
- ✅ Real-world shell scenario
- ✅ Backspace at prompt boundary

### Example (examples/shell_prompt/)

- ✅ Interactive shell demo
- ✅ Demonstrates all 4 features
- ✅ Shows boundary protection in action

---

## API Reference

### SetCursorPosition

```go
func (t TextArea) SetCursorPosition(row, col int) TextArea
```

Sets cursor to specific position with bounds checking. Position is clamped to valid range.

**Returns**: New TextArea instance (immutable)

**Example**:
```go
ta = ta.SetCursorPosition(1, 5)
```

---

### OnMovement

```go
func (t TextArea) OnMovement(validator MovementValidator) TextArea
```

Sets validator that is called BEFORE cursor movements. Return `false` to block movement.

**Returns**: New TextArea instance (immutable)

**Example**:
```go
ta = ta.OnMovement(func(from, to CursorPos) bool {
    return to.Row >= 0 && to.Col >= 0
})
```

---

### OnCursorMoved

```go
func (t TextArea) OnCursorMoved(handler CursorMovedHandler) TextArea
```

Sets observer that is called AFTER successful cursor movement. Cannot block movement.

**Returns**: New TextArea instance (immutable)

**Example**:
```go
ta = ta.OnCursorMoved(func(from, to CursorPos) {
    fmt.Printf("Moved from (%d,%d) to (%d,%d)\n", from.Row, from.Col, to.Row, to.Col)
})
```

---

### OnBoundaryHit

```go
func (t TextArea) OnBoundaryHit(handler BoundaryHitHandler) TextArea
```

Sets handler that is called when cursor movement is blocked. Provides user feedback.

**Returns**: New TextArea instance (immutable)

**Example**:
```go
ta = ta.OnBoundaryHit(func(attemptedPos CursorPos, reason string) {
    flash("Cannot move: " + reason)
})
```

---

## Performance

All callback checks are O(1):
- nil check for each callback
- Single function call if set
- No overhead if callbacks not used (100% backward compatible)

**Benchmarks** (TODO - add after implementation):
- Navigation with callbacks: < 1μs overhead
- Navigation without callbacks: 0μs overhead (same as before)

---

## Backward Compatibility

✅ **100% Backward Compatible**

All features are opt-in:
- Existing code works unchanged
- No callbacks = no overhead
- No breaking changes to existing API

```go
// Old code (still works!)
ta := textarea.New().SetValue("test")
ta, _ = ta.Update(textarea.KeyMsg{Type: textarea.KeyRight})
// Works exactly as before
```

---

## Future Enhancements (Out of Scope for v0.2.0)

- SetBounds() - declarative constraints (e.g., "row >= 5 && col >= 2")
- Mouse event validation
- Paste event validation
- Multi-cursor support
- Custom boundary hit animations

---

## FAQ

### Q: Does OnMovement() validate typing?

**A**: No. OnMovement() only validates cursor navigation (arrow keys, Home/End, etc.). Typing is handled by EditingService. If you need to prevent typing in certain areas, use EditingService validation.

### Q: Can I change validator dynamically?

**A**: Yes! Just call OnMovement() again with new validator. TextArea is immutable, so you get a new instance.

```go
ta = ta.OnMovement(newValidator)
```

### Q: Does SetCursorPosition() trigger validator?

**A**: No. Explicit API calls bypass validator (you're in control). Validator only applies to user navigation.

### Q: What if I set both validator and observer?

**A**: They work together:
1. Validator checked BEFORE movement
2. If blocked: OnBoundaryHit fired, OnCursorMoved NOT fired
3. If allowed: movement applied, OnCursorMoved fired

### Q: Can I have multiple validators?

**A**: Not directly. But you can chain them:

```go
ta = ta.OnMovement(func(from, to CursorPos) bool {
    return validator1(from, to) && validator2(from, to)
})
```

### Q: Does this work with Vi keybindings?

**A**: Yes! Once Vi keybindings are implemented (future), they'll use the same NavigationService, so validation will work automatically.

---

## Credits

**Implementation**: Claude Code (phoenix-tui-architect agent)
**Specification**: Andy (GoSh shell requirements)
**Date**: 2025-01-19
**Version**: Phoenix v0.2.0

---

**Status**: ✅ IMPLEMENTED - Ready for GoSh integration!
