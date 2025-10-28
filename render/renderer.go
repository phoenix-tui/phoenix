// Package render provides high-performance differential rendering for Phoenix TUI.
//
// The render package implements a sophisticated rendering pipeline:.
//   - Virtual buffer representation
//   - Differential rendering (only changed cells)
//   - ANSI sequence optimization
//   - Zero-allocation hot paths
//   - 60 FPS target performance
//
// Example usage:.
//
//	import "github.com/phoenix-tui/phoenix/render"
//
//	// Create renderer
//	renderer := render.New(80, 24, os.Stdout)
//	defer renderer.Close()
//
//	// Create buffer
//	buf := renderer.Buffer()
//	buf.SetString(0, 0, "Hello, World!", render.StyleDefault())
//
//	// Render (differential)
//	if err := renderer.Render(buf); err != nil {
//		log.Fatal(err)
//	}
package render

import (
	"io"
	"os"

	"github.com/phoenix-tui/phoenix/render/internal/application"
	model2 "github.com/phoenix-tui/phoenix/render/internal/domain/model"
	value2 "github.com/phoenix-tui/phoenix/render/internal/domain/value"
)

// Renderer is the public API for high-performance rendering.
type Renderer struct {
	app *application.Renderer
}

// New creates a new renderer with specified dimensions.
// Output is typically os.Stdout for terminal rendering.
func New(width, height int, output io.Writer) *Renderer {
	return &Renderer{
		app: application.NewRenderer(width, height, output),
	}
}

// NewDefault creates a renderer with default terminal dimensions (80x24).
func NewDefault(output io.Writer) *Renderer {
	return New(80, 24, output)
}

// NewStdout creates a renderer writing to stdout.
func NewStdout(width, height int) *Renderer {
	return New(width, height, os.Stdout)
}

// WithBufferSize sets a custom output buffer size (for optimization).
// Default is 4096 bytes.
func WithBufferSize(width, height int, output io.Writer, bufferSize int) *Renderer {
	return &Renderer{
		app: application.NewRendererWithOptions(width, height, output, bufferSize),
	}
}

// Render renders a buffer to the terminal using differential rendering.
// Only changed cells are rendered for optimal performance.
func (r *Renderer) Render(buf *Buffer) error {
	return r.app.Render(buf.internal)
}

// Clear clears the screen.
func (r *Renderer) Clear() error {
	return r.app.Clear()
}

// Resize resizes the renderer dimensions.
func (r *Renderer) Resize(width, height int) error {
	return r.app.Resize(width, height)
}

// HideCursor hides the terminal cursor.
func (r *Renderer) HideCursor() error {
	return r.app.HideCursor()
}

// ShowCursor shows the terminal cursor.
func (r *Renderer) ShowCursor() error {
	return r.app.ShowCursor()
}

// Width returns current width.
func (r *Renderer) Width() int {
	return r.app.Width()
}

// Height returns current height.
func (r *Renderer) Height() int {
	return r.app.Height()
}

// Buffer creates a new buffer for rendering.
// Buffers are pooled internally for zero-allocation rendering.
func (r *Renderer) Buffer() *Buffer {
	return &Buffer{
		internal: r.app.GetBuffer(),
		renderer: r,
	}
}

// Close flushes and closes the renderer.
// Always call Close() when done to restore terminal state.
func (r *Renderer) Close() error {
	return r.app.Close()
}

// Buffer represents a virtual terminal buffer.
// All modifications are made to the buffer, then rendered with Renderer.Render().
type Buffer struct {
	internal *model2.Buffer
	renderer *Renderer
}

// NewBuffer creates a standalone buffer (not pooled).
func NewBuffer(width, height int) *Buffer {
	return &Buffer{
		internal: model2.NewBuffer(width, height),
	}
}

// Width returns buffer width.
func (b *Buffer) Width() int {
	return b.internal.Width()
}

// Height returns buffer height.
func (b *Buffer) Height() int {
	return b.internal.Height()
}

// Set sets a cell at position.
func (b *Buffer) Set(x, y int, char rune, style Style) {
	pos := value2.NewPosition(x, y)
	cell := model2.NewCell(char, style.internal)
	b.internal.Set(pos, cell)
}

// SetString writes a string at position with style.
// Handles Unicode grapheme clusters correctly.
// Returns the number of cells written.
func (b *Buffer) SetString(x, y int, text string, style Style) int {
	pos := value2.NewPosition(x, y)
	return b.internal.SetString(pos, text, style.internal)
}

