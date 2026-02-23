# phoenix/clipboard

Cross-platform clipboard support for Phoenix TUI Framework with automatic OSC 52 detection for SSH sessions.

## Features

- **Cross-platform Support**: Windows, macOS, Linux (X11 & Wayland)
- **OSC 52 Support**: Automatic clipboard sync over SSH connections
- **Smart Provider Selection**: Auto-detects best clipboard method
- **DDD Architecture**: Clean, testable, extensible design
- **Type-Safe**: Leverages Go 1.25+ features
- **Zero External TUI Dependencies**: Built from scratch for Phoenix

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/phoenix-tui/phoenix/clipboard/api"
)

func main() {
    // Write to clipboard
    err := api.Write("Hello from Phoenix!")
    if err != nil {
        log.Fatal(err)
    }

    // Read from clipboard
    text, err := api.Read()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(text) // Output: Hello from Phoenix!
}
```

## Installation

```bash
go get github.com/phoenix-tui/phoenix/clipboard
```

## Usage

### Basic Usage

```go
// Using global convenience functions
api.Write("text to copy")
text, err := api.Read()

// Check availability
if api.IsAvailable() {
    fmt.Println("Clipboard is available")
}

// Get provider name
fmt.Println(api.GetProviderName()) // e.g., "Windows Native", "OSC52"
```

### Instance-Based Usage

```go
// Create clipboard instance
clipboard, err := api.New()
if err != nil {
    log.Fatal(err)
}

// Write text
err = clipboard.Write("Hello, World!")
if err != nil {
    log.Fatal(err)
}

// Read text
text, err := clipboard.Read()
if err != nil {
    log.Fatal(err)
}

fmt.Println(text)
```

### Custom Configuration

```go
// Using builder pattern
clipboard, err := api.NewBuilder().
    WithOSC52(true).                    // Enable OSC 52
    WithOSC52Timeout(5*time.Second).    // Set timeout
    WithNative(true).                    // Enable native clipboard
    Build()

if err != nil {
    log.Fatal(err)
}

err = clipboard.Write("Custom config")
```

### OSC 52 for SSH Sessions

OSC 52 automatically enables clipboard sync over SSH:

```go
clipboard, err := api.New()

// Check if running in SSH session
if clipboard.IsSSH() {
    fmt.Println("OSC 52 will be used for clipboard sync")
}

// Write - automatically uses OSC 52 in SSH
clipboard.Write("Text syncs to local clipboard!")
```

### Working with Domain Models

```go
import (
    "github.com/phoenix-tui/phoenix/clipboard/domain/model"
    "github.com/phoenix-tui/phoenix/clipboard/domain/value"
)

// Create text content
content, err := model.NewTextContent("Hello")

// Check content properties
fmt.Println(content.MIMEType())  // text/plain
fmt.Println(content.Encoding())  // utf-8
fmt.Println(content.Size())      // 5
fmt.Println(content.IsText())    // true

// Create binary content
binaryContent, err := model.NewBinaryContent([]byte{0x00, 0x01})

// Transform content (immutable)
htmlContent := content.WithMIMEType(value.MIMETypeHTML)
base64Content := htmlContent.WithEncoding(value.EncodingBase64)
```

## Architecture

Phoenix clipboard follows Domain-Driven Design (DDD) with Hexagonal Architecture:

```
clipboard/
├── domain/              # Pure business logic (95%+ test coverage)
│   ├── model/          # ClipboardContent aggregate
│   ├── value/          # MIMEType, Encoding value objects
│   └── service/        # Provider interface, ClipboardService
├── infrastructure/      # Technical implementations (80%+ coverage)
│   ├── osc52/          # OSC 52 provider
│   ├── native/         # Platform-specific providers
│   │   ├── clipboard_windows.go   # Windows user32.dll
│   │   ├── clipboard_darwin.go    # macOS pbcopy/pbpaste
│   │   └── clipboard_linux.go     # Linux xclip/wl-clipboard
│   └── platform/       # Platform detection
├── application/         # Use cases (90%+ coverage)
│   └── clipboard_manager.go
└── api/                # Public API (100% coverage)
    └── clipboard.go
```

### Provider Selection

Phoenix uses a prioritized fallback chain:

1. **SSH Session Detected** → OSC 52 (primary)
2. **Native Platform Clipboard**:
   - Windows: `user32.dll` APIs
   - macOS: `pbcopy`/`pbpaste`
   - Linux: `xclip`, `xsel`, or `wl-clipboard`
3. **OSC 52 Fallback** (for non-SSH terminals that support it)

### Why OSC 52?

OSC 52 is an ANSI escape sequence that allows terminal applications to set the system clipboard. This is essential for SSH sessions where the remote server doesn't have direct access to the local clipboard.

**Benefits**:
- Works over SSH connections
- No X11 forwarding required
- Syncs remote clipboard → local machine
- Supported by modern terminals (iTerm2, Windows Terminal, xterm, tmux, etc.)

**Limitations**:
- Read operations not widely supported (write-only in most terminals)
- Requires terminal support
- 5-second timeout by default

## Platform Support

### Windows
- Native clipboard via `user32.dll` (OpenClipboard, GetClipboardData, SetClipboardData)
- UTF-16 encoding support
- Always available

### macOS
- Uses `pbcopy` and `pbpaste` commands
- Requires macOS 10.0+
- Always available on macOS

### Linux
- **Wayland**: `wl-copy` / `wl-paste` (preferred)
- **X11**: `xclip` or `xsel` (fallback)
- Auto-detects available clipboard tool
- Requires external tool installation

Install clipboard tools on Linux:
```bash
# Wayland
sudo apt install wl-clipboard

