# Value Objects - phoenix/style/domain/value

This package contains immutable value objects for the Phoenix style system, following Domain-Driven Design principles.

## Overview

Value objects represent domain concepts that are defined by their attributes, not their identity. They are immutable and compared by value equality.

## Color Value Object

The `Color` type represents an RGB color value with support for multiple input formats and terminal capability adaptation.

### Construction

```go
import "github.com/phoenix-tui/phoenix/style/domain/value"

// RGB constructor (0-255 for each component)
red := value.RGB(255, 0, 0)
green := value.RGB(0, 255, 0)
blue := value.RGB(0, 0, 255)

// Hex constructor (supports #RRGGBB and #RGB formats)
purple, err := value.Hex("#FF00FF")
if err != nil {
    // Invalid hex string
}

// Short hex format (expands #RGB to #RRGGBB)
white, _ := value.Hex("#FFF")  // Same as #FFFFFF

// From ANSI 256-color code
color := value.FromANSI256(196)  // Converts ANSI code to RGB
```

### Methods

```go
// Get RGB components
r, g, b := color.RGB()

// Get hex representation
hex := color.Hex()  // Returns "#FF00FF"

// Equality comparison
if color1.Equal(color2) {
    // Colors are equal
}

// String representation (for debugging)
fmt.Println(color)  // "Color(r=255, g=0, b=255, hex=#FF00FF)"
```

### Design Decisions

- **Immutable**: Once created, colors cannot be changed
- **RGB Internal**: Stored as RGB (0-255) internally for consistency
- **No Color Space**: Only sRGB is supported (TUI standard)
- **Value Semantics**: Colors are compared by value, not identity

## TerminalCapability Enum

The `TerminalCapability` type represents the color support level of a terminal.

### Constants

```go
const (
    NoColor     // Monochrome terminal (no color support)
    ANSI16      // 16 basic colors (8 normal + 8 bright)
    ANSI256     // 256 colors (16 + 216 cube + 24 grayscale)
    TrueColor   // 24-bit RGB (16.7 million colors)
)
```

### Methods

```go
capability := value.TrueColor

// Check color support
if capability.SupportsColor() {
    // Terminal supports some color
}

if capability.SupportsTrueColor() {
    // Terminal supports RGB colors
}

if capability.Supports256Color() {
    // Terminal supports 256 colors or better
}

if capability.Supports16Color() {
    // Terminal supports at least 16 colors
}

// String representation
fmt.Println(capability.String())  // "TrueColor"
```

## Border Value Object

The `Border` type represents a border style with 8 Unicode characters for drawing boxes.

### Pre-defined Borders

```go
// Rounded corners (modern style)
border := value.RoundedBorder
// ╭─────╮
// │     │
// ╰─────╯

// Thick/bold lines
border := value.ThickBorder
// ┏━━━━━┓
// ┃     ┃
// ┗━━━━━┛

// Double lines (formal style)
border := value.DoubleBorder
// ╔═════╗
// ║     ║
// ╚═════╝

// Normal single lines (most compatible)
border := value.NormalBorder
// ┌─────┐
// │     │
// └─────┘

// ASCII characters (maximum compatibility)
border := value.ASCIIBorder
// +-----+
// |     |
// +-----+

// No border (invisible)
border := value.HiddenBorder
```

### Custom Borders

```go
customBorder := value.Border{
    Top:         "═",
    Bottom:      "═",
    Left:        "║",
    Right:       "║",
    TopLeft:     "╔",
    TopRight:    "╗",
    BottomLeft:  "╚",
    BottomRight: "╝",
}
```

### Methods

```go
// Check if borders are equal
if border1.Equal(border2) {
    // Borders are identical
}

// Check if border is hidden
if border.IsHidden() {
    // All characters are empty
}

// String representation
fmt.Println(border.String())  // "Border(╭─╮ / │ │ / ╰─╯)"
```

### Unicode Correctness

All pre-defined borders use single Unicode codepoints (one rune each). This ensures correct width calculation and rendering.

## Design Patterns

### Immutability

All value objects are immutable. Methods that would modify the object instead return a new instance:

```go
color := value.RGB(255, 0, 0)  // Red
// color.r = 128  // Not possible - fields are unexported
newColor := value.RGB(128, 0, 0)  // Create new color instead
```

### Value Equality

Value objects are compared by their values, not their identity:

```go
color1 := value.RGB(255, 0, 0)
color2 := value.RGB(255, 0, 0)

// These are different instances but equal values
fmt.Println(color1.Equal(color2))  // true
```

### No External Dependencies

The value package has zero external dependencies (only Go stdlib). This keeps the domain layer pure and easy to test.

