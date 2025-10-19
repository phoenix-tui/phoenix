# phoenix/layout - Terminal Layout Engine

Phoenix Layout implements a CSS-like box model and Flexbox layout for terminal user interfaces with support for padding, margins, borders, alignment, and flexible layouts.

**Status**: ✅ Week 9-10 Complete - Box Model + Flexbox Ready
**Module**: `github.com/phoenix-tui/phoenix/layout`
**Coverage**: 98.3% (API), 99.2% (domain/model), 95.6% (domain/service)

## Features

- ✅ **CSS Box Model** - Content, padding, border, margin layers
- ✅ **Flexbox Layout** - Row/column layouts with flexible sizing (NEW!)
- ✅ **Unicode-Aware** - Correct width calculation for CJK and emoji
- ✅ **Fluent API** - Chainable method calls for easy styling
- ✅ **Type-Safe** - Compile-time guarantees
- ✅ **DDD Architecture** - Clean, testable, maintainable
- ✅ **98%+ Coverage** - Comprehensive test suite

## Quick Start

```go
package main

import (
    "fmt"
    layout "github.com/phoenix-tui/phoenix/layout/api"
)

func main() {
    box := layout.NewBox("Hello World").
        PaddingAll(1).
        Border().
        Render()

    fmt.Println(box)
}
```

Output:
```
┌─────────────┐
│             │
│ Hello World │
│             │
└─────────────┘
```

## API Reference

### Creating Boxes

```go
box := layout.NewBox("Hello")
```

### Size Constraints

```go
box.Width(20).Height(10)        // Exact size
box.MinWidth(10).MinHeight(5)   // Minimum
box.MaxWidth(80).MaxHeight(24)  // Maximum
```

### Padding, Border, Margin

```go
box.PaddingAll(2)               // All sides
box.PaddingVH(1, 3)             // Vertical, horizontal
box.Padding(1, 2, 1, 2)         // Top, right, bottom, left

box.Border()                    // Enable border
box.NoBorder()                  // Disable border

box.MarginAll(3)                // All sides
box.MarginVH(2, 5)              // Vertical, horizontal
```

**Note**: Borders automatically add +1 cell aesthetic padding.

### Alignment

```go
box.AlignLeft()                 // Default
box.AlignCenter()               // Center horizontal + vertical
box.AlignRight()                // Right horizontal

box.AlignTop()                  // Top vertical
box.AlignMiddle()               // Middle vertical
box.AlignBottom()               // Bottom vertical

box.Align(value.AlignCenter, value.AlignTop)  // Custom
```

### Rendering & Layout

```go
output := box.Render()          // Generate string
fmt.Println(box)                // Uses String()

pos := box.Layout(80, 24)       // Position in parent
fmt.Printf("At (%d, %d)\n", pos.X(), pos.Y())
```

## Flexbox Layout (NEW!)

### Creating Flexbox Containers

```go
// Horizontal layout (row)
flex := layout.Row().
    Gap(2).
    JustifyStart().
    AlignStretch().
    AddRaw("Item 1").
    AddRaw("Item 2").
    AddRaw("Item 3").
    Render(80, 24)

// Vertical layout (column)
flex := layout.Column().
    Gap(1).
    JustifyCenter().
    AlignCenter().
    AddRaw("Header").
    AddRaw("Content").
    AddRaw("Footer").
    Render(80, 24)
```

### Justify Content (Main Axis)

Controls how items are distributed along the main axis (row = horizontal, column = vertical):

```go
flex.JustifyStart()        // Pack at start (default)
flex.JustifyEnd()          // Pack at end
flex.JustifyCenter()       // Center items
flex.JustifySpaceBetween() // Equal spacing between items
```

Visual (Row):
```
JustifyStart:       [1][2][3]         (space)
JustifyEnd:         (space)         [1][2][3]
JustifyCenter:      (space)  [1][2][3]  (space)
JustifySpaceBetween: [1]    (gap)    [2]    (gap)    [3]
```

### Align Items (Cross Axis)

Controls how items align along the cross axis (row = vertical, column = horizontal):

