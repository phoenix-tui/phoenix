# Rich Text Clipboard Example

This example demonstrates the HTML and RTF clipboard functionality in Phoenix TUI's clipboard module.

## Features

- **Write HTML to clipboard** with styling (bold, italic, underline, colors)
- **Write RTF to clipboard** with styling
- **Read HTML** from clipboard and display both formatted and plain text
- **Read RTF** from clipboard and display both formatted and plain text
- **Strip formatting** - read rich text and convert to plain text
- **Convert formats** - convert between HTML and RTF

## Running the Example

```bash
cd examples/richtext-clipboard
go run main.go
```

## Usage

### Main Menu

The application presents a menu with the following options:

1. **Write HTML to clipboard** - Enter text and apply formatting, then copy as HTML
2. **Write RTF to clipboard** - Enter text and apply formatting, then copy as RTF
3. **Read HTML from clipboard** - Display HTML content from clipboard
4. **Read RTF from clipboard** - Display RTF content from clipboard
5. **Read as plain text** - Read and strip all HTML tags
6. **Convert HTML to RTF** - Paste HTML and convert to RTF
7. **Convert RTF to HTML** - Paste RTF and convert to HTML
8. **Quit** - Exit the application

### Writing Rich Text

When in write mode (option 1 or 2):

- **Type** your text normally
- **Ctrl+B** - Toggle bold
- **Ctrl+I** - Toggle italic
- **Ctrl+U** - Toggle underline
- **Ctrl+R** - Set red color (#FF0000)
- **Ctrl+G** - Set green color (#00FF00)
- **Ctrl+L** - Clear color
- **Enter** - Write to clipboard with current styles
- **Esc** - Return to menu without writing

### Reading Rich Text

When in read mode (option 3, 4, or 5):

- The application displays the clipboard content
- For HTML/RTF, it shows both the raw markup and decoded text with styles
- Press any key to return to the menu

### Converting Formats

When in convert mode (option 6 or 7):

- Type or paste HTML or RTF content
- Press **Enter** to convert
- The result is displayed below
- Press **Esc** to return to menu

## Architecture

The example demonstrates:

- **Domain Services**: Uses `RichTextCodec` for HTML/RTF encoding/decoding
- **Value Objects**: Uses `TextStyles` to represent text formatting
- **Application Layer**: Uses `ClipboardManager` for clipboard operations
- **Public API**: Uses `Clipboard` type for clean, type-safe interface
- **TEA Pattern**: Implements Elm Architecture with Model-Update-View

## Supported Formatting

### HTML

- `<strong>` or `<b>` - Bold text
- `<em>` or `<i>` - Italic text
- `<u>` - Underline text
- `<span style="color:#RRGGBB">` - Colored text (hex format)

### RTF

- `\b` - Bold text
- `\i` - Italic text
- `\ul` - Underline text
- `\cf1` - Colored text (with color table)

## Security

The example demonstrates safe HTML handling:

- All user input is escaped before being written as HTML
- When reading HTML, dangerous tags are stripped
- Only basic formatting tags are supported (no scripts, images, etc.)

## Notes

- The clipboard backend currently stores HTML/RTF as plain text
- Future versions will support native HTML/RTF clipboard formats on supported platforms
- The codec handles all encoding/decoding transparently
- Round-trip conversions (HTML → RTF → HTML) preserve basic formatting

## Testing the Example

Try these workflows:

1. **Write and read HTML**:
   - Option 1, type "Hello World", press Ctrl+B (bold), Enter
   - Option 3 to read back the HTML and see the formatting

2. **Convert formats**:
   - Option 6, paste `<strong>Bold text</strong>`, Enter
   - See the RTF output with `\b` codes

3. **Strip formatting**:
   - Copy some HTML with tags
   - Option 5 to read as plain text (tags removed)

4. **Round trip**:
   - Write HTML → Read HTML → Convert to RTF → Convert back to HTML
   - Verify formatting is preserved
