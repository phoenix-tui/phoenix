# Phoenix TUI - Hover Highlight Demo

This example demonstrates **phoenix/mouse** hover detection capabilities with interactive button highlighting.

## Features

✅ **6 Interactive Buttons** - Arranged in a responsive grid layout
✅ **Hover Detection** - Buttons change appearance when mouse hovers
✅ **Click Handling** - Track which button was last clicked
✅ **Real-time State** - Visual feedback showing current hover and click state
✅ **Responsive Layout** - Buttons reposition on terminal resize
✅ **Phoenix TEA Pattern** - Demonstrates Elm Architecture with mouse events

## What This Demonstrates

### 1. Hover State Management
```go
// Process hover across multiple component areas
eventType := m.mouse.ProcessHover(pos, areas)

switch eventType {
case mouse.EventHoverEnter:
    m.hoveredID = m.mouse.CurrentHoverComponent()
case mouse.EventHoverLeave:
    m.hoveredID = ""
}
```

### 2. Bounding Box Definition
```go
// Define hover detection areas for each button
button{
    id:   "button1",
    text: "Button 1",
    area: mouse.NewBoundingBox(x, y, width, height),
}
```

### 3. Visual Feedback
- **Normal State**: White text on dark gray background
- **Hovered State**: Black text on bright cyan background (bold)
- **Status Display**: Shows currently hovered and last clicked button

### 4. Tea Integration
- Proper use of `tea.MouseMsg` for mouse events
- `tea.WindowSizeMsg` for responsive layout
- Clean shutdown with mouse cleanup

## Building and Running

### Build
```bash
cd examples/hover-highlight
go build
```

### Run
```bash
# From examples/hover-highlight/
./hover-highlight

# Or from tui root:
go run ./examples/hover-highlight
```

## Usage

- **Move mouse** over buttons to see hover highlighting
- **Click** buttons to register selection (shown in status)
- **Press 'r'** to reset state (clear hover and click)
- **Press 'q'** or **Ctrl+C** to quit

## UI Layout

```
Phoenix TUI - Hover Highlight Demo
════════════════════════════════════

  ╭────────────╮  ╭────────────╮  ╭────────────╮
  │  Button 1  │  │  Button 2  │  │  Button 3  │
  ╰────────────╯  ╰────────────╯  ╰────────────╯

     ╭────────────╮  ╭────────────╮
     │  Button 4  │  │  Button 5  │
     ╰────────────╯  ╰────────────╯

          ╭────────────╮
          │  Button 6  │
          ╰────────────╯

Hovered: button2
Last clicked: button1

Move mouse to hover • Click to select • 'r' to reset • 'q' to quit
```

**Visual Feedback:**
- **Normal button**: Single-line border (`╭─╮│╰─╯`)
- **Hovered button**: Double-line border (`╔═╗║╚═╝`) - more prominent!
- **Status section**: Shows current hover and last click with styled text

## Code Structure

### Model
```go
type model struct {
    mouse         *mouse.Mouse    // Phoenix mouse handler
    buttons       []button        // UI buttons with hover areas
    hoveredID     string          // Current hover state
    lastClickedID string          // Last click state
    width, height int             // Terminal dimensions
    ready         bool            // Initialization flag
}
```

### Button Definition
```go
type button struct {
    id   string             // Unique identifier
    text string             // Display text
    x, y int               // Position
    area mouse.BoundingBox  // Hover detection area
}
```

### Key Methods

- **`handleMouseEvent()`** - Processes hover and click events
- **`layoutButtons()`** - Calculates responsive button positions
- **`drawButton()`** - Renders button with hover highlighting
- **`View()`** - Main rendering function (Elm Architecture)

## Performance Notes

- **Hover detection**: O(n) where n = number of buttons (negligible for typical UIs)
- **Event processing**: Single pass through component areas
- **No allocations** in hot path (hover detection)
- **Efficient state tracking**: Only updates on actual state changes

## Integration with Phoenix Components

This example shows **low-level** hover detection suitable for custom components.

For higher-level UI components (TextInput, Button, etc.), use **phoenix/components** which have built-in hover support:

```go
// High-level component with hover (coming in phoenix/components v0.2.0)
btn := components.NewButton().
    Text("Click Me").
    OnHover(func(hovered bool) { /* ... */ }).
    OnClick(func() { /* ... */ })
```

## Terminal Compatibility

- ✅ **Modern terminals** - Full support (iTerm2, Windows Terminal, Alacritty, etc.)
- ✅ **Mouse tracking** - Requires terminal mouse support (most modern terminals)
- ⚠️ **SSH sessions** - Ensure mouse events forwarded (depends on SSH client)
- ❌ **Very old terminals** - May not support mouse tracking

## Next Steps

- Try modifying button layout (add more buttons, different arrangements)
- Experiment with different hover styles (colors, borders, animations)
- Implement drag-and-drop (use `mouse.IsDragging()`)
- Add right-click context menus (check `msg.Button == tea.MouseButtonRight`)
- Integrate with phoenix/components for higher-level abstractions

## Phoenix Mouse API Reference

See [phoenix/mouse documentation](../../mouse/README.md) for full API details:

- `mouse.New()` - Create mouse handler
- `mouse.Enable()` / `Disable()` - Control mouse tracking
- `mouse.ProcessHover()` - Process hover across component areas
- `mouse.IsHovering()` - Check if any component hovered
- `mouse.CurrentHoverComponent()` - Get hovered component ID
- `mouse.NewBoundingBox()` - Define hover detection areas

---

**Phoenix TUI Framework** - Modern TUI library for Go
**GitHub**: https://github.com/phoenix-tui/phoenix
**Version**: v0.1.0-beta.5 (Week 15 - Mouse & Clipboard)
