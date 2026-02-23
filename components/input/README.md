# Phoenix TextInput Component

**Universal text input component** for Phoenix TUI Framework. Designed for any application (shell, editor, form, chat), not just shell-specific use cases.

## Key Features

- **Grapheme-aware cursor movement** - Proper handling of emoji, CJK, combining characters
- **Horizontal scrolling** - Long input support for narrow terminals
- **Selection support** - Visual highlighting for copy/paste operations
- **Public cursor API** ⭐ **KEY DIFFERENTIATOR** - Fine-grained cursor control
- **Validation hooks** - Custom validation with clear error states
- **Extensible keybindings** - Add your own key handlers
- **Immutable design** - All operations return new instances
- **High test coverage** - Extensive domain and API coverage

## Installation

```bash
go get github.com/phoenix-tui/phoenix/components/input
```

## Quick Start

```go
package main

import (
    "github.com/phoenix-tui/phoenix/components/input/api"
    "github.com/phoenix-tui/phoenix/tea"
)

func main() {
    input := input.New(40).
        Placeholder("Enter your name...").
        Focused(true)

    p := tea.NewProgram(input)
    p.Run()
}
```

## API Reference

### Constructor

```go
// New creates a TextInput with specified visible width
input := input.New(40)
```

### Fluent Configuration

```go
input.Placeholder("Enter text...")        // Set placeholder text
input.Content("initial")                  // Set initial content
input.Focused(true)                       // Set focus state
input.Width(80)                           // Set visible width
input.Validator(func(s string) error {...}) // Set validation function
input.KeyBindings(customHandler)          // Set custom key handler
```

### Public Cursor API ⭐ **KEY DIFFERENTIATOR**

```go
// Get cursor position (grapheme offset)
pos := input.CursorPosition()

// Split content around cursor (for custom rendering)
before, at, after := input.ContentParts()

// Set content and cursor atomically (race-free)
input = input.SetContent("new text", 5)
```

**Why this matters:**
- **Custom cursor rendering** - gosh uses this for shell prompt styling
- **Syntax highlighting** - Highlight around cursor position
- **Autocomplete** - Insert completions at exact cursor location
- **History navigation** - Atomic content+cursor updates prevent races
- **Multi-line editing** - Split content for line-aware operations

### Accessors

```go
input.Value()           // Get current content
input.IsValid()         // Check validation status
input.IsFocused()       // Get focus state
```

### tea.Model Implementation

```go
input.Init()                     // Initialize (returns nil)
input.Update(msg tea.Msg)        // Handle messages
input.View()                     // Render to string
```

## Built-in Keybindings

| Key | Action |
|-----|--------|
| Left/Right Arrow | Move cursor by grapheme |
| Home / Ctrl-A | Move to start |
| End / Ctrl-E | Move to end |
| Backspace | Delete before cursor |
| Delete | Delete after cursor |
| Ctrl-U | Clear all content |
| Ctrl-A (string) | Select all |
| Printable chars | Insert at cursor |

## Validation

### Built-in Validators

```go
import "github.com/phoenix-tui/phoenix/components/input/api"

// Common validators
input.NotEmpty()                    // Content cannot be empty
input.MinLength(5)                  // Minimum length
input.MaxLength(100)                // Maximum length
input.Range(5, 100)                 // Length range
input.Chain(validator1, validator2) // Multiple validators
```

### Custom Validators

```go
emailValidator := func(s string) error {
    if !strings.Contains(s, "@") {
        return errors.New("must be valid email")
    }
    return nil
}

input := input.New(40).Validator(emailValidator)

// Check validation
if !input.IsValid() {
    // Show error
}
```

### Validation Errors

```go
import "github.com/phoenix-tui/phoenix/components/input/api"

// Exported error types
input.ErrEmpty          // Content cannot be empty
input.ErrTooShort       // Content is too short
input.ErrTooLong        // Content is too long
input.ErrInvalidFormat  // Content has invalid format
```

## Custom Keybindings

```go
// Create custom handler
customHandler := func(input *model.TextInput, msg tea.KeyMsg) *model.TextInput {
    if msg.String() == "ctrl+d" {
        // Duplicate content
        content := input.Content()
        return input.WithContent(content + content)
    }
    return nil // Not handled, fall through to defaults
}

// Apply custom bindings
input := input.New(40).KeyBindings(
    input.CustomKeyBindings(customHandler),
)
```

**Note:** Custom handlers are tried first. Return `nil` to fall through to default bindings.

## Unicode Handling

TextInput is **grapheme-aware** using `github.com/rivo/uniseg`:

```go
// All these work correctly:
input.Content("Hello 世界 👋")        // CJK + emoji
input.Content("👨‍👩‍👧‍👦 Family")       // Complex emoji (single grapheme)
input.Content("café")                // Combining characters (é = e + ́)
input.Content("🇺🇸 Flag")             // Flag emoji (single grapheme)
```

Cursor movement operates on **grapheme clusters**, not bytes or runes:
- Moving left from "👨‍👩‍👧‍👦" moves over the entire family emoji (not individual parts)
- Emoji with skin tones (👋🏽) are treated as single units
- CJK characters are handled correctly

## Examples

### 1. Basic Input

```go
input := input.New(40).
    Placeholder("Enter your name...").
    Focused(true)
```

See: `examples/basic/main.go`

### 2. Validated Input

```go
emailValidator := func(s string) error {
    if !regexp.MustCompile(`^.+@.+\..+$`).MatchString(s) {
        return errors.New("invalid email")
    }
    return nil
}

input := input.New(50).
    Placeholder("user@example.com").
    Validator(emailValidator).
    Focused(true)

// Check validation
if input.IsValid() {
    // Process valid email
}
```