```go
flex.AlignStretch()  // Stretch to fill (default)
flex.AlignStart()    // Align at start
flex.AlignEnd()      // Align at end
flex.AlignCenter()   // Center items
```

Visual (Row):
```
AlignStretch:  ┌───┐ ┌───┐ ┌───┐
               │ 1 │ │ 2 │ │ 3 │  ← All same height
               └───┘ └───┘ └───┘

AlignStart:    ┌───┐ ┌───┐ ┌───┐
               │ 1 │ │ 2 │ │ 3 │  ← Aligned to top
               │   │ └───┘ │   │
               └───┘       └───┘
```

### Gap Spacing

Add spacing between items:

```go
flex.Gap(3)  // 3 cells between each item
```

### Shell Layout Examples

#### Horizontal Split (Prompt + Input)

```go
prompt := layout.NewBox("$ ").Width(2)
input := layout.NewBox("echo 'Hello World'")

shell := layout.Row().
    Gap(0).
    JustifyStart().
    Add(prompt).
    Add(input).
    Render(80, 1)
// Output: $ echo 'Hello World'
```

#### Vertical Split (History + Input)

```go
history := layout.NewBox("Command history...\n$ ls\n$ cd projects")
input := layout.NewBox("$ ").Border()

shell := layout.Column().
    Gap(0).
    JustifyStart().
    Add(history).
    Add(input).
    Render(80, 24)
```

#### Three-Column Layout

```go
col1 := layout.NewBox("Left").Border()
col2 := layout.NewBox("Center").Border()
col3 := layout.NewBox("Right").Border()

layout := layout.Row().
    Gap(2).
    JustifySpaceBetween().
    Add(col1).
    Add(col2).
    Add(col3).
    Render(80, 10)
```

## Examples

See [examples/](examples/) for complete examples:

```bash
# Run examples
go run ./examples/basic          # Basic box model usage
go run ./examples/dialog         # Dialog boxes and modals
go run ./examples/shell_layouts/horizontal  # Horizontal layouts for shells
go run ./examples/shell_layouts/vertical    # Vertical layouts for shells
go run ./examples/shell_layouts/complex     # Nested and complex layouts
```

## Box Model

```
┌─────────────────────────┐
│       Margin            │  ← Outside spacing
│  ┌──────────────────┐   │
│  │     Border       │   │  ← Visual boundary
│  │  ┌───────────┐   │   │
│  │  │  Padding  │   │   │  ← Inside spacing
│  │  │ ┌───────┐ │   │   │
│  │  │ │Content│ │   │   │  ← Text content
│  │  │ └───────┘ │   │   │
│  │  └───────────┘   │   │
│  └──────────────────┘   │
└─────────────────────────┘
```

## Architecture

```
layout/
├── domain/          # DDD business logic
│   ├── model/      # Box, Node entities
│   ├── value/      # Size, Position, Spacing
│   └── service/    # Measure, Layout, Render services
├── api/            # Public fluent API (100% coverage)
└── examples/       # Usage examples
```

## Status

### Week 9 (Complete)
- ✅ Day 1-2: Value Objects + Domain Models (Size, Position, Spacing, Alignment, Box, Node)
- ✅ Day 3: Layout Engine Services (Measure, Layout, Render)
- ✅ Day 4: Public API + Examples + Documentation

### Week 10 (Complete)
- ✅ Day 5: Flexbox Domain Model (FlexContainer, Direction, JustifyContent, AlignItems)
- ✅ Day 6: Flexbox Layout Engine (Row/column distribution, gap support)
- ✅ Day 7: Flexbox Public API + Shell Layout Examples

## Roadmap

### v0.2.0 (Future)
- Grid layout
- Absolute positioning
- Z-index stacking

## Testing

```bash
go test ./...              # All tests
go test ./... -cover       # With coverage
go test ./api/... -v       # API tests verbose
```

**Current Coverage**:
- API: 98.3%
- Domain/Model: 99.2%
- Domain/Service: 95.6%
- Domain/Value: 98.5%
- **Overall: 97.9%** (exceeds 95% target!)

## License

MIT (planned)

---

Part of **Phoenix TUI Framework** - Modern Go TUI library
