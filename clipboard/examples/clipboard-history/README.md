# Clipboard History Manager

Interactive TUI application demonstrating clipboard history tracking and management capabilities of the Phoenix clipboard library.

## Features

- **History Tracking**: Automatically track all clipboard operations
- **Multiple Formats**: Support for plain text, HTML, and RTF
- **History Navigation**: Browse history with keyboard shortcuts
- **Restore Entries**: Restore any entry from history back to clipboard
- **Memory Monitoring**: View memory usage and entry counts
- **Expiration Management**: Automatic cleanup of expired entries
- **Enable/Disable**: Toggle history tracking on the fly

## Building

```bash
cd clipboard/examples/clipboard-history
go build
```

## Running

```bash
./clipboard-history      # On Linux/macOS
clipboard-history.exe   # On Windows
```

## Keyboard Controls

### Normal Mode

| Key | Action |
|-----|--------|
| `t` | Copy plain text |
| `h` | Copy HTML |
| `r` | Copy RTF |
| `v` | Refresh/view history |
| `R` | Restore selected entry to clipboard |
| `c` | Clear all history |
| `e` | Enable/disable history tracking |
| `m` | Show memory usage and remove expired |
| `↑/k` | Navigate up in history |
| `↓/j` | Navigate down in history |
| `Home/g` | Jump to first entry |
| `End/G` | Jump to last entry |
| `h/?` | Show help |
| `q/Ctrl+C` | Quit |

### Input Mode

When copying text/HTML/RTF:

| Key | Action |
|-----|--------|
| `Enter` | Submit input and copy to clipboard |
| `Esc` | Cancel input |
| `Backspace` | Delete last character |

## Usage Examples

### 1. Basic Text Tracking

```
1. Press 't' to enter text input mode
2. Type: "Hello, World!"
3. Press Enter to copy
4. Press 'v' to view history
5. See your entry in the list
```

### 2. Multiple Format Support

```
1. Press 't' and copy plain text
2. Press 'h' and copy HTML
3. Press 'r' and copy RTF
4. Press 'v' to see all entries with different MIME types
```

### 3. Restoring from History

```
1. Use arrow keys to select an entry
2. Press 'R' to restore it to clipboard
3. The entry is now in your system clipboard
```

### 4. Memory Management

```
1. Press 'm' to view:
   - Total entries
   - Memory usage (bytes)
   - Expired entries removed
```

### 5. Toggle History Tracking

```
1. Press 'e' to disable history
2. Copy some text (not tracked)
3. Press 'e' to re-enable history
4. New copies will be tracked again
```

## Configuration

The example uses these default settings:

- **Max Size**: 100 entries (FIFO eviction when full)
- **Max Age**: 24 hours (entries older than this are expired)
- **Auto-cleanup**: On demand via 'm' key

## Architecture Highlights

This example demonstrates:

1. **Domain-Driven Design**: History tracking using rich domain models
2. **TEA Pattern**: Clean Model-View-Update architecture
3. **Immutability**: History entries are immutable value objects
4. **Thread Safety**: Concurrent-safe history operations
5. **Memory Efficiency**: Copy-on-read, FIFO eviction

## API Usage

```go
// Enable history tracking
clip.EnableHistory(100, 24*time.Hour)

// Copy content (automatically tracked)
clip.Write("Hello")
clip.WriteHTML("<p>Hello</p>")
clip.WriteRTF("{\\rtf1 Hello}")

// View history
entries := clip.GetHistory()
for _, entry := range entries {
    fmt.Printf("ID: %s, Type: %s, Size: %d\n",
        entry.ID, entry.MIMEType, entry.Size)
}

// Restore from history
clip.RestoreFromHistory(entry.ID)

// Memory management
size := clip.GetHistorySize()
totalBytes := clip.GetHistoryTotalSize()
removed := clip.RemoveExpiredHistory()

// Clear history
clip.ClearHistory()

// Disable tracking
clip.DisableHistory()
```

## Implementation Details

### History Storage

- Entries stored in memory (not persisted)
- FIFO eviction when max size exceeded
- Automatic expiration based on age
- Thread-safe concurrent access

### Content Handling

- Plain text: stored as UTF-8 bytes
- HTML: stored with `text/html` MIME type
- RTF: stored with `text/rtf` MIME type
- Images: stored with appropriate `image/*` MIME type

### Performance

- O(1) add operation
- O(n) search/filter operations
- O(1) size/total size queries
- Minimal memory overhead per entry

## Testing

Run tests for the clipboard history feature:

```bash
cd clipboard/internal/domain/model
go test -v -run TestHistoryEntry

cd ../service
go test -v -run TestClipboardHistory
```

## License

Part of the Phoenix TUI Framework.
