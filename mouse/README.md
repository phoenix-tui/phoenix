# phoenix/mouse

**Cross-platform mouse event handling for Phoenix TUI Framework.**

Part of the Phoenix TUI Framework.

**Module**: `github.com/phoenix-tui/phoenix/mouse`

---

## Features

- **All Mouse Buttons**: Left, Right, Middle, Scroll Wheel
- **Click Detection**: Single, Double, Triple clicks (automatic)
- **Drag & Drop**: State tracking with threshold detection
- **Scroll Wheel**: Up/Down with configurable delta
- **Modifier Keys**: Shift, Ctrl, Alt detection
- **Multi-Protocol**: SGR (modern), X10 (legacy), URxvt (alternative)
- **DDD Architecture**: Rich domain models with behavior
- **Zero Dependencies**: Built from scratch (stdlib only)

---

## Quick Start

### Installation

```bash
go get github.com/phoenix-tui/phoenix/mouse
```

### Basic Usage

```go
package main

import (
    "github.com/phoenix-tui/phoenix/mouse"
    "github.com/phoenix-tui/phoenix/tea"
)

type Model struct {
    lastClick string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        switch msg.Button {
        case mouse.ButtonLeft:
            m.lastClick = "Left click at " + msg.Position.String()
        case mouse.ButtonRight:
            m.lastClick = "Right click at " + msg.Position.String()
        case mouse.ButtonMiddle:
            m.lastClick = "Middle click at " + msg.Position.String()
        }
    }
    return m, nil
}
```

---

## Mouse Buttons

### All Supported Buttons

```go
mouse.ButtonNone       // No button (motion only)
mouse.ButtonLeft       // Left mouse button
mouse.ButtonRight      // Right mouse button
mouse.ButtonMiddle     // Middle button (scroll wheel press)
mouse.ButtonWheelUp    // Scroll wheel up
mouse.ButtonWheelDown  // Scroll wheel down
```

### Button Detection

```go
case tea.MouseMsg:
    button := msg.Button

    if button.IsButton() {
        // Left, Right, or Middle button
        switch button {
        case mouse.ButtonLeft:
            fmt.Println("Left click")
        case mouse.ButtonRight:
            fmt.Println("Right click - Context menu?")
        case mouse.ButtonMiddle:
            fmt.Println("Middle click - Paste?")
        }
    }

    if button.IsWheel() {
        // Scroll wheel event
        if button == mouse.ButtonWheelUp {
            viewport.ScrollUp(3)
        } else {
            viewport.ScrollDown(3)
        }
    }
```

---

## Event Types

```go
mouse.EventPress        // Button pressed (start of click/drag)
mouse.EventRelease      // Button released (end of click/drag)
mouse.EventClick        // Single click (automatic detection)
mouse.EventDoubleClick  // Double click (< 500ms between clicks)
mouse.EventTripleClick  // Triple click (< 500ms between clicks)
mouse.EventDrag         // Mouse drag (motion with button pressed)
mouse.EventMotion       // Mouse motion (no button pressed)
mouse.EventScroll       // Scroll wheel event
```

---

## Use Cases

### 1. TextInput Cursor Positioning (Left Click)

```go
// In phoenix/components/input
func (i *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        if msg.Type == mouse.EventClick && msg.Button == mouse.ButtonLeft {
            // Calculate cursor position from mouse X coordinate
            cursorPos := i.calculateCursorFromX(msg.Position.X())
            return i.WithCursor(cursorPos), nil
        }
    }
    return i, nil
}
```

### 2. Context Menu (Right Click)

```go
type Model struct {
    contextMenu *ContextMenu
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        // Right click = show context menu
        if msg.Type == mouse.EventClick && msg.Button == mouse.ButtonRight {
            m.contextMenu = NewContextMenu(msg.Position)
            return m, nil
        }

        // Left click outside = close context menu
        if msg.Type == mouse.EventClick && msg.Button == mouse.ButtonLeft {
            if m.contextMenu != nil && !m.contextMenu.Contains(msg.Position) {
                m.contextMenu = nil
            }
        }
    }
    return m, nil
}
```

### 3. Viewport Scrolling (Wheel)

```go
// In phoenix/components/viewport
func (v *Viewport) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        if msg.Type == mouse.EventScroll {
            switch msg.Button {
            case mouse.ButtonWheelUp:
                return v.ScrollUp(3), nil  // 3 lines per scroll
            case mouse.ButtonWheelDown:
                return v.ScrollDown(3), nil
            }
        }
    }
    return v, nil
}
```

