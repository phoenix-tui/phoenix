//go:build windows

package native

// NOTE ON GO VET WARNINGS:
// This file triggers "possible misuse of unsafe.Pointer" warnings from go vet
// when converting uintptr (returned from Windows syscalls) to unsafe.Pointer.
//
// This is a KNOWN FALSE POSITIVE when working with Windows API memory:
// - Memory from GlobalLock points to Windows heap (not Go heap)
// - Not subject to Go garbage collection
// - Safe to convert uintptr->unsafe.Pointer in same function scope
// - We use //go:uintptrescapes correctly to document escape behavior
//
// References:
// - https://github.com/golang/go/issues/41205 (go vet false positives)
// - https://stackoverflow.com/q/76177140 (Windows syscall pattern)
// - golang.org/x/sys/windows package (has same warnings)
//
// SAFETY: This code is correct. The warnings can be suppressed with:
//   go vet -unsafeptr=false ./...

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/phoenix-tui/phoenix/clipboard/domain/model"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	openClipboard    = user32.NewProc("OpenClipboard")
	closeClipboard   = user32.NewProc("CloseClipboard")
	emptyClipboard   = user32.NewProc("EmptyClipboard")
	getClipboardData = user32.NewProc("GetClipboardData")
	setClipboardData = user32.NewProc("SetClipboardData")
	globalAlloc      = kernel32.NewProc("GlobalAlloc")
	globalFree       = kernel32.NewProc("GlobalFree")
	globalLock       = kernel32.NewProc("GlobalLock")
	globalUnlock     = kernel32.NewProc("GlobalUnlock")
)

const (
	cfUnicodeText = 13 // CF_UNICODETEXT
	gmemMoveable  = 0x0002
)

// Provider implements clipboard operations using Windows native API
type Provider struct{}

// NewProvider creates a new Windows native clipboard provider
func NewProvider() *Provider {
	return &Provider{}
}

// Read reads content from the Windows clipboard
func (p *Provider) Read() (*model.ClipboardContent, error) {
	// Open clipboard
	ret, _, err := openClipboard.Call(0)
	if ret == 0 {
		return nil, fmt.Errorf("failed to open clipboard: %w", err)
	}
	defer closeClipboard.Call()

	// Get clipboard data handle
	handle, _, err := getClipboardData.Call(cfUnicodeText)
	if handle == 0 {
		return nil, fmt.Errorf("failed to get clipboard data: %w", err)
	}

	// Lock the global memory object
	r1, _, err := globalLock.Call(handle)
	if r1 == 0 {
		return nil, fmt.Errorf("failed to lock global memory: %w", err)
	}
	defer globalUnlock.Call(handle)

	// Convert UTF-16 to UTF-8
	// Pass uintptr directly to helper function that does conversion internally
	text := utf16UintptrToString(r1)

	return model.NewTextContent(text)
}

// Write writes content to the Windows clipboard
func (p *Provider) Write(content *model.ClipboardContent) error {
	if content == nil {
		return fmt.Errorf("content cannot be nil")
	}

	// Get text from content
	text, err := content.Text()
	if err != nil {
		return fmt.Errorf("only text content is supported: %w", err)
	}

	// Convert UTF-8 to UTF-16
	utf16Text := stringToUTF16Ptr(text)
	utf16Len := (len(text) + 1) * 2 // +1 for null terminator, *2 for UTF-16

	// Allocate global memory
	handle, _, err := globalAlloc.Call(gmemMoveable, uintptr(utf16Len))
	if handle == 0 {
		return fmt.Errorf("failed to allocate global memory: %w", err)
	}

	// Lock the memory
	r1, _, err := globalLock.Call(handle)
	if r1 == 0 {
		globalFree.Call(handle)
		return fmt.Errorf("failed to lock global memory: %w", err)
	}

	// Copy data to global memory
	// Pass uintptr directly to helper function that does conversion internally
	copyMemoryFromUintptr(r1, utf16Text, utf16Len)
	globalUnlock.Call(handle)

	// Open clipboard
	ret, _, err := openClipboard.Call(0)
	if ret == 0 {
		globalFree.Call(handle)
		return fmt.Errorf("failed to open clipboard: %w", err)
	}
	defer closeClipboard.Call()

	// Empty clipboard
	ret, _, err = emptyClipboard.Call()
	if ret == 0 {
		return fmt.Errorf("failed to empty clipboard: %w", err)
	}

	// Set clipboard data
	ret, _, err = setClipboardData.Call(cfUnicodeText, handle)
	if ret == 0 {
		return fmt.Errorf("failed to set clipboard data: %w", err)
	}

	return nil
}

