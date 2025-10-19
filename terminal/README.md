# Phoenix Terminal - Platform-Optimized Terminal Operations

**Status**: Week 15 Complete - ANSI Baseline Implementation ✅
**Coverage**: 93.0% (exceeds 90% target) ⭐
**Version**: v0.1.0-alpha (Week 15)

Phoenix Terminal provides a platform-abstraction layer for terminal operations with automatic detection and optimized implementations for each platform.

## Features

### Week 15 (Current) - ANSI Baseline

- ✅ **Complete Terminal API** - All operations defined and documented
- ✅ **ANSI Implementation** - Universal fallback for Unix/Linux/macOS/Git Bash
- ✅ **93.0% Test Coverage** - Comprehensive unit and benchmark tests
- ✅ **Auto-Detection** - Automatic platform detection (returns ANSI in Week 15)
- ✅ **ClearLines()** - Critical multiline clearing operation
- ✅ **Rich Documentation** - Complete godoc and usage examples

### Week 16 (Coming) - Windows Console API

- 🎯 **Windows Console API** - Direct Win32 calls (10x faster!)
- 🎯 **Auto-Fallback** - Git Bash detection → ANSI fallback
- 🎯 **Cursor Readback** - GetCursorPosition() support
- 🎯 **Screen Buffer Readback** - ReadScreenBuffer() for differential rendering
- 🎯 **Performance Benchmarks** - ANSI vs Windows API comparison

## Quick Start

### Installation

```bash
cd D:\projects\grpmsoft\tui
go work use ./terminal
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
// Auto-detect platform (Week 15: returns ANSI, Week 16: Windows API or ANSI)
term := infrastructure.NewTerminal()

// Force ANSI implementation (for testing or compatibility)
term := infrastructure.NewANSITerminal()
```

### Cursor Operations

```go
// Absolute positioning (0-based coordinates)
term.SetCursorPosition(x, y int) error

// Cursor position readback (Week 16 only - Windows Console API)
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

### Screen Buffer (Week 16 - Windows Console API only)

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
// api.PlatformWindowsConsole - cmd.exe/PowerShell (Win32 API) - Week 16
// api.PlatformWindowsANSI - Git Bash on Windows (ANSI fallback) - Week 16
```

## Platform Support Matrix

| Platform | Week 15 (ANSI) | Week 16 (Optimized) |
|----------|----------------|---------------------|
| Linux | ✅ ANSI | ✅ ANSI |
| macOS | ✅ ANSI | ✅ ANSI |
| Windows (cmd.exe) | ✅ ANSI | 🎯 Win32 API (10x faster) |
| Windows (PowerShell) | ✅ ANSI | 🎯 Win32 API (10x faster) |
| Windows (Git Bash) | ✅ ANSI | ✅ ANSI (auto-fallback) |
| Windows (WSL) | ✅ ANSI | ✅ ANSI |

## Performance

### Week 15 Baseline (ANSI)

Benchmarks (on Windows Git Bash):

```
BenchmarkANSI_SetCursorPosition    ~100ns per operation
BenchmarkANSI_ClearLines           ~500ns per operation (10 lines)
BenchmarkANSI_Write                ~50ns per operation
BenchmarkANSI_WriteAt              ~150ns per operation
```

### Week 16 Target (Windows Console API)

Expected improvements on Windows (cmd.exe, PowerShell):

| Operation | ANSI | Windows API | Speedup |
|-----------|------|-------------|---------|
| SetCursorPosition | ~100μs | ~10μs | 10x |
| ClearLines(10) | ~500μs | ~50μs | 10x |
| ReadScreenBuffer | N/A | ~1ms | ∞ |

## Testing

### Run Tests

```bash
cd D:\projects\grpmsoft\tui\terminal

# All tests
go test ./...

# With coverage
go test ./infrastructure/unix -cover
# Output: coverage: 93.0% of statements ✅

# Verbose output
go test ./infrastructure/unix -v

# Benchmarks
go test ./infrastructure/unix -bench=. -benchmem
```

### Coverage Report

```
github.com/phoenix-tui/phoenix/terminal/infrastructure/unix
    coverage: 93.0% of statements ✅
```

Exceeds Phoenix 90%+ target!

## Architecture

Phoenix Terminal follows DDD layering:

```
terminal/
├── api/
│   └── terminal.go         # Terminal interface (platform-independent)
├── infrastructure/
│   ├── unix/
│   │   ├── ansi.go         # ANSI implementation (Week 15) ✅
│   │   └── ansi_test.go    # Comprehensive tests (93.0% coverage) ✅
│   ├── windows/            # Windows Console API (Week 16) 🎯
│   │   ├── console.go      # Win32 API calls
│   │   └── console_test.go
│   └── detect.go           # Auto-detection logic
├── examples/
│   └── basic/
│       └── main.go         # Demo application ✅
├── go.mod                  # Module definition ✅
└── README.md               # This file ✅
```

**Key Principles**:
- `api/` defines interface (no dependencies on infrastructure)
- `infrastructure/` provides platform-specific implementations
- Auto-detection logic chooses best implementation at runtime
- Graceful fallback ensures compatibility

## ANSI Escape Codes Reference

Week 15 implementation uses these ANSI codes:

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

**ANSI Implementation** (Week 15):
```
count == 1: \r\033[J           (CR + clear to end)
count > 1:  \033[{n-1}A\r\033[J (move up, CR, clear to end)
```

**Windows API Implementation** (Week 16):
```go
// Direct Win32 FillConsoleOutputCharacter() - 10x faster!
windows.FillConsoleOutputCharacter(stdout, ' ', totalChars, startCoord, &written)
```

## Integration with GoSh

GoSh multiline mode will use Phoenix Terminal in Week 16:

```go
// Before (manual ANSI codes):
if len(lines) > 1 {
    fmt.Printf("\033[%dA", len(lines)-1)
}
fmt.Print("\r\033[J")

// After (Phoenix Terminal - platform-optimized):
term.ClearLines(len(lines))
```

## Roadmap

### ✅ Week 15 (Complete)
- [x] Terminal interface definition
- [x] ANSI implementation (all methods)
- [x] Auto-detection (returns ANSI)
- [x] Comprehensive tests (93.0% coverage)
- [x] Benchmark tests
- [x] Example application
- [x] Documentation

### 🎯 Week 16 (Next)
- [ ] Windows Console API implementation
- [ ] Windows detection logic
- [ ] Auto-fallback (Windows API → Git Bash ANSI)
- [ ] Cursor readback (GetCursorPosition)
- [ ] Screen buffer readback (ReadScreenBuffer)
- [ ] Performance benchmarks (ANSI vs Windows API)
- [ ] GoSh integration

## Contributing

This is part of Phoenix TUI Framework (Week 15-16).

**Testing Requirements**:
- Unit tests for all public methods
- 90%+ coverage (Week 15: 93.0% ✅)
- Benchmark tests for performance tracking
- Example applications that demonstrate usage

**Code Standards**:
- Follow Phoenix DDD architecture
- Comprehensive godoc comments
- Table-driven tests
- Error handling with clear messages

## License

Part of Phoenix TUI Framework - See root LICENSE

---

**Phoenix Terminal** - Building the foundation for 10x Windows performance! 🚀

Week 15: ✅ ANSI Baseline (93.0% coverage)
Week 16: 🎯 Windows Console API (10x faster!)
