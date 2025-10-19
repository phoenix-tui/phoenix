// Package model provides rich domain models for textarea.
package model

// KillRing implements Emacs-style kill ring (clipboard with history).
// This is a rich domain model that encapsulates kill ring behavior.
// All operations return new instances (immutable).
type KillRing struct {
	items   []string
	maxSize int
	index   int // Current yank position
}

// NewKillRing creates kill ring with max size.
func NewKillRing(maxSize int) *KillRing {
	if maxSize <= 0 {
		maxSize = 10 // Default
	}
	return &KillRing{
		items:   make([]string, 0, maxSize),
		maxSize: maxSize,
		index:   0,
	}
}

// Kill adds text to kill ring (returns new ring).
func (k *KillRing) Kill(text string) *KillRing {
	if text == "" {
		return k
	}

	copy := k.Copy()

	// Add to ring (circular buffer)
	if len(copy.items) >= copy.maxSize {
		// Ring is full, remove oldest
		copy.items = copy.items[1:]
	}

	copy.items = append(copy.items, text)
	copy.index = len(copy.items) - 1 // Point to latest

	return copy
}

// Yank returns current kill ring item.
func (k *KillRing) Yank() string {
	if len(k.items) == 0 {
		return ""
	}
	if k.index < 0 || k.index >= len(k.items) {
		return ""
	}
	return k.items[k.index]
}

// YankPop rotates kill ring backward (for Emacs M-y).
// Returns new ring with index decremented.
func (k *KillRing) YankPop() *KillRing {
	copy := k.Copy()

	if len(copy.items) == 0 {
		return copy
	}

	copy.index--
	if copy.index < 0 {
		copy.index = len(copy.items) - 1
	}

	return copy
}

// IsEmpty returns true if kill ring has no items.
func (k *KillRing) IsEmpty() bool {
	return len(k.items) == 0
}

// Copy returns deep copy of kill ring.
func (k *KillRing) Copy() *KillRing {
	itemsCopy := make([]string, len(k.items))
	copy(itemsCopy, k.items)

	return &KillRing{
		items:   itemsCopy,
		maxSize: k.maxSize,
		index:   k.index,
	}
}
