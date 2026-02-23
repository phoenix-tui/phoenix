# Phoenix Terminal - Platform-Optimized Terminal Operations

**Module**: `github.com/phoenix-tui/phoenix/terminal`

Phoenix Terminal provides a platform-abstraction layer for terminal operations with automatic detection and optimized implementations for each platform.

## Features

- **Complete Terminal API** - All operations defined and documented
- **ANSI Implementation** - Universal fallback for Unix/Linux/macOS/Git Bash
- **Auto-Detection** - Automatic platform detection
- **ClearLines()** - Critical multiline clearing operation
- **Rich Documentation** - Complete godoc and usage examples
- **Extensive test coverage** - Comprehensive unit and benchmark tests

### Planned

- **Windows Console API** - Direct Win32 calls for faster performance
- **Auto-Fallback** - Git Bash detection with ANSI fallback
- **Cursor Readback** - GetCursorPosition() support
- **Screen Buffer Readback** - ReadScreenBuffer() for differential rendering

## Quick Start

### Installation

```bash
go get github.com/phoenix-tui/phoenix/terminal
```

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/phoenix-tui/phoenix/terminal/infrastructure"
)

func main() {
    // Auto-detect best terminal implementation
    term := infrastructure.NewTerminal()

    // Platform detection
    fmt.Printf("Platform: %s\n", term.Platform())
    fmt.Printf("Supports readback: %v\n", term.SupportsReadback())

    // Cursor operations
    term.HideCursor()
    term.SetCursorPosition(10, 5)
    term.Write("Hello, Phoenix!")
    term.ShowCursor()

    // Multiline clearing (critical for shells like GoSh)
    term.ClearLines(3) // Clear 3 lines
}
```

### Run Example

```bash
cd terminal/examples/basic
go run main.go
```

## API Reference

### Constructors

```go
// Auto-detect platform (returns best available implementation)
term := infrastructure.NewTerminal()

// Force ANSI implementation (for testing or compatibility)
term := infrastructure.NewANSITerminal()
```

### Cursor Operations

```go
// Absolute positioning (0-based coordinates)
term.SetCursorPosition(x, y int) error

// Cursor position readback (Windows Console API only)
x, y, err := term.GetCursorPosition() // Returns error on ANSI

// Relative movements
term.MoveCursorUp(n int) error
term.MoveCursorDown(n int) error
term.MoveCursorLeft(n int) error
term.MoveCursorRight(n int) error

// Save/restore position
term.SaveCursorPosition() error
term.RestoreCursorPosition() error
```

### Cursor Visibility & Style

```go
// Visibility
term.HideCursor() error
term.ShowCursor() error

// Style (Block, Underline, Bar)
term.SetCursorStyle(api.CursorBlock) error
term.SetCursorStyle(api.CursorUnderline) error
term.SetCursorStyle(api.CursorBar) error
```

### Screen Operations

```go
// Clear operations
term.Clear() error              // Clear entire screen
term.ClearLine() error          // Clear current line
term.ClearFromCursor() error    // Clear from cursor to end

// Clear N lines (CRITICAL for multiline input like GoSh)
term.ClearLines(count int) error
```

### Output

```go
// Write at current position
term.Write(s string) error

// Write at specific position (optimized on Windows Console API)
term.WriteAt(x, y int, s string) error
```

### Screen Buffer (Windows Console API only)

```go
// Read screen buffer for differential rendering
buffer, err := term.ReadScreenBuffer() // [][]rune
// Returns error on ANSI (not supported)
```

### Terminal Info

```go
// Get terminal dimensions
width, height, err := term.Size()

// Color support
colors := term.ColorDepth() // 16, 256, or 16777216 (24-bit RGB)
```

### Capabilities Discovery

```go
// Check platform capabilities
supportsDirectPos := term.SupportsDirectPositioning() // false on ANSI, true on Windows API
supportsReadback := term.SupportsReadback()          // false on ANSI, true on Windows API
supportsTrueColor := term.SupportsTrueColor()        // true if 24-bit RGB

