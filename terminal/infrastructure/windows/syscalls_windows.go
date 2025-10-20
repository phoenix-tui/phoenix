//go:build windows.
// +build windows.

package windows

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	// Console cursor functions.
	procGetConsoleCursorInfo = kernel32.NewProc("GetConsoleCursorInfo")
	procSetConsoleCursorInfo = kernel32.NewProc("SetConsoleCursorInfo")

	// Console output functions.
	procFillConsoleOutputCharacter = kernel32.NewProc("FillConsoleOutputCharacterW")
	procFillConsoleOutputAttribute = kernel32.NewProc("FillConsoleOutputAttribute")
	procReadConsoleOutput          = kernel32.NewProc("ReadConsoleOutputW")
	procWriteConsoleOutput         = kernel32.NewProc("WriteConsoleOutputW")
	// procScrollConsoleScreenBuffer removed (unused - reserved for future scrolling feature)
)

// ConsoleCursorInfo represents cursor information.
type ConsoleCursorInfo struct {
	Size    uint32
	Visible int32
}

// GetConsoleCursorInfo retrieves cursor information.
func GetConsoleCursorInfo(handle windows.Handle, info *ConsoleCursorInfo) error {
	r1, _, err := procGetConsoleCursorInfo.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(info)),
	)
	if r1 == 0 {
		return err
	}
	return nil
}

// SetConsoleCursorInfo sets cursor information.
func SetConsoleCursorInfo(handle windows.Handle, info *ConsoleCursorInfo) error {
	r1, _, err := procSetConsoleCursorInfo.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(info)),
	)
	if r1 == 0 {
		return err
	}
	return nil
}

// FillConsoleOutputCharacter fills console buffer with character.
func FillConsoleOutputCharacter(
	handle windows.Handle,
	char rune,
	length uint32,
	coord windows.Coord,
	written *uint32,
) error {
	r1, _, err := procFillConsoleOutputCharacter.Call(
		uintptr(handle),
		uintptr(char),
		uintptr(length),
		uintptr(*(*uint32)(unsafe.Pointer(&coord))), // COORD is 4 bytes (2x int16)
		uintptr(unsafe.Pointer(written)),
	)
	if r1 == 0 {
		return err
	}
	return nil
}

// FillConsoleOutputAttribute fills console buffer with attribute.
func FillConsoleOutputAttribute(
	handle windows.Handle,
	attr uint16,
	length uint32,
	coord windows.Coord,
	written *uint32,
) error {
	r1, _, err := procFillConsoleOutputAttribute.Call(
		uintptr(handle),
		uintptr(attr),
		uintptr(length),
		uintptr(*(*uint32)(unsafe.Pointer(&coord))),
		uintptr(unsafe.Pointer(written)),
	)
	if r1 == 0 {
		return err
	}
	return nil
}

// CharInfo represents character and attribute.
type CharInfo struct {
	Char       uint16
	Attributes uint16
}

// SmallRect represents screen rectangle.
type SmallRect struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
}

// ReadConsoleOutput reads from console screen buffer.
func ReadConsoleOutput(
	handle windows.Handle,
	buffer []CharInfo,
	bufferSize windows.Coord,
	bufferCoord windows.Coord,
	readRegion *SmallRect,
) error {
	r1, _, err := procReadConsoleOutput.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(*(*uint32)(unsafe.Pointer(&bufferSize))),
		uintptr(*(*uint32)(unsafe.Pointer(&bufferCoord))),
		uintptr(unsafe.Pointer(readRegion)),
	)
	if r1 == 0 {
		return err
	}
	return nil
}

// WriteConsoleOutput writes to console screen buffer.
func WriteConsoleOutput(
	handle windows.Handle,
	buffer []CharInfo,
	bufferSize windows.Coord,
	bufferCoord windows.Coord,
	writeRegion *SmallRect,
) error {
	r1, _, err := procWriteConsoleOutput.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(*(*uint32)(unsafe.Pointer(&bufferSize))),
		uintptr(*(*uint32)(unsafe.Pointer(&bufferCoord))),
		uintptr(unsafe.Pointer(writeRegion)),
	)
	if r1 == 0 {
		return err
	}
	return nil
}

// Coord is an alias for windows.Coord for consistency.
type Coord = windows.Coord
