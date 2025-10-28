package value

import (
	"testing"
)

func TestNewEncoding(t *testing.T) {
	tests := []struct {
		name      string
		encoding  string
		wantError bool
	}{
		{"valid utf-8", "utf-8", false},
		{"valid utf-16", "utf-16", false},
		{"valid base64", "base64", false},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewEncoding(tt.encoding)

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

			if result.String() != tt.encoding {
				t.Errorf("expected %s, got %s", tt.encoding, result.String())
			}
		})
	}
}

func TestEncoding_IsTextEncoding(t *testing.T) {
	tests := []struct {
		name     string
		encoding Encoding
		want     bool
	}{
		{"utf-8", EncodingUTF8, true},
		{"utf-16", EncodingUTF16, true},
		{"base64", EncodingBase64, false},
		{"binary", EncodingBinary, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.encoding.IsTextEncoding(); got != tt.want {
				t.Errorf("IsTextEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncoding_IsBinaryEncoding(t *testing.T) {
	tests := []struct {
		name     string
		encoding Encoding
		want     bool
	}{
		{"utf-8", EncodingUTF8, false},
		{"binary", EncodingBinary, true},
		{"base64", EncodingBase64, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.encoding.IsBinaryEncoding(); got != tt.want {
				t.Errorf("IsBinaryEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncoding_NeedsDecoding(t *testing.T) {
	tests := []struct {
		name     string
		encoding Encoding
		want     bool
	}{
		{"utf-8", EncodingUTF8, false},
		{"base64", EncodingBase64, true},
		{"binary", EncodingBinary, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.encoding.NeedsDecoding(); got != tt.want {
				t.Errorf("NeedsDecoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncoding_Equals(t *testing.T) {
	enc1 := EncodingUTF8
	enc2 := EncodingUTF8
	enc3 := EncodingUTF16

	if !enc1.Equals(enc2) {
		t.Errorf("expected equal encodings to be equal")
	}

	if enc1.Equals(enc3) {
		t.Errorf("expected different encodings to not be equal")
	}
}
