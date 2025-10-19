package service

import (
	"github.com/rivo/uniseg"
)

// CursorMovementService handles grapheme-aware cursor movement.
// It uses uniseg for proper Unicode grapheme cluster segmentation,
// ensuring correct handling of emoji, combining characters, and CJK text.
type CursorMovementService struct{}

// NewCursorMovementService creates a new cursor movement service.
func NewCursorMovementService() *CursorMovementService {
	return &CursorMovementService{}
}

// MoveLeft moves the cursor left by one grapheme cluster.
// Returns the new cursor position (clamped to 0).
func (s *CursorMovementService) MoveLeft(content string, currentPos int) int {
	if currentPos <= 0 {
		return 0
	}

	// Count grapheme clusters from start until we reach currentPos-1
	targetPos := currentPos - 1
	var byteOffset int
	graphemeCount := 0

	gr := uniseg.NewGraphemes(content)
	for gr.Next() {
		if graphemeCount == targetPos {
			return targetPos
		}
		byteOffset += len(gr.Bytes())
		graphemeCount++
	}

	// If we've processed all graphemes and haven't reached target,
	// return the count of graphemes we found
	if graphemeCount < targetPos {
		return graphemeCount
	}

	return targetPos
}

// MoveRight moves the cursor right by one grapheme cluster.
// Returns the new cursor position (clamped to content length in graphemes).
func (s *CursorMovementService) MoveRight(content string, currentPos int) int {
	// Count total graphemes in content
	maxPos := s.GraphemeCount(content)
	if currentPos >= maxPos {
		return maxPos
	}

	return currentPos + 1
}

// GraphemeCount returns the number of grapheme clusters in the content.
func (s *CursorMovementService) GraphemeCount(content string) int {
	count := 0
	gr := uniseg.NewGraphemes(content)
	for gr.Next() {
		count++
	}
	return count
}

// SplitAtCursor splits content at the cursor position into three parts:
// - before: text before cursor
// - at: grapheme cluster at cursor (empty if at end)
// - after: text after cursor
func (s *CursorMovementService) SplitAtCursor(content string, pos int) (before, at, after string) {
	if pos < 0 {
		pos = 0
	}

	gr := uniseg.NewGraphemes(content)
	var byteOffset int
	graphemeCount := 0

	// Find byte offset for the grapheme cluster at position 'pos'
	for gr.Next() {
		if graphemeCount == pos {
			// Found the position
			before = content[:byteOffset]
			atBytes := gr.Bytes()
			at = string(atBytes)
			after = content[byteOffset+len(atBytes):]
			return
		}
		byteOffset += len(gr.Bytes())
		graphemeCount++
	}

	// Cursor is at or beyond end
	before = content
	at = ""
	after = ""
	return
}

// ByteOffsetToGraphemeOffset converts a byte offset to a grapheme offset.
// Useful for converting positions from string operations to cursor positions.
func (s *CursorMovementService) ByteOffsetToGraphemeOffset(content string, byteOffset int) int {
	if byteOffset <= 0 {
		return 0
	}
	if byteOffset >= len(content) {
		return s.GraphemeCount(content)
	}

	gr := uniseg.NewGraphemes(content)
	currentByteOffset := 0
	graphemeCount := 0

	for gr.Next() {
		nextByteOffset := currentByteOffset + len(gr.Bytes())
		if nextByteOffset > byteOffset {
			// The byte offset falls within this grapheme cluster
			return graphemeCount
		}
		currentByteOffset = nextByteOffset
		graphemeCount++
	}

	return graphemeCount
}

// GraphemeOffsetToByteOffset converts a grapheme offset to a byte offset.
// Useful for converting cursor positions to string indices.
func (s *CursorMovementService) GraphemeOffsetToByteOffset(content string, graphemeOffset int) int {
	if graphemeOffset <= 0 {
		return 0
	}

	gr := uniseg.NewGraphemes(content)
	byteOffset := 0
	graphemeCount := 0

	for gr.Next() {
		if graphemeCount == graphemeOffset {
			return byteOffset
		}
		byteOffset += len(gr.Bytes())
		graphemeCount++
	}

	// Offset is at or beyond end
	return len(content)
}
