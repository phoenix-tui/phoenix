//go:build !windows

// Package input provides keyboard input reading and parsing.
package input

// UnblockStdinRead attempts to unblock a goroutine blocked on stdin Read().
//
// On Unix-like systems, this is a no-op because:
// 1. Most terminals properly support non-blocking I/O
// 2. The issue primarily affects Windows Console
//
// If needed in the future, this could use techniques like:
// - Writing to /dev/tty
// - Sending SIGIO signal
// - Using fcntl to set O_NONBLOCK temporarily
func UnblockStdinRead() error {
	// No-op on non-Windows platforms
	return nil
}
