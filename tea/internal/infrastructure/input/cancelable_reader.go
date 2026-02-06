// Package input provides keyboard input reading and parsing.
package input

import (
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// CancelableReader wraps an io.Reader with cancellation support.
// Essential for ExecProcess to guarantee stdin release before child process starts.
//
// The zero value is not usable; use NewCancelableReader.
//
// Architecture: Pipe-based relay for platform-agnostic cancellation.
//
// The reader uses an os.Pipe() internally. A relay goroutine copies data from
// the underlying reader (e.g., os.Stdin) to the pipe's write end. The readLoop
// goroutine reads from the pipe's read end and delivers results via channel.
//
// Cancellation is instant and reliable on ALL platforms (including MSYS/mintty):
//   - Cancel() closes the pipe's write end
//   - The readLoop's Read() on the pipe's read end returns io.EOF immediately
//   - The relay goroutine is unblocked via SetReadDeadline (if supported)
//     or WriteConsoleInputW (Windows Console), or exits on next stdin event
//
// Fallback: If os.Pipe() fails, falls back to direct read with
// platform-specific UnblockStdinRead().
type CancelableReader struct {
	r io.Reader

	// Cancellation state
	canceled atomic.Bool
	done     chan struct{}
	doneOnce sync.Once

	// Read result channel - buffered to allow background reader to continue
	readCh chan readResult

	// Cleanup synchronization
	readerDone chan struct{}

	// Pipe-based relay (platform-agnostic cancellation)
	pipeReader     *os.File     // Read end - readLoop reads from this
	pipeWriter     *os.File     // Write end - relay writes to this, Cancel() closes
	pipeWriterOnce sync.Once    // Protects pipeWriter from double-close
	relayDone      chan struct{} // Signals when relay goroutine exits
	usePipe        bool         // true if pipe relay is active
}

// readResult holds the result of a single read operation.
type readResult struct {
	data []byte
	n    int
	err  error
}

// NewCancelableReader creates a reader that can be canceled.
// Uses pipe-based relay for instant, platform-agnostic cancellation.
//
// Example:
//
//	cr := NewCancelableReader(os.Stdin)
//	defer cr.Cancel() // Ensure cleanup
//
//	buf := make([]byte, 1024)
//	n, err := cr.Read(buf)
//	if err == io.EOF {
//	    // Either actual EOF or canceled
//	}
func NewCancelableReader(r io.Reader) *CancelableReader {
	cr := &CancelableReader{
		r:          r,
		done:       make(chan struct{}),
		readCh:     make(chan readResult, 1), // Buffered for non-blocking send
		readerDone: make(chan struct{}),
	}

	// Try pipe-based relay (works on all platforms including MSYS/mintty)
	pipeReader, pipeWriter, err := os.Pipe()
	if err == nil {
		cr.pipeReader = pipeReader
		cr.pipeWriter = pipeWriter
		cr.relayDone = make(chan struct{})
		cr.usePipe = true

		go cr.relayLoop()
		go cr.readLoopPipe()
	} else {
		// Fallback: direct read + UnblockStdinRead (legacy behavior)
		// This path is only hit when os.Pipe() fails (extremely rare)
		cr.usePipe = false
		go cr.readLoop()
	}

	return cr
}

// closePipeWriter closes the pipe write end exactly once.
// Safe to call from both Cancel() and relayLoop defer concurrently.
func (cr *CancelableReader) closePipeWriter() {
	cr.pipeWriterOnce.Do(func() {
		cr.pipeWriter.Close()
	})
}

// relayLoop copies data from the underlying reader (stdin) to the pipe writer.
// Exits when: the pipe writer is closed (Cancel), the underlying reader returns
// EOF, or the done channel is closed.
//
// Note: After Cancel(), the relay goroutine may briefly remain blocked in
// cr.r.Read() if SetReadDeadline is not supported. This is acceptable because:
//   - It no longer writes to the pipe (pipe writer is closed)
//   - It will exit on the next stdin event
//   - It does not interfere with the child process reading stdin
//     (the child reads os.Stdin directly, relay reads cr.r which may be different)
//
// When cr.r IS os.Stdin, Cancel() uses SetReadDeadline and UnblockStdinRead
// to ensure the relay exits promptly.
func (cr *CancelableReader) relayLoop() {
	defer close(cr.relayDone)
	defer cr.closePipeWriter() // Signal EOF to readLoopPipe on any exit

	buf := make([]byte, 4096)

	for {
		// Check cancellation before blocking read
		select {
		case <-cr.done:
			return
		default:
		}

		// Read from underlying reader (may block on real stdin)
		n, err := cr.r.Read(buf)

		if n > 0 {
			// Write to pipe (may fail if pipe closed by Cancel)
			_, writeErr := cr.pipeWriter.Write(buf[:n])
			if writeErr != nil {
				// Pipe closed by Cancel() — expected during cancellation
				return
			}
		}

		if err != nil {
			// EOF, DeadlineExceeded, or other error
			return
		}
	}
}

// readLoopPipe reads from the pipe reader (NOT directly from stdin).
// Cancellation is instant: Cancel() closes the pipe writer, causing
// pipeReader.Read() to return io.EOF immediately.
func (cr *CancelableReader) readLoopPipe() {
	defer close(cr.readerDone)

	buf := make([]byte, 256)

	for {
		// Check cancellation before read
		select {
		case <-cr.done:
			return
		default:
		}

		// Read from pipe (instantly cancellable via pipe close)
		n, err := cr.pipeReader.Read(buf)

		// Prepare result (copy data to avoid race)
		result := readResult{
			n:   n,
			err: err,
		}
		if n > 0 {
			result.data = make([]byte, n)
			copy(result.data, buf[:n])
		}

		// Send result or exit if canceled
		select {
		case cr.readCh <- result:
			// Result delivered
		case <-cr.done:
			// Canceled while trying to send
			return
		}

		// Exit on error (EOF from pipe close, etc.)
		if err != nil {
			return
		}
	}
}

// readLoop is the legacy fallback for environments where os.Pipe() fails.
// Reads directly from the underlying reader. Cancellation relies on
// UnblockStdinRead() which may not work on all platforms (e.g., MSYS/mintty).
func (cr *CancelableReader) readLoop() {
	defer close(cr.readerDone)

	buf := make([]byte, 256)

	for {
		select {
		case <-cr.done:
			return
		default:
		}

		n, err := cr.r.Read(buf)

		result := readResult{
			n:   n,
			err: err,
		}
		if n > 0 {
			result.data = make([]byte, n)
			copy(result.data, buf[:n])
		}

		select {
		case cr.readCh <- result:
		case <-cr.done:
			return
		}

		if err != nil {
			return
		}
	}
}

// Read reads data with cancellation support.
// Returns io.EOF immediately if Cancel() has been called.
//
// This method blocks until:
//   - Data is available from underlying reader
//   - An error occurs (including EOF)
//   - Cancel() is called (returns io.EOF)
func (cr *CancelableReader) Read(p []byte) (int, error) {
	// Fast path: already canceled
	if cr.canceled.Load() {
		return 0, io.EOF
	}

	select {
	case result := <-cr.readCh:
		if result.n > 0 {
			n := copy(p, result.data)
			return n, result.err
		}
		return 0, result.err

	case <-cr.done:
		// Canceled while waiting
		return 0, io.EOF
	}
}

// Cancel stops the reader immediately.
// Any blocked Read() calls will return io.EOF.
// Safe to call multiple times.
//
// IMPORTANT: Must be called before ExecProcess to release stdin.
//
// Cancellation strategy (in priority order):
//  1. Close pipe writer → readLoopPipe returns EOF immediately
//  2. SetReadDeadline(now) → unblocks relay goroutine's stdin Read
//  3. WriteConsoleInputW → fallback for true Windows Console
//  4. Timeout → relay exits on next stdin event (safety net)
func (cr *CancelableReader) Cancel() {
	cr.doneOnce.Do(func() {
		cr.canceled.Store(true)
		close(cr.done)

		if cr.usePipe {
			// Phase 1: Close pipe writer → readLoopPipe returns EOF immediately
			cr.closePipeWriter()

			// Phase 2: Try to unblock relay goroutine's stdin Read
			cr.unblockRelay()
		} else {
			// Legacy fallback: inject fake input for direct-read mode
			_ = UnblockStdinRead()
		}
	})
}

// unblockRelay attempts to unblock the relay goroutine which may be
// blocked in cr.r.Read() (reading from real stdin).
//
// Uses two techniques:
//  1. SetReadDeadline — works on os.File with poller support (Go 1.19+)
//  2. UnblockStdinRead — injects fake input on Windows Console
//
// If neither works (e.g., MSYS/mintty without deadline support), the relay
// goroutine will exit on the next stdin event. This is acceptable because
// the readLoopPipe has already exited (pipe closed), so no data is lost.
func (cr *CancelableReader) unblockRelay() {
	// Try SetReadDeadline to immediately unblock the relay's stdin Read
	if f, ok := cr.r.(*os.File); ok {
		// Set deadline to past time — pending Read() returns os.ErrDeadlineExceeded
		_ = f.SetReadDeadline(time.Now())
		// NOTE: This may silently fail on some platforms (returns error).
		// If it fails, relay will exit on next stdin input — acceptable.
	}

	// Also try platform-specific unblock as secondary fallback
	// On Windows Console: injects fake key event via WriteConsoleInputW
	// On Unix/MSYS: no-op (SetReadDeadline above is the primary mechanism)
	_ = UnblockStdinRead()
}

// IsCanceled returns true if the reader has been canceled.
func (cr *CancelableReader) IsCanceled() bool {
	return cr.canceled.Load()
}

// WaitForShutdown waits for background goroutines to exit.
// Call after Cancel() to ensure clean shutdown.
// Returns immediately if reader was never started or already stopped.
func (cr *CancelableReader) WaitForShutdown() {
	// Wait for readLoop/readLoopPipe to exit
	select {
	case <-cr.readerDone:
		// Already done
	case <-cr.done:
		// Wait for reader to notice cancellation
		<-cr.readerDone
	}

	// Wait for relay goroutine (pipe mode only)
	if cr.usePipe && cr.relayDone != nil {
		select {
		case <-cr.relayDone:
			// Relay exited cleanly
		case <-time.After(50 * time.Millisecond):
			// Relay may still be blocked in stdin Read — acceptable.
			// It will exit on next input event or when process exits.
		}
	}

	// Cleanup pipe reader
	if cr.pipeReader != nil {
		_ = cr.pipeReader.Close()
	}

	// Reset stdin deadline for next CancelableReader instance
	if cr.usePipe {
		if f, ok := cr.r.(*os.File); ok {
			_ = f.SetReadDeadline(time.Time{}) // Clear deadline
		}
	}
}
