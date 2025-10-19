package model

import (
	"testing"

	"github.com/phoenix-tui/phoenix/clipboard/domain/value"
)

func TestNewClipboardContent(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		mimeType  value.MIMEType
		encoding  value.Encoding
		wantError bool
	}{
		{
			"valid content",
			[]byte("test data"),
			value.MIMETypePlainText,
			value.EncodingUTF8,
			false,
		},
		{
			"empty data",
			[]byte{},
			value.MIMETypePlainText,
			value.EncodingUTF8,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewClipboardContent(tt.data, tt.mimeType, tt.encoding)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("expected non-nil result")
				return
			}

			if result.MIMEType() != tt.mimeType {
				t.Errorf("expected MIME type %s, got %s", tt.mimeType, result.MIMEType())
			}

			if result.Encoding() != tt.encoding {
				t.Errorf("expected encoding %s, got %s", tt.encoding, result.Encoding())
			}
		})
	}
}

func TestNewTextContent(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		wantError bool
	}{
		{"valid text", "Hello, World!", false},
		{"empty text", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewTextContent(tt.text)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.MIMEType() != value.MIMETypePlainText {
				t.Errorf("expected plain text MIME type")
			}

			if result.Encoding() != value.EncodingUTF8 {
				t.Errorf("expected UTF-8 encoding")
			}

			text, err := result.Text()
			if err != nil {
				t.Errorf("unexpected error getting text: %v", err)
			}

			if text != tt.text {
				t.Errorf("expected text %s, got %s", tt.text, text)
			}
		})
	}
}

func TestNewBinaryContent(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		wantError bool
	}{
		{"valid binary", []byte{0x00, 0x01, 0x02}, false},
		{"empty binary", []byte{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewBinaryContent(tt.data)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.MIMEType() != value.MIMETypeBinary {
				t.Errorf("expected binary MIME type")
			}

			if result.Encoding() != value.EncodingBinary {
				t.Errorf("expected binary encoding")
			}
		})
	}
}

func TestClipboardContent_Data(t *testing.T) {
	original := []byte("test data")
	content, err := NewTextContent(string(original))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data := content.Data()

	// Verify data is correct
	if string(data) != string(original) {
		t.Errorf("expected data %s, got %s", original, data)
	}

	// Verify immutability (data is copied)
	data[0] = 'X'
	data2 := content.Data()
	if data2[0] == 'X' {
		t.Errorf("data should be immutable")
	}
}

func TestClipboardContent_Text(t *testing.T) {
	tests := []struct {
		name      string
		content   *ClipboardContent
		wantError bool
		wantText  string
	}{
		{
			"text content",
			mustNewTextContent("Hello"),
			false,
			"Hello",
		},
		{
			"binary content",
			mustNewBinaryContent([]byte{0x00}),
			true,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := tt.content.Text()

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if text != tt.wantText {
				t.Errorf("expected text %s, got %s", tt.wantText, text)
			}
		})
	}
}

func TestClipboardContent_Size(t *testing.T) {
	content := mustNewTextContent("Hello")

	if content.Size() != 5 {
		t.Errorf("expected size 5, got %d", content.Size())
	}
}

func TestClipboardContent_IsEmpty(t *testing.T) {
	content := mustNewTextContent("Hello")

	if content.IsEmpty() {
		t.Errorf("expected non-empty content")
	}
}

func TestClipboardContent_IsText(t *testing.T) {
	textContent := mustNewTextContent("Hello")
	binaryContent := mustNewBinaryContent([]byte{0x00})

	if !textContent.IsText() {
		t.Errorf("expected text content to be text")
	}

	if binaryContent.IsText() {
		t.Errorf("expected binary content to not be text")
	}
}

func TestClipboardContent_IsBinary(t *testing.T) {
	textContent := mustNewTextContent("Hello")
	binaryContent := mustNewBinaryContent([]byte{0x00})

	if textContent.IsBinary() {
		t.Errorf("expected text content to not be binary")
	}

	if !binaryContent.IsBinary() {
		t.Errorf("expected binary content to be binary")
	}
}

func TestClipboardContent_WithEncoding(t *testing.T) {
	original := mustNewTextContent("Hello")
	modified := original.WithEncoding(value.EncodingBase64)

	if modified.Encoding() != value.EncodingBase64 {
		t.Errorf("expected base64 encoding")
	}

	// Verify immutability
	if original.Encoding() != value.EncodingUTF8 {
		t.Errorf("original should not be modified")
	}
}

func TestClipboardContent_WithMIMEType(t *testing.T) {
	original := mustNewTextContent("Hello")
	modified := original.WithMIMEType(value.MIMETypeHTML)

	if modified.MIMEType() != value.MIMETypeHTML {
		t.Errorf("expected HTML MIME type")
	}

	// Verify immutability
	if original.MIMEType() != value.MIMETypePlainText {
		t.Errorf("original should not be modified")
	}
}

// Helper functions

func mustNewTextContent(text string) *ClipboardContent {
	content, err := NewTextContent(text)
	if err != nil {
		panic(err)
	}
	return content
}

func mustNewBinaryContent(data []byte) *ClipboardContent {
	content, err := NewBinaryContent(data)
	if err != nil {
		panic(err)
	}
	return content
}
