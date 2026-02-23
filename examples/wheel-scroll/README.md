# Mouse Wheel Scrolling Demo

Demonstrates the configurable mouse wheel scrolling feature in Phoenix Viewport component.

## Features

- **Configurable scroll speed**: Adjust lines scrolled per wheel tick (1, 3, 5, 10 lines)
- **Real-time statistics**: Track scroll position, lines scrolled, and wheel events
- **Interactive controls**: Change scroll speed on the fly
- **Visual feedback**: Progress indicator and help overlay

## Usage

```bash
go run main.go
```

## Controls

| Key | Action |
|-----|--------|
| **Mouse Wheel Up/Down** | Scroll content up/down |
| **1** | Slow scroll (1 line per wheel tick) |
| **2** | Default scroll (3 lines per wheel tick) |
| **3** | Fast scroll (5 lines per wheel tick) |
| **4** | Very fast scroll (10 lines per wheel tick) |
| **r** | Reset to top |
| **h** | Toggle help overlay |
| **q** | Quit |

## API Example

```go
// Create viewport with custom wheel scroll speed
vp := viewport.New(80, 24).
    MouseEnabled(true).
    SetWheelScrollLines(5) // 5 lines per wheel tick

// Change scroll speed dynamically
vp = vp.SetWheelScrollLines(10) // Now 10 lines per wheel tick

// Default is 3 lines per wheel tick
vp := viewport.New(80, 24).MouseEnabled(true) // Uses default: 3
```

## Implementation Details

### Scroll Speed Configuration

The viewport component allows configurable wheel scrolling via `SetWheelScrollLines(lines int)`:

- **Default**: 3 lines per wheel tick
- **Minimum**: 1 line (values < 1 are clamped to 1)
- **Recommended**: 1-10 lines for smooth scrolling
- **Large values**: 10+ lines for fast navigation

### Bounds Handling

The viewport automatically clamps scroll position to valid bounds:
- Top boundary: offset = 0
- Bottom boundary: offset = totalLines - viewportHeight

### Immutability

All viewport operations return a new instance (immutable):

```go
vp1 := viewport.New(80, 24).SetWheelScrollLines(3)
vp2 := vp1.SetWheelScrollLines(5)

// vp1 still uses 3 lines/tick
// vp2 uses 5 lines/tick
```

## Performance

Wheel scrolling is highly optimized:
- **O(1)** time complexity per wheel event
- **Zero allocations** in hot path (after initial setup)
- Handles rapid wheel events smoothly
- Tested with 10,000+ line content

## Testing

The feature has comprehensive test coverage:
- Custom scroll values (1, 5, 10+ lines)
- Boundary conditions (top/bottom)
- Small content (no scrolling needed)
- Multiple sequential wheel events
- Immutability verification
- Fluent API chaining

See `components/viewport/viewport_test.go` for full test suite.

## Features Demonstrated

1. Configurable wheel scroll speed (SetWheelScrollLines API)
2. Dynamic speed adjustment (change on the fly)
3. Real-time feedback (statistics display)
4. Edge case handling (bounds, small content)

Part of Phoenix TUI Framework
