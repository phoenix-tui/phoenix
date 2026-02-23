# Phoenix List Component

A universal selectable list component for Phoenix TUI Framework. Build file pickers, menus, searchable lists, and more with full keyboard navigation and custom rendering.

## Features

- **Single & Multi-Selection** - Radio button or checkbox behavior
- **Keyboard Navigation** - Arrow keys, Vim keys (j/k), page up/down, home/end
- **Custom Rendering** - Full control over item display with custom render functions
- **Filtering** - Built-in and custom filter functions for searchable lists
- **Scrolling** - Automatic viewport scrolling for long lists
- **Immutable API** - All operations return new instances (follows Elm Architecture)
- **Type-Safe** - Fully typed with Go generics support
- **Well-Tested** - High test coverage

## Installation

```bash
go get github.com/phoenix-tui/phoenix/components/list
```

## Quick Start

### Basic File Picker

```go
package main

import (
    "github.com/phoenix-tui/phoenix/components/list/api"
    "github.com/phoenix-tui/phoenix/components/list/domain/value"
    tea "github.com/phoenix-tui/phoenix/tea/api"
)

func main() {
    files := []interface{}{"file1.txt", "file2.go", "file3.md"}
    labels := []string{"file1.txt", "file2.go", "file3.md"}

    l := list.New(files, labels, value.SelectionModeSingle).Height(10)

    p := tea.New(l)
    p.Run()
}
```

### Multi-Selection Todo List

```go
todos := []interface{}{"Buy milk", "Write code", "Read book"}
labels := []string{"Buy milk", "Write code", "Read book"}

l := list.NewMultiSelect(todos, labels).Height(5)

// User can press Space to toggle multiple items, Ctrl+A to select all
```

### Custom Item Rendering

```go
type Person struct {
    Name string
    Age  int
}

people := []interface{}{
    Person{Name: "Alice", Age: 30},
    Person{Name: "Bob", Age: 25},
}
labels := []string{"Alice", "Bob"}

l := list.NewSingleSelect(people, labels).
    Height(5).
    ItemRenderer(func(item interface{}, index int, selected, focused bool) string {
        p := item.(Person)
        prefix := "  "
        if selected {
            prefix = "✓ "
        }
        if focused {
            prefix = "→ "
        }
        return fmt.Sprintf("%s%s (age %d)", prefix, p.Name, p.Age)
    })
```

### Filtered List

```go
files := []interface{}{"main.go", "main_test.go", "README.md", "config.yaml"}
labels := []string{"main.go", "main_test.go", "README.md", "config.yaml"}

l := list.NewSingleSelect(files, labels).
    Height(10).
    Filter(func(item interface{}, query string) bool {
        filename := item.(string)
        return strings.HasSuffix(filename, ".go")
    })

// Apply filter
l.domain = l.domain.SetFilterQuery("go") // Shows only .go files
```

## API Reference

### Constructors

#### `New(values []interface{}, labels []string, mode SelectionMode) *List`
Creates a new list with the given items and selection mode.

#### `NewSingleSelect(values []interface{}, labels []string) *List`
Convenience constructor for single-selection lists.

#### `NewMultiSelect(values []interface{}, labels []string) *List`
Convenience constructor for multi-selection lists.

### Configuration Methods

All configuration methods return a new `*List` (immutable API):

#### `Height(height int) *List`
Sets the visible height of the list.

#### `ItemRenderer(renderer func(item interface{}, index int, selected, focused bool) string) *List`
Sets a custom item renderer function. Receives:
- `item` - The item value
- `index` - Item index in filtered list
- `selected` - Whether item is selected
- `focused` - Whether item has focus

#### `Filter(filterFunc func(item interface{}, query string) bool) *List`
Sets a custom filter function. Receives:
- `item` - The item value
- `query` - The filter query string

Returns `true` to include the item in filtered results.

#### `ShowFilter(show bool) *List`
Enables/disables the filter input display at the bottom of the list.

