# phoenix/core Architecture

> **Foundation library for Phoenix TUI Framework**
> **Week 3-4 Implementation** | **DDD + Rich Model Architecture**

---

## 📋 Overview

`phoenix/core` provides terminal primitives with **zero dependencies** (except `github.com/rivo/uniseg` for Unicode). This is the foundation that all other Phoenix libraries build upon.

### Responsibilities

- ✅ Terminal capabilities detection (ANSI, colors, features)
- ✅ Raw mode setup/teardown with guaranteed cleanup
- ✅ Signal handling (Ctrl+C, window resize)
- ✅ Unicode width calculation (grapheme cluster aware)
- ✅ Platform abstraction (Unix/Windows)
- ✅ Basic ANSI escape sequences

### Non-Responsibilities

- ❌ Styling (that's `phoenix/style`)
- ❌ Layout (that's `phoenix/layout`)
- ❌ Event loop (that's `phoenix/tea`)
- ❌ High-level components (that's `phoenix/components`)

---

## 🏗️ DDD Architecture

### Layer Structure

```
core/
├── domain/              # Pure business logic (95%+ coverage target)
│   ├── model/          # Rich domain models (Terminal, RawMode)
│   ├── value/          # Value objects (Capabilities, Cell, Position)
│   └── service/        # Domain services (SignalHandler, CapabilitiesDetector)
├── infrastructure/      # Technical implementations (80%+ coverage)
│   ├── platform/       # Platform-specific code (Unix/Windows)
│   └── parser/         # ANSI sequence parsing
├── api/                # Public interface (85%+ coverage)
│   └── core.go        # Exported API
└── testdata/           # Test fixtures
```

---

## 🎯 Domain Layer Design

### Aggregate Root: Terminal

**Terminal** is the main aggregate root representing terminal state and capabilities.

```go
// domain/model/terminal.go
package model

import (
    "github.com/phoenix-tui/phoenix/core/domain/value"
)

// Terminal represents a terminal instance with its capabilities and state.
// This is the aggregate root for terminal operations.
type Terminal struct {
    capabilities *value.Capabilities
    rawMode      *RawMode
    size         value.Size
}

// NewTerminal creates a new Terminal with detected capabilities.
func NewTerminal(caps *value.Capabilities) *Terminal {
    return &Terminal{
        capabilities: caps,
        rawMode:      nil,  // Not in raw mode initially
        size:         value.Size{Width: 80, Height: 24},  // Default
    }
}

// Capabilities returns terminal capabilities (immutable).
func (t *Terminal) Capabilities() *value.Capabilities {
    return t.capabilities
}

// Size returns current terminal size.
func (t *Terminal) Size() value.Size {
    return t.size
}

// IsRawMode returns true if terminal is in raw mode.
func (t *Terminal) IsRawMode() bool {
    return t.rawMode != nil
}

// WithSize returns new Terminal with updated size (immutable).
func (t *Terminal) WithSize(size value.Size) *Terminal {
    newTerm := *t
    newTerm.size = size
    return &newTerm
}

// WithRawMode returns new Terminal in raw mode (immutable).
func (t *Terminal) WithRawMode(rawMode *RawMode) *Terminal {
    newTerm := *t
    newTerm.rawMode = rawMode
    return &newTerm
}

// WithoutRawMode returns new Terminal not in raw mode (immutable).
func (t *Terminal) WithoutRawMode() *Terminal {
    newTerm := *t
    newTerm.rawMode = nil
    return &newTerm
}
```

### Value Objects

#### Capabilities

```go
// domain/value/capabilities.go
package value

// ColorDepth represents terminal color support level.
type ColorDepth int

const (
    ColorDepthNone     ColorDepth = 0   // No colors
    ColorDepth8        ColorDepth = 8   // 8 colors
    ColorDepth256      ColorDepth = 256 // 256 colors
    ColorDepthTrueColor ColorDepth = 16777216 // 24-bit RGB
)

// Capabilities represents terminal capabilities (immutable value object).
type Capabilities struct {
    ansiSupport   bool
    colorDepth    ColorDepth
    mouseSupport  bool
    altScreen     bool
    cursorControl bool
}

// NewCapabilities creates capabilities with validation.
func NewCapabilities(ansi bool, colors ColorDepth, mouse bool, alt bool, cursor bool) *Capabilities {
    return &Capabilities{
        ansiSupport:   ansi,
        colorDepth:    colors,
        mouseSupport:  mouse,
        altScreen:     alt,
        cursorControl: cursor,
    }
}

// SupportsANSI returns true if terminal supports ANSI escape sequences.
func (c *Capabilities) SupportsANSI() bool {
    return c.ansiSupport
}

// ColorDepth returns terminal color support level.
func (c *Capabilities) ColorDepth() ColorDepth {
    return c.colorDepth
}

// SupportsMouse returns true if terminal supports mouse events.
func (c *Capabilities) SupportsMouse() bool {
    return c.mouseSupport
}

// SupportsAltScreen returns true if terminal supports alternate screen buffer.
func (c *Capabilities) SupportsAltScreen() bool {
    return c.altScreen
}

// SupportsCursorControl returns true if terminal supports cursor positioning.
func (c *Capabilities) SupportsCursorControl() bool {
    return c.cursorControl
}
```

#### Position & Size

```go
// domain/value/position.go
package value

// Position represents a position in the terminal (row, column).
// This is an immutable value object.
type Position struct {
    Row int
    Col int
}

// NewPosition creates a position with validation.
func NewPosition(row, col int) Position {
    if row < 0 {
        row = 0
    }
    if col < 0 {
        col = 0
    }
    return Position{Row: row, Col: col}
}

// Size represents terminal dimensions (width, height in cells).
type Size struct {
    Width  int
    Height int
}

// NewSize creates size with validation.
func NewSize(width, height int) Size {
    if width < 1 {
        width = 1
    }
    if height < 1 {
        height = 1
    }
    return Size{Width: width, Height: height}
}

// Area returns total number of cells (width * height).
func (s Size) Area() int {
    return s.Width * s.Height
}
```

#### Cell

```go
// domain/value/cell.go
package value

// Cell represents a single terminal cell with content and visual width.
// This is an immutable value object.
type Cell struct {
    content string  // Grapheme cluster (not rune!)
    width   int     // Visual width in terminal columns
}

// NewCell creates a cell from a grapheme cluster.
func NewCell(content string, width int) Cell {
    return Cell{
        content: content,
        width:   width,
    }
}

// Content returns the grapheme cluster content.
func (c Cell) Content() string {
    return c.content
}

// Width returns visual width in terminal columns.
func (c Cell) Width() int {
    return c.width
}

// IsEmpty returns true if cell has no content.
func (c Cell) IsEmpty() bool {
    return c.content == "" || c.content == " "
}
```

### Entity: RawMode

```go
// domain/model/rawmode.go
package model

import (
    "errors"
    "syscall"
)

// RawMode represents terminal raw mode state (entity with lifecycle).
type RawMode struct {
    enabled       bool
    originalState interface{}  // Platform-specific (syscall.Termios on Unix)
}

// NewRawMode creates raw mode entity with original state.
func NewRawMode(originalState interface{}) *RawMode {
    return &RawMode{
        enabled:       false,
        originalState: originalState,
    }
}

// IsEnabled returns true if raw mode is active.
func (r *RawMode) IsEnabled() bool {
    return r.enabled
}

// OriginalState returns the saved terminal state for restoration.
func (r *RawMode) OriginalState() interface{} {
    return r.originalState
}

// Enable marks raw mode as enabled (immutable).
func (r *RawMode) Enable() *RawMode {
    newMode := *r
    newMode.enabled = true
    return &newMode
}

// Disable marks raw mode as disabled (immutable).
func (r *RawMode) Disable() *RawMode {
    newMode := *r
    newMode.enabled = false
    return &newMode
}
```

### Domain Services

#### CapabilitiesDetector

```go
// domain/service/capabilities_detector.go
package service

import (
    "os"
    "github.com/phoenix-tui/phoenix/core/domain/value"
)

// CapabilitiesDetector detects terminal capabilities from environment.
// This is a domain service because it embodies business rules.
type CapabilitiesDetector struct{}

// NewCapabilitiesDetector creates a capabilities detector.
func NewCapabilitiesDetector() *CapabilitiesDetector {
    return &CapabilitiesDetector{}
}

// Detect analyzes environment and returns terminal capabilities.
func (cd *CapabilitiesDetector) Detect() *value.Capabilities {
    term := os.Getenv("TERM")
    colorterm := os.Getenv("COLORTERM")

    // Detect ANSI support
    ansi := cd.detectANSI(term)

    // Detect color depth
    colors := cd.detectColorDepth(term, colorterm)

    // Detect mouse support (requires ANSI)
    mouse := ansi

    // Detect alt screen (requires ANSI)
    altScreen := ansi

    // Detect cursor control (requires ANSI)
    cursor := ansi

    return value.NewCapabilities(ansi, colors, mouse, altScreen, cursor)
}

// detectANSI checks if terminal supports ANSI escape sequences.
func (cd *CapabilitiesDetector) detectANSI(term string) bool {
    // Business rule: these terminals support ANSI
    switch term {
    case "xterm", "xterm-256color", "xterm-truecolor":
        return true
    case "screen", "screen-256color":
        return true
    case "tmux", "tmux-256color":
        return true
    case "vt100", "vt220":
        return true
    case "dumb":
        return false  // Explicitly no ANSI
    default:
        return term != ""  // Assume yes if TERM is set
    }
}

// detectColorDepth determines color support level.
func (cd *CapabilitiesDetector) detectColorDepth(term, colorterm string) value.ColorDepth {
    // Business rule: COLORTERM=truecolor → 24-bit
    if colorterm == "truecolor" || colorterm == "24bit" {
        return value.ColorDepthTrueColor
    }

    // Business rule: TERM contains "256color" → 256 colors
    if contains(term, "256color") {
        return value.ColorDepth256
    }

    // Business rule: TERM contains "color" → 8 colors
    if contains(term, "color") {
        return value.ColorDepth8
    }

    // Business rule: specific terminals
    switch term {
    case "xterm", "screen", "tmux":
        return value.ColorDepth8
    case "dumb", "":
        return value.ColorDepthNone
    default:
        return value.ColorDepth8  // Conservative default
    }
}

func contains(s, substr string) bool {
    return len(s) >= len(substr) &&
           s[:len(substr)] == substr ||
           (len(s) > len(substr) && s[len(s)-len(substr):] == substr)
}
```

#### SignalHandler

```go
// domain/service/signal_handler.go
package service

import (
    "os"
    "os/signal"
)

// SignalType represents different signal types.
type SignalType int

const (
    SignalInterrupt SignalType = iota  // Ctrl+C (SIGINT)
    SignalResize                        // Window resize (SIGWINCH)
)

// SignalEvent represents a signal event.
type SignalEvent struct {
    Type SignalType
    Size value.Size  // For resize events
}

// SignalHandler handles OS signals (domain service).
type SignalHandler struct {
    signals chan os.Signal
    events  chan SignalEvent
}

// NewSignalHandler creates a signal handler.
func NewSignalHandler() *SignalHandler {
    return &SignalHandler{
        signals: make(chan os.Signal, 1),
        events:  make(chan SignalEvent, 10),
    }
}

// Start begins listening for signals.
func (sh *SignalHandler) Start() {
    signal.Notify(sh.signals, os.Interrupt, syscall.SIGWINCH)

    go func() {
        for sig := range sh.signals {
            event := sh.convertSignal(sig)
            sh.events <- event
        }
    }()
}

// Events returns channel for signal events.
func (sh *SignalHandler) Events() <-chan SignalEvent {
    return sh.events
}

// Stop stops listening for signals.
func (sh *SignalHandler) Stop() {
    signal.Stop(sh.signals)
    close(sh.signals)
    close(sh.events)
}

// convertSignal converts OS signal to domain event.
func (sh *SignalHandler) convertSignal(sig os.Signal) SignalEvent {
    if sig == os.Interrupt {
        return SignalEvent{Type: SignalInterrupt}
    }
    // For SIGWINCH, would need to query terminal size
    return SignalEvent{Type: SignalResize}
}
```

---

## 🔧 Infrastructure Layer

### Platform Abstraction

Platform-specific code lives in `infrastructure/platform/` with build tags:

```go
// infrastructure/platform/terminal.go
package platform

import "github.com/phoenix-tui/phoenix/core/domain/value"

// TerminalOperations defines platform-specific operations.
type TerminalOperations interface {
    // GetSize returns current terminal size.
    GetSize() (value.Size, error)

    // EnableRawMode puts terminal into raw mode.
    EnableRawMode() (originalState interface{}, err error)

    // DisableRawMode restores terminal to original state.
    DisableRawMode(originalState interface{}) error
}
```

**Unix Implementation**:

```go
// infrastructure/platform/terminal_unix.go
//go:build unix

package platform

import (
    "os"
    "syscall"
    "unsafe"
    "github.com/phoenix-tui/phoenix/core/domain/value"
)

// UnixTerminal implements TerminalOperations for Unix systems.
type UnixTerminal struct{}

// NewTerminal creates platform-specific terminal operations.
func NewTerminal() TerminalOperations {
    return &UnixTerminal{}
}

// GetSize returns terminal size using TIOCGWINSZ ioctl.
func (ut *UnixTerminal) GetSize() (value.Size, error) {
    var ws struct {
        Row    uint16
        Col    uint16
        Xpixel uint16
        Ypixel uint16
    }

    _, _, errno := syscall.Syscall(
        syscall.SYS_IOCTL,
        uintptr(os.Stdout.Fd()),
        uintptr(syscall.TIOCGWINSZ),
        uintptr(unsafe.Pointer(&ws)),
    )

    if errno != 0 {
        return value.Size{Width: 80, Height: 24}, errno
    }

    return value.NewSize(int(ws.Col), int(ws.Row)), nil
}

// EnableRawMode puts terminal into raw mode.
func (ut *UnixTerminal) EnableRawMode() (interface{}, error) {
    fd := int(os.Stdin.Fd())

    // Get current terminal settings
    var original syscall.Termios
    if err := termiosGet(fd, &original); err != nil {
        return nil, err
    }

    // Create new settings for raw mode
    raw := original

    // Input flags: disable BREAK, CR-to-NL, input parity check, strip 8th bit
    raw.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON

    // Output flags: disable post-processing
    raw.Oflag &^= syscall.OPOST

    // Control flags: set 8-bit chars
    raw.Cflag |= syscall.CS8

    // Local flags: disable echo, canonical mode, extended functions, signals
    raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG

    // Set minimum bytes to read and timeout
    raw.Cc[syscall.VMIN] = 1
    raw.Cc[syscall.VTIME] = 0

    // Apply raw mode settings
    if err := termiosSet(fd, &raw); err != nil {
        return nil, err
    }

    return original, nil
}

// DisableRawMode restores terminal to original state.
func (ut *UnixTerminal) DisableRawMode(originalState interface{}) error {
    original, ok := originalState.(syscall.Termios)
    if !ok {
        return errors.New("invalid original state type")
    }

    fd := int(os.Stdin.Fd())
    return termiosSet(fd, &original)
}

func termiosGet(fd int, termios *syscall.Termios) error {
    _, _, errno := syscall.Syscall(
        syscall.SYS_IOCTL,
        uintptr(fd),
        uintptr(syscall.TCGETS),
        uintptr(unsafe.Pointer(termios)),
    )
    if errno != 0 {
        return errno
    }
    return nil
}

func termiosSet(fd int, termios *syscall.Termios) error {
    _, _, errno := syscall.Syscall(
        syscall.SYS_IOCTL,
        uintptr(fd),
        uintptr(syscall.TCSETS),
        uintptr(unsafe.Pointer(termios)),
    )
    if errno != 0 {
        return errno
    }
    return nil
}
```

**Windows Implementation**:

```go
// infrastructure/platform/terminal_windows.go
//go:build windows

package platform

import (
    "golang.org/x/sys/windows"
    "github.com/phoenix-tui/phoenix/core/domain/value"
)

// WindowsTerminal implements TerminalOperations for Windows.
type WindowsTerminal struct{}

// NewTerminal creates platform-specific terminal operations.
func NewTerminal() TerminalOperations {
    return &WindowsTerminal{}
}

// GetSize returns terminal size using GetConsoleScreenBufferInfo.
func (wt *WindowsTerminal) GetSize() (value.Size, error) {
    handle := windows.Handle(os.Stdout.Fd())

    var info windows.ConsoleScreenBufferInfo
    if err := windows.GetConsoleScreenBufferInfo(handle, &info); err != nil {
        return value.Size{Width: 80, Height: 24}, err
    }

    width := int(info.Window.Right - info.Window.Left + 1)
    height := int(info.Window.Bottom - info.Window.Top + 1)

    return value.NewSize(width, height), nil
}

// EnableRawMode enables virtual terminal processing on Windows.
func (wt *WindowsTerminal) EnableRawMode() (interface{}, error) {
    handle := windows.Handle(os.Stdin.Fd())

    // Get current console mode
    var originalMode uint32
    if err := windows.GetConsoleMode(handle, &originalMode); err != nil {
        return nil, err
    }

    // Set raw mode: disable line input, echo input, processed input
    rawMode := originalMode
    rawMode &^= windows.ENABLE_LINE_INPUT | windows.ENABLE_ECHO_INPUT | windows.ENABLE_PROCESSED_INPUT

    if err := windows.SetConsoleMode(handle, rawMode); err != nil {
        return nil, err
    }

    return originalMode, nil
}

// DisableRawMode restores original console mode.
func (wt *WindowsTerminal) DisableRawMode(originalState interface{}) error {
    originalMode, ok := originalState.(uint32)
    if !ok {
        return errors.New("invalid original state type")
    }

    handle := windows.Handle(os.Stdin.Fd())
    return windows.SetConsoleMode(handle, originalMode)
}
```

---

## 📦 API Layer

### Public Interface

```go
// api/core.go
package core

import (
    "github.com/phoenix-tui/phoenix/core/domain/model"
    "github.com/phoenix-tui/phoenix/core/domain/service"
    "github.com/phoenix-tui/phoenix/core/domain/value"
    "github.com/phoenix-tui/phoenix/core/infrastructure/platform"
)

// Terminal is the main public API for phoenix/core.
type Terminal struct {
    model      *model.Terminal
    ops        platform.TerminalOperations
    detector   *service.CapabilitiesDetector
    signals    *service.SignalHandler
}

// NewTerminal creates a new terminal with detected capabilities.
func NewTerminal() (*Terminal, error) {
    detector := service.NewCapabilitiesDetector()
    caps := detector.Detect()

    termModel := model.NewTerminal(caps)
    ops := platform.NewTerminal()

    // Get initial size
    size, err := ops.GetSize()
    if err == nil {
        termModel = termModel.WithSize(size)
    }

    return &Terminal{
        model:    termModel,
        ops:      ops,
        detector: detector,
        signals:  service.NewSignalHandler(),
    }, nil
}

// Capabilities returns terminal capabilities.
func (t *Terminal) Capabilities() *value.Capabilities {
    return t.model.Capabilities()
}

// Size returns current terminal size.
func (t *Terminal) Size() value.Size {
    return t.model.Size()
}

// EnableRawMode puts terminal into raw mode.
func (t *Terminal) EnableRawMode() error {
    originalState, err := t.ops.EnableRawMode()
    if err != nil {
        return err
    }

    rawMode := model.NewRawMode(originalState).Enable()
    t.model = t.model.WithRawMode(rawMode)
    return nil
}

// DisableRawMode restores terminal to original state.
func (t *Terminal) DisableRawMode() error {
    if !t.model.IsRawMode() {
        return nil  // Already disabled
    }

    rawMode := t.model.RawMode()
    if err := t.ops.DisableRawMode(rawMode.OriginalState()); err != nil {
        return err
    }

    t.model = t.model.WithoutRawMode()
    return nil
}

// StartSignalHandling begins listening for OS signals.
func (t *Terminal) StartSignalHandling() {
    t.signals.Start()
}

// Signals returns channel for signal events.
func (t *Terminal) Signals() <-chan service.SignalEvent {
    return t.signals.Events()
}

// Close cleans up resources (ensures raw mode is disabled).
func (t *Terminal) Close() error {
    t.signals.Stop()
    return t.DisableRawMode()
}
```

---

## 🧪 Testing Strategy

### Domain Tests (95%+ coverage target)

```go
// domain/model/terminal_test.go
package model_test

func TestTerminal_WithSize(t *testing.T) {
    // Arrange
    caps := value.NewCapabilities(true, value.ColorDepth256, true, true, true)
    term := model.NewTerminal(caps)
    newSize := value.NewSize(100, 50)

    // Act
    newTerm := term.WithSize(newSize)

    // Assert
    assert.Equal(t, newSize, newTerm.Size())
    assert.NotEqual(t, newSize, term.Size())  // Original unchanged
}

func TestTerminal_RawMode(t *testing.T) {
    caps := value.NewCapabilities(true, value.ColorDepth256, true, true, true)
    term := model.NewTerminal(caps)

    // Initially not in raw mode
    assert.False(t, term.IsRawMode())

    // Enable raw mode
    rawMode := model.NewRawMode("original").Enable()
    termRaw := term.WithRawMode(rawMode)

    assert.True(t, termRaw.IsRawMode())
    assert.False(t, term.IsRawMode())  // Original unchanged
}
```

### Infrastructure Tests (80%+ coverage)

```go
// infrastructure/platform/terminal_unix_test.go
//go:build unix

package platform_test

func TestUnixTerminal_GetSize(t *testing.T) {
    term := platform.NewTerminal()

    size, err := term.GetSize()

    assert.NoError(t, err)
    assert.Greater(t, size.Width, 0)
    assert.Greater(t, size.Height, 0)
}

func TestUnixTerminal_RawMode(t *testing.T) {
    term := platform.NewTerminal()

    // Enable
    original, err := term.EnableRawMode()
    assert.NoError(t, err)
    assert.NotNil(t, original)

    // Disable
    err = term.DisableRawMode(original)
    assert.NoError(t, err)
}
```

### API Tests (85%+ coverage)

```go
// api/core_test.go
package core_test

func TestTerminal_Lifecycle(t *testing.T) {
    term, err := core.NewTerminal()
    require.NoError(t, err)
    defer term.Close()

    // Check capabilities
    caps := term.Capabilities()
    assert.NotNil(t, caps)

    // Check size
    size := term.Size()
    assert.Greater(t, size.Width, 0)

    // Enable raw mode
    err = term.EnableRawMode()
    assert.NoError(t, err)

    // Disable raw mode
    err = term.DisableRawMode()
    assert.NoError(t, err)
}
```

---

## 📊 Implementation Order

### Week 3: Foundation

**Day 1-2: Domain Layer**
1. ✅ Value objects (Position, Size, Cell, Capabilities)
2. ✅ Domain models (Terminal, RawMode)
3. ✅ Domain services (CapabilitiesDetector)
4. ✅ Unit tests (95%+ coverage)

**Day 3-4: Infrastructure Layer**
5. ✅ Platform interface
6. ✅ Unix implementation
7. ✅ Windows implementation (basic)
8. ✅ Integration tests

**Day 5-7: API Layer**
9. ✅ Public API design
10. ✅ API implementation
11. ✅ API tests
12. ✅ Examples

### Week 4: Advanced Features

**Day 1-3: Unicode Support**
1. ✅ Integrate uniseg library
2. ✅ Width calculation service
3. ✅ Grapheme cluster handling
4. ✅ Tests with emoji/CJK

**Day 4-5: Signal Handling**
5. ✅ SignalHandler service
6. ✅ Signal event types
7. ✅ Goroutine-safe channels
8. ✅ Tests with fake signals

**Day 6-7: Polish & Documentation**
9. ✅ ANSI parser (basic)
10. ✅ Documentation
11. ✅ Examples
12. ✅ Benchmarks

---

## 🎯 Success Criteria

- ✅ **95%+ domain test coverage**
- ✅ **80%+ infrastructure test coverage**
- ✅ **85%+ API test coverage**
- ✅ **Zero external dependencies** (except uniseg)
- ✅ **Cross-platform** (Unix + Windows)
- ✅ **Raw mode cleanup** guaranteed (defer pattern)
- ✅ **Immutable domain models** (all With* methods)
- ✅ **Rich behavior** (domain knows business rules)

---

## 🚀 Next Steps

1. Start implementing domain layer (value objects first)
2. Add comprehensive tests for each component
3. Implement platform abstraction
4. Create public API
5. Write examples and documentation

**Ready to begin Week 3 implementation!** 🔥

---

*Architecture Version: 1.0*
*Created: 2025-10-14*
*Status: Week 3-4 Implementation Guide*
