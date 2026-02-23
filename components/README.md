# phoenix/components - UI Component Library

Rich, reusable TUI components for Phoenix Framework with DDD architecture, type safety, and comprehensive testing.

**Module**: `github.com/phoenix-tui/phoenix/components`
**License**: MIT (planned)

---

## What is phoenix/components?

Phoenix Components is a rich library of terminal UI widgets built on top of Phoenix/tea (Elm Architecture). Each component follows Domain-Driven Design principles with clear separation between business logic, presentation, and infrastructure.

**Features**:
- **7 Production-Ready Components** - Input, TextArea, List, Viewport, Table, Modal, Progress
- **DDD Architecture** - Domain layer with high test coverage
- **Type-Safe** - Generic constraints for compile-time safety
- **Fluent API** - Chainable method calls for easy styling
- **Unicode-Aware** - Perfect CJK and emoji support
- **Extensive Test Suite** - High coverage across all components
- **Zero External TUI Dependencies** - Built from scratch

---

## Components Overview

### 1. Input - Single-Line Text Input
**Module**: `github.com/phoenix-tui/phoenix/components/input`

Single-line text input with placeholder, validation, password mode, and character filtering.

```go
import input "github.com/phoenix-tui/phoenix/components/input/api"

// Create input field
field := input.New(40).
    Placeholder("Enter your name...").
    Validate(func(s string) error {
        if len(s) < 3 {
            return errors.New("name too short")
        }
        return nil
    }).
    Focused(true)
```

**Features**:
- ✅ Placeholder text with custom styling
- ✅ Character limits (max length)
- ✅ Character filtering (e.g., numbers only)
- ✅ Password mode (mask input)
- ✅ Real-time validation
- ✅ Cursor control and navigation
- ✅ Copy/paste support (Ctrl+C/Ctrl+V)

[📖 Full Documentation](./input/README.md)

---

### 2. TextArea - Multiline Text Editor
**Module**: `github.com/phoenix-tui/phoenix/components/input/textarea`

Multiline text editor with advanced cursor control, perfect for shells, code editors, and chat interfaces.

```go
import textarea "github.com/phoenix-tui/phoenix/components/input/textarea/api"

// Create text area with cursor control
ta := textarea.New().
    Width(80).
    Height(20).
    Placeholder("Type here...").
    OnMovement(func(from, to textarea.CursorPos) bool {
        // Validate cursor movement (e.g., protect prompt area)
        if to.Row == 0 && to.Col < 2 {
            return false  // Block movement
        }
        return true
    }).
    OnCursorMoved(func(from, to textarea.CursorPos) {
        // React to cursor changes (e.g., update syntax highlighting)
        if from.Row != to.Row {
            refreshSyntaxHighlight(to.Row)
        }
    })
```

**Features**:
- ✅ **SetCursorPosition(row, col)** - Programmatic cursor positioning
- ✅ **OnMovement(validator)** - Movement validation (protect areas)
- ✅ **OnCursorMoved(handler)** - Cursor movement observer
- ✅ **OnBoundaryHit(handler)** - Boundary hit feedback
- ✅ Multiline editing with word wrap
- ✅ Line numbers and gutters
- ✅ Scrolling support
- ✅ Selection and copy/paste
- ✅ Emacs keybindings

**Use Cases**:
- Shell REPLs (GoSh, custom shells)
- Code editors with syntax highlighting
- SQL clients with multiline queries
- Chat interfaces with history
- Log viewers with scroll

[📖 Full Documentation](./input/textarea/README.md) · [📖 Cursor Control API](./input/textarea/CURSOR_CONTROL_API.md)

---

### 3. List - Scrollable Item List
**Module**: `github.com/phoenix-tui/phoenix/components/list`

Scrollable list with filtering, multi-select, and custom item rendering.

```go
import list "github.com/phoenix-tui/phoenix/components/list/api"

// Create list with items
l := list.New([]string{"Item 1", "Item 2", "Item 3"}).
    Title("Select an option").
    Filter(true).  // Enable filtering
    Height(10)
```

**Features**:
- ✅ Scrollable with keyboard navigation
- ✅ Filtering with fuzzy matching
- ✅ Multi-select mode
- ✅ Custom item rendering
- ✅ Active/selected styling
- ✅ Pagination support

[📖 Full Documentation](./list/README.md)

---

### 4. Viewport - Scrollable Content Area
**Module**: `github.com/phoenix-tui/phoenix/components/viewport`

Scrollable content area for displaying long text, logs, or chat history with configurable wheel scrolling and drag support.

```go
import viewport "github.com/phoenix-tui/phoenix/components/viewport"

// Create scrollable viewport with custom wheel scroll speed
vp := viewport.New(80, 24).
    SetContent(longText).
    MouseEnabled(true).          // Enables wheel + drag scrolling
    SetWheelScrollLines(5)       // 5 lines per wheel tick (default: 3)

// Use in tea.Model
p := tea.New(model, tea.WithMouseAllMotion[Model]())
```

