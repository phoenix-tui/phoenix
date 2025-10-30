package service

import (
	"sync"
	"testing"
	"time"

	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClipboardHistory(t *testing.T) {
	t.Run("creates history with valid parameters", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		assert.NotNil(t, history)
		assert.Equal(t, 100, history.MaxSize())
		assert.Equal(t, 24*time.Hour, history.MaxAge())
		assert.True(t, history.IsEmpty())
	})

	t.Run("treats negative max size as unlimited", func(t *testing.T) {
		history := NewClipboardHistory(-1, 0)

		assert.Equal(t, 0, history.MaxSize())
	})

	t.Run("treats negative max age as no expiration", func(t *testing.T) {
		history := NewClipboardHistory(0, -1*time.Hour)

		assert.Equal(t, time.Duration(0), history.MaxAge())
	})

	t.Run("creates empty history", func(t *testing.T) {
		history := NewClipboardHistory(0, 0)

		assert.True(t, history.IsEmpty())
		assert.Equal(t, 0, history.Size())
	})
}

func TestClipboardHistory_Add(t *testing.T) {
	t.Run("adds entry successfully", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("test content"), value.MIMETypePlainText)

		require.NoError(t, err)
		assert.Equal(t, 1, history.Size())
	})

	t.Run("adds multiple entries", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content 1"), value.MIMETypePlainText)
		require.NoError(t, err)

		err = history.Add([]byte("content 2"), value.MIMETypeHTML)
		require.NoError(t, err)

		err = history.Add([]byte("content 3"), value.MIMETypeRTF)
		require.NoError(t, err)

		assert.Equal(t, 3, history.Size())
	})

	t.Run("rejects empty content", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte{}, value.MIMETypePlainText)

		assert.Error(t, err)
		assert.Equal(t, 0, history.Size())
	})

	t.Run("rejects empty MIME type", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content"), "")

		assert.Error(t, err)
		assert.Equal(t, 0, history.Size())
	})

	t.Run("enforces max size with FIFO eviction", func(t *testing.T) {
		history := NewClipboardHistory(3, 24*time.Hour)

		// Add 4 entries
		err := history.Add([]byte("content 1"), value.MIMETypePlainText)
		require.NoError(t, err)
		entry1ID := history.GetAll()[0].ID()

		err = history.Add([]byte("content 2"), value.MIMETypePlainText)
		require.NoError(t, err)

		err = history.Add([]byte("content 3"), value.MIMETypePlainText)
		require.NoError(t, err)

		err = history.Add([]byte("content 4"), value.MIMETypePlainText)
		require.NoError(t, err)

		// Should only have 3 entries (max size)
		assert.Equal(t, 3, history.Size())

		// First entry should be evicted
		assert.False(t, history.Contains(entry1ID))
	})

	t.Run("unlimited size when max size is 0", func(t *testing.T) {
		history := NewClipboardHistory(0, 24*time.Hour)

		// Add many entries
		for i := 0; i < 100; i++ {
			err := history.Add([]byte("content"), value.MIMETypePlainText)
			require.NoError(t, err)
		}

		assert.Equal(t, 100, history.Size())
	})
}

func TestClipboardHistory_Get(t *testing.T) {
	t.Run("gets entry by ID", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("test content"), value.MIMETypePlainText)
		require.NoError(t, err)

		entries := history.GetAll()
		require.Len(t, entries, 1)

		entry, err := history.Get(entries[0].ID())

		require.NoError(t, err)
		assert.Equal(t, entries[0].ID(), entry.ID())
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		entry, err := history.Get("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, entry)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("finds correct entry among multiple", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content 1"), value.MIMETypePlainText)
		require.NoError(t, err)

		time.Sleep(2 * time.Millisecond) // Ensure different timestamps for deterministic sort
		err = history.Add([]byte("content 2"), value.MIMETypeHTML)
		require.NoError(t, err)

		entries := history.GetAll()
		// GetAll() returns entries sorted by newest first: [0]=HTML (newest), [1]=PlainText (oldest)
		targetID := entries[0].ID()

		entry, err := history.Get(targetID)

		require.NoError(t, err)
		assert.Equal(t, targetID, entry.ID())
		assert.Equal(t, value.MIMETypeHTML, entry.MIMEType())
	})
}

func TestClipboardHistory_GetAll(t *testing.T) {
	t.Run("returns empty slice for empty history", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		entries := history.GetAll()

		assert.Empty(t, entries)
	})

	t.Run("returns all entries sorted by newest first", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content 1"), value.MIMETypePlainText)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		err = history.Add([]byte("content 2"), value.MIMETypePlainText)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		err = history.Add([]byte("content 3"), value.MIMETypePlainText)
		require.NoError(t, err)

		entries := history.GetAll()

		require.Len(t, entries, 3)
		// Newest should be first
		text, _ := entries[0].Text()
		assert.Equal(t, "content 3", text)
		text, _ = entries[2].Text()
		assert.Equal(t, "content 1", text)
	})

	t.Run("returns copy not affecting original", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content"), value.MIMETypePlainText)
		require.NoError(t, err)

		entries := history.GetAll()
		entries[0] = nil // Modify copy

		// Original should be unchanged
		assert.Equal(t, 1, history.Size())
	})
}

