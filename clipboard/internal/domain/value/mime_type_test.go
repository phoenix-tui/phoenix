package value

import (
	"testing"
)

func TestMIMEType_IsImage(t *testing.T) {
	tests := []struct {
		name     string
		mimeType MIMEType
		want     bool
	}{
		{"PNG image", MIMETypeImagePNG, true},
		{"JPEG image", MIMETypeImageJPEG, true},
		{"GIF image", MIMETypeImageGIF, true},
		{"BMP image", MIMETypeImageBMP, true},
		{"Custom image type", MIMEType("image/webp"), true},
		{"Plain text", MIMETypePlainText, false},
		{"HTML", MIMETypeHTML, false},
		{"Binary", MIMETypeBinary, false},
		{"Deprecated Image constant", MIMETypeImage, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mimeType.IsImage(); got != tt.want {
				t.Errorf("MIMEType.IsImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMIMEType_IsBinary_WithImages(t *testing.T) {
	tests := []struct {
		name     string
		mimeType MIMEType
		want     bool
	}{
		{"PNG is binary", MIMETypeImagePNG, true},
		{"JPEG is binary", MIMETypeImageJPEG, true},
		{"GIF is binary", MIMETypeImageGIF, true},
		{"BMP is binary", MIMETypeImageBMP, true},
		{"Binary type", MIMETypeBinary, true},
		{"Plain text is not binary", MIMETypePlainText, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mimeType.IsBinary(); got != tt.want {
				t.Errorf("MIMEType.IsBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMIMEType_BackwardsCompatibility(t *testing.T) {
	// Ensure MIMETypeImage still works (deprecated but compatible)
	if MIMETypeImage != MIMETypeImagePNG {
		t.Errorf("MIMETypeImage should equal MIMETypeImagePNG for backwards compatibility")
	}

	if !MIMETypeImage.IsImage() {
		t.Error("MIMETypeImage.IsImage() should return true")
	}

	if !MIMETypeImage.IsBinary() {
		t.Error("MIMETypeImage.IsBinary() should return true")
	}
}