**Features**:
- **Configurable Wheel Scrolling** - Adjust scroll speed (1-10+ lines per tick)
- **Drag Scrolling** - Click and drag to scroll (natural touch behavior)
- Mouse wheel support (default: 3 lines per tick, customizable)
- Keyboard navigation (arrows, page up/down, Home/End, Ctrl+U/D)
- Dynamic content updates with FollowMode (tail -f style)
- Precise scroll position control (SetYOffset)
- Line wrapping and truncation support
- Bounds checking (won't scroll past content)
- Immutable operations (functional updates)

**Wheel Scrolling Configuration**:
- `SetWheelScrollLines(lines int)` - Configure scroll amount per wheel tick
- Default: 3 lines per tick
- Recommended: 1-10 lines for smooth scrolling
- Minimum enforced: 1 line (values < 1 are clamped)
- Dynamic adjustment: Change scroll speed on the fly

**Drag Scrolling**:
- Left mouse button drag to scroll
- Natural direction: drag down → content scrolls up
- Works with 10,000+ lines smoothly
- Automatic bounds clamping

[📖 Full Documentation](./viewport/)
[🎮 Drag Scroll Example](../examples/drag-scroll/)
[🎮 Wheel Scroll Example](../examples/wheel-scroll/)

---

### 5. Table - Data Table with Sorting
**Module**: `github.com/phoenix-tui/phoenix/components/table`

Data table with columns, rows, sorting, and pagination.

```go
import table "github.com/phoenix-tui/phoenix/components/table/api"

// Create table with columns
t := table.New().
    Columns([]table.Column{
        {Title: "Name", Width: 20},
        {Title: "Age", Width: 5},
        {Title: "Email", Width: 30},
    }).
    Rows([]table.Row{
        {"Alice", "30", "alice@example.com"},
        {"Bob", "25", "bob@example.com"},
    }).
    Height(10)
```

**Features**:
- ✅ Column headers with alignment
- ✅ Sortable columns (ascending/descending)
- ✅ Row selection
- ✅ Custom cell rendering
- ✅ Pagination
- ✅ Resizable columns

[📖 API Documentation](./table/api/)

---

### 6. Modal - Dialog Boxes
**Module**: `github.com/phoenix-tui/phoenix/components/modal`

Modal dialogs for confirmations, alerts, and user prompts.

```go
import modal "github.com/phoenix-tui/phoenix/components/modal/api"

// Create confirmation modal
m := modal.New().
    Title("Confirm Action").
    Content("Are you sure you want to delete this file?").
    Buttons([]modal.Button{
        {Label: "Cancel", Action: modal.ActionCancel},
        {Label: "Delete", Action: modal.ActionConfirm, Primary: true},
    }).
    Width(50)
```

**Features**:
- ✅ Customizable title and content
- ✅ Multiple button layouts
- ✅ Keyboard navigation (Tab, Enter, Esc)
- ✅ Centered positioning
- ✅ Overlay/backdrop
- ✅ Action callbacks

[📖 Full Documentation](./modal/README.md)

---

### 7. Progress - Progress Bars and Spinners
**Module**: `github.com/phoenix-tui/phoenix/components/progress`

Progress indicators for loading states and long-running operations.

```go
import progress "github.com/phoenix-tui/phoenix/components/progress/api"

// Progress bar
bar := progress.NewBar(0.75).  // 75% progress
    Width(40).
    ShowPercentage(true)

// Spinner
spinner := progress.NewSpinner().
    Style(progress.SpinnerDots)
```

**Features**:
- ✅ Progress bars with percentage
- ✅ Multiple spinner styles (dots, line, pulse)
- ✅ Customizable colors and styles
- ✅ Indeterminate mode
- ✅ Multi-progress (multiple bars)

[📖 Full Documentation](./progress/README.md)

---

## Installation

### Individual Components (Recommended)

```bash
# Install specific components
go get github.com/phoenix-tui/phoenix/components/input@latest
go get github.com/phoenix-tui/phoenix/components/list@latest
go get github.com/phoenix-tui/phoenix/components/modal@latest
```

### Full Library

```bash
go get github.com/phoenix-tui/phoenix/components@latest
```

---

## Quick Start

### Example: Todo App with Multiple Components

```go
package main

import (
    "fmt"
    "os"

    input "github.com/phoenix-tui/phoenix/components/input/api"
    list "github.com/phoenix-tui/phoenix/components/list/api"
    tea "github.com/phoenix-tui/phoenix/tea/api"
)

type model struct {
    input *input.Input
    list  *list.List
    todos []string
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        case "enter":
            // Add new todo
            if m.input.Value() != "" {
                m.todos = append(m.todos, m.input.Value())
                m.list = list.New(m.todos).Height(10)
                m.input = input.New(40).Focused(true)
            }
            return m, nil
        }
    }

    // Update components
    var cmd tea.Cmd
    updatedInput, cmd := m.input.Update(msg)
    m.input = updatedInput.(*input.Input)
    return m, cmd
}

func (m model) View() string {
    return fmt.Sprintf(
        "Todo List\n\n%s\n\n%s\n\nPress Enter to add, Ctrl-C to quit",
        m.input.View(),
        m.list.View(),
    )
}

func main() {
    p := tea.NewProgram(model{
        input: input.New(40).Placeholder("Add a todo...").Focused(true),
        list:  list.New([]string{}).Height(10),
        todos: []string{},
    })

    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

---

## Examples

Each component includes working examples in `examples/` subdirectories:

### Running Examples

```bash
# From repository root
go run ./components/input/examples/basic
go run ./components/input/textarea/examples/shell_prompt
go run ./components/list/examples/filtered
go run ./components/modal/examples/confirmation
go run ./components/progress/examples/multi_progress
```

See [README_EXAMPLES.md](./README_EXAMPLES.md) for complete examples documentation.

---

## Architecture

All Phoenix components follow Domain-Driven Design:

```
component/
├── domain/          # Business logic
│   ├── model/      # Entities and aggregates
│   ├── value/      # Value objects
│   └── service/    # Domain services
├── infrastructure/  # Technical implementation
│   ├── renderer/   # View rendering
│   └── keybindings/# Keyboard handling
└── api/            # Public interface
    └── component.go
```

**Why DDD?**
- Pure business logic in domain (easy to test)
- Infrastructure swappable (ANSI -> native -> web?)
- API layer provides fluent interface
- High test coverage consistently achieved

---

## Component Integration with phoenix/tea

All Phoenix components integrate seamlessly with the Elm Architecture:

### 1. Component as Model
```go
type model struct {
    input *input.Input
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    updated, cmd := m.input.Update(msg)
    m.input = updated.(*input.Input)
    return m, cmd
}

func (m model) View() string {
    return m.input.View()
}
```

### 2. Component Messages
```go
// Components send messages to notify state changes
switch msg := msg.(type) {
case input.ValidationMsg:
    // Handle validation result
case list.SelectionMsg:
    // Handle item selection
}
```

### 3. Component Composition
```go
// Compose multiple components
func (m model) View() string {
    return fmt.Sprintf(
        "%s\n%s\n%s",
        m.header.View(),
        m.list.View(),
        m.footer.View(),
    )
}
```

---

## Testing

Phoenix components have excellent test coverage:

```bash
# Run all component tests
go test ./...

# Run with coverage
go test -cover ./...

# Specific component
go test ./input/... -v
```

All components have extensive test coverage across API, domain, and infrastructure layers.

---

## Future Components

- Form component (validation, submission)
- Tree component (hierarchical data)
- Menu component (dropdowns, context menus)
- Tabs component (multi-panel views)
- Chart component (graphs, plots)
- File picker component
- Calendar component

### Long-term

- Theme system with presets
- Animation framework
- Accessibility improvements
- Component composition patterns

---

## Comparison with Charm/Bubbles

| Feature | Phoenix Components | Bubbles |
|---------|-------------------|---------|
| Architecture | ✅ DDD (testable) | ⚠️ Monolithic |
| Type Safety | ✅ Generic constraints | ⚠️ interface{} |
| Test Coverage | High across all components | Variable |
| Unicode Support | Perfect (CJK/emoji) | Broken ([lipgloss#562](https://github.com/charmbracelet/lipgloss/issues/562)) |
| Dependencies | Zero external TUI deps | Charm ecosystem |
| Component Count | 7 components | 10+ components |

**Why Phoenix Components?**
- Modern DDD architecture (clean, testable)
- Perfect Unicode support (no emoji bugs)
- Type-safe API (compile-time guarantees)
- Zero dependency on broken Charm ecosystem
- Built for extensibility and customization

---

## Contributing

Phoenix Components is part of the Phoenix TUI Framework project. Contributions welcome!

See [../../docs/dev/](../../docs/dev/) for:
- Architecture documentation
- Development roadmap
- API design principles
- Component design patterns

---

## License

MIT (planned)

---

## Resources

- **Examples**: [README_EXAMPLES.md](./README_EXAMPLES.md) - Complete examples guide
- **Architecture**: [../../docs/dev/ARCHITECTURE.md](../../docs/dev/ARCHITECTURE.md)
- **API Design**: [../../docs/dev/API_DESIGN.md](../../docs/dev/API_DESIGN.md)
- **Roadmap**: [../../docs/dev/ROADMAP.md](../../docs/dev/ROADMAP.md)
- **Phoenix TEA**: [../tea/](../tea/) - Elm Architecture implementation
- **Phoenix Layout**: [../layout/](../layout/) - Box model and flexbox
- **Phoenix Style**: [../style/](../style/) - CSS-like styling

---

*Built with Domain-Driven Design and Modern Go*
