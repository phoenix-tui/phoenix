//go:build !windows

// Package input provides keyboard input reading and parsing.
package input

// UnblockStdinRead attempts to unblock a goroutine blocked on stdin Read().
//
// On Unix-like systems, this is a no-op because the pipe-based CancelableReader
// handles cancellation by closing the os.Pipe writer, which causes an immediate
// EOF on the read end. SetReadDeadline on os.Stdin unblocks the relay goroutine.
//
// This function is kept as a fallback interface for symmetry with the Windows
// implementation (WriteConsoleInputW).
func UnblockStdinRead() error {
	// No-op on non-Windows platforms
	return nil
}