### 4. Drag & Drop Files

```go
type Model struct {
    dragState mouse.DragState
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        switch msg.Type {
        case mouse.EventPress:
            if msg.Button == mouse.ButtonLeft {
                m.dragState.Start(msg.Position)
            }

        case mouse.EventDrag:
            m.dragState.Update(msg.Position)

        case mouse.EventRelease:
            if m.dragState.IsDrag() {
                // Extract file path from drop event
                filePath := extractDroppedFile(msg)
                if filePath != "" {
                    return m.InsertFile(filePath), nil
                }
            }
            m.dragState.Reset()
        }
    }
    return m, nil
}
```

### 5. Middle Click Paste

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        // Middle click = paste from clipboard (X11 primary selection)
        if msg.Type == mouse.EventClick && msg.Button == mouse.ButtonMiddle {
            return m, tea.Batch(
                clipboard.ReadPrimary(),  // X11 primary selection
                tea.Printf("Pasted at %s", msg.Position),
            )
        }
    }
    return m, nil
}
```

---

## Click Detection

Phoenix automatically detects single, double, and triple clicks:

```go
case tea.MouseMsg:
    switch msg.Type {
    case mouse.EventClick:
        fmt.Println("Single click")

    case mouse.EventDoubleClick:
        fmt.Println("Double click - Select word?")

    case mouse.EventTripleClick:
        fmt.Println("Triple click - Select line?")
    }
```

**Detection Rules**:
- **Double click**: Two clicks within 500ms at same position (+-1 cell tolerance)
- **Triple click**: Three clicks within 500ms at same position
- Automatic timeout after 500ms

---

## Modifier Keys

```go
case tea.MouseMsg:
    mods := msg.Modifiers

    if mods.HasShift() {
        fmt.Println("Shift + Click")
    }

    if mods.HasCtrl() {
        fmt.Println("Ctrl + Click - Multi-select?")
    }

    if mods.HasAlt() {
        fmt.Println("Alt + Click")
    }

    // Common combinations
    if mods.HasCtrl() && msg.Button == mouse.ButtonLeft {
        // Ctrl + Left Click = Open in new tab?
    }

    if mods.HasShift() && msg.Button == mouse.ButtonLeft {
        // Shift + Left Click = Extend selection
    }
