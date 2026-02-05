//go:build windows

// Package input provides keyboard input reading and parsing.
package input

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// Windows Console API constants
const (
	keyEvent = 0x0001 // KEY_EVENT in INPUT_RECORD.EventType
)

// keyEventRecord matches Windows KEY_EVENT_RECORD structure
type keyEventRecord struct {
	KeyDown         int32
	RepeatCount     uint16
	VirtualKeyCode  uint16
	VirtualScanCode uint16
	UnicodeChar     uint16
	ControlKeyState uint32
}

// inputRecord matches Windows INPUT_RECORD structure
type inputRecord struct {
	EventType uint16
	_         uint16 // padding
	Event     [16]byte
}

var (
	kernel32              = windows.NewLazySystemDLL("kernel32.dll")
	procWriteConsoleInput = kernel32.NewProc("WriteConsoleInputW")
)

// UnblockStdinRead injects a fake key event into the console input buffer
// to unblock any goroutine blocked on stdin Read().
//
// This is necessary because Go's os.Stdin.Read() is a blocking syscall
// that cannot be interrupted by closing a channel. The only way to unblock
// it is to inject data into the input buffer.
//
// The injected event is a null character (VK 0) key-up event, which:
// - Unblocks the Read() call
// - Is typically ignored by applications
// - Doesn't produce visible output
//
// Returns nil if stdin is not a console (e.g., redirected or in tests).
// This is not an error - just means unblocking is not needed.
//
// NOTE: On MSYS2/mintty (Git Bash), stdin is a pty, not a Windows Console.
// GetConsoleMode() fails → this function becomes a no-op. For MSYS, the
// pipe-based CancelableReader + SetReadDeadline is the primary cancellation
// mechanism. This function serves as a secondary fallback for true Windows
// Console environments.
func UnblockStdinRead() error {
	// Get stdin handle
	handle, err := windows.GetStdHandle(windows.STD_INPUT_HANDLE)
	if err != nil {
		return err
	}

	// Check if stdin is actually a console
	// In tests or when stdin is redirected, this will fail
	var mode uint32
	err = windows.GetConsoleMode(handle, &mode)
	if err != nil {
		// Not a console - return nil (not an error, just not applicable)
		return nil
	}

	// Create a minimal key event (key up, null character)
	// Using key-up to minimize side effects
	var keyEvent keyEventRecord
	keyEvent.KeyDown = 0 // Key up event
	keyEvent.RepeatCount = 1
	keyEvent.VirtualKeyCode = 0 // VK_NULL - produces no character
	keyEvent.VirtualScanCode = 0
	keyEvent.UnicodeChar = 0
	keyEvent.ControlKeyState = 0

	// Create INPUT_RECORD
	var record inputRecord
	record.EventType = 0x0001 // KEY_EVENT
	*(*keyEventRecord)(unsafe.Pointer(&record.Event[0])) = keyEvent

	// Write to console input buffer
	var written uint32
	r1, _, err := procWriteConsoleInput.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&record)),
		1, // number of records
		uintptr(unsafe.Pointer(&written)),
	)

	if r1 == 0 {
		return err
	}

	return nil
}