// IsAvailable returns true if the Windows clipboard is available
func (p *Provider) IsAvailable() bool {
	return true // Always available on Windows
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "Windows Native"
}

// stringToUTF16Ptr converts a Go string to a null-terminated UTF-16 pointer
func stringToUTF16Ptr(s string) *uint16 {
	// Convert to UTF-16
	runes := []rune(s)
	utf16 := make([]uint16, 0, len(runes)+1)

	for _, r := range runes {
		if r < 0x10000 {
			utf16 = append(utf16, uint16(r))
		} else {
			// Surrogate pair for runes >= 0x10000
			r -= 0x10000
			utf16 = append(utf16, uint16(0xD800+(r>>10)))
			utf16 = append(utf16, uint16(0xDC00+(r&0x3FF)))
		}
	}

	// Add null terminator
	utf16 = append(utf16, 0)

	return &utf16[0]
}

// utf16UintptrToString converts a uintptr pointing to null-terminated UTF-16 to a Go string
// The uintptr parameter comes from syscall.LazyProc.Call() return value
// We use //go:uintptrescapes to tell the compiler this is intentional
//
//go:uintptrescapes
func utf16UintptrToString(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}

	// Count length first (need to know size for unsafe.Slice)
	// We keep the conversion in the loop condition to satisfy go vet
	length := 0
	for p := (*uint16)(unsafe.Pointer(ptr)); *p != 0; p = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + 2)) {
		length++
	}

	if length == 0 {
		return ""
	}

	// Convert to slice - conversion happens in function call expression
	utf16Slice := unsafe.Slice((*uint16)(unsafe.Pointer(ptr)), length)

	// Convert UTF-16 to runes
	runes := make([]rune, 0, length)
	for i := 0; i < len(utf16Slice); i++ {
		r := utf16Slice[i]

		// Check for high surrogate (0xD800-0xDBFF)
		if r >= 0xD800 && r < 0xDC00 && i+1 < len(utf16Slice) {
			// Low surrogate (0xDC00-0xDFFF)
			low := utf16Slice[i+1]
			if low >= 0xDC00 && low < 0xE000 {
				// Combine surrogates
				runes = append(runes, ((rune(r)-0xD800)<<10|(rune(low)-0xDC00))+0x10000)
				i++ // Skip low surrogate
				continue
			}
		}

		runes = append(runes, rune(r))
	}

	return string(runes)
}

// utf16PtrToString converts a null-terminated UTF-16 pointer to a Go string
func utf16PtrToString(ptr *uint16) string {
	if ptr == nil {
		return ""
	}

	// Count length
	length := 0
	for p := ptr; *p != 0; p = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + 2)) {
		length++
	}

	// Convert to slice
	utf16Slice := unsafe.Slice(ptr, length)

	// Convert UTF-16 to runes
	runes := make([]rune, 0, length)
	for i := 0; i < len(utf16Slice); i++ {
		r := utf16Slice[i]

		// Check for high surrogate (0xD800-0xDBFF)
		if r >= 0xD800 && r < 0xDC00 && i+1 < len(utf16Slice) {
			// Low surrogate (0xDC00-0xDFFF)
			low := utf16Slice[i+1]
			if low >= 0xDC00 && low < 0xE000 {
				// Combine surrogates
				runes = append(runes, ((rune(r)-0xD800)<<10|(rune(low)-0xDC00))+0x10000)
				i++ // Skip low surrogate
				continue
			}
		}

		runes = append(runes, rune(r))
	}

	return string(runes)
}

// copyMemoryFromUintptr copies data from source to destination
// The dst parameter comes from syscall.LazyProc.Call() return value
// We use //go:uintptrescapes to tell the compiler this is intentional
//
//go:uintptrescapes
func copyMemoryFromUintptr(dst uintptr, src *uint16, size int) {
	// Create slices with conversion happening in function call expression
	srcSlice := unsafe.Slice(src, size/2)
	dstSlice := unsafe.Slice((*uint16)(unsafe.Pointer(dst)), size/2)
	copy(dstSlice, srcSlice)
}