```

---

## Mouse Protocols

Phoenix supports multiple mouse protocols with automatic fallback:

### SGR (1006) - Modern (Preferred)

```
\x1b[<0;10;5M   # Left press at (10, 5)
\x1b[<0;10;5m   # Left release at (10, 5)
\x1b[<2;10;5M   # Right press at (10, 5)
\x1b[<64;10;5M  # Scroll up at (10, 5)
\x1b[<65;10;5M  # Scroll down at (10, 5)
```

**Advantages**:
- Supports press/release separately
- Works with large terminals (> 223 columns)
- Modern terminal support (iTerm2, Windows Terminal, Alacritty, Kitty)

### X10 (1000) - Legacy

```
\x1b[M !!       # Click at (1, 1)
\x1b[M`"#       # Click at (96, 35)
```

**Limitations**:
- No press/release distinction
- Limited to 223 columns
- Old terminal compatibility only

### URxvt (1015) - Alternative

```
\x1b[0;10;5M    # Click at (10, 5)
```

---

## Mouse Modes

Enable different tracking modes:

```go
// Button events only (press/release)
mouse.Enable(mouse.ModeButton)

// Button + motion while pressed (drag detection)
mouse.Enable(mouse.ModeMotion)

// All motion events (even without button)
mouse.Enable(mouse.ModeAllMotion)

// Prefer SGR protocol
mouse.Enable(mouse.ModeButton | mouse.ModeSGR)

// Disable mouse tracking
mouse.Disable()
```

---

## Integration with Components

### TextInput - Cursor Positioning

Phoenix TextInput supports mouse click positioning:

```go
input := textinput.New().
    WithPlaceholder("Click to position cursor...").
    WithMouseSupport(true)  // Enable mouse clicks

// User clicks at X coordinate -> cursor moves to that position
```

### Viewport - Scrolling

Phoenix Viewport supports scroll wheel:

```go
viewport := viewport.New(80, 20).
    WithContent(longText).
    WithMouseScroll(true)  // Enable scroll wheel

// User scrolls wheel -> viewport scrolls up/down
```

### List - Item Selection

Phoenix List supports click selection:

```go
list := list.New(items).
    WithMouseSelect(true)  // Enable click selection

// User clicks item -> item selected
// User right-clicks item -> context menu?
```

---

## Advanced: Drag State

Track drag & drop operations:

```go
type DragState struct {
    Active    bool
    Start     Position
    Current   Position
    Delta     Position  // Current - Start
    Distance  int       // Euclidean distance
}

// Usage
dragState := mouse.NewDragState()

case mouse.EventPress:
    dragState.Start(msg.Position)

case mouse.EventDrag:
    dragState.Update(msg.Position)
    if dragState.Distance() > 10 {
        // Drag threshold exceeded
        fmt.Printf("Dragging from %s to %s\n",
            dragState.Start(), dragState.Current())
    }

case mouse.EventRelease:
    if dragState.IsDrag() {
        // Handle drop
        handleDrop(dragState.Start(), dragState.Current())
    }
    dragState.Reset()
```

---

## Architecture

Phoenix mouse follows Domain-Driven Design:

```
mouse/
├── domain/           # Pure business logic
│   ├── value/       # Button, Position, EventType, Modifiers
│   ├── model/       # MouseEvent (aggregate), DragState (entity)
│   └── service/     # ClickDetector, DragTracker, ScrollCalculator
├── infrastructure/   # Technical implementation
│   ├── parser/      # SGR, X10, URxvt parsers
│   ├── ansi/        # ANSI escape sequences
│   └── platform/    # Terminal mode management
├── application/      # Use cases
│   ├── event_processor.go   # Enriches events (click detection)
│   └── mouse_handler.go     # Coordinates parsing/processing
└── api/             # Public interface
    └── mouse.go     # Clean, fluent API
```

---

## Platform Support

| Platform | Protocol | Support |
|----------|----------|---------|
| **iTerm2** | SGR (1006) | Full |
| **Windows Terminal** | SGR (1006) | Full |
| **Alacritty** | SGR (1006) | Full |
| **Kitty** | SGR (1006) | Full |
| **xterm** | SGR (1006) | Full |
| **tmux** | SGR (1006) | Full |
| **Terminal.app** (macOS) | SGR (1006) | Full |
| **GNOME Terminal** | SGR (1006) | Full |
| **Old xterm** | X10 (1000) | Limited |
| **PuTTY** | X10 (1000) | Limited |

---

## Testing

```bash
# Run tests
cd mouse && go test ./...

# With coverage
go test -cover ./...

# Specific test
go test -run TestClickDetector_DoubleClick ./domain/service
```

---

## Examples

See `examples/` directory:

- `click_positioning.go` - TextInput cursor positioning with left/right click
- `scroll.go` - Viewport scrolling with wheel
- `drag_drop.go` - File drag & drop insertion
- `full_demo.go` - Combined demo (all features)

---

## Performance

- **Zero allocations** in hot paths (event parsing)
- **Sub-microsecond** click detection
- **Efficient protocol parsing** (no regex, direct byte parsing)

---

## Comparison

### vs Bubbletea

| Feature | Phoenix | Bubbletea |
|---------|---------|-----------|
| **Click Detection** | Built-in | Manual |
| **Drag Tracking** | Built-in | Manual |
| **Right Click** | Full support | Supported |
| **Middle Click** | Full support | Supported |
| **Protocol Abstraction** | DDD | Raw events |
| **Domain Model** | Rich | Anemic |

---

## FAQ

**Q: Does right click work?**
A: Yes! Full support for left, right, and middle buttons.

**Q: How do I show a context menu on right click?**
A: Check `msg.Button == mouse.ButtonRight` and render your context menu component.

**Q: Can I detect double-click?**
A: Yes! Phoenix automatically detects double and triple clicks.

**Q: Does scroll wheel work in SSH?**
A: Yes, if your terminal supports SGR protocol (most modern terminals do).

**Q: How do I enable mouse support?**
A: Call `mouse.Enable()` at program start, or use `tea.WithMouseCellMotion()` in tea.Program options.

---

## License

MIT

---

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

---

## Credits

Part of the **Phoenix TUI Framework** - Next-generation terminal UI for Go.
