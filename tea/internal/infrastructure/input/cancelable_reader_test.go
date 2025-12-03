package input

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestCancelableReader_BasicRead(t *testing.T) {
	data := "hello world"
	cr := NewCancelableReader(strings.NewReader(data))
	defer cr.Cancel()

	buf := make([]byte, 256)
	n, err := cr.Read(buf)

	if err != nil && err != io.EOF {
		t.Errorf("unexpected error: %v", err)
	}
	if string(buf[:n]) != data {
		t.Errorf("got %q, want %q", string(buf[:n]), data)
	}
}

func TestCancelableReader_MultipleReads(t *testing.T) {
	data := "hello world this is a longer string for testing"
	cr := NewCancelableReader(strings.NewReader(data))
	defer cr.Cancel()

	var result bytes.Buffer
	buf := make([]byte, 256) // Buffer large enough for CancelableReader's internal buffer

	for {
		n, err := cr.Read(buf)
		if n > 0 {
			result.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if result.String() != data {
		t.Errorf("got %q, want %q", result.String(), data)
	}
}

func TestCancelableReader_CancelBeforeRead(t *testing.T) {
	// Use a reader that would block forever
	pr, _ := io.Pipe()
	cr := NewCancelableReader(pr)

	// Cancel before reading
	cr.Cancel()

	buf := make([]byte, 256)
	n, err := cr.Read(buf)

	if n != 0 {
		t.Errorf("expected 0 bytes, got %d", n)
	}
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestCancelableReader_CancelDuringRead(t *testing.T) {
	// Use a reader that blocks
	pr, _ := io.Pipe()
	cr := NewCancelableReader(pr)

	// Start read in goroutine
	done := make(chan struct{})
	var readErr error
	go func() {
		buf := make([]byte, 256)
		_, readErr = cr.Read(buf)
		close(done)
	}()

	// Give read time to start blocking
	time.Sleep(50 * time.Millisecond)

	// Cancel
	cr.Cancel()

	// Wait for read to complete
	select {
	case <-done:
		// Success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Read did not unblock after Cancel")
	}

	if readErr != io.EOF {
		t.Errorf("expected io.EOF after cancel, got %v", readErr)
	}
}

func TestCancelableReader_IsCanceled(t *testing.T) {
	cr := NewCancelableReader(strings.NewReader("test"))

	if cr.IsCanceled() {
		t.Error("should not be canceled initially")
	}

	cr.Cancel()

	if !cr.IsCanceled() {
		t.Error("should be canceled after Cancel()")
	}
}

func TestCancelableReader_CancelMultipleTimes(t *testing.T) {
	cr := NewCancelableReader(strings.NewReader("test"))

	// Should not panic
	cr.Cancel()
	cr.Cancel()
	cr.Cancel()

	if !cr.IsCanceled() {
		t.Error("should be canceled")
	}
}

func TestCancelableReader_WaitForShutdown(t *testing.T) {
	pr, pw := io.Pipe()
	cr := NewCancelableReader(pr)

	// Write some data to ensure reader is active
	go func() {
		pw.Write([]byte("test"))
		pw.Close() // Close to unblock reader
	}()

	// Read the data
	buf := make([]byte, 256)
	cr.Read(buf)

	// Cancel and wait
	cr.Cancel()

	done := make(chan struct{})
	go func() {
		cr.WaitForShutdown()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("WaitForShutdown did not complete")
	}
}

func TestCancelableReader_ConcurrentReads(t *testing.T) {
	// Create a reader with enough data
	data := strings.Repeat("x", 1000)
	cr := NewCancelableReader(strings.NewReader(data))
	defer cr.Cancel()

	var wg sync.WaitGroup
	errors := make(chan error, 10)

	// Start multiple concurrent readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := make([]byte, 100)
			_, err := cr.Read(buf)
			if err != nil && err != io.EOF {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("concurrent read error: %v", err)
	}
}

func TestCancelableReader_ReadAfterEOF(t *testing.T) {
	cr := NewCancelableReader(strings.NewReader(""))
	defer cr.Cancel()

	buf := make([]byte, 256)

	// First read should return EOF
	_, err := cr.Read(buf)
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}
}

func TestCancelableReader_SmallBuffer(t *testing.T) {
	cr := NewCancelableReader(strings.NewReader("test"))
	defer cr.Cancel()

	// Read with small buffer - should still work
	buf := make([]byte, 2)
	n, err := cr.Read(buf)

	// Should return up to buffer size
	if n > len(buf) {
		t.Errorf("read more than buffer size: got %d, max %d", n, len(buf))
	}
	if err != nil && err != io.EOF {
		t.Errorf("unexpected error: %v", err)
	}
}

// BenchmarkCancelableReader_Read measures read performance.
func BenchmarkCancelableReader_Read(b *testing.B) {
	data := strings.Repeat("x", 4096)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cr := NewCancelableReader(strings.NewReader(data))
		buf := make([]byte, 4096)
		cr.Read(buf)
		cr.Cancel()
	}
}

// BenchmarkCancelableReader_Cancel measures cancel latency.
func BenchmarkCancelableReader_Cancel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pr, _ := io.Pipe()
		cr := NewCancelableReader(pr)
		cr.Cancel()
	}
}