func TestClipboardHistory_GetRecent(t *testing.T) {
	t.Run("returns recent entries", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		for i := 1; i <= 5; i++ {
			err := history.Add([]byte("content"), value.MIMETypePlainText)
			require.NoError(t, err)
			time.Sleep(5 * time.Millisecond)
		}

		recent := history.GetRecent(3)

		assert.Len(t, recent, 3)
	})

	t.Run("returns all entries when count exceeds size", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content 1"), value.MIMETypePlainText)
		require.NoError(t, err)

		err = history.Add([]byte("content 2"), value.MIMETypePlainText)
		require.NoError(t, err)

		recent := history.GetRecent(10)

		assert.Len(t, recent, 2)
	})

	t.Run("returns all entries when count is 0", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content 1"), value.MIMETypePlainText)
		require.NoError(t, err)

		err = history.Add([]byte("content 2"), value.MIMETypePlainText)
		require.NoError(t, err)

		recent := history.GetRecent(0)

		assert.Len(t, recent, 2)
	})

	t.Run("returns all entries when count is negative", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content 1"), value.MIMETypePlainText)
		require.NoError(t, err)

		err = history.Add([]byte("content 2"), value.MIMETypePlainText)
		require.NoError(t, err)

		recent := history.GetRecent(-5)

		assert.Len(t, recent, 2)
	})

	t.Run("returns entries sorted by newest first", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("oldest"), value.MIMETypePlainText)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		err = history.Add([]byte("newest"), value.MIMETypePlainText)
		require.NoError(t, err)

		recent := history.GetRecent(2)

		text, _ := recent[0].Text()
		assert.Equal(t, "newest", text)
	})
}

func TestClipboardHistory_GetByMIMEType(t *testing.T) {
	t.Run("returns entries matching MIME type", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("text 1"), value.MIMETypePlainText)
		require.NoError(t, err)

		err = history.Add([]byte("<p>html</p>"), value.MIMETypeHTML)
		require.NoError(t, err)

		err = history.Add([]byte("text 2"), value.MIMETypePlainText)
		require.NoError(t, err)

		textEntries := history.GetByMIMEType(value.MIMETypePlainText)

		assert.Len(t, textEntries, 2)
		for _, entry := range textEntries {
			assert.Equal(t, value.MIMETypePlainText, entry.MIMEType())
		}
	})

	t.Run("returns empty slice when no matches", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("text"), value.MIMETypePlainText)
		require.NoError(t, err)

		entries := history.GetByMIMEType(value.MIMETypeImagePNG)

		assert.Empty(t, entries)
	})

	t.Run("returns entries sorted by newest first", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("text 1"), value.MIMETypePlainText)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		err = history.Add([]byte("text 2"), value.MIMETypePlainText)
		require.NoError(t, err)

		entries := history.GetByMIMEType(value.MIMETypePlainText)

		require.Len(t, entries, 2)
		text, _ := entries[0].Text()
		assert.Equal(t, "text 2", text)
	})
}

func TestClipboardHistory_Clear(t *testing.T) {
	t.Run("removes all entries", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content 1"), value.MIMETypePlainText)
		require.NoError(t, err)

		err = history.Add([]byte("content 2"), value.MIMETypePlainText)
		require.NoError(t, err)

		history.Clear()

		assert.True(t, history.IsEmpty())
		assert.Equal(t, 0, history.Size())
	})

	t.Run("clear on empty history does nothing", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		history.Clear()

		assert.True(t, history.IsEmpty())
	})
}

func TestClipboardHistory_RemoveExpired(t *testing.T) {
	t.Run("removes expired entries", func(t *testing.T) {
		history := NewClipboardHistory(100, 1*time.Second)

		// Add old entry
		err := history.Add([]byte("old content"), value.MIMETypePlainText)
		require.NoError(t, err)

		// Wait for expiration
		time.Sleep(1100 * time.Millisecond)

		// Add new entry
		err = history.Add([]byte("new content"), value.MIMETypePlainText)
		require.NoError(t, err)

		removed := history.RemoveExpired()

		assert.Equal(t, 1, removed)
		assert.Equal(t, 1, history.Size())
	})

	t.Run("removes nothing when max age is 0", func(t *testing.T) {
		history := NewClipboardHistory(100, 0)

		err := history.Add([]byte("content"), value.MIMETypePlainText)
		require.NoError(t, err)

		removed := history.RemoveExpired()

		assert.Equal(t, 0, removed)
		assert.Equal(t, 1, history.Size())
	})

	t.Run("removes nothing when no entries expired", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content"), value.MIMETypePlainText)
		require.NoError(t, err)

		removed := history.RemoveExpired()

		assert.Equal(t, 0, removed)
		assert.Equal(t, 1, history.Size())
	})

	t.Run("removes all expired entries", func(t *testing.T) {
		history := NewClipboardHistory(100, 500*time.Millisecond)

		err := history.Add([]byte("old 1"), value.MIMETypePlainText)
		require.NoError(t, err)

		err = history.Add([]byte("old 2"), value.MIMETypePlainText)
		require.NoError(t, err)

		// Wait for expiration
		time.Sleep(600 * time.Millisecond)

		removed := history.RemoveExpired()

		assert.Equal(t, 2, removed)
		assert.True(t, history.IsEmpty())
	})
}

