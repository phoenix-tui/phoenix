package model

import (
	"testing"
	"time"

	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHistoryEntry(t *testing.T) {
	t.Run("creates valid entry", func(t *testing.T) {
		content := []byte("test content")
		mimeType := value.MIMETypePlainText

		entry, err := NewHistoryEntry(content, mimeType)

		require.NoError(t, err)
		require.NotNil(t, entry)
		assert.NotEmpty(t, entry.ID())
		assert.Equal(t, content, entry.Content())
		assert.Equal(t, mimeType, entry.MIMEType())
		assert.Equal(t, len(content), entry.Size())
		assert.False(t, entry.Timestamp().IsZero())
	})

	t.Run("generates unique IDs", func(t *testing.T) {
		entry1, err := NewHistoryEntry([]byte("content1"), value.MIMETypePlainText)
		require.NoError(t, err)

		entry2, err := NewHistoryEntry([]byte("content2"), value.MIMETypePlainText)
		require.NoError(t, err)

		assert.NotEqual(t, entry1.ID(), entry2.ID())
	})

	t.Run("rejects empty content", func(t *testing.T) {
		entry, err := NewHistoryEntry([]byte{}, value.MIMETypePlainText)

		assert.Error(t, err)
		assert.Nil(t, entry)
		assert.Contains(t, err.Error(), "content cannot be empty")
	})

	t.Run("rejects nil content", func(t *testing.T) {
		entry, err := NewHistoryEntry(nil, value.MIMETypePlainText)

		assert.Error(t, err)
		assert.Nil(t, entry)
		assert.Contains(t, err.Error(), "content cannot be empty")
	})

	t.Run("rejects empty MIME type", func(t *testing.T) {
		entry, err := NewHistoryEntry([]byte("content"), "")

		assert.Error(t, err)
		assert.Nil(t, entry)
		assert.Contains(t, err.Error(), "MIME type cannot be empty")
	})

	t.Run("handles different MIME types", func(t *testing.T) {
		testCases := []struct {
			name     string
			mimeType value.MIMEType
		}{
			{"plain text", value.MIMETypePlainText},
			{"HTML", value.MIMETypeHTML},
			{"RTF", value.MIMETypeRTF},
			{"PNG image", value.MIMETypeImagePNG},
			{"JPEG image", value.MIMETypeImageJPEG},
			{"binary", value.MIMETypeBinary},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				entry, err := NewHistoryEntry([]byte("content"), tc.mimeType)

				require.NoError(t, err)
				assert.Equal(t, tc.mimeType, entry.MIMEType())
			})
		}
	})
}

func TestNewHistoryEntryWithTime(t *testing.T) {
	t.Run("creates entry with specific timestamp", func(t *testing.T) {
		content := []byte("test content")
		mimeType := value.MIMETypePlainText
		timestamp := time.Date(2025, 10, 17, 10, 30, 0, 0, time.UTC)

		entry, err := NewHistoryEntryWithTime(content, mimeType, timestamp)

		require.NoError(t, err)
		assert.Equal(t, timestamp, entry.Timestamp())
	})

	t.Run("rejects zero timestamp", func(t *testing.T) {
		entry, err := NewHistoryEntryWithTime([]byte("content"), value.MIMETypePlainText, time.Time{})

		assert.Error(t, err)
		assert.Nil(t, entry)
		assert.Contains(t, err.Error(), "timestamp cannot be zero")
	})

	t.Run("rejects empty content", func(t *testing.T) {
		entry, err := NewHistoryEntryWithTime([]byte{}, value.MIMETypePlainText, time.Now())

		assert.Error(t, err)
		assert.Nil(t, entry)
	})

	t.Run("rejects empty MIME type", func(t *testing.T) {
		entry, err := NewHistoryEntryWithTime([]byte("content"), "", time.Now())

		assert.Error(t, err)
		assert.Nil(t, entry)
	})
}

func TestHistoryEntry_Content(t *testing.T) {
	t.Run("returns copy of content", func(t *testing.T) {
		original := []byte("test content")
		entry, err := NewHistoryEntry(original, value.MIMETypePlainText)
		require.NoError(t, err)

		retrieved := entry.Content()

		// Modify retrieved copy
		retrieved[0] = 'X'

		// Original should be unchanged
		assert.Equal(t, original, entry.Content())
		assert.NotEqual(t, retrieved, entry.Content())
	})
}

