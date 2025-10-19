package value

import "fmt"

// MIMEType represents a MIME type for clipboard content
type MIMEType string

const (
	// MIMETypePlainText represents plain text content
	MIMETypePlainText MIMEType = "text/plain"

	// MIMETypeHTML represents HTML content
	MIMETypeHTML MIMEType = "text/html"

	// MIMETypeRTF represents rich text format
	MIMETypeRTF MIMEType = "text/rtf"

	// MIMETypeImage represents image content
	MIMETypeImage MIMEType = "image/png"

	// MIMETypeBinary represents binary content
	MIMETypeBinary MIMEType = "application/octet-stream"
)

// NewMIMEType creates a new MIME type value object
func NewMIMEType(mimeType string) (MIMEType, error) {
	if mimeType == "" {
		return "", fmt.Errorf("MIME type cannot be empty")
	}
	return MIMEType(mimeType), nil
}

// String returns the string representation of the MIME type
func (m MIMEType) String() string {
	return string(m)
}

// IsText returns true if the MIME type represents text content
func (m MIMEType) IsText() bool {
	switch m {
	case MIMETypePlainText, MIMETypeHTML, MIMETypeRTF:
		return true
	default:
		return false
	}
}

// IsBinary returns true if the MIME type represents binary content
func (m MIMEType) IsBinary() bool {
	return m == MIMETypeBinary || m == MIMETypeImage
}

// Equals compares two MIME types for equality
func (m MIMEType) Equals(other MIMEType) bool {
	return m == other
}
