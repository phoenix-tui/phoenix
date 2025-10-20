// Package model contains clipboard domain models and business entities.
package model

import (
	"fmt"

	"github.com/phoenix-tui/phoenix/clipboard/domain/value"
)

// ClipboardContent represents the content stored in the clipboard.
// This is the aggregate root in DDD terms.
type ClipboardContent struct {
	data     []byte
	mimeType value.MIMEType
	encoding value.Encoding
}

// NewClipboardContent creates a new clipboard content aggregate.
func NewClipboardContent(data []byte, mimeType value.MIMEType, encoding value.Encoding) (*ClipboardContent, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("clipboard content cannot be empty")
	}

	return &ClipboardContent{
		data:     data,
		mimeType: mimeType,
		encoding: encoding,
	}, nil
}

// NewTextContent creates clipboard content from text.
func NewTextContent(text string) (*ClipboardContent, error) {
	if text == "" {
		return nil, fmt.Errorf("text content cannot be empty")
	}

	return &ClipboardContent{
		data:     []byte(text),
		mimeType: value.MIMETypePlainText,
		encoding: value.EncodingUTF8,
	}, nil
}

// NewBinaryContent creates clipboard content from binary data.
func NewBinaryContent(data []byte) (*ClipboardContent, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("binary content cannot be empty")
	}

	return &ClipboardContent{
		data:     data,
		mimeType: value.MIMETypeBinary,
		encoding: value.EncodingBinary,
	}, nil
}

// Data returns the raw data.
func (c *ClipboardContent) Data() []byte {
	// Return a copy to preserve immutability
	result := make([]byte, len(c.data))
	copy(result, c.data)
	return result
}

// MIMEType returns the MIME type.
func (c *ClipboardContent) MIMEType() value.MIMEType {
	return c.mimeType
}

// Encoding returns the encoding.
func (c *ClipboardContent) Encoding() value.Encoding {
	return c.encoding
}

// Text returns the content as text if it's text-based.
func (c *ClipboardContent) Text() (string, error) {
	if !c.mimeType.IsText() {
		return "", fmt.Errorf("content is not text (MIME type: %s)", c.mimeType)
	}

	return string(c.data), nil
}

// Size returns the size of the content in bytes.
func (c *ClipboardContent) Size() int {
	return len(c.data)
}

// IsEmpty returns true if the content is empty.
func (c *ClipboardContent) IsEmpty() bool {
	return len(c.data) == 0
}

// IsText returns true if the content is text-based.
func (c *ClipboardContent) IsText() bool {
	return c.mimeType.IsText()
}

// IsBinary returns true if the content is binary.
func (c *ClipboardContent) IsBinary() bool {
	return c.mimeType.IsBinary()
}

// WithEncoding returns a new ClipboardContent with a different encoding.
func (c *ClipboardContent) WithEncoding(encoding value.Encoding) *ClipboardContent {
	return &ClipboardContent{
		data:     c.data,
		mimeType: c.mimeType,
		encoding: encoding,
	}
}

// WithMIMEType returns a new ClipboardContent with a different MIME type.
func (c *ClipboardContent) WithMIMEType(mimeType value.MIMEType) *ClipboardContent {
	return &ClipboardContent{
		data:     c.data,
		mimeType: mimeType,
		encoding: c.encoding,
	}
}
