package service

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	_ "image/gif"  // Register GIF format
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format

	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/value"
)

// ImageCodec provides image encoding and decoding functionality.
// It's a domain service that handles image format conversions.
type ImageCodec struct{}

// NewImageCodec creates a new ImageCodec instance.
func NewImageCodec() *ImageCodec {
	return &ImageCodec{}
}

// EncodePNG encodes image data to PNG format.
// The input imageData should be raw pixel data or an already decoded image.
func (c *ImageCodec) EncodePNG(img image.Image) ([]byte, error) {
	if img == nil {
		return nil, fmt.Errorf("image cannot be nil")
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return buf.Bytes(), nil
}

// EncodeJPEG encodes image data to JPEG format with the specified quality (1-100).
// Default quality is 90 if quality <= 0 or quality > 100.
func (c *ImageCodec) EncodeJPEG(img image.Image, quality int) ([]byte, error) {
	if img == nil {
		return nil, fmt.Errorf("image cannot be nil")
	}

	// Validate and set default quality
	if quality <= 0 || quality > 100 {
		quality = 90
	}

	var buf bytes.Buffer
	opts := &jpeg.Options{Quality: quality}
	if err := jpeg.Encode(&buf, img, opts); err != nil {
		return nil, fmt.Errorf("failed to encode JPEG: %w", err)
	}

	return buf.Bytes(), nil
}

// EncodeGIF encodes image data to GIF format.
func (c *ImageCodec) EncodeGIF(img image.Image) ([]byte, error) {
	if img == nil {
		return nil, fmt.Errorf("image cannot be nil")
	}

	var buf bytes.Buffer
	opts := &gif.Options{}
	if err := gif.Encode(&buf, img, opts); err != nil {
		return nil, fmt.Errorf("failed to encode GIF: %w", err)
	}

	return buf.Bytes(), nil
}

// Decode decodes image data from any supported format (PNG, JPEG, GIF, BMP).
// Returns the decoded image and the detected format.
func (c *ImageCodec) Decode(data []byte) (image.Image, value.MIMEType, error) {
	if len(data) == 0 {
		return nil, "", fmt.Errorf("image data cannot be empty")
	}

	reader := bytes.NewReader(data)
	img, format, err := image.Decode(reader)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Convert format string to MIME type
	mimeType, err := c.formatToMIMEType(format)
	if err != nil {
		return nil, "", err
	}

	return img, mimeType, nil
}

// DecodePNG decodes PNG image data.
func (c *ImageCodec) DecodePNG(data []byte) (image.Image, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("PNG data cannot be empty")
	}

	reader := bytes.NewReader(data)
	img, err := png.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}

	return img, nil
}

// DecodeJPEG decodes JPEG image data.
func (c *ImageCodec) DecodeJPEG(data []byte) (image.Image, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("JPEG data cannot be empty")
	}

	reader := bytes.NewReader(data)
	img, err := jpeg.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JPEG: %w", err)
	}

	return img, nil
}

// DecodeGIF decodes GIF image data.
func (c *ImageCodec) DecodeGIF(data []byte) (image.Image, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("GIF data cannot be empty")
	}

	reader := bytes.NewReader(data)
	img, err := gif.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode GIF: %w", err)
	}

	return img, nil
}

// DetectFormat detects the image format from the data bytes.
// It uses magic bytes to identify the format.
func (c *ImageCodec) DetectFormat(data []byte) (value.MIMEType, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("image data cannot be empty")
	}

	// Check PNG magic bytes (89 50 4E 47 0D 0A 1A 0A)
	if len(data) >= 8 {
		if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
			return value.MIMETypeImagePNG, nil
		}
	}

	// Check JPEG magic bytes (FF D8 FF)
	if len(data) >= 3 {
		if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
			return value.MIMETypeImageJPEG, nil
		}
	}

	// Check GIF magic bytes (47 49 46 38)
	if len(data) >= 6 {
		if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x38 {
			return value.MIMETypeImageGIF, nil
		}
	}

	// Check BMP magic bytes (42 4D)
	if len(data) >= 2 {
		if data[0] == 0x42 && data[1] == 0x4D {
			return value.MIMETypeImageBMP, nil
		}
	}

	// Try using image.DecodeConfig to detect format
	reader := bytes.NewReader(data)
	_, format, err := image.DecodeConfig(reader)
	if err != nil {
		return "", fmt.Errorf("unable to detect image format: %w", err)
	}

	return c.formatToMIMEType(format)
}

// ConvertFormat converts an image from one format to another.
func (c *ImageCodec) ConvertFormat(data []byte, targetFormat value.MIMEType) ([]byte, error) {
	// Decode source image
	img, _, err := c.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode source image: %w", err)
	}

	// Encode to target format
	switch targetFormat {
	case value.MIMETypeImagePNG:
		return c.EncodePNG(img)
	case value.MIMETypeImageJPEG:
		return c.EncodeJPEG(img, 90)
	case value.MIMETypeImageGIF:
		return c.EncodeGIF(img)
	default:
		return nil, fmt.Errorf("unsupported target format: %s", targetFormat)
	}
}

// formatToMIMEType converts Go's image format string to a MIME type.
func (c *ImageCodec) formatToMIMEType(format string) (value.MIMEType, error) {
	switch format {
	case "png":
		return value.MIMETypeImagePNG, nil
	case "jpeg":
		return value.MIMETypeImageJPEG, nil
	case "gif":
		return value.MIMETypeImageGIF, nil
	case "bmp":
		return value.MIMETypeImageBMP, nil
	default:
		return "", fmt.Errorf("unsupported image format: %s", format)
	}
}
