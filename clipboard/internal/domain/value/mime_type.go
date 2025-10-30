package value

import (
	"fmt"
	"strings"
)

// MIMEType represents a MIME type for clipboard content.
type MIMEType string

const (
	// MIMETypePlainText represents plain text content.
	MIMETypePlainText MIMEType = "text/plain"

	// MIMETypeHTML represents HTML content.
	MIMETypeHTML MIMEType = "text/html"

	// MIMETypeRTF represents rich text format.
	MIMETypeRTF MIMEType = "text/rtf"

	// MIMETypeImagePNG represents PNG image content.
	MIMETypeImagePNG MIMEType = "image/png"

	// MIMETypeImageJPEG represents JPEG image content.
	MIMETypeImageJPEG MIMEType = "image/jpeg"

	// MIMETypeImageGIF represents GIF image content.
	MIMETypeImageGIF MIMEType = "image/gif"

	// MIMETypeImageBMP represents BMP image content.
	MIMETypeImageBMP MIMEType = "image/bmp"

	// MIMETypeBinary represents binary content.
	MIMETypeBinary MIMEType = "application/octet-stream"

	// MIMETypeImage is deprecated. Use MIMETypeImagePNG instead.
	// Kept for backwards compatibility.
	MIMETypeImage MIMEType = MIMETypeImagePNG
)

// NewMIMEType creates a new MIME type value object.
func NewMIMEType(mimeType string) (MIMEType, error) {
	if mimeType == "" {
		return "", fmt.Errorf("MIME type cannot be empty")
	}
	return MIMEType(mimeType), nil
}

// String returns the string representation of the MIME type.
func (m MIMEType) String() string {
	return string(m)
}

// IsText returns true if the MIME type represents text content.
func (m MIMEType) IsText() bool {
	switch m {
	case MIMETypePlainText, MIMETypeHTML, MIMETypeRTF:
		return true
	default:
		return false
	}
}

// IsImage returns true if the MIME type represents image content.
func (m MIMEType) IsImage() bool {
	switch m {
	case MIMETypeImagePNG, MIMETypeImageJPEG, MIMETypeImageGIF, MIMETypeImageBMP:
		return true
	default:
		// Also check for any image/* MIME type
		return strings.HasPrefix(string(m), "image/")
	}
}

// IsBinary returns true if the MIME type represents binary content.
func (m MIMEType) IsBinary() bool {
	return m.IsImage() || m == MIMETypeBinary
}

// Equals compares two MIME types for equality.
func (m MIMEType) Equals(other MIMEType) bool {
	return m == other
}