## Padding Value Object

The `Padding` type represents CSS-like padding (space inside content).

### Construction

```go
// Individual values (top, right, bottom, left - clockwise from top)
padding := value.NewPadding(1, 2, 3, 4)

// Uniform padding (same on all sides)
padding := value.UniformPadding(5)

// Vertical and horizontal (top/bottom, left/right)
padding := value.VerticalHorizontal(2, 4)  // 2 vertical, 4 horizontal
```

### Methods

```go
// Get individual sides
top := padding.Top()
right := padding.Right()
bottom := padding.Bottom()
left := padding.Left()

// Get totals
horizontal := padding.Horizontal()  // left + right
vertical := padding.Vertical()      // top + bottom
vert, horz := padding.Total()       // both

// Equality
if padding1.Equal(padding2) {
    // Paddings are equal
}
```

### Negative Values

Negative padding values are automatically clamped to 0:

```go
padding := value.NewPadding(-1, -2, -3, -4)  // All become 0
```

## Margin Value Object

The `Margin` type represents CSS-like margin (space outside content).

### Construction

```go
// Same constructors as Padding
margin := value.NewMargin(1, 2, 3, 4)
margin := value.UniformMargin(5)
margin := value.VerticalHorizontalMargin(2, 4)
```

### Type Safety

Margin and Padding are separate types for compile-time safety:

```go
padding := value.UniformPadding(2)
margin := value.UniformMargin(3)

// This won't compile (different types):
// if padding.Equal(margin) { ... }
```

## Size Value Object

The `Size` type represents width/height constraints with min/max bounds.

### Construction

```go
// Empty size (no constraints)
size := value.NewSize()

// With dimensions
size := value.WithWidth(100)
size := value.WithHeight(50)
```

### Setting Constraints

All setters return new instances (immutability):

```go
size := value.NewSize().
    SetWidth(100).
    SetHeight(50).
    SetMinWidth(10).
    SetMaxWidth(200).
    SetMinHeight(5).
    SetMaxHeight(100)
```

### Getting Values

```go
// Returns (value int, isSet bool)
width, isSet := size.Width()
if isSet {
    fmt.Println("Width:", width)
}

// Similar for all constraints
height, isSet := size.Height()
minWidth, isSet := size.MinWidth()
maxWidth, isSet := size.MaxWidth()
minHeight, isSet := size.MinHeight()
maxHeight, isSet := size.MaxHeight()
```

### Validation

```go
// Clamps value to min/max constraints
validatedWidth := size.ValidateWidth(150)
validatedHeight := size.ValidateHeight(75)
```

### Edge Case: Min > Max

If min > max, validation applies min first, then max:

```go
size := value.NewSize().SetMinWidth(100).SetMaxWidth(50)
validated := size.ValidateWidth(200)  // Returns 50 (clamped to max)
```

## Alignment Value Object

The `Alignment` type represents horizontal and vertical alignment.

### Horizontal Alignment

```go
const (
    AlignLeft
    AlignCenter
    AlignRight
)
```

### Vertical Alignment

```go
const (
    AlignTop
    AlignMiddle
    AlignBottom
)
```

### Construction

```go
// Explicit alignment
alignment := value.NewAlignment(value.AlignCenter, value.AlignMiddle)

// Convenience constructors (all 9 combinations)
alignment := value.LeftTop()
alignment := value.LeftMiddle()
alignment := value.LeftBottom()
alignment := value.CenterTop()
alignment := value.CenterMiddle()
alignment := value.CenterBottom()
alignment := value.RightTop()
alignment := value.RightMiddle()
alignment := value.RightBottom()
```

### Methods

```go
// Get components
horizontal := alignment.Horizontal()  // Returns HorizontalAlignment
vertical := alignment.Vertical()      // Returns VerticalAlignment

// Equality
if alignment1.Equal(alignment2) {
    // Alignments are equal
}

// String representation
fmt.Println(alignment.String())  // "Alignment(Center, Middle)"
```

## Testing

All value objects have comprehensive test coverage across every type (Color, TerminalCapability, Border, Padding, Margin, Size, Alignment).

Run tests:

```bash
go test ./domain/value/... -cover
```

## Related Packages

- **domain/service**: `ColorAdapter`, `SpacingCalculator`, `TextAligner`
- **infrastructure/ansi**: `ANSICodeGenerator` for low-level ANSI codes
- **api**: Public API for creating styles with these value objects

---

**Package**: `github.com/phoenix-tui/phoenix/style/domain/value`
**Architecture**: DDD Value Objects (immutable, pure domain logic)
