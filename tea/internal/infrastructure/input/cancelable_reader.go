// Package input provides keyboard input reading and parsing.
package input

import (
	"io"
	"sync"
	"sync/atomic"
)

// CancelableReader wraps an io.Reader with cancellation support.
// Essential for ExecProcess to guarantee stdin release before child process starts.
//
// The zero value is not usable; use NewCancelableReader.
//
// Problem solved:
// Standard io.Reader.Read() is blocking and cannot be interrupted.
// When ExecProcess needs to run an interactive command (vim, ssh, python),
// it must fully release stdin. Without CancelableReader, the inputReader
// goroutine remains blocked in Read(), causing race conditions.
//
// Solution:
// CancelableReader uses a background goroutine to perform blocking reads,
// then delivers results via channel. Cancel() closes the done channel,
// causing any pending Read() to return io.EOF immediately.
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
}

// readResult holds the result of a single read operation.
type readResult struct {
	data []byte
	n    int
	err  error
}

// NewCancelableReader creates a reader that can be canceled.
// The background reader goroutine starts immediately.
//
// Example:
//
//	cr := NewCancelableReader(os.Stdin)
//	defer cr.Cancel() // Ensure cleanup
//
//	// Read with cancellation support
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

	// Start background reader goroutine
	go cr.readLoop()

	return cr
}

// readLoop continuously reads from underlying reader.
// Results are sent to readCh for the main Read to consume.
// Exits when canceled or on read error.
func (cr *CancelableReader) readLoop() {
	defer close(cr.readerDone)

	// Use small buffer for responsive cancellation
	// Larger reads would block longer before checking done channel
	buf := make([]byte, 256)

	for {
		// Check cancellation before blocking read
		select {
		case <-cr.done:
			return
		default:
		}

		// Read from underlying reader (blocking)
		n, err := cr.r.Read(buf)

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

		// Exit on error (EOF, etc.)
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
			// Return actual bytes copied (respects caller's buffer size)
			// Note: any data beyond len(p) is lost. This is acceptable
			// because CancelableReader is wrapped with bufio.Reader.
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
func (cr *CancelableReader) Cancel() {
	cr.doneOnce.Do(func() {
		cr.canceled.Store(true)
		close(cr.done)
	})
}

// IsCanceled returns true if the reader has been canceled.
func (cr *CancelableReader) IsCanceled() bool {
	return cr.canceled.Load()
}

// WaitForShutdown waits for the background reader goroutine to exit.
// Call after Cancel() to ensure clean shutdown.
// Returns immediately if reader was never started or already stopped.
func (cr *CancelableReader) WaitForShutdown() {
	select {
	case <-cr.readerDone:
		// Already done
	case <-cr.done:
		// Wait for reader to notice cancellation
		<-cr.readerDone
	}
}
