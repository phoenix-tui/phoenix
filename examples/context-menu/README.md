# Context Menu Example

This example demonstrates **Phoenix context menu positioning** with smart edge detection.

## Features

- **Smart Positioning**: Menu automatically adjusts position to stay fully visible
- **Edge Detection**: Menu shifts left/up when near screen edges
- **Interactive Menu**: 7 menu items with visual selection feedback
- **Keyboard Navigation**: Use arrow keys and Enter to navigate
- **Mouse Interaction**: Click menu items or close by clicking outside

## What It Demonstrates

### MenuPositioner Service
- Calculates optimal menu position at cursor
- Detects right/bottom edge overflow
- Shifts menu to keep it fully visible
- Handles corner cases (bottom-right corner)
- Pins oversized menus to top-left

### Mouse Integration
- Right-click to open context menu
- Left-click to select menu items
- Click outside menu to close
- Full mouse event handling

### TEA Pattern
- Elm Architecture implementation
- Clean model/update/view separation
- Message-based event handling
- Immutable state updates

## Usage

### Build
```bash
cd examples/context-menu
go mod tidy
go build -o context-menu.exe
```

### Run
```bash
./context-menu.exe
```

### Interact
1. **Right-click** anywhere to open context menu
2. Menu will appear at cursor position
3. If near edge, menu adjusts position automatically
4. **Click** menu items to activate
5. **↑/↓** keys + **Enter** to navigate and select
6. Click **outside menu** to close
7. Press **q** or **ESC** to quit

## Expected Behavior

### Positioning Examples

**Normal (center of screen)**:
```
Right-click at (40, 12)
→ Menu appears at (40, 12)
```

**Right edge**:
```
Right-click at (70, 12)
→ Menu shifts left to (60, 12) to fit
  (Menu width = 20, Screen width = 80)
```

**Bottom edge**:
```
Right-click at (40, 20)
→ Menu shifts up to (40, 17) to fit
  (Menu height = 7, Screen height = 24)
```

**Bottom-right corner**:
```
Right-click at (70, 20)
→ Menu shifts both left and up to (60, 17)
```

### Menu Items
- Copy
- Paste
- Cut
- --- (separator)
- Properties
- Delete
- Refresh

## Code Structure

### Model
```go
type model struct {
    mouse         *mouse.Mouse  // Mouse handler
    menu          contextMenu   // Context menu
    lastClick     mouse.Position // Last click position
    statusMessage string        // Status message
    width         int           // Terminal width
    height        int           // Terminal height
    ready         bool          // UI ready flag
}
```

### Context Menu
```go
type contextMenu struct {
    items    []menuItem // Menu items
    x        int        // Adjusted X position
    y        int        // Adjusted Y position
    width    int        // Menu width
    height   int        // Menu height (items count)
    selected int        // Selected item index
    visible  bool       // Visibility flag
}
```

### Smart Positioning
```go
// Calculate optimal position
safePos := m.mouse.CalculateMenuPosition(
    cursorPos,
    m.menu.width,
    m.menu.height,
    m.width,
    m.height,
)
```

## Testing Tips

1. **Edge Testing**: Right-click near all four edges
2. **Corner Testing**: Right-click in all four corners
3. **Center Testing**: Right-click in center (should stay at cursor)
4. **Keyboard Testing**: Use ↑/↓ + Enter for selection
5. **Close Testing**: Click outside menu to close

## Implementation Details

### MenuPositioner Algorithm
1. **Preferred**: Position at cursor
2. **Right overflow**: Shift left (`x = screenWidth - menuWidth`)
3. **Bottom overflow**: Shift up (`y = screenHeight - menuHeight`)
4. **Corner overflow**: Shift both directions
5. **Oversized**: Pin to (0,0) if menu larger than screen

### Visual Feedback
- **Selected item**: Cyan background + bold
- **Normal items**: Blue background + white text
- **Separator**: Horizontal line (─)
- **Status message**: Shows adjustment info

## Phoenix Libraries Used

- `phoenix/mouse` - Mouse event handling + MenuPositioner
- `phoenix/style` - Terminal styling
- `phoenix/tea` - Elm Architecture (MVU pattern)

## Learning Points

1. **Smart positioning prevents UI overflow**
2. **MenuPositioner is stateless** (pure domain service)
3. **Position calculation is deterministic** (same input = same output)
4. **Edge detection is automatic** (no manual checks needed)
5. **Clean API** (`CalculateMenuPosition` single method)

## Next Steps

After understanding this example:
1. Try creating multi-level menus (submenus)
2. Add dynamic menu items based on context
3. Implement menu animations (fade in/out)
4. Add custom menu themes/styles
5. Create reusable menu component

## Related Examples

- `hover-highlight/` - Hover detection
- `drag-scroll/` - Drag scrolling
- `wheel-scroll/` - Mouse wheel scrolling

---

**Part of Phoenix TUI Framework Week 15 (Advanced Mouse Features)**
**Author**: Phoenix TUI Contributors
**License**: MIT
