// Package application orchestrates the rendering pipeline for terminal output.
package application

import (
	"io"
	"sync"

	"github.com/phoenix-tui/phoenix/render/domain/model"
	"github.com/phoenix-tui/phoenix/render/domain/service"
	"github.com/phoenix-tui/phoenix/render/domain/value"
	"github.com/phoenix-tui/phoenix/render/infrastructure/ansi"
)

// Renderer orchestrates the rendering pipeline.
// This is the main application service that coordinates:.
// - Diffing (finding changes)
// - Optimization (reducing ANSI sequences)
// - Writing (outputting to terminal).
type Renderer struct {
	diffService     *service.DiffService
	optimizeService *service.OptimizeService
	writer          *ansi.Writer

	currentBuffer  *model.Buffer
	previousBuffer *model.Buffer

	width  int
	height int

	bufferPool sync.Pool // For zero-allocation rendering
	mu         sync.Mutex
}

// NewRenderer creates a new renderer with specified dimensions.
func NewRenderer(width, height int, output io.Writer) *Renderer {
	r := &Renderer{
		diffService:     service.NewDiffService(),
		optimizeService: service.NewOptimizeService(),
		writer:          ansi.NewWriter(output),
		currentBuffer:   model.NewBuffer(width, height),
		previousBuffer:  model.NewBuffer(width, height),
		width:           width,
		height:          height,
	}

	// Initialize buffer pool for zero-allocation.
	r.bufferPool.New = func() interface{} {
		return model.NewBuffer(width, height)
	}

	return r
}

// NewRendererWithOptions creates a renderer with custom options.
func NewRendererWithOptions(width, height int, output io.Writer, bufferSize int) *Renderer {
	r := &Renderer{
		diffService:     service.NewDiffService(),
		optimizeService: service.NewOptimizeService(),
		writer:          ansi.NewWriterWithBuffer(output, bufferSize),
		currentBuffer:   model.NewBuffer(width, height),
		previousBuffer:  model.NewBuffer(width, height),
		width:           width,
		height:          height,
	}

	r.bufferPool.New = func() interface{} {
		return model.NewBuffer(width, height)
	}

	return r
}

// Render renders a buffer to the terminal.
// This is the main rendering method - implements differential rendering for performance.
func (r *Renderer) Render(buffer *model.Buffer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if buffer == nil {
		return nil
	}

	// Check if dimensions changed.
	if buffer.Width() != r.width || buffer.Height() != r.height {
		return r.renderFullScreen(buffer)
	}

	// 1. Compute diff
	ops := r.diffService.Diff(r.previousBuffer, buffer)

	// Early exit if no changes.
	if len(ops) == 0 {
		return nil
	}

	// 2. Optimize operations
	ops = r.optimizeService.RemoveRedundant(ops, r.writer.CurrentStyle())

	// 3. Check if full redraw is more efficient
	totalCells := r.width * r.height
	if r.optimizeService.ShouldFullRedraw(len(ops), totalCells) {
		return r.renderFullScreen(buffer)
	}

	// 4. Apply operations
	if err := r.applyOperations(ops); err != nil {
		return err
	}

	// 5. Flush output
	if err := r.writer.Flush(); err != nil {
		return err
	}

	// 6. Update previous buffer
	r.previousBuffer = buffer.Clone()

	return nil
}

// renderFullScreen renders entire buffer (for resize or large changes).
//
//nolint:gocognit // Fullscreen rendering orchestrates multiple operations
func (r *Renderer) renderFullScreen(buffer *model.Buffer) error {
	// Clear screen.
	if err := r.writer.Clear(); err != nil {
		return err
	}

	// Move to home.
	if err := r.writer.MoveCursor(0, 0); err != nil {
		return err
	}

	// Render all non-empty cells.
	for y := 0; y < buffer.Height(); y++ {
		for x := 0; x < buffer.Width(); x++ {
			cell := buffer.Get(value.NewPosition(x, y))
			if !cell.IsEmpty() {
				if err := r.writer.MoveCursor(x, y); err != nil {
					return err
				}
				if err := r.writer.WriteCell(cell); err != nil {
					return err
				}
			}
		}
	}

	// Flush.
	if err := r.writer.Flush(); err != nil {
		return err
	}

	// Update state.
	r.width = buffer.Width()
	r.height = buffer.Height()
	r.previousBuffer = buffer.Clone()

	return nil
}