See: `examples/validated/main.go`

### 3. Styled Input

```go
// Custom rendering with borders
func renderWithBorder(input *input.Input, focused bool) string {
    border := "─"
    if focused {
        border = "═"
    }

    return fmt.Sprintf(
        "┌%s┐\n│ %s │\n└%s┘",
        strings.Repeat(border, 40),
        input.View(),
        strings.Repeat(border, 40),
    )
}
```

See: `examples/styled/main.go`

### 4. Cursor API Demo ⭐

```go
// Get cursor information
pos := input.CursorPosition()
before, at, after := input.ContentParts()

fmt.Printf("Cursor at %d: before=%q, at=%q, after=%q\n", pos, before, at, after)

// Set content and cursor atomically
input = input.SetContent("Hello World", 6) // cursor after "Hello "

// Custom cursor rendering
before, at, after := input.ContentParts()
customRender := before + "[" + at + "]" + after  // Render cursor as brackets
```

See: `examples/cursor_api/main.go`

## How gosh Will Use TextInput

TextInput is **universal** - gosh will wrap it in `ShellInput` with shell-specific features:

```go
// gosh/ui/components/shell_input.go
type ShellInput struct {
    input       *input.Input       // Phoenix TextInput
    history     *History           // Command history
    completer   *Completer         // Tab completion
    emacs       *EmacsBindings     // Emacs keybindings (Ctrl-A, Ctrl-E, etc.)
}

func (s *ShellInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up":
            // History previous (uses SetContent API)
            cmd := s.history.Previous()
            s.input = s.input.SetContent(cmd, len(cmd))

        case "tab":
            // Autocomplete at cursor
            pos := s.input.CursorPosition()
            before, _, _ := s.input.ContentParts()
            completion := s.completer.Complete(before)
            s.input = s.input.SetContent(before+completion, pos+len(completion))

        // ... custom Emacs bindings ...
        }
    }

    // Delegate to base input
    return s.input.Update(msg)
}
```

**Benefits of this approach:**
- TextInput stays universal (useful for forms, editors, chats)
- gosh adds shell-specific behavior via composition
- Other apps can add their own wrappers
- Public cursor API enables all shell features

## Architecture (DDD Layers)

```
input/
├── domain/              # Pure business logic
│   ├── model/          # TextInput aggregate root
│   ├── value/          # Cursor, Selection value objects
│   └── service/        # CursorMovement, Validation services
├── infrastructure/      # Technical implementation
│   └── keybindings.go  # Default and custom key handlers
├── api/                # Public interface
│   └── input.go        # Fluent API, tea.Model integration
└── examples/           # Usage demonstrations
    ├── basic.go
    ├── validated.go
    ├── styled.go
    └── cursor_api.go
```

**Key architectural decisions:**
- **Immutable operations** - All `With*()` and `Move*()` methods return new instances
- **Grapheme-aware** - Uses `github.com/rivo/uniseg` for proper Unicode handling
- **Service injection** - CursorMovementService, ValidationService injected into domain
- **Strategy pattern** - KeyBindingHandler interface for extensible key handling

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# High coverage across all layers
```

## Comparison with Bubbles textinput

| Feature | Phoenix TextInput | Bubbles textinput |
|---------|------------------|-------------------|
| Unicode handling | ✅ Grapheme-aware | ⚠️ Rune-based (emoji issues) |
| Public cursor API | ✅ **CursorPosition(), ContentParts(), SetContent()** | ❌ Private cursor field |
| Immutability | ✅ All operations return new instance | ✅ Similar |
| Selection | ✅ Built-in | ❌ Not supported |
| Validation | ✅ Built-in hooks | ⚠️ Manual check |
| Custom keybindings | ✅ KeyBindingHandler interface | ⚠️ Override entire Update() |
| Scrolling | ✅ Horizontal scrolling | ✅ Similar |
| Architecture | ✅ DDD with clear layers | ⚠️ Monolithic |
| Test coverage | ✅ High | ⚠️ Lower |

**Key differentiator:** Public cursor API enables applications to:
- Customize cursor rendering (shell prompts, syntax highlighting)
- Implement atomic content+cursor updates (history navigation)
- Build advanced features (autocomplete, multi-line, syntax aware)

## Migration from Bubbles

```go
// Bubbles textinput
import "github.com/charmbracelet/bubbles/textinput"

ti := textinput.New()
ti.Placeholder = "Enter text..."
ti.Focus()
value := ti.Value()

// Phoenix TextInput
import "github.com/phoenix-tui/phoenix/components/input/api"

input := input.New(40).
    Placeholder("Enter text...").
    Focused(true)
value := input.Value()
```

**Major differences:**
- Fluent API instead of field assignment
- `New(width)` instead of width calculation
- `Focused(bool)` instead of `Focus()/Blur()`
- `Value()` instead of `Value()` (same!)

## Future Enhancements

- **Styling integration** - Full phoenix/style support
- **Multi-line mode** - Textarea variant
- **Password mode** - Masked input
- **Input masks** - Format-aware input (phone numbers, dates)
- **Suggestions** - Dropdown completion
- **Undo/Redo** - History stack

## Contributing

TextInput is part of Phoenix TUI Framework. See main project README for contribution guidelines.

## License

Apache-2.0 (same as Phoenix TUI Framework)

---

**Examples:** 4 complete (basic, validated, styled, cursor_api)
