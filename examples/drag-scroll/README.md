# Drag Scroll Example

This example demonstrates **click-and-drag scrolling** in the Phoenix Viewport component.

## Features Demonstrated

1. **Drag Scrolling**: Click and hold the left mouse button, then drag up/down to scroll content
2. **Mouse Wheel Scrolling**: Use mouse wheel for traditional scrolling
3. **Keyboard Scrolling**: Arrow keys, Page Up/Down, Home/End still work
4. **Visual Feedback**: Real-time scroll position indicator and scroll arrows
5. **Bounds Checking**: Content won't scroll past top or bottom
6. **Immutability**: All viewport operations return new instances

## How to Run

```bash
go run main.go
```

## Controls

### Mouse
- **Left Click + Drag**: Scroll content by dragging
- **Mouse Wheel Up/Down**: Scroll 3 lines at a time

### Keyboard
- **Arrow Up/Down**: Scroll 1 line at a time
- **Page Up/Down**: Scroll by viewport height
- **Ctrl+U/D**: Scroll by half viewport height
- **Home**: Jump to top
- **End**: Jump to bottom
- **R**: Reset to top
- **E**: Jump to end
- **Q / Ctrl+C**: Quit

## Implementation Details

### Enabling Drag Scrolling

```go
v := viewport.New(width, height).
    MouseEnabled(true) // Enables BOTH wheel and drag scrolling

p := tea.New(
    model,
    tea.WithMouseAllMotion(), // Required for drag events
)
```

### How It Works

1. **Press**: When left mouse button is pressed, viewport records:
   - Starting Y position (where drag began)
   - Current scroll offset (baseline for delta calculation)

2. **Motion**: During drag (motion with button held):
   - Calculate delta: `currentY - startY`
   - Calculate new offset: `startOffset - delta`
   - Apply with bounds checking (domain layer clamps to valid range)

3. **Release**: When mouse button is released:
   - Clear drag state
   - Content remains at final scroll position

### Scroll Direction

- **Drag down** (+Y) → Content scrolls **up** (lower offset)
- **Drag up** (-Y) → Content scrolls **down** (higher offset)

This matches natural touch/trackpad behavior: moving your hand down reveals earlier content.

## Architecture

The drag scrolling implementation follows Phoenix's DDD architecture:

- **API Layer** (`viewport.go`): Handles tea.MouseMsg events, tracks drag state
- **Domain Layer** (`model/viewport.go`): Enforces scroll bounds, immutability
- **Infrastructure Layer** (`infrastructure/keybindings.go`): Key detection utilities

### Drag State

The viewport maintains drag state:
```go
type Viewport struct {
    // ... other fields
    isDragging   bool  // Is a drag currently in progress?
    dragStartY   int   // Y position where drag started
    scrollStartY int   // Scroll offset when drag started
}
```

### Immutability

All drag operations return new viewport instances:
```go
// Start drag
v2 := v.withDragState(true, msg.Y, v.domain.ScrollOffset())

// Update scroll during drag
v3 := v2.withDomain(v.domain.WithScrollOffset(newOffset))

// End drag
v4 := v3.withDragState(false, 0, 0)
```

## Testing

The drag scroll implementation has **14 comprehensive tests** covering:

- Basic drag operations (start, motion, release)
- Bounds checking (top, bottom)
- Edge cases (small content, empty content, large content)
- Button filtering (left only, not right/middle)
- Multiple drags
- Immutability
- Disabled state

**Test coverage**: 98.6% (exceeds 95% target)

Run tests:
```bash
cd ../../components/viewport
go test -v -run TestViewport_DragScroll
```

## Performance

Drag scrolling is **extremely fast**:
- No allocations in hot path (reuses existing viewport methods)
- Simple delta calculation: `O(1)`
- Bounds checking handled by domain layer: `O(1)`
- Works smoothly with 10,000+ lines of content

## Comparison with Other TUI Frameworks

### Bubbletea/Bubbles Viewport

Bubbletea does NOT have built-in drag scrolling support. Users must:
1. Manually track mouse press/release events
2. Calculate scroll deltas themselves
3. Handle bounds checking manually

### Phoenix Viewport

Phoenix provides drag scrolling out-of-the-box:
- Just enable `MouseEnabled(true)` and `WithMouseAllMotion()`
- Drag tracking handled automatically
- Bounds checking built-in
- Immutable operations

## Next Steps

Try modifying this example:
- Add horizontal scrolling (X-axis drag)
- Implement momentum scrolling (physics simulation)
- Add visual drag indicator (highlight during drag)
- Customize scroll sensitivity

## Related Examples

- `examples/viewport-basic/` - Basic viewport usage
- `examples/hover-detection/` - Mouse hover events
- `examples/click-detection/` - Click and double-click detection

## Documentation

- [Viewport API Documentation](../../components/viewport/README.md)
- [Mouse Module Documentation](../../mouse/README.md)
- [TEA (Elm Architecture) Guide](../../tea/README.md)

---

**Phoenix TUI Framework** - Next-generation terminal interfaces for Go
