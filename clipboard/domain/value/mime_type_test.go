package value

import (
	"testing"
)

func TestNewMIMEType(t *testing.T) {
	tests := []struct {
		name      string
		mimeType  string
		wantError bool
	}{
		{"valid plain text", "text/plain", false},
		{"valid html", "text/html", false},
		{"valid custom", "application/json", false},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewMIMEType(tt.mimeType)

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

			if result.String() != tt.mimeType {
				t.Errorf("expected %s, got %s", tt.mimeType, result.String())
			}
		})
	}
}

func TestMIMEType_IsText(t *testing.T) {
	tests := []struct {
		name     string
		mimeType MIMEType
		want     bool
	}{
		{"plain text", MIMETypePlainText, true},
		{"html", MIMETypeHTML, true},
		{"rtf", MIMETypeRTF, true},
		{"binary", MIMETypeBinary, false},
		{"image", MIMETypeImage, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mimeType.IsText(); got != tt.want {
				t.Errorf("IsText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMIMEType_IsBinary(t *testing.T) {
	tests := []struct {
		name     string
		mimeType MIMEType
		want     bool
	}{
		{"plain text", MIMETypePlainText, false},
		{"binary", MIMETypeBinary, true},
		{"image", MIMETypeImage, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mimeType.IsBinary(); got != tt.want {
				t.Errorf("IsBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMIMEType_Equals(t *testing.T) {
	mime1 := MIMETypePlainText
	mime2 := MIMETypePlainText
	mime3 := MIMETypeHTML

	if !mime1.Equals(mime2) {
		t.Errorf("expected equal MIME types to be equal")
	}

	if mime1.Equals(mime3) {
		t.Errorf("expected different MIME types to not be equal")
	}
}