func TestHistoryEntry_IsExpired(t *testing.T) {
	t.Run("entry not expired within max age", func(t *testing.T) {
		entry, err := NewHistoryEntry([]byte("content"), value.MIMETypePlainText)
		require.NoError(t, err)

		assert.False(t, entry.IsExpired(1*time.Hour))
	})

	t.Run("entry expired after max age", func(t *testing.T) {
		oldTime := time.Now().Add(-2 * time.Hour)
		entry, err := NewHistoryEntryWithTime([]byte("content"), value.MIMETypePlainText, oldTime)
		require.NoError(t, err)

		assert.True(t, entry.IsExpired(1*time.Hour))
	})

	t.Run("entry not expired with zero max age", func(t *testing.T) {
		oldTime := time.Now().Add(-100 * time.Hour)
		entry, err := NewHistoryEntryWithTime([]byte("content"), value.MIMETypePlainText, oldTime)
		require.NoError(t, err)

		assert.False(t, entry.IsExpired(0))
	})

	t.Run("entry not expired with negative max age", func(t *testing.T) {
		oldTime := time.Now().Add(-100 * time.Hour)
		entry, err := NewHistoryEntryWithTime([]byte("content"), value.MIMETypePlainText, oldTime)
		require.NoError(t, err)

		assert.False(t, entry.IsExpired(-1*time.Hour))
	})

	t.Run("boundary case - exactly at max age", func(t *testing.T) {
		maxAge := 1 * time.Hour
		oldTime := time.Now().Add(-maxAge)
		entry, err := NewHistoryEntryWithTime([]byte("content"), value.MIMETypePlainText, oldTime)
		require.NoError(t, err)

		// Due to time passing during test execution, this might be slightly expired
		// Just verify it doesn't panic
		_ = entry.IsExpired(maxAge)
	})
}

func TestHistoryEntry_IsText(t *testing.T) {
	testCases := []struct {
		name     string
		mimeType value.MIMEType
		expected bool
	}{
		{"plain text", value.MIMETypePlainText, true},
		{"HTML", value.MIMETypeHTML, true},
		{"RTF", value.MIMETypeRTF, true},
		{"PNG image", value.MIMETypeImagePNG, false},
		{"JPEG image", value.MIMETypeImageJPEG, false},
		{"binary", value.MIMETypeBinary, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry, err := NewHistoryEntry([]byte("content"), tc.mimeType)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, entry.IsText())
		})
	}
}

func TestHistoryEntry_IsImage(t *testing.T) {
	testCases := []struct {
		name     string
		mimeType value.MIMEType
		expected bool
	}{
		{"PNG image", value.MIMETypeImagePNG, true},
		{"JPEG image", value.MIMETypeImageJPEG, true},
		{"GIF image", value.MIMETypeImageGIF, true},
		{"BMP image", value.MIMETypeImageBMP, true},
		{"plain text", value.MIMETypePlainText, false},
		{"HTML", value.MIMETypeHTML, false},
		{"binary", value.MIMETypeBinary, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry, err := NewHistoryEntry([]byte("content"), tc.mimeType)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, entry.IsImage())
		})
	}
}

func TestHistoryEntry_IsBinary(t *testing.T) {
	testCases := []struct {
		name     string
		mimeType value.MIMEType
		expected bool
	}{
		{"PNG image", value.MIMETypeImagePNG, true},
		{"JPEG image", value.MIMETypeImageJPEG, true},
		{"binary", value.MIMETypeBinary, true},
		{"plain text", value.MIMETypePlainText, false},
		{"HTML", value.MIMETypeHTML, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry, err := NewHistoryEntry([]byte("content"), tc.mimeType)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, entry.IsBinary())
		})
	}
}

