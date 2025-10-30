// Package service contains clipboard domain services and business logic.
package service

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/model"
	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/value"
)

// ClipboardHistory is a domain service that manages clipboard history.
// It maintains a FIFO list of clipboard entries with configurable size and age limits.
type ClipboardHistory struct {
	entries []*model.HistoryEntry
	maxSize int
	maxAge  time.Duration
	mu      sync.RWMutex
}

// NewClipboardHistory creates a new clipboard history service.
// maxSize sets the maximum number of entries (0 = unlimited).
// maxAge sets the maximum age of entries (0 = no expiration).
func NewClipboardHistory(maxSize int, maxAge time.Duration) *ClipboardHistory {
	if maxSize < 0 {
		maxSize = 0 // Treat negative as unlimited
	}
	if maxAge < 0 {
		maxAge = 0 // Treat negative as no expiration
	}

	return &ClipboardHistory{
		entries: make([]*model.HistoryEntry, 0, maxSize),
		maxSize: maxSize,
		maxAge:  maxAge,
	}
}

// Add adds a new entry to the history.
// If maxSize is exceeded, the oldest entry is removed (FIFO).
// Returns the created entry or an error.
func (h *ClipboardHistory) Add(content []byte, mimeType value.MIMEType) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Create new entry
	entry, err := model.NewHistoryEntry(content, mimeType)
	if err != nil {
		return fmt.Errorf("failed to create history entry: %w", err)
	}

	// Add to history
	h.entries = append(h.entries, entry)

	// Enforce max size (FIFO eviction)
	if h.maxSize > 0 && len(h.entries) > h.maxSize {
		// Remove oldest entry
		h.entries = h.entries[1:]
	}

	return nil
}

// Get returns the entry with the given ID.
// Returns an error if the entry is not found.
func (h *ClipboardHistory) Get(id string) (*model.HistoryEntry, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, entry := range h.entries {
		if entry.ID() == id {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("history entry not found: %s", id)
}

// GetAll returns all entries sorted by timestamp (newest first).
func (h *ClipboardHistory) GetAll() []*model.HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Create a copy to avoid external modification
	result := make([]*model.HistoryEntry, len(h.entries))
	copy(result, h.entries)

	// Sort by timestamp (newest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp().After(result[j].Timestamp())
	})

	return result
}

// GetRecent returns the N most recent entries.
// If count is 0 or negative, returns all entries.
// If count exceeds available entries, returns all entries.
func (h *ClipboardHistory) GetRecent(count int) []*model.HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if count <= 0 {
		count = len(h.entries)
	}

	// Get all entries sorted by timestamp
	result := make([]*model.HistoryEntry, len(h.entries))
	copy(result, h.entries)

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp().After(result[j].Timestamp())
	})

	// Limit to requested count
	if count < len(result) {
		result = result[:count]
	}

	return result
}

// GetByMIMEType returns all entries matching the given MIME type.
// Results are sorted by timestamp (newest first).
func (h *ClipboardHistory) GetByMIMEType(mimeType value.MIMEType) []*model.HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var result []*model.HistoryEntry

	for _, entry := range h.entries {
		if entry.MIMEType() == mimeType {
			result = append(result, entry)
		}
	}

	// Sort by timestamp (newest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp().After(result[j].Timestamp())
	})

	return result
}

// Clear removes all entries from the history.
func (h *ClipboardHistory) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.entries = make([]*model.HistoryEntry, 0, h.maxSize)
}

// RemoveExpired removes all entries older than maxAge.
// If maxAge is 0 or negative, no entries are removed.
// Returns the number of entries removed.
func (h *ClipboardHistory) RemoveExpired() int {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.maxAge <= 0 {
		return 0 // No expiration
	}

	var kept []*model.HistoryEntry
	removed := 0

	for _, entry := range h.entries {
		if !entry.IsExpired(h.maxAge) {
			kept = append(kept, entry)
		} else {
			removed++
		}
	}

	h.entries = kept
	return removed
}

// Size returns the number of entries in the history.
func (h *ClipboardHistory) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.entries)
}

// TotalSize returns the total memory usage of all entries in bytes.
func (h *ClipboardHistory) TotalSize() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	total := 0
	for _, entry := range h.entries {
		total += entry.Size()
	}
	return total
}

// MaxSize returns the maximum number of entries allowed.
// Returns 0 if unlimited.
func (h *ClipboardHistory) MaxSize() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.maxSize
}

// MaxAge returns the maximum age of entries.
// Returns 0 if no expiration.
func (h *ClipboardHistory) MaxAge() time.Duration {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.maxAge
}

// SetMaxSize updates the maximum number of entries.
// If the new size is smaller than the current size, oldest entries are removed.
func (h *ClipboardHistory) SetMaxSize(maxSize int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if maxSize < 0 {
		maxSize = 0 // Treat negative as unlimited
	}

	h.maxSize = maxSize

	// Enforce new size limit
	if h.maxSize > 0 && len(h.entries) > h.maxSize {
		// Remove oldest entries
		h.entries = h.entries[len(h.entries)-h.maxSize:]
	}
}

// SetMaxAge updates the maximum age of entries.
// Does not automatically remove expired entries - call RemoveExpired() to clean up.
func (h *ClipboardHistory) SetMaxAge(maxAge time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if maxAge < 0 {
		maxAge = 0 // Treat negative as no expiration
	}

	h.maxAge = maxAge
}

// IsEmpty returns true if there are no entries in the history.
func (h *ClipboardHistory) IsEmpty() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.entries) == 0
}

// Contains checks if an entry with the given ID exists.
func (h *ClipboardHistory) Contains(id string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, entry := range h.entries {
		if entry.ID() == id {
			return true
		}
	}
	return false
}