// applyOperations applies diff operations to terminal.
//
//nolint:gocognit // Operation application requires sequential checks
func (r *Renderer) applyOperations(ops []service.DiffOp) error {
	for _, op := range ops {
		switch op.Type {
		case service.OpTypeSet:
			if err := r.writer.MoveCursor(op.Position.X(), op.Position.Y()); err != nil {
				return err
			}
			if err := r.writer.WriteCell(op.Cell); err != nil {
				return err
			}

		case service.OpTypeClear:
			if err := r.writer.MoveCursor(op.Position.X(), op.Position.Y()); err != nil {
				return err
			}
			if err := r.writer.ClearLine(); err != nil {
				return err
			}

		case service.OpTypeMoveCursor:
			if err := r.writer.MoveCursor(op.Position.X(), op.Position.Y()); err != nil {
				return err
			}
		}
	}

	return nil
}

// Clear clears the screen and resets state.
func (r *Renderer) Clear() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.clearUnlocked()
}

// clearUnlocked clears without locking (internal use).
func (r *Renderer) clearUnlocked() error {
	if err := r.writer.Clear(); err != nil {
		return err
	}

	if err := r.writer.MoveCursor(0, 0); err != nil {
		return err
	}

	if err := r.writer.Flush(); err != nil {
		return err
	}

	// Reset buffers.
	r.currentBuffer.Clear()
	r.previousBuffer.Clear()

	return nil
}

// Resize resizes the renderer and buffers.
func (r *Renderer) Resize(width, height int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.width = width
	r.height = height

	// Recreate buffers.
	r.currentBuffer = model.NewBuffer(width, height)
	r.previousBuffer = model.NewBuffer(width, height)

	// Update buffer pool.
	r.bufferPool.New = func() interface{} {
		return model.NewBuffer(width, height)
	}

	return r.clearUnlocked() // Use unlocked version to avoid deadlock
}

// HideCursor hides the terminal cursor.
func (r *Renderer) HideCursor() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.writer.HideCursor(); err != nil {
		return err
	}

	return r.writer.Flush()
}

// ShowCursor shows the terminal cursor.
func (r *Renderer) ShowCursor() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.writer.ShowCursor(); err != nil {
		return err
	}

	return r.writer.Flush()
}

// Close flushes and closes the renderer.
func (r *Renderer) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Show cursor on exit.
	_ = r.writer.ShowCursor()

	return r.writer.Close()
}

// Width returns current width.
func (r *Renderer) Width() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.width
}

// Height returns current height.
func (r *Renderer) Height() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.height
}

// GetBuffer returns a buffer from the pool (for zero-allocation).
func (r *Renderer) GetBuffer() *model.Buffer {
	buf := r.bufferPool.Get().(*model.Buffer)
	buf.Clear()
	return buf
}

// PutBuffer returns a buffer to the pool.
func (r *Renderer) PutBuffer(buf *model.Buffer) {
	if buf != nil {
		r.bufferPool.Put(buf)
	}
}

// Stats returns rendering statistics.
type Stats struct {
	TotalRenders  int
	CellsChanged  int
	LastRenderOps int
	BufferedBytes int
}

// GetStats returns current rendering statistics.
func (r *Renderer) GetStats() Stats {
	r.mu.Lock()
	defer r.mu.Unlock()

	return Stats{
		BufferedBytes: r.writer.BufferedSize(),
	}
}