func TestClipboardHistory_Size(t *testing.T) {
	t.Run("returns 0 for empty history", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		assert.Equal(t, 0, history.Size())
	})

	t.Run("returns correct size", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		for i := 0; i < 5; i++ {
			err := history.Add([]byte("content"), value.MIMETypePlainText)
			require.NoError(t, err)
		}

		assert.Equal(t, 5, history.Size())
	})

	t.Run("updates after clear", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content"), value.MIMETypePlainText)
		require.NoError(t, err)

		history.Clear()

		assert.Equal(t, 0, history.Size())
	})
}

func TestClipboardHistory_TotalSize(t *testing.T) {
	t.Run("returns 0 for empty history", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		assert.Equal(t, 0, history.TotalSize())
	})

	t.Run("returns sum of all entry sizes", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("1234"), value.MIMETypePlainText) // 4 bytes
		require.NoError(t, err)

		err = history.Add([]byte("123456"), value.MIMETypePlainText) // 6 bytes
		require.NoError(t, err)

		assert.Equal(t, 10, history.TotalSize())
	})

	t.Run("updates after clear", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content"), value.MIMETypePlainText)
		require.NoError(t, err)

		history.Clear()

		assert.Equal(t, 0, history.TotalSize())
	})
}

func TestClipboardHistory_SetMaxSize(t *testing.T) {
	t.Run("updates max size", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		history.SetMaxSize(50)

		assert.Equal(t, 50, history.MaxSize())
	})

	t.Run("enforces new max size by removing oldest", func(t *testing.T) {
		history := NewClipboardHistory(10, 24*time.Hour)

		for i := 0; i < 5; i++ {
			err := history.Add([]byte("content"), value.MIMETypePlainText)
			require.NoError(t, err)
		}

		history.SetMaxSize(3)

		assert.Equal(t, 3, history.Size())
	})

	t.Run("treats negative as unlimited", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		history.SetMaxSize(-1)

		assert.Equal(t, 0, history.MaxSize())
	})
}

func TestClipboardHistory_SetMaxAge(t *testing.T) {
	t.Run("updates max age", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		history.SetMaxAge(12 * time.Hour)

		assert.Equal(t, 12*time.Hour, history.MaxAge())
	})

	t.Run("treats negative as no expiration", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		history.SetMaxAge(-1 * time.Hour)

		assert.Equal(t, time.Duration(0), history.MaxAge())
	})
}

func TestClipboardHistory_IsEmpty(t *testing.T) {
	t.Run("returns true for new history", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		assert.True(t, history.IsEmpty())
	})

	t.Run("returns false after adding entry", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content"), value.MIMETypePlainText)
		require.NoError(t, err)

		assert.False(t, history.IsEmpty())
	})

	t.Run("returns true after clear", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content"), value.MIMETypePlainText)
		require.NoError(t, err)

		history.Clear()

		assert.True(t, history.IsEmpty())
	})
}

func TestClipboardHistory_Contains(t *testing.T) {
	t.Run("returns true for existing entry", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content"), value.MIMETypePlainText)
		require.NoError(t, err)

		id := history.GetAll()[0].ID()

		assert.True(t, history.Contains(id))
	})

	t.Run("returns false for non-existent entry", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		assert.False(t, history.Contains("non-existent-id"))
	})

	t.Run("returns false after clear", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		err := history.Add([]byte("content"), value.MIMETypePlainText)
		require.NoError(t, err)

		id := history.GetAll()[0].ID()

		history.Clear()

		assert.False(t, history.Contains(id))
	})
}

func TestClipboardHistory_Concurrency(t *testing.T) {
	t.Run("handles concurrent adds", func(t *testing.T) {
		history := NewClipboardHistory(1000, 24*time.Hour)

		var wg sync.WaitGroup
		numGoroutines := 10
		addsPerGoroutine := 10

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < addsPerGoroutine; j++ {
					_ = history.Add([]byte("content"), value.MIMETypePlainText)
				}
			}()
		}

		wg.Wait()

		assert.Equal(t, numGoroutines*addsPerGoroutine, history.Size())
	})

	t.Run("handles concurrent reads and writes", func(t *testing.T) {
		history := NewClipboardHistory(100, 24*time.Hour)

		var wg sync.WaitGroup

		// Writer goroutines
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					_ = history.Add([]byte("content"), value.MIMETypePlainText)
					time.Sleep(1 * time.Millisecond)
				}
			}()
		}

		// Reader goroutines
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					_ = history.GetAll()
					_ = history.Size()
					time.Sleep(1 * time.Millisecond)
				}
			}()
		}

		wg.Wait()

		// Should not panic and should have some entries
		assert.Greater(t, history.Size(), 0)
	})
}