// Platform type
platform := term.Platform()
// api.PlatformUnix - Linux/macOS/Git Bash (ANSI)
// api.PlatformWindowsConsole - cmd.exe/PowerShell (Win32 API)
// api.PlatformWindowsANSI - Git Bash on Windows (ANSI fallback)
```

## Platform Support Matrix

| Platform | ANSI | Optimized (Windows API) |
|----------|------|-------------------------|
| Linux | Supported | N/A |
| macOS | Supported | N/A |
| Windows (cmd.exe) | Supported | Planned (Win32 API) |
| Windows (PowerShell) | Supported | Planned (Win32 API) |
| Windows (Git Bash) | Supported | ANSI (auto-fallback) |
| Windows (WSL) | Supported | N/A |

## Performance

### ANSI Baseline

Benchmarks (on Windows Git Bash):

```
BenchmarkANSI_SetCursorPosition    ~100ns per operation
BenchmarkANSI_ClearLines           ~500ns per operation (10 lines)
BenchmarkANSI_Write                ~50ns per operation
BenchmarkANSI_WriteAt              ~150ns per operation
```

### Windows Console API (Planned)

Expected improvements on Windows (cmd.exe, PowerShell):

| Operation | ANSI | Windows API | Speedup |
|-----------|------|-------------|---------|
| SetCursorPosition | ~100us | ~10us | 10x |
| ClearLines(10) | ~500us | ~50us | 10x |
| ReadScreenBuffer | N/A | ~1ms | N/A |

## Testing

### Run Tests

```bash
cd terminal

# All tests
go test ./...

# With coverage
go test ./infrastructure/unix -cover

# Verbose output
go test ./infrastructure/unix -v

# Benchmarks
go test ./infrastructure/unix -bench=. -benchmem
```

## Architecture

Phoenix Terminal follows DDD layering:

```
terminal/
├── api/
│   └── terminal.go         # Terminal interface (platform-independent)
├── infrastructure/
│   ├── unix/
│   │   ├── ansi.go         # ANSI implementation
│   │   └── ansi_test.go    # Comprehensive tests
│   ├── windows/            # Windows Console API (planned)
│   │   ├── console.go      # Win32 API calls
│   │   └── console_test.go
│   └── detect.go           # Auto-detection logic
├── examples/
│   └── basic/
│       └── main.go         # Demo application
├── go.mod                  # Module definition
└── README.md               # This file
```

**Key Principles**:
- `api/` defines interface (no dependencies on infrastructure)
- `infrastructure/` provides platform-specific implementations
- Auto-detection logic chooses best implementation at runtime
- Graceful fallback ensures compatibility

## ANSI Escape Codes Reference

The ANSI implementation uses these escape codes:

```
Cursor Positioning:
  \033[{row};{col}H         Absolute position (1-based!)
  \033[{n}A                 Move up n lines
  \033[{n}B                 Move down n lines
  \033[{n}C                 Move right n columns
  \033[{n}D                 Move left n columns
  \033[s                    Save cursor position
  \033[u                    Restore cursor position

Cursor Visibility:
  \033[?25h                 Show cursor
  \033[?25l                 Hide cursor
  \033[{n} q                Set cursor style (2=block, 4=underline, 6=bar)

Screen Clearing:
  \033[2J                   Clear entire screen
  \033[H                    Move to home (1,1)
  \033[2K                   Clear current line
  \033[J                    Clear from cursor to end
  \033[{n}A\r\033[J         ClearLines(n) - move up then clear
```

## ClearLines() - Critical Operation

The `ClearLines(count)` operation is essential for multiline input (like GoSh shell):

```go
// Before rendering new multiline content, clear old content
term.ClearLines(oldLineCount)

// Then render new content
term.Write(newContent)
```

**ANSI Implementation**:
```
count == 1: \r\033[J           (CR + clear to end)
count > 1:  \033[{n-1}A\r\033[J (move up, CR, clear to end)
```

**Windows API Implementation** (planned):
```go
// Direct Win32 FillConsoleOutputCharacter() - significantly faster
windows.FillConsoleOutputCharacter(stdout, ' ', totalChars, startCoord, &written)
```

## Integration with GoSh

GoSh multiline mode uses Phoenix Terminal for platform-optimized clearing:

```go
// Before (manual ANSI codes):
if len(lines) > 1 {
    fmt.Printf("\033[%dA", len(lines)-1)
}
fmt.Print("\r\033[J")

// After (Phoenix Terminal - platform-optimized):
term.ClearLines(len(lines))
```

## Contributing

This is part of Phoenix TUI Framework.

**Testing Requirements**:
- Unit tests for all public methods
- High test coverage target
- Benchmark tests for performance tracking
- Example applications that demonstrate usage

**Code Standards**:
- Follow Phoenix DDD architecture
- Comprehensive godoc comments
- Table-driven tests
- Error handling with clear messages

## License

Part of Phoenix TUI Framework - See root LICENSE
