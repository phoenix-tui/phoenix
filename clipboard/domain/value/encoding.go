package value

import "fmt"

// Encoding represents the encoding of clipboard content
type Encoding string

const (
	// EncodingUTF8 represents UTF-8 encoding
	EncodingUTF8 Encoding = "utf-8"

	// EncodingUTF16 represents UTF-16 encoding (Windows)
	EncodingUTF16 Encoding = "utf-16"

	// EncodingBase64 represents Base64 encoding (OSC 52)
	EncodingBase64 Encoding = "base64"

	// EncodingBinary represents raw binary encoding
	EncodingBinary Encoding = "binary"
)

// NewEncoding creates a new encoding value object
func NewEncoding(encoding string) (Encoding, error) {
	if encoding == "" {
		return "", fmt.Errorf("encoding cannot be empty")
	}
	return Encoding(encoding), nil
}

// String returns the string representation of the encoding
func (e Encoding) String() string {
	return string(e)
}

// IsTextEncoding returns true if the encoding is for text content
func (e Encoding) IsTextEncoding() bool {
	switch e {
	case EncodingUTF8, EncodingUTF16:
		return true
	default:
		return false
	}
}

// IsBinaryEncoding returns true if the encoding is for binary content
func (e Encoding) IsBinaryEncoding() bool {
	return e == EncodingBinary
}

// NeedsDecoding returns true if the content needs to be decoded before use
func (e Encoding) NeedsDecoding() bool {
	return e == EncodingBase64
}

// Equals compares two encodings for equality
func (e Encoding) Equals(other Encoding) bool {
	return e == other
}