# X11
sudo apt install xclip
# or
sudo apt install xsel
```

## Testing

Run all tests:
```bash
cd clipboard
go test ./...
```

Run with coverage:
```bash
go test -cover ./...
```

Expected coverage:
- Domain layer: 95%+
- Application layer: 90%+
- Infrastructure layer: 80%+
- API layer: 100%

## Examples

### Basic Example
```bash
cd examples/basic
go run main.go
```

### OSC 52 Example
```bash
cd examples/osc52
go run main.go
```

### Multiple Formats Example
```bash
cd examples/formats
go run main.go
```

## Key Differentiators vs Other Libraries

### vs Charm's clipboard
1. **Auto-detects OSC 52 support** (Charm requires manual configuration)
2. **Intelligent fallback chain** (try OSC 52 → native → error)
3. **DDD architecture** (testable, extensible, maintainable)
4. **Type-safe domain models** (ClipboardContent with behavior, not just data)
5. **Cross-platform native APIs** (not just shell commands)
6. **Zero Charm dependencies** (built from scratch for Phoenix)

### vs atotto/clipboard
1. **OSC 52 support** (essential for SSH sessions)
2. **Domain-driven design** (clean architecture)
3. **Smart provider detection** (automatic fallback)
4. **Rich domain models** (MIME types, encodings, transformations)
5. **Builder pattern** (flexible configuration)

## Performance Characteristics

### Memory Management
- Efficient binary content handling (zero-copy where possible)
- Content validation prevents empty/invalid data
- Builder pattern reduces object creation overhead
- Minimal allocations in hot paths (Read/Write operations)

### Best Practices for Large Content
1. **Streaming**: For files > 1MB, consider chunking before clipboard operations
2. **Validation**: Check content size before writing (some platforms have limits)
3. **Timeout Configuration**: Adjust OSC52 timeout for slow connections
   ```go
   clipboard, _ := NewBuilder().
       WithOSC52Timeout(10 * time.Second).
       Build()
   ```
4. **Provider Selection**: Use native-only for large local content (faster than OSC52)
   ```go
   clipboard, _ := NewBuilder().
       WithOSC52(false).  // Disable OSC52 for local-only usage
       WithNative(true).
       Build()
   ```

### Thread Safety
- **Clipboard instance**: Safe for concurrent Read/Write from multiple goroutines
- **Global functions**: Thread-safe lazy initialization
- **Provider implementations**: Handle platform-specific locking internally
- **History**: Concurrent access protected by internal synchronization

### Performance Notes
- Native clipboard access: < 1ms typical latency
- OSC 52 operations: 5-50ms depending on terminal and content size
- History tracking: O(1) add, O(n) for GetHistory() where n = entry count
- Content validation: O(1) for size check, O(n) for encoding validation

## Design Principles

### 1. Rich Domain Models
```go
// NOT anemic (just data):
type Content struct {
    Data []byte
}

// Phoenix rich model (data + behavior):
type ClipboardContent struct {
    data     []byte
    mimeType MIMEType
    encoding Encoding
}

func (c *ClipboardContent) Text() (string, error) { ... }
func (c *ClipboardContent) WithEncoding(enc Encoding) *ClipboardContent { ... }
```

### 2. Immutability
All domain objects are immutable. Transformations return new instances:
```go
original := content.WithMIMEType(MIMETypePlainText)
modified := original.WithEncoding(EncodingBase64)
// original unchanged
```

### 3. Hexagonal Architecture
Domain layer defines interfaces (`Provider`), infrastructure implements them. Easy to:
- Test with mocks
- Add new clipboard providers
- Swap implementations without changing domain logic

### 4. Type Safety
Leverages Go 1.25+ features:
- Value objects for MIME types and encodings
- Compile-time safety
- No string-based "magic values"

## Integration with Phoenix Components

Phoenix clipboard integrates seamlessly with other Phoenix components:

```go
// In a TextInput component (phoenix/components/input)
func (m *TextInputModel) handlePaste() tea.Cmd {
    return func() tea.Msg {
        text, err := clipboard.Read()
        if err != nil {
            return PasteErrorMsg{err}
        }
        return PasteMsg{text}
    }
}
```

## Roadmap

### Future Enhancements
- [ ] Rich text format support
- [ ] Image clipboard support
- [ ] Custom MIME types
- [ ] Clipboard monitoring (watch for changes)
- [ ] Async clipboard operations
- [ ] Clipboard history

## Contributing

Phoenix is part of the Phoenix TUI Framework project. See the main repository for contribution guidelines.

## License

MIT License - see LICENSE file for details

## Credits

Part of [Phoenix TUI Framework](https://github.com/phoenix-tui/phoenix) - Next-generation TUI framework for Go.

Inspired by but independent of:
- Charm's clipboard (atotto/clipboard wrapper)
- atotto/clipboard (cross-platform clipboard library)
- golang-design/clipboard (advanced clipboard features)

Built from scratch with DDD architecture and modern Go practices.

## Support

- GitHub Issues: [phoenix-tui/phoenix/issues](https://github.com/phoenix-tui/phoenix/issues)
- Documentation: [phoenix-tui.dev](https://phoenix-tui.dev)
- Examples: `clipboard/examples/`

---

*Phoenix TUI Framework - Building the future of terminal interfaces*
