// Package model contains clipboard domain models and business entities.
package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	value2 "github.com/phoenix-tui/phoenix/clipboard/internal/domain/value"
)

// HistoryEntry represents a single clipboard history item.
// This is an entity in DDD terms with identity (UUID).
type HistoryEntry struct {
	id        string
	content   []byte
	mimeType  value2.MIMEType
	timestamp time.Time
	size      int
}

// NewHistoryEntry creates a new history entry with the given content and MIME type.
// The entry is assigned a unique ID and timestamped with the current time.
func NewHistoryEntry(content []byte, mimeType value2.MIMEType) (*HistoryEntry, error) {
	if len(content) == 0 {
		return nil, fmt.Errorf("history entry content cannot be empty")
	}

	if mimeType == "" {
		return nil, fmt.Errorf("history entry MIME type cannot be empty")
	}

	// Generate unique ID
	id := uuid.New().String()

	// Make a copy of content to preserve immutability
	contentCopy := make([]byte, len(content))
	copy(contentCopy, content)

	return &HistoryEntry{
		id:        id,
		content:   contentCopy,
		mimeType:  mimeType,
		timestamp: time.Now(),
		size:      len(content),
	}, nil
}

// NewHistoryEntryWithTime creates a new history entry with a specific timestamp.
// This is useful for testing and restoring history.
func NewHistoryEntryWithTime(content []byte, mimeType value2.MIMEType, timestamp time.Time) (*HistoryEntry, error) {
	if len(content) == 0 {
		return nil, fmt.Errorf("history entry content cannot be empty")
	}

	if mimeType == "" {
		return nil, fmt.Errorf("history entry MIME type cannot be empty")
	}

	if timestamp.IsZero() {
		return nil, fmt.Errorf("history entry timestamp cannot be zero")
	}

	// Generate unique ID
	id := uuid.New().String()

	// Make a copy of content to preserve immutability
	contentCopy := make([]byte, len(content))
	copy(contentCopy, content)

	return &HistoryEntry{
		id:        id,
		content:   contentCopy,
		mimeType:  mimeType,
		timestamp: timestamp,
		size:      len(content),
	}, nil
}

// ID returns the unique identifier of this history entry.
func (h *HistoryEntry) ID() string {
	return h.id
}

// Content returns a copy of the content to preserve immutability.
func (h *HistoryEntry) Content() []byte {
	result := make([]byte, len(h.content))
	copy(result, h.content)
	return result
}

// MIMEType returns the MIME type of the content.
func (h *HistoryEntry) MIMEType() value2.MIMEType {
	return h.mimeType
}

// Timestamp returns the time when this entry was created.
func (h *HistoryEntry) Timestamp() time.Time {
	return h.timestamp
}

// Size returns the size of the content in bytes.
func (h *HistoryEntry) Size() int {
	return h.size
}

// IsExpired returns true if the entry is older than the given max age.
func (h *HistoryEntry) IsExpired(maxAge time.Duration) bool {
	if maxAge <= 0 {
		return false // No expiration
	}
	return time.Since(h.timestamp) > maxAge
}

// IsText returns true if the content is text-based.
func (h *HistoryEntry) IsText() bool {
	return h.mimeType.IsText()
}

// IsImage returns true if the content is an image.
func (h *HistoryEntry) IsImage() bool {
	return h.mimeType.IsImage()
}

// IsBinary returns true if the content is binary.
func (h *HistoryEntry) IsBinary() bool {
	return h.mimeType.IsBinary()
}

// Text returns the content as text if it's text-based.
func (h *HistoryEntry) Text() (string, error) {
	if !h.IsText() {
		return "", fmt.Errorf("history entry is not text (MIME type: %s)", h.mimeType)
	}
	return string(h.content), nil
}

// Age returns how old this entry is.
func (h *HistoryEntry) Age() time.Duration {
	return time.Since(h.timestamp)
}

// Equals compares two history entries by ID.
func (h *HistoryEntry) Equals(other *HistoryEntry) bool {
	if other == nil {
		return false
	}
	return h.id == other.id
}
