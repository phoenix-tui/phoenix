# Phoenix Render

High-performance differential rendering engine for Phoenix TUI Framework.

**Module**: `github.com/phoenix-tui/phoenix/render`

## Features

- **Differential Rendering**: Only renders changed cells (significantly faster than full redraws)
- **ANSI Optimization**: Batches escape sequences for minimal output
- **Zero-Allocation Hot Paths**: Buffer pooling for garbage-free rendering
- **Unicode Perfect**: Correct grapheme cluster handling (emoji, CJK, etc.)
- **Thread-Safe**: Safe for concurrent rendering
- **Performance**: Designed for 60 FPS interactive applications

## Installation

```bash
go get github.com/phoenix-tui/phoenix/render@latest
```

## Quick Start

```go
package main

import (
    "os"
    render "github.com/phoenix-tui/phoenix/render/api"
)

func main() {
    // Create renderer
    renderer := render.New(80, 24, os.Stdout)
    defer renderer.Close()

    // Create buffer
    buf := renderer.Buffer()
    defer buf.Release()

    // Write content
    buf.SetString(0, 0, "Hello, Phoenix!", render.FgRed.WithBold(true))

    // Render (differential)
    if err := renderer.Render(buf); err != nil {
        panic(err)
    }
}
```

## API Reference

See full API documentation in this file.

## Performance

Target: 60 FPS (< 16.67ms per frame)

Run benchmarks:
```bash
go test -bench=. -benchmem ./benchmarks
```

## Testing

```bash
# Run all tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## License

MIT License - see [LICENSE](../LICENSE) for details.