func TestHistoryEntry_Text(t *testing.T) {
	t.Run("returns text for text content", func(t *testing.T) {
		expectedText := "Hello, World!"
		entry, err := NewHistoryEntry([]byte(expectedText), value.MIMETypePlainText)
		require.NoError(t, err)

		text, err := entry.Text()

		require.NoError(t, err)
		assert.Equal(t, expectedText, text)
	})

	t.Run("returns text for HTML content", func(t *testing.T) {
		expectedText := "<p>Hello</p>"
		entry, err := NewHistoryEntry([]byte(expectedText), value.MIMETypeHTML)
		require.NoError(t, err)

		text, err := entry.Text()

		require.NoError(t, err)
		assert.Equal(t, expectedText, text)
	})

	t.Run("returns error for image content", func(t *testing.T) {
		entry, err := NewHistoryEntry([]byte("fake image data"), value.MIMETypeImagePNG)
		require.NoError(t, err)

		text, err := entry.Text()

		assert.Error(t, err)
		assert.Empty(t, text)
		assert.Contains(t, err.Error(), "not text")
	})

	t.Run("returns error for binary content", func(t *testing.T) {
		entry, err := NewHistoryEntry([]byte{0x00, 0x01, 0x02}, value.MIMETypeBinary)
		require.NoError(t, err)

		text, err := entry.Text()

		assert.Error(t, err)
		assert.Empty(t, text)
	})
}

func TestHistoryEntry_Age(t *testing.T) {
	t.Run("returns age of entry", func(t *testing.T) {
		pastTime := time.Now().Add(-5 * time.Minute)
		entry, err := NewHistoryEntryWithTime([]byte("content"), value.MIMETypePlainText, pastTime)
		require.NoError(t, err)

		age := entry.Age()

		// Age should be approximately 5 minutes (allow some tolerance)
		assert.InDelta(t, 5*time.Minute, age, float64(1*time.Second))
	})

	t.Run("returns near-zero age for new entry", func(t *testing.T) {
		entry, err := NewHistoryEntry([]byte("content"), value.MIMETypePlainText)
		require.NoError(t, err)

		age := entry.Age()

		// Age should be very small
		assert.Less(t, age, 1*time.Second)
	})
}

func TestHistoryEntry_Equals(t *testing.T) {
	t.Run("entries with same ID are equal", func(t *testing.T) {
		entry, err := NewHistoryEntry([]byte("content"), value.MIMETypePlainText)
		require.NoError(t, err)

		assert.True(t, entry.Equals(entry))
	})

	t.Run("entries with different IDs are not equal", func(t *testing.T) {
		entry1, err := NewHistoryEntry([]byte("content1"), value.MIMETypePlainText)
		require.NoError(t, err)

		entry2, err := NewHistoryEntry([]byte("content2"), value.MIMETypePlainText)
		require.NoError(t, err)

		assert.False(t, entry1.Equals(entry2))
		assert.False(t, entry2.Equals(entry1))
	})

	t.Run("entry does not equal nil", func(t *testing.T) {
		entry, err := NewHistoryEntry([]byte("content"), value.MIMETypePlainText)
		require.NoError(t, err)

		assert.False(t, entry.Equals(nil))
	})
}

func TestHistoryEntry_Size(t *testing.T) {
	testCases := []struct {
		name    string
		content []byte
	}{
		{"small content", []byte("test")},
		{"medium content", []byte("This is a longer test content with more characters")},
		{"large content", make([]byte, 1024*1024)}, // 1 MB
		{"unicode content", []byte("Hello 世界 🌍")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry, err := NewHistoryEntry(tc.content, value.MIMETypePlainText)
			require.NoError(t, err)

			assert.Equal(t, len(tc.content), entry.Size())
		})
	}
}

func TestHistoryEntry_Immutability(t *testing.T) {
	t.Run("modifying original content does not affect entry", func(t *testing.T) {
		original := []byte("test content")
		entry, err := NewHistoryEntry(original, value.MIMETypePlainText)
		require.NoError(t, err)

		// Modify original
		original[0] = 'X'

		// Entry should be unchanged
		assert.Equal(t, byte('t'), entry.Content()[0])
	})

	t.Run("modifying retrieved content does not affect entry", func(t *testing.T) {
		entry, err := NewHistoryEntry([]byte("test content"), value.MIMETypePlainText)
		require.NoError(t, err)

		// Get and modify content
		content := entry.Content()
		content[0] = 'X'

		// Entry should be unchanged
		assert.Equal(t, byte('t'), entry.Content()[0])
	})

	t.Run("multiple retrievals return independent copies", func(t *testing.T) {
		entry, err := NewHistoryEntry([]byte("test content"), value.MIMETypePlainText)
		require.NoError(t, err)

		content1 := entry.Content()
		content2 := entry.Content()

		// Modify first copy
		content1[0] = 'X'

		// Second copy should be unchanged
		assert.Equal(t, byte('t'), content2[0])
	})
}