// SetLine writes a string at the beginning of line y with style.
// Clears the rest of the line.
func (b *Buffer) SetLine(y int, text string, style Style) {
	b.internal.SetLine(y, text, style.internal)
}

// Clear clears all cells to empty.
func (b *Buffer) Clear() {
	b.internal.Clear()
}

// Fill fills all cells with character and style.
func (b *Buffer) Fill(char rune, style Style) {
	b.internal.Fill(char, style.internal)
}

// Release returns the buffer to the pool (if pooled).
// Only call this if buffer was created by Renderer.Buffer().
func (b *Buffer) Release() {
	if b.renderer != nil {
		b.renderer.app.PutBuffer(b.internal)
	}
}

// Style represents ANSI text styling.
type Style struct {
	internal value2.Style
}

// StyleDefault returns an empty style.
func StyleDefault() Style {
	return Style{internal: value2.NewStyle()}
}

// StyleFg returns a style with foreground color.
func StyleFg(r, g, b uint8) Style {
	color := value2.NewColor(r, g, b)
	return Style{internal: value2.NewStyleWithFg(color)}
}

// StyleBg returns a style with background color.
func StyleBg(r, g, b uint8) Style {
	color := value2.NewColor(r, g, b)
	return Style{internal: value2.NewStyleWithBg(color)}
}

// StyleColors returns a style with foreground and background colors.
func StyleColors(fgR, fgG, fgB, bgR, bgG, bgB uint8) Style {
	fg := value2.NewColor(fgR, fgG, fgB)
	bg := value2.NewColor(bgR, bgG, bgB)
	return Style{internal: value2.NewStyleWithColors(fg, bg)}
}

// WithFg returns a new style with foreground color.
func (s Style) WithFg(r, g, b uint8) Style {
	color := value2.NewColor(r, g, b)
	return Style{internal: s.internal.WithFg(color)}
}

// WithBg returns a new style with background color.
func (s Style) WithBg(r, g, b uint8) Style {
	color := value2.NewColor(r, g, b)
	return Style{internal: s.internal.WithBg(color)}
}

// WithBold returns a new style with bold setting.
func (s Style) WithBold(bold bool) Style {
	return Style{internal: s.internal.WithBold(bold)}
}

// WithItalic returns a new style with italic setting.
func (s Style) WithItalic(italic bool) Style {
	return Style{internal: s.internal.WithItalic(italic)}
}

// WithUnderline returns a new style with underline setting.
func (s Style) WithUnderline(underline bool) Style {
	return Style{internal: s.internal.WithUnderline(underline)}
}

// WithReverse returns a new style with reverse video setting.
func (s Style) WithReverse(reverse bool) Style {
	return Style{internal: s.internal.WithReverse(reverse)}
}

// Predefined colors for convenience.
var (
	ColorBlack   = value2.ColorBlack
	ColorRed     = value2.ColorRed
	ColorGreen   = value2.ColorGreen
	ColorYellow  = value2.ColorYellow
	ColorBlue    = value2.ColorBlue
	ColorMagenta = value2.ColorMagenta
	ColorCyan    = value2.ColorCyan
	ColorWhite   = value2.ColorWhite
)

// Predefined styles for convenience.
var (
	// Text styles.
	Bold      = StyleDefault().WithBold(true)
	Italic    = StyleDefault().WithItalic(true)
	Underline = StyleDefault().WithUnderline(true)
	Reverse   = StyleDefault().WithReverse(true)

	// Common foreground colors.
	FgRed     = StyleFg(255, 0, 0)
	FgGreen   = StyleFg(0, 255, 0)
	FgBlue    = StyleFg(0, 0, 255)
	FgYellow  = StyleFg(255, 255, 0)
	FgMagenta = StyleFg(255, 0, 255)
	FgCyan    = StyleFg(0, 255, 255)
	FgWhite   = StyleFg(255, 255, 255)

	// Common background colors.
	BgRed     = StyleBg(255, 0, 0)
	BgGreen   = StyleBg(0, 255, 0)
	BgBlue    = StyleBg(0, 0, 255)
	BgYellow  = StyleBg(255, 255, 0)
	BgMagenta = StyleBg(255, 0, 255)
	BgCyan    = StyleBg(0, 255, 255)
	BgWhite   = StyleBg(255, 255, 255)
)