#### `KeyBindings(bindings []KeyBinding) *List`
Sets custom key bindings. See [Key Bindings](#key-bindings) section.

### Query Methods

#### `SelectedItems() []interface{}`
Returns the currently selected item values.

#### `FocusedItem() interface{}`
Returns the currently focused item value, or `nil` if no items.

#### `SelectedIndices() []int`
Returns the indices of selected items (in filtered list).

#### `FocusedIndex() int`
Returns the index of the focused item (in filtered list).

### tea.Model Methods

#### `Init() tea.Cmd`
Initializes the list component.

#### `Update(msg tea.Msg) (*List, tea.Cmd)`
Handles messages and returns updated list.

#### `View() string`
Renders the list to a string.

## Key Bindings

### Default Key Bindings

| Key | Action |
|-----|--------|
| `↑`, `k` | Move up |
| `↓`, `j` | Move down |
| `PgUp`, `Ctrl+U` | Page up |
| `PgDown`, `Ctrl+D` | Page down |
| `Home`, `g` | Move to start |
| `End`, `G` | Move to end |
| `Space` | Toggle selection |
| `Enter` | Confirm selection |
| `Ctrl+A` | Select all (multi-select only) |
| `Esc` | Clear selection |
| `q`, `Ctrl+C` | Quit |

### Custom Key Bindings

```go
import "github.com/phoenix-tui/phoenix/components/list/infrastructure"

customBindings := []infrastructure.KeyBinding{
    {Key: "w", Action: "move_up"},
    {Key: "s", Action: "move_down"},
    {Key: "a", Action: "move_to_start"},
    {Key: "d", Action: "move_to_end"},
}

l := list.NewSingleSelect(items, labels).
    KeyBindings(customBindings)
```

### Available Actions

- `move_up` - Move focus up
- `move_down` - Move focus down
- `page_up` - Move up one page
- `page_down` - Move down one page
- `move_to_start` - Move to first item
- `move_to_end` - Move to last item
- `toggle_selection` - Toggle selection of focused item
- `select_all` - Select all items (multi-select only)
- `clear_selection` - Clear all selections
- `clear_filter` - Clear filter query
- `confirm` - Confirm selection
- `quit` - Quit program

## Selection Modes

### Single Selection

Like radio buttons - only one item can be selected at a time. Selecting a new item automatically deselects the previous one.

```go
import "github.com/phoenix-tui/phoenix/components/list/domain/value"

l := list.New(items, labels, value.SelectionModeSingle)
```

### Multi Selection

Like checkboxes - multiple items can be selected simultaneously.

```go
l := list.New(items, labels, value.SelectionModeMulti)
```

## Examples

See the `examples/` directory for complete working examples:

- **basic.go** - Simple file picker
- **multi_select.go** - Todo list with multi-selection
- **filtered.go** - File list with custom filter
- **custom_render.go** - People directory with structured data rendering

Run an example:

```bash
cd examples
go run basic.go
```

## Architecture

Phoenix List follows Domain-Driven Design (DDD) principles:

```
list/
├── domain/           # Pure business logic
│   ├── model/       # List aggregate (rich model with behavior)
│   ├── value/       # Value objects (SelectionMode, Item)
│   └── service/     # Domain services (Navigation, Filter)
├── infrastructure/   # Technical implementation
│   └── keybindings.go
└── api/             # Public interface (fluent, type-safe)
    └── list.go
```

### Design Principles

1. **Rich Domain Models** - Business logic encapsulated in domain models, not anemic data structures
2. **Immutability** - All operations return new instances following Elm Architecture
3. **Type Safety** - Leverage Go generics for compile-time safety
4. **Testability** - High test coverage across all layers
5. **Zero External Dependencies** - Only depends on `phoenix/tea` for Elm Architecture

## How gosh Will Use This

Phoenix List is a **universal** component. The [gosh](https://github.com/grpmsoft/gosh) shell will extend it:

```go
// gosh will wrap List for shell-specific features
type HistorySearch struct {
    list *list.List
    fuzzyMatcher *matcher.FuzzyMatcher  // gosh's fuzzy matching
}

type CompletionMenu struct {
    list *list.List
    commandDB *commands.Database  // gosh's command database
}
```

Phoenix List provides the foundation; apps like gosh add domain-specific logic on top.

## Testing

Run tests:

```bash
go test ./...
```

Run with coverage:

```bash
go test ./... -cover
```

View coverage report:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

High test coverage across all layers (domain, infrastructure, API).

## Contributing

Phoenix TUI Framework follows strict quality standards:

- **DDD First** - Rich domain models with behavior
- **High Test Coverage** - Especially domain layer
- **Immutable API** - All operations return new instances
- **Type Safety** - Use Go 1.25+ features
- **Documentation** - Clear examples and API docs

## License

MIT License - See LICENSE file for details

## Related Components

- [phoenix/tea](../../tea) - Elm Architecture implementation
- [phoenix/style](../../style) - CSS-like styling
- [phoenix/layout](../../layout) - Box Model & Flexbox layout
- [phoenix/core](../../core) - Terminal primitives & Unicode

Phoenix TUI Framework is under active development. API may change before a stable release.
